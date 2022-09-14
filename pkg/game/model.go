package game

import (
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type MoveNotation string

const MOVE_REGEX string = "[a-i][(0-9)|10][a-i][(0-9)|10]"

type Move struct {
	Move      MoveNotation `firestore:"move"`
	Timestamp time.Time    `firestore:"timestamp"`
}

func (m MoveNotation) String() string {
	return string(m)
}

type GameType string

const (
	JANGGI GameType = "janggi"
	SHOGI  GameType = "shogi"
)

func (g GameType) String() string {
	return string(g)
}

type TimeLimit int

const (
	BULLET int = 1
	BLITZ  int = 3
	BLITZ2 int = 5
	RAPID  int = 10
	RAPID2 int = 15
	RAPID3 int = 20
	RAPID4 int = 30
)

// NOTE: In janggi, the game always starts with red
type GameDocument struct {
	ID        format.GameID `firestore:"id"`
	WinnerID  format.UserID `firestore:"winner_id"`
	PlayerOne format.UserID `firestore:"player_one"`
	PlayerTwo format.UserID `firestore:"player_two"`
	Moves     []Move        `firestore:"moves"`
	Draw      bool          `firestore:"draw"`
	Aborted   bool          `firestore:"aborted"`
	Type      GameType      `firestore:"type"`
	TimeLimit TimeLimit     `firestore:"time_limit"`
	Timestamp time.Time     `firestore:"timestamp"`
}
