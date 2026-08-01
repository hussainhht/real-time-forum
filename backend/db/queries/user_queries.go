package queries

import (
	"database/sql"
	"realtime/backend/global"
	"realtime/backend/global/structures"
	"time"
)

const insert_user = `
INSERT INTO users (
	username,
	first_name,
	last_name,
	age,
	phone_number,
	gender,
	email,
	password
)
values (?,?,?,?,?,?,?,?)
`

func InsertUser(username string, firstName string, lastName string, age int, phoneNumber string, gender string, email string, password string) error {
	_, err := global.Database.Exec(insert_user, username, firstName, lastName, age, phoneNumber, gender, email, password)
	if err != nil {
		return err
	}
	return nil
}

const get_user_by_username_or_email = `
SELECT id, username, email, password
FROM users
WHERE username = ? or email = ?
`

func GetUserByUsernameOrEmail(identifier string) (*structures.Users, error) {
	var user structures.Users
	row := global.Database.QueryRow(get_user_by_username_or_email, identifier, identifier)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

const insert_session = `
INSERT INTO sessions (
	id,
	user_id,
	expires_at
)
values (?,?,?)
`

func InsertSession(sessionID string, userID int, expiresAt time.Time) error {
	_, err := global.Database.Exec(insert_session, sessionID, userID, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

const get_user_id_by_Session = `
SELECT user_id, expires_at
FROM sessions
WHERE id = ?
`

const get_user_by_id = `
SELECT id, username, email
FROM users
WHERE id = ?
`

func GetUserBySession(sessionID string) (*structures.Users, error) {

	var userID int
	var expiresAt time.Time
	row := global.Database.QueryRow(get_user_id_by_Session, sessionID)
	err := row.Scan(&userID, &expiresAt)
	if err != nil {
		return nil, err
	}
	if time.Now().After(expiresAt) {
		return nil, sql.ErrNoRows //? is it the way you use it here?
	}

	var user structures.Users
	row = global.Database.QueryRow(get_user_by_id, userID)
	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

const deleteSession = `
DELETE FROM sessions
WHERE id = ?
`

func DeleteSession(sessionID string) error {
	_, err := global.Database.Exec(deleteSession, sessionID)
	return err
}

const getUserForChatQuery = `
SELECT id, username
FORM users
WHERE id !=?
ORDER BY LOWER(username) ASC
`

func GetUsersForChat(currentUserID int) ([]structures.ChatUser, error) {
	rows, err := global.Database.Query(
		getUserForChatQuery,
		currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []structures.ChatUser{}

	for rows.Next(){
		var user structures.ChatUser

		err := rows.Scan(
			&user.ID,
			&user.Username,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err := rows.Err();err != nil {
		return nil, err
		
	}

	return users, nil


}
