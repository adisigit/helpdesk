package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAuth(public, protected huma.API, h *handler.AuthHandler) {
	huma.Register(public, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/api/auth/register",
		Summary:       "Register a requester account",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusCreated,
	}, h.Register)

	huma.Register(public, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Summary:     "Login and get a JWT",
		Tags:        []string{"auth"},
	}, h.Login)

	huma.Register(public, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Summary:     "Clear the session cookie",
		Tags:        []string{"auth"},
	}, h.Logout)

	huma.Register(protected, huma.Operation{
		OperationID: "me",
		Method:      http.MethodGet,
		Path:        "/api/me",
		Summary:     "Current user profile",
		Tags:        []string{"auth"},
	}, h.Me)
}
