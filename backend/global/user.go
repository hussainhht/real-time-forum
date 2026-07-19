package global

import "time"

type Users struct {
	ID          int       	`json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Age         int       `json:"age"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	// CreatedAt   time.Time `json:"created_at"`
}

type LoginRequest struct {
	Identifier 	string 	`json:"identifier"`
	Password   	string 	`json:"password"`
}

type Sessions struct {
	ID         	string    	`json:"id"`
	UserID     	int       	`json:"user_id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	ExpiresAt 	time.Time 	`json:"expires_at"`
}

