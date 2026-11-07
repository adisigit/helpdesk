package auth

import (
	"helpdesk/backend/internal/database/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT_GenerateAndParse(t *testing.T) {
	j := NewJWT("0123456789abcdef-secret", time.Hour)
	id := uuid.New()

	tok, _, err := j.Generate(id, models.RoleAgent)
	require.NoError(t, err)

	claims, err := j.Parse(tok)
	require.NoError(t, err)
	got, _ := claims.UserID()
	assert.Equal(t, id, got)
	assert.Equal(t, models.RoleAgent, claims.Role)
}

func TestJWT_RejectsBadTokens(t *testing.T) {
	j := NewJWT("0123456789abcdef-secret", time.Hour)
	tok, _, _ := j.Generate(uuid.New(), models.RoleAdmin)

	_, err := j.Parse(tok + "X")
	assert.Error(t, err, "tampered")

	_, err = NewJWT("another-secret-0123456789", time.Hour).Parse(tok)
	assert.Error(t, err, "wrong secret")

	expired, _, _ := NewJWT("0123456789abcdef-secret", -time.Minute).Generate(uuid.New(), models.RoleAdmin)
	_, err = j.Parse(expired)
	assert.Error(t, err, "expired")
}
