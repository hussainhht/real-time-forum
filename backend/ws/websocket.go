package ws

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/utilities"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// CheckOrigin left at default (same-origin only) since this is
	// served from the same host as the frontend.
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate BEFORE upgrading, same as the HTTP handlers.
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you must login first")
		return
	}

	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid session")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ws] Failed to upgrade to websocket:", err)
		return
	}
	defer conn.Close()

	GlobalHub.Register(user.ID, conn)
	log.Printf("[ws] REGISTERED user %s (id %d) | online now: %v", user.Username, user.ID, GlobalHub.OnlineUserIDs())

	GlobalHub.Broadcast(Event{
		Type:    "user_online",
		Content: map[string]any{"user_id": user.ID, "username": user.Username},
	}, user.ID)
	log.Printf("[ws] BROADCAST user_online for %s (id %d)", user.Username, user.ID)

	defer func() {
		GlobalHub.Unregister(user.ID, conn)
		log.Printf("[ws] UNREGISTERED user %s (id %d) | online now: %v", user.Username, user.ID, GlobalHub.OnlineUserIDs())
		GlobalHub.Broadcast(Event{
			Type:    "user_offline",
			Content: map[string]any{"user_id": user.ID, "username": user.Username},
		}, user.ID)
		log.Printf("[ws] BROADCAST user_offline for %s (id %d)", user.Username, user.ID)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[ws] ReadMessage error for user %s (id %d): %v", user.Username, user.ID, err)
			break
		}
	}
}