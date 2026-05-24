package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type FriendHandler struct {
	repo *repository.FriendRepository
}

func NewFriendHandler(repo *repository.FriendRepository) *FriendHandler {
	return &FriendHandler{repo: repo}
}

func (h *FriendHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}
	users, err := h.repo.SearchUsers(query, 20, c.GetInt64("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	if users == nil {
		users = []*repository.UserSearchResult{}
	}
	c.JSON(http.StatusOK, users)
}

func (h *FriendHandler) SendRequest(c *gin.Context) {
	userID := c.GetInt64("user_id")
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userID == friendID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.SendRequest(userID, friendID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败"})
		return
	}
	h.repo.CreateNotification(friendID, userID, "friend_request", "请求添加你为好友", 0)
	c.JSON(http.StatusOK, gin.H{"message": "请求已发送"})
}

func (h *FriendHandler) AcceptRequest(c *gin.Context) {
	userID := c.GetInt64("user_id")
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.AcceptRequest(userID, friendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "操作失败"})
		return
	}
	h.repo.CreateNotification(friendID, userID, "friend_accepted", "通过了你的好友请求", 0)
	c.JSON(http.StatusOK, gin.H{"message": "已同意"})
}

func (h *FriendHandler) RemoveFriend(c *gin.Context) {
	userID := c.GetInt64("user_id")
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.Remove(userID, friendID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
func (h *FriendHandler) RejectRequest(c *gin.Context) {
	userID := c.GetInt64("user_id")
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.RejectRequest(userID, friendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已拒绝"})
}

func (h *FriendHandler) ListFriends(c *gin.Context) {
	userID := c.GetInt64("user_id")
	friends, err := h.repo.ListFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if friends == nil {
		friends = []*model.User{}
	}
	c.JSON(http.StatusOK, friends)
}

func (h *FriendHandler) ListRequests(c *gin.Context) {
	userID := c.GetInt64("user_id")
	requests, err := h.repo.ListRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if requests == nil {
		requests = []*model.User{}
	}
	c.JSON(http.StatusOK, requests)
}

func (h *FriendHandler) ListNotifications(c *gin.Context) {
	userID := c.GetInt64("user_id")
	notes, err := h.repo.ListNotifications(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if notes == nil {
		notes = []*model.Notification{}
	}
	c.JSON(http.StatusOK, notes)
}

func (h *FriendHandler) MarkRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	h.repo.MarkNotificationsRead(userID)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
