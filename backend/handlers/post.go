package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"strconv"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
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

	var newPost structures.NewPost
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

	var post structures.Post

	post.Title = newPost.Title
	post.Content = newPost.Content
	post.UserID = user.ID

	postID, err := queries.InsertPost(post.UserID, post.Title, post.Content)
	if err != nil {
		http.Error(w, "post couldn't be inserted to the database", http.StatusInternalServerError)
		return
	}

	post.ID = int(postID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func FeedPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	_, err = queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	posts, err := queries.GetPostFeed()
	if err != nil {
		http.Error(w, "there is a problem with database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Println(err)
	}
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	_, err = queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID <= 0 {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	post, err := queries.GetPostByID(postID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("could not get post:", err)
		http.Error(w, "could not get post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Println("could not encode post:", err)
	}
}