package tcpsvr

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"net"
)

func RunServer(endpoint string) error {

	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf(" Listening and serving TCP on %s\n", endpoint)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		go ConnectionHandler(c)

	}

}

// ConnectionHandler is to handle each individual connection.
func ConnectionHandler(con net.Conn) {
	for {
		str, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		} else {
			log.Println("recv: ", str)
		}

		switch str {
			case "help\r\n":
				con.Write([]byte("An echo service with delimit newline. \r\n'help' - simple help.\r\n 'exit' - close connection.\r\n"))
				break
			case "exit\r\n":
				con.Close()
				break
			default:
				_, err = con.Write([]byte(str))
				if err != nil {
					log.Println(err)
				}
		}


	}

}
