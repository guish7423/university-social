package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type ChatHandler struct {
	msgRepo *repository.MessageRepository
}

func NewChatHandler(msgRepo *repository.MessageRepository) *ChatHandler {
	return &ChatHandler{msgRepo: msgRepo}
}

func (h *ChatHandler) Send(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	msg, err := h.msgRepo.Send(userID, req.ReceiverID, req.Content, req.MsgType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) Conversations(c *gin.Context) {
	userID := c.GetInt64("user_id")
	convs, err := h.msgRepo.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, convs)
}

func (h *ChatHandler) Messages(c *gin.Context) {
	userID := c.GetInt64("user_id")
	otherUserID := c.Param("userId")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	msgs, err := h.msgRepo.GetMessages(strconv.FormatInt(userID, 10), otherUserID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (h *ChatHandler) MarkRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	otherUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.msgRepo.MarkRead(userID, otherUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ChatHandler) UnreadCount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	count, err := h.msgRepo.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
