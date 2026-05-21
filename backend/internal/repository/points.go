package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guish/university-social/internal/model"
)

type PointsRepository struct {
	db *sql.DB
}

func NewPointsRepository(db *sql.DB) *PointsRepository {
	return &PointsRepository{db: db}
}

func (r *PointsRepository) AddPoints(userID int64, action string, points int, description string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO points_logs (user_id, action, points, description, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		userID, action, points, description,
	)
	if err != nil {
		return fmt.Errorf("insert points log: %w", err)
	}

	_, err = tx.Exec(`UPDATE users SET points = points + $1 WHERE id = $2`, points, userID)
	if err != nil {
		return fmt.Errorf("update user points: %w", err)
	}

	return tx.Commit()
}

func (r *PointsRepository) GetBalance(userID int64) (*model.PointsBalance, error) {
	var points int
	err := r.db.QueryRow(`SELECT COALESCE(points, 0) FROM users WHERE id = $1`, userID).Scan(&points)
	if err != nil {
		return nil, fmt.Errorf("get points: %w", err)
	}

	name, id, next, threshold := model.CalculateLevel(points)
	b := &model.PointsBalance{
		Points:  points,
		Level:   name,
		LevelID: id,
	}
	if next != "" {
		b.NextLevel = next
		b.NextThreshold = threshold
	}
	return b, nil
}

func (r *PointsRepository) ClaimDailyLogin(userID int64) (bool, error) {
	var lastLogin *time.Time
	err := r.db.QueryRow(`SELECT last_login_date FROM users WHERE id = $1`, userID).Scan(&lastLogin)
	if err != nil {
		return false, fmt.Errorf("get last login: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	if lastLogin != nil && lastLogin.Truncate(24*time.Hour).Equal(today) {
		return false, nil
	}

	err = r.AddPoints(userID, "daily_login", 1, "每日登录奖励")
	if err != nil {
		return false, err
	}

	_, err = r.db.Exec(`UPDATE users SET last_login_date = $1 WHERE id = $2`, time.Now(), userID)
	if err != nil {
		return false, fmt.Errorf("update last login date: %w", err)
	}

	return true, nil
}

func (r *PointsRepository) GetLogs(userID string, offset, limit int) ([]model.PointsLog, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM points_logs WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count logs: %w", err)
	}

	rows, err := r.db.Query(
		`SELECT id, user_id, action, points, description, created_at
		 FROM points_logs WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query logs: %w", err)
	}
	defer rows.Close()

	var logs []model.PointsLog
	for rows.Next() {
		var l model.PointsLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Points, &l.Description, &l.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *PointsRepository) GetLeaderboard(limit int) ([]model.LeaderboardEntry, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.nickname, u.school, u.avatar, u.points
		 FROM users u ORDER BY u.points DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []model.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e model.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Nickname, &e.School, &e.Avatar, &e.Points); err != nil {
			return nil, fmt.Errorf("scan leaderboard: %w", err)
		}
		name, _, _, _ := model.CalculateLevel(e.Points)
		e.Level = name
		e.Rank = rank
		entries = append(entries, e)
		rank++
	}
	return entries, nil
}

func (r *PointsRepository) AwardPointsIfNotExists(userID int64, action string, points int, description string, targetID string) error {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM points_logs WHERE user_id = $1 AND action = $2 AND description = $3`,
		userID, action, description,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check points log: %w", err)
	}
	if count > 0 {
		return nil
	}
	return r.AddPoints(userID, action, points, description)
}
