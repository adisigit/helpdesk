package router

import (
	"helpdesk/backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterTickets(protected huma.API, h *handler.TicketHandler) {
	huma.Register(protected, huma.Operation{
		OperationID:   "createTicket",
		Method:        http.MethodPost,
		Path:          "/api/tickets",
		Summary:       "Create a ticket",
		Tags:          []string{"tickets"},
		DefaultStatus: http.StatusCreated,
	}, h.Create)

	huma.Register(protected, huma.Operation{
		OperationID: "getTicket",
		Method:      http.MethodGet,
		Path:        "/api/tickets/{id}",
		Summary:     "Get ticket detail (owner, agent, admin)",
		Tags:        []string{"tickets"},
	}, h.Get)

	huma.Register(protected, huma.Operation{
		OperationID: "listTickets",
		Method:      http.MethodGet,
		Path:        "/api/tickets",
		Summary:     "List tickets (filter, search, sort, pagination)",
		Description: "Requester hanya melihat tiket miliknya. Agent dan admin melihat semua.",
		Tags:        []string{"tickets"},
	}, h.List)

	huma.Register(protected, huma.Operation{
		OperationID: "updateTicket",
		Method:      http.MethodPatch,
		Path:        "/api/tickets/{id}",
		Summary:     "Change status, priority or assignee",
		Description: "Status mengikuti tabel transisi: di luar tabel 422, peran tidak berhak 403. Prioritas dan assignee hanya untuk agent/admin.",
		Tags:        []string{"tickets"},
	}, h.Update)
}
