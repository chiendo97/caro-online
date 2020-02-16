package server

import (
	"log"
	"time"

	"github.com/chiendo97/caro-online/internal/socket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type msgServer struct {
	conn   *websocket.Conn
	gameId string
}

func InitMessage(conn *websocket.Conn, gameId string) msgServer {
	return msgServer{
		conn:   conn,
		gameId: gameId,
	}
}

type CoreServer struct {
	hubs          map[string]*Hub
	availableHubs chan string

	FindGame   chan msgServer
	JoinGame   chan msgServer
	CreateGame chan msgServer

	register   chan *Hub
	unregister chan *Hub

	done chan int
}

func InitCore() *CoreServer {
	return &CoreServer{
		hubs:          make(map[string]*Hub),
		availableHubs: make(chan string, 5),

		FindGame:   make(chan msgServer),
		JoinGame:   make(chan msgServer),
		CreateGame: make(chan msgServer),

		register:   make(chan *Hub),
		unregister: make(chan *Hub),

		done: make(chan int),
	}
}

func (core *CoreServer) createHub(msg msgServer) string {
	var gameId = uuid.New().String()[:8]

	_, ok := core.hubs[gameId]

	if ok {
		log.Panicln("Key duplicate: ", gameId, msg.conn.RemoteAddr())
	}

	var hub = InitHub(core, gameId)
	go hub.run()

	core.hubs[gameId] = hub

	core.subscribe(hub)

	return hub.key
}

func (core *CoreServer) joinHub(msg msgServer) {

	hub, ok := core.hubs[msg.gameId]

	if !ok {
		log.Println("core: hub not available - ", msg.gameId, msg.conn.RemoteAddr())
		msg.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	var s = socket.InitSocket(msg.conn, hub)
	hub.register <- &s

	go s.Read()
	go s.Write()

}
func (core *CoreServer) findHub(msg msgServer) {
	var tick = time.After(3 * time.Second)

	for {
		select {
		case gameId := <-core.availableHubs:
			msg.gameId = gameId
			_, ok := core.hubs[msg.gameId]
			if ok {
				core.joinHub(msg)
				return
			}
		case <-tick:
			var gameId = core.createHub(msg)
			msg.gameId = gameId
			core.joinHub(msg)
			return
		}
	}

}

func (core *CoreServer) subscribe(hub *Hub) {
	core.availableHubs <- hub.key
}

func (core *CoreServer) unsubscribe(hub *Hub) {

	delete(core.hubs, hub.key)
}

func (core *CoreServer) Run() {
	for {
		select {
		case <-core.done:
			return
		case hub := <-core.register:
			log.Printf("core: hub (%s) subscribe.", hub.key)

			core.subscribe(hub)
		case hub := <-core.unregister:
			log.Printf("core: detele hub (%s)", hub.key)

			core.unsubscribe(hub)
		case msg := <-core.FindGame:
			log.Printf("core: socket (%s) find game", msg.conn.RemoteAddr())

			core.findHub(msg)
		case msg := <-core.JoinGame:
			log.Printf("core: socket (%s) join hub (%s)", msg.gameId, msg.gameId)

			core.joinHub(msg)
		case msg := <-core.CreateGame:
			log.Printf("core: socket (%s) create hub", msg.conn.RemoteAddr())

			core.createHub(msg)
		}
	}
}
