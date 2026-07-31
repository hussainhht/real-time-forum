package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
)

func LikeHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Methode not allowed", http.StatusMethodNotAllowed)
		return

	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you must login", http.StatusUnauthorized)
		return
	}

	UserId, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var like global.Like
	err = json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, "invbaled request body", http.StatusBadRequest)
		return
	}
	if like.PostID <= 0 {
		http.Error(w, "invaled post id", http.StatusBadRequest)
		return
	}

	like.UserID = UserId.ID
	err = queries.InsertLike(like.PostID, like.UserID)
	if err != nil {
		http.Error(w, "could not insert like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]any{
		"massge": "line added successfully",
		"liked":  true,
	})

}
