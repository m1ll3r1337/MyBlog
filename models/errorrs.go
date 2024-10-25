package models

import "errors"

var (
	ErrNotFound = errors.New("recourse could not be found")
	ErrEmailTaken = errors.New("models: email already taken")
)
