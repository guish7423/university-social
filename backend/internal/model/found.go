package model

import "time"

type LostItem struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ItemType    string    `json:"item_type"`
	Location    string    `json:"location"`
	Contact     string    `json:"contact"`
	Images      []string  `json:"images"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Nickname    string    `json:"nickname,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
}

type CreateLostItemRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Category    string   `json:"category" binding:"required,oneof=lost found"`
	ItemType    string   `json:"item_type"`
	Location    string   `json:"location"`
	Contact     string   `json:"contact"`
	Images      []string `json:"images"`
}

type UpdateLostItemStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1 2"`
}
