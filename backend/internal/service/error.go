package service

import "errors"

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrPasswordTooLong    = errors.New("password must be at most 72 bytes")
	ErrUserNotFound       = errors.New("user not found")

	ErrTicketNotFound  = errors.New("ticket not found")
	ErrForbidden       = errors.New("you are not allowed to perform this action")
	ErrInvalidInput    = errors.New("title and description must not be blank")
	ErrInvalidFilter   = errors.New("invalid filter value")
	ErrInvalidAssignee = errors.New("assignee must be an existing agent or admin")
	ErrNothingToUpdate = errors.New("no fields to update")
	ErrConflict        = errors.New("ticket was modified by someone else, please reload")

	ErrEmptyComment = errors.New("comment body must not be blank")
	ErrInvalidRole  = errors.New("invalid role")
	ErrLastAdmin    = errors.New("cannot remove the last admin")
)
