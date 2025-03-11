package models

type Post struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"`
	Version int `json:"version"`
	UpdatedAt string `json:"updated_at"`
	User User `json:"user"`
}


type PostWithMetadata struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"`
	Version int `json:"version"`
	UpdatedAt string `json:"updated_at"`
}