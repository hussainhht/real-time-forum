package queries

import "realtime/backend/global"

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
