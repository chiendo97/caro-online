package main

import (
	"fmt"
	"helloworld/caro/socket"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("Client is running")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var option int

	fmt.Println(`Options:
	- Option 0: create new hub
	- Option 1: join current hub
	`)
	fmt.Print("Option: ")
	fmt.Scanln(&option)

	var host string

	switch option {
	case 0:
		host = "ws://localhost:8080/create_hub"
	case 1:
		host = "ws://localhost:8080/join_hub"
	default:
		log.Fatalln("Invalid option: ", option)
	}

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

	for {
		select {
		case <-interrupt:
			log.Fatalln("Exit")
		}

	}
}
