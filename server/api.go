package main

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Hello
	router.GET("/", Hello)
	// Health Endpoint
	router.GET("/health", Health)
	// Trusted Knock
	router.POST("/trustedknock", TrustedHmacAuthentication(), TrustedKnock)

	return router
}
