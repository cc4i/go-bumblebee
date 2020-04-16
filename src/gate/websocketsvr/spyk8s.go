package websocketsvr

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

var sws *websocket.Conn

func Spy(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Printf("Client (%s) is connected at (%s) by WebSocket", ws.RemoteAddr(), ws.LocalAddr())

	ws.SetCloseHandler(func(code int, text string) error {
		return ws.Close()
	})
	defer ws.Close()

	MessageHandler(ws)

}

func MessageHandler(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		////
		message = Forward2Spy(string(message))
		////

		err = ws.WriteMessage(mt, []byte(message))
		if err != nil {
			log.Println("write:", err)
			break
		}

	}
}

func Forward2Spy(command string) []byte {
	endpoint := os.Getenv("SPY_SERVICE_ENDPOINT")
	url := "ws://" + endpoint + "/spy"
	err := Reconnect(url, false)
	if err != nil {
		return []byte(err.Error())
	}
	err = sws.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("write to spy:", err)
		if e, ok := err.(*net.OpError); ok || websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			log.Println(e)
			err = Reconnect(url, true)
			if err != nil {
				return []byte(err.Error())
			}
		}
	}
	_, message, err := sws.ReadMessage()
	if err != nil {
		log.Println("read from spy:", err)
	}
	log.Printf("recv from spy: %s", message)

	return message
}

func Reconnect(url string, force bool) error {
	if sws == nil || force {
		if force {
			if sws != nil {
				sws.Close()
				sws = nil
			}
		}
		tws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Connect to %s with error: %s", url, err)
			return err
		}
		sws = tws
		sws.SetCloseHandler(func(code int, text string) error {
			log.Printf("Close connection from %s ", sws.RemoteAddr())
			sws, _, err = websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Printf("ReConnect to %s with error: %s", url, err)
			} else {
				log.Printf("ReConnect to %s with success", url)
			}
			return err
		})
		return err
	}
	return nil
}
