package handler

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// SendProductManual 接收前端传来的路径，直接发邮件
func SendProductManual(c echo.Context) error {
	recipient := c.FormValue("recipient")
	attachmentPath := c.FormValue("attachment")

	if recipient == "" || attachmentPath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "邮箱或文件路径不能为空"})
	}

	// 1. 把前端转码的路径（%20, %E4等）还原成真正的中文或空格！
	decodedPath, err := url.PathUnescape(attachmentPath)
	if err == nil {
		attachmentPath = decodedPath
	}

	// 2. 架构级安全防线：防止越权读取服务器底层文件
	if strings.Contains(attachmentPath, "..") {
		fmt.Printf("【安全拦截】检测到非法文件读取尝试: %s\n", attachmentPath)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "非法的文件请求"})
	}

	// 去掉路径开头的 "/"
	absPath := strings.TrimPrefix(attachmentPath, "/")
	
	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("【文件报错】前台传来的路径在服务器上找不到: %s\n", absPath)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "服务器上未找到该说明书文件"})
	}

	// 发送邮件
	go sendAttachmentEmailAsync(recipient, absPath)

	return c.JSON(http.StatusOK, map[string]string{"message": "邮件发送任务已成功提交"})
}

// sendAttachmentEmailAsync 专门负责打包附件并发送邮件的底层协程
func sendAttachmentEmailAsync(targetEmail, manualPath string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_USER")
	senderPass := os.Getenv("SMTP_PASS")

	fileData, err := os.ReadFile(manualPath)
	if err != nil { 
		fmt.Printf("【邮件报错】读取附件文件失败: %v\n", err)
		return 
	}
	fileName := filepath.Base(manualPath)

	// 将 Base64 巨型字符串，强制按 76 个字符进行切片换行！
	b64Raw := base64.StdEncoding.EncodeToString(fileData)
	var chunkedB64 strings.Builder
	for i := 0; i < len(b64Raw); i += 76 {
		end := i + 76
		if end > len(b64Raw) {
			end = len(b64Raw)
		}
		chunkedB64.WriteString(b64Raw[i:end] + "\r\n")
	}

	boundary := "my-boundary-12345"
	subject := "Your Requested Technical Manual - TONFY"
	
	header := fmt.Sprintf("From: TONFY CMS <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/mixed; boundary=%s\r\n"+
		"\r\n", senderEmail, targetEmail, subject, boundary)

	body := fmt.Sprintf("--%s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"Dear Customer,\r\n\r\n"+
		"Thank you for your interest in TONFY products. Please find the technical manual you requested attached to this email.\r\n\r\n"+
		"If you have any further questions, please feel free to contact us.  info@tonfy.com\r\n\r\n"+
		"Best Regards,\r\n"+
		"TONFY Team\r\n"+
		"\r\n", boundary)

	attachment := fmt.Sprintf("--%s\r\n"+
		"Content-Type: application/octet-stream\r\n"+
		"Content-Transfer-Encoding: base64\r\n"+
		"Content-Disposition: attachment; filename=\"%s\"\r\n"+
		"\r\n%s\r\n"+
		"--%s--", boundary, fileName, chunkedB64.String(), boundary)

	fullMessage := header + body + attachment

	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: smtpHost}
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		fmt.Printf("【邮件报错】TLS连接失败: %v\n", err)
		return 
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("【邮件报错】创建SMTP客户端失败: %v\n", err)
		return 
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("【邮件报错】SMTP账号认证失败: %v\n", err)
		return 
	}
	if err := client.Mail(senderEmail); err != nil {
		fmt.Printf("【邮件报错】设置发件人失败: %v\n", err)
		return 
	}
	if err := client.Rcpt(targetEmail); err != nil {
		fmt.Printf("【邮件报错】设置收件人失败: %v\n", err)
		return 
	}

	w, err := client.Data()
	if err != nil {
		fmt.Printf("【邮件报错】开启写入流失败: %v\n", err)
		return 
	}
	if _, err := w.Write([]byte(fullMessage)); err != nil {
		fmt.Printf("【邮件报错】写入邮件内容失败: %v\n", err)
		return 
	}
	if err := w.Close(); err != nil {
		fmt.Printf("【邮件报错】关闭邮件流失败: %v\n", err)
		return 
	}
	
	fmt.Printf("已将产品资料发至客户: %s\n", targetEmail)
}