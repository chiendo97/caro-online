package main

import (
	"fmt"
	"github.com/chiendo97/caro-online/game"
	s "github.com/chiendo97/caro-online/socket"
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
		// game:    game.InitGame(),
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
			case s.MsgMessage:
				log.Println("Server:", msg.Msg)

			case s.GameMessage:
				hub.game = msg.Game
				hub.game.Render()

				if hub.game.Status == 0 || hub.game.Status == 1 {

					if hub.game.WhoAmI == hub.game.Status {
						fmt.Println("You won !!!")
					} else {
						fmt.Println("Your opponent won, good luck next !!")
					}
				} else if hub.game.Status == 2 {
					fmt.Println("Game tie!!")
				} else {
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
				}

			default:
				log.Panicln("Invalid request kind:", msg)
			}
		}
	}
}
