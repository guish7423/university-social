package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/repository"
)

type PointsHandler struct {
	repo *repository.PointsRepository
}

func NewPointsHandler(repo *repository.PointsRepository) *PointsHandler {
	return &PointsHandler{repo: repo}
}

func (h *PointsHandler) GetBalance(c *gin.Context) {
	userID := c.GetInt64("user_id")
	balance, err := h.repo.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取积分失败"})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *PointsHandler) ClaimDailyLogin(c *gin.Context) {
	userID := c.GetInt64("user_id")
	claimed, err := h.repo.ClaimDailyLogin(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "领取失败"})
		return
	}
	if !claimed {
		c.JSON(http.StatusOK, gin.H{"claimed": false, "message": "今日已领取"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"claimed": true, "points": 1, "message": "每日登录 +1 积分"})
}

func (h *PointsHandler) GetLogs(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	logs, total, err := h.repo.GetLogs(strconv.FormatInt(userID, 10), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}

func (h *PointsHandler) GetLeaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}
	entries, err := h.repo.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"leaderboard": entries})
}
