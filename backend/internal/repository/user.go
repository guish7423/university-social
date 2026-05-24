package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByOpenID(openID string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, open_id, union_id, nickname, avatar, phone, school, is_verified, role, created_at, updated_at FROM users WHERE open_id = $1",
		openID,
	).Scan(&u.ID, &u.OpenID, &u.UnionID, &u.Nickname, &u.Avatar, &u.Phone, &u.School, &u.IsVerified, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find by openid: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Create(u *model.User) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (open_id, union_id, nickname, avatar, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		 RETURNING id`,
		u.OpenID, u.UnionID, u.Nickname, u.Avatar, "user",
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, open_id, union_id, nickname, avatar, phone, school, is_verified, role, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.OpenID, &u.UnionID, &u.Nickname, &u.Avatar, &u.Phone, &u.School, &u.IsVerified, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) UpdateProfile(id int64, nickname, avatar, school string) error {
	_, err := r.db.Exec(
		"UPDATE users SET nickname=$1, avatar=$2, school=$3, updated_at=NOW() WHERE id=$4",
		nickname, avatar, school, id,
	)
	return err
}

func (r *UserRepository) GetPasswordHash(id int64) (string, error) {
	var hash string
	err := r.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", id).Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("get password hash: %w", err)
	}
	return hash, nil
}

func (r *UserRepository) UpdatePassword(id int64, hash string) error {
	_, err := r.db.Exec("UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2", hash, id)
	return err
}
