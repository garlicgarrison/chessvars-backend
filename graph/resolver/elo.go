package resolver

import (
	"context"
	"fmt"

	"github.com/garlicgarrison/chessvars-backend/pkg/elo"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
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
			return services.Elo.GetElos(ctx, elo.GetElosRequest{
				UserID: userID,
			})
		}),
	}
}

func (e *Elo) Janggi(ctx context.Context) (int, error) {
	reply, err := e.getter.Call(ctx)
	fmt.Printf("[Janggi] resolver reply -- %v", reply)
	if err != nil {
		return 1200, err
	}

	for _, el := range reply.Elos {
		if el.Game == elo.JANGGI {
			return el.Elo, nil
		}
	}

	return 1200, err
}

func (e *Elo) Shogi(ctx context.Context) (int, error) {
	reply, err := e.getter.Call(ctx)
	if err != nil {
		return 1200, err
	}

	for _, el := range reply.Elos {
		if el.Game == elo.SHOGI {
			return el.Elo, nil
		}
	}

	return 1200, err
}
