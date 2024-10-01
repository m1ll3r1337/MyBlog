package db

import (
	"blog/models"
	"database/sql"
	"errors"
)

func (s *PostgresStore) CreateCommentsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Comments (
			  id SERIAL PRIMARY KEY, 
			  content TEXT NOT NULL,
			  user_id INT,
			  post_id INT,
			  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
			  FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateComment(comment *models.Comment) (int, error) {
	query := `INSERT INTO Comments (content, user_id, post_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRow(query,
		comment.Content,
		comment.UserID,
		comment.PostID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetComments() ([]models.Comment, error) {
	rows, err := s.db.Query(`SELECT id, content, user_id, post_id FROM Comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *PostgresStore) GetCommentByID(id int) (*models.Comment, error) {
	query := `SELECT id, content, user_id, post_id FROM Comments WHERE id = $1`
	var comment models.Comment
	err := s.db.QueryRow(query, id).Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (s *PostgresStore) UpdateComment(id int, comment *models.Comment) error {
	_, err := s.db.Query(`UPDATE comments SET content=$2 WHERE id=$1`,
		id, comment.Content)
	return err
}

func (s *PostgresStore) DeleteComment(id int) error {
	_, err := s.db.Query(`delete from comments where id = $1`, id)
	return err
}