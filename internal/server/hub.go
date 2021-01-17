package server

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
)

type Hub struct {
	key     string
	core    *coreServer
	game    game.Game
	players map[socket.Socket]game.Player

	register   chan socket.Socket
	unregister chan socket.Socket
	message    chan socket.Message
}

func newHub(core *coreServer, key string, conns ...*websocket.Conn) *Hub {
	var hub = &Hub{
		key:     key,
		core:    core,
		game:    game.InitGame(key),
		players: make(map[socket.Socket]game.Player),

		register:   make(chan socket.Socket),
		unregister: make(chan socket.Socket),
		message:    make(chan socket.Message),
	}

	for _, conn := range conns {
		s := socket.NewSocket(conn, hub)
		hub.onEnter(s)
	}

	go hub.run()

	return hub
}

func (hub *Hub) OnMessage(msg socket.Message) {
	hub.message <- msg
}

func (hub *Hub) OnLeave(s socket.Socket) {
	hub.unregister <- s
}

func (hub *Hub) OnEnter(s socket.Socket) {
	hub.register <- s
}

func (hub *Hub) onMessage(msg socket.Message) {
	if msg.Type != socket.Move {
		logrus.Errorf("hub %s: No msg kind case %s", hub.key, msg)
		return
	}

	var err error
	hub.game, err = hub.game.TakeMove(msg.Move)
	if err != nil {
		logrus.Errorf("hub %s: game error - %s", hub.key, err)
	}

	hub.broadcast()
}

func (hub *Hub) onLeave(s socket.Socket) {
	if _, found := hub.players[s]; !found {
		return
	}

	s.Stop()
	delete(hub.players, s)
}

func (hub *Hub) onEnter(s socket.Socket) {
	if len(hub.players) == 2 {
		logrus.Debugf("hub %s: room is full %s", hub.key, s.GetSocketIPAddress())
		s.Stop()
		return
	}

	var player = game.XPlayer
	for _, otherId := range hub.players {
		switch otherId {
		case game.XPlayer:
			player = game.OPlayer
		case game.OPlayer:
			player = game.XPlayer
		}
	}

	hub.players[s] = player

	go func() {
		err1, err2 := s.Run()
		if err1 != nil || err2 != nil {
			logrus.Errorf("Socket run err: %v:%v", err1, err2)
		}
	}()

	hub.broadcast()

	logrus.Debugf("hub %s: take new socket %s as player %d", hub.key, s.GetSocketIPAddress(), player)
}

func (hub *Hub) broadcast() {
	if len(hub.players) < 2 {
		var msg = socket.NewAnnouncementMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for s := range hub.players {
			s.SendMessage(msg)
		}
	} else {
		for s, player := range hub.players {
			s.SendMessage(socket.NewGameMsg(player, hub.game))
		}
	}
}

func (hub *Hub) run() {
	defer hub.core.OnLeave(hub)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case s := <-hub.register:
			hub.onEnter(s)
		case s := <-hub.unregister:
			hub.onLeave(s)
		case msg := <-hub.message:
			hub.onMessage(msg)
		case <-ticker.C:
			if hub.game.Status != game.Running {
				for s := range hub.players {
					s.Stop()
				}
			}
			if len(hub.players) == 0 {
				return
			}
		}
	}
}
