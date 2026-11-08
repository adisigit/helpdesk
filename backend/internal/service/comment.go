package service

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type CommentService struct {
	tickets  repository.TicketRepository
	comments repository.CommentRepository
}

func NewCommentService(t repository.TicketRepository, c repository.CommentRepository) *CommentService {
	return &CommentService{
		tickets:  t,
		comments: c,
	}
}

func (s *CommentService) Add(ctx context.Context, actor Actor, ticketID uuid.UUID, body string, internal bool) (*models.Comment, error) {
	t, err := s.tickets.FindByID(ctx, ticketID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if !canAccess(actor, t) {
		return nil, ErrTicketNotFound
	}
	if internal && !actor.Role.IsStaff() {
		return nil, ErrForbidden
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, ErrEmptyComment
	}
	c := &models.Comment{
		TicketID:   t.ID,
		AuthorID:   actor.ID,
		Body:       body,
		IsInternal: internal,
	}
	if err := s.comments.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CommentService) ListFor(ctx context.Context, actor Actor, ticketID uuid.UUID) ([]models.Comment, error) {
	return s.comments.ListByTicket(ctx, ticketID, actor.Role.IsStaff())
}
