package model

import "time"

type CircleActivity struct {
	ID        int64     `json:"id"`
	CircleID  int64     `json:"circle_id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}
