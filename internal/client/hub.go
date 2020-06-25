package client

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	message chan socket.Message

	player game.Player
	game   game.Game

	socket socket.Socket

	inputLock    bool
	inputChannel chan chan interface{}

	done chan int
}

// InitHub init new client hub
func InitHub(c *websocket.Conn) *Hub {
	var hub = Hub{
		message:      make(chan socket.Message),
		inputChannel: InpupChannel(),
		done:         make(chan int),
	}

	hub.socket = socket.InitAndRunSocket(c, &hub)

	return &hub
}

func (hub *Hub) HandleMsg(msg socket.Message) {
	hub.message <- msg
}

func (hub *Hub) UnRegister(s socket.Socket) {
	logrus.Info("Server disconnect")
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

		switch hub.game.Status {
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

						hub.socket.SendMessage(msg)
					}

				}()
			} else {
				fmt.Println("Enemy turn.")
			}
		case game.XWin, game.OWin:
			if hub.player == hub.game.Status.GetPlayer() {
				fmt.Println("You won !!!")
			} else {
				fmt.Println("Your opponent won, good luck next !!")
			}
		case game.Tie:
			fmt.Println("Game tie!!")
		}

	default:
		logrus.Warn("Invalid msg:", msg)
	}
}

func (hub *Hub) Run() {

	for {
		select {
		case msg := <-hub.message:
			hub.handleMsg(msg)
		case <-hub.done:
			return
		}
	}
}

func (hub *Hub) Stop() {
	hub.socket.CloseMessage()
	close(hub.done)
}
