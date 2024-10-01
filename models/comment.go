package models

import (
	"database/sql"
	"errors"
)

type Comment struct {
	ID int `json:"id"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type CommentService struct {
	DB *sql.DB
}

func NewComment(content string, userID, postID int) *Comment {
	return &Comment{
		Content: content,
		PostID: postID,
		UserID: userID,
	}
}

func (cs *CommentService) CreateCommentsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Comments (
			  id SERIAL PRIMARY KEY, 
			  content TEXT NOT NULL,
			  user_id INT,
			  post_id INT,
			  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
			  FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE
	)`
	_, err := cs.DB.Exec(query)
	return err
}

func (cs *CommentService) CreateComment(comment *Comment) (int, error) {
	query := `INSERT INTO Comments (content, user_id, post_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := cs.DB.QueryRow(query,
		comment.Content,
		comment.UserID,
		comment.PostID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (cs *CommentService) GetComments() ([]Comment, error) {
	rows, err := cs.DB.Query(`SELECT id, content, user_id, post_id FROM Comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err = rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cs *CommentService) GetCommentByID(id int) (*Comment, error) {
	query := `SELECT id, content, user_id, post_id FROM Comments WHERE id = $1`
	var comment Comment
	err := cs.DB.QueryRow(query, id).Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (cs *CommentService) UpdateComment(id int, comment *Comment) error {
	_, err := cs.DB.Query(`UPDATE comments SET content=$2 WHERE id=$1`,
		id, comment.Content)
	return err
}

func (cs *CommentService) DeleteComment(id int) error {
	_, err := cs.DB.Query(`delete from comments where id = $1`, id)
	return err
}