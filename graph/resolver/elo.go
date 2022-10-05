package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/elo"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Elo struct {
	services *Services
	userID   format.UserID

	getter[*elo.Elos, func(context.Context) (*elo.Elos, error)]
}

func NewElo(services *Services, userID format.UserID) *Elo {
	return &Elo{
		services: services,
		userID:   userID,
		getter: NewGetter(func(ctx context.Context) (*elo.Elos, error) {
			reply, err := services.Elo.GetElos(ctx, elo.GetElosRequest{
				UserID: userID,
			})
			if err != nil {
				return nil, err
			}

			return reply, nil
		}),
	}
}

func (e *Elo) Janggi(ctx context.Context) (int, error) {
	reply, err := e.getter.Call(ctx)
	if err != nil {
		return 1200, err
	}

	for _, elo := range reply.Elos {
		if elo.Game == game.JANGGI {
			return elo.Elo, nil
		}
	}

	return 1200, err
}

func (e *Elo) Shogi(ctx context.Context) (int, error) {
	reply, err := e.getter.Call(ctx)
	if err != nil {
		return 1200, err
	}

	for _, elo := range reply.Elos {
		if elo.Game == game.SHOGI {
			return elo.Elo, nil
		}
	}

	return 1200, err
}