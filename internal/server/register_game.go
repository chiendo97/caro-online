package server

import "github.com/sirupsen/logrus"

func (core *coreServer) OnLeave(hub *Hub) {
	core.mux.Lock()
	defer core.mux.Unlock()

	logrus.Infof("core: delete hub (%s)", hub.key)
	delete(core.hubs, hub.key)
}
