package game

import (
	"time"

	"github.com/pafkiuq/backend/pkg/format"
)

type Move string

const MOVE_REGEX string = "[a-i][(0-9)|10][a-i][(0-9)|10]"

// NOTE: In janggi, the game always starts with red
type GameDocument struct {
	ID        format.GameID `firestore:"id"`
	WinnerID  format.UserID `firestore:"winner_id"`
	PlayerOne format.UserID `firestore:"player_one"`
	PlayerTwo format.UserID `firestore:"player_two"`
	Moves     []Move        `firestore:"moves"`
	Draw      bool          `firestore:"draw"`
	Aborted   bool          `firestore:"aborted"`
	Timestamp time.Time     `firestore:"timestamp"`
}
