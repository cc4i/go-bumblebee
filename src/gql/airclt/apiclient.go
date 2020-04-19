package airclt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

func AirOfCity(ctx context.Context, city string)([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AirOfCity")
	defer span.Finish()

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	url := "http://" + endpoint + "/air/" + city

	body, err := HttpGet(ctx, url)
	if err != nil {
		err = errors.Wrapf(err, "Failed to call air service.")
		return body, err
	}
	return body, nil
}

func HttpGet(ctx context.Context, url string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "HttpGet")
	defer span.Finish()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := Client.Do(req)
	if err != nil {
		err = errors.Wrapf(err,"API call was failed with %s with Err: %s. ", url, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrapf(err,"Read buffer failed.")
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
