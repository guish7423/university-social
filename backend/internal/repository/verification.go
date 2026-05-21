package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) CreateCode(userID int64, email, school, code string) error {
	_, err := r.db.Exec(
		`UPDATE verification_codes SET status = 'expired' WHERE user_id = $1 AND status = 'pending'`,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO verification_codes (user_id, email, code, school, status, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, 'pending', NOW() + INTERVAL '10 minutes', NOW())`,
		userID, email, code, school,
	)
	return err
}

func (r *VerificationRepository) VerifyCode(userID int64, code string) (*model.VerificationCode, error) {
	var v model.VerificationCode
	err := r.db.QueryRow(
		`SELECT id, user_id, email, code, school, status, expires_at, created_at
		 FROM verification_codes
		 WHERE user_id = $1 AND code = $2 AND status = 'pending' AND expires_at > NOW()
		 ORDER BY id DESC LIMIT 1`,
		userID, code,
	).Scan(&v.ID, &v.UserID, &v.Email, &v.Code, &v.School, &v.Status, &v.ExpiresAt, &v.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *VerificationRepository) MarkVerified(codeID int64) error {
	_, err := r.db.Exec(
		`UPDATE verification_codes SET status = 'verified' WHERE id = $1`,
		codeID,
	)
	return err
}

func (r *VerificationRepository) MarkVerifiedUser(userID int64, school, email string) error {
	_, err := r.db.Exec(
		`UPDATE users SET is_verified = TRUE, school = $1, email = $2, updated_at = NOW() WHERE id = $3`,
		school, email, userID,
	)
	return err
}

func (r *VerificationRepository) CheckVerified(userID int64) (bool, error) {
	var verified bool
	err := r.db.QueryRow(`SELECT is_verified FROM users WHERE id = $1`, userID).Scan(&verified)
	if err != nil {
		return false, fmt.Errorf("check verified: %w", err)
	}
	return verified, nil
}
