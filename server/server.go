package main

import (
	"helloworld/caro/socket"
	"log"

	"github.com/gorilla/websocket"
)

type msgServer struct {
	socket *websocket.Conn
	gameId string
}

type coreServer struct {
	hubs          map[string]*Hub
	availableHubs chan string

	findGame   chan msgServer
	joinGame   chan msgServer
	createGame chan msgServer

	register   chan *Hub
	unregister chan *Hub

	done chan int
}

func InitServer() *coreServer {
	return &coreServer{}
}

// func (core *coreServer) createHub() {
// }

// func (core *coreServer) joinHub() {

// }
// func (core *coreServer) findHub() {
// }

func (core *coreServer) run() {

	go func() {
		for {
			select {
			case <-core.done:
				close(core.findGame)
				close(core.createGame)
				close(core.joinGame)
				close(core.register)
				close(core.unregister)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-core.register:
			}
		}
	}()
	go func() {
		for {
			select {
			case <-core.unregister:
			}
		}
	}()
	go func() {
		for {
			select {
			case msg, ok := <-core.findGame:

				if !ok {
					return
				}

				log.Println("Finding game")

				var gameId = <-core.availableHubs
				msg.gameId = gameId

				core.joinGame <- msg
			}
		}
	}()
	go func() {
		for {
			select {
			case msg, ok := <-core.joinGame:

				if !ok {
					return
				}

				log.Println("Joining game: ", msg.gameId)

				hub, ok := core.hubs[msg.gameId]

				if !ok {
					log.Println("Hub not available:", msg.gameId, msg.socket.RemoteAddr())
					msg.socket.WriteMessage(websocket.CloseMessage, []byte{})
					continue
				}

				var s = socket.InitSocket(msg.socket, hub)
				hub.register <- &s

				go s.Read()
				go s.Write()
			}
		}
	}()
	go func() {
		for {
			select {
			case msg, ok := <-core.createGame:

				if !ok {
					return
				}

				log.Println("Creating game")

				var gameId = "asdfasdf"

				_, ok = core.hubs[gameId]

				if ok {
					log.Panicln("Key duplicate: ", gameId, msg.socket.RemoteAddr())
				}

				var hub = InitHub(gameId)
				go hub.run()

				hubs[gameId] = hub

				var s = socket.InitSocket(msg.socket, hub)
				hub.register <- &s

				go s.Read()
				go s.Write()
			}
		}
	}()
}
