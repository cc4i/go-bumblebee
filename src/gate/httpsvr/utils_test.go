package httpsvr

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpGet(t *testing.T) {

	content := "Hello, client"
	Client = &http.Client{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, content)
	}))
	defer ts.Close()

	body, err := HttpGet(context.TODO(), ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Test HttpGet to get local content", func(t *testing.T) {
		assert.Equal(t, content+"\n", string(body))
	})
}

func TestToJson(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Test simple json", `{"name":"Rio Chen", "address":"Planet B, NZ"}`, "{\n\t\"name\": \"Rio Chen\",\n\t\"address\": \"Planet B, NZ\"\n}"},
		{"Test abnormal json", `{name:"Rio Chen", address:"Planet B, NZ"}`, ""},
	}

	for _, test := range tests {
		log.Println(test.output)
		t.Run(test.name, func(t *testing.T) {
			pretty := ToJson([]byte(test.input))
			assert.Equal(t, test.output, pretty)
		})
	}
}
