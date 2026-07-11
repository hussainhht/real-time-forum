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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//todo: take the sationId for the user and get from it the user id

	//todo : take the catgory id from the post and insart it 

	var post global.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w,"invalid request body",http.StatusBadRequest)
		return 
	}

	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	if post.Title != "" || post.Content != ""{
		http.Error(w,"missing required fields", http.StatusBadRequest)
		return
	}
	if len(post.Title) <= 100{
		http.Error(w, "maximum Title reached 100",http.StatusBadRequest)
		return
	}
	if len(post.Content) >= 10000{
		http.Error(w, "maximum Content reached 10000 " , http.StatusBadRequest)
		return
	}


	err = queries.InsertPost(post.UserID,post.Title,post.Content,post.CategoryID)
	if err != nil {
		http.Error(w,"there is a problem with database" , http.StatusInternalServerError)
		return 
	}

}
