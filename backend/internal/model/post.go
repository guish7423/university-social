package model

import "time"

type Post struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Content       string    `json:"content"`
	Images        []string  `json:"images"`
	TopicID       *int64    `json:"topic_id,omitempty"`
	School        string    `json:"school"`
	LikeCount     int       `json:"like_count"`
	CommentCount  int       `json:"comment_count"`
	IsFeatured    bool      `json:"is_featured"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Author        *User     `json:"author,omitempty"`
	IsLiked       bool      `json:"is_liked,omitempty"`
	IsBlocked     bool      `json:"is_blocked,omitempty"`
}

type CreatePostRequest struct {
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images"`
	TopicID *int64   `json:"topic_id"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author,omitempty"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type Like struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Topic struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	PostCount int       `json:"post_count"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdatePostRequest struct {
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images"`
	TopicID *int64   `json:"topic_id"`
}
