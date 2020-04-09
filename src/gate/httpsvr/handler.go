package httpsvr

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func AirOfCity(c *gin.Context) {
	city := c.Param("city")

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	url := "http://" + endpoint + "/air/" + city

	body, err := HttpGet(url)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(body))
}
