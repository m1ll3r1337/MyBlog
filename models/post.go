package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Post struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	UserID int `json:"user_id"`
}

type PostService struct{
	DB *sql.DB
}

func NewPost(title, content string, userID int ) *Post {
	return &Post{
		Title: title,
		Content: content,
		UserID: userID,
	}
}

func (ps *PostService) Create(title, content string, userID int) (*Post, error) {
	query := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`
	post := Post{
		Title: title,
		Content: content,
		UserID: userID,
	}
	err := ps.DB.QueryRow(query,
		title,
		content,
		userID,
	).Scan(&post.ID)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return &post, nil
}

func (ps *PostService) GetAll() ([]Post, error) {
	rows, err := ps.DB.Query(`SELECT id, title, content, user_id FROM Posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (ps *PostService) GetByID(id int) (*Post, error) {
	query := `SELECT id, title, content, user_id FROM Posts WHERE id = $1`
	var post Post
	err := ps.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (ps *PostService) Update(post *Post) error {
	_, err := ps.DB.Query(`UPDATE posts SET title=$2, content=$3 WHERE id=$1`,
		post.ID, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

func (ps *PostService) Delete(id int) error {
	_, err := ps.DB.Exec(`delete from posts where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}

