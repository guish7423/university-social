package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type BlockRepository struct {
	db *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

func (r *BlockRepository) Block(blockerID, blockedID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO user_blocks (blocker_id, blocked_id, created_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (blocker_id, blocked_id) DO NOTHING`,
		blockerID, blockedID,
	)
	return err
}

func (r *BlockRepository) Unblock(blockerID, blockedID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM user_blocks WHERE blocker_id=$1 AND blocked_id=$2`,
		blockerID, blockedID,
	)
	return err
}

func (r *BlockRepository) IsBlocking(blockerID, blockedID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM user_blocks WHERE blocker_id=$1 AND blocked_id=$2`,
		blockerID, blockedID,
	).Scan(&count)
	return count > 0, err
}

func (r *BlockRepository) BlockedIDs(userID int64) ([]int64, error) {
	rows, err := r.db.Query(
		`SELECT blocked_id FROM user_blocks WHERE blocker_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *BlockRepository) BlockedByIDs(userID int64) ([]int64, error) {
	rows, err := r.db.Query(
		`SELECT blocker_id FROM user_blocks WHERE blocked_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *BlockRepository) ListBlockedUsers(userID int64) ([]*model.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.open_id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		 FROM user_blocks ub
		 JOIN users u ON u.id = ub.blocked_id
		 WHERE ub.blocker_id = $1
		 ORDER BY ub.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list blocked users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.OpenID, &u.Nickname, &u.Avatar, &u.School, &u.IsVerified, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
