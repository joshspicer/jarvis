/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"github.com/gin-gonic/gin"
)

func SetupAgentRouter() *gin.Engine {
	router := gin.Default()

	// Health Endpoint
	router.GET("/health", Health)

	return router
}
