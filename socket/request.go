package socket

import (
	"fmt"
	"helloworld/caro/game"
)

const (
	MoveMessage = "move"
	GameMessage = "game"
	MsgMessage  = "msg"
)

type Message struct {
	Kind string

	Move game.Move
	Game game.Game
	Msg  string
}

func (msg Message) String() string {
	switch msg.Kind {
	case MoveMessage:
		return fmt.Sprint("Move msg: ", msg.Move)
	case GameMessage:
		return fmt.Sprint("Game msg: ", msg.Game)
	case MsgMessage:
		return fmt.Sprint("Msg msg: ", msg.Msg)
	}
	return fmt.Sprint("Unknown msg kind: ", msg.Kind)
}

func GenerateMoveMsg(move game.Move) Message {
	return Message{
		Kind: MoveMessage,
		Move: move,
	}
}

func GenerateGameMsg(game game.Game) Message {
	return Message{
		Kind: GameMessage,
		Game: game,
	}
}

func GenerateErrMsg(err string) Message {
	return Message{
		Kind: MsgMessage,
		Msg:  err,
	}
}
