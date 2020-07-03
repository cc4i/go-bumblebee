package apis

import (
	"air/aqi"
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockClient struct{}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func TestAirOfCity(t *testing.T) {

	tests := []struct {
		name     string
		env      string
		uri      string
		input    string
		output   string
		httpCode int
	}{
		{"Test normal response with 200",
			"127.0.0.1:1234",
			"/air/v1/city/beijing",
			`{"status":"ok","data":{"aqi":78,"idx":1451,"attributions":[{"url":"http://www.bjmemc.com.cn/","name":"Beijing Environmental Protection Monitoring Center (北京市环境保护监测中心)"},{"url":"https://waqi.info/","name":"quality Index Project"}],"city":{"geo":[39.954592,116.468117],"name":"Beijing (北京)","url":"https://aqicn.org/city/beijing"},"dominentpol":"pm25","iaqi":{"co":{"v":4.6},"h":{"v":47},"no2":{"v":7.4},"o3":{"v":70.1},"p":{"v":1002},"pm10":{"v":47},"pm25":{"v":78},"so2":{"v":0.6},"t":{"v":25},"w":{"v":0.5}},"time":{"s":"2020-05-22 18:00:00","tz":"+08:00","v":1590170400},"debug":{"sync":"2020-05-22T20:14:45+09:00"}}}`,
			`{"server_version":"v1","air_quality":{"index_city_v_hash":"bf2e795904dfd281bb96432e9ed400e3886d4e8d","index_city":"Beijing_1451","idx":1451,"aqi":78,"city":"Beijing","city_cn":"北京","lat":"39.9546","lng":"116.468","co":"4.6","h":"47","no2":"7.4","o3":"70.1","p":"1002","pm10":"47","pm25":"78","so2":"0.6","t":"25","w":"0.5","s":"2020-05-22 18:00:00","tz":"+08:00","v":1590170400}}`,
			200,
		},
		{"Test URI not exist with 404",
			"127.0.0.1:1234",
			"/city/beijing",
			``,
			``,
			404,
		},
	}

	for _, test := range tests {
		aqi.Client = &MockClient{}
		os.Setenv("AIR_SERVICE_ENDPOINT", test.env)
		body := ioutil.NopCloser(bytes.NewReader([]byte(test.input)))
		GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		}

		r := Router(context.TODO())
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.uri, nil)

		r.ServeHTTP(recorder, req)

		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.httpCode, recorder.Code)
			if test.httpCode == 200 {
				assert.Equal(t, test.output, recorder.Body.String())
			}
		})
	}

}
