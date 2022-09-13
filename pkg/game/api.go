package game

import (
	"context"
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type Service interface {
	CreateGame(context.Context, CreateGameRequest) (*CreateGameResponse, error)
	GetGame(context.Context, GetGameRequest) (*GetGameResponse, error)
	EditGame(context.Context, EditGameRequest) (*EditGameResponse, error)
}

type Game struct {
	ID        format.GameID `json:"game_id"`
	WinnerID  format.UserID `json:"winner_id"`
	PlayerOne format.UserID `json:"player_one"`
	PlayerTwo format.UserID `json:"player_two"`
	Moves     []Move        `json:"moves"`
	Draw      bool          `json:"draw"`
	Aborted   bool          `json:"aborted"`
	Type      GameType      `json:"type"`
	Timestamp time.Time     `json:"timestamp"`
}

type GetGameRequest struct {
	GameID format.GameID `json:"game_id"`
}

type GetGameResponse = Game

type CreateGameRequest struct {
	UserID format.UserID `json:"user_id"`
	Type   GameType      `json:"type"`
}

type CreateGameResponse = Game

type EditGameRequest struct {
	UserID format.UserID `json:"user_id"`
	Move   MoveNotation  `json:"move"`
}

type EditGameResponse = Game
