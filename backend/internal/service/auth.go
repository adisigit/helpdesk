package service

import (
	"context"
	"errors"
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	users repository.UserRepository
	jwt   *auth.JWT
}

func NewAuthService(users repository.UserRepository, jwt *auth.JWT) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
	User      *models.User
}

func normalizeEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	if len([]byte(password)) > auth.MaxPasswordBytes {
		return nil, ErrPasswordTooLong
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Name:         strings.TrimSpace(name),
		Email:        normalizeEmail(email),
		PasswordHash: hash,
		Role:         models.RoleRequester,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	u, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.CheckPassword(u.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	token, exp, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, ExpiresAt: exp, User: u}, nil
}

func (s *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	return u, err
}
