package structures

type Post struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_ID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NewPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type FeedPost struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Username string `json:"username"`

	//? we can add here created at to know when this post create
}

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

type FeedComment struct {
	ID     int `json:"id"`
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`


	// CreatedAt string `json:"created_at"` //? we can added it if we want to know when this comment added

}

type NewComment struct {
	Content string `json:"content"`
}
