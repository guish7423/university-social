package model

import "time"

type Product struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	OriginalPrice *float64  `json:"original_price,omitempty"`
	Category      string    `json:"category"`
	Condition     string    `json:"condition"`
	Images        []string  `json:"images"`
	Contact       string    `json:"contact"`
	Status        int       `json:"status"`
	LikeCount     int       `json:"like_count"`
	CommentCount  int       `json:"comment_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsLiked       bool      `json:"is_liked,omitempty"`
	IsOwner       bool      `json:"is_owner,omitempty"`
}

type CreateProductRequest struct {
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required"`
	OriginalPrice *float64 `json:"original_price"`
	Category      string   `json:"category" binding:"required"`
	Condition     string   `json:"condition"`
	Images        []string `json:"images"`
	Contact       string   `json:"contact"`
}

type ProductComment struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
