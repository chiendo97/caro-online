package main

import (
	"helloworld/caro/socket"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var hubs = make(map[string]*Hub)
var upgrader = websocket.Upgrader{}

func createHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Starting new hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	key := "asdfasdf"

	_, ok := hubs[key]

	if ok {
		log.Panicln("Key duplicate:", key)
	}

	var hub = InitHub()
	go hub.run()

	hubs[key] = &hub

	var s = socket.InitSocket(conn, &hub)
	hub.register <- &s

	go s.Read()
	go s.Write()
}

func joinHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Joining hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	key := "asdfasdf"

	hub, ok := hubs[key]

	if !ok {
		return
	}

	var s = socket.InitSocket(conn, hub)
	hub.register <- &s

	go s.Read()
	go s.Write()

}

func main() {
	log.Println("Server is running")

	http.HandleFunc("/create_hub", createHubHandler)
	http.HandleFunc("/join_hub", joinHubHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
