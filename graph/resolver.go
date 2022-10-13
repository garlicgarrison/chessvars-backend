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
		Value: []chan *resolver.Move
		NOTE: whenever a game changes moves,
		everyone that subscribes should get a push
	*/
	// moveChannels map[format.GameID](map[format.UserID]chan *resolver.Move)

	GamesMovesMap sync.Map
}

type Observers struct {
	MoveObservers sync.Map
}

type MoveObserver struct {
	UserID format.UserID
	Move   chan *resolver.Move
}

func NewResolver(cfg Config) (*Resolver, error) {
	return &Resolver{
		Services:      cfg.Services,
		GamesMovesMap: sync.Map{},
	}, nil
}

func (r *Resolver) getObserverMap(gameID format.GameID) *Observers {
	game, _ := r.GamesMovesMap.LoadOrStore(gameID, &Observers{
		MoveObservers: sync.Map{},
	})

	return game.(*Observers)
}
