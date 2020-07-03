package airclt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gql/graph/model"
	"io/ioutil"
	"net/http"
	"os"
)

type ResponseAirQuality struct {
	ServerVersion string         `json:"server_version"`
	Air           OriginAirQuality `json:"air_quality"`
}

type OriginAirQuality struct {
	IndexCityVHash string `json:"index_city_v_hash"`
	IndexCity      string `json:"index_city"`
	StationIndex   int    `json:"idx"`
	AQI            int    `json:"aqi"`
	City           string `json:"city"`
	CityCN         string `json:"city_cn"`
	Latitude       string `json:"lat"`
	Longitude      string `json:"lng"`
	Co             string `json:"co"`
	H              string `json:"h"`
	No2            string `json:"no2"`
	O3             string `json:"o3"`
	P              string `json:"p"`
	Pm10           string `json:"pm10"`
	Pm25           string `json:"pm25"`
	So2            string `json:"so2"`
	T              string `json:"t"`
	W              string `json:"w"`
	S              string `json:"s"`  //Local measurement time
	Tz             string `json:"tz"` //Station timezone
	V              int    `json:"v"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HttpClient
)

func init() {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	Client = &http.Client{Transport: tr}

}

func AirOfCity(ctx context.Context, city string) (*model.AirQuality, error) {
	span, sctx := opentracing.StartSpanFromContext(ctx, "AirOfCity")
	defer span.Finish()

	var air model.AirQuality
	var o ResponseAirQuality

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	url := "http://" + endpoint + "/air/city/" + city

	buf, err := HttpGet(sctx, url)
	if err != nil {
		log.Debugf("Call air service from  %s", url)
		err = errors.Wrapf(err, "Failed to call air service.")
		return &air, err
	}

	if err = json.Unmarshal(buf, &o); err != nil {
		log.Debugf("Unmarshal: %s", string(buf))
		err = errors.Wrapf(err, "Failed to unmarshal data.")
		return &air, err
	}
	// copy data from origin
	air = model.AirQuality{
		IndexCityVHash: o.Air.IndexCityVHash,
		IndexCity:      o.Air.IndexCity,
		StationIndex:   o.Air.StationIndex,
		Aqi:            o.Air.AQI,
		City:           o.Air.City,
		CityCn:         o.Air.CityCN,
		Latitude:       o.Air.Latitude,
		Longitude:      o.Air.Longitude,
		Co:             o.Air.Co,
		H:              o.Air.H,
		No2:            o.Air.No2,
		O3:             o.Air.O3,
		P:              o.Air.P,
		Pm10:           o.Air.Pm10,
		Pm25:           o.Air.Pm25,
		So2:            o.Air.So2,
		T:              o.Air.T,
		W:              o.Air.W,
		S:              o.Air.S,
		Tz:             o.Air.Tz,
		V:              o.Air.V,
	}

	return &air, nil
}

func HttpGet(ctx context.Context, url string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "HttpGet")
	defer span.Finish()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := Client.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "API call was failed with %s with Err: %s. ", url, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "Read buffer failed.")
		return nil, err
	}

	return body, nil
}

func ToJson(body []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Errorf("Parsed body with error: %s", err)
	}
	return prettyJSON.String()

}
