package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Game struct {
	services *Services
	gameID   format.GameID

	getter[*game.Game, func(context.Context) (*game.Game, error)]
}

func NewGame(services *Services, gameID format.GameID) *Game {
	return &Game{
		services: services,
		gameID:   gameID,
		getter: NewGetter(func(ctx context.Context) (*game.Game, error) {
			game, err := services.Game.GetGame(ctx, game.GetGameRequest{
				GameID: gameID,
			})

			if err != nil {
				return nil, err
			}

			return game, nil
		}),
	}
}

func NewGameWithData(services *Services, data *game.Game) *Game {
	return &Game{
		services: services,
		gameID:   data.ID,
		getter: NewGetter(func(ctx context.Context) (*game.Game, error) {
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

	return NewUser(g.services, game.PlayerOne), nil
}

func (g *Game) PlayerTwo(ctx context.Context) (*User, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
	}

	return NewUser(g.services, game.PlayerTwo), nil
}

func (g *Game) Winner(ctx context.Context) (*User, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return nil, err
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

func (g *Game) TimeLimit(ctx context.Context) (int, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return 0, err
	}

	return int(game.TimeLimit), nil
}

func (g *Game) Timestamp(ctx context.Context) (string, error) {
	game, err := g.getter.Call(ctx)
	if err != nil {
		return "", err
	}

	return game.Timestamp.String(), nil
}
