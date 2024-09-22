package db

import (
	"blog/models"
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*models.User) (int, error)
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(int, *models.User) error
	DeleteUser(int) error

	CreatePost(*models.Post) (int, error)
	GetPosts() ([]models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	UpdatePost(int, *models.Post) error
	DeletePost(int) error

	CreateComment(*models.Comment) (int, error)
	GetComments() ([]models.Comment, error)
	GetCommentByID(id int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
	UpdateComment(int, *models.Comment) error
	DeleteComment(int) error

	UploadImage(data []byte) error
	GetAllImages() ([]models.Image, error)
	GetImageByID(id int) (*models.Image, error)
	GetImagesByPostID(postID int) ([]*models.Image, error)

	LinkPostToImage(postID, imageID int) error
}



type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetImagesByPostID(postID int) ([]*models.Image, error) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=balls sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() {
	s.CreateUsersTable()
	s.CreatePostsTable()
	s.CreateCommentsTable()
	s.CreateImagesTable()
	s.CreatePostsImagesTable()
}

func (s *PostgresStore) CreatePostsImagesTable() error {
	query := `CREATE TABLE IF NOT EXISTS Posts_Images (
			  post_id INT,
			  image_id INT,
			  PRIMARY KEY (post_id, image_id),
			  FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
		 	  FOREIGN KEY (image_id) REFERENCES Images(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) LinkPostToImage(postID, imageID int) error {
	query := `INSERT INTO Posts_Images (post_id, image_id) VALUES ($1, $2)`
	_, err := s.db.Exec(query, postID, imageID)
	return err
}






