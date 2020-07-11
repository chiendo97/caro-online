package server

import (
	"context"
	"sync"
	"time"

	"github.com/chiendo97/caro-online/internal/socket"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CoreServer interface {
	Run() error
	Stop()

	FindGame(msg msgStruct)
	JoinGame(msg msgStruct)
	CreateGame(msg msgStruct)
}

type coreServer struct {
	hubs  map[string]*Hub
	hubs2 map[string]*Hub

	availHubC chan string
	hubWg     sync.WaitGroup

	findC   chan msgStruct
	joinC   chan msgStruct
	createC chan msgStruct

	regC   chan *Hub
	unregC chan *Hub

	done chan int
}

func InitCoreServer() CoreServer {

	var core = coreServer{
		hubs:      make(map[string]*Hub),
		hubs2:     make(map[string]*Hub),
		availHubC: make(chan string, 5),

		findC:   make(chan msgStruct),
		joinC:   make(chan msgStruct),
		createC: make(chan msgStruct),

		regC:   make(chan *Hub),
		unregC: make(chan *Hub),

		done: make(chan int),
	}

	return &core
}

func (core *coreServer) Test(hub *Hub) {

	if _, ok := core.hubs[hub.key]; ok {
		log.Infof("core: delete hub (%s)", hub.key)
		close(hub.doneC)
		delete(core.hubs, hub.key)
	}
}

func (core *coreServer) Run() error {

	log.Infof("Core start")
	defer log.Infof("Core stop")

	for {
		select {
		case <-core.done:
			// stop core
			for _, hub := range core.hubs {
				core.Test(hub)
			}
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				select {
				case <-time.After(time.Second):
					log.Debugf("Hubs: %v", core.hubs2)
				case <-ctx.Done():
					return
				}
			}()
			core.hubWg.Wait()
			cancel()
			return nil
		case hub := <-core.regC:
			// hub -> core
			log.Infof("core: hub (%s) subscribe.", hub.key)

			go func() {
				core.availHubC <- hub.key
			}()
		case hub := <-core.unregC:
			// hub <- core
			core.Test(hub)
		case msg := <-core.findC:
			// find hub from core
			log.Infof("core: socket (%s) find game", msg.conn.RemoteAddr())

			go func() {
				for {
					select {
					case gameID := <-core.availHubC:
						if _, ok := core.hubs[gameID]; !ok {
							continue
						}
						msg.gameId = gameID
						core.joinC <- msg
						return
					case <-time.After(1 * time.Second):
						core.createC <- msg
						return
					}
				}
			}()
		case msg := <-core.joinC:
			// join hub from core
			log.Infof("core: socket (%s) join hub (%s)", msg.conn.RemoteAddr(), msg.gameId)

			hub, ok := core.hubs[msg.gameId]

			if !ok {
				log.Info("core: hub not found - ", msg.gameId, msg.conn.RemoteAddr())
				msg.conn.WriteMessage(websocket.CloseMessage, []byte{})
			} else {
				go func() {
					hub.Register(socket.InitSocket(msg.conn, hub))
				}()
			}

		case msg := <-core.createC:
			// core create a new hub

			var gameId = uuid.New().String()[:8]

			_, ok := core.hubs[gameId]
			if ok {
				log.Error("Key duplicate: ", gameId, msg.conn.RemoteAddr())
			}

			var hub = initHub(core, gameId)

			core.hubs[gameId] = hub
			core.hubs2[gameId] = hub

			go func() {
				core.availHubC <- hub.key
			}()

			core.hubWg.Add(1)
			go func() {
				err := hub.run()
				if err != nil {
					log.Errorf("hub run error: %v", err)
				}
				core.hubWg.Done()
				delete(core.hubs2, hub.key)
			}()

			go func() {
				hub.Register(socket.InitSocket(msg.conn, hub))
			}()

			log.Infof("core: socket (%s) create hub (%s)", msg.conn.RemoteAddr(), gameId)
		}
	}
}

func (core *coreServer) Stop() {
	close(core.done)
}
