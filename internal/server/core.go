package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type CoreServer interface {
	FindGame(conn *websocket.Conn)
	JoinGame(conn *websocket.Conn, gameId string)
	CreateGame(conn *websocket.Conn)
}

type coreServer struct {
	idGenerator int64
	mux         sync.Mutex
	hubs        map[string]*Hub
	playerC     chan *websocket.Conn
}

func InitCoreServer() CoreServer {
	var core = &coreServer{
		hubs:    make(map[string]*Hub),
		playerC: make(chan *websocket.Conn),
	}

	go core.findGame()

	return core
}
