package api

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/internal/server"
)

type service struct {
	upgrader websocket.Upgrader
	server   *http.Server
	core     server.CoreServer
}

func InitService(core server.CoreServer, port int) *service {
	s := &service{
		upgrader: websocket.Upgrader{},
		core:     core,
		server:   &http.Server{Addr: fmt.Sprintf(":%d", port)},
	}
	s.buildAPI()
	return s
}

func (s *service) buildAPI() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to caro online. Please come to https://github.com/chiendo97/caro-online for introduction")
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/create_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		s.core.CreateGame(conn)
	})
	http.HandleFunc("/join_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		var gameID = r.URL.Query().Get("hub")
		s.core.JoinGame(conn, gameID)
	})
	http.HandleFunc("/find_hub", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)

		s.core.FindGame(conn)
	})
}

func (s *service) ListenAndServe(port int) error {
	logrus.Info("Server is running on port ", port)
	defer logrus.Info("Server stop")

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *service) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
