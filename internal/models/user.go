package models


type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password int64 `json:"-"`
	CreatedAt string `json:"created_at"`
	// UpdatedAt string `json:"updated_at"`
}