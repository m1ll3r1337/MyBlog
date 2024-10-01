package models

import "database/sql"

type Session struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	// Token is only set when creating a new session. When look up a session this will be left empty
	Token string `json:"token"`
	TokenHash string `json:"token_hash"`
}

type SessionService struct {
	DB *sql.DB
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

func (ss *SessionService) CreateSession(userID int) (*Session, error) {
	// TODO: Create the session token
	// TODO: Implement SessionService.Create
	return nil, nil
}
