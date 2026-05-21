package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(userID int64, req *model.CreateReportRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO reports (reporter_id, content_type, content_id, reason, detail)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, req.ContentType, req.ContentID, req.Reason, req.Detail,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create report: %w", err)
	}
	return id, nil
}

func (r *ReportRepository) List(status string, offset, limit int) ([]*model.Report, error) {
	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = r.db.Query(
			`SELECT id, reporter_id, content_type, content_id, reason, detail,
			        status, resolved_by, resolved_at, created_at
			 FROM reports WHERE status = $1
			 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			status, limit, offset,
		)
	} else {
		rows, err = r.db.Query(
			`SELECT id, reporter_id, content_type, content_id, reason, detail,
			        status, resolved_by, resolved_at, created_at
			 FROM reports
			 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
			limit, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()

	var reports []*model.Report
	for rows.Next() {
		r := &model.Report{}
		if err := rows.Scan(&r.ID, &r.ReporterID, &r.ContentType, &r.ContentID,
			&r.Reason, &r.Detail, &r.Status, &r.ResolvedBy, &r.ResolvedAt, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, r)
	}
	if reports == nil {
		reports = []*model.Report{}
	}
	return reports, nil
}

func (r *ReportRepository) Resolve(id int64, adminID int64) error {
	res, err := r.db.Exec(
		`UPDATE reports SET status = 'resolved', resolved_by = $1, resolved_at = NOW()
		 WHERE id = $2 AND status = 'pending'`,
		adminID, id,
	)
	if err != nil {
		return fmt.Errorf("resolve report: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found or already resolved")
	}
	return nil
}

func (r *ReportRepository) Dismiss(id int64, adminID int64) error {
	res, err := r.db.Exec(
		`UPDATE reports SET status = 'dismissed', resolved_by = $1, resolved_at = NOW()
		 WHERE id = $2 AND status = 'pending'`,
		adminID, id,
	)
	if err != nil {
		return fmt.Errorf("dismiss report: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found or already resolved")
	}
	return nil
}

func (r *ReportRepository) CountByStatus(status string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM reports WHERE status = $1`, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count reports: %w", err)
	}
	return count, nil
}
