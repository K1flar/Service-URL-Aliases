package repository

import "fmt"

var (
	ErrURLExists   = fmt.Errorf("alias URL must be unique")
	ErrURLNotFound = fmt.Errorf("url not found")

	ErrUserExists   = fmt.Errorf("user already exists")
	ErrUserNotFound = fmt.Errorf("user not found")
)
