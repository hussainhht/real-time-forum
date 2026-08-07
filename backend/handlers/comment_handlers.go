package handlers

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"realtime/backend/global/utilities"
	"strconv"
	"strings"
)

func CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	postExists, err := queries.PostExists(postID)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !postExists {
		utilities.ErrorJSON(w, http.StatusNotFound, "post not found")
		return
	}

	var newComment structures.NewComment
	err = utilities.ReadJSON(r, &newComment)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newComment.Content = strings.TrimSpace(newComment.Content)
	if newComment.Content == "" {
		utilities.ErrorJSON(w, http.StatusBadRequest, "missing comment content")
		return
	}
	if len(newComment.Content) >= 10000 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "maximum Content reached 10000")
		return
	}

	commentID, err := queries.InsertComment(postID, user.ID, newComment.Content)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "there is a problem with database")
		return
	}

	comment := structures.Comment{
		ID:      int(commentID),
		PostID:  postID,
		UserID:  user.ID,
		Content: newComment.Content,
	}

	utilities.WriteJSON(w, http.StatusCreated, comment)
}

func FeedCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// w.Header().Set("Allow", http.MethodGet)
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	_, err = queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id")) 		//! check how is this exactly working, READ ABOUT IT
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	comments, err := queries.FeedCommentsByPostID(postID)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not load comments")
		return
	}

	err = utilities.WriteJSON(w, http.StatusOK, comments)
	if err != nil {
		log.Println(err)
	}
}