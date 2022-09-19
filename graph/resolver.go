package graph

import (
	"sync"

	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Config struct {
	*resolver.Services
}

type Resolver struct {
	// backend services
	*resolver.Services

	/*
		Key: gameID
		Value: *resolver.Move
	*/
	moveChannels map[format.GameID]chan *resolver.Move

	mutex sync.Mutex
}

func NewResolver(cfg Config) (*Resolver, error) {
	return &Resolver{Services: cfg.Services}, nil
}
