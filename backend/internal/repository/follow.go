package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Follow(followerID, followingID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO follows (follower_id, following_id, created_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (follower_id, following_id) DO NOTHING`,
		followerID, followingID,
	)
	return err
}

func (r *FollowRepository) Unfollow(followerID, followingID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM follows WHERE follower_id=$1 AND following_id=$2`,
		followerID, followingID,
	)
	return err
}

func (r *FollowRepository) IsFollowing(followerID, followingID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM follows WHERE follower_id=$1 AND following_id=$2`,
		followerID, followingID,
	).Scan(&count)
	return count > 0, err
}

func (r *FollowRepository) ListFollowing(userID int64, offset, limit int) ([]*model.Follow, error) {
	rows, err := r.db.Query(
		`SELECT f.id, f.follower_id, f.following_id, f.created_at,
		        u.id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		 FROM follows f
		 JOIN users u ON u.id = f.following_id
		 WHERE f.follower_id = $1
		 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list following: %w", err)
	}
	defer rows.Close()

	return scanFollowsWithTarget(rows)
}

func (r *FollowRepository) ListFollowers(userID int64, offset, limit int) ([]*model.Follow, error) {
	rows, err := r.db.Query(
		`SELECT f.id, f.follower_id, f.following_id, f.created_at,
		        u.id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		 FROM follows f
		 JOIN users u ON u.id = f.follower_id
		 WHERE f.following_id = $1
		 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list followers: %w", err)
	}
	defer rows.Close()

	return scanFollowsWithSource(rows)
}

func (r *FollowRepository) Counts(userID int64) (*model.FollowCounts, error) {
	var counts model.FollowCounts
	err := r.db.QueryRow(
		`SELECT
			(SELECT COUNT(*) FROM follows WHERE follower_id = $1),
			(SELECT COUNT(*) FROM follows WHERE following_id = $1)`,
		userID,
	).Scan(&counts.Following, &counts.Followers)
	if err != nil {
		return nil, err
	}
	return &counts, nil
}

func (r *FollowRepository) FollowingIDs(userID int64) ([]int64, error) {
	rows, err := r.db.Query(
		`SELECT following_id FROM follows WHERE follower_id = $1`, userID,
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

func scanFollowsWithTarget(rows *sql.Rows) ([]*model.Follow, error) {
	var follows []*model.Follow
	for rows.Next() {
		var f model.Follow
		var u model.User
		if err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.CreatedAt,
			&u.ID, &u.Nickname, &u.Avatar, &u.School, &u.IsVerified, &u.Role); err != nil {
			return nil, err
		}
		f.Following = &u
		follows = append(follows, &f)
	}
	return follows, nil
}

func scanFollowsWithSource(rows *sql.Rows) ([]*model.Follow, error) {
	var follows []*model.Follow
	for rows.Next() {
		var f model.Follow
		var u model.User
		if err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.CreatedAt,
			&u.ID, &u.Nickname, &u.Avatar, &u.School, &u.IsVerified, &u.Role); err != nil {
			return nil, err
		}
		f.Follower = &u
		follows = append(follows, &f)
	}
	return follows, nil
}
