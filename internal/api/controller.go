package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chiendo97/caro-online/internal/server"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	//     defer wg.Done()
	//     err := s.core.Run()
	//     if err != nil {
	//         logrus.Errorf("Core run error: %v", err)
	//         return
	//     }
	// }()

	// defer func() {
	//     s.core.Stop()
	//     wg.Wait()
	// }()

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	logrus.Info("Server stop")

	return nil
}

func (s *service) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
