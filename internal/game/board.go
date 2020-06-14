package game

type Board struct {
	Width  int
	Height int

	Cells [][]Cell
}

func initBoard(w, h int) Board {

	b := Board{
		Width:  w,
		Height: h,
		Cells:  [][]Cell{},
	}

	for i := 0; i < w; i++ {
		cells := []Cell{}
		for j := 0; j < h; j++ {
			cells = append(cells, initCell())
		}
		b.Cells = append(b.Cells, cells)
	}

	return b
}

func (b Board) isValidPosition(x, y int) bool {

	if x < 0 || x >= b.Width {
		return false
	}
	if y < 0 || y >= b.Height {
		return false
	}

	return true
}

func (b Board) isWin(x, y int) bool {

	if b.Cells[x][y].isFill() == false {
		return false
	}

	for i := 0; i < 8; i++ {

		var u = x
		var v = y

		var isWinner = true

		for j := 1; j < defaultWinLength; j++ {
			u = u + row[i]
			v = v + col[i]

			if b.isValidPosition(u, v) == false {
				isWinner = false
				break
			}

			if b.Cells[u][v].isFill() == false {
				isWinner = false
				break
			}

			if b.Cells[u][v].Player != b.Cells[x][y].Player {
				isWinner = false
				break
			}
		}

		if isWinner {
			return true
		}
	}

	return false
}

func (b Board) getStatus() Status {

	var anyAvaiableMove = false

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {

			if !b.Cells[i][j].isFill() && !anyAvaiableMove {
				anyAvaiableMove = true
			}

			if b.Cells[i][j].isFill() {

				var winner = b.isWin(i, j)

				if winner {
					switch b.Cells[i][j].Player {
					case XPlayer:
						return XWin
					case OPlayer:
						return OWin
					}
				}
			}

		}
	}

	if anyAvaiableMove {
		return Running
	} else {
		return Tie
	}
}
