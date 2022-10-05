package game

import (
	"context"
	"errors"
	"fmt"
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
		moves = append(moves, MoveResponse(m))
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

func (s *service) validateMove(userID format.UserID, game *GameDocument) bool {
	if game.PlayerOne == userID && len(game.Moves)%2 != 1 ||
		game.PlayerTwo == userID && len(game.Moves)%2 != 0 ||
		game.Aborted ||
		game.Draw ||
		game.WinnerID != "" {
		return false
	}
	return true
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

	_, err := s.getGameRef(gameID).Create(ctx, gameDoc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) GetGame(ctx context.Context, request GetGameRequest) (*GetGameResponse, error) {
	gameSnap, err := s.getGameRef(request.GameID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var game GameDocument
	err = gameSnap.DataTo(&game)
	if err != nil {
		return nil, err
	}

	return s.populateGame(&game), nil
}

func (s *service) EditGame(ctx context.Context, request EditGameRequest) (*EditGameResponse, error) {
	now := time.Now()

	var game GameDocument
	err := s.fs.RunTransaction(ctx, func(_ context.Context, t *firestore.Transaction) error {
		gameSnap, err := t.Get(s.getGameRef(request.GameID))
		if err != nil {
			return err
		}

		err = gameSnap.DataTo(&game)
		if err != nil {
			return err
		}

		/*
			This makes sure that a move is even allowed to be made.
			Later, this should have a move validator to make sure a move is legal,
			and to see if a move is a winning or drawing move with game logic.
			For now, we can trust the client for game logic.
		*/
		if !s.validateMove(request.UserID, &game) {
			return fmt.Errorf("move not validated")
		}

		isPlayer1 := game.PlayerOne == request.UserID
		switch request.Status {
		case LOSS:
			if isPlayer1 {
				game.WinnerID = game.PlayerTwo
			} else {
				game.WinnerID = game.PlayerOne
			}
		case WIN:
			if isPlayer1 {
				game.WinnerID = game.PlayerOne
			} else {
				game.WinnerID = game.PlayerTwo
			}
		case DRAW:
			game.Draw = true
		case INGAME:
			break
		}

		newMoves := append(game.Moves, Move{
			Move:      request.Move,
			Timestamp: now,
		})
		game.Moves = newMoves

		return t.Set(
			s.getGameRef(request.GameID),
			game,
		)
	})
	if err != nil {
		return nil, err
	}

	return s.populateGame(&game), nil
}

func (s *service) JoinGame(ctx context.Context, request JoinGameRequest) (*EditGameResponse, error) {
	var game GameDocument
	err := s.fs.RunTransaction(ctx, func(_ context.Context, t *firestore.Transaction) error {
		gameSnap, err := t.Get(s.getGameRef(request.GameID))
		if err != nil {
			return err
		}

		err = gameSnap.DataTo(&game)
		if err != nil {
			return err
		}

		if game.Aborted {
			return fmt.Errorf("game was aborted")
		}

		if game.PlayerOne == "" {
			game.PlayerOne = request.UserID
		} else if game.PlayerTwo == "" {
			game.PlayerTwo = request.UserID
		} else {
			return fmt.Errorf("this game is no longer available")
		}

		return t.Set(
			s.getGameRef(request.GameID),
			game,
		)
	})
	if err != nil {
		return nil, err
	}

	return s.populateGame(&game), nil
}
