package handler

import (
	"fmt"
	"net/smtp"
	"time"
	"wk-go/internal/database"

	"github.com/gin-gonic/gin"
)

type resetToken struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	UID       int    `gorm:"column:uid" json:"uid"`
	Token     string `gorm:"column:token" json:"token"`
	ExpiresAt int64  `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
}

func (resetToken) TableName() string { return "qingka_wangke_reset_token" }

// 查配置
func getConfigVal(key string) string {
	var cfg struct {
		K string `gorm:"column:k"`
	}
	database.DB.Table("qingka_wangke_config").Where("v = ?", key).First(&cfg)
	return cfg.K
}

// POST /api/auth/forget-password
func ForgetPassword(c *gin.Context) {
	var req struct {
		Account string `json:"account" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	type User struct {
		UID   int    `gorm:"column:uid"`
		User  string `gorm:"column:user"`
		Email string `gorm:"column:email"`
	}
	var user User
	if err := database.DB.Table("qingka_wangke_user").Where("user = ?", req.Account).First(&user).Error; err != nil {
		Fail(c, "账号不存在")
		return
	}
	if user.Email == "" {
		Fail(c, "该账号未绑定邮箱，请联系管理员重置密码")
		return
	}

	// 生成 token
	token := randomString(32)
	rt := resetToken{
		UID:       user.UID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		CreatedAt: time.Now().Unix(),
	}
	database.DB.Create(&rt)

	// 发邮件
	smtpHost := getConfigVal("smtp_host")
	smtpPort := getConfigVal("smtp_port")
	smtpUser := getConfigVal("smtp_user")
	smtpPass := getConfigVal("smtp_pass")
	smtpFrom := getConfigVal("smtp_from")
	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		Fail(c, "邮箱服务未配置")
		return
	}
	if smtpFrom == "" {
		smtpFrom = smtpUser
	}

	resetURL := fmt.Sprintf("http://%s/#/auth/reset-password?token=%s", c.Request.Host, token)
	subject := "密码重置 - " + getConfigVal("sitename")
	body := fmt.Sprintf(`<p>您好，%s：</p><p>您正在申请重置密码，请点击以下链接（30分钟内有效）：</p><p><a href="%s">%s</a></p><p>如非本人操作请忽略此邮件。</p>`, user.User, resetURL, resetURL)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s",
		smtpFrom, user.Email, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, smtpFrom, []string{user.Email}, []byte(msg)); err != nil {
		Fail(c, "邮件发送失败: "+err.Error())
		return
	}

	OK(c, gin.H{"msg": "重置邮件已发送，请查收"})
}

// POST /api/auth/reset-password
func ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		NewPass  string `json:"new_pass" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if len(req.NewPass) < 6 {
		Fail(c, "新密码至少6位")
		return
	}

	var rt resetToken
	if err := database.DB.Where("token = ?", req.Token).First(&rt).Error; err != nil {
		Fail(c, "重置链接无效")
		return
	}
	if time.Now().Unix() > rt.ExpiresAt {
		Fail(c, "重置链接已过期")
		return
	}

	database.DB.Table("qingka_wangke_user").Where("uid = ?", rt.UID).Update("pass", req.NewPass)
	database.DB.Delete(&rt)
	OK(c, gin.H{"msg": "密码重置成功"})
}
