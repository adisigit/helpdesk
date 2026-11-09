package service

import (
	"context"
	"errors"
	"fmt"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/domain"
	"helpdesk/backend/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type Actor struct {
	ID   uuid.UUID
	Role models.Role
}

type TicketService struct {
	tickets repository.TicketRepository
	users   repository.UserRepository
}

func NewTicketService(t repository.TicketRepository, u repository.UserRepository) *TicketService {
	return &TicketService{
		tickets: t,
		users:   u,
	}
}

func canAccess(a Actor, t *models.Ticket) bool {
	return a.Role.IsStaff() || t.RequesterID == a.ID
}

type CreateTicketInput struct {
	Title       string
	Description string
	Category    string
	Priority    models.Priority
}

func (s *TicketService) Create(ctx context.Context, actor Actor, in CreateTicketInput) (*models.Ticket, error) {
	title := strings.TrimSpace(in.Title)
	desc := strings.TrimSpace(in.Description)
	if title == "" || desc == "" {
		return nil, ErrInvalidInput
	}
	if in.Priority == "" {
		in.Priority = models.PriorityMedium
	}
	t := &models.Ticket{
		Title:       title,
		Description: desc,
		Category:    in.Category,
		Priority:    in.Priority,
		Status:      models.StatusOpen,
		RequesterID: actor.ID,
	}
	if err := s.tickets.Create(ctx, t); err != nil {
		return nil, err
	}
	return s.tickets.FindByID(ctx, t.ID)
}

func (s *TicketService) Get(ctx context.Context, actor Actor, id uuid.UUID) (*models.Ticket, error) {
	t, err := s.tickets.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if !canAccess(actor, t) {
		return nil, ErrTicketNotFound
	}
	return t, nil
}

type ListTicketsInput struct {
	Status   models.TicketStatus
	Priority models.Priority
	Assignee string
	Query    string
	Sort     string
	Page     int
	PageSize int
}

type TicketPage struct {
	Items    []models.Ticket
	Total    int64
	Page     int
	PageSize int
}

func (s *TicketService) List(ctx context.Context, actor Actor, in ListTicketsInput) (*TicketPage, error) {
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	f := repository.TicketFilter{
		Status:   in.Status,
		Priority: in.Priority,
		Query:    strings.TrimSpace(in.Query),
		Sort:     in.Sort,
		Page:     max(in.Page, 1),
		PageSize: min(pageSize, 100),
	}

	if !actor.Role.IsStaff() {
		id := actor.ID
		f.RequesterID = &id
	}

	switch in.Assignee {
	case "":
	case "none":
		f.Unassigned = true
	case "me":
		id := actor.ID
		f.AssigneeID = &id
	default:
		id, err := uuid.Parse(in.Assignee)
		if err != nil {
			return nil, ErrInvalidFilter
		}
		f.AssigneeID = &id
	}

	items, total, err := s.tickets.List(ctx, f)
	if err != nil {
		return nil, err
	}
	return &TicketPage{Items: items, Total: total, Page: f.Page, PageSize: f.PageSize}, nil
}

type UpdateTicketInput struct {
	Status     *models.TicketStatus
	Priority   *models.Priority
	AssigneeID *uuid.UUID
}

func (s *TicketService) Update(ctx context.Context, actor Actor, id uuid.UUID, in UpdateTicketInput) (*models.Ticket, error) {
	t, err := s.Get(ctx, actor, id)
	if err != nil {
		return nil, err
	}
	if in.Status == nil && in.Priority == nil && in.AssigneeID == nil {
		return nil, ErrNothingToUpdate
	}
	if !actor.Role.IsStaff() && (in.Priority != nil || in.AssigneeID != nil) {
		return nil, ErrForbidden
	}

	from := t.Status
	if in.Status != nil {
		isOwner := t.RequesterID == actor.ID
		if err := domain.CanTransition(from, *in.Status, actor.Role, isOwner); err != nil {
			return nil, fmt.Errorf("%w: %s -> %s", err, from, *in.Status)
		}
		t.Status = *in.Status
	}

	if in.Priority != nil {
		t.Priority = *in.Priority
	}

	if in.AssigneeID != nil {
		if actor.Role == models.RoleAgent {
			if *in.AssigneeID != actor.ID {
				return nil, ErrForbidden
			}
			if t.AssigneeID != nil && *t.AssigneeID != actor.ID {
				return nil, ErrForbidden
			}
		}
		a, err := s.users.FindByID(ctx, *in.AssigneeID)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidAssignee
		}
		if err != nil {
			return nil, err
		}
		if !a.Role.IsStaff() {
			return nil, ErrInvalidAssignee
		}
		t.AssigneeID = in.AssigneeID
	} else if in.Status != nil && *in.Status == models.StatusInProgress && t.AssigneeID == nil {
		self := actor.ID
		t.AssigneeID = &self
	}

	if err := s.tickets.Update(ctx, t, from); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return s.tickets.FindByID(ctx, t.ID)
}
