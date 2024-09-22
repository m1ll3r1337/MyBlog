package db

import (
	"blog/models"
	"database/sql"
	"errors"
)

func (s *PostgresStore) CreatePostsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Posts (
			  id SERIAL PRIMARY KEY, 
			  title VARCHAR(255) NOT NULL,
			  content TEXT NOT NULL, 
			  user_id INT,
			  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreatePost(post *models.Post) (int, error) {
	query := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRow(query,
		post.Title,
		post.Content,
		post.UserID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetPosts() ([]models.Post, error) {
	rows, err := s.db.Query(`SELECT id, title, content, user_id FROM Posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostgresStore) GetPostByID(id int) (*models.Post, error) {
	query := `SELECT id, title, content, user_id FROM Posts WHERE id = $1`
	var post models.Post
	err := s.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (s *PostgresStore) UpdatePost(id int, post *models.Post) error {
	_, err := s.db.Query(`UPDATE posts SET title=$2, content=$3 WHERE id=$1`,
		id, post.Title, post.Content)
	return err
}

func (s *PostgresStore) DeletePost(id int) error {
	_, err := s.db.Query(`delete from posts where id = $1`, id)
	return err
}