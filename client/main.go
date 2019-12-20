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

	var socket socket.Socket
	var hub = InitHub()
	var c *websocket.Conn
	var err error

	switch option {
	case 0:
		c, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/create_hub", nil)
	case 1:
		c, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/join_hub", nil)
	default:
		log.Fatalln("Invalid option: ", option)
	}

	if err != nil {
		log.Fatal("Dial error: ", err)
	}

	socket = InitSocket(c, &hub)
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
