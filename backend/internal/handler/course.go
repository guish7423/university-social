package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type CourseHandler struct {
	courseRepo *repository.CourseRepository
}

func NewCourseHandler(courseRepo *repository.CourseRepository) *CourseHandler {
	return &CourseHandler{courseRepo: courseRepo}
}

func (h *CourseHandler) SearchCourses(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	if query == "" {
		c.JSON(http.StatusOK, []*model.Course{})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	courses, err := h.courseRepo.SearchCourses(query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	if courses == nil {
		courses = []*model.Course{}
	}
	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetCourseDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	detail, err := h.courseRepo.GetCourseDetail(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *CourseHandler) CreateRating(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req model.CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写评分"})
		return
	}

	ratingID, err := h.courseRepo.CreateRating(userID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评价失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": ratingID})
}

func (h *CourseHandler) GetUserRating(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	rating, err := h.courseRepo.GetUserRating(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, rating)
}

func (h *CourseHandler) ListRatings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	order := c.DefaultQuery("order", "newest")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ratings, err := h.courseRepo.ListRatings(id, order, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}
	if ratings == nil {
		ratings = []*model.CourseRating{}
	}
	c.JSON(http.StatusOK, ratings)
}

func (h *CourseHandler) MarkHelpful(c *gin.Context) {
	userID := c.GetInt64("user_id")
	ratingID, err := strconv.ParseInt(c.Param("ratingId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.courseRepo.MarkHelpful(ratingID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记"})
}
