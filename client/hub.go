package main

import (
	"fmt"
	"helloworld/caro/game"
	s "helloworld/caro/socket"
	"log"
)

type Hub struct {
	message chan s.Message

	game game.Game

	socket *s.Socket
}

func InitHub() Hub {
	return Hub{
		message: make(chan s.Message),
		game:    game.InitGame(),
	}
}

func (hub *Hub) ReceiveMsg(msg s.Message) {
	hub.message <- msg
}

func (hub *Hub) Unregister(s *s.Socket) {
	log.Fatalln("Server down")
}

func (hub *Hub) run() {

	for {
		select {
		case msg := <-hub.message:

			switch msg.Kind {
			case "error":
				log.Println("Msg msg: ", msg.Msg)

			case "game":
				hub.game = msg.Game
				hub.game.Render()

				switch hub.game.WhoAmI {
				case hub.game.XFirst:
					// input
					var x, y int
					fmt.Print("Your turn: ")
					fmt.Scanln(&x, &y)

					var msg = s.GenerateMoveMsg(game.Move{
						X:    x,
						Y:    y,
						Turn: hub.game.WhoAmI,
					})

					select {
					case hub.socket.Message <- msg:
					default:
						log.Panicln("Cant send move to socket")
					}

				default:
					fmt.Println("Enemy turn.")
				}
			default:
				log.Panicln("Invalid msg kind", msg)
			}
		}
	}
}
