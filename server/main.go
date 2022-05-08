package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Health Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	return r
}

func main() {
	r := setupRouter()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}
	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	r.Run(fmt.Sprintf(":%s", PORT))
}
