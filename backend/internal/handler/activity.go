package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type ActivityHandler struct {
	activityRepo *repository.ActivityRepository
}

func NewActivityHandler(activityRepo *repository.ActivityRepository) *ActivityHandler {
	return &ActivityHandler{activityRepo: activityRepo}
}

func (h *ActivityHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写必填字段"})
		return
	}
	id, err := h.activityRepo.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *ActivityHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	activityType := c.DefaultQuery("type", "")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	activities, err := h.activityRepo.List(activityType, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	if activities == nil {
		activities = []*model.Activity{}
	}
if activities == nil { activities = []*model.Activity{} }
	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	a, err := h.activityRepo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.activityRepo.Delete(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func (h *ActivityHandler) Join(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.activityRepo.Join(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "报名成功"})
}

func (h *ActivityHandler) Leave(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	h.activityRepo.Leave(id, userID)
	c.JSON(http.StatusOK, gin.H{"message": "已取消报名"})
}

func (h *ActivityHandler) Participants(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	users, err := h.activityRepo.ListParticipants(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if users == nil {
		users = []*model.User{}
	}
if users == nil { users = []*model.User{} }
	c.JSON(http.StatusOK, users)
}

func (h *ActivityHandler) CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入内容"})
		return
	}
	commentID, err := h.activityRepo.CreateComment(userID, id, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": commentID})
}

func (h *ActivityHandler) ListComments(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	comments, err := h.activityRepo.ListComments(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
	if comments == nil {
		comments = []*model.ActivityComment{}
	}
if comments == nil { comments = []*model.ActivityComment{} }
	c.JSON(http.StatusOK, comments)
}

func (h *ActivityHandler) MyActivities(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	created, _ := h.activityRepo.ListByUser(userID, userID, offset, limit)
	joined, _ := h.activityRepo.ListJoined(userID, 0, 10)

	if created == nil {
		created = []*model.Activity{}
	}
	if joined == nil {
		joined = []*model.Activity{}
	}
	c.JSON(http.StatusOK, gin.H{"created": created, "joined": joined})
}

func (h *ActivityHandler) Types(c *gin.Context) {
	types := h.activityRepo.Types()
if types == nil { types = []string{} }
	c.JSON(http.StatusOK, types)
}
