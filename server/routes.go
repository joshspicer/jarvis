package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCommand(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}
