package elo

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Service interface {
	CreateElo(context.Context, CreateEloRequest) (*CreateEloResponse, error)
	GetGame(context.Context, GetEloRequest) (*GetEloResponse, error)
	UpdateElo(context.Context, UpdateEloRequest) (*UpdateEloResponse, error)
}

type GameStatus string

const (
	WIN  GameStatus = "win"
	LOSS GameStatus = "loss"
	DRAW GameStatus = "draw"
)

type CreateEloRequest struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
}

type GetEloRequest struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
}

type UpdateEloRequest struct {
	UserID      format.UserID `json:"user_id"`
	OtherUserID format.UserID `json:"other_user_id"`
	Game        game.GameType `json:"game"`
	Status      GameStatus    `json:"status"`
}

type Elo struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
	Elo    int           `json:"elo"`
}

type CreateEloResponse = Elo
type UpdateEloResponse = Elo
type GetEloResponse = Elo
