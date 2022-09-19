package elo

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

const DEFAULT_ELO int = 1200

type EloDocument struct {
	UserID   format.UserID `firestore:"user_id"`
	GameType game.GameType `firestore:"game_type"`
	Elo      int           `firestore:"elo"`
}
