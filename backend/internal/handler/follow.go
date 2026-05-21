package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type FollowHandler struct {
	followRepo      *repository.FollowRepository
	notificationRepo *repository.FriendRepository
}

func NewFollowHandler(followRepo *repository.FollowRepository, notificationRepo *repository.FriendRepository) *FollowHandler {
	return &FollowHandler{followRepo: followRepo, notificationRepo: notificationRepo}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	followingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.followRepo.Follow(userID, followingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	// 创建通知
	h.notificationRepo.CreateNotification(followingID, userID, "follow", "关注了你", 0)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	followingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.followRepo.Unfollow(userID, followingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {
	userID := c.GetInt64("user_id")
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	isFollowing, _ := h.followRepo.IsFollowing(userID, targetID)
	c.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}

func (h *FollowHandler) Following(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		targetID = c.GetInt64("user_id")
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	follows, err := h.followRepo.ListFollowing(targetID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if follows == nil {
		follows = []*model.Follow{}
	}
	c.JSON(http.StatusOK, follows)
}

func (h *FollowHandler) Followers(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		targetID = c.GetInt64("user_id")
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	follows, err := h.followRepo.ListFollowers(targetID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if follows == nil {
		follows = []*model.Follow{}
	}
	c.JSON(http.StatusOK, follows)
}

func (h *FollowHandler) Counts(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		targetID = c.GetInt64("user_id")
	}
	counts, err := h.followRepo.Counts(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, counts)
}
