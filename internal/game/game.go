package game

import (
	"fmt"
)

const (
	defaultWidth  = 20
	defaultHeight = 20

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

func InitGame(gameId string) Game {
	var game = Game{
		GameID: gameId,
		Board:  initBoard(defaultWidth, defaultHeight),
		Status: Running,
		Player: XPlayer,
	}

	return game
}

func (g Game) String() string {
	return fmt.Sprintf("Status (%d) Player (%d)", g.Status, g.Player)
}

func (g Game) Copy() Game {
	var newGame = Game{
		GameID: g.GameID,
		Status: g.Status,
		Player: g.Player,
		Board:  initBoard(g.Board.Width, g.Board.Height),
	}
	for i := range newGame.Board.Cells {
		copy(newGame.Board.Cells[i], g.Board.Cells[i])
	}
	return newGame
}

func (g Game) IsValidMove(move Move) error {
	x := move.X
	y := move.Y

	if x < 0 || x >= g.Board.Width {
		return fmt.Errorf("Invalid X")
	}
	if y < 0 || y >= g.Board.Height {
		return fmt.Errorf("Invalid Y")
	}
	if g.Board.Cells[x][y].isFill() {
		return fmt.Errorf("Cell is already filled")
	}
	return nil
}

func (g Game) TakeMove(move Move) (Game, error) {

	var board = g.Board

	x := move.X
	y := move.Y
	p := move.Player

	// check valid x, y, t
	if err := g.IsValidMove(move); err != nil {
		return g, err
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
