package service

import "errors"

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrPasswordTooLong    = errors.New("password must be at most 72 bytes")
	ErrUserNotFound       = errors.New("user not found")
)
