package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/guish/university-social/internal/config"
	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type AuthHandler struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewAuthHandler(repo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{repo: repo, config: cfg}
}

func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少登录参数"})
		return
	}

	session, err := h.code2Session(req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "微信登录失败"})
		return
	}

	user, err := h.repo.FindByOpenID(session.OpenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	if user == nil {
		id, err := h.repo.Create(&model.User{
			OpenID:   session.OpenID,
			UnionID:  session.UnionID,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
			return
		}
		user, _ = h.repo.FindByID(id)
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) DevLogin(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	nickname := req.Nickname
	if nickname == "" {
		nickname = "测试用户"
	}
	openID := "dev_" + nickname
	user, err := h.repo.FindByOpenID(openID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	if user == nil {
		id, err := h.repo.Create(&model.User{
			OpenID:   openID,
			Nickname: nickname,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
			return
		}
		user, _ = h.repo.FindByID(id)
	}
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}
	c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	user, err := h.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		School   string `json:"school"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.UpdateProfile(userID, req.Nickname, req.Avatar, req.School); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (h *AuthHandler) code2Session(code string) (*WechatSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		h.config.WechatAppID, h.config.WechatSecret, code,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("weixin api request create: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("weixin api request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var session WechatSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("weixin api response parse: %w", err)
	}
	if session.ErrCode != 0 {
		return nil, fmt.Errorf("weixin api error: %s", session.ErrMsg)
	}
	return &session, nil
}

func (h *AuthHandler) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Add(-time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.config.JWTSecret))
}
