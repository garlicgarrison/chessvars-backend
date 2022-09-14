package resolver

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
)

type Services struct {
	Users users.Service
	Game  game.Service
}
