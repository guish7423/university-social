package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/guish/university-social/internal/model"
	"github.com/guish/university-social/internal/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsClient struct {
	hub    *WsHub
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

type WsHub struct {
	mu         sync.RWMutex
	clients    map[int64]*WsClient
	msgRepo    *repository.MessageRepository
}

func NewWsHub(msgRepo *repository.MessageRepository) *WsHub {
	return &WsHub{
		clients: make(map[int64]*WsClient),
		msgRepo: msgRepo,
	}
}

func (h *WsHub) Register(client *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.userID] = client
}

func (h *WsHub) Unregister(client *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[client.userID] == client {
		delete(h.clients, client.userID)
	}
}

func (h *WsHub) SendToUser(userID int64, data []byte) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		select {
		case client.send <- data:
		default:
		}
	}
}

func (h *WsHub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *WsHub) HandleMessage(senderID int64, msg *model.WsMessage) {
	saved, err := h.msgRepo.Send(senderID, msg.ReceiverID, msg.Content, msg.MsgType)
	if err != nil {
		log.Printf("save message: %v", err)
		return
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"type": "message",
		"data": saved,
	})

	h.SendToUser(msg.ReceiverID, payload)
	h.SendToUser(senderID, payload)
}

func (h *WsHub) HandleTyping(senderID, receiverID int64) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type":      "typing",
		"sender_id": senderID,
	})

	h.SendToUser(receiverID, payload)
}

func (h *WsHub) HandleClient(ws *websocket.Conn, userID int64) {
	client := &WsClient{
		hub:    h,
		userID: userID,
		conn:   ws,
		send:   make(chan []byte, 64),
	}
	h.Register(client)
	defer h.Unregister(client)

	go client.writePump()

	payload, _ := json.Marshal(map[string]string{"type": "connected"})
	client.send <- payload

	client.readPump()
}

func (c *WsClient) readPump() {
	defer c.conn.Close()
	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg model.WsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		msg.SenderID = c.userID

		switch msg.Type {
		case "message":
			c.hub.HandleMessage(c.userID, &msg)
		case "typing":
			c.hub.HandleTyping(c.userID, msg.ReceiverID)
		case "ping":
			c.send <- []byte(`{"type":"pong"}`)
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case data, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func WsHandler(hub *WsHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("ws upgrade: %v", err)
			return
		}
		hub.HandleClient(ws, userID)
	}
}
