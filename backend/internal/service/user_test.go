package service

import (
	"context"
	"helpdesk/backend/internal/database/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeRole_RejectsInvalidRole(t *testing.T) {
	svc := NewUserService(newFakeUsers())
	_, err := svc.ChangeRole(context.Background(), uuid.New(), models.Role("superuser"))
	assert.ErrorIs(t, err, ErrInvalidRole)
}

func TestChangeRole_UnknownUser(t *testing.T) {
	svc := NewUserService(newFakeUsers())
	_, err := svc.ChangeRole(context.Background(), uuid.New(), models.RoleAgent)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestChangeRole_Success(t *testing.T) {
	users := newFakeUsers()
	u := &models.User{Email: "a@x.io", Role: models.RoleRequester}
	require.NoError(t, users.Create(context.Background(), u))

	got, err := NewUserService(users).ChangeRole(context.Background(), u.ID, models.RoleAgent)
	require.NoError(t, err)
	assert.Equal(t, models.RoleAgent, got.Role)
}
