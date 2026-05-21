package model

import "time"

type Course struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Teacher    string    `json:"teacher"`
	School     string    `json:"school"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

type CourseRating struct {
	ID             int64     `json:"id"`
	CourseID       int64     `json:"course_id"`
	UserID         int64     `json:"user_id"`
	TeachingRating int       `json:"teaching_rating"`
	Difficulty     int       `json:"difficulty"`
	Grading        int       `json:"grading"`
	Comment        string    `json:"comment"`
	IsAnonymous    bool      `json:"is_anonymous"`
	HelpfulCount   int       `json:"helpful_count"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Nickname       string    `json:"nickname,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
}

type CreateRatingRequest struct {
	TeachingRating int    `json:"teaching_rating" binding:"required,min=1,max=5"`
	Difficulty     int    `json:"difficulty" binding:"required,min=1,max=5"`
	Grading        int    `json:"grading" binding:"required,min=1,max=5"`
	Comment        string `json:"comment"`
	IsAnonymous    bool   `json:"is_anonymous"`
}

type CourseDetail struct {
	Course
	AvgTeaching float64        `json:"avg_teaching"`
	AvgDifficulty float64      `json:"avg_difficulty"`
	AvgGrading   float64       `json:"avg_grading"`
	RatingCount  int           `json:"rating_count"`
	UserRating   *CourseRating `json:"user_rating,omitempty"`
}
