package model

import "time"

type University struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain,omitempty"`
	LogoURL   string    `json:"logo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	MemberCount int     `json:"member_count,omitempty"`
	IsMember  bool      `json:"is_member,omitempty"`
}

type UniversityMember struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	UniversityID int64     `json:"university_id"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	User         *User     `json:"user,omitempty"`
}
