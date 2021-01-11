package server

func (core *coreServer) OnLeave(hub *Hub) {
	core.mux.Lock()
	defer core.mux.Unlock()

	if _, found := core.hubs[hub.key]; found {
		core.leaveHub(hub)
	}
}
