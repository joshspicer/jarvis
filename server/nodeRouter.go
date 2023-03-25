package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/serial"
)

func NodeRouter() *gin.Engine {
	router := gin.Default()

	// Static
	router.StaticFile("robots.txt", "./static/robots.txt")

	// router.GET("/", Hello)
	router.GET("/health", health)
	router.GET("/info", Auth(), getNodeInfo)

	return router
}

func getNodeInfo(c *gin.Context) {

	// ioWR := c.MustGet("ioWR").(io.ReadWriter)

	// Make new IO ReadWriter
	m, err := serial.New(serial.WithPort("/dev/ttyUSB2"), serial.WithBaud(115200))
	if err != nil {
		log.Println(err)
		return
	}
	defer m.Close()
	var mio io.ReadWriter = m

	modem := at.New(mio, at.WithTimeout(2*time.Second))
	info, err := modem.Command("I")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting node info",
			"error":   err.Error(),
		})
		return
	}
	signalStrength, err := modem.Command("+CSQ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting signal strength",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Hello World",
		"modemInfo": info,
		"signal":    signalStrength,
	})
}
