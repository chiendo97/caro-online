package server

import "github.com/sirupsen/logrus"

func (core *coreServer) Register(hub *Hub) {

	logrus.Infof("core: hub (%s) subscribe.", hub.key)

	core.availHub <- hub.key
}

func (core *coreServer) OnLeave(hub *Hub) {

	core.mux.Lock()
	defer core.mux.Unlock()

	if _, ok := core.hubs[hub.key]; ok {
		core.leaveHub(hub)
	}
}
