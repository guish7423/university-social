package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/repository"
)

type SensitiveHandler struct {
	repo *repository.SensitiveWordRepository
}

func NewSensitiveHandler(repo *repository.SensitiveWordRepository) *SensitiveHandler {
	return &SensitiveHandler{repo: repo}
}

type SensitiveWordInput struct {
	Word string `json:"word" binding:"required"`
}

func (h *SensitiveHandler) List(c *gin.Context) {
	words, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取敏感词列表失败"})
		return
	}
	if words == nil {
		words = []string{}
	}
	c.JSON(http.StatusOK, words)
}

func (h *SensitiveHandler) Add(c *gin.Context) {
	var input SensitiveWordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入敏感词"})
		return
	}
	if err := h.repo.Add(input.Word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "已添加"})
}

func (h *SensitiveHandler) Remove(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.Remove(word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
