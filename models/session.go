package models

import (
	"blog/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	// Token is only set when creating a new session. When look up a session this will be left empty
	Token string `json:"token"`
	TokenHash string `json:"token_hash"`
}

type SessionService struct {
	DB *sql.DB
	BytesPerToken int
}

func (ss *SessionService) CreateSessionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Sessions (
			  id SERIAL PRIMARY KEY, 
			  user_id int UNIQUE not null, 
			  token_hash text unique not null
    )`
	_, err := ss.DB.Exec(query)
	return err
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID: userID,
		Token: token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`UPDATE Sessions SET token_hash = $2 WHERE user_id = $1`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if errors.Is(err, sql.ErrNoRows) {
		row = ss.DB.QueryRow(`INSERT INTO Sessions (user_id, token_hash) VALUES ($1, $2)
                                           returning id`, userID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`SELECT user_id FROM Sessions WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	row = ss.DB.QueryRow(`SELECT email, password FROM Users WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, err
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
