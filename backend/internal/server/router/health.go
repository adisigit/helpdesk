package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHealth(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/api/health",
		Summary:     "Health check",
		Tags:        []string{"system"},
	}, handler.Health)
}
