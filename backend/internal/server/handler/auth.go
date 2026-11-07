package handler

import (
	"context"
	"errors"
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/server/middleware"
	"helpdesk/backend/internal/service"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type CookieConfig struct {
	Secure bool
	TTL    time.Duration
}

type AuthHandler struct {
	svc    *service.AuthService
	cookie CookieConfig
}

func NewAuthHandler(svc *service.AuthService, c CookieConfig) *AuthHandler {
	return &AuthHandler{svc: svc, cookie: c}
}

func (h *AuthHandler) sessionCookie(value string, maxAge int) http.Cookie {
	return http.Cookie{
		Name:     auth.CookieName,
		Value:    value,
		Path:     "/api",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: http.SameSiteLaxMode,
	}
}

type UserBody struct {
	ID    string `json:"id" format:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role" enum:"requester,agent,admin"`
}

func toUserBody(u *models.User) UserBody {
	return UserBody{ID: u.ID.String(), Name: u.Name, Email: u.Email, Role: string(u.Role)}
}

type RegisterInput struct {
	Body struct {
		Name     string `json:"name" minLength:"2" maxLength:"100" example:"Budi Santoso"`
		Email    string `json:"email" format:"email" maxLength:"254" example:"budi@example.com"`
		Password string `json:"password" minLength:"8" maxLength:"72"`
	}
}

type RegisterOutput struct{ Body UserBody }

func (h *AuthHandler) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	u, err := h.svc.Register(ctx, in.Body.Name, in.Body.Email, in.Body.Password)
	switch {
	case errors.Is(err, service.ErrEmailTaken):
		return nil, huma.Error409Conflict("email already registered")
	case errors.Is(err, service.ErrPasswordTooLong):
		return nil, huma.Error422UnprocessableEntity("password must be at most 72 bytes")
	case err != nil:
		return nil, err
	}
	return &RegisterOutput{Body: toUserBody(u)}, nil
}

type LoginInput struct {
	Body struct {
		Email    string `json:"email" format:"email" example:"admin@helpdesk.com"`
		Password string `json:"password" minLength:"1"`
	}
}

type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      struct {
		AccessToken string    `json:"access_token" doc:"Untuk klien API/Swagger. Browser memakai cookie httpOnly."`
		TokenType   string    `json:"token_type" example:"Bearer"`
		ExpiresAt   time.Time `json:"expires_at"`
		User        UserBody  `json:"user"`
	}
}

func (h *AuthHandler) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	res, err := h.svc.Login(ctx, in.Body.Email, in.Body.Password)
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return nil, huma.Error401Unauthorized("invalid email or password")
	case err != nil:
		return nil, err
	}
	out := &LoginOutput{}
	out.SetCookie = h.sessionCookie(res.Token, int(h.cookie.TTL.Seconds()))
	out.Body.AccessToken = res.Token
	out.Body.TokenType = "Bearer"
	out.Body.ExpiresAt = res.ExpiresAt
	out.Body.User = toUserBody(res.User)
	return out, nil
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (h *AuthHandler) Logout(ctx context.Context, _ *struct{}) (*LogoutOutput, error) {
	return &LogoutOutput{SetCookie: h.sessionCookie("", -1)}, nil
}

type MeOutput struct{ Body UserBody }

func (h *AuthHandler) Me(ctx context.Context, _ *struct{}) (*MeOutput, error) {
	p, ok := middleware.PrincipalFrom(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("authentication required")
	}
	u, err := h.svc.GetUser(ctx, p.UserID)
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return nil, huma.Error401Unauthorized("user no longer exists")
	case err != nil:
		return nil, err
	}
	return &MeOutput{Body: toUserBody(u)}, nil
}
