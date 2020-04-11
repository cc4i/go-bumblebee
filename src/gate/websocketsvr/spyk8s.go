package websocketsvr

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"os"
)

var sws *websocket.Conn

func Spy(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
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
	if sws == nil {
		wsUrl := os.Getenv("SPY_SERVICE_ENDPOINT")
		ws, res, err := websocket.DefaultDialer.Dial("ws://"+wsUrl+"/spy", nil)
		if err != nil {
			log.Printf("Connect to %s with error: %s/code: %d", wsUrl, err, res.StatusCode)
		}
		sws = ws
	}
	err := sws.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("write to spy:", err)
	}
	_, message, err := sws.ReadMessage()
	if err != nil {
		log.Println("read from spy:", err)
	}
	log.Printf("recv from spy: %s", message)

	return message
}
