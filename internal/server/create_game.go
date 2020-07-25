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

	_, ok := core.hubs[gameId]
	if ok {
		logrus.Error("Key duplicate: ", gameId, conn.RemoteAddr())
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
			logrus.Errorf("hub run error: %v", err)
		}
		core.hubWG.Done()
	}()

	go func() {
		hub.OnEnter(socket.InitSocket(conn, hub))
	}()

	logrus.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
}
