package queries

import (
	"database/sql"
	"realtime/backend/global"
	"realtime/backend/global/structures"
)

// * ------------------------------------------------
// *----------POST SECTION -------------------------
// *------------------------------------------------
const insert_post = `INSERT INTO posts (userID ,title, content) VALUES (?,?,?)`

func InsertPost(UserID int, Title string, Content string) (int64, error) {
	result, err := global.Database.Exec(insert_post, UserID, Title, Content)

	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
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
SELECT posts.id, posts.title, posts.content, users.username
FROM posts
JOIN users ON posts.userID = users.id
ORDER BY posts.id DESC
`

func GetPostFeed() ([]structures.FeedPost, error) {
	rows, err := global.Database.Query(get_posts_list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structures.FeedPost
	for rows.Next() {
		var p structures.FeedPost
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Username)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil { //to check if an error occurred during iteration and did not appear.
		return nil, err
	}

	return posts, nil
}

const getPost = `
SELECT
	posts.id,
	posts.title,
	posts.content,
	users.username
FROM posts
JOIN users ON posts.userID = users.id
WHERE posts.id = ?
`

func GetPostByID(postID int) (*structures.FeedPost, error) {
	var post structures.FeedPost
	err := global.Database.QueryRow(
		getPost,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Username,
	)

	if err != nil {
		return nil, err
	}

	return &post, err
}

//* ------------------------------------------------
//*----------COMMENT SECTION -------------------------
//*------------------------------------------------

const insert_comment = `
INSERT INTO comments (post_id, user_id, content) VALUES (?,?,?)`

func InsertComment(PostId int, UserId int, Content string) (int64, error) {
	result, err := global.Database.Exec(insert_comment, PostId, UserId, Content)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()

}

const feedCommentsByPostID = `
	SELECT
		comments.id,
		comments.post_id,
		comments.user_id,
		users.username,
		comments.content
	FROM comments
	JOIN users ON comments.user_id = users.id
	WHERE comments.post_id = ?
	ORDER BY comments.id ASC
`

func FeedCommentsByPostID(postID int) ([]structures.FeedComment, error) {
	rows, err := global.Database.Query(
		feedCommentsByPostID,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]structures.FeedComment, 0)
	for rows.Next() {

		var comment structures.FeedComment

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Username,
			&comment.Content,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}

	return comments, nil

}

