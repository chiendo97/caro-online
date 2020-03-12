package client

import (
	"fmt"
	"log"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	s "github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	message chan s.Message

	game game.Game

	Socket *s.Socket
}

// InitHub init new client hub
func InitHub() *Hub {
	var hub = Hub{
		message: make(chan s.Message),
	}

	go hub.Run()

	return &hub
}

func (hub *Hub) ReceiveMsg(msg s.Message) {
	hub.message <- msg
}

func (hub *Hub) Unregister(s *s.Socket) {
	log.Fatalln("Server down")
}

func (hub *Hub) handleMsg(msg s.Message) {

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
				case hub.Socket.Message <- msg:
				case <-time.After(3 * time.Second):
					log.Panicln("Cant send move to socket")
				}

			default:
				fmt.Println("Enemy turn.")
			}
		}

	default:
		log.Panicln("Invalid msg:", msg)
	}
}

func (hub *Hub) Run() {

	for {
		select {
		case msg := <-hub.message:
			hub.handleMsg(msg)
		}
	}
}
