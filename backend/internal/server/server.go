package server

import (
	"helpdesk/backend/internal/config"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *http.Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	api := humagin.New(r, huma.DefaultConfig("Helpdesk API", "1.0.0"))
	registerRoutes(api, db)

	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
