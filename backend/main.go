package main

import (
	"github.com/CBYeuler/atlas-workspace/backend/internal/config"
	"github.com/CBYeuler/atlas-workspace/backend/internal/database"
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

	r.Run(":8080")
}
