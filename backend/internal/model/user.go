package model

import "time"

type User struct {
	ID          int64     `json:"id" db:"id"`
	OpenID      string    `json:"open_id" db:"open_id"`
	UnionID     string    `json:"union_id,omitempty" db:"union_id"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Avatar      string    `json:"avatar" db:"avatar"`

	Phone       string    `json:"phone,omitempty" db:"phone"`
	Email       string    `json:"email,omitempty" db:"email"`
	School      string    `json:"school,omitempty" db:"school"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsBanned    bool      `json:"is_banned" db:"is_banned"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UserProfile struct {
	User
	PostCount  int `json:"post_count"`
	FriendCount int `json:"friend_count"`
}
