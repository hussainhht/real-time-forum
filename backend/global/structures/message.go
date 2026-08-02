package structures

import "time"

// type Messages struct { //for a chat list
// 	ID         int64  `json:"id"`
// 	SenderId   int    `json:"sender_id"`
// 	ReceiverId int    `json:"receiver_id"`
// 	Content    string `json:"content"`
// }

type Message struct { //a single message
	ID         int64  `json:"id"`
	SenderId   int    `json:"sender_id"`
	ReceiverId int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type SendMessage struct {
	ReceiverId int    `json:"receiver_id"`
	Content    string `json:"content"`
}
