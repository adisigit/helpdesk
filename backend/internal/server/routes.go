package server

import (
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
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
	ticketRepo := repository.NewticketReporitory(db)
	commentRepo := repository.NewCommentRepository(db)

	commentSvc := service.NewCommentService(ticketRepo, commentRepo)

	authH := handler.NewAuthHandler(service.NewAuthService(userRepo, jwt), ck)
	ticketH := handler.NewTicketHandler(service.NewTicketService(ticketRepo, userRepo), commentSvc)
	commentH := handler.NewCommentHandler(commentSvc)
	userH := handler.NewUserHandler(service.NewUserService(userRepo))

	authMW := middleware.Auth(api, jwt)
	protected := router.ProtectedGroup(api, authMW)
	adminOnly := router.ProtectedGroup(api, authMW, middleware.RequireRole(api, models.RoleAdmin))

	router.RegisterHealth(api)
	router.RegisterAuth(api, protected, authH)
	router.RegisterTickets(protected, ticketH)
	router.RegisterComments(protected, commentH)
	router.RegisterUsers(adminOnly, userH)
}
