package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"strings"
)

func CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")

	var Comment global.Comment

	err = json.NewDecoder(r.Body).Decode(&Comment)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	Comment.Content = strings.TrimSpace(Comment.Content)
	if Comment.Content == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)

	}

}
