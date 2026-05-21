package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type InviteCodeHandler struct {
	inviteRepo *repository.InviteCodeRepository
	userRepo   *repository.UserRepository
	pointsRepo *repository.PointsRepository
}

func NewInviteCodeHandler(inviteRepo *repository.InviteCodeRepository, userRepo *repository.UserRepository, pointsRepo *repository.PointsRepository) *InviteCodeHandler {
	return &InviteCodeHandler{inviteRepo: inviteRepo, userRepo: userRepo, pointsRepo: pointsRepo}
}

func (h *InviteCodeHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Count = 1
	}

	if req.Count < 1 {
		req.Count = 1
	}
	if req.Count > 5 {
		req.Count = 5
	}

	total, used, err := h.inviteRepo.CreatorStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询邀请码失败"})
		return
	}
	if total-used >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已拥有足够邀请码"})
		return
	}
	remaining := 5 - (total - used)
	if req.Count > remaining {
		req.Count = remaining
	}

	codes, err := h.inviteRepo.Create(userID, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建邀请码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"codes": codes})
}

func (h *InviteCodeHandler) MyCodes(c *gin.Context) {
	userID := c.GetInt64("user_id")
	codes, err := h.inviteRepo.ListByCreator(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取邀请码列表失败"})
		return
	}
	total, used, _ := h.inviteRepo.CreatorStats(userID)
	c.JSON(http.StatusOK, gin.H{
		"codes":           codes,
		"total_created":   total,
		"total_used":      used,
		"remaining_slots": 5 - (total - used),
	})
}

func (h *InviteCodeHandler) Redeem(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.RedeemInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供邀请码"})
		return
	}

	code, err := h.inviteRepo.FindByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询邀请码失败"})
		return
	}
	if code == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请码不存在"})
		return
	}
	if code.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码已使用"})
		return
	}

	err = h.inviteRepo.Redeem(code.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "兑换邀请码失败"})
		return
	}

	h.pointsRepo.AddPoints(code.CreatorID, "invite_redeemed", 5, "邀请注册奖励")
	c.JSON(http.StatusOK, gin.H{"message": "兑换成功"})
}
