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
	mux     sync.Mutex
	hubs    map[string]*Hub
	players map[*websocket.Conn]bool
}

func InitCoreServer() CoreServer {
	var core = &coreServer{
		hubs:    make(map[string]*Hub),
		players: make(map[*websocket.Conn]bool),
	}

	return core
}
