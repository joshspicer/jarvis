package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	// Make new IO ReadWriter
	m, err := serial.New(serial.WithPort("/dev/ttyUSB2"), serial.WithBaud(115200))
	if err != nil {
		log.Println(err)
		return
	}
	defer m.Close()
	var mio io.ReadWriter = m

	modem := at.New(mio, at.WithTimeout(5*time.Second))

	//var info string = "ERR"
	var signalStrength string = "ERR"
	var gpsLatitude float32 = 12.517572
	var gpsLongitude float32 = -69.9649462

	//infoArr, err := modem.Command("I")
	//if err == nil {
	// info = infoArr[0]
	//}

	// +CSQ: 23,99
	signalStrengthArr, err := modem.Command("+CSQ")
	if err == nil {
		signalStrength = strings.Join(signalStrengthArr, " ")
	}

	// +QGPSLOC: 174111.0,4736.8397N,12218.8860W,1.5,98.7,3,224.01,0.0,0.0,260323,05
	gpsLocationArr, err := modem.Command("+QGPSLOC?")
	if err == nil {
		// Convert AT latitude and longitude response
		var tmpLatitude = convertGpsLocation(gpsLocationArr[1])
		var tmpLongitude = convertGpsLocation(gpsLocationArr[2])

		if tmpLatitude != -1 && tmpLongitude != -1 {
			gpsLatitude = tmpLatitude
			gpsLongitude = tmpLongitude
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"modemInfo":          "modemInfo",
		"signal":             signalStrength,
		"accessoriesBattery": -1, // Mock
		"gpsLatitude":        gpsLatitude,
		"gpsLongitude":       gpsLongitude,
	})
}

func convertGpsLocation(gpsLocation string) float32 {
	// 4736.8397N
	// 12218.8860W

	// Remove and store the last character
	compassDirection := gpsLocation[len(gpsLocation)-1:]
	gpsLocation = gpsLocation[:len(gpsLocation)-1]

	// Move the decimal point two places to the left
	// 4736.8397 -> 47.368397
	// 12218.8860 -> 122.18886
	gpsLocation = gpsLocation[:len(gpsLocation)-2] + "." + gpsLocation[len(gpsLocation)-2:]
	// If compass direction is South or West, make the number negative
	if compassDirection == "S" || compassDirection == "W" {
		gpsLocation = "-" + gpsLocation
	}

	// Convert to float
	gpsLocationFloat, err := strconv.ParseFloat(gpsLocation, 32)
	if err != nil {
		log.Println(err)
		return -1
	}
	return float32(gpsLocationFloat)
}
