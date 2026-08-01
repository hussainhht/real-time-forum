package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"strconv"
	"strings"
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
		http.Error(w, "Method request body", http.StatusBadRequest)
		return

	}

	sendMessage.Content = strings.TrimSpace(sendMessage.Content)

	if sendMessage.ReceiverId <= 0 {
		http.Error(w, "Method receiver id", http.StatusBadRequest)
		return

	}
	if sendMessage.ReceiverId == sender.ID {
		http.Error(w, "you cant massage your self", http.StatusBadRequest)
		return

	}
	if sendMessage.Content == "" || len(sendMessage.Content) > 1000 {
		http.Error(w, "message content not valid", http.StatusBadRequest)
		return

	}

	messageID, err := queries.InsertMessage(sender.ID, sendMessage.ReceiverId, sendMessage.Content)
	if err != nil {
		http.Error(w, "could not save message", http.StatusInternalServerError)
		return

	}

	response := structures.Messages{
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

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		//todo : make the json error
		return
	}

	currentUser, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		//todo add error by json
		return
	}

	otherUserId, err := strconv.Atoi(r.PathValue("userId")) //this come from the path 
	if err != nil {
		//todo : add jeson error
		return
	}

	if otherUserId == currentUser.ID { //canot be the same user 
		// todo : json error
		return

	}

	beforeID := 0
	beforeValue := r.URL.Query().Get("before") //Query String  in the url  /api/messages/5?before=41&limit=10

	if beforeValue != "" {
		beforeID, err := strconv.Atoi(beforeValue)
		if err != nil || beforeID <= 0 { //invalid before message id
			//todo : write json Error
			return
		}

	}

	//get the message from DB
	messages, err := queries.GetMessagesBetweenUsers(
		currentUser.ID,
		otherUserId,
		beforeID,
	)
	if err != nil {
		//todo : write json Error
		return
	}


	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println("could not encode message :", err)
	}


}
