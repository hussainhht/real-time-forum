package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"strconv"
	"strings"
)

func CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //! check router
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

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	postExists, err := queries.PostExists(postID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !postExists {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	var Comment global.Comment

	var newComment global.NewComment
	err = json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	Comment.Content = strings.TrimSpace(Comment.Content)
	if Comment.Content == "" {
		http.Error(w, "missing comment content", http.StatusBadRequest)
	}
	if len(newComment.Content) >= 10000 {
		http.Error(w, "maximum Content reached 10000", http.StatusBadRequest)
		return
	}

	commentID, err := queries.InsertComment(postID, user.ID, newComment.Content)
	if err != nil {
		http.Error(w, "there is a problem with database", http.StatusInternalServerError)
		return
	}

	comment := global.Comment{
		ID:      int(commentID),
		PostID:  postID,
		UserID:  user.ID,
		Content: newComment.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)

}
