package db

import "blog/models"

func (s *PostgresStore) CreateSessionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS Sessions (
			  id SERIAL PRIMARY KEY, 
			  user_id int UNIQUE not null, 
			  token_hash text unique not null
    )`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateSession(userID int) (*models.Session, error) {
	// TODO: Create the session token
	// TODO: Implement SessionService.Create
	return nil, nil
}
