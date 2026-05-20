package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type CircleHandler struct {
	repo *repository.CircleRepository
}

func NewCircleHandler(repo *repository.CircleRepository) *CircleHandler {
	return &CircleHandler{repo: repo}
}

func (h *CircleHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	circles, err := h.repo.List(offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取圈子列表失败"})
		return
	}
	if circles == nil {
		circles = []*model.Circle{}
	}
	c.JSON(http.StatusOK, circles)
}

func (h *CircleHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	circle, err := h.repo.GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取圈子失败"})
		return
	}
	if circle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "圈子不存在"})
		return
	}
	c.JSON(http.StatusOK, circle)
}

func (h *CircleHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateCircleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	id, err := h.repo.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建圈子失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *CircleHandler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req model.CreateCircleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.repo.Update(id, userID, &req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改或圈子不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *CircleHandler) Join(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.repo.Join(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "加入成功"})
}

func (h *CircleHandler) Leave(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.repo.Leave(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已退出"})
}

func (h *CircleHandler) ListMembers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	members, err := h.repo.ListMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员列表失败"})
		return
	}
	if members == nil {
		members = []*model.CircleMember{}
	}
	c.JSON(http.StatusOK, members)
}

func (h *CircleHandler) CreatePost(c *gin.Context) {
	userID := c.GetInt64("user_id")
	circleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req model.CreateCirclePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	id, err := h.repo.CreatePost(circleID, userID, req.Content, req.ImageUrls)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发帖失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *CircleHandler) ListPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.repo.ListPosts(id, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}
	if posts == nil {
		posts = []*model.CirclePost{}
	}
	c.JSON(http.StatusOK, posts)
}

func (h *CircleHandler) TogglePostLike(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	liked, err := h.repo.ToggleLike(postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *CircleHandler) CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}

	id, err := h.repo.CreateComment(postID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *CircleHandler) ListComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	comments, err := h.repo.ListComments(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
	if comments == nil {
		comments = []*model.CirclePostComment{}
	}
	c.JSON(http.StatusOK, comments)
}
