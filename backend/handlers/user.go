package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
)

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json") //this for write a json

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false, // * we use the json because we want to return a bool value like authenticated: false when i use http.Error it return only text i can't control what i send in the body >> insted if i use jsone and w.writeheder i can control what i send in the body like bool value or intger value even and opject.
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
	w.WriteHeader(http.StatusOK) //200  
	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}

func GetChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // cahnge this to json later
		return
	}
	w.Header().Set("Content-Type", "application/json")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "you must login first",
		})
		return
	}

	currentUser, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {

		//todo: write the error json
		return
	}
	users, err := queries.GetUsersForChat(currentUser.ID)
	if err != nil {
		//todo : write the json error
		return
	}

	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		//todo: write json error without return

	}

}
