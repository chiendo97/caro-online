package server

import (
	"fmt"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (core *coreServer) FindGame(conn *websocket.Conn) {
	exporterCounter.WithLabelValues("FindGame").Inc()
	core.playerC <- conn
}

func (core *coreServer) findGame() {
	playerQueue := make([]*websocket.Conn, 0, 2)
	for player := range core.playerC {
		playerQueue = append(playerQueue, player)
		if len(playerQueue) == 2 {
			gameId := fmt.Sprintf("%d", atomic.AddInt64(&core.idGenerator, 1))
			hub := newHub(core, gameId, playerQueue...)

			core.mux.Lock()
			core.hubs[gameId] = hub
			core.mux.Unlock()

			logrus.Infof("core: create hub (%s)", gameId)

			playerQueue = playerQueue[:0]
		}
	}
}
