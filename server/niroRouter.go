package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NiroRouter() *gin.Engine {
	router := gin.Default()

	// Static
	router.StaticFile("robots.txt", "./static/robots.txt")

	// router.GET("/", Hello)
	router.GET("/health", health)
	router.GET("/info", Auth(), getNiroInfo)

	return router
}

func getNiroInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
