package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chiendo97/caro-online/internal/server"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func findHubHandler(core *server.CoreServer, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = server.InitMessage(conn, key)
	core.FindGame <- msg
}

func createHubHandler(core *server.CoreServer, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = server.InitMessage(conn, key)
	core.CreateGame <- msg
}

func joinHubHandler(core *server.CoreServer, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = r.URL.Query().Get("hub")
	log.Println("web: joining hub - ", key)
	var msg = server.InitMessage(conn, key)
	core.JoinGame <- msg
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to caro online. Please come to https://github.com/chiendo97/caro-online for introduction")
}

func main() {

	var core = server.InitCore()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create_hub", func(w http.ResponseWriter, r *http.Request) {
		createHubHandler(core, w, r)
	})
	http.HandleFunc("/join_hub", func(w http.ResponseWriter, r *http.Request) {
		joinHubHandler(core, w, r)
	})
	http.HandleFunc("/find_hub", func(w http.ResponseWriter, r *http.Request) {
		findHubHandler(core, w, r)
	})

	go core.Run()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is Running on %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
