package main

import (
	"github.com/CBYeuler/atlas-workspace/backend/internal/config"
	"github.com/CBYeuler/atlas-workspace/backend/internal/database"
	"github.com/CBYeuler/atlas-workspace/backend/internal/handlers"
	"github.com/CBYeuler/atlas-workspace/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()
	r := gin.Default()

	// Health Check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// Services
	authService := services.NewAuthService()
	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}
	}

	r.Run(":8080")
}
