package game

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type MoveNotation string

const MOVE_REGEX string = "^[a-i]([1-9]|10)[a-i]([1-9]|10)$"

type Move struct {
	Move      MoveNotation `firestore:"move"`
	Timestamp time.Time    `firestore:"timestamp"`
}

func (m MoveNotation) String() string {
	return string(m)
}

func ParseMoveNotation(smove string) (MoveNotation, error) {
	ok, err := regexp.Match(MOVE_REGEX, bytes.NewBufferString(smove).Bytes())
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("not a valid move")
	}

	return MoveNotation(smove), nil
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
	BULLET TimeLimit = 1
	BLITZ  TimeLimit = 3
	BLITZ2 TimeLimit = 5
	RAPID  TimeLimit = 10
	RAPID2 TimeLimit = 15
	RAPID3 TimeLimit = 20
	RAPID4 TimeLimit = 30
)

type GameStatus string

const (
	INGAME GameStatus = "ingame"
	WIN    GameStatus = "win"
	LOSS   GameStatus = "loss"
	DRAW   GameStatus = "draw"
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
