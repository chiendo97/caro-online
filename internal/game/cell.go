package game

type Cell struct {
	Player Player
}

func initCell() Cell {
	return Cell{
		Player: EPlayer,
	}
}

func (c Cell) isFill() bool {
	return c.Player != EPlayer
}
