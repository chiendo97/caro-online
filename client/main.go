package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chiendo97/caro-online/socket"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("Client is running")

	var addr = os.Getenv("host")
	if addr == "" {
		addr = "localhost:8080"
	}

	// === Take options
	var args = os.Args
	var host string

	switch len(args) {
	case 1:
		host = "ws://" + addr + "/find_hub"
	case 2:
		host = "ws://" + addr + "/create_hub"
	case 3:
		var param = args[2]
		host = "ws://" + addr + "/join_hub" + "?" + "hub=" + param
	default:
		log.Fatalln("Invalid option")
	}

	// === Init socket and hub
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

	// === take interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			log.Fatalln("Exit client")
		}

	}

}
