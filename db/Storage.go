package db

import (
	"blog/models"
	"database/sql"
	"errors"
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

func (s *PostgresStore) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS Users (
			  id SERIAL PRIMARY KEY, 
			  username VARCHAR(255) NOT NULL,
			  email VARCHAR(255) UNIQUE NOT NULL,
			  password VARCHAR(255) NOT NULL
    )`
	_, err := s.db.Exec(query)
	return err
}

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

func (s *PostgresStore) CreateImagesTable() error {
	query := `CREATE TABLE IF NOT EXISTS Images (
			  id SERIAL PRIMARY KEY, 
			  data BYTEA NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
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

func (s *PostgresStore) UploadImage(fileBytes []byte) error {
	_, err := s.db.Exec("INSERT INTO Images (data) VALUES ($1)", fileBytes)
	return err
}

func (s *PostgresStore) CreateUser(user *models.User) (int, error) {
	query := `INSERT INTO Users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRow(query,
		user.Username,
		user.Email,
		user.Password,
		).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetUsers() ([]models.User, error) {
	rows, err := s.db.Query(`SELECT id, username, email FROM Users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *PostgresStore) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, email FROM Users WHERE id = $1`
	var user models.User
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
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

func (s *PostgresStore) UpdateUser(id int, user *models.User) error {
	_, err := s.db.Query(`UPDATE users SET username=$2, email=$3 WHERE id=$1`,
		id, user.Username, user.Email)
	return err
}

func (s *PostgresStore) DeleteUser(id int) error {
	_, err := s.db.Query(`delete from users where id = $1`, id)
	return err
}

func (s *PostgresStore) DeletePost(id int) error {
	_, err := s.db.Query(`delete from posts where id = $1`, id)
	return err
}

func (s *PostgresStore) DeleteComment(id int) error {
	_, err := s.db.Query(`delete from comments where id = $1`, id)
	return err
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



func (s *PostgresStore) LinkPostToImage(postID, imageID int) error {
	query := `INSERT INTO Posts_Images (post_id, image_id) VALUES ($1, $2)`
	_, err := s.db.Exec(query, postID, imageID)
	return err
}






