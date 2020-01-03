package main

import (
	"helloworld/caro/game"
	s "helloworld/caro/socket"
	"log"
)

type Hub struct {
	key string

	game    game.Game
	players map[*s.Socket]int

	message    chan s.Message
	register   chan *s.Socket
	unregister chan *s.Socket

	done chan int
}

func InitHub(key string) *Hub {

	return &Hub{
		key:     key,
		message: make(chan s.Message),

		game:    game.InitGame(),
		players: make(map[*s.Socket]int),

		register:   make(chan *s.Socket),
		unregister: make(chan *s.Socket),

		done: make(chan int),
	}
}

func sendMessage(socket *s.Socket, msg s.Message) {
	select {
	case socket.Message <- msg:
		log.Println("Send msg to socket", socket.Conn.RemoteAddr().String(), msg)
	}
}

func (hub *Hub) broadcastGame() {

	for socket := range hub.players {
		var game = hub.game
		game.WhoAmI = hub.players[socket]
		var msg = s.GenerateGameMsg(game)
		sendMessage(socket, msg)
	}
}

func (hub *Hub) broadcast(msg s.Message) {

	for socket := range hub.players {
		sendMessage(socket, msg)
	}
}

func (hub *Hub) ReceiveMsg(msg s.Message) {
	select {
	case hub.message <- msg:
	}
}

func (hub *Hub) Unregister(s *s.Socket) {
	select {
	case hub.unregister <- s:
	}
}

func (hub *Hub) run() {

	for {
		select {
		case msg := <-hub.message:
			switch msg.Kind {
			case s.MoveMessage:
				var move = msg.Move
				game, err := hub.game.TakeMove(move)

				if err != nil {
					log.Println("Error with game:", err)
					hub.broadcastGame()
					continue
				}

				hub.game = game
				hub.broadcastGame()

			case s.MsgMessage:
				hub.broadcast(msg)

			default:
				log.Panicln("No msg kind case", msg)
				return
			}
		case socket := <-hub.register:
			log.Println("Hub receive new socket:", socket.Conn.RemoteAddr())

			var id int

			switch len(hub.players) {
			case 0:
				id = 0
			case 1:
				for _, otherID := range hub.players {
					id = 1 - otherID
				}
			default:
				log.Println("Room is full", socket.Conn.RemoteAddr().String())
				close(socket.Message)
				continue
			}

			hub.players[socket] = id
			log.Println("Hub take new socket:", id)

			switch len(hub.players) {
			case 2:
				hub.broadcastGame()
			default:
				hub.broadcast(s.GenerateErrMsg("Waiting for other players"))
			}

		case socket := <-hub.unregister:
			if _, ok := hub.players[socket]; ok {
				log.Println("Player left:", socket.Conn.RemoteAddr())
				delete(hub.players, socket)
				close(socket.Message)
				hub.broadcast(s.GenerateErrMsg("Waiting for other players"))
			}
		case <-hub.done:
			log.Panicln("Hub immediately stop")
			for socket := range hub.players {
				close(socket.Message)
			}
			return
		}
	}
}
