package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAuth(api huma.API, h *handler.AuthHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/api/auth/register",
		Summary:       "Register a requester account",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusCreated,
	}, h.Register)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Summary:     "Login and get a JWT",
		Tags:        []string{"auth"},
	}, h.Login)

	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Summary:     "Clear the session cookie",
		Tags:        []string{"auth"},
	}, h.Logout)
}
