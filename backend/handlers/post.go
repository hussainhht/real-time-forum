package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}
	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}


	var newPost global.NewPost
	err = json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newPost.Title = strings.TrimSpace(newPost.Title)
	newPost.Content = strings.TrimSpace(newPost.Content)

	if newPost.Title == "" || newPost.Content == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}
	if len(newPost.Title) > 100 {
		http.Error(w, "maximum Title reached (100 characters)", http.StatusBadRequest)
		return
	}
	if len(newPost.Content) >= 10000 {
		http.Error(w, "maximum Content reached (10000 characters)", http.StatusBadRequest)
		return
	}
	if newPost.CategoryID <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	exists, err := queries.CategoryExists(newPost.CategoryID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}


	var post global.Post

	post.Title = newPost.Title
	post.Content = newPost.Content
	post.UserID = user.ID
	post.CategoryID = newPost.CategoryID

	postID, err := queries.InsertPost(post.UserID, post.Title, post.Content, post.CategoryID)
	if err != nil {
		http.Error(w, "post couldn't be inserted to the database", http.StatusInternalServerError)
		return
	}

	post.ID = int(postID)
	w.WriteHeader(http.StatusCreated)	//201
	json.NewEncoder(w).Encode(post)		//? do not forget to be sure about what these are doing..
}
