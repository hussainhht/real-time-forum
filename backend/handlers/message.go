package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"strconv"
	"strings"
	"time"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you must login first", http.StatusUnauthorized)
		return
	}

	sender, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var sendMessage structures.SendMessage

	err = json.NewDecoder(r.Body).Decode(&sendMessage)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sendMessage.Content = strings.TrimSpace(sendMessage.Content)

	if sendMessage.ReceiverId <= 0 {
		http.Error(w, "invalid receiver id", http.StatusBadRequest)
		return
	}
	if sendMessage.ReceiverId == sender.ID {
		http.Error(w, "you cant message your self", http.StatusBadRequest)
		return
	}
	if sendMessage.Content == "" || len(sendMessage.Content) > 1000 {
		http.Error(w, "message content not valid", http.StatusBadRequest)
		return
	}

	messageID, err := queries.InsertMessage(sender.ID, sendMessage.ReceiverId, sendMessage.Content, time.Now())
	if err != nil {
		http.Error(w, "could not save message", http.StatusInternalServerError)
		return
	}

	response := structures.Message{
		ID:         messageID,
		SenderId:   sender.ID,
		ReceiverId: sendMessage.ReceiverId,
		Content:    sendMessage.Content,
		CreatedAt:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("could not encode message:", err)
	}
}

func FeedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you must login first", http.StatusUnauthorized)
		return
	}

	currentUser, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	otherUserId, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if otherUserId == currentUser.ID {
		http.Error(w, "you cant message your self", http.StatusBadRequest)
		return
	}

	beforeID := 0
	beforeValue := r.URL.Query().Get("before") //?  /messages/5?before=41

	if beforeValue != "" {
		parsed, err := strconv.Atoi(beforeValue)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid before value", http.StatusBadRequest)
			return
		}
		beforeID = parsed
	}

	messages, err := queries.FeedMessagesBetweenUsers(
		currentUser.ID,
		otherUserId,
		beforeID,
	)
	if err != nil {
		http.Error(w, "could not load messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println("could not encode messages:", err)
	}
}
