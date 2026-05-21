package model

import "time"

type Activity struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ActivityType     string    `json:"activity_type"`
	Location         string    `json:"location"`
	MaxParticipants  int       `json:"max_participants"`
	StartTime        time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	Images           []string  `json:"images"`
	Status           int       `json:"status"`
	ParticipantCount int       `json:"participant_count"`
	CommentCount     int       `json:"comment_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsParticipant    bool      `json:"is_participant,omitempty"`
	IsOwner          bool      `json:"is_owner,omitempty"`
}

type CreateActivityRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description"`
	ActivityType    string   `json:"activity_type" binding:"required"`
	Location        string   `json:"location"`
	MaxParticipants int      `json:"max_participants"`
	StartTime       string   `json:"start_time" binding:"required"`
	EndTime         string   `json:"end_time"`
	Images          []string `json:"images"`
}

type ActivityComment struct {
	ID         int64     `json:"id"`
	ActivityID int64     `json:"activity_id"`
	UserID     int64     `json:"user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
