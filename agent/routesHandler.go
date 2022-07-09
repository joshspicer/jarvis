/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}
