package aqi

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client             HttpClient
	AQIServer          = "https://api.waqi.info"
	AQIServerToken     = "b0e78ca32d058a9170b6907c5214c0e946534cc9"
	IpStackServer      = "http://api.ipstack.com"
	IpStackServerToken = "ad7c6834f8dba51e8943d96d3742fcc5"

	//api.ipstack.com/127.0.0.1?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	//https://ipapi.co/json
	//https://ipapi.co/8.8.8.8/json
)

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
	Status string     `json:"Status"`
	Data   OriginData `json:"Data"`
}

type OriginData struct {
	AQI          int        `json:"AQIServer"`
	StationIndex int        `json:"StationIndex"`
	City         OriginCity `json:"City"`
	IAQI         OriginIAQI `json:"IAQI"`
	OriginTime   OriginTime `json:"OriginTime"`
}

type OriginCity struct {
	Geo  []float64 `json:"Geo"`
	Name string    `json:"Name"`
}

type OriginIAQI struct {
	Co   OValue `json:"Co"`
	H    OValue `json:"H"`
	No2  OValue `json:"No2"`
	O3   OValue `json:"O3"`
	P    OValue `json:"P"`
	Pm10 OValue `json:"Pm10"`
	Pm25 OValue `json:"Pm25"`
	So2  OValue `json:"So2"`
	T    OValue `json:"T"`
	W    OValue `json:"W"`
}

type OValue struct {
	V float64 `json:"V"`
}

type OriginTime struct {
	S  string `json:"S"`  //Local measurement time
	TZ string `json:"TZ"` //Station timezone
	V  int    `json:"V"`
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

func AirOfGeo(c *gin.Context) {
	///air/geo/:lat/:lng ->//feed/geo::lat;:lng/?token=:token
	//Auckland: -36.916839599609375, 174.70875549316406

	lat := c.Param("lat")
	lng := c.Param("lng")
	url := AQIServer + "/feed/geo:" + lat + ";" + lng + "/?token=" + AQIServerToken
	if buf, err := HttpGet(url); err != nil {
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
func AirOfIP(c *gin.Context) {
	//http://api.ipstack.com/127.0.0.1?access_key=ad7c6834f8dba51e8943d96d3742fcc5
	var ipStack IpStack
	ip := c.Param("ip")
	url := IpStackServer + "/" + ip + "?access_key=" + IpStackServerToken
	if buf, err := HttpGet(url); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)

	} else {
		if err = json.Unmarshal(buf, &ipStack); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err)

		} else {

			if air, err := byCity(ipStack.City); err != nil {
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, air)
			}
		}
	}
}

func byCity(city string) (*AirQuality, error) {
	// TODO: Making more secure
	url := AQIServer + "/feed/" + city + "/?token=" + AQIServerToken
	// ---
	buf, err := HttpGet(url)
	if err != nil {
		log.Errorf("Fail to call AQIServer service from %s", url)
		return &AirQuality{}, err
	}

	return convertAir(buf)

}

func AirOfCity(c *gin.Context) {
	city := c.Param("city")
	air, err := byCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, air)
	}

}

func convertAir(content []byte) (*AirQuality, error) {
	var originAir OriginAirQuality
	var newAir AirQuality
	var apiError ApiError

	err := json.Unmarshal(content, &originAir)
	if err != nil {
		log.Println(err)
		return &newAir, err
	}
	if originAir.Status == "error" {
		err = json.Unmarshal(content, &apiError)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Convert data was failed due to <%s>. ", apiError.Data)
		return &newAir, err

	}
	newAir = Copy2AirQuality(originAir)

	return &newAir, nil

}

func HttpGet(url string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := Client.Do(req)
	if err != nil {
		log.Printf("API call was failed from %s with Err: %s. \n", url, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read buffer failed.\n")
		return nil, err
	}
	log.Println("origin response : ", string(body))

	return body, nil
}
