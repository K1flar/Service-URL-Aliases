package service

import "fmt"

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserExists   = fmt.Errorf("user alredy exists")

	ErrURLNotFound          = fmt.Errorf("url not found")
	ErrURLExists            = fmt.Errorf("alias URL must be unique")
	ErrURLForbiddenToDelete = fmt.Errorf("forbidden to delete")
)
