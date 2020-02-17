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
	hubs          map[string]*hubG
	availableHubs chan string

	FindGame   chan msgServer
	JoinGame   chan msgServer
	CreateGame chan msgServer

	register   chan *hubG
	unregister chan *hubG

	done chan int
}

func InitCore() *CoreServer {
	return &CoreServer{
		hubs:          make(map[string]*hubG),
		availableHubs: make(chan string, 5),

		FindGame:   make(chan msgServer),
		JoinGame:   make(chan msgServer),
		CreateGame: make(chan msgServer),

		register:   make(chan *hubG),
		unregister: make(chan *hubG),

		done: make(chan int),
	}
}

func (core *CoreServer) createHub(msg msgServer) {
	var gameId = uuid.New().String()[:8]

	_, ok := core.hubs[gameId]

	if ok {
		log.Panicln("Key duplicate: ", gameId, msg.conn.RemoteAddr())
	}

	var hub = initHub(core, gameId)
	go hub.run()

	core.hubs[gameId] = hub

	core.subscribe(hub)

	hub.register <- socket.InitSocket(msg.conn, hub)
}

func (core *CoreServer) joinHub(msg msgServer) {

	hub, ok := core.hubs[msg.gameId]

	if !ok {
		log.Println("core: hub not found - ", msg.gameId, msg.conn.RemoteAddr())
		msg.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	hub.register <- socket.InitSocket(msg.conn, hub)
}
func (core *CoreServer) findHub(msg msgServer) {

	go func() {
		for {
			select {
			case gameID := <-core.availableHubs:
				if _, ok := core.hubs[gameID]; ok {
					msg.gameId = gameID
					select {
					case core.JoinGame <- msg:
						return
					case <-time.After(3 * time.Second):
						log.Panicf("core: Can't join game %s", msg.conn.RemoteAddr())
					}
				}
			case <-time.After(3 * time.Second):
				select {
				case core.CreateGame <- msg:
					return
				case <-time.After(3 * time.Second):
					log.Panicf("core: Can't create game %s", msg.conn.RemoteAddr())
				}
			}
		}
	}()
}

func (core *CoreServer) subscribe(hub *hubG) {
	go func() {
		select {
		case core.availableHubs <- hub.key:
		case <-time.After(3 * time.Second):
			log.Panicf("core: Can't subscribe hub %s", hub.key)
		}
	}()
}

func (core *CoreServer) unsubscribe(hub *hubG) {

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
