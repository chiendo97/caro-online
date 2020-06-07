package client

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	message chan socket.Message

	game game.Game

	Socket *socket.Socket

	inputLock    bool
	inputChannel chan chan interface{}
}

// InitAndRunHub init new client hub
func InitAndRunHub() *Hub {
	var hub = Hub{
		message:      make(chan socket.Message),
		inputChannel: InpupChannel(),
	}
}

func (hub *Hub) ReceiveMsg(msg socket.Message) {
	hub.message <- msg
}

func (hub *Hub) Unregister(s *socket.Socket) {
	log.Fatalln("Server disconnect")
}

func (hub *Hub) handleMsg(msg socket.Message) {

	hub.inputLock = false

	switch msg.Kind {
	case socket.MsgMessage:
		log.Printf("\nServer: %s\n", msg.Msg)

	case socket.GameMessage:
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
				hub.inputLock = true
				fmt.Printf("Your turn: ")
				go func() {
					var x, y int
					input := make(chan interface{})
					hub.inputChannel <- input
					xs := <-input
					hub.inputChannel <- input
					ys := <-input
					x, _ = strconv.Atoi(xs.(string))
					y, _ = strconv.Atoi(ys.(string))

					if hub.inputLock == true {
						var msg = socket.GenerateMoveMsg(game.Move{
							X:    x,
							Y:    y,
							Turn: hub.game.WhoAmI,
						})

						hub.Socket.Message <- msg
					}

				}()

			default:
				log.Panicln("Invalid request kind:", msg)
			}
		}
	}
}
