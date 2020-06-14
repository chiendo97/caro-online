package api

import (
	"fmt"
	"net/http"

	"github.com/chiendo97/caro-online/internal/server"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type service struct {
	upgrader websocket.Upgrader
	core     server.CoreServer
}

func InitService(core server.CoreServer) *service {
	s := &service{
		upgrader: websocket.Upgrader{},
		core:     core,
	}
	s.buildAPI()
	return s
}

func (s *service) buildAPI() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to caro online. Please come to https://github.com/chiendo97/caro-online for introduction")
	})
	http.HandleFunc("/create_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		var key = ""
		var msg = server.InitMessage(conn, key)
		s.core.CreateGame(msg)
	})
	http.HandleFunc("/join_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		var key = r.URL.Query().Get("hub")
		var msg = server.InitMessage(conn, key)
		s.core.JoinGame(msg)
	})
	http.HandleFunc("/find_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		var key = ""
		var msg = server.InitMessage(conn, key)
		s.core.FindGame(msg)
	})
}

func (s *service) ListenAndServe(port int) error {
	log.Info("Server is running on port ", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
