package models


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

type Session struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	// Token is only set when creating a new session. When look up a session this will be left empty
	Token string `json:"token"`
	TokenHash string `json:"token_hash"`
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


