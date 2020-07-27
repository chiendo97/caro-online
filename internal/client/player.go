package client

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/chiendo97/caro-online/internal/game"
)

type Player struct{}

var reader = bufio.NewReader(os.Stdin)

func (play *Player) GetMove(p game.Player, g game.Game) (game.Move, error) {

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			return game.Move{}, err
		}

		text = strings.Replace(text, "\n", "", -1)

		steps := strings.Split(text, ":")

		if len(steps) != 2 {
			continue
		}

		x, err := strconv.Atoi(steps[0])
		if err != nil {
			continue
		}

		y, err := strconv.Atoi(steps[1])
		if err != nil {
			continue
		}

		return game.Move{
			X:      x,
			Y:      y,
			Player: p,
		}, nil
	}

}
