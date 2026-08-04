package handlers

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"realtime/backend/global/utilities"
	"strconv"
	"strings"
	"time"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// w.Header().Set("Allow", http.MethodPost)		//! i dont think we need this, because we are already sending the error message in the next line
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

	messageID, err := queries.InsertMessage(sender.ID, sendMessage.ReceiverId, sendMessage.Content, time.Now())
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not save message")
		return
	}

	response := structures.Message{
		ID:         messageID,
		SenderId:   sender.ID,
		ReceiverId: sendMessage.ReceiverId,
		Content:    sendMessage.Content,
		CreatedAt:  time.Now(),
	}

	err = utilities.WriteJSON(w, http.StatusCreated, response)
	if err != nil {
		log.Println("could not encode message:", err)		//? is this true to be only in the terminal?
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

	beforeID := 0 //? last message loaded, if 0 then load the latest messages
	beforeValue := r.URL.Query().Get("before") //?  /messages/5?before=41			//! also check how is this exactly working, READ ABOUT IT

	if beforeValue != "" { //! when will this stop?
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
		log.Println("could not encode messages:", err) //! only in terminal?
	}
}