package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"


	"github.com/guish/university-social/internal/repository"
)

type UniversityHandler struct {
	uniRepo *repository.UniversityRepository
}

func NewUniversityHandler(uniRepo *repository.UniversityRepository) *UniversityHandler {
	return &UniversityHandler{uniRepo: uniRepo}
}

func (h *UniversityHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	universities, err := h.uniRepo.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取高校列表失败"})
		return
	}
	c.JSON(http.StatusOK, universities)
}

func (h *UniversityHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	university, err := h.uniRepo.GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "高校不存在"})
		return
	}
	c.JSON(http.StatusOK, university)
}

func (h *UniversityHandler) Join(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.uniRepo.Join(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *UniversityHandler) Leave(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.uniRepo.Leave(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *UniversityHandler) Members(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	members, err := h.uniRepo.GetMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员列表失败"})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *UniversityHandler) MyUniversity(c *gin.Context) {
	userID := c.GetInt64("user_id")
	university, err := h.uniRepo.GetUserUniversity(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未加入任何高校"})
		return
	}
	c.JSON(http.StatusOK, university)
}
