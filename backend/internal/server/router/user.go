package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterUsers(adminOnly huma.API, h *handler.UserHandler) {
	huma.Register(adminOnly, huma.Operation{
		OperationID: "listUsers",
		Method:      http.MethodGet,
		Path:        "/api/users",
		Summary:     "List users (admin)",
		Tags:        []string{"users"},
	}, h.List)

	huma.Register(adminOnly, huma.Operation{
		OperationID: "changeUserRole",
		Method:      http.MethodPatch,
		Path:        "/api/users/{id}/role",
		Summary:     "Change a user's role (admin)",
		Description: "Menurunkan admin terakhir ditolak dengan 409.",
		Tags:        []string{"users"},
	}, h.ChangeRole)

	huma.Register(adminOnly, huma.Operation{
		OperationID: "listAssignees",
		Method:      http.MethodGet,
		Path:        "/api/assignees",
		Summary:     "List users who can be assigned to tickets (admin)",
		Tags:        []string{"users"},
	}, h.ListAssignees)
}
