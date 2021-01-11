package server

import (
	"time"

	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (core *coreServer) FindGame(conn *websocket.Conn) {
	core.mux.Lock()
	defer core.mux.Unlock()

	logrus.Infof("core: socket (%s) find game", conn.RemoteAddr())

	core.players[conn] = true

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if found := core.findPlayer(conn); !found {
				continue
			}
			break
		}
	}()
}

func (core *coreServer) findPlayer(conn *websocket.Conn) bool {

	// Check connection
	err := conn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		return true
	}

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, found := core.players[conn]; !found {
		return true
	}

	for player := range core.players {
		if conn == player {
			continue
		}

		logrus.Debugf("Match: %v=%v", conn.RemoteAddr(), player.RemoteAddr())

		delete(core.players, conn)
		delete(core.players, player)

		var gameId = uuid.New().String()[:8]
		var hub = initHub(core, gameId)

		core.hubs[gameId] = hub

		go func() {
			hub.OnEnter(socket.InitSocket(conn, hub))
			hub.OnEnter(socket.InitSocket(player, hub))
		}()

		logrus.Infof("core: socket (%s) create hub (%s)", conn.RemoteAddr(), gameId)
		return true
	}

	return false
}
