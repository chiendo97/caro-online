package server

import (
	"fmt"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	core    *coreServer
	key     string
	game    game.Game
	players map[*socket.Socket]game.Player

	msgC   chan socket.Message
	regC   chan *socket.Socket
	unregC chan *socket.Socket
	doneC  chan int
}

func (hub *Hub) HandleMsg(msg socket.Message) {
	hub.msgC <- msg
}

func (hub *Hub) UnRegister(s *socket.Socket) {
	hub.unregC <- s
}

func (hub *Hub) Register(s *socket.Socket) {
	hub.regC <- s
}

func initHub(core *coreServer, key string) *Hub {

	var hub = Hub{
		core: core,
		key:  key,
		msgC: make(chan socket.Message),

		game:    game.InitGame(key),
		players: make(map[*socket.Socket]game.Player),

		regC:   make(chan *socket.Socket),
		unregC: make(chan *socket.Socket),

		doneC: make(chan int),
	}

	go hub.run()

	return &hub
}

func (hub *Hub) broadcast() {

	if len(hub.players) < 2 {
		var msg = socket.GenerateAnnouncementMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for socket := range hub.players {
			socket.SendMessage(msg)
			logrus.Infof("hub: send (%s) to (%s)", msg, socket.GetSocketIPAddress())
		}
	} else {
		for s, player := range hub.players {
			var game = hub.game
			var msg = socket.GenerateGameMsg(player, game)
			s.SendMessage(msg)
			logrus.Infof("hub: send (%s) to (%s)", msg, s.GetSocketIPAddress())
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
		logrus.Infof("hub %s: game error - %s", hub.key, err)
	} else {
		hub.game = g
	}

	hub.broadcast()

	if hub.game.GetStatus() != game.Running {
		time.Sleep(5 * time.Second)
		hub.core.UnRegister(hub)
		close(hub.doneC)
	}
}

func (hub *Hub) subscribe(s *socket.Socket) {
	if len(hub.players) == 2 {
		logrus.Infof("hub %s: room is full %s", hub.key, s.GetSocketIPAddress())
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

	logrus.Infof("hub %s: take new socket %s as player %d", hub.key, s.GetSocketIPAddress(), player)
	hub.players[s] = player

	hub.broadcast()
}

func (hub *Hub) unsubscribe(s *socket.Socket) {
	if _, ok := hub.players[s]; ok {
		logrus.Infof("hub %s: Player %s left", hub.key, s.GetSocketIPAddress())

		delete(hub.players, s)
		s.CloseMessage()

		hub.broadcast()

		if len(hub.players) == 1 {
			hub.core.Register(hub)
		} else {
			hub.core.UnRegister(hub)
		}

	}
}

func (hub *Hub) run() {

	for {
		select {
		case msg := <-hub.msgC:
			hub.handleMsg(msg)
		case socket := <-hub.regC:
			hub.subscribe(socket)
		case socket := <-hub.unregC:
			hub.unsubscribe(socket)
		case <-hub.doneC:
			for socket := range hub.players {
				socket.CloseMessage()
			}
			return
		}
	}
}
