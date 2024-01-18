package repository

import "fmt"

var (
	ErrURLExists   = fmt.Errorf("Alias URL must be unique")
	ErrURLNotFound = fmt.Errorf("URL not found")
)
