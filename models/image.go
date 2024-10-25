package models

import (
	"database/sql"
	"errors"
)

type Image struct {
	ID    int    `json:"id"`
	Data []byte `json:"bytea"`
}

type ImageService struct {
	DB *sql.DB
}

func (is *ImageService) CreateImage(data []byte) (int, error) {
	query := `INSERT INTO Images (data) VALUES ($1) RETURNING id`
	var id int
	err := is.DB.QueryRow(query, data).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (is *ImageService) GetAllImages() ([]Image, error) {
	rows, err := is.DB.Query(`SELECT id, data FROM Images`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.Data); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (is *ImageService) GetImageByID(id int) (*Image, error) {
	query := `SELECT id, data FROM Images WHERE id = $1`
	var img Image
	err := is.DB.QueryRow(query, id).Scan(&img.ID, &img.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}

func (is *ImageService) ByPostID(postID int) (*Image, error) {
	query := `SELECT id, data FROM Images WHERE post_id = $1`
	var img Image
	err := is.DB.QueryRow(query, postID).Scan(&img.ID, &img.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}

func (is *ImageService) UploadImage(fileBytes []byte) error {
	_, err := is.DB.Exec("INSERT INTO Images (data) VALUES ($1)", fileBytes)
	return err
}
