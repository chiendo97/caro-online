package server

import "github.com/sirupsen/logrus"

func (core *coreServer) leaveHub(hub *Hub) {

	logrus.Infof("core: delete hub (%s)", hub.key)

	// hub.Stop()
	delete(core.hubs, hub.key)
}

// func (core *coreServer) leaveAllHubs() {

//     core.mux.Lock()
//     defer core.mux.Unlock()

//     for _, hub := range core.hubs {
//         core.leaveHub(hub)
//     }
// }
