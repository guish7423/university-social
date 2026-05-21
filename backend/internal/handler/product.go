package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

type ProductHandler struct {
	productRepo *repository.ProductRepository
}

func NewProductHandler(productRepo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

func (h *ProductHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写必填字段"})
		return
	}

	if req.Condition == "" {
		req.Condition = "轻微使用"
	}

	id, err := h.productRepo.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *ProductHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	category := c.DefaultQuery("category", "")
	sort := c.DefaultQuery("sort", "")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, err := h.productRepo.List(category, sort, 0, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	if products == nil {
		products = []*model.Product{}
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	product, err := h.productRepo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.productRepo.Delete(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func (h *ProductHandler) MarkSold(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.productRepo.MarkSold(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已出"})
}

func (h *ProductHandler) MyProducts(c *gin.Context) {
	userID := c.GetInt64("user_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, err := h.productRepo.ListByUser(userID, userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	if products == nil {
		products = []*model.Product{}
	}

	total, active := h.productRepo.Counts(userID)
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"active":   active,
	})
}

func (h *ProductHandler) Search(c *gin.Context) {
	userID := c.GetInt64("user_id")
	query := c.DefaultQuery("q", "")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if query == "" {
		c.JSON(http.StatusOK, []*model.Product{})
		return
	}

	products, err := h.productRepo.SearchProducts(query, offset, limit, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	if products == nil {
		products = []*model.Product{}
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) ToggleLike(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	liked, err := h.productRepo.ToggleLike(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *ProductHandler) CreateComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req model.ProductCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入内容"})
		return
	}

	commentID, err := h.productRepo.CreateComment(userID, id, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": commentID})
}

func (h *ProductHandler) ListComments(c *gin.Context) {
	userID := c.GetInt64("user_id")
	paramID := c.Param("id")
	if paramID == "" {
		paramID = c.Param("productId")
	}
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// verify product exists
	_, err = h.productRepo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	comments, err := h.productRepo.ListComments(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}
	if comments == nil {
		comments = []*model.ProductComment{}
	}
	c.JSON(http.StatusOK, comments)
}

func (h *ProductHandler) Categories(c *gin.Context) {
	cats := h.productRepo.Categories()
	c.JSON(http.StatusOK, cats)
}

func (h *ProductHandler) Trending(c *gin.Context) {
	userID := c.GetInt64("user_id")
	products, err := h.productRepo.Trending(10, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取热门商品失败: %v", err)})
		return
	}
	c.JSON(http.StatusOK, products)
}
