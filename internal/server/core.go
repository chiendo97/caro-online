package server

import (
	"time"

	"github.com/chiendo97/caro-online/internal/socket"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CoreServer interface {
	Run()

	Register(hub *Hub)
	UnRegister(hub *Hub)

	FindGame(msg msgStruct)
	JoinGame(msg msgStruct)
	CreateGame(msg msgStruct)
}

func InitCoreServer() CoreServer {

	var core = coreServer{
		hubs:      make(map[string]*Hub),
		availHubC: make(chan string, 5),

		findC:   make(chan msgStruct),
		joinC:   make(chan msgStruct),
		createC: make(chan msgStruct),

		regC:   make(chan *Hub),
		unregC: make(chan *Hub),

		done: make(chan int),
	}

	return &core
}

type coreServer struct {
	hubs      map[string]*Hub
	availHubC chan string

	findC   chan msgStruct
	joinC   chan msgStruct
	createC chan msgStruct

	regC   chan *Hub
	unregC chan *Hub

	done chan int
}

func (core *coreServer) Register(hub *Hub) {
	core.regC <- hub
}
func (core *coreServer) UnRegister(hub *Hub) {
	core.unregC <- hub
}
func (core *coreServer) FindGame(msg msgStruct) {
	core.findC <- msg
}
func (core *coreServer) JoinGame(msg msgStruct) {
	core.joinC <- msg
}
func (core *coreServer) CreateGame(msg msgStruct) {
	core.createC <- msg
}

func (core *coreServer) createHub(msg msgStruct) {
	var gameId = uuid.New().String()[:8]

	_, ok := core.hubs[gameId]

	if ok {
		log.Error("Key duplicate: ", gameId, msg.conn.RemoteAddr())
	}

	var hub = initHub(core, gameId)

	core.hubs[gameId] = hub

	core.subscribe(hub)

	go func() {
		hub.Register(socket.InitAndRunSocket(msg.conn, hub))
	}()
}

func (core *coreServer) joinHub(msg msgStruct) {

	hub, ok := core.hubs[msg.gameId]

	if !ok {
		log.Info("core: hub not found - ", msg.gameId, msg.conn.RemoteAddr())
		msg.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	go func() {
		hub.Register(socket.InitAndRunSocket(msg.conn, hub))
	}()
}

func (core *coreServer) findHub(msg msgStruct) {

	go func() {
		for {
			select {
			case gameID := <-core.availHubC:
				if _, ok := core.hubs[gameID]; !ok {
					continue
				}
				msg.gameId = gameID
				core.joinC <- msg
				return
			case <-time.After(3 * time.Second):
				core.createC <- msg
				return
			}
		}
	}()
}

func (core *coreServer) subscribe(hub *Hub) {
	go func() {
		core.availHubC <- hub.key
	}()
}

func (core *coreServer) unsubscribe(hub *Hub) {

	delete(core.hubs, hub.key)
}

func (core *coreServer) Run() {
	for {
		select {
		case <-core.done:
			return
		case hub := <-core.regC:
			log.Infof("core: hub (%s) subscribe.", hub.key)

			core.subscribe(hub)
		case hub := <-core.unregC:
			log.Infof("core: detele hub (%s)", hub.key)

			core.unsubscribe(hub)
		case msg := <-core.findC:
			log.Infof("core: socket (%s) find game", msg.conn.RemoteAddr())

			core.findHub(msg)
		case msg := <-core.joinC:
			log.Infof("core: socket (%s) join hub (%s)", msg.gameId, msg.gameId)

			core.joinHub(msg)
		case msg := <-core.createC:
			log.Infof("core: socket (%s) create hub", msg.conn.RemoteAddr())

			core.createHub(msg)
		}
	}
}
