package server

import (
	"sync"
	"time"

	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type CoreServer interface {
	Run() error
	Stop()

	FindGame(conn *websocket.Conn)
	JoinGame(conn *websocket.Conn, gameId string)
	CreateGame(conn *websocket.Conn)
}

type coreServer struct {
	mux   sync.Mutex
	hubWG sync.WaitGroup

	hubs     map[string]*Hub
	availHub chan string

	done chan int
}

func InitCoreServer() CoreServer {

	var core = coreServer{

		hubs:     make(map[string]*Hub),
		availHub: make(chan string, 5),

		done: make(chan int),
	}

	return &core
}

func (core *coreServer) leaveHub(hub *Hub) {

	log.Infof("core: delete hub (%s)", hub.key)

	close(hub.doneC)
	delete(core.hubs, hub.key)
}

func (core *coreServer) leaveAllHubs() {

	core.mux.Lock()
	defer core.mux.Unlock()

	for _, hub := range core.hubs {
		core.leaveHub(hub)
	}
}

func (core *coreServer) Stop() {
	close(core.done)
}

func (core *coreServer) Register(hub *Hub) {

	log.Infof("core: hub (%s) subscribe.", hub.key)

	core.availHub <- hub.key
}

func (core *coreServer) UnRegister(hub *Hub) {

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, ok := core.hubs[hub.key]; ok {
		core.leaveHub(hub)
	}
}

func (core *coreServer) FindGame(conn *websocket.Conn) {

	log.Infof("core: socket (%s) find game", conn.RemoteAddr())

	go func() {
		for {
			select {
			case gameID := <-core.availHub:
				core.mux.Lock()
				if _, ok := core.hubs[gameID]; !ok {
					core.mux.Unlock()
					continue
				}
				core.mux.Unlock()
				go core.JoinGame(conn, gameID)
				return
			case <-time.After(1 * time.Second):
				go core.CreateGame(conn)
				return
			}
		}
	}()
}

func (core *coreServer) JoinGame(conn *websocket.Conn, gameId string) {

	core.mux.Lock()
	defer core.mux.Unlock()

	log.Infof("core: socket (%s) join hub (%s)", conn.RemoteAddr(), gameId)

	hub, ok := core.hubs[gameId]

	if !ok {
		log.Warn("core: hub not found - ", gameId, conn.RemoteAddr())
		conn.WriteMessage(websocket.CloseMessage, []byte{})

		return
	}

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()
}

func (core *coreServer) CreateGame(conn *websocket.Conn) {

	core.mux.Lock()
	defer core.mux.Unlock()

	var gameId = uuid.New().String()[:8]

	_, ok := core.hubs[gameId]
	if ok {
		log.Error("Key duplicate: ", gameId, conn.RemoteAddr())
	}

	var hub = initHub(core, gameId)

	core.hubs[gameId] = hub

	go func() {
		core.availHub <- hub.key
	}()

	core.hubWG.Add(1)
	go func() {
		err := hub.Run()
		if err != nil {
			log.Errorf("hub run error: %v", err)
		}
		core.hubWG.Done()
	}()

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()

	log.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
}

func (core *coreServer) Run() error {

	log.Infof("Core start")
	defer log.Infof("Core stop")

	for {
		select {
		case <-core.done:
			core.leaveAllHubs()

			core.hubWG.Wait()
			return nil
		}
	}

}
