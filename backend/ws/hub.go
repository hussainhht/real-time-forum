package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

// client wraps a websocket connection with its own write mutex.
// gorilla/websocket only supports one concurrent writer per connection;
// without this, SendToUser and Broadcast could write to the same conn
// from different goroutines at the same time and corrupt the frame
// stream (seen client-side as "Invalid frame header" / truncated JSON).
type client struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func (c *client) writeJSON(v any) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteJSON(v)
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int]*client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]*client),
	}
}

var GlobalHub = NewHub()

func (h *Hub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, exists := h.clients[userID]; exists {
		old.conn.Close()
	}
	h.clients[userID] = &client{conn: conn}
}

func (h *Hub) Unregister(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if current, exists := h.clients[userID]; exists && current.conn == conn {
		delete(h.clients, userID)
	}
}

func (h *Hub) IsOnline(userID int) bool {
	h.mu.RLock()         //!check
	defer h.mu.RUnlock() //?read about it RW Mutex
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) OnlineUserIDs() []int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]int, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

func (h *Hub) SendToUser(userID int, event Event) bool {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	if err := c.writeJSON(event); err != nil {
		log.Println("ws: failed to send to user", userID, ":", err)
		return false
	}
	return true
}

func (h *Hub) Broadcast(event Event, skipUserID int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, c := range h.clients {
		if userID == skipUserID {
			continue
		}
		if err := c.writeJSON(event); err != nil {
			log.Println("ws: failed to broadcast to user", userID, ":", err)
		}
	}
}