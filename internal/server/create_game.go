package server

import (
	"fmt"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/socket"
)

func (core *coreServer) CreateGame(conn *websocket.Conn) {
	exporterCounter.WithLabelValues("CreateGame").Inc()

	core.mux.Lock()
	defer core.mux.Unlock()

	atomic.AddInt64(&core.idGenerator, 1)
	var gameId = fmt.Sprintf("%d", core.idGenerator)

	if _, found := core.hubs[gameId]; found {
		logrus.Error("Key duplicate: ", gameId, conn.RemoteAddr())
		return
	}

	var hub = initHub(core, gameId)
	core.hubs[gameId] = hub

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()

	logrus.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
}
