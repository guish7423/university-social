package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type WhisperHandler struct {
	repo         *repository.WhisperRepository
	sensitiveRepo *repository.SensitiveWordRepository
}

func NewWhisperHandler(repo *repository.WhisperRepository, sensitiveRepo *repository.SensitiveWordRepository) *WhisperHandler {
	return &WhisperHandler{repo: repo, sensitiveRepo: sensitiveRepo}
}

func (h *WhisperHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateWhisperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}
	filtered, hit := h.sensitiveRepo.ReplaceSensitive(req.Content)
	if hit {
		req.Content = filtered
	}
	id, err := h.repo.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *WhisperHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	whispers, err := h.repo.List(offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取树洞失败"})
		return
	}
	c.JSON(http.StatusOK, whispers)
}

func (h *WhisperHandler) Get(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	whisper, err := h.repo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "树洞不存在"})
		return
	}
	c.JSON(http.StatusOK, whisper)
}

func (h *WhisperHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.Delete(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *WhisperHandler) ToggleLike(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	liked, err := h.repo.ToggleLike(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *WhisperHandler) CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	whisperID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var req model.WhisperCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}
	filtered, hit := h.sensitiveRepo.ReplaceSensitive(req.Content)
	if hit {
		req.Content = filtered
	}
	id, err := h.repo.CreateComment(whisperID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *WhisperHandler) ListComments(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	comments, err := h.repo.ListComments(int(id), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *WhisperHandler) Summary(c *gin.Context) {
	s, err := h.repo.Summary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}
	c.JSON(http.StatusOK, s)
}
