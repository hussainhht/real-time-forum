package queries

import (
	"database/sql"
	"realtime/backend/global"
)

const insert_post = `INSERT INTO posts (userID ,title, content, category_id) VALUES (?,?,?,?)`

func InsertPost(UserID int, Title string, Content string, CategoryId int) (int64, error) {
	result, err := global.Database.Exec(insert_post, UserID, Title, Content, CategoryId)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId() //this function already returns (int64, error)
}

const insert_comment = `
INSERT INTO comments (post_id, user_id, content) VALUES (?,?,?)`

func InsertComment(PostId int, UserId int, Content string) (int64, error) {
	result, err := global.Database.Exec(insert_comment, PostId, UserId, Content)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()

}

const category_exists = `SELECT 1 FROM categories WHERE id = ?`

func CategoryExists(categoryID int) (bool, error) {
	var one int
	err := global.Database.QueryRow(category_exists, categoryID).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
