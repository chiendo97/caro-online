package game

import (
	"fmt"
	"math/rand"
)

type Bot interface {
	GetMove(Player, Game) (Move, error)
}

type RandomBot struct{}

func (bot *RandomBot) GetMove(p Player, g Game) (Move, error) {
	var x, y int
	count := 0
	for {
		if count == g.Board.Height*g.Board.Width {
			return Move{}, fmt.Errorf("Can not find any move")
		}
		x = rand.Intn(g.Board.Height)
		y = rand.Intn(g.Board.Width)
		if err := g.IsValidMove(Move{
			X:      x,
			Y:      y,
			Player: p,
		}); err != nil {
			count++
			continue
		}
		break
	}
	return Move{
		X:      x,
		Y:      y,
		Player: p,
	}, nil
}
