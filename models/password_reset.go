package models

import (
	"blog/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	//Token is only used when a PasswordReset is being created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	//Amount of time a PasswordReset is valid for
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := &PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: service.Hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = service.DB.QueryRow(`insert into password_resets (user_id, token_hash, expires_at) 
			VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET	token_hash = $2, expires_at = $3
			RETURNING id`, pwReset.UserID, pwReset.Token, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return pwReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := service.Hash(token)
	var user User
	var pwReset PasswordReset
	row := service.DB.QueryRow(`SELECT password_resets.id, password_resets.expires_at,
       users.id, users.email, users.password FROM password_resets JOIN users on users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1`, tokenHash)
	err := row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}
	err = service.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (service *PasswordResetService) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(tokenHash[:])
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`delete from password_resets where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}