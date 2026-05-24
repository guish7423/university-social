package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type BlockHandler struct {
	blockRepo *repository.BlockRepository
}

func NewBlockHandler(blockRepo *repository.BlockRepository) *BlockHandler {
	return &BlockHandler{blockRepo: blockRepo}
}

func (h *BlockHandler) Block(c *gin.Context) {
	userID := c.GetInt64("user_id")
	blockedID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userID == blockedID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.blockRepo.Block(userID, blockedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BlockHandler) Unblock(c *gin.Context) {
	userID := c.GetInt64("user_id")
	blockedID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.blockRepo.Unblock(userID, blockedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BlockHandler) BlockStatus(c *gin.Context) {
	userID := c.GetInt64("user_id")
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	isBlocked, _ := h.blockRepo.IsBlocking(userID, targetID)
	c.JSON(http.StatusOK, gin.H{"is_blocked": isBlocked})
}

func (h *BlockHandler) ListBlocked(c *gin.Context) {
	userID := c.GetInt64("user_id")
	users, err := h.blockRepo.ListBlockedUsers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if users == nil {
		users = []*model.User{}
	}
	c.JSON(http.StatusOK, users)
}
