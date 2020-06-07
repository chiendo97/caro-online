package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chiendo97/caro-online/internal/client"

	"github.com/chiendo97/caro-online/internal/socket"

	"github.com/gorilla/websocket"
)

func main() {
	var addr = os.Getenv("host")
	var port = os.Getenv("port")
	if addr == "" {
		if port == "" {
			port = "8080"
		}
		addr = "localhost:" + port
	}

	log.Printf("Client is connecting to %s", addr)

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

	var hub = client.InitHub()
	var socket = socket.InitSocket(c, &hub)
	hub.Socket = socket

	go hub.Run()

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
