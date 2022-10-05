package elo

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Service interface {
	CreateElo(context.Context, CreateEloRequest) (*CreateEloResponse, error)
	GetElo(context.Context, GetEloRequest) (*GetEloResponse, error)
	GetElos(context.Context, GetElosRequest) (*Elos, error)
	UpdateElo(context.Context, UpdateEloRequest) (*UpdateEloResponse, error)
}

type CreateEloRequest struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
}

type GetEloRequest struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
}

type GetElosRequest struct {
	UserID format.UserID `json:"user_id"`
}

type Elos struct {
	Elos []*Elo `json:"elos"`
}

type UpdateEloRequest struct {
	UserID      format.UserID   `json:"user_id"`
	OtherUserID format.UserID   `json:"other_user_id"`
	Game        game.GameType   `json:"game"`
	Status      game.GameStatus `json:"status"`
}

type Elo struct {
	UserID format.UserID `json:"user_id"`
	Game   game.GameType `json:"game"`
	Elo    int           `json:"elo"`
}

type CreateEloResponse = Elo
type UpdateEloResponse = Elo
type GetEloResponse = Elo
