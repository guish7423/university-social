package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type PostHandler struct {
	repo         *repository.PostRepository
	sensitiveRepo *repository.SensitiveWordRepository
	pointsRepo   *repository.PointsRepository
}

func NewPostHandler(repo *repository.PostRepository, sensitiveRepo *repository.SensitiveWordRepository, pointsRepo *repository.PointsRepository) *PostHandler {
	return &PostHandler{repo: repo, sensitiveRepo: sensitiveRepo, pointsRepo: pointsRepo}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreatePostRequest
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
	h.pointsRepo.AddPoints(userID, "create_post", 3, "发布动态奖励")
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	school := c.Query("school")
	topicStr := c.Query("topic_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var topicID *int64
	if topicStr != "" {
		if id, err := strconv.ParseInt(topicStr, 10, 64); err == nil {
			topicID = &id
		}
	}

	posts, err := h.repo.List(school, topicID, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取动态失败"})
		return
	}
if posts == nil { posts = []*model.Post{} }
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	post, err := h.repo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "动态不存在"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
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

func (h *PostHandler) CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var req model.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}
	filtered, hit := h.sensitiveRepo.ReplaceSensitive(req.Content)
	if hit {
		req.Content = filtered
	}
	id, err := h.repo.CreateComment(postID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
	h.pointsRepo.AddPoints(userID, "create_comment", 1, "评论奖励")
}

func (h *PostHandler) ListComments(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	comments, err := h.repo.ListComments(int(postID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
if comments == nil { comments = []*model.Comment{} }
	c.JSON(http.StatusOK, comments)
}

func (h *PostHandler) ToggleLike(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	liked, err := h.repo.ToggleLike(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *PostHandler) ListTopics(c *gin.Context) {
	topics, err := h.repo.ListTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取话题失败"})
		return
	}
if topics == nil { topics = []*model.Topic{} }
	c.JSON(http.StatusOK, topics)
}

func (h *PostHandler) ShareCard(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	post, err := h.repo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "动态不存在"})
		return
	}

	cardTitle := post.Content
	if len([]rune(cardTitle)) > 50 {
		cardTitle = string([]rune(cardTitle)[:50]) + "..."
	}

	var cardImage string
	if len(post.Images) > 0 {
		cardImage = post.Images[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"title":       cardTitle,
		"description": "校园社 - 发现校园新鲜事",
		"image":       cardImage,
		"path":        "/pages/post/detail?id=" + strconv.FormatInt(post.ID, 10),
		"author":      post.Author,
		"like_count":  post.LikeCount,
	})
}

func (h *PostHandler) SearchPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	query := c.Query("q")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	posts, err := h.repo.SearchPosts(query, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	if posts == nil { posts = []*model.Post{} }
if posts == nil { posts = []*model.Post{} }
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) TrendingPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	posts, err := h.repo.TrendingPosts(offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热帖失败"})
		return
	}
	if posts == nil { posts = []*model.Post{} }
if posts == nil { posts = []*model.Post{} }
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) FollowingPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	posts, err := h.repo.FollowingPosts(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取关注动态失败"})
		return
	}
	if posts == nil { posts = []*model.Post{} }
if posts == nil { posts = []*model.Post{} }
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) UserPosts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	posts, err := h.repo.ListUserPosts(targetID, userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}
	if posts == nil { posts = []*model.Post{} }
if posts == nil { posts = []*model.Post{} }
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}
	if err := h.repo.Update(id, userID, &req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "编辑失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
