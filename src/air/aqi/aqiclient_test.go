package aqi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
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
			"/air/city/beijing",
			`{"status":"ok","data":{"aqi":63,"idx":1451,"attributions":[{"url":"http://www.bjmemc.com.cn/","name":"Beijing Environmental Protection Monitoring Center (北京市环境保护监测中心)"},{"url":"https://china.usembassy-china.org.cn/embassy-consulates/beijing/air-quality-monitor/","name":"U.S Embassy Beijing Air Quality Monitor (美国驻北京大使馆空气质量监测)"},{"url":"https://waqi.info/","name":"World Air Quality Index Project"}],"city":{"geo":[39.954592,116.468117],"name":"Beijing (北京)","url":"https://aqicn.org/city/beijing"},"dominentpol":"pm25","iaqi":{"co":{"v":4.6},"h":{"v":19},"no2":{"v":5.5},"o3":{"v":37.8},"p":{"v":1020},"pm10":{"v":56},"pm25":{"v":63},"so2":{"v":3.6},"t":{"v":15},"w":{"v":3.6}},"time":{"s":"2020-04-08 17:00:00","tz":"+08:00","v":1586365200},"debug":{"sync":"2020-04-08T18:28:14+09:00"}}}`,
			`{"index_city_v_hash":"ce8df35a1ae16beefc8d8a45be6d3a4ac224e008","index_city":"Beijing_0","idx":0,"aqi":0,"city":"Beijing","city_cn":"北京","lat":"39.9546","lng":"116.468","co":"4.6","h":"19","no2":"5.5","o3":"37.8","p":"1020","pm10":"56","pm25":"63","so2":"3.6","t":"15","w":"3.6","s":"","tz":"","v":0}`,
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
		Client = &MockClient{}
		os.Setenv("AIR_SERVICE_ENDPOINT", test.env)
		body := ioutil.NopCloser(bytes.NewReader([]byte(test.input)))
		GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		}

		r := Router()
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

func TestHttpGet(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"Test response with plain text", "Hello Go!"},
		{"Test response with json", `{"index_city_v_hash":"5a367bf029843359937b1830a85970b175faffea","index_city":"Beijing_1451","idx":1451,"aqi":63,"city":"Beijing","city_cn":"北京","lat":"39.9546","lng":"116.468","co":"2.8","h":"32","no2":"10.1","o3":"26.4","p":"1020","pm10":"57","pm25":"63","so2":"1.6","t":"11","w":"2.5","s":"2020-04-07 22:00:00","tz":"+08:00","v":1586296800}`},
	}
	for _, test := range tests {
		Client = &http.Client{}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.content)
		}))
		defer ts.Close()

		body, err := HttpGet(ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.content+"\n", string(body))
		})
	}

}

func TestCopy2AirQuality(t *testing.T) {
	var originAir OriginAirQuality

	content := `{"status":"ok","data":{"aqi":63,"idx":1451,"attributions":[{"url":"http://www.bjmemc.com.cn/","name":"Beijing Environmental Protection Monitoring Center (北京市环境保护监测中心)"},{"url":"https://china.usembassy-china.org.cn/embassy-consulates/beijing/air-quality-monitor/","name":"U.S Embassy Beijing Air Quality Monitor (美国驻北京大使馆空气质量监测)"},{"url":"https://waqi.info/","name":"World Air Quality Index Project"}],"city":{"geo":[39.954592,116.468117],"name":"Beijing (北京)","url":"https://aqicn.org/city/beijing"},"dominentpol":"pm25","iaqi":{"co":{"v":4.6},"h":{"v":19},"no2":{"v":5.5},"o3":{"v":37.8},"p":{"v":1020},"pm10":{"v":56},"pm25":{"v":63},"so2":{"v":3.6},"t":{"v":15},"w":{"v":3.6}},"time":{"s":"2020-04-08 17:00:00","tz":"+08:00","v":1586365200},"debug":{"sync":"2020-04-08T18:28:14+09:00"}}}`

	json.Unmarshal([]byte(content), &originAir)

	newAir := Copy2AirQuality(originAir)

	assert.Equal(t, "Beijing", newAir.City)
	assert.Equal(t, "北京", newAir.CityCN)
	assert.Equal(t, "63", newAir.Pm25)

}

func TestSplitName(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		output   string
		outputCN string
	}{
		{"Test normal input", "Beijing (北京)", "Beijing", "北京"},
		{"Test input without Chinese", "Beijing", "Beijing", ""},
		{"Test input without space", "Beijing(北京)", "Beijing", "北京"},
		{"Test input with half bracket", "Beijing(北京", "Beijing", "北京"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, outputCN := SplitName(test.input)
			assert.Equal(t, output, test.output)
			assert.Equal(t, outputCN, test.outputCN)
		})
	}
}
