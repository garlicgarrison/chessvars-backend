package elo

import (
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type GameType string

const (
	JANGGI GameType = "janggi"
	SHOGI  GameType = "shogi"
)

func (g GameType) String() string {
	return string(g)
}

type GameStatus string

const (
	INGAME GameStatus = "ingame"
	WIN    GameStatus = "win"
	LOSS   GameStatus = "loss"
	DRAW   GameStatus = "draw"
)

const (
	DEFAULT_ELO          int    = 1200
	CURRENT_ELO_DOCUMENT string = "current"
)

type EloDocument struct {
	UserID    format.UserID `firestore:"user_id"`
	GameType  GameType      `firestore:"game_type"`
	Elo       int           `firestore:"elo"`
	Timestamp time.Time     `firestore:"timestamp"`
}
