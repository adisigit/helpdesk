package service

import (
	"context"
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUsers struct {
	byEmail map[string]*models.User
}

func newFakeUsers() *fakeUsers {
	return &fakeUsers{byEmail: map[string]*models.User{}}
}

func (f *fakeUsers) Create(_ context.Context, u *models.User) error {
	if _, ok := f.byEmail[u.Email]; ok {
		return repository.ErrDuplicate
	}
	u.ID = uuid.New()
	f.byEmail[u.Email] = u
	return nil
}

func (f *fakeUsers) FindByEmail(_ context.Context, email string) (*models.User, error) {
	if u, ok := f.byEmail[email]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUsers) FindByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	for _, u := range f.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUsers) List(_ context.Context, _ repository.UserFilter) ([]models.User, int64, error) {
	return nil, 0, nil
}

func (f *fakeUsers) SetRole(_ context.Context, id uuid.UUID, role models.Role) (*models.User, error) {
	for _, u := range f.byEmail {
		if u.ID == id {
			u.Role = role
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func newSvc() (*AuthService, *auth.JWT) {
	j := auth.NewJWT("0123456789abcdef-secret", time.Hour)
	return NewAuthService(newFakeUsers(), j), j
}

func TestRegister_AlwaysRequesterAndHashesPassword(t *testing.T) {
	svc, _ := newSvc()
	u, err := svc.Register(context.Background(), "Budi", "Budi@Example.com", "rahasia123")
	require.NoError(t, err)
	assert.Equal(t, models.RoleRequester, u.Role)
	assert.Equal(t, "budi@example.com", u.Email)
	assert.NotEqual(t, "rahasia123", u.PasswordHash)
}

func TestRegister_DuplicateEmailIsCaseInsensitive(t *testing.T) {
	svc, _ := newSvc()
	_, err := svc.Register(context.Background(), "Budi", "budi@example.com", "rahasia123")
	require.NoError(t, err)
	_, err = svc.Register(context.Background(), "Budi 2", "BUDI@example.com", "rahasia123")
	assert.ErrorIs(t, err, ErrEmailTaken)
}

func TestLogin_WrongPasswordAndUnknownEmailRejected(t *testing.T) {
	svc, _ := newSvc()
	_, err := svc.Register(context.Background(), "Budi", "budi@example.com", "rahasia123")
	require.NoError(t, err)

	_, err = svc.Login(context.Background(), "budi@example.com", "salah-total")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	_, err = svc.Login(context.Background(), "tidak-ada@example.com", "rahasia123")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_ReturnsTokenWithRole(t *testing.T) {
	svc, j := newSvc()
	_, err := svc.Register(context.Background(), "Budi", "budi@example.com", "rahasia123")
	require.NoError(t, err)

	res, err := svc.Login(context.Background(), "budi@example.com", "rahasia123")
	require.NoError(t, err)
	claims, err := j.Parse(res.Token)
	require.NoError(t, err)
	assert.Equal(t, models.RoleRequester, claims.Role)
}
