package queries

import (
	"realtime/backend/global"
	"realtime/backend/global/structures"
)

const Insert_messge = `INSERT INTO messages (
	sender_id,
	receiver_id,
	content

) VALUES (?,?,?)
`

func InsertMessage(senderID, receiverID int, content string) (int64, error) {
	result, err := global.Database.Exec(
		Insert_messge,
		senderID,
		receiverID,
		content,
	)
	if err != nil {
		return 0, err
	}

	messgeID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return messgeID, nil
}

const getMessageBetweenUsers = `
SELECT
	id,
	sender_id,
	reciver_id,
	content,
	created_at
	FROM messages
	WHERE (
		(sender_id = ? AND receiver_id = ?)
		OR	 
		(sender_id = ? AND receiver_id = ?) 
	)
		AND (? = 0 OR id < ?)
		ORDER BY id DESC 
		LIMIT 10
		`

func GetMessagesBetweenUsers(currentUserID, otherUserID, beforeID int) (
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

	for rows.Next(){
		var message structures.Message
		err := rows.Scan(
			&message.ID,
			&message.SenderId,
			&message.ReceiverId,
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

	// to make the message from nold to new 
	for i, j:= 0, len(messages)-1; i< j;i,j = i+1, j-1{
		messages[i], messages[j] = messages[j] , messages[i]
	}

	return messages, nil
}
