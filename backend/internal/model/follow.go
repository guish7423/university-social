package model

import "time"

type Follow struct {
	ID          int64     `json:"id"`
	FollowerID  int64     `json:"follower_id"`
	FollowingID int64     `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
	Following   *User     `json:"following,omitempty"`
	Follower    *User     `json:"follower,omitempty"`
}

type FollowCounts struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}
