package httpsvr

import (
	"context"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Ping(ctx context.Context, c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Ping")
	defer span.Finish()

	c.String(http.StatusOK, "pong")
}

func AirOfCity(ctx context.Context, c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AirOfCity")
	defer span.Finish()

	city := c.Param("city")

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	url := "http://" + endpoint + "/air/" + city

	body, err := HttpGet(ctx, url)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(body))
}
