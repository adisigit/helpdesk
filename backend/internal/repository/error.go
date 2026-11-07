package repository

import "errors"

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("record already exists")
	ErrConflict  = errors.New("record was modified concurrently")
)
