package aqi

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	Client             HttpClient
	AQIServer          string // "https://api.waqi.info"
	AQIServerToken     string // "b0e78ca32d058a9170b6907c5214c0e946534cc9"
	IpStackServer      string // "http://api.ipstack.com"
	IpStackServerToken string // "ad7c6834f8dba51e8943d96d3742fcc5"

	//api.ipstack.com/127.0.0.1?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	//https://ipapi.co/json
	//https://ipapi.co/8.8.8.8/json
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type IpStack struct {
	Ip            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Cip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

type ApiError struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type AirQuality struct {
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
	TZ             string `json:"tz"` //Station timezone
	V              int    `json:"v"`
}

type OriginAirQuality struct {
	Status string     `json:"status"`
	Data   OriginData `json:"data"`
}

type OriginData struct {
	AQI          int        `json:"aqi"`
	StationIndex int        `json:"idx"`
	City         OriginCity `json:"city"`
	IAQI         OriginIAQI `json:"iaqi"`
	OriginTime   OriginTime `json:"time"`
}

type OriginCity struct {
	Geo  []float64 `json:"geo"`
	Name string    `json:"name"`
}

type OriginIAQI struct {
	Co   OValue `json:"co"`
	H    OValue `json:"h"`
	No2  OValue `json:"no2"`
	O3   OValue `json:"o3"`
	P    OValue `json:"p"`
	Pm10 OValue `json:"pm10"`
	Pm25 OValue `json:"pm25"`
	So2  OValue `json:"so2"`
	T    OValue `json:"t"`
	W    OValue `json:"w"`
}

type OValue struct {
	V float64 `json:"v"`
}

type OriginTime struct {
	S  string `json:"s"`  //Local measurement time
	TZ string `json:"tz"` //Station timezone
	V  int    `json:"v"`
}

// Log handle for third parties
var FL = log.New()

func init() {
	// Initial ENVs
	AQIServer = os.Getenv("AQI_SERVER_URL")
	AQIServerToken = os.Getenv("AQI_SERVER_TOKEN")
	IpStackServer = os.Getenv("IP_STACK_SERVER_URL")
	IpStackServerToken = os.Getenv("IP_STACK_SERVER_TOKEN")
	if AQIServer == "" || IpStackServer == "" {
		log.Fatal("Retrieving servers' address were failed. Check out environments setting & Reboot.")
	}
	if AQIServerToken == "" || IpStackServerToken == "" {
		log.Fatal("Retrieving servers' token were failed. Check out environments setting & Reboot.")
	}

	// Initial http client
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	Client = &http.Client{Transport: tr}

	// Initial file logs for third parties
	fl, err := os.OpenFile("third-parties.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Initial logs file for third parties was failed.")
	}
	FL.SetOutput(fl)
	FL.SetFormatter(&log.JSONFormatter{})
}

func SplitName(name string) (city string, citycn string) {
	ns := strings.Split(name, "(")
	if len(ns) != 2 {
		log.Println("Input name <", name, "> wasn't matched with convention. eg: ", "Beijing (北京)")
		return name, ""
	}
	city = strings.Trim(ns[0], " ")
	citycn = strings.Trim(ns[1], ")")
	return city, citycn
}

func Copy2AirQuality(src OriginAirQuality) AirQuality {

	var dest AirQuality
	dest.StationIndex = src.Data.StationIndex
	dest.AQI = src.Data.AQI
	c, cn := SplitName(src.Data.City.Name)
	dest.City = c
	dest.CityCN = cn
	dest.Latitude = strconv.FormatFloat(src.Data.City.Geo[0], 'g', 6, 64)
	dest.Longitude = strconv.FormatFloat(src.Data.City.Geo[1], 'g', 6, 64)

	dest.Co = strconv.FormatFloat(src.Data.IAQI.Co.V, 'g', 6, 64)
	dest.H = strconv.FormatFloat(src.Data.IAQI.H.V, 'g', 6, 64)
	dest.No2 = strconv.FormatFloat(src.Data.IAQI.No2.V, 'g', 6, 64)
	dest.O3 = strconv.FormatFloat(src.Data.IAQI.O3.V, 'g', 6, 64)
	dest.P = strconv.FormatFloat(src.Data.IAQI.P.V, 'g', 6, 64)
	dest.Pm10 = strconv.FormatFloat(src.Data.IAQI.Pm10.V, 'g', 6, 64)
	dest.Pm25 = strconv.FormatFloat(src.Data.IAQI.Pm25.V, 'g', 6, 64)
	dest.So2 = strconv.FormatFloat(src.Data.IAQI.So2.V, 'g', 6, 64)
	dest.T = strconv.FormatFloat(src.Data.IAQI.T.V, 'g', 6, 64)
	dest.W = strconv.FormatFloat(src.Data.IAQI.W.V, 'g', 6, 64)

	dest.S = src.Data.OriginTime.S
	dest.TZ = src.Data.OriginTime.TZ
	dest.V = src.Data.OriginTime.V

	dest.IndexCity = "" + dest.City + "_" + strconv.Itoa(dest.StationIndex)

	h := sha1.New()
	h.Write([]byte(dest.IndexCity + "_" + strconv.Itoa(dest.V)))
	dest.IndexCityVHash = hex.EncodeToString(h.Sum(nil))
	return dest

}

func HttpGet(ctx context.Context, url string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "HttpGet")
	defer span.Finish()

	log.Printf("Request to : %s\n", url)
	FL.WithFields(log.Fields{
		"request_to": url,
	}).Info("Third party url")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := Client.Do(req)
	if err != nil {
		log.Printf("API call was failed from %s with Err: %s. \n", url, err)
		FL.WithFields(log.Fields{
			"request_to": url,
			"error":      err.Error(),
		}).Error("Request failed")
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read buffer failed.\n")
		FL.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Read buffer failed")
		return nil, err
	}
	log.Printf("Response from %s : %s\n", url, string(body))
	FL.WithFields(log.Fields{
		"response_from": url,
		"origin_body":   string(body),
	}).Info("Third party response")

	return body, nil
}

func CityByIP(ctx context.Context, ip string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CityByIP")
	defer span.Finish()

	var ipStack IpStack
	url := IpStackServer + "/" + ip + "?access_key=" + IpStackServerToken
	buf, err := HttpGet(ctx, url)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if err = json.Unmarshal(buf, &ipStack); err != nil {
		log.Error(err)
		return "", err
	}
	return ipStack.City, nil

}
