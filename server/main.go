package main

import (
	"fmt"
	"helloworld/caro/socket"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var hubs = make(map[string]*Hub)
var upgrader = websocket.Upgrader{}

func findHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Finding hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key string

	for k, v := range hubs {

		if len(v.players) < 2 { // hub not enough players

			key = k
			break
		}
	}

	hub, ok := hubs[key]

	if !ok {
		// TODO: generate random key
		key := "asdfasdf"

		_, ok := hubs[key]

		if ok {
			log.Panicln("Key duplicate:", key, conn.RemoteAddr())
		}

		hub = InitHub(key)
		go hub.run()

		hubs[key] = hub
	}

	var s = socket.InitSocket(conn, hub)
	hub.register <- &s

	go s.Read()
	go s.Write()
}

func createHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Starting new hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	// TODO: generate random key
	key := "asdfasdf"

	_, ok := hubs[key]

	if ok {
		log.Panicln("Key duplicate:", key, conn.RemoteAddr())
	}

	var hub = InitHub(key)
	go hub.run()

	hubs[key] = hub

	var s = socket.InitSocket(conn, hub)
	hub.register <- &s

	go s.Read()
	go s.Write()

	hub.message <- socket.GenerateErrMsg("Hub key: " + key)
}

func joinHubHandler(w http.ResponseWriter, r *http.Request) {

	var key = r.URL.Query().Get("hub")

	log.Println("Joining hub: ", key)

	conn, _ := upgrader.Upgrade(w, r, nil)

	hub, ok := hubs[key]

	if !ok {
		log.Println("No available hub:", key, conn.RemoteAddr())
		conn.WriteMessage(websocket.CloseMessage, []byte{})
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
	http.HandleFunc("/find_hub", findHubHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to caro-online")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
