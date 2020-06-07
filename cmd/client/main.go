package main

import (
	"fmt"
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
	if port == "" {
		port = "8080"
	}
	if addr == "" {
		addr = "localhost"
	}

	log.Printf("Client is connecting to %s:%s", addr, port)

	// === Take options
	var args = os.Args
	var host string

	switch len(args) {
	case 1:
		host = fmt.Sprintf("ws://%s:%s/find_hub", addr, port)
	case 2:
		host = fmt.Sprintf("ws://%s:%s/create_hub", addr, host)
	case 3:
		var hubID = args[2]
		host = fmt.Sprintf("ws://%s:%s/join_hub?hub=%s", addr, host, hubID)
	default:
		log.Fatalln("Invalid option")
	}

	// === Init socket and hub
	c, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("Dial error: ", err)
	}

	hub := client.InitAndRunHub()
	hub.Socket = socket.InitAndRunSocket(c, hub)

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
