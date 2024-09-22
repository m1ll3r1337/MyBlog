package db

import (
	"blog/models"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

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

func (s *PostgresStore) CreateUser(user *models.User) (int, error) {
	query := `INSERT INTO Users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRow(query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&id)
	if err != nil {
		return -1, err
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

func (s *PostgresStore) UpdateUser(id int, user *models.User) error {
	_, err := s.db.Query(`UPDATE users SET username=$2, email=$3 WHERE id=$1`,
		id, user.Username, user.Email)
	return err
}

func (s *PostgresStore) DeleteUser(id int) error {
	_, err := s.db.Query(`delete from users where id = $1`, id)
	return err
}

func (s *PostgresStore) Authenticate (email, password string) (*models.User, error){
	email = strings.ToLower(email)
	user := &models.User{
		Email: email,
	}
	row := s.db.QueryRow(`SELECT id, password from Users WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return user, nil
}