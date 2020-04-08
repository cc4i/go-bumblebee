package aqi

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)


type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HttpClient
)

type ApiError struct {
	Status string `json:"status"`
	Data string `json:"data"`

}

type AirQuality struct {
	IndexCityVHash string `json:"index_city_v_hash"`
	IndexCity string `json:"index_city"`
	StationIndex int `json:"idx"`
	AQI int `json:"aqi"`
	City string `json:"city"`
	CityCN string `json:"city_cn"`
	Latitude string `json:"lat"`
	Longitude string `json:"lng"`
	Co string `json:"co"`
	H string `json:"h"`
	No2 string `json:"no2"`
	O3 string `json:"o3"`
	P string `json:"p"`
	Pm10 string `json:"pm10"`
	Pm25 string `json:"pm25"`
	So2 string `json:"so2"`
	T string `json:"t"`
	W string `json:"w"`
	S string `json:"s"` //Local measurement time
	TZ string `json:"tz"` //Station timezone
	V int `json:"v"`
}


type OriginAirQuality struct {
	Status string `json:"status"`
	Data OriginData `json:"data"`
}

type OriginData struct {
	AQI int `json:"aqi"`
	StationIndex int `json:"idx"`
	City OriginCity `json:"city"`
	IAQI OriginIAQI	`json:"iaqi"`
	OriginTime OriginTime `json:"time"`
}

type OriginCity struct {
	Geo []float64 `json:"geo"`
	Name string `json:"name"`
}

type OriginIAQI struct {
	Co OValue `json:"co"`
	H OValue `json:"h"`
	No2 OValue `json:"no2"`
	O3 OValue `json:"o3"`
	P OValue `json:"p"`
	Pm10 OValue `json:"pm10"`
	Pm25 OValue `json:"pm25"`
	So2 OValue `json:"so2"`
	T OValue `json:"t"`
	W OValue `json:"w"`
}

type OValue struct {
	V float64 `json:"v"`
}

type OriginTime struct {
	S string `json:"s"` //Local measurement time
	TZ string `json:"tz"` //Station timezone
	V int `json:"v"`
}


func init() {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	Client = &http.Client{Transport: tr}

}

func SplitName(name string) (city string, citycn string) {
	ns := strings.Split(name, "(")
	if len(ns) !=2 {
		log.Println("Input name <", name, "> wasn't matched with convention. eg: ", "Beijing (北京)")
		return name, ""
	}
	city = strings.Trim(ns[0]," ")
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
	dest.Latitude = strconv.FormatFloat(src.Data.City.Geo[0], 'g',6,64)
	dest.Longitude = strconv.FormatFloat(src.Data.City.Geo[1], 'g',6,64)

	dest.Co = strconv.FormatFloat(src.Data.IAQI.Co.V, 'g',6,64)
	dest.H = strconv.FormatFloat(src.Data.IAQI.H.V, 'g',6,64)
	dest.No2 = strconv.FormatFloat(src.Data.IAQI.No2.V, 'g',6,64)
	dest.O3 = strconv.FormatFloat(src.Data.IAQI.O3.V, 'g',6,64)
	dest.P = strconv.FormatFloat(src.Data.IAQI.P.V, 'g',6,64)
	dest.Pm10 = strconv.FormatFloat(src.Data.IAQI.Pm10.V, 'g',6,64)
	dest.Pm25 = strconv.FormatFloat(src.Data.IAQI.Pm25.V, 'g',6,64)
	dest.So2 = strconv.FormatFloat(src.Data.IAQI.So2.V, 'g',6,64)
	dest.T = strconv.FormatFloat(src.Data.IAQI.T.V, 'g',6,64)
	dest.W = strconv.FormatFloat(src.Data.IAQI.W.V, 'g',6,64)

	dest.S = src.Data.OriginTime.S
	dest.TZ = src.Data.OriginTime.TZ
	dest.V = src.Data.OriginTime.V

	dest.IndexCity = "" + dest.City + "_" + strconv.Itoa(dest.StationIndex)

	h:=sha1.New()
	h.Write([]byte(dest.IndexCity + "_" + strconv.Itoa(dest.V)))
	dest.IndexCityVHash = hex.EncodeToString(h.Sum(nil))
	return dest
}

func AirOfCity(c *gin.Context) {
	city := c.Param("city")

	// TODO: Making more secure
	token := "b0e78ca32d058a9170b6907c5214c0e946534cc9"
	host := "https://api.waqi.info"
	url := host + "/feed/" + city + "/?token=" + token
	// ---

	body, err := HttpGet(url)
	if err !=nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, convertAir(body))

}



func convertAir(content []byte) AirQuality {
	var originAir OriginAirQuality
	var newAir AirQuality
	var apiError ApiError


	err := json.Unmarshal(content, &originAir)
	if err!=nil {
		log.Println(err)
	}
	if originAir.Status=="error" {
		err2 := json.Unmarshal(content, &apiError)
		if err2!=nil {
			log.Println(err2)
		}
		log.Printf("Convert data was failed due to <%s>. ",  apiError.Data)
		return newAir

	}
	newAir = Copy2AirQuality(originAir)


	return newAir

}

func HttpGet(url string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := Client.Do(req)
	if err !=nil {
		log.Printf("API call was failed from %s with Err: %s. \n", url, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read buffer failed.\n")
		return nil, err
	}

	return body, nil
}