package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/repository"
)

type FavoriteHandler struct {
	favRepo  *repository.FavoriteRepository
	userRepo *repository.UserRepository
}

func NewFavoriteHandler(favRepo *repository.FavoriteRepository, userRepo *repository.UserRepository) *FavoriteHandler {
	return &FavoriteHandler{favRepo: favRepo, userRepo: userRepo}
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.favRepo.Add(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.favRepo.Remove(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *FavoriteHandler) Status(c *gin.Context) {
	userID := c.GetInt64("user_id")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	isFav, _ := h.favRepo.IsFavorited(userID, postID)
	c.JSON(http.StatusOK, gin.H{"is_favorited": isFav})
}

func (h *FavoriteHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	posts, total, err := h.favRepo.ListByUser(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	type postResp struct {
		ID            int64            `json:"id"`
		UserID        int64            `json:"user_id"`
		Author        *userBrief       `json:"author"`
		Content       string           `json:"content"`
		Images        []string         `json:"images"`
		LikeCount     int              `json:"like_count"`
		CommentCount  int              `json:"comment_count"`
		CreatedAt     string           `json:"created_at"`
		IsFavorited   bool             `json:"is_favorited"`
	}

	result := make([]postResp, 0, len(posts))
	for _, p := range posts {
		author := &userBrief{}
		u, err := h.userRepo.FindByID(p.UserID)
		if err == nil && u != nil {
			author = &userBrief{Nickname: u.Nickname, Avatar: u.Avatar}
		}
		result = append(result, postResp{
			ID:           p.ID,
			UserID:       p.UserID,
			Author:       author,
			Content:      p.Content,
			Images:       p.Images,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			IsFavorited:  true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": result,
		"total": total,
		"page":  page,
	})
}

type userBrief struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
