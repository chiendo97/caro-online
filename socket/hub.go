package socket

import (
	"fmt"
	"log"
)

type Hub struct {
	Message chan Message
}

func (hub Hub) SendMsg(sockets []*Socket) {
	for {
		select {
		case msg, ok := <-hub.Message:
			if !ok {
				log.Panicln("Message chan close")
				return
			}
			log.Println(msg)
			for _, socket := range sockets {
				select {
				case socket.Message <- msg:
					fmt.Println("sending", msg)
				default:
					log.Panicln("cant send msg", msg)
				}

			}

		default:
			log.Panicln("Cant send msg")
		}
	}
}

func (hub Hub) ReceiveMsg(msg Message) {
	switch msg.Kind {
	default:
		log.Panicln("No case for msg", msg)
	}
}
