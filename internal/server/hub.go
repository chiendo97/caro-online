package server

import (
	"fmt"
	"time"

	"github.com/chiendo97/caro-online/internal/game"
	soc "github.com/chiendo97/caro-online/internal/socket"
	log "github.com/sirupsen/logrus"
)

type Hub struct {
	core    *coreServer
	key     string
	game    game.Game
	players map[*soc.Socket]game.Player

	msgC   chan soc.Message
	regC   chan *soc.Socket
	unregC chan *soc.Socket
	doneC  chan int
}

func (hub *Hub) HandleMsg(msg soc.Message) {
	hub.msgC <- msg
}

func (hub *Hub) UnRegister(s *soc.Socket) {
	hub.unregC <- s
}

func (hub *Hub) Register(s *soc.Socket) {
	hub.regC <- s
}

func initHub(core *coreServer, key string) *Hub {

	var hub = Hub{
		core: core,
		key:  key,
		msgC: make(chan soc.Message),

		game:    game.InitGame(key),
		players: make(map[*soc.Socket]game.Player),

		regC:   make(chan *soc.Socket),
		unregC: make(chan *soc.Socket),

		doneC: make(chan int),
	}

	go hub.run()

	return &hub
}

func (hub *Hub) broadcast() {

	if len(hub.players) < 2 {
		var msg = soc.GenerateAnnouncementMsg(fmt.Sprintf("hub %s: wait for players", hub.key))
		for socket := range hub.players {
			socket.SendMessage(msg)
			log.Infof("hub: send (%s) to (%s)", msg, socket.GetSocketIPAddress())
		}
	} else {
		for socket, player := range hub.players {
			var game = hub.game
			var msg = soc.GenerateGameMsg(player, game)
			socket.SendMessage(msg)
			log.Infof("hub: send (%s) to (%s)", msg, socket.GetSocketIPAddress())
		}
	}

}

func (hub *Hub) handleMsg(msg soc.Message) {

	if msg.Type != soc.MoveMessageType {
		log.Errorf("hub %s: No msg kind case %s", hub.key, msg)
		return
	}

	g, err := hub.game.TakeMove(msg.Move)

	if err != nil {
		log.Infof("hub %s: game error - %s", hub.key, err)
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

func (hub *Hub) subscribe(socket *soc.Socket) {
	if len(hub.players) == 2 {
		log.Infof("hub %s: room is full %s", hub.key, socket.GetSocketIPAddress())
		socket.CloseMessage()
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

	log.Infof("hub %s: take new socket %s as player %d", hub.key, socket.GetSocketIPAddress(), player)
	hub.players[socket] = player

	hub.broadcast()
}

func (hub *Hub) unsubscribe(socket *soc.Socket) {
	if _, ok := hub.players[socket]; ok {
		log.Infof("hub %s: Player %s left", hub.key, socket.GetSocketIPAddress())

		delete(hub.players, socket)
		socket.CloseMessage()

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
