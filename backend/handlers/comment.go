package handlers

import (
  	"encoding/json"
	"net/http"
	"realtime/backend/global"
	"strings"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var Comment global.Comment

	err := json.NewDecoder(r.Body).Decode(&Comment)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	Comment.Content = strings.TrimSpace(Comment.Content)
	if Comment.Content == ""{
		http.Error(w,"missing required fields",http.StatusBadRequest )
		
	}

	


}
