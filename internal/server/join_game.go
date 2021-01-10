package server

import (
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (core *coreServer) JoinGame(conn *websocket.Conn, gameId string) {

	core.mux.Lock()
	defer core.mux.Unlock()

	logrus.Infof("core: socket (%s) join hub (%s)", conn.RemoteAddr(), gameId)

	hub, ok := core.hubs[gameId]

	if !ok {
		logrus.Warn("core: hub not found - ", gameId, conn.RemoteAddr())
		// conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()
}
