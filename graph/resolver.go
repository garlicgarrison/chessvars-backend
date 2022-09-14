package graph

import (
	"errors"

	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Config struct {
	UsersService users.Service
	GameService  game.Service
}

type Resolver struct {
	UsersService users.Service
	GameService  game.Service
}

func NewResolver(cfg Config) (*Resolver, error) {
	if cfg.UsersService == nil {
		return nil, errors.New("users service required")
	}

	return &Resolver{
		UsersService: cfg.UsersService,
	}, nil
}
