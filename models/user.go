package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func New(username, email, password string) (*User, error) {
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

func (us *UserService) Create(email, username, password string) (*User, error) {
	email = strings.ToLower(email)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	password = string(hash)

	user := User{
		Username: username,
		Email:    email,
		Password: password,
	}
	row := us.DB.QueryRow(`insert into users (username, email, password) values ($1, $2, $3) returning id;`,
		username, email, password)
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := &User{
		Email: email,
	}
	row := us.DB.QueryRow(`SELECT id, password from Users WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotAuthenticated
		}
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrNotAuthenticated
		}
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
