package models

import (
	"database/sql"
	"errors"
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

func (ps *PostService) CreatePostsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Posts (
			  id SERIAL PRIMARY KEY, 
			  title VARCHAR(255) NOT NULL,
			  content TEXT NOT NULL, 
			  user_id INT,
			  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
	)`
	_, err := ps.DB.Exec(query)
	return err
}

func (ps *PostService) CreatePost(post *Post) (int, error) {
	query := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := ps.DB.QueryRow(query,
		post.Title,
		post.Content,
		post.UserID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ps *PostService) GetPosts() ([]Post, error) {
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

func (ps *PostService) GetPostByID(id int) (*Post, error) {
	query := `SELECT id, title, content, user_id FROM Posts WHERE id = $1`
	var post Post
	err := ps.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (ps *PostService) UpdatePost(id int, post *Post) error {
	_, err := ps.DB.Query(`UPDATE posts SET title=$2, content=$3 WHERE id=$1`,
		id, post.Title, post.Content)
	return err
}

func (ps *PostService) DeletePost(id int) error {
	_, err := ps.DB.Query(`delete from posts where id = $1`, id)
	return err
}

func (ps *PostService) LinkPostToImage(postID, imageID int) error {
	query := `INSERT INTO Posts_Images (post_id, image_id) VALUES ($1, $2)`
	_, err := ps.DB.Exec(query, postID, imageID)
	return err
}

func (ps *PostService) CreatePostsImagesTable() error {
	query := `CREATE TABLE IF NOT EXISTS Posts_Images (
			  post_id INT,
			  image_id INT,
			  PRIMARY KEY (post_id, image_id),
			  FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
		 	  FOREIGN KEY (image_id) REFERENCES Images(id) ON DELETE CASCADE
	)`
	_, err := ps.DB.Exec(query)
	return err
}