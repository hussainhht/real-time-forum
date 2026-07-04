package global

import "time"

type Messages struct {
	ID         int       `json:"id"`
	SenderId   int       `json:"sender_id"`
	ReceiverId int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
