package client

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d\t", filename, f.Line)
		},
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.ErrorLevel)

	log.SetReportCaller(true)
}

type Hub struct {
	message chan socket.Message

	player game.Player
	game   game.Game
	bot    Bot

	socket socket.Socket

	inputLock    bool
	inputChannel chan chan interface{}
}

// InitHub init new client hub
func InitHub(c *websocket.Conn, bot Bot) *Hub {
	var hub = Hub{
		bot:          bot,
		message:      make(chan socket.Message),
		inputChannel: InpupChannel(),
	}

	hub.socket = socket.InitSocket(c, &hub)

	return &hub
}

func (hub *Hub) HandleMsg(msg socket.Message) {
	hub.message <- msg
}

func (hub *Hub) UnRegister(s socket.Socket) {
	log.Info("Server disconnect")
}

func (hub *Hub) handleMsg(msg socket.Message) {

	hub.inputLock = false

	switch msg.Type {
	case socket.AnnouncementMessageType:
		log.Infof("Server: %s\n", msg.Announcement)

	case socket.GameMessageType:
		hub.player = msg.Player
		hub.game = msg.Game

		// hub.game.Render()

		switch hub.game.Status {
		case game.Running:
			if hub.player == hub.game.Player {
				hub.inputLock = true
				//fmt.Printf("Your turn: ")
				go func() {
					var x, y int
					move, _ := hub.bot.GetMove(hub.player, hub.game)
					x = move.X
					y = move.Y
					// input := make(chan interface{})
					// hub.inputChannel <- input
					// xs := <-input
					// hub.inputChannel <- input
					// ys := <-input
					// x, _ = strconv.Atoi(xs.(string))
					// y, _ = strconv.Atoi(ys.(string))

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
				//fmt.Println("Enemy turn.")
			}
		case game.XWin, game.OWin:
			if hub.player == hub.game.Status.GetPlayer() {
				//fmt.Println("You won !!!")
			} else {
				//fmt.Println("Your opponent won, good luck next !!")
			}
		case game.Tie:
			//fmt.Println("Game tie!!")
		}

	default:
		log.Warn("Invalid msg:", msg)
	}
}

func (hub *Hub) Run() error {

	errC := make(chan error)

	go func() {
		err1, err2 := hub.socket.Run()
		if err1 != nil || err2 != nil {
			errC <- fmt.Errorf("%v:%v", err1, err2)
		} else {
			errC <- nil
		}
	}()

	for {
		select {
		case msg := <-hub.message:
			hub.handleMsg(msg)
		case err := <-errC:
			log.Info("Hub shutdown")
			return err
		}
	}
}

func (hub *Hub) Stop() {
	hub.socket.CloseMessage()
}
