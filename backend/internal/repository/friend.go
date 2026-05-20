package repository

import (
	"database/sql"
	"fmt"

	"github.com/guish/university-social/internal/model"
)
type UserSearchResult struct {
	model.User
	FriendStatus string `json:"friend_status"`
}

type FriendRepository struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) SendRequest(userID, friendID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO friends (user_id, friend_id, status, created_at, updated_at)
		 VALUES ($1, $2, 'pending', NOW(), NOW())
		 ON CONFLICT (user_id, friend_id) DO UPDATE SET updated_at = NOW()`,
		userID, friendID,
	)
	return err
}

func (r *FriendRepository) AcceptRequest(userID, friendID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`UPDATE friends SET status='accepted', updated_at=NOW()
		 WHERE user_id=$1 AND friend_id=$2 AND status='pending'`,
		friendID, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("request not found")
	}

	_, err = tx.Exec(
		`INSERT INTO friends (user_id, friend_id, status, created_at, updated_at)
		 VALUES ($1, $2, 'accepted', NOW(), NOW())`,
		userID, friendID,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *FriendRepository) Remove(userID, friendID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM friends WHERE (user_id=$1 AND friend_id=$2) OR (user_id=$2 AND friend_id=$1)",
		userID, friendID,
	)
	return err
}

func (r *FriendRepository) RejectRequest(userID, friendID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM friends WHERE user_id=$1 AND friend_id=$2 AND status='pending'`,
		friendID, userID,
	)
	return err
}
func (r *FriendRepository) ListFriends(userID int64) ([]*model.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.open_id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		FROM friends f JOIN users u ON u.id = f.friend_id
		WHERE f.user_id = $1 AND f.status = 'accepted'
		UNION
		SELECT u.id, u.open_id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		FROM friends f JOIN users u ON u.id = f.user_id
		WHERE f.friend_id = $1 AND f.status = 'accepted'`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
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

func (r *FriendRepository) ListRequests(userID int64) ([]*model.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.open_id, u.nickname, u.avatar, u.school, u.is_verified, u.role
		FROM friends f JOIN users u ON u.id = f.user_id
		WHERE f.friend_id = $1 AND f.status = 'pending'`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list requests: %w", err)
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

func (r *FriendRepository) SearchUsers(query string, limit int, currentUserID int64) ([]*UserSearchResult, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.open_id, u.nickname, u.avatar, u.school, u.is_verified, u.role,
		CASE
			WHEN f1.status = 'accepted' THEN 'friend'
			WHEN f2.id IS NOT NULL THEN 'pending_sent'
			WHEN f3.id IS NOT NULL THEN 'pending_received'
			ELSE ''
		END as friend_status
		FROM users u
		LEFT JOIN friends f1 ON ((f1.user_id = $2 AND f1.friend_id = u.id) OR (f1.user_id = u.id AND f1.friend_id = $2)) AND f1.status = 'accepted'
		LEFT JOIN friends f2 ON f2.user_id = $2 AND f2.friend_id = u.id AND f2.status = 'pending'
		LEFT JOIN friends f3 ON f3.user_id = u.id AND f3.friend_id = $2 AND f3.status = 'pending'
		WHERE u.id != $2 AND u.nickname ILIKE $1 LIMIT $3
	`, "%"+query+"%", currentUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}
	defer rows.Close()

	var users []*UserSearchResult
	for rows.Next() {
		var r UserSearchResult
		if err := rows.Scan(&r.ID, &r.OpenID, &r.Nickname, &r.Avatar, &r.School, &r.IsVerified, &r.Role, &r.FriendStatus); err != nil {
			return nil, err
		}
		users = append(users, &r)
	}
	return users, nil
}

func (r *FriendRepository) IsFriend(userID, friendID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM friends
		 WHERE ((user_id=$1 AND friend_id=$2) OR (user_id=$2 AND friend_id=$1))
		 AND status='accepted'`,
		userID, friendID,
	).Scan(&count)
	return count > 0, err
}

func (r *FriendRepository) CreateNotification(userID, fromUserID int64, ntype, content string, refID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO notifications (user_id, from_user_id, type, content, ref_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW())`,
		userID, fromUserID, ntype, content, refID,
	)
	return err
}

func (r *FriendRepository) ListNotifications(userID int64, limit int) ([]*model.Notification, error) {
	rows, err := r.db.Query(
		`SELECT n.id, n.user_id, n.from_user_id, n.type, n.content, n.ref_id, n.is_read, n.created_at,
			u.nickname, u.avatar
		FROM notifications n JOIN users u ON u.id = n.from_user_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notes []*model.Notification
	for rows.Next() {
		var n model.Notification
		var fromUser model.User
		if err := rows.Scan(&n.ID, &n.UserID, &n.FromUserID, &n.Type, &n.Content, &n.RefID, &n.IsRead, &n.CreatedAt,
			&fromUser.Nickname, &fromUser.Avatar); err != nil {
			return nil, err
		}
		n.FromUser = &fromUser
		notes = append(notes, &n)
	}
	return notes, nil
}

func (r *FriendRepository) MarkNotificationsRead(userID int64) error {
	_, err := r.db.Exec("UPDATE notifications SET is_read=true WHERE user_id=$1", userID)
	return err
}
