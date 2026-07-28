package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
)

func CurrentUserHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed",http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type","application/json") //this for write a json 

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
		})
		return
	}

	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user":map[string]any{
			"id": user.ID,
			"username": user.Username,
			"email": user.Email,
		},

	})

	
}