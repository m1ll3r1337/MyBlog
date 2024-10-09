package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	DB *sql.DB
}

func NewUser(username, email, password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Username: username,
		Email:    email,
		Password: string(encpw),
	}, nil
}

func (us *UserService) CreateUser(user *User) (int, error) {
	query := `INSERT INTO Users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := us.DB.QueryRow(query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (us *UserService) GetUsers() ([]User, error) {
	rows, err := us.DB.Query(`SELECT id, username, email FROM Users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserService) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, email FROM Users WHERE id = $1`
	var user User
	err := us.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (us *UserService) UpdateUser(id int, user *User) error {
	_, err := us.DB.Query(`UPDATE users SET username=$2, email=$3 WHERE id=$1`,
		id, user.Username, user.Email)
	return err
}

func (us *UserService) DeleteUser(id int) error {
	_, err := us.DB.Query(`delete from users where id = $1`, id)
	return err
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := &User{
		Email: email,
	}
	row := us.DB.QueryRow(`SELECT id, password from Users WHERE email = $1`, email)
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

func (us *UserService) UpdatePassword(userID int, password string) error {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %w", err)
	}
	passwordHash := string(encpw)
	_, err = us.DB.Exec(`UPDATE users set password = $2 WHERE id = $1`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %w", err)
	}
	return nil
}
