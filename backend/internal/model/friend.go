package model

import "time"

type Friend struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FriendID  int64     `json:"friend_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Friend    *User     `json:"friend,omitempty"`
}

type FriendRequest struct {
	FriendID int64 `json:"friend_id" binding:"required"`
}

type Notification struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	FromUserID int64     `json:"from_user_id"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	RefID      int64     `json:"ref_id"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	FromUser   *User     `json:"from_user,omitempty"`
}
