package handlers

import (
	"log"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global/structures"
	"realtime/backend/global/utilities"
	"realtime/backend/ws"
	"regexp"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var loginRequest structures.LoginRequest
	err := utilities.ReadJSON(r, &loginRequest) // decoding
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if loginRequest.Identifier == "" || loginRequest.Password == "" {
		utilities.ErrorJSON(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, err := queries.GetUserByUsernameOrEmail(loginRequest.Identifier)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	err = queries.DeleteSessionsByUserID(user.ID)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "error clearing old sessions")
		return
	}

	sessionID, err := generateSession()
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "error on generating the session")
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = queries.InsertSession(sessionID, user.ID, expiresAt)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "error insert session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true, //?! this do not allow the cookie to be accessed by JavaScript, which helps protect against cross-site scripting (XSS) attacks.
		Path:     "/",
		SameSite: http.SameSiteLaxMode, //?! this helps protect against cross-site request forgery (CSRF) attacks by restricting how cookies are sent with cross-site requests.
	})

	utilities.WriteJSON(w, http.StatusOK, map[string]string{"username": user.Username})
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
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var User structures.Users

	err := utilities.ReadJSON(r, &User)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	User.FirstName = strings.TrimSpace(User.FirstName)
	User.LastName = strings.TrimSpace(User.LastName)
	User.Gender = strings.TrimSpace(User.Gender)
	User.Username = strings.TrimSpace(User.Username)
	User.Email = strings.TrimSpace(User.Email)
	User.ConfirmPassword = strings.TrimSpace(User.ConfirmPassword)

	if User.Age <= 0 ||
		User.Email == "" ||
		User.FirstName == "" ||
		User.LastName == "" ||
		User.Gender == "" ||
		User.Username == "" ||
		User.Password == "" ||
		User.ConfirmPassword == "" {
		utilities.ErrorJSON(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if len(User.Password) < 8 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "the password can't be less then 8 characters")
		return
	}

	if User.Password != User.ConfirmPassword {
		utilities.ErrorJSON(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	if strings.ContainsAny(User.Password, " \t\n\r") {
		utilities.ErrorJSON(w, http.StatusBadRequest, "password cannot contain spaces")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "can't hash password")
		return
	}
	User.Password = string(hashedPassword)

	if strings.Contains(User.Username, " ") {
		utilities.ErrorJSON(w, http.StatusBadRequest, "Name cant have space ")
		return
	}

	if len(User.FirstName) > 20 || len(User.LastName) > 20 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "Name too long")
		return
	}

	if User.Gender != "male" && User.Gender != "female" {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid gender")
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if len(User.Email) < 3 || len(User.Email) > 254 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid email")
		return
	}
	if !emailRegex.MatchString(User.Email) {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid email")
		return
	}

	if User.Age < 0 || User.Age > 100 {
		utilities.ErrorJSON(w, http.StatusBadRequest, "invalid age")
		return
	}

	err = queries.InsertUser(User.Username, User.FirstName, User.LastName, User.Age, User.Gender, User.Email, User.Password)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "faild to create user")
		log.Println("Error inserting user:", err) // Log the error for debugging purposes
		return
	}

	ws.GlobalHub.Broadcast(ws.Event{
		Type: "new_user_registered",
		Content: map[string]any{
			"username":   User.Username,
			"first_name": User.FirstName,
			"last_name":  User.LastName,
		},
	}, -1)

	utilities.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user registered successfully"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.ErrorJSON(w, http.StatusUnauthorized, "not logged in")
		return
	}

	err = queries.DeleteSession(sessionCookie.Value)
	if err != nil {
		utilities.ErrorJSON(w, http.StatusInternalServerError, "could not logout")
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

	utilities.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}
