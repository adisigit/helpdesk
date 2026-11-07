package repository

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	Create(ctx context.Context, t *models.Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	List(ctx context.Context, f TicketFilter) ([]models.Ticket, int64, error)
	Update(ctx context.Context, t *models.Ticket, expectedStatus models.TicketStatus) error
}

type ticketRepo struct {
	db *gorm.DB
}

func NewticketReporitory(db *gorm.DB) TicketRepository {
	return &ticketRepo{
		db: db,
	}
}

func (r *ticketRepo) Create(ctx context.Context, t *models.Ticket) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(t).Error
}

func (r *ticketRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var t models.Ticket
	err := r.db.WithContext(ctx).Preload("Requester").Preload("Assignee").First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type TicketFilter struct {
	RequesterID *uuid.UUID
	Status      models.TicketStatus
	Priority    models.Priority
	AssigneeID  *uuid.UUID
	Unassigned  bool
	Query       string
	Sort        string
	Page        int
	PageSize    int
}

var sortExprs = map[string]string{
	"created_at": "tickets.created_at",
	"updated_at": "tickets.updated_at",
	"title":      "LOWER(tickets.title)",
	"priority":   "CASE tickets.priority WHEN 'low' THEN 1 WHEN 'medium' THEN 2 WHEN 'high' THEN 3 WHEN 'urgent' THEN 4 END",
}

func orderClause(sort string) string {
	dir := "ASC"
	if strings.HasPrefix(sort, "-") {
		dir = "DESC"
		sort = sort[1:]
	}
	expr, ok := sortExprs[sort]
	if !ok {
		expr, dir = sortExprs["created_at"], "DESC"
	}
	return expr + " " + dir + ", tickets.id ASC"
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func applyFilter(db *gorm.DB, f TicketFilter) *gorm.DB {
	db = db.Model(&models.Ticket{})
	if f.RequesterID != nil {
		db = db.Where("tickets.requester_id = ?", *f.RequesterID)
	}
	if f.Status != "" {
		db = db.Where("tickets.status = ?", f.Status)
	}
	if f.Priority != "" {
		db = db.Where("tickets.priority = ?", f.Priority)
	}
	if f.AssigneeID != nil {
		db = db.Where("tickets.assignee_id = ?", *f.AssigneeID)
	} else if f.Unassigned {
		db = db.Where("tickets.assignee_id IS NULL")
	}
	if f.Query != "" {
		db = db.Where("tickets.title ILIKE ?", "%"+likeEscaper.Replace(f.Query)+"%")
	}
	return db
}

func (r *ticketRepo) List(ctx context.Context, f TicketFilter) ([]models.Ticket, int64, error) {
	var total int64
	if err := applyFilter(r.db.WithContext(ctx), f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.Ticket
	err := applyFilter(r.db.WithContext(ctx), f).
		Preload("Requester").Preload("Assignee").
		Order(orderClause(f.Sort)).
		Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).
		Find(&items).Error
	return items, total, err
}

func (r *ticketRepo) Update(ctx context.Context, t *models.Ticket, expectedStatus models.TicketStatus) error {
	var assignee any
	if t.AssigneeID != nil {
		assignee = *t.AssigneeID
	}
	res := r.db.WithContext(ctx).Model(&models.Ticket{}).
		Where("id = ? AND status = ?", t.ID, expectedStatus).
		Updates(map[string]any{
			"status":      t.Status,
			"priority":    t.Priority,
			"assignee_id": assignee,
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrConflict
	}
	return nil
}
