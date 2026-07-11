package queries

import (
	"realtime/backend/global"
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

func GetUserByUsernameOrEmail(identifier string) (*global.Users, error) {
	var user global.Users
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
	created_at,
	expires_at
)
values (?,?,?,?)
`
func InsertSession(sessionID string, userID int, expiresAt time.Time) error {
	_, err := global.Database.Exec(insert_session, userID, time.Now(), expiresAt)
	if err != nil {
		return err
	}
	return nil
}