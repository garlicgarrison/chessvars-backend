package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	game_pb "github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Game struct {
	services *Services
	gameID   format.GameID

	getter[*game_pb.Game, func(context.Context) (*game_pb.Game, error)]
}

type TimeLimit string

const (
	BULLET TimeLimit = "BULLET"
	BLITZ  TimeLimit = "BLITZ"
	BLITZ2 TimeLimit = "BLITZ2"
	RAPID  TimeLimit = "RAPID"
	RAPID2 TimeLimit = "RAPID2"
	RAPID3 TimeLimit = "RAPID3"
	RAPID4 TimeLimit = "RAPID4"
	NONE   TimeLimit = "NONE"
)

func NewGame(services *Services, gameID format.GameID) *Game {
	return &Game{
		services: services,
		gameID:   gameID,
		getter: NewGetter(func(ctx context.Context) (*game_pb.Game, error) {
			game, err := services.Game.GetGame(ctx, game_pb.GetGameRequest{
				GameID: gameID,
			})

			if err != nil {
				return nil, err
			}

			return game, nil
		}),
	}
}

func NewGameWithData(services *Services, data *game_pb.Game) *Game {
	return &Game{
		services: services,
		gameID:   data.ID,
		getter: NewGetter(func(ctx context.Context) (*game_pb.Game, error) {
			return data, nil
		}),
	}
}

func (g *Game) ID(ctx context.Context) (string, error) {
	return g.gameID.String(), nil
}

func (g *Game) Moves(ctx context.Context) ([]*Move, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
	}

	toRet := make([]*Move, 0)
	for _, move := range game.Moves {
		toRet = append(toRet, NewMove(g.services, &move))
	}

	return toRet, nil
}

func (g *Game) PlayerOne(ctx context.Context) (*User, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
	}

	if game.PlayerOne == "" {
		return nil, nil
	}

	return NewUser(g.services, game.PlayerOne), nil
}

func (g *Game) PlayerTwo(ctx context.Context) (*User, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
	}

	if game.PlayerTwo == "" {
		return nil, nil
	}

	return NewUser(g.services, game.PlayerTwo), nil
}

func (g *Game) Winner(ctx context.Context) (*User, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
	}

	if game.WinnerID == "" {
		return nil, nil
	}

	return NewUser(g.services, game.WinnerID), nil
}

func (g *Game) Draw(ctx context.Context) (bool, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return false, err
	}

	return game.Draw, nil
}

func (g *Game) Aborted(ctx context.Context) (bool, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return false, err
	}

	return game.Aborted, nil
}

func (g *Game) Type(ctx context.Context) (string, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return "", err
	}

	return game.Type.String(), nil
}

func (g *Game) TimeLimit(ctx context.Context) (TimeLimit, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return NONE, err
	}

	var timeLimit TimeLimit
	switch game.TimeLimit {
	case game_pb.BULLET:
		timeLimit = BULLET
	case game_pb.BLITZ:
		timeLimit = BLITZ
	case game_pb.BLITZ2:
		timeLimit = BLITZ2
	case game_pb.RAPID:
		timeLimit = RAPID
	case game_pb.RAPID2:
		timeLimit = RAPID2
	case game_pb.RAPID3:
		timeLimit = RAPID3
	case game_pb.RAPID4:
		timeLimit = RAPID4
	default:
		timeLimit = NONE
	}

	return timeLimit, nil
}

func (g *Game) Timestamp(ctx context.Context) (string, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return "", err
	}

	return game.Timestamp.String(), nil
}
