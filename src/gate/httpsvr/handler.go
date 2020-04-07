package httpsvr

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)



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

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}


func AirOfCity(c *gin.Context) {
	city := c.Param("city")

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	url := "http://" + endpoint + "/air/" + city

	body, err := HttpGet(url)
	if err !=nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(body))

}

func HttpGet(url string) ([]byte, error) {

	//config := &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	//tr := &http.Transport{TLSClientConfig: config}
	//client := &http.Client{Transport: tr}
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


func ToJson( body []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Println("JSON parse error: ", err)
	}
	return prettyJSON.String()

}