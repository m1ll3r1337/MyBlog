package main

import	"golang.org/x/crypto/bcrypt"


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

type CreateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	Title string `json:"title"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
}

type CreateCommentRequest struct {
	PostID int `json:"post_id"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
}

type CreateImageRequest struct {
	ID int `json:"id"`
	Data []byte `json:"bytea"`
}

func NewUser(username, email, password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Username: username,
		Email: email,
		Password: string(encpw),
	}, nil
}

func NewPost(title, content string, userID int ) *Post {
	return &Post{
		Title: title,
		Content: content,
		UserID: userID,
	}
}

func NewComment(content string, userID, postID int) *Comment {
	return &Comment{
		Content: content,
		PostID: postID,
		UserID: userID,
	}
}


