package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	"github.com/chiendo97/caro-online/internal/socket"
	"github.com/sirupsen/logrus"
)

type Hub struct {
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
		core: core,
		key:  key,
		msgC: make(chan socket.Message),

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

/*
	core close -> close all hub.
	hub close -> close all socket (player)

	// in case, client send close msg -> need UnRegister
	socket close -> hub UnRegister
	hub UnRegister -> close that socket (player)

	// conclusion WRONG
	only UnRegister when client send close msg.
	not read closed network
	read close -> UnRegister

	// 2nd conclusion
	UnRegister when ???

	Client send close msg.
	Server read close msg.
	UnRegister

	// 3rd conclusion:
	UnRegister should not close socket.
	When server read close msg, close socket itself

	// 4th conclusion
	FUCKING DEADLOCK
	UnRegister call to for select
	it stucks because of closing all socket
	socket call to UnRegister

*/

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

	if hub.game.Status != game.Running {
		time.Sleep(5 * time.Second)
		hub.core.UnRegister(hub)
		close(hub.doneC)
	}
}

func (hub *Hub) subscribe(s socket.Socket) {
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

	hub.players[s] = player

	hub.playerWG.Add(1)
	go func() {
		err1, err2 := s.Run()

		if err1 != nil || err2 != nil {
			logrus.Errorf("%v:%v", err1, err2)
		}
		hub.playerWG.Done()
	}()

	hub.broadcast()

	logrus.Infof("hub %s: take new socket %s as player %d", hub.key, s.GetSocketIPAddress(), player)
}

func (hub *Hub) unsubscribe(s socket.Socket) {
	if _, ok := hub.players[s]; ok {

		delete(hub.players, s)

		hub.broadcast()

		if len(hub.players) == 1 {
			hub.core.Register(hub)
		} else {
			hub.core.UnRegister(hub)
		}

		logrus.Infof("hub %s: Player %s left", hub.key, s.GetSocketIPAddress())
	}
}

func (hub *Hub) run() error {

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
			hub.playerWG.Wait()
			return nil
		}
	}
}
