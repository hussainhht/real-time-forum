package global

type Post struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_ID"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
}

type NewPost struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
}

type FeedPost struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Category string `json:"category"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

type NewComment struct {
	Content string `json:"content"`
}

type Like struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}

type Dislike struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}
