package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterComments(protected huma.API, h *handler.CommentHandler) {
	huma.Register(protected, huma.Operation{
		OperationID:   "addComment",
		Method:        http.MethodPost,
		Path:          "/api/tickets/{id}/comments",
		Summary:       "Add a comment or internal note",
		Description:   "is_internal=true hanya boleh agent/admin (requester mendapat 403).",
		Tags:          []string{"comments"},
		DefaultStatus: http.StatusCreated,
	}, h.Add)
}
