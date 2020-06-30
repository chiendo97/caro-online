package client

import (
	"math/rand"

	"github.com/chiendo97/caro-online/internal/game"
)

type Bot interface {
	GetMove(game.Player, game.Game) (game.Move, error)
}

type RandomBot struct{}

func (bot *RandomBot) GetMove(p game.Player, g game.Game) (game.Move, error) {
	var x, y int
	for {
		x = rand.Intn(g.Board.Height)
		y = rand.Intn(g.Board.Width)
		if err := g.IsValidMove(game.Move{
			X:      x,
			Y:      y,
			Player: p,
		}); err != nil {
			continue
		}
		break
	}
	return game.Move{
		X:      x,
		Y:      y,
		Player: p,
	}, nil
}
