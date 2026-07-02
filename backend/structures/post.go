package structures

type Post struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryId int    `json:"category_id"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Like struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}

type Dislike struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}
