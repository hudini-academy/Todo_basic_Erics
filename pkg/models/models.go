package models

import (
	"errors"
	"time"
)

type Todo struct {
	Id      int
	Title   string
	Created time.Time
	Expires time.Time
	Tags 	string
	Tag []string
}

// hold the user data.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
