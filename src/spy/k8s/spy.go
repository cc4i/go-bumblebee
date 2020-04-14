package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

func initK8s() error {
	if k8s == nil {
		clientset, err := kubernetes.NewForConfig(getConfig())
		if err != nil {
			return errors.Wrapf(err, "Failed to initial Clientset %s. ", err.Error())
		}
		k8s = &K8sContext{Clientset: clientset}
	}
	return nil

}

func (web *WebContext) Handler(path string, c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Printf("Fail to upgrade with error : %s", err)
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
	err := initK8s()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
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
		case "ns -o json":
			n, err := k8s.GetNamespaces()
			if err != nil {
				buf.WriteString(err.Error())

			} else {
				s, _ := json.Marshal(n)
				buf.WriteString(string(s))
			}
			break
		case "ns":
			nl, err := k8s.GetNamespaces()
			if err != nil {
				buf.WriteString(err.Error())

			} else {
				buf.WriteString("Name\t\tStatus\t\tCreated\t\tAge\t\tLabels\n")
				for _, n := range nl {
					buf.WriteString(n.Name + "\t\t")
					buf.WriteString(n.Status + "\t\t")
					buf.WriteString(n.Created.String() + "\t\t")

					h := fmt.Sprintf("%.2f ms\t\t", float64(n.Age/1000))
					buf.WriteString(h)
					for _, label := range n.Labels {
						buf.WriteString(label.Key + "=" + label.Value + ";")
					}

					buf.WriteString("\n")
				}
			}
			break
		default:
			buf.WriteString(" < " + command + "> is not valid command.")
		}
		err = ws.WriteMessage(mt, buf.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}

	}

}
