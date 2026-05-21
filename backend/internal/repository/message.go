package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guish/university-social/internal/model"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Send(senderID, receiverID int64, content string, msgType int) (*model.Message, error) {
	var msg model.Message
	err := r.db.QueryRow(
		`INSERT INTO messages (sender_id, receiver_id, content, msg_type)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, sender_id, receiver_id, content, msg_type, created_at, read_at`,
		senderID, receiverID, content, msgType,
	).Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.MsgType, &msg.CreatedAt, &msg.ReadAt)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	return &msg, nil
}

func (r *MessageRepository) GetConversations(userID int64) ([]*model.Conversation, error) {
	rows, err := r.db.Query(
		`SELECT
			CASE WHEN m.sender_id = $1 THEN m.receiver_id ELSE m.sender_id END AS other_user_id,
			u.nickname, u.avatar,
			m.content AS last_content,
			m.created_at AS last_time,
			(SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND sender_id = other_user_id AND read_at IS NULL) AS unread_count
		FROM messages m
		JOIN users u ON u.id = CASE WHEN m.sender_id = $1 THEN m.receiver_id ELSE m.sender_id END
		WHERE m.id IN (
			SELECT MAX(id) FROM messages
			WHERE sender_id = $1 OR receiver_id = $1
			GROUP BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)
		)
		ORDER BY m.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("get conversations: %w", err)
	}
	defer rows.Close()

	var convs []*model.Conversation
	for rows.Next() {
		var c model.Conversation
		if err := rows.Scan(&c.OtherUserID, &c.OtherNickname, &c.OtherAvatar, &c.LastContent, &c.LastTime, &c.UnreadCount); err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}
		convs = append(convs, &c)
	}
	return convs, nil
}

func (r *MessageRepository) GetMessages(userID, otherUserID string, limit, offset int) ([]*model.Message, error) {
	rows, err := r.db.Query(
		`SELECT m.id, m.sender_id, m.receiver_id, m.content, m.msg_type, m.created_at, m.read_at,
			u.nickname, u.avatar
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE (m.sender_id = $1 AND m.receiver_id = $2) OR (m.sender_id = $2 AND m.receiver_id = $1)
		ORDER BY m.created_at DESC LIMIT $3 OFFSET $4`,
		userID, otherUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}
	defer rows.Close()

	var msgs []*model.Message
	for rows.Next() {
		var m model.Message
		var u model.User
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.MsgType, &m.CreatedAt, &m.ReadAt, &u.Nickname, &u.Avatar); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		m.Sender = &u
		msgs = append(msgs, &m)
	}
	return msgs, nil
}

func (r *MessageRepository) MarkRead(userID, otherUserID int64) error {
	_, err := r.db.Exec(
		`UPDATE messages SET read_at = $1 WHERE receiver_id = $2 AND sender_id = $3 AND read_at IS NULL`,
		time.Now(), userID, otherUserID)
	if err != nil {
		return fmt.Errorf("mark read: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND read_at IS NULL`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}
	return count, nil
}
