package game

import (
	"errors"
	"fmt"
)

const (
	defaultWidth  = 10
	defaultHeight = 10

	defaultWinLength = 5

	_Icon = "_"
	xIcon = "X"
	oIcon = "O"

	_Type = -1
	xType = 0
	oType = 1
)

type Game struct {
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
	return fmt.Sprint("Status: ", g.Status, " XFirst: ", g.XFirst, " WhoAmI: ", g.WhoAmI)
}

type Cell struct {
	Icon   string
	IsFill bool
}

func InitGame() Game {
	return Game{
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

	var b = g.Board

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
	if b.Cells[x][y].IsFill {
		return g, errors.New("Cell is already filled")
	}

	switch t {
	case xType:
		b.Cells[x][y] = Cell{Icon: xIcon}
	case oType:
		b.Cells[x][y] = Cell{Icon: oIcon}
	default:
		return g, errors.New("Invalid type: " + string(t))
	}
	b.Cells[x][y].IsFill = true

	g.Board = b
	g.XFirst = 1 - g.XFirst

	return g, nil
}

func (g Game) Render() {

	fmt.Printf("Game: \n\n")

	var b = g.Board

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {
			fmt.Print(b.Cells[i][j].Icon, " ")
		}
		fmt.Println()
	}
}

func (b Board) CheckWinner() int {

	// -1 playing
	// 0 x win
	// 1 0 win
	// 2 tie

	return -1
}
