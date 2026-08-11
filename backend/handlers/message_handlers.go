package handlers

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"realtime/backend/global/utilities"
	"realtime/backend/ws"
	"strconv"
	"strings"
	"time"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you must login first")
		return
	}

	sender, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid session")
		return
	}

	var sendMessage structures.SendMessage

	err = utilities.ReadJSON(r, &sendMessage)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sendMessage.Content = strings.TrimSpace(sendMessage.Content)

	if sendMessage.ReceiverId <= 0 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid receiver id")
		return
	}
	if sendMessage.ReceiverId == sender.ID {
		utilities.ErrorJSON(w, http.StatusBadRequest, "you cant message your self")
		return
	}
	if sendMessage.Content == "" || len(sendMessage.Content) > 1000 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "message content not valid")
		return
	}

	now := time.Now()

	messageID, err := queries.InsertMessage(sender.ID, sendMessage.ReceiverId, sendMessage.Content, now)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not save message")
		return
	}

	response := structures.Message{
		ID:         messageID,
		SenderId:   sender.ID,
		ReceiverId: sendMessage.ReceiverId,
		Username:   sender.Username,
		Content:    sendMessage.Content,
		CreatedAt:  now,
	}

	ws.GlobalHub.SendToUser(sendMessage.ReceiverId, ws.Event{
		Type:    "new_message",
		Content: response,
	})

	err = utilities.WriteJSON(w, http.StatusCreated, response)
	if err != nil {
		log.Println("could not encode message:", err)
	}
}

func FeedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you must login first")
		return
	}

	currentUser, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid session")
		return
	}

	otherUserId, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if otherUserId == currentUser.ID {
		utilities.ErrorJSON(w, http.StatusBadRequest, "you cant message your self")
		return
	}

	beforeID := 0
	beforeValue := r.URL.Query().Get("before")

	if beforeValue != "" {
		parsed, err := strconv.Atoi(beforeValue)
		if err != nil || parsed <= 0 {
			utilities.ErrorJSON(w, http.StatusBadRequest, "invalid before value")
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
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not load messages")
		return
	}

	if err := utilities.WriteJSON(w, http.StatusOK, messages); err != nil {
		log.Println("could not encode messages:", err)
	}
}
