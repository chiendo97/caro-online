package server

import (
	"time"

	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (core *coreServer) FindGame(conn *websocket.Conn) {

	logrus.Infof("core: socket (%s) find game", conn.RemoteAddr())

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

func (core *coreServer) findHub(conn *websocket.Conn, gameID string) bool {

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, ok := core.hubs[gameID]; !ok {
		return false
	}

	go core.JoinGame(conn, gameID)
	return true
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

		logrus.Warnf("DEBUG: %v=%v", conn.RemoteAddr(), player.RemoteAddr())

		delete(core.players, player)
		delete(core.players, conn)

		var gameId = uuid.New().String()[:8]
		var hub = initHub(core, gameId)

		core.hubs[gameId] = hub

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
			hub.OnEnter(socket.InitSocket(player, hub))
		}()

		logrus.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
		return true
	}

	return false
}
