package model

import "time"

type Circle struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Cover       string    `json:"cover"`
	CreatorID   int64     `json:"creator_id"`
	MemberCount int       `json:"member_count"`
	PostCount   int       `json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
	IsMember    bool      `json:"is_member,omitempty"`
	Creator     *User     `json:"creator,omitempty"`
}

type CircleMember struct {
	ID        int64     `json:"id"`
	CircleID  int64     `json:"circle_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

type CirclePost struct {
	ID            int64     `json:"id"`
	CircleID      int64     `json:"circle_id"`
	UserID        int64     `json:"user_id"`
	Content       string    `json:"content"`
	ImageUrls     []string  `json:"image_urls"`
	LikeCount     int       `json:"like_count"`
	CommentCount  int       `json:"comment_count"`
	CreatedAt     time.Time `json:"created_at"`
	Author        *User     `json:"author,omitempty"`
	IsLiked       bool      `json:"is_liked,omitempty"`
}

type CreateCircleRequest struct {
	Name        string `json:"name" binding:"required,max=64"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Cover       string `json:"cover"`
}

type CreateCirclePostRequest struct {
	Content   string   `json:"content" binding:"required"`
	ImageUrls []string `json:"image_urls"`
}

type CirclePostComment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author,omitempty"`
}
