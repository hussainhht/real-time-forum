package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest global.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if loginRequest.Identifier == "" || loginRequest.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByUsernameOrEmail(loginRequest.Identifier)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID, err := generateSession()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = queries.InsertSession(sessionID, user.ID, expiresAt)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: sessionID,
		Expires: expiresAt,
		HttpOnly: true, //only on the website when it is using an http request, not accessible by JavaScript
		Path: "/", //works only under the root path which starts with /
		SameSite: http.SameSiteLaxMode, // hide the session from the other tabs
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"username": user.Username})
}

func generateSession() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
