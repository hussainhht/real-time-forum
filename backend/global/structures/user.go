package structures

import "time"

type Users struct {
	ID          int    	`json:"id"`
	FirstName   string 	`json:"first_name"`
	LastName    string 	`json:"last_name"`
	Age         int    	`json:"age"`
	Gender      string 	`json:"gender"`
	Email       string 	`json:"email"`
	Username    string 	`json:"username"`
	Password    string 	`json:"password"`
	ConfirmPassword string 	`json:"confirm_password"`
	// CreatedAt   time.Time `json:"created_at"`
}

type LoginRequest struct {
	Identifier string 	`json:"identifier"`
	Password   string 	`json:"password"`
}

type Sessions struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChatUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	// Online   bool   `json:"online"` this need to add when we make WS for know who is online and who is ofline 
}
