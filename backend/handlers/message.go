package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"strings"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return

	}

	ssessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you must login first", http.StatusUnauthorized)
		return
	}

	sender, err := queries.GetUserBySession(ssessionCookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return

	}

	var sendMessage global.SendMessage

	err = json.NewDecoder(r.Body).Decode(&sendMessage)
	if err != nil {
		http.Error(w, "invaled request body", http.StatusBadRequest)
		return

	}

	sendMessage.Content = strings.TrimSpace(sendMessage.Content)

	if sendMessage.ReceiverId <= 0 {
		http.Error(w, "invaled receiver id", http.StatusBadRequest)
		return

	}
	if sendMessage.ReceiverId == sender.ID {
		http.Error(w, "you cant massage your self", http.StatusBadRequest)
		return

	}
	if sendMessage.Content == "" || len(sendMessage.Content) > 1000 {
		http.Error(w, "message content not valed", http.StatusBadRequest)
		return

	}

	messageID, err := queries.InsertMessage(sender.ID, sendMessage.ReceiverId, sendMessage.Content)
	if err != nil {
		http.Error(w, "cold not save messge", http.StatusInternalServerError)
		return

	}

	response := global.Messages{
		ID:         messageID,
		SenderId:   sender.ID,
		ReceiverId: sendMessage.ReceiverId,
		Content:    sendMessage.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return

	}

}
