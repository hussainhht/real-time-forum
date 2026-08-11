package handlers

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/utilities"
	"realtime/backend/ws"
)

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"authenticated": false,
		})
		return
	}

	user, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"authenticated": false,
		})
		return
	}

	utilities.WriteJSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"user": map[string]any{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

func GetChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "you must login first")
		return
	}

	currentUser, err := queries.GetUserBySession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid session")
		return
	}

	users, err := queries.GetUsersForChat(currentUser.ID)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not load users")
		return
	}

	for i := range users {
		users[i].Online = ws.GlobalHub.IsOnline(users[i].ID)
	}

	if err := utilities.WriteJSON(w, http.StatusOK, users); err != nil {
		log.Println("could not encode chat users:", err)
	}
}