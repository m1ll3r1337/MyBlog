package main

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
}

type Comment struct {
	ID int `json:"id"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Image struct {
	ID    int    `json:"id"`
	Data []byte `json:"bytea"`
}

