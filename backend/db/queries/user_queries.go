package queries

import "realtime/backend/global"

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

const insert_session = `
INSERT INTO sessions (
	id,
	user_id,
	created_at,
	expires_at
)
values (?,?,?,?)
`

const get_user_by_username_or_email = `
SELECT id, username, username, email, password
FROM users
WHERE username = ? or email = ?
`

func GetUserByUsernameOrEmail() {
	
}
