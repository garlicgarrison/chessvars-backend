package graph

import (
	"context"
	"errors"

	"github.com/garlicgarrison/chessvars-backend/middleware"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func GetAuthUserID(ctx context.Context) (format.UserID, bool) {
	userID, ok := ctx.Value(middleware.AUTH_USER_CONTEXT_KEY).(format.UserID)
	return userID, ok
}

type Config struct {
	UsersService users.Service
	GameService  game.Service
}

type Resolver struct {
	UsersService users.Service
	GameService  game.Service
}

type Services struct {
	Users users.Service
	Game  game.Service
}

func NewResolver(cfg Config) (*Resolver, error) {
	if cfg.UsersService == nil {
		return nil, errors.New("users service required")
	}

	return &Resolver{
		UsersService: cfg.UsersService,
	}, nil
}
