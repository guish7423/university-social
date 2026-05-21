package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/config"
	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type VerificationHandler struct {
	verifRepo *repository.VerificationRepository
	cfg       *config.Config
}

func NewVerificationHandler(verifRepo *repository.VerificationRepository, cfg *config.Config) *VerificationHandler {
	return &VerificationHandler{verifRepo: verifRepo, cfg: cfg}
}

func generateCode() string {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}

func (h *VerificationHandler) SendCode(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供邮箱和学校"})
		return
	}

	verified, err := h.verifRepo.CheckVerified(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	if verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已通过认证"})
		return
	}

	code := generateCode()
	if err := h.verifRepo.CreateCode(userID, req.Email, req.School, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证码失败"})
		return
	}

	if h.cfg.SMTPHost != "" {
		auth := smtp.PlainAuth("", h.cfg.SMTPUser, h.cfg.SMTPPass, h.cfg.SMTPHost)
		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: 校园社 - 邮箱验证\r\n\r\n您的验证码是: %s\r\n有效期为10分钟。", h.cfg.SMTPFrom, req.Email, code)
		addr := fmt.Sprintf("%s:%s", h.cfg.SMTPHost, h.cfg.SMTPPort)
		err := smtp.SendMail(addr, auth, h.cfg.SMTPFrom, []string{req.Email}, []byte(msg))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件发送失败"})
			return
		}
	} else {
		fmt.Printf("[DEV] Verification code for user %d: %s (email: %s)\n", userID, code, req.Email)
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

func (h *VerificationHandler) Verify(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供验证码"})
		return
	}

	v, err := h.verifRepo.VerifyCode(userID, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证失败"})
		return
	}
	if v == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期"})
		return
	}

	if err := h.verifRepo.MarkVerified(v.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记失败"})
		return
	}
	if err := h.verifRepo.MarkVerifiedUser(userID, v.School, v.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "认证成功", "verified": true})
}

func (h *VerificationHandler) Status(c *gin.Context) {
	userID := c.GetInt64("user_id")
	verified, err := h.verifRepo.CheckVerified(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"verified": verified})
}
