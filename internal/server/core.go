package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type CoreServer interface {
	Run() error
	Stop()

	FindGame(conn *websocket.Conn)
	JoinGame(conn *websocket.Conn, gameId string)
	CreateGame(conn *websocket.Conn)
}

type coreServer struct {
	mux sync.Mutex
	// hubWG sync.WaitGroup

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

func (core *coreServer) Run() error {

	logrus.Infof("Core start")
	defer logrus.Infof("Core stop")

	// for {
	//     select {
	//     case <-core.done:
	//         core.leaveAllHubs()

	// core.hubWG.Wait()
	//         return nil
	//     }
	// }

	return nil
}

func (core *coreServer) Stop() {
	// close(core.done)
}
