package service

import "fmt"

var (
	ErrUserNotFound     = fmt.Errorf("user not found")
	ErrUserAlredyExists = fmt.Errorf("user alredy exists")
)
