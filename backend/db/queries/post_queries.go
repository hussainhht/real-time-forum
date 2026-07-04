package queries

import (
	"realtime/backend/global"
)



func InsertPost(Title string, Content string, CategoryId int) error {
	_, err := global.Database.Exec(`INSERT INTO posts (title, content, category_id)
	VALUES (?,?,?)
	`, Title, Content, CategoryId)

	if err != nil {
		return err
	}

	return nil
}

func InsertComment(PostId int, UserId int, Content string) error {
	_, err := global.Database.Exec(`INSERT INTO comments ( post_id, user_id, content)
	VALUES (?,?,?)
	`, PostId, UserId, Content)

	if err != nil {
		return err
	}

	return nil

}


