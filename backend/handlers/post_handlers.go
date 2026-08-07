package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"realtime/backend/global/utilities"
	"strconv"
	"strings"
)

var validCategories = map[string]bool{
	"General":       true,
	"Technology":    true,
	"Sports":        true,
	"Entertainment": true,
	"Science":       true,
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

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

	var newPost structures.NewPost
	err = utilities.ReadJSON(r, &newPost)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newPost.Title = strings.TrimSpace(newPost.Title)
	newPost.Content = strings.TrimSpace(newPost.Content)
	newPost.Category = strings.TrimSpace(newPost.Category)

	if newPost.Title == "" || newPost.Content == "" || newPost.Category == "" {
		utilities.ErrorJSON(w, http.StatusBadRequest, "missing required fields")
		return
	}
	if !validCategories[newPost.Category] {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid category")
		return
	}
	if len(newPost.Title) > 100 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "maximum Title reached (100 characters)")
		return
	}
	if len(newPost.Content) >= 10000 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "maximum Content reached (10000 characters)")
		return
	}

	var post structures.Post

	post.Title = newPost.Title
	post.Content = newPost.Content
	post.Category = newPost.Category
	post.UserID = user.ID

	postID, err := queries.InsertPost(post.UserID, post.Title, post.Content, post.Category)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "post couldn't be inserted to the database")
		return
	}

	post.ID = int(postID)

	utilities.WriteJSON(w, http.StatusCreated, post)
}

func FeedPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	posts, err := queries.PostsFeed()
	if err != nil {
		log.Println("PostsFeed error:", err)
		utilities.ErrorJSON(w, http.StatusInternalServerError, "there is a problem with database")
		return
	}

	if err := utilities.WriteJSON(w, http.StatusOK, posts); err != nil {
		log.Println(err)
	}
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID <= 0 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := queries.GetPostByID(postID)
	if errors.Is(err, sql.ErrNoRows) {
		utilities.ErrorJSON(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		log.Println("could not get post:", err)
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not get post")
		return
	}

	if err := utilities.WriteJSON(w, http.StatusOK, post); err != nil {
		log.Println("could not encode post:", err)
	}
}