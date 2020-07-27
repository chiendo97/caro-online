package client

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	mux sync.Mutex

	socket socket.Socket

	player game.Player
	game   game.Game
	bot    game.Bot

	doneC chan int
}

// InitHub init new client hub
func InitHub(c *websocket.Conn, bot game.Bot) *Hub {
	var hub = Hub{
		bot: bot,
	}

	hub.socket = socket.InitSocket(c, &hub)

	return &hub
}

func (hub *Hub) OnLeave(s socket.Socket) {

	hub.mux.Lock()
	defer hub.mux.Unlock()

	if hub.socket != nil && hub.socket == s {
		hub.socket.CloseMessage()
		hub.socket = nil
		logrus.Infof("Server disconnect")
	}
}

func (hub *Hub) OnMessage(msg socket.Message) {

	hub.mux.Lock()
	defer hub.mux.Unlock()

	switch msg.Type {
	case socket.AnnouncementMessageType:
		logrus.Infof("Server: %s\n", msg.Announcement)

	case socket.GameMessageType:
		hub.player = msg.Player
		hub.game = msg.Game

		// TODO: option to render
		// hub.game.Render()

		switch hub.game.Status {
		case game.Running:
			if hub.player == hub.game.Player {
				logrus.Infof("Your turn: \n")
				go func() {

					move, err := hub.bot.GetMove(hub.player, hub.game)
					if err != nil {
						logrus.Errorf("GetMove err: %v", err)
					}

					hub.socket.SendMessage(socket.GenerateMoveMsg(move))
				}()
			} else {
				logrus.Infof("Enemy turn.")
			}
		default:
			switch hub.game.Status {
			case game.XWin, game.OWin:
				if hub.player == hub.game.Status.GetPlayer() {
					logrus.Infof("You won !!!")
				} else {
					logrus.Infof("Your opponent won, good luck next !!")
				}
			case game.Tie:
				logrus.Infof("Game tie!!")
			}
		}

	default:
		logrus.Warn("Invalid msg:", msg)
	}
}

func (hub *Hub) Run() error {

	err1, err2 := hub.socket.Run()
	if err1 != nil || err2 != nil {
		return fmt.Errorf("%v:%v", err1, err2)
	}

	return nil
}

func (hub *Hub) Stop() {
	hub.OnLeave(hub.socket)
}
