package server

import (
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/repository"
	"helpdesk/backend/internal/server/handler"
	"helpdesk/backend/internal/server/middleware"
	"helpdesk/backend/internal/server/router"
	"helpdesk/backend/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func registerRoutes(api huma.API, db *gorm.DB, jwt *auth.JWT, ck handler.CookieConfig) {
	userRepo := repository.NewUserRepository(db)
	authH := handler.NewAuthHandler(service.NewAuthService(userRepo, jwt), ck)

	protected := router.ProtectedGroup(api, middleware.Auth(api, jwt))

	router.RegisterHealth(api)
	router.RegisterAuth(api, protected, authH)
}
