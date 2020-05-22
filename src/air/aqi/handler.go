package aqi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var (
	RedisServer         string
	RedisServerPassword string
	RedisClient         *redis.Client
)

func init() {
	RedisServer = os.Getenv("REDIS_SERVER_ADDRESS")
	RedisServerPassword = os.Getenv("REDIS_SERVER_PASSWORD")
	if RedisServer == "" {
		log.Error("Initial redis server address was failed.")
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     RedisServer,
			Password: RedisServerPassword,
		})
		pong, err := RedisClient.Ping(context.TODO()).Result()
		if err != nil {
			log.Error(pong, err)
		}
	}
}

func CachingAirQuality(ctx context.Context, city string, quality AirQuality) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CachingAirQuality")
	defer span.Finish()

	buf, _ := json.Marshal(&quality)
	if RedisClient == nil {
		return errors.New("redis client was nil")
	}
	pipeline := RedisClient.Pipeline()
	pipeline.Set(ctx, "air_quality_cache-"+city, "expired-600s", 600*time.Second)
	pipeline.HSet(ctx, "air_quality_cache", city, buf)
	_, err := pipeline.Exec(ctx)

	return err
}

func CachedAirQuality(ctx context.Context, city string) (AirQuality, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CachedAirQuality")
	defer span.Finish()

	var air = AirQuality{}

	if RedisClient == nil {
		return air, errors.New("redis client was nil")
	}
	r, err := RedisClient.Get(ctx, "air_quality_cache-"+city).Result()
	log.Infof("air_quality_cache-%s = %s", city, r)
	if err == nil {
		buf, err := RedisClient.HGet(ctx, "air_quality_cache", city).Bytes()
		if err != nil {
			return air, err
		}
		err = json.Unmarshal(buf, &air)
	} else {
		err = errors.New("Cached Air Quality of " + city + "has been expired.")
	}

	return air, err

}

func AirOfGeo(ctx context.Context, c *gin.Context) {
	///air/geo/:lat/:lng ->//feed/geo::lat;:lng/?token=:token
	//Auckland: -36.916839599609375, 174.70875549316406
	span, _ := opentracing.StartSpanFromContext(ctx, "http-AirOfGeo")
	defer span.Finish()

	lat := c.Param("lat")
	lng := c.Param("lng")
	url := AQIServer + "/feed/geo:" + lat + ";" + lng + "/?token=" + AQIServerToken
	if buf, err := HttpGet(ctx, url); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)

	} else {
		if air, err := convertAir(buf); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, air)
		}
	}
}

func byCity(ctx context.Context, city string) (AirQuality, error) {

	url := AQIServer + "/feed/" + city + "/?token=" + AQIServerToken
	// ---
	buf, err := HttpGet(ctx, url)
	if err != nil {
		log.Errorf("Fail to call AQIServer service from %s", url)
		return AirQuality{}, err
	}
	return convertAir(buf)
}

func AirOfCity(ctx context.Context, c *gin.Context) {
	span, sctx := opentracing.StartSpanFromContext(ctx, "http-AirOfCity")
	defer span.Finish()

	city := c.Param("city")

	air, err := CachedAirQuality(sctx, city)
	if err != nil {
		log.Error(err)
		log.Infof("No cache for %s and looking for fresh value.\n ", city)
		air, err := byCity(sctx, city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			if err := CachingAirQuality(sctx, city, air); err != nil {
				log.Errorf("Caching air quality data was failed. -> %s, %s\n", air.City, err)
			}
			log.Infof("Air Quality of %s was cached.\n ", city)
			c.JSON(http.StatusOK, air)
		}

	} else {
		log.Infof("Return cached Air Quality of %s.\n ", city)
		c.JSON(http.StatusOK, air)
	}

}

func convertAir(content []byte) (AirQuality, error) {
	var originAir OriginAirQuality
	var newAir AirQuality
	var apiError ApiError

	err := json.Unmarshal(content, &originAir)
	if err != nil {
		log.Println(err)
		return newAir, err
	}
	if originAir.Status == "error" {
		err = json.Unmarshal(content, &apiError)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Convert data was failed due to <%s>. ", apiError.Data)
		return newAir, err

	}
	newAir = Copy2AirQuality(originAir)

	return newAir, nil

}

func AirOfIP(ctx context.Context, c *gin.Context) {
	span, sctx := opentracing.StartSpanFromContext(ctx, "http-AirOfIP")
	defer span.Finish()

	ip := c.Param("ip")
	if city, err := CityByIP(sctx, ip); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		if air, err := byCity(sctx, city); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, air)
		}

	}
}
