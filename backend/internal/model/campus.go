package model

import "time"

type CampusCalendar struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	EventDate   string    `json:"event_date"`
	EventType   string    `json:"event_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CampusDirectory struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Office     string    `json:"office"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
}

type EmptyClassroom struct {
	ID        int64     `json:"id"`
	Building  string    `json:"building"`
	Room      string    `json:"room"`
	Capacity  int       `json:"capacity"`
	Weekday   int       `json:"weekday"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Semester  string    `json:"semester"`
	CreatedAt time.Time `json:"created_at"`
}
