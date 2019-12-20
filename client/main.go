package main

import (
	"helloworld/caro/socket"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("Client is running")

	// {{{ === Take options
	var args = os.Args
	var host string

	switch len(args) {
	case 1:
		log.Fatalln("No option")
	case 2:
		host = "ws://localhost:8080/create_hub"
	case 3:
		host = "ws://localhost:8080/join_hub"
	default:
		log.Fatalln("Invalid option")
	}
	// }}}

	// {{{ === Init socket and hub
	var c, _, err = websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("Dial error: ", err)
	}

	var hub = InitHub()
	var socket = socket.InitSocket(c, &hub)
	hub.socket = &socket

	go hub.run()
	go socket.Read()
	go socket.Write()
	//}}}

	// {{{ === take interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			log.Fatalln("Exit")
		}

	}
	// }}}

}
