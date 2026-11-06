package database

import (
	"helpdesk/backend/internal/database/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var seedUsers = []struct {
	Name, Email string
	Role        models.Role
}{
	{"Admin Helpdesk", "admin@helpdesk.com", models.RoleAdmin},
	{"Agent Helpdesk", "agent@helpdesk.com", models.RoleAgent},
	{"Requester Helpdesk", "requester@helpdesk.com", models.RoleRequester},
}

func Seed(db *gorm.DB, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	for _, s := range seedUsers {
		u := models.User{Name: s.Name, Email: s.Email, PasswordHash: string(hash), Role: s.Role}
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&u).Error
		if err != nil {
			return err
		}
	}
	return nil
}
