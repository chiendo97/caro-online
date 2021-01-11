package server

import (
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (core *coreServer) CreateGame(conn *websocket.Conn) {
	core.mux.Lock()
	defer core.mux.Unlock()

	var gameId = uuid.New().String()[:8]

	if _, found := core.hubs[gameId]; found {
		logrus.Error("Key duplicate: ", gameId, conn.RemoteAddr())
	}

	var hub = initHub(core, gameId)
	core.hubs[gameId] = hub

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()

	logrus.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
}
