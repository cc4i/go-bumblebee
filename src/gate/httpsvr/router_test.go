package httpsvr

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
