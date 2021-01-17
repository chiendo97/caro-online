package socket

import (
	"fmt"

	"github.com/chiendo97/caro-online/internal/game"
)

type MessageType int

const (
	Move MessageType = iota
	Game
	Announce
)

type Message struct {
	Type MessageType

	Player   game.Player
	Move     game.Move
	Game     game.Game
	Announce string
}

func (msg Message) String() string {
	switch msg.Type {
	case Move:
		return fmt.Sprint("Move: ", msg.Move)
	case Game:
		return fmt.Sprint("Game: ", msg.Game)
	case Announce:
		return fmt.Sprint("Announcement: ", msg.Announce)
	}

	return fmt.Sprint("Unknown msg kind: ", msg.Type)
}

func NewMoveMsg(move game.Move) Message {
	return Message{
		Type: Move,
		Move: move,
	}
}

func NewGameMsg(p game.Player, g game.Game) Message {
	return Message{
		Type:   Game,
		Player: p,
		Game:   g,
	}
}

func NewAnnouncementMsg(message string) Message {
	return Message{
		Type:     Announce,
		Announce: message,
	}
}
