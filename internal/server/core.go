package server

import (
	"log"
	"time"

	"github.com/chiendo97/caro-online/internal/socket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CoreServer struct {
	hubs          map[string]*hubG
	availableHubs chan string

	FindGame   chan msgStruct
	JoinGame   chan msgStruct
	CreateGame chan msgStruct

	register   chan *hubG
	unregister chan *hubG

	done chan int
}

func InitAndRunCore() *CoreServer {

	var core = CoreServer{
		hubs:          make(map[string]*hubG),
		availableHubs: make(chan string, 5),

		FindGame:   make(chan msgStruct),
		JoinGame:   make(chan msgStruct),
		CreateGame: make(chan msgStruct),

		register:   make(chan *hubG),
		unregister: make(chan *hubG),

		done: make(chan int),
	}

	go core.run()

	return &core
}

func (core *CoreServer) createHub(msg msgStruct) {
	var gameId = uuid.New().String()[:8]

	_, ok := core.hubs[gameId]

	if ok {
		log.Panicln("Key duplicate: ", gameId, msg.conn.RemoteAddr())
	}

	var hub = initHub(core, gameId)

	core.hubs[gameId] = hub

	core.subscribe(hub)

	go func() {
		hub.register <- socket.InitAndRunSocket(msg.conn, hub)
	}()
}

func (core *CoreServer) joinHub(msg msgStruct) {

	hub, ok := core.hubs[msg.gameId]

	if !ok {
		log.Println("core: hub not found - ", msg.gameId, msg.conn.RemoteAddr())
		msg.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	go func() {
		hub.register <- socket.InitAndRunSocket(msg.conn, hub)
	}()
}

func (core *CoreServer) findHub(msg msgStruct) {

	go func() {
		for {
			select {
			case gameID := <-core.availableHubs:
				if _, ok := core.hubs[gameID]; !ok {
					continue
				}
				msg.gameId = gameID
				core.JoinGame <- msg
				return
			case <-time.After(3 * time.Second):
				core.CreateGame <- msg
				return
			}
		}
	}()
}

func (core *CoreServer) subscribe(hub *hubG) {
	go func() {
		core.availableHubs <- hub.key
	}()
}

func (core *CoreServer) unsubscribe(hub *hubG) {

	delete(core.hubs, hub.key)
}

func (core *CoreServer) run() {
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
