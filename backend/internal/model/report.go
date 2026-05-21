package model

import "time"

type Report struct {
	ID          int64      `json:"id"`
	ReporterID  int64      `json:"reporter_id"`
	ContentType string     `json:"content_type"`
	ContentID   int64      `json:"content_id"`
	Reason      string     `json:"reason"`
	Detail      string     `json:"detail,omitempty"`
	Status      string     `json:"status"`
	ResolvedBy  *int64     `json:"resolved_by,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateReportRequest struct {
	ContentType string `json:"content_type" binding:"required"`
	ContentID   int64  `json:"content_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	Detail      string `json:"detail"`
}
