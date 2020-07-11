package server

import (
	"fmt"
	"sync"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	mux      sync.Mutex
	playerWG sync.WaitGroup

	core    *coreServer
	key     string
	game    game.Game
	players map[socket.Socket]game.Player

	doneC chan int
}

func initHub(core *coreServer, key string) *Hub {

	var hub = Hub{
		core:    core,
		key:     key,
		game:    game.InitGame(key),
		players: make(map[socket.Socket]game.Player),

		doneC: make(chan int),
	}

	return &hub
}

func (hub *Hub) OnMessage(msg socket.Message) {

	hub.mux.Lock()
	defer hub.mux.Unlock()

	if msg.Type != socket.MoveMessageType {
		logrus.Errorf("hub %s: No msg kind case %s", hub.key, msg)
		return
	}

	g, err := hub.game.TakeMove(msg.Move)

	if err != nil {
		logrus.Debugf("hub %s: game error - %s", hub.key, err)
	} else {
		hub.game = g
	}

	hub.broadcast()

	if hub.game.Status != game.Running {
		go hub.core.UnRegister(hub)
	}
}

func (hub *Hub) UnRegister(s socket.Socket) {
	hub.mux.Lock()
	defer hub.mux.Unlock()

	if _, ok := hub.players[s]; ok {
		delete(hub.players, s)

		hub.broadcast()

		if len(hub.players) == 1 {
			go hub.core.Register(hub)
		} else {
			go hub.core.UnRegister(hub)
		}

		s.CloseMessage()
	}
}

func (hub *Hub) OnEnter(s socket.Socket) {

	hub.mux.Lock()
	defer hub.mux.Unlock()

	if len(hub.players) == 2 {
		logrus.Debugf("hub %s: room is full %s", hub.key, s.GetSocketIPAddress())
		s.CloseMessage()
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

	hub.playerWG.Add(1)
	go func() {
		err1, err2 := s.Run()

		if err1 != nil || err2 != nil {
			logrus.Errorf("Socket run err: %v:%v", err1, err2)
		}
		hub.playerWG.Done()
	}()

	hub.broadcast()

	logrus.Debugf("hub %s: take new socket %s as player %d", hub.key, s.GetSocketIPAddress(), player)
}

func (hub *Hub) broadcast() {

	if len(hub.players) < 2 {
		var msg = socket.GenerateAnnouncementMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for socket := range hub.players {
			socket.SendMessage(msg)
			// logrus.Debugf("hub %s: send (%s) to (%s)", hub.key, msg, socket.GetSocketIPAddress())
		}
	} else {
		for s, player := range hub.players {
			var game = hub.game
			var msg = socket.GenerateGameMsg(player, game)
			s.SendMessage(msg)
			// logrus.Debugf("hub %s: send (%s) to (%s)", hub.key, msg, s.GetSocketIPAddress())
		}
	}

}

func (hub *Hub) Run() error {

	for {
		select {
		case <-hub.doneC:
			for socket := range hub.players {
				hub.UnRegister(socket)
			}
			hub.playerWG.Wait()
			return nil
		}
	}
}
