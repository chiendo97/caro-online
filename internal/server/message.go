package server

import "github.com/gorilla/websocket"

type msgStruct struct {
	conn   *websocket.Conn
	gameId string
}

func InitMessage(conn *websocket.Conn, gameId string) msgStruct {
	return msgStruct{
		conn:   conn,
		gameId: gameId,
	}
}
