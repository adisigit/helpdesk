package handler

import (
	"context"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/service"
	"time"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

type CommentBody struct {
	ID         string    `json:"id" format:"uuid"`
	Author     UserBody  `json:"author"`
	Body       string    `json:"body"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func toCommentBody(c *models.Comment) CommentBody {
	return CommentBody{
		ID:         c.ID.String(),
		Author:     toUserBody(&c.Author),
		Body:       c.Body,
		IsInternal: c.IsInternal,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func toCommentBodies(cs []models.Comment) []CommentBody {
	out := make([]CommentBody, 0, len(cs))
	for i := range cs {
		out = append(out, toCommentBody(&cs[i]))
	}
	return out
}

type AddCommentInput struct {
	ID   string `path:"id" format:"uuid"`
	Body struct {
		Body       string `json:"body" minLength:"1" maxLength:"5000"`
		IsInternal bool   `json:"is_internal,omitempty" default:"false"`
	}
}

type CommentOutout struct {
	Body CommentBody
}

func (h *CommentHandler) Add(ctx context.Context, in *AddCommentInput) (*CommentOutout, error) {
	actor, err := actorFrom(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID(in.ID)
	if err != nil {
		return nil, err
	}
	c, err := h.svc.Add(ctx, actor, id, in.Body.Body, in.Body.IsInternal)
	if err != nil {
		return nil, mapTicketErr(err)
	}
	return &CommentOutout{
		Body: toCommentBody(c),
	}, nil
}
