package server

func (core *coreServer) Register(hub *Hub) {
	core.regC <- hub
}
func (core *coreServer) UnRegister(hub *Hub) {
	core.unregC <- hub
}
func (core *coreServer) FindGame(msg msgStruct) {
	core.findC <- msg
}
func (core *coreServer) JoinGame(msg msgStruct) {
	core.joinC <- msg
}
func (core *coreServer) CreateGame(msg msgStruct) {
	core.createC <- msg
}
