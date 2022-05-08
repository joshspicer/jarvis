package main

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health Endpoint
	router.GET("/health", HealthCommand)

	return router
}
