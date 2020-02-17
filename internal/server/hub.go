package server

import (
	"fmt"
	"log"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	s "github.com/chiendo97/caro-online/internal/socket"
)

type hubG struct {
	core    *CoreServer
	key     string
	game    game.Game
	players map[*s.Socket]int

	message    chan s.Message
	register   chan *s.Socket
	unregister chan *s.Socket
	done       chan int
}

func (hub *hubG) ReceiveMsg(msg s.Message) {
	hub.message <- msg
}

func (hub *hubG) Unregister(s *s.Socket) {
	hub.unregister <- s
}

func initHub(core *CoreServer, key string) *hubG {

	return &hubG{
		core:    core,
		key:     key,
		message: make(chan s.Message),

		game:    game.InitGame(key),
		players: make(map[*s.Socket]int),

		register:   make(chan *s.Socket),
		unregister: make(chan *s.Socket),

		done: make(chan int),
	}
}

func sendMessage(socket *s.Socket, msg s.Message) {
	socket.Message <- msg
	// go func() {
	// 	select {
	// 	case socket.Message <- msg:
	// 		log.Printf("hub: send (%s) to (%s)", msg, socket.Conn.RemoteAddr())
	// 	}
	// }()
}

func (hub *hubG) broadcast() {

	if len(hub.players) < 2 {
		var msg = s.GenerateErrMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for socket := range hub.players {
			sendMessage(socket, msg)
		}
	} else {
		for socket := range hub.players {
			var game = hub.game
			game.WhoAmI = hub.players[socket]
			sendMessage(socket, s.GenerateGameMsg(game))
		}
	}

}

func (hub *hubG) handleMsg(msg s.Message) {

	if msg.Kind != s.MoveMessage {
		log.Panicf("hub %s: No msg kind case %s", hub.key, msg)
		return
	}

	game, err := hub.game.TakeMove(msg.Move)

	if err != nil {
		log.Printf("hub %s: game error - %s", hub.key, err)
	} else {
		hub.game = game
	}

	hub.broadcast()

	if hub.game.Status != -1 {
		go func() {
			select {
			case <-time.After(5 * time.Second):
				hub.done <- 1
				hub.core.unregister <- hub
			}
		}()
	}
}

func (hub *hubG) subscribe(socket *s.Socket) {
	if len(hub.players) == 2 {
		log.Printf("hub %s: room is full %s", hub.key, socket.Conn.RemoteAddr().String())
		close(socket.Message)
		return
	}

	var id = 0
	for _, otherId := range hub.players {
		id = 1 - otherId
	}

	log.Printf("hub %s: take new socket %s as player %d", hub.key, socket.Conn.RemoteAddr(), id)
	hub.players[socket] = id

	hub.broadcast()
}

func (hub *hubG) unsubscribe(socket *s.Socket) {
	if _, ok := hub.players[socket]; ok {
		log.Printf("hub %s: Player %s left", hub.key, socket.Conn.RemoteAddr())

		delete(hub.players, socket)
		close(socket.Message)

		hub.broadcast()

		go func() {
			if len(hub.players) == 1 {
				select {
				case hub.core.register <- hub:
				}
			} else {
				select {
				case hub.core.unregister <- hub:
				}
			}
		}()

	}
}

func (hub *hubG) run() {

	for {
		select {
		case msg := <-hub.message:
			hub.handleMsg(msg)
		case socket := <-hub.register:
			hub.subscribe(socket)
		case socket := <-hub.unregister:
			hub.unsubscribe(socket)
		case <-hub.done:
			for socket := range hub.players {
				close(socket.Message)
			}
			return
		}
	}
}
