package model

import "time"

type VerificationCode struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Code      string    `json:"-"`
	School    string    `json:"school"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SendCodeRequest struct {
	Email  string `json:"email" binding:"required"`
	School string `json:"school" binding:"required"`
}

type VerifyCodeRequest struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email" binding:"required"`
}
