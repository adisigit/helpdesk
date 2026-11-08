package service

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(u repository.UserRepository) *UserService {
	return &UserService{users: u}
}

type ListUsersInput struct {
	Role     models.Role
	Query    string
	Page     int
	PageSize int
}

type UserPage struct {
	Items    []models.User
	Total    int64
	Page     int
	PageSize int
}

func (s *UserService) List(ctx context.Context, in ListUsersInput) (*UserPage, error) {
	if in.Role != "" && !in.Role.Valid() {
		return nil, ErrInvalidFilter
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	f := repository.UserFilter{
		Role:     in.Role,
		Query:    strings.TrimSpace(in.Query),
		Page:     max(in.Page, 1),
		PageSize: min(pageSize, 100),
	}
	items, total, err := s.users.List(ctx, f)
	if err != nil {
		return nil, err
	}
	return &UserPage{Items: items, Total: total, Page: f.Page, PageSize: f.PageSize}, nil
}

func (s *UserService) ChangeRole(ctx context.Context, id uuid.UUID, role models.Role) (*models.User, error) {
	if !role.Valid() {
		return nil, ErrInvalidRole
	}
	u, err := s.users.SetRole(ctx, id, role)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return nil, ErrUserNotFound
	case errors.Is(err, repository.ErrLastAdmin):
		return nil, ErrLastAdmin
	case err != nil:
		return nil, err
	}
	return u, nil
}
