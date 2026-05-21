package model

import "time"

type InviteCode struct {
	ID            int64      `json:"id"`
	Code          string     `json:"code"`
	CreatorID     int64      `json:"creator_id"`
	InvitedUserID *int64     `json:"invited_user_id,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UsedAt        *time.Time `json:"used_at,omitempty"`
	Creator       *User      `json:"creator,omitempty"`
	InvitedUser   *User      `json:"invited_user,omitempty"`
}

type CreateInviteCodeRequest struct {
	Count int `json:"count" binding:"min=1,max=5"`
}

type RedeemInviteCodeRequest struct {
	Code string `json:"code" binding:"required"`
}
