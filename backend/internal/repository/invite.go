package repository

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/guish/university-social/internal/model"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type InviteCodeRepository struct {
	db *sql.DB
}

func NewInviteCodeRepository(db *sql.DB) *InviteCodeRepository {
	return &InviteCodeRepository{db: db}
}

func generateCode(n int) (string, error) {
	code := make([]byte, n)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[idx.Int64()]
	}
	return string(code), nil
}

func (r *InviteCodeRepository) Create(userID int64, count int) ([]model.InviteCode, error) {
	codes := make([]model.InviteCode, 0, count)
	for i := 0; i < count; i++ {
		codeStr, err := generateCode(8)
		if err != nil {
			return nil, fmt.Errorf("generate code: %w", err)
		}
		var c model.InviteCode
		err = r.db.QueryRow(
			`INSERT INTO invite_codes (code, creator_id, status, created_at)
			 VALUES ($1, $2, 'active', NOW()) RETURNING id, code, creator_id, status, created_at`,
			codeStr, userID,
		).Scan(&c.ID, &c.Code, &c.CreatorID, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("create invite code: %w", err)
		}
		codes = append(codes, c)
	}
	return codes, nil
}

func (r *InviteCodeRepository) ListByCreator(userID int64) ([]model.InviteCode, error) {
	rows, err := r.db.Query(
		`SELECT id, code, creator_id, invited_user_id, status, created_at, used_at
		 FROM invite_codes WHERE creator_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []model.InviteCode
	for rows.Next() {
		var c model.InviteCode
		err := rows.Scan(&c.ID, &c.Code, &c.CreatorID, &c.InvitedUserID, &c.Status, &c.CreatedAt, &c.UsedAt)
		if err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, nil
}

func (r *InviteCodeRepository) FindByCode(code string) (*model.InviteCode, error) {
	var c model.InviteCode
	err := r.db.QueryRow(
		`SELECT id, code, creator_id, invited_user_id, status, created_at, used_at
		 FROM invite_codes WHERE code = $1`, code,
	).Scan(&c.ID, &c.Code, &c.CreatorID, &c.InvitedUserID, &c.Status, &c.CreatedAt, &c.UsedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *InviteCodeRepository) Redeem(codeID int64, invitedUserID int64) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE invite_codes SET status = 'used', invited_user_id = $1, used_at = $2 WHERE id = $3`,
		invitedUserID, now, codeID,
	)
	return err
}

func (r *InviteCodeRepository) CreatorStats(userID int64) (total int, used int, err error) {
	err = r.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(CASE WHEN status = 'used' THEN 1 ELSE 0 END), 0)
		 FROM invite_codes WHERE creator_id = $1`, userID,
	).Scan(&total, &used)
	return
}
