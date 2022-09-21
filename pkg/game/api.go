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
	JoinGame(context.Context, JoinGameRequest) (*EditGameResponse, error)
}

type MoveResponse struct {
	Move      MoveNotation `json:"move"`
	Timestamp time.Time    `json:"timestamp"`
}

type Game struct {
	ID        format.GameID  `json:"game_id"`
	WinnerID  format.UserID  `json:"winner_id"`
	PlayerOne format.UserID  `json:"player_one"`
	PlayerTwo format.UserID  `json:"player_two"`
	Moves     []MoveResponse `json:"moves"`
	Draw      bool           `json:"draw"`
	Aborted   bool           `json:"aborted"`
	TimeLimit TimeLimit      `json:"time_limit"`
	Type      GameType       `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
}

type GetGameRequest struct {
	GameID format.GameID `json:"game_id"`
}

type GetGameResponse = Game

type CreateGameRequest struct {
	UserID    format.UserID `json:"user_id"`
	TimeLimit TimeLimit     `json:"time_limit"`
	Type      GameType      `json:"type"`
}

type CreateGameResponse = Game

type EditGameRequest struct {
	UserID format.UserID `json:"user_id"`
	GameID format.GameID `json:"game_id"`
	Move   MoveNotation  `json:"move"`

	/* For now, we trust the client on status of the game */
	Status GameStatus `json:"status"`
}

type JoinGameRequest struct {
	GameID format.GameID `json:"game_id"`
	UserID format.UserID `json:"user_id"`
}

type EditGameResponse = Game
