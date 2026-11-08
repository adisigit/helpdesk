package handler

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/domain"
	"helpdesk/backend/internal/server/middleware"
	"helpdesk/backend/internal/service"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	svc      *service.TicketService
	comments *service.CommentService
}

func NewTicketHandler(svc *service.TicketService, comments *service.CommentService) *TicketHandler {
	return &TicketHandler{
		svc:      svc,
		comments: comments,
	}
}

func actorFrom(ctx context.Context) (service.Actor, error) {
	p, ok := middleware.PrincipalFrom(ctx)
	if !ok {
		return service.Actor{}, huma.Error401Unauthorized("authentication required")
	}
	return service.Actor{ID: p.UserID, Role: p.Role}, nil
}

func parseID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, huma.Error422UnprocessableEntity("invalid id")
	}
	return id, nil
}

func mapTicketErr(err error) error {
	switch {
	case errors.Is(err, service.ErrTicketNotFound):
		return huma.Error404NotFound("ticket not found")
	case errors.Is(err, service.ErrForbidden), errors.Is(err, domain.ErrTransitionForbidden):
		return huma.Error403Forbidden(err.Error())
	case errors.Is(err, domain.ErrInvalidTransition),
		errors.Is(err, service.ErrInvalidAssignee),
		errors.Is(err, service.ErrInvalidFilter),
		errors.Is(err, service.ErrInvalidInput),
		errors.Is(err, service.ErrEmptyComment):
		return huma.Error422UnprocessableEntity(err.Error())
	case errors.Is(err, service.ErrNothingToUpdate):
		return huma.Error400BadRequest(err.Error())
	case errors.Is(err, service.ErrConflict):
		return huma.Error409Conflict(err.Error())
	}
	return err
}

