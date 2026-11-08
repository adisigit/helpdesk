package repository

import (
	"context"
	"helpdesk/backend/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository interface {
	Create(ctx context.Context, c *models.Comment) error
	ListByTicket(ctx context.Context, ticketId uuid.UUID, includeInternal bool) ([]models.Comment, error)
}

type commentRepo struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) Create(ctx context.Context, c *models.Comment) error {
	if err := r.db.WithContext(ctx).Omit(clause.Associations).Create(c).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).First(&c.Author, "id = ?", c.AuthorID).Error
}

func (r *commentRepo) ListByTicket(ctx context.Context, ticketId uuid.UUID, includeInternal bool) ([]models.Comment, error) {
	q := r.db.WithContext(ctx).Preload("Author").Where("ticket_id = ?", ticketId)
	if !includeInternal {
		q = q.Where("is_internal = ?", false)
	}
	var out []models.Comment
	err := q.Order("created_at ASC, id ASC").Find(&out).Error
	return out, err
}
