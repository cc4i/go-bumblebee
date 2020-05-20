package httpsvr

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	opentracing "github.com/opentracing/opentracing-go"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
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

func HttpGet(ctx context.Context, url string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "HttpGet")
	defer span.Finish()

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
