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
}

// NewHub init new client hub
func NewHub(c *websocket.Conn, bot game.Bot) *Hub {
	var hub = Hub{
		bot: bot,
	}

	hub.socket = socket.NewSocket(c, &hub)

	return &hub
}

func (hub *Hub) OnLeave(s socket.Socket) {
	hub.socket.Stop()
}

func (hub *Hub) OnMessage(msg socket.Message) {

	hub.mux.Lock()
	defer hub.mux.Unlock()

	switch msg.Type {
	case socket.Announce:
		logrus.Debugf("Server: %s", msg.Announce)

	case socket.Game:
		hub.player = msg.Player
		hub.game = msg.Game

		// TODO: option to render
		// hub.game.Render()

		switch hub.game.Status {
		case game.Running:
			if hub.player == hub.game.Player {
				logrus.Debugf("Your turn: \n")
				move, err := hub.bot.GetMove(hub.player, hub.game)
				if err != nil {
					logrus.Errorf("GetMove err: %v", err)
				}

				msg := socket.NewMoveMsg(move)
				hub.socket.SendMessage(msg)
			} else {
				logrus.Debugf("Enemy turn.")
			}
		default:
			switch hub.game.Status {
			case game.XWin, game.OWin:
				if hub.player == hub.game.Status.GetPlayer() {
					logrus.Debugf("You won !!!")
				} else {
					logrus.Debugf("Your opponent won, good luck next !!")
				}
			case game.Tie:
				logrus.Debugf("Game tie!!")
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
