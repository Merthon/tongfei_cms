package handler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

func SubmitContact(c echo.Context) error {
	// 1. 精准提取前台 HTML 的 name 属性
	company := c.FormValue("company")
	name := c.FormValue("name")
	tel := c.FormValue("tel")         // 对应 HTML 里的 name="tel"
	email := c.FormValue("email")
	content := c.FormValue("content") // 对应 HTML 里的 name="content"
	city := c.FormValue("city")
	industry := c.FormValue("industry")

	// 2. 存入 SQLite 数据库
	contact := model.ContactMessage{
		Company:  company,
		Name:     name,
		Phone:    tel, // 把 tel 存入 Phone
		Email:    email,
		City:     city,
		Industry: industry,
		Message:  content, // 把 content 存入 Message
	}

	if err := repository.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存留言失败"})
	}

	// 3. 异步发送邮件提醒
	go sendContactEmailAsync(company, name, tel, email, city, industry, content)

	return c.JSON(http.StatusOK, map[string]string{"message": "留言成功"})
}

// 邮件发送协程
func sendContactEmailAsync(company, name, tel, email, city, industry, content string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_USER")
	senderPass := os.Getenv("SMTP_PASS")
	receiverEmailStr := os.Getenv("HR_EMAIL")
	receivers := strings.Split(receiverEmailStr, ",")

	subject := "【TONFY 官网】收到新的客户询盘！"
	
	// 把客户填写的信息排版的邮件正文
	body := fmt.Sprintf("您好：\r\n\r\n官网有客户提交了新的 Contact Us 留言。\r\n\r\n"+
		"■ 客户姓名：%s\r\n"+
		"■ 公司名称：%s\r\n"+
		"■ 所在城市：%s\r\n"+
		"■ 所属行业：%s\r\n"+
		"■ 联系电话：%s\r\n"+
		"■ 联系邮箱：%s\r\n\r\n"+
		"■ 留言内容：\r\n%s\r\n\r\n请及时联系跟进！",
		name, company, city, industry, tel, email, content)

	header := make(map[string]string)
	header["From"] = "TONFY CMS <" + senderEmail + ">"
	header["To"] = receiverEmailStr
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: smtpHost}
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil { return }
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil { return }
	defer client.Quit()

	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	if err = client.Auth(auth); err != nil { return }
	if err = client.Mail(senderEmail); err != nil { return }
	for _, email := range receivers {
		email = strings.TrimSpace(email)
		if email != "" {
			client.Rcpt(email) // 群发指令
		}
	}

	w, err := client.Data()
	if err != nil { return }
	w.Write([]byte(msg))
	w.Close()
	
	fmt.Printf("【异步任务-成功】已发送包含 7 个字段的客户询盘邮件\n")
}