package game

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type Config struct {
	Firestore firestore.Firestore
}

type service struct {
	fs firestore.Firestore
}

func NewService(cfg Config) (Service, error) {
	if cfg.Firestore == nil {
		return nil, errors.New("firestore required")
	}

	return &service{
		fs: cfg.Firestore,
	}, nil
}

func (s *service) populateGame(game *GameDocument) *Game {
	moves := make([]MoveResponse, 0)
	for _, m := range game.Moves {
		moves = append(moves, MoveResponse{
			Move:      m.Move,
			Timestamp: m.Timestamp,
		})
	}

	return &Game{
		ID:        game.ID,
		WinnerID:  game.WinnerID,
		PlayerOne: game.PlayerOne,
		PlayerTwo: game.PlayerTwo,
		Moves:     moves,
		Draw:      game.Draw,
		Aborted:   game.Aborted,
		TimeLimit: game.TimeLimit,
		Type:      game.Type,
		Timestamp: game.Timestamp,
	}
}

func (s *service) CreateGame(ctx context.Context, request CreateGameRequest) (*CreateGameResponse, error) {
	gameID := format.NewGameID()
	now := time.Now()

	gameDoc := GameDocument{
		ID:        gameID,
		TimeLimit: request.TimeLimit,
		Timestamp: now,
	}

	switch request.Type {
	case JANGGI:
		gameDoc.Type = JANGGI
	case SHOGI:
		gameDoc.Type = SHOGI
	}

	// Decides if user if player 1 or player 2 randomly
	if rand.Intn(2) == 0 {
		gameDoc.PlayerOne = request.UserID
	} else {
		gameDoc.PlayerTwo = request.UserID
	}

	_, err := s.getGameRef(JANGGI, gameID).Create(ctx, gameDoc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) GetGame(ctx context.Context, request GetGameRequest) (*GetGameResponse, error) {
	return nil, nil
}

func (s *service) EditGame(ctx context.Context, request EditGameRequest) (*EditGameResponse, error) {
	return nil, nil
}
