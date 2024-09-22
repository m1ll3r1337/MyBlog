package db

import (
	"blog/models"
	"database/sql"
	"errors"
)

func (s *PostgresStore) CreateImagesTable() error {
	query := `CREATE TABLE IF NOT EXISTS Images (
			  id SERIAL PRIMARY KEY, 
			  data BYTEA NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateImage(data []byte) (int, error) {
	query := `INSERT INTO Images (data) VALUES ($1) RETURNING id`
	var id int
	err := s.db.QueryRow(query, data).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAllImages() ([]models.Image, error) {
	rows, err := s.db.Query(`SELECT id, data FROM Images`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		if err := rows.Scan(&img.ID, &img.Data); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (s *PostgresStore) GetImageByID(id int) (*models.Image, error) {
	query := `SELECT id, data FROM Images WHERE id = $1`
	var img models.Image
	err := s.db.QueryRow(query, id).Scan(&img.ID, &img.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (s *PostgresStore) UploadImage(fileBytes []byte) error {
	_, err := s.db.Exec("INSERT INTO Images (data) VALUES ($1)", fileBytes)
	return err
}
