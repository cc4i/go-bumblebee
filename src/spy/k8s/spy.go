package k8s

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

type WebContext struct{}

type ServiceHandlers interface {
	Handler(path string, c *gin.Context)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var k8s *K8sContext

func initK8s() {
	if k8s == nil {
		clientset, err := kubernetes.NewForConfig(getConfig())
		if err != nil {
			log.Error(err.Error())
		}
		k8s = &K8sContext{Clientset: clientset}
	}

}

func (web *WebContext) Handler(path string, c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	switch path {
	case "/wsping":
		WSPing(ws)
		break
	case "/spy":
		Spy(ws)
		break

	}

}

func WSPing(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = ws.WriteMessage(mt, []byte("pong"))
		if err != nil {
			log.Println("write:", err)
			break
		}

	}
}

func Spy(ws *websocket.Conn) {
	initK8s()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		command := string(message)
		var buf bytes.Buffer
		switch command {
		case "ns":
			buf.WriteString(k8s.GetNamespaces())
			break
		default:
			buf.WriteString(" < "+command + "> is not valid.")
		}
		err = ws.WriteMessage(mt, buf.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}

	}

}
