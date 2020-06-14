package socket

import (
	"fmt"

	"github.com/chiendo97/caro-online/internal/game"
)

type MessageType int

const (
	MoveMessageType MessageType = iota
	GameMessageType
	AnnouncementMessageType
)

const (
	MoveMessage = "move"
	GameMessage = "game"
	MsgMessage  = "msg"
)

type Message struct {
	Type MessageType

	Player       game.Player
	Move         game.Move
	Game         game.Game
	Announcement string
}

func (msg Message) String() string {
	switch msg.Type {
	case MoveMessageType:
		return fmt.Sprint("Move: ", msg.Move)
	case GameMessageType:
		return fmt.Sprint("Game: ", msg.Game)
	case AnnouncementMessageType:
		return fmt.Sprint("Announcement: ", msg.Announcement)
	}

	return fmt.Sprint("Unknown msg kind: ", msg.Type)
}

func GenerateMoveMsg(move game.Move) Message {
	return Message{
		Type: MoveMessageType,
		Move: move,
	}
}

func GenerateGameMsg(p game.Player, g game.Game) Message {
	return Message{
		Type:   GameMessageType,
		Player: p,
		Game:   g,
	}
}

func GenerateAnnouncementMsg(message string) Message {
	return Message{
		Type:         AnnouncementMessageType,
		Announcement: message,
	}
}
