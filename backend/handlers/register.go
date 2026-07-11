package handlers

import (
	"encoding/json"
	"net/http"
	"realtime/backend/db/queries"
	"realtime/backend/global"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var User global.Users

	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		http.Error(w, "invaled request body", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)              //201
	w.Write([]byte("user registered successfully")) //this will go to the page

}
