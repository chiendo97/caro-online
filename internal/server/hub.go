package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	debug    map[string]int
	core     *coreServer
	key      string
	game     game.Game
	players  map[socket.Socket]game.Player
	playerWG sync.WaitGroup

	msgC   chan socket.Message
	regC   chan socket.Socket
	unregC chan socket.Socket
	doneC  chan int
}

func initHub(core *coreServer, key string) *Hub {

	var hub = Hub{
		debug: make(map[string]int),
		core:  core,
		key:   key,
		msgC:  make(chan socket.Message),

		game:    game.InitGame(key),
		players: make(map[socket.Socket]game.Player),

		regC:   make(chan socket.Socket),
		unregC: make(chan socket.Socket),

		doneC: make(chan int),
	}

	return &hub
}

func (hub *Hub) HandleMsg(msg socket.Message) {
	hub.msgC <- msg
}

func (hub *Hub) UnRegister(s socket.Socket) {
	hub.unregC <- s
}

func (hub *Hub) Register(s socket.Socket) {
	hub.regC <- s
}

func (hub *Hub) broadcast() {

	if len(hub.players) < 2 {
		var msg = socket.GenerateAnnouncementMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for socket := range hub.players {
			socket.SendMessage(msg)
			logrus.Debugf("hub: send (%s) to (%s)", msg, socket.GetSocketIPAddress())
		}
	} else {
		for s, player := range hub.players {
			var game = hub.game
			var msg = socket.GenerateGameMsg(player, game)
			s.SendMessage(msg)
			logrus.Debugf("hub: send (%s) to (%s)", msg, s.GetSocketIPAddress())
		}
	}

}

func (hub *Hub) handleMsg(msg socket.Message) {

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

func (hub *Hub) subscribe(s socket.Socket) {
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
		hub.debug[s.GetSocketIPAddress()] = 1
		err1, err2 := s.Run()

		if err1 != nil || err2 != nil {
			logrus.Errorf("Socket run err: %v:%v", err1, err2)
		}
		hub.playerWG.Done()
		delete(hub.debug, s.GetSocketIPAddress())
	}()

	hub.broadcast()

	logrus.Debugf("hub %s: take new socket %s as player %d", hub.key, s.GetSocketIPAddress(), player)
}

func (hub *Hub) unsubscribe(s socket.Socket) {
	s.CloseMessage()

	if _, ok := hub.players[s]; ok {

		delete(hub.players, s)

		hub.broadcast()

		if len(hub.players) == 1 {
			go hub.core.Register(hub)
		} else {
			go hub.core.UnRegister(hub)
		}
	}

}

func (hub *Hub) run() error {

	logrus.Debugf("Hub %v start", hub.key)
	defer logrus.Debugf("Hub %v stop", hub.key)

	var debug = ""
	var ctx, cancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				logrus.Debugf("hub:(%v) debug:(%v) Sockets:(%v)", hub.key, debug, hub.debug)
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case msg := <-hub.msgC:
			debug = fmt.Sprintf("1 %v", msg)
			hub.handleMsg(msg)
		case socket := <-hub.regC:
			debug = fmt.Sprintf("2 %v", socket.GetSocketIPAddress())
			hub.subscribe(socket)
		case socket := <-hub.unregC:
			debug = fmt.Sprintf("3 %v", socket.GetSocketIPAddress())
			hub.unsubscribe(socket)
		case <-hub.doneC:
			debug = fmt.Sprintf("4 ")
			for socket := range hub.players {
				logrus.Debugf("Hub %v calling stop %s", hub.key, socket.GetSocketIPAddress())
				hub.unsubscribe(socket)
			}
			hub.playerWG.Wait()
			cancel()
			return nil
		}
		debug = ""
	}
}