type TicketBody struct {
	ID                 string    `json:"id" format:"uuid"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Category           string    `json:"category"`
	Priority           string    `json:"priority" enum:"low,medium,high,urgent"`
	Status             string    `json:"status" enum:"open,in_progress,resolved,closed"`
	Requester          UserBody  `json:"requester"`
	Assignee           *UserBody `json:"assignee,omitempty"`
	AllowedTransitions []string  `json:"allowed_transitions" doc:"Status tujuan yang boleh dipilih pemanggil saat ini"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func toTicketBody(t *models.Ticket, a service.Actor) TicketBody {
	targets := domain.AllowedTargets(t.Status, a.Role, t.RequesterID == a.ID)
	allowed := make([]string, 0, len(targets))
	for _, s := range targets {
		allowed = append(allowed, string(s))
	}
	b := TicketBody{
		ID:                 t.ID.String(),
		Title:              t.Title,
		Description:        t.Description,
		Category:           t.Category,
		Priority:           string(t.Priority),
		Status:             string(t.Status),
		Requester:          toUserBody(&t.Requester),
		AllowedTransitions: allowed,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
	if t.Assignee != nil {
		u := toUserBody(t.Assignee)
		b.Assignee = &u
	}
	return b
}

type TicketOutput struct{ Body TicketBody }

type CreateTicketInput struct {
	Body struct {
		Title       string `json:"title" minLength:"3" maxLength:"150" example:"VPN tidak bisa connect"`
		Description string `json:"description" minLength:"5" maxLength:"5000"`
		Category    string `json:"category" enum:"hardware,software,network,account,other"`
		Priority    string `json:"priority,omitempty" enum:"low,medium,high,urgent" default:"medium"`
	}
}

func (h *TicketHandler) Create(ctx context.Context, in *CreateTicketInput) (*TicketOutput, error) {
	actor, err := actorFrom(ctx)
	if err != nil {
		return nil, err
	}
	t, err := h.svc.Create(ctx, actor, service.CreateTicketInput{
		Title:       in.Body.Title,
		Description: in.Body.Description,
		Category:    in.Body.Category,
		Priority:    models.Priority(in.Body.Priority),
	})
	if err != nil {
		return nil, mapTicketErr(err)
	}
	return &TicketOutput{Body: toTicketBody(t, actor)}, nil
}

type GetTicketInput struct {
	ID string `path:"id" format:"uuid"`
}

type TicketDetailBody struct {
	TicketBody
	Comments []CommentBody `json:"comments" doc:"Requester tidak menerima catatan internal"`
}

type TicketDetailOutput struct{ Body TicketDetailBody }

func (h *TicketHandler) Get(ctx context.Context, in *GetTicketInput) (*TicketDetailOutput, error) {
	actor, err := actorFrom(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID(in.ID)
	if err != nil {
		return nil, err
	}
	t, err := h.svc.Get(ctx, actor, id) // pemeriksaan akses
	if err != nil {
		return nil, mapTicketErr(err)
	}
	cs, err := h.comments.ListFor(ctx, actor, t.ID)
	if err != nil {
		return nil, err
	}
	out := &TicketDetailOutput{}
	out.Body.TicketBody = toTicketBody(t, actor)
	out.Body.Comments = toCommentBodies(cs)
	return out, nil
}

type ListTicketsInput struct {
	Status   string `query:"status" enum:"open,in_progress,resolved,closed"`
	Priority string `query:"priority" enum:"low,medium,high,urgent"`
	Assignee string `query:"assignee" doc:"UUID agent, 'me', atau 'none' (belum ditugaskan)"`
	Q        string `query:"q" maxLength:"100" doc:"Cari berdasarkan judul"`
	Sort     string `query:"sort" enum:"created_at,-created_at,updated_at,-updated_at,priority,-priority,title,-title" default:"-created_at"`
	Page     int    `query:"page" minimum:"1" default:"1"`
	PageSize int    `query:"page_size" minimum:"1" maximum:"100" default:"20"`
}

type ListTicketsOutput struct {
	Body struct {
		Items      []TicketBody `json:"items"`
		Total      int64        `json:"total"`
		Page       int          `json:"page"`
		PageSize   int          `json:"page_size"`
		TotalPages int          `json:"total_pages"`
	}
}

func (h *TicketHandler) List(ctx context.Context, in *ListTicketsInput) (*ListTicketsOutput, error) {
	actor, err := actorFrom(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.svc.List(ctx, actor, service.ListTicketsInput{
		Status:   models.TicketStatus(in.Status),
		Priority: models.Priority(in.Priority),
		Assignee: in.Assignee,
		Query:    in.Q,
		Sort:     in.Sort,
		Page:     in.Page,
		PageSize: in.PageSize,
	})
	if err != nil {
		return nil, mapTicketErr(err)
	}

	out := &ListTicketsOutput{}
	out.Body.Items = make([]TicketBody, 0, len(res.Items))
	for i := range res.Items {
		out.Body.Items = append(out.Body.Items, toTicketBody(&res.Items[i], actor))
	}
	out.Body.Total = res.Total
	out.Body.Page = res.Page
	out.Body.PageSize = res.PageSize
	out.Body.TotalPages = int((res.Total + int64(res.PageSize) - 1) / int64(res.PageSize))
	return out, nil
}

type UpdateTicketInput struct {
	ID   string `path:"id" format:"uuid"`
	Body struct {
		Status     *string `json:"status,omitempty" enum:"open,in_progress,resolved,closed"`
		Priority   *string `json:"priority,omitempty" enum:"low,medium,high,urgent"`
		AssigneeID *string `json:"assignee_id,omitempty" format:"uuid"`
	}
}

func (h *TicketHandler) Update(ctx context.Context, in *UpdateTicketInput) (*TicketOutput, error) {
	actor, err := actorFrom(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID(in.ID)
	if err != nil {
		return nil, err
	}

	var upd service.UpdateTicketInput
	if b := in.Body.Status; b != nil {
		s := models.TicketStatus(*b)
		upd.Status = &s
	}
	if b := in.Body.Priority; b != nil {
		p := models.Priority(*b)
		upd.Priority = &p
	}
	if b := in.Body.AssigneeID; b != nil {
		aid, err := uuid.Parse(*b)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid assignee_id")
		}
		upd.AssigneeID = &aid
	}

	t, err := h.svc.Update(ctx, actor, id, upd)
	if err != nil {
		return nil, mapTicketErr(err)
	}
	return &TicketOutput{Body: toTicketBody(t, actor)}, nil
}
