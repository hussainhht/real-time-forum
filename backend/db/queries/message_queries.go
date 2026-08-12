package queries

import (
	"realtime/backend/global"
	"realtime/backend/global/structures"
	"time"
)

const insert_message = `INSERT INTO messages (
	sender_id,
	receiver_id,
	content,
	created_at
) VALUES (?,?,?,?)
`

func InsertMessage(senderID, receiverID int, content string, createdAt time.Time) (int64, error) {
	result, err := global.Database.Exec(
		insert_message,
		senderID,
		receiverID,
		content,
		createdAt,
	)
	if err != nil {
		return 0, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

const getMessageBetweenUsers = `
SELECT
	messages.id,
	messages.sender_id,
	messages.receiver_id,
	users.username,
	messages.content,
	messages.created_at
	FROM messages
	JOIN users ON users.id = messages.sender_id
	WHERE (
		(messages.sender_id = ? AND messages.receiver_id = ?)
		OR
		(messages.sender_id = ? AND messages.receiver_id = ?)
	)
	AND (? = 0 OR messages.id < ?)
	ORDER BY messages.id DESC
	LIMIT 10
`

func FeedMessagesBetweenUsers(currentUserID, otherUserID, beforeID int) ( //before is the last inserted message id
	[]structures.Message, error,
) {
	rows, err := global.Database.Query(
		getMessageBetweenUsers,
		currentUserID,
		otherUserID,

		otherUserID,
		currentUserID,
		beforeID,
		beforeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []structures.Message{}

	for rows.Next() {
		var message structures.Message
		err := rows.Scan(
			&message.ID,
			&message.SenderId,
			&message.ReceiverId,
			&message.Username,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// to make the message from old to new
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

const get_username_by_id = `SELECT username FROM users WHERE id = ?`

func GetUsernameByID(userID int) (string, error) {
	var username string
	err := global.Database.QueryRow(get_username_by_id, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}