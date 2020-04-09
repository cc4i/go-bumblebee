package websocketsvr

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func Realtime(c *gin.Context) {
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

		//TODO
		//TODO 1. realtime what? who's backend? what components?
		//
		for {
			err = ws.WriteMessage(mt, []byte("realtime message ..."+time.Now().String()))
			if err != nil {
				log.Println("write:", err)
				break
			}
			time.Sleep(2 * time.Second)
		}

	}
}
