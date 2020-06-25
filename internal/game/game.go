package game

import (
	"fmt"
)

const (
	defaultWidth  = 10
	defaultHeight = 10

	defaultWinLength = 2
)

var row = []int{-1, -1, -1, 0, 1, 1, 1, 0}
var col = []int{-1, 0, 1, 1, 1, 0, -1, -1}

type Game struct {
	GameID string
	Board  Board

	Status Status
	Player Player
}

type Move struct {
	X, Y   int
	Player Player
}

func (g Game) String() string {
	return fmt.Sprintf("Status (%d) Player (%d)", g.Status, g.Player)
}

func InitGame(gameId string) Game {
	var game = Game{
		GameID: gameId,
		Board:  initBoard(defaultWidth, defaultHeight),
		Status: Running,
		Player: XPlayer,
	}

	return game
}

func (g Game) TakeMove(move Move) (Game, error) {

	var board = g.Board

	x := move.X
	y := move.Y
	p := move.Player

	// check valid x, y, t
	if x < 0 || x >= g.Board.Width {
		return g, fmt.Errorf("Invalid X")
	}
	if y < 0 || y >= g.Board.Height {
		return g, fmt.Errorf("Invalid Y")
	}
	if board.Cells[x][y].isFill() {
		return g, fmt.Errorf("Cell is already filled")
	}

	board.Cells[x][y] = Cell{Player: p}

	g.Board = board
	g.Player = g.Player.changeTurn()
	g.Status = board.getStatus()

	return g, nil
}

func (g Game) Render() {

	fmt.Printf("Game: %s\n\n", g.GameID)

	var b = g.Board

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {
			fmt.Printf("%s ", b.Cells[i][j].Player)
		}
		fmt.Println()
	}
}
