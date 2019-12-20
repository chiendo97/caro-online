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

func createHubHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Starting new hub")

	conn, _ := upgrader.Upgrade(w, r, nil)

	key := "asdfasdf"

	_, ok := hubs[key]

	if ok {
		log.Panicln("Key duplicate:", key, conn.RemoteAddr())
	}

	var hub = InitHub()
	go hub.run()

	hubs[key] = &hub

	var s = socket.InitSocket(conn, &hub)
	hub.register <- &s

	go s.Read()
	go s.Write()

	select {
	case s.Message <- socket.GenerateErrMsg("hub: " + key):
	}
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		keys, ok := r.URL.Query()["keys"]

		if !ok {
			log.Println("Key missing")
		}

		if len(keys) < 1 {
			log.Println("Key missing")
		}

		fmt.Fprint(w, keys)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
