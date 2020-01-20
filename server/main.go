package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func findHubHandler(core *coreServer, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = InitMessage(conn, key)
	core.findGame <- msg
}

func createHubHandler(core *coreServer, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	var key = ""
	var msg = InitMessage(conn, key)
	core.createGame <- msg
}

func joinHubHandler(core *coreServer, w http.ResponseWriter, r *http.Request) {

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

	var core = initCore()

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

	go core.run()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
