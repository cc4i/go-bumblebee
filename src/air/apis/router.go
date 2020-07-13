package apis

import (
	v1 "air/apis/v1"
	v2 "air/apis/v2"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

var prefixV1 = "/air/v1"
var prefixV2 = "/air/v2"

// API routes/path definition
func Router(ctx context.Context) *gin.Engine {
	r := gin.Default()
	routeV1(ctx, r)
	routeV2(ctx, r)

	// Basic health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Metrics for Prometheus
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)

	})

	return r
}

func routeV1(ctx context.Context, r *gin.Engine) {
	// Get air quality data of given city
	// Example: ?c="city"
	r.GET(prefixV1+"/city/:city", func(c *gin.Context) {
		v1.AirOfCity(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	//
	// Destination: /feed/here/?token=:token
	//				api.ipstack.com/:ip?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	r.GET(prefixV1+"/ip/:ip", func(c *gin.Context) {
		v1.AirOfIP(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	// Example: ?lat=&lng=
	//
	// Destination: /feed/geo::lat;:lng/?token=:token
	r.GET(prefixV1+"/geo/:lat/:lng", func(c *gin.Context) {
		v1.AirOfGeo(ctx, c)
	})

	// Get all the stations within a given lat/lng bounds
	// Example: ?lat1=39.379436&lng1=116.091230&lat2=40.235643&lng2=116.784382
	//
	// Destination: /map/bounds/?token=:token&latlng=:latlng
	r.POST(prefixV1+"/station/bounds", func(c *gin.Context) {
		//TODO
		//lat1:lng1:lat2:lng2
		c.String(http.StatusHTTPVersionNotSupported, "Working in progress")
	})

	// Search stations by city name
	// Example: ?c="city"
	//
	// Destination: /search/?keyword=:keyword&token=:token
	r.POST(prefixV1+"/station/city", func(c *gin.Context) {
		//TODO
		c.String(http.StatusHTTPVersionNotSupported, "Working in progress")
	})

	// Get AQIServer standard: Air Quality Index scale as defined by the US-EPA 2016 standard
	r.GET(prefixV1+"/aqi", func(c *gin.Context) {
		c.String(http.StatusOK, `
			{ 
				"version": "v1",
				{
					"Standard": "Air Quality Index scale as defined by the US-EPA 2016 standard.",
					"Definitions": [
						{
							"AQIServer": "0-50",
							"Level": "Good",
							"Implication": "Air quality is considered satisfactory, and air pollution poses little or no risk",
							"Caution": "None"
						},
						{
							"AQIServer": "51 -100",
							"Level": "Moderate",
							"Implication": "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "101-150",
							"Level": "Unhealthy for Sensitive Groups",
							"Implication": "Members of sensitive groups may experience health effects. The general public is not likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "151-200",
							"Level": "Unhealthy",
							"Implication": "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion"
						},
						{
							"AQIServer": "201-300",
							"Level": "Very Unhealthy",
							"Implication": "Health warnings of emergency conditions. The entire population is more likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion."
						},
						{
							"AQIServer": "300+",
							"Level": "Hazardous",
							"Implication": "Health alert: everyone may experience more serious health effects",
							"Caution": "Everyone should avoid all outdoor exertion"
						}
					]
				}
			}
		`)
	})

	// Version from ENV for test purpose
	r.GET(prefixV1+"/version", func(c *gin.Context) {
		version := os.Getenv("AIR_VERSION")
		if version == "" {
			version = "v0.0.0"
		}
		c.String(http.StatusOK, version)
	})
}

func routeV2(ctx context.Context, r *gin.Engine) {
	r.GET(prefixV2+"/city/:city", func(c *gin.Context) {
		v2.AirOfCity(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	//
	// Destination: /feed/here/?token=:token
	//				api.ipstack.com/:ip?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	r.GET(prefixV2+"/ip/:ip", func(c *gin.Context) {
		v2.AirOfIP(ctx, c)
	})

	// Get the nearest station close to the user location, based on the IP information.
	// Example: ?lat=&lng=
	//
	// Destination: /feed/geo::lat;:lng/?token=:token
	r.GET(prefixV2+"/geo/:lat/:lng", func(c *gin.Context) {
		v2.AirOfGeo(ctx, c)
	})

	r.GET(prefixV2+"/aqi", func(c *gin.Context) {
		c.String(http.StatusOK, `
			{ 
				"version": "v2",
				{
					"Standard": "Air Quality Index scale as defined by the US-EPA 2016 standard.DEMO HAHA!!!",
					"Definitions": [
						{
							"AQIServer": "0-50",
							"Level": "Good",
							"Implication": "Air quality is considered satisfactory, and air pollution poses little or no risk",
							"Caution": "None"
						},
						{
							"AQIServer": "51 -100",
							"Level": "Moderate",
							"Implication": "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "101-150",
							"Level": "Unhealthy for Sensitive Groups",
							"Implication": "Members of sensitive groups may experience health effects. The general public is not likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "151-200",
							"Level": "Unhealthy",
							"Implication": "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion"
						},
						{
							"AQIServer": "201-300",
							"Level": "Very Unhealthy",
							"Implication": "Health warnings of emergency conditions. The entire population is more likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion."
						},
						{
							"AQIServer": "300+",
							"Level": "Hazardous",
							"Implication": "Health alert: everyone may experience more serious health effects",
							"Caution": "Everyone should avoid all outdoor exertion"
						}
					]
				}
			}
		`)
	})

	// Version from ENV for test purpose
	r.GET(prefixV2+"/version", func(c *gin.Context) {
		version := os.Getenv("AIR_VERSION")
		if version == "" {
			version = "v0.0.0"
		}
		c.String(http.StatusOK, version)
	})

}
