package model

import "time"

type Message struct {
	ID         int64      `json:"id"`
	SenderID   int64      `json:"sender_id"`
	ReceiverID int64      `json:"receiver_id"`
	Content    string     `json:"content"`
	MsgType    int        `json:"msg_type"`
	CreatedAt  time.Time  `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	Sender     *User      `json:"sender,omitempty"`
}

type SendMessageRequest struct {
	ReceiverID int64  `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	MsgType    int    `json:"msg_type"`
}

type Conversation struct {
	OtherUserID   int64     `json:"other_user_id"`
	OtherNickname string    `json:"other_nickname"`
	OtherAvatar   string    `json:"other_avatar"`
	LastContent   string    `json:"last_content"`
	LastTime      time.Time `json:"last_time"`
	UnreadCount   int       `json:"unread_count"`
}

type WsMessage struct {
	Type      string `json:"type"`
	SenderID  int64  `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	Content   string `json:"content,omitempty"`
	MsgType   int    `json:"msg_type,omitempty"`
	MessageID int64  `json:"message_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
