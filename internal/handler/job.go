package handler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
	"strings"

	"tonfy_CMS/internal/model"      
	"tonfy_CMS/internal/repository" 

	"github.com/labstack/echo/v4"
)

// ========= 核心接口：接收简历 =============

func SubmitApplication(c echo.Context) error {
	// 1. 提取表单文本数据
	position := c.FormValue("position")
	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	content := c.FormValue("content")

	// 2. 处理简历文件上传
	file, err := c.FormFile("resume")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "必须上传简历文件"})
	}

	// 自动创建存放简历的目录 (如果不存在的话)
	uploadDir := "uploads/resumes"
	os.MkdirAll(uploadDir, os.ModePerm)

	// 为了防止同名文件覆盖，给文件名加个时间戳前缀
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "文件读取失败"})
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "文件保存失败"})
	}
	defer dst.Close()
	io.Copy(dst, src)

	// 3. 数据写入 SQLite
	application := model.JobApplication{
		Position:      position,
		ApplicantName: name,
		Email:         email,
		Phone:         phone,
		CoverLetter:   content,
		ResumeFileUrl: "/" + savePath, // 例如: /uploads/resumes/123456_cv.pdf
	}

	if err := repository.DB.Create(&application).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "数据保存失败"})
	}

	// 启动 Goroutine 异步发送邮件
	// 所以我们把提取出来的纯文本字符串传进去。
	go sendEmailNotificationAsync(position, name, email)

	// 主线程瞬间返回
	return c.JSON(http.StatusOK, map[string]string{"message": "投递成功"})
}

// ============= 协程专用的发邮件函数 ===============
func sendEmailNotificationAsync(position, candidateName, candidateEmail string) {
	// 配置你的 SMTP 信息
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_USER")
	senderPass := os.Getenv("SMTP_PASS")
	receiverEmail := os.Getenv("HR_EMAIL")

	// 1. 组装支持中文的邮件头和正文
	subject := "【TONFY 招聘系统】收到一份新简历！"
	body := fmt.Sprintf("HR 您好：\r\n\r\n您在官网收到了一份新的职位申请。\r\n\r\n应聘岗位：%s\r\n候选人姓名：%s\r\n联系邮箱：%s\r\n\r\n请尽快登录后台管理系统下载并查看简历！", position, candidateName, candidateEmail)

	header := make(map[string]string)
	header["From"] = "TONFY CMS <" + senderEmail + ">"
	header["To"] = receiverEmail
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8" // 解决中文乱码

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 使用 TLS 拨号，强行打通 465 加密端口
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		fmt.Printf("【异步任务-报错】TLS 连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("【异步任务-报错】创建 SMTP 客户端失败: %v\n", err)
		return
	}
	defer client.Quit()

	// 3. 进行账号密码认证
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Printf("【异步任务-报错】SMTP 账号认证失败: %v\n", err)
		return
	}

	// 4. 设置发件人和收件人
	if err = client.Mail(senderEmail); err != nil {
		fmt.Printf("【异步任务-报错】设置发件人失败: %v\n", err)
		return
	}
	if err = client.Rcpt(receiverEmail); err != nil {
		fmt.Printf("【异步任务-报错】设置收件人失败: %v\n", err)
		return
	}

	// 5. 写入并发送内容
	w, err := client.Data()
	if err != nil {
		fmt.Printf("【异步任务-报错】获取 Data 写入流失败: %v\n", err)
		return
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Printf("【异步任务-报错】写入邮件内容失败: %v\n", err)
		return
	}
	err = w.Close()
	if err != nil {
		fmt.Printf("【异步任务-报错】关闭邮件流失败: %v\n", err)
		return
	}

	fmt.Printf("【异步任务-成功】已发送简历提醒邮件至: %s\n", receiverEmail)
}

// ============== 前台接口：拉取职位列表 ================

// GetFrontJobs 获取前台展示的职位 
func GetFrontJobs(c echo.Context) error {
	var jobs []model.Job
	if err := repository.DB.Where("is_active = ?", true).Order("sort_order DESC, created_at ASC").Find(&jobs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取职位失败"})
	}
	return c.JSON(http.StatusOK, jobs)
}

// =================== 后台接口：职位管理 (CRUD) ==============

// GetAdminJobs 获取后台职位列表 (包含已隐藏的)
func GetAdminJobs(c echo.Context) error {
	var jobs []model.Job
	if err := repository.DB.Order("sort_order DESC, created_at ASC").Find(&jobs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取职位失败"})
	}
	return c.JSON(http.StatusOK, jobs)
}

func CreateJob(c echo.Context) error {
	var job model.Job
	if err := c.Bind(&job); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	if err := repository.DB.Create(&job).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存失败"})
	}
	return c.JSON(http.StatusOK, job)
}

func UpdateJob(c echo.Context) error {
	id := c.Param("id")
	var job model.Job
	if err := repository.DB.First(&job, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "职位不存在"})
	}
	if err := c.Bind(&job); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	repository.DB.Save(&job)
	return c.JSON(http.StatusOK, job)
}

func DeleteJob(c echo.Context) error {
	id := c.Param("id")
	repository.DB.Delete(&model.Job{}, id)
	return c.JSON(http.StatusOK, map[string]string{"message": "删除成功"})
}

func UpdateJobsSort(c echo.Context) error {
	var payload struct {
		IDs []uint `json:"ids"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数错误"})
	}

	tx := repository.DB.Begin()
	total := len(payload.IDs)
	for index, id := range payload.IDs {
		sortOrder := total - index
		if err := tx.Model(&model.Job{}).Where("id = ?", id).Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新失败"})
		}
	}
	tx.Commit()
	return c.JSON(http.StatusOK, map[string]string{"message": "排序已保存"})
}

// ===================== 后台接口：简历管理 ==================

// GetApplications 获取所有收到的候选人简历 
func GetApplications(c echo.Context) error {
	var apps []model.JobApplication
	if err := repository.DB.Order("created_at DESC").Find(&apps).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取简历失败"})
	}
	return c.JSON(http.StatusOK, apps)
}

// UpdateApplicationStatus 修改简历状态
func UpdateApplicationStatus(c echo.Context) error {
	id := c.Param("id")
	var payload struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数错误"})
	}
	if err := repository.DB.Model(&model.JobApplication{}).Where("id = ?", id).Update("status", payload.Status).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "状态更新失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "状态已更新"})
}

// DeleteApplication 删除简历
func DeleteApplication(c echo.Context) error {
	id := c.Param("id")
	
	// 获取它的文件路径
	var app model.JobApplication
	if err := repository.DB.First(&app, id).Error; err == nil {
		// 物理删除文件。app.ResumeFileUrl 是 "/uploads/resumes/xxx.pdf"
		// 用 strings.TrimPrefix 去掉开头的 "/"，变成相对路径才能删
		filePath := strings.TrimPrefix(app.ResumeFileUrl, "/")
		os.Remove(filePath) // 尝试删除文件，就算文件不存在报错了也不影响数据库删除
	}

	//  删除数据库记录
	repository.DB.Delete(&model.JobApplication{}, id)
	return c.JSON(http.StatusOK, map[string]string{"message": "简历及文件已彻底删除"})
}