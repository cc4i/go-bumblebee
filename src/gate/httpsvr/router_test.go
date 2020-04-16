package httpsvr

import (
	"bytes"
	"context"
	"gate/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Test return http_code", "", "200"},
		{"Test the content of response", "", "pong"},
	}

	r := Router(context.TODO())
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(recorder, req)

	t.Run(tests[0].name, func(t *testing.T) {
		code, _ := strconv.Atoi(tests[0].output)
		assert.Equal(t, code, recorder.Code)
	})

	t.Run(tests[1].name, func(t *testing.T) {
		assert.Equal(t, tests[1].output, recorder.Body.String())
	})

}

func TestAirOfCity(t *testing.T) {

	content := `{"index_city_v_hash":"5a367bf029843359937b1830a85970b175faffea","index_city":"Beijing_1451","idx":1451,"aqi":63,"city":"Beijing","city_cn":"北京","lat":"39.9546","lng":"116.468","co":"2.8","h":"32","no2":"10.1","o3":"26.4","p":"1020","pm10":"57","pm25":"63","so2":"1.6","t":"11","w":"2.5","s":"2020-04-07 22:00:00","tz":"+08:00","v":1586296800}`
	os.Setenv("AIR_SERVICE_ENDPOINT", "127.0.0.1")

	Client = &mocks.MockClient{}
	body := ioutil.NopCloser(bytes.NewReader([]byte(content)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}

	r := Router(context.TODO())
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/air/beijing", nil)

	r.ServeHTTP(recorder, req)

	t.Run("Test return http_code = 200", func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("Test return content", func(t *testing.T) {
		assert.Equal(t, content, recorder.Body.String())
	})
}
