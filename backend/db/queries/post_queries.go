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

const post_exists = `SELECT 1 FROM posts WHERE id = ?`

func PostExists(postID int) (bool, error) {
	var one int
	err := global.Database.QueryRow(post_exists, postID).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

const get_posts_list = `
SELECT posts.id, posts.title, posts.content, users.username, categories.name
FROM posts
JOIN users ON posts.userID = users.id
JOIN categories ON posts.category_id = categories.id
ORDER BY posts.id DESC
`

func GetFeed() ([]global.FeedPost, error) {
	rows, err := global.Database.Query(get_posts_list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []global.FeedPost
	for rows.Next() {
		var p global.FeedPost
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Username, &p.Category)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil { //to check if an error occurred during iteration and not appearing.
		return nil, err
	}

	return posts, nil
}

const insert_like = `INSERT INTO likes (post_id,user_id)
 VALUES (?, ?)`

func InsertLike(postID int, userId int) error {
	_, err := global.Database.Exec(insert_like, postID, userId)

	if err != nil {
		return err
	}
	return nil

}
