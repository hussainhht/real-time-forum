package global

type Messages struct {
	ID         int64  `json:"id"`
	SenderId   int    `json:"sender_id"`
	ReceiverId int    `json:"receiver_id"`
	Content    string `json:"content"`
	// CreatedAt  time.Time `json:"created_at"`
}

type SendMessage struct {
	ReceiverId int    `json:"receiver_id"`
	Content    string `json:"content"`
}
