package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/repository"
)

type CampusHandler struct {
	repo *repository.CampusRepository
}

func NewCampusHandler(repo *repository.CampusRepository) *CampusHandler {
	return &CampusHandler{repo: repo}
}

func (h *CampusHandler) ListCalendar(c *gin.Context) {
	items, err := h.repo.ListCalendar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取校历失败"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CampusHandler) ListDirectory(c *gin.Context) {
	dept := c.Query("department")
	items, err := h.repo.ListDirectory(dept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取黄页失败"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CampusHandler) ListEmptyRooms(c *gin.Context) {
	building := c.Query("building")
	weekday, _ := strconv.Atoi(c.DefaultQuery("weekday", "0"))
	items, err := h.repo.ListEmptyRooms(building, weekday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取空教室失败"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CampusHandler) ListBuildings(c *gin.Context) {
	buildings, err := h.repo.ListBuildings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教学楼列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buildings": buildings})
}
