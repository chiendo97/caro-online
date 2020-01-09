package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var core = initCore()
var upgrader = websocket.Upgrader{}

func findHubHandler(w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = InitMessage(conn, key)
	core.findGame <- msg
}

func createHubHandler(w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = InitMessage(conn, key)
	core.createGame <- msg
}

func joinHubHandler(w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = r.URL.Query().Get("hub")
	log.Println("web: joining hub - ", key)
	var msg = InitMessage(conn, key)
	core.joinGame <- msg
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to caro online")
	return
}

func main() {
	log.Println("Server is running")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create_hub", createHubHandler)
	http.HandleFunc("/join_hub", joinHubHandler)
	http.HandleFunc("/find_hub", findHubHandler)

	go core.run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
