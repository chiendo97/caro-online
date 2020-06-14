package client

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	message chan socket.Message

	player game.Player
	game   game.Game

	Socket socket.SocketI

	inputLock    bool
	inputChannel chan chan interface{}
}

// InitAndRunHub init new client hub
func InitAndRunHub() *Hub {
	var hub = Hub{
		message:      make(chan socket.Message),
		inputChannel: InpupChannel(),
	}

	go hub.Run()

	return &hub
}

func (hub *Hub) HandleMsg(msg socket.Message) {
	hub.message <- msg
}

func (hub *Hub) UnRegister(s *socket.Socket) {
	logrus.Fatal("Server disconnect")
}

func (hub *Hub) handleMsg(msg socket.Message) {

	hub.inputLock = false

	switch msg.Type {
	case socket.AnnouncementMessageType:
		logrus.Infof("Server: %s\n", msg.Announcement)

	case socket.GameMessageType:
		hub.player = msg.Player
		hub.game = msg.Game
		hub.game.Render()

		switch hub.game.GetStatus() {
		case game.Running:
			if hub.player == hub.game.Player {
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
							X:      x,
							Y:      y,
							Player: hub.player,
						})

						hub.Socket.SendMessage(msg)
					}

				}()
			} else {
				fmt.Println("Enemy turn.")
			}
		case game.XWin, game.OWin:
			if hub.player == hub.game.GetStatus().GetPlayer() {
				fmt.Println("You won !!!")
			} else {
				fmt.Println("Your opponent won, good luck next !!")
			}
		case game.Tie:
			fmt.Println("Game tie!!")
		}

	default:
		logrus.Infof("Invalid msg:", msg)
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
