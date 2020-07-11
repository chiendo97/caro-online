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

	players map[*websocket.Conn]bool

	done chan int
}

func InitCoreServer() CoreServer {

	var core = coreServer{

		hubs:     make(map[string]*Hub),
		availHub: make(chan string, 5),

		players: make(map[*websocket.Conn]bool),

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

func (core *coreServer) findPlayer(conn *websocket.Conn) bool {

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, ok := core.players[conn]; !ok {
		return true
	}

	err := conn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		return true
	}

	// log.Warnf("%v=%v", conn.RemoteAddr(), len(core.players))
	for player := range core.players {
		if conn == player {
			continue
		}

		log.Warnf("DEBUG: %v=%v", conn.RemoteAddr(), player.RemoteAddr())

		delete(core.players, player)
		delete(core.players, conn)

		var gameId = uuid.New().String()[:8]
		var hub = initHub(core, gameId)

		core.hubs[gameId] = hub

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
			hub.OnEnter(socket.InitSocket(player, hub))
		}()

		log.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
		return true
	}

	return false
}

func (core *coreServer) findHub(conn *websocket.Conn, gameID string) bool {

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, ok := core.hubs[gameID]; !ok {
		return false
	}

	go core.JoinGame(conn, gameID)
	return true
}

func (core *coreServer) FindGame(conn *websocket.Conn) {

	log.Infof("core: socket (%s) find game", conn.RemoteAddr())

	core.mux.Lock()
	defer core.mux.Unlock()

	core.players[conn] = true

	go func() {

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if ok := core.findPlayer(conn); ok {
					return
				}
			case gameID := <-core.availHub:
				if ok := core.findHub(conn, gameID); ok {
					return
				}
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
		core.Register(hub)
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
