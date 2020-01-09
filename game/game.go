package game

import (
	"errors"
	"fmt"
)

const (
	defaultWidth  = 10
	defaultHeight = 10

	defaultWinLength = 2

	_Icon = "_"
	xIcon = "X"
	oIcon = "O"

	_Type = -1
	xType = 0
	oType = 1
)

var row = []int{-1, -1, -1, 0, 1, 1, 1, 0}
var col = []int{-1, 0, 1, 1, 1, 0, -1, -1}

type Game struct {
	GameId string
	Board  Board
	Status int

	XFirst int
	WhoAmI int
}

type Move struct {
	X, Y, Turn int
}

type Board struct {
	Width  int
	Height int

	Cells [][]Cell
}

func (g Game) String() string {
	return fmt.Sprintf("Status (%d) XFirst (%d) WhoAmI (%d)", g.Status, g.XFirst, g.WhoAmI)
}

type Cell struct {
	Icon   string
	IsFill bool
}

func InitGame(gameId string) Game {
	return Game{
		GameId: gameId,
		Board:  initBoard(defaultWidth, defaultHeight),
		Status: -1,
		XFirst: 0,
	}
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

func initCell() Cell {
	cell := Cell{}
	cell.Icon = _Icon
	cell.IsFill = false
	return cell
}

func (g Game) TakeMove(move Move) (Game, error) {

	var board = g.Board

	x := move.X
	y := move.Y
	t := move.Turn

	// check valid x, y, t
	if x < 0 || x >= g.Board.Width {
		return g, errors.New("Invalid X")
	}
	if y < 0 || y >= g.Board.Height {
		return g, errors.New("Invalid Y")
	}
	if board.Cells[x][y].IsFill {
		return g, errors.New("Cell is already filled")
	}

	switch t {
	case xType:
		board.Cells[x][y] = Cell{Icon: xIcon}
	case oType:
		board.Cells[x][y] = Cell{Icon: oIcon}
	default:
		return g, errors.New("Invalid type: " + string(t))
	}
	board.Cells[x][y].IsFill = true

	g.Board = board
	g.XFirst = 1 - g.XFirst
	g.Status = board.getStatus()

	return g, nil
}

func (g Game) Render() {

	fmt.Printf("Game: %s\n\n", g.GameId)

	var b = g.Board

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {
			fmt.Print(b.Cells[i][j].Icon, " ")
		}
		fmt.Println()
	}
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

func (b Board) getWinner(x, y int) int {

	if b.Cells[x][y].IsFill == false {
		return -1
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

			if b.Cells[u][v].IsFill == false {
				isWinner = false
				break
			}

			if b.Cells[u][v].Icon != b.Cells[x][y].Icon {
				isWinner = false
				break
			}
		}

		if isWinner {
			if b.Cells[x][y].Icon == xIcon {
				return xType
			} else {
				return oType
			}
		}
	}

	return -1
}

func (b Board) getStatus() int {

	// -1 playing
	// 0 x win
	// 1 0 win
	// 2 tie

	var anyAvaiableMove = false

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {

			if !b.Cells[i][j].IsFill && !anyAvaiableMove {
				anyAvaiableMove = true
			}

			if b.Cells[i][j].IsFill {

				var winner = b.getWinner(i, j)

				if winner != -1 {
					return winner
				}
			}

		}
	}

	if anyAvaiableMove {
		return -1
	} else {
		return 2
	}
}
