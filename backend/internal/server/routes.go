package server

import (
	"helpdesk/backend/internal/server/router"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func registerRoutes(api huma.API, db *gorm.DB) {
	_ = db
	router.RegisterHealth(api)
}
