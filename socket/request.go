package socket

import (
	"fmt"
	"helloworld/caro/game"
)

type Message struct {
	Kind string
	Move game.Move
	Game game.Game
	Msg  string
}

func (msg Message) String() string {
	switch msg.Kind {
	case "move":
		return fmt.Sprint("Move msg: ", msg.Move)
	case "game":
		return fmt.Sprint("Game msg: ", msg.Game)
	case "error":
		return fmt.Sprint("Msg msg: ", msg.Msg)
	}
	return fmt.Sprint("Empty msg")
}

func GenerateMoveMsg(move game.Move) Message {
	return Message{
		Kind: "move",
		Move: move,
	}
}

func GenerateGameMsg(game game.Game) Message {
	return Message{
		Kind: "game",
		Game: game,
	}
}

func GenerateErrMsg(err string) Message {
	return Message{
		Kind: "error",
		Msg:  err,
	}
}
