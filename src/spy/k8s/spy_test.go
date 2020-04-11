package k8s

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestPing(t *testing.T) {
	//content := "ping"
	//body := ioutil.NopCloser(bytes.NewReader([]byte(content)))

	r := Router()
	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ping"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}

	mt, response, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message over ws connection %v", err)
	}

	assert.Equal(t, websocket.TextMessage, mt)
	assert.Equal(t, "pong", string(response))
}