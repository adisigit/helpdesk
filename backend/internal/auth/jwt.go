package auth

import (
	"errors"
	"helpdesk/backend/internal/database/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

func (c *Claims) UserID() (uuid.UUID, error) { return uuid.Parse(c.Subject) }

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl}
}

func (j *JWT) Generate(userID uuid.UUID, role models.Role) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(j.ttl)
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	return signed, exp, err
}

func (j *JWT) Parse(raw string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(raw, claims,
		func(*jwt.Token) (any, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !claims.Role.Valid() {
		return nil, errors.New("invalid role claim")
	}
	return claims, nil
}
