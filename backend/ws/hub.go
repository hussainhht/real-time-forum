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

type Hub struct {
	mu      sync.RWMutex
	clients map[int]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]*websocket.Conn),
	}
}

var GlobalHub = NewHub()

func (h *Hub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, exists := h.clients[userID]; exists {
		old.Close()
	}
	h.clients[userID] = conn
}

func (h *Hub) Unregister(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if current, exists := h.clients[userID]; exists && current == conn {
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
	conn, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	if err := conn.WriteJSON(event); err != nil {
		log.Println("ws: failed to send to user", userID, ":", err)
		return false
	}
	return true
}

func (h *Hub) Broadcast(event Event, skipUserID int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, conn := range h.clients {
		if userID == skipUserID {
			continue
		}
		if err := conn.WriteJSON(event); err != nil {
			log.Println("ws: failed to broadcast to user", userID, ":", err)
		}
	}
}
