package repository

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	List(ctx context.Context, f UserFilter) ([]models.User, int64, error)
	SetRole(ctx context.Context, id uuid.UUID, role models.Role) (*models.User, error)
	ListStaff(ctx context.Context) ([]models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicate
	}
	return err
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return wrap(&u, err)
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return wrap(&u, err)
}

func wrap(u *models.User, err error) (*models.User, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

type UserFilter struct {
	Role     models.Role
	Query    string
	Page     int
	PageSize int
}

func (r *userRepo) List(ctx context.Context, f UserFilter) ([]models.User, int64, error) {
	base := func() *gorm.DB {
		db := r.db.WithContext(ctx).Model(&models.User{})
		if f.Role != "" {
			db = db.Where("role = ?", f.Role)
		}
		if f.Query != "" {
			like := "%" + likeEscaper.Replace(f.Query) + "%"
			db = db.Where("(name ILIKE ? or email ILIKE ?)", like, like)
		}
		return db
	}
	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.User
	err := base().Order("created_at ASC, id ASC").Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).Find(&items).Error
	return items, total, err
}

func (r *userRepo) SetRole(ctx context.Context, id uuid.UUID, role models.Role) (*models.User, error) {
	var out models.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		lock := clause.Locking{Strength: "UPDATE"}
		var admins []models.User
		if err := tx.Clauses(lock).Where("role = ?", models.RoleAdmin).Order("id").Find(&admins).Error; err != nil {
			return err
		}
		if err := tx.Clauses(lock).First(&out, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
		}
		if out.Role == models.RoleAdmin && role != models.RoleAdmin && len(admins) <= 1 {
			return ErrLastAdmin
		}
		return tx.Model(&out).Update("role", role).Error
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *userRepo) ListStaff(ctx context.Context) ([]models.User, error) {
	var out []models.User
	err := r.db.WithContext(ctx).
		Where("role IN ?", []string{string(models.RoleAgent), string(models.RoleAdmin)}).
		Order("name ASC, id ASC").
		Find(&out).Error
	return out, err
}
