package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest structures.LoginRequest
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
		http.Error(w, "error on genrate session", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = queries.InsertSession(sessionID, user.ID, expiresAt)
	if err != nil {
		http.Error(w, "error insert session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,                 //only on the website when it is using an http request, not accessible by JavaScript
		Path:     "/",                  //works only under the root path which starts with /
		SameSite: http.SameSiteLaxMode, // hide the session from the other tabs
	})

	w.Header().Set("Content-Type", "application/json")
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var User structures.Users

	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	User.FirstName = strings.TrimSpace(User.FirstName)
	User.LastName = strings.TrimSpace(User.LastName)
	User.PhoneNumber = strings.TrimSpace(User.PhoneNumber)
	User.Gender = strings.TrimSpace(User.Gender)
	User.Username = strings.TrimSpace(User.Username)
	User.Email = strings.TrimSpace(User.Email)

	if User.Age <= 0 ||
		User.Email == "" ||
		User.FirstName == "" ||
		User.LastName == "" ||
		User.Gender == "" ||
		User.Username == "" ||
		User.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	if len(User.Password) < 8 {
		http.Error(w, "the password can't be less then 8 characters ", http.StatusBadRequest)
		return
	}

	if strings.ContainsAny(User.Password, " \t\n\r") {
		http.Error(w, "password cannot contain spaces", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "can't hash password", http.StatusInternalServerError)
		return
	}
	User.Password = string(hashedPassword)

	err = queries.InsertUser(User.Username, User.FirstName, User.LastName, User.Age, User.PhoneNumber, User.Gender, User.Email, User.Password)
	if err != nil {
		http.Error(w, "faild to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	err = queries.DeleteSession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "could not logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successful",
	})
}