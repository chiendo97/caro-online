package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var core = InitServer()
var upgrader = websocket.Upgrader{}

func findHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Finding hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = "asdf"
	var msg = InitMessage(conn, key)
	core.findGame <- msg
}

func createHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Starting new hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	// TODO: generate random key
	key := "asdfasdf"

	var msg = InitMessage(conn, key)
	core.createGame <- msg
}

func joinHubHandler(w http.ResponseWriter, r *http.Request) {

	var key = r.URL.Query().Get("hub")

	log.Println("Joining hub: ", key)

	conn, _ := upgrader.Upgrade(w, r, nil)

	var msg = InitMessage(conn, key)
	core.joinGame <- msg
}

func main() {
	log.Println("Server is running")

	http.HandleFunc("/create_hub", createHubHandler)
	http.HandleFunc("/join_hub", joinHubHandler)
	http.HandleFunc("/find_hub", findHubHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to caro-online")
	})

	core.run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
