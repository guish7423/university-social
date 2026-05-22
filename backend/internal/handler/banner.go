package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/cache"
	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type BannerHandler struct {
	repo *repository.BannerRepository
	rdb  *cache.Cache
}

func NewBannerHandler(repo *repository.BannerRepository, rdb *cache.Cache) *BannerHandler {
	return &BannerHandler{repo: repo, rdb: rdb}
}

const cacheKeyBanners = "cache:banners"
const cacheKeyAnnouncements = "cache:announcements"

func (h *BannerHandler) ListActive(c *gin.Context) {
	ctx := context.Background()
	if h.rdb != nil {
		var cached []*model.Banner
		if err := h.rdb.Get(ctx, cacheKeyBanners, &cached); err == nil {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	banners, err := h.repo.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Banner失败"})
		return
	}
	if banners == nil {
		banners = []*model.Banner{}
	}
	if h.rdb != nil {
		h.rdb.Set(ctx, cacheKeyBanners, banners, 5*time.Minute)
	}
	c.JSON(http.StatusOK, banners)
}

func (h *BannerHandler) AdminList(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	banners, err := h.repo.ListAll(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Banner列表失败"})
		return
	}
	if banners == nil {
		banners = []*model.Banner{}
	}
	c.JSON(http.StatusOK, banners)
}

func (h *BannerHandler) Create(c *gin.Context) {
	var req model.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	id, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyBanners)
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *BannerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner不存在"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyBanners)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BannerHandler) ToggleActive(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.ToggleActive(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyBanners)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BannerHandler) ListAnnouncementsActive(c *gin.Context) {
	ctx := context.Background()
	if h.rdb != nil {
		var cached []*model.Announcement
		if err := h.rdb.Get(ctx, cacheKeyAnnouncements, &cached); err == nil {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	anns, err := h.repo.ListAnnouncementsActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败"})
		return
	}
	if anns == nil {
		anns = []*model.Announcement{}
	}
	if h.rdb != nil {
		h.rdb.Set(ctx, cacheKeyAnnouncements, anns, 5*time.Minute)
	}
	c.JSON(http.StatusOK, anns)
}

func (h *BannerHandler) AdminListAnnouncements(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	anns, err := h.repo.ListAllAnnouncements(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告列表失败"})
		return
	}
	if anns == nil {
		anns = []*model.Announcement{}
	}
	c.JSON(http.StatusOK, anns)
}

func (h *BannerHandler) CreateAnnouncement(c *gin.Context) {
	var req model.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	id, err := h.repo.CreateAnnouncement(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyAnnouncements)
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *BannerHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.DeleteAnnouncement(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyAnnouncements)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BannerHandler) ToggleAnnouncementActive(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.ToggleAnnouncementActive(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	if h.rdb != nil {
		h.rdb.Delete(context.Background(), cacheKeyAnnouncements)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BannerHandler) SetFeatured(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.SetPostFeatured(postID, true); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已设为精华"})
}

func (h *BannerHandler) UnsetFeatured(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.SetPostFeatured(postID, false); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已取消精华"})
}

func (h *BannerHandler) ListFeatured(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	posts, err := h.repo.ListFeaturedPosts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取精华帖失败"})
		return
	}
	if posts == nil {
		posts = []*model.Post{}
	}
	c.JSON(http.StatusOK, posts)
}
