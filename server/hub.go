package main

import (
	"log"
	"time"

	"github.com/chiendo97/caro-online/game"
	s "github.com/chiendo97/caro-online/socket"
)

type Hub struct {
	core *coreServer
	key  string

	game    game.Game
	players map[*s.Socket]int

	message    chan s.Message
	register   chan *s.Socket
	unregister chan *s.Socket

	done chan int
}

func (hub *Hub) ReceiveMsg(msg s.Message) {
	hub.message <- msg
}

func (hub *Hub) Unregister(s *s.Socket) {
	hub.unregister <- s
}

func InitHub(core *coreServer, key string) *Hub {

	return &Hub{
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
	log.Printf("hub: send (%s) to (%s)", msg, socket.Conn.RemoteAddr())
}

func (hub *Hub) broadcast() {

	if len(hub.players) < 2 {
		hub.transmitMsgToAll(s.GenerateErrMsg("hub (" + hub.key + ") wait for players"))
		return
	}

	for socket := range hub.players {
		var game = hub.game
		game.WhoAmI = hub.players[socket]
		sendMessage(socket, s.GenerateGameMsg(game))
	}
}

func (hub *Hub) transmitMsgToAll(msg s.Message) {
	for socket := range hub.players {
		sendMessage(socket, msg)
	}
}

func (hub *Hub) handleMsg(msg s.Message) {

	if msg.Kind != s.MoveMessage {
		log.Panicln("hub: No msg kind case", msg)
		return
	}

	game, err := hub.game.TakeMove(msg.Move)

	if err != nil {
		log.Println("hub: game error - ", err)
	} else {
		hub.game = game
	}

	hub.broadcast()

	if hub.game.Status != -1 {
		var tick = time.After(5 * time.Second)
		<-tick
		close(hub.done)
		hub.core.unregister <- hub
	}
}

func (hub *Hub) subscribe(socket *s.Socket) {
	if len(hub.players) == 2 {
		log.Println("hub: room is full", socket.Conn.RemoteAddr().String())
		close(socket.Message)
		return
	}

	var id = 0
	for _, otherId := range hub.players {
		id = 1 - otherId
	}

	log.Printf("hub: take new socket (%s) as (%d)", socket.Conn.RemoteAddr(), id)
	hub.players[socket] = id

	hub.broadcast()
}

func (hub *Hub) unsubscribe(socket *s.Socket) {
	if _, ok := hub.players[socket]; ok {
		log.Println("hub: Player left:", socket.Conn.RemoteAddr())

		delete(hub.players, socket)
		close(socket.Message)

		hub.broadcast()

		if len(hub.players) == 1 {
			hub.core.register <- hub
		} else {
			hub.core.unregister <- hub
		}
	}
}

func (hub *Hub) run() {

	// defer func() {
	// 	close(hub.message)
	// 	close(hub.register)
	// 	close(hub.unregister)
	// }()

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
