package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/game"
)

type Move struct {
	services *Services
	move     *game.MoveResponse

	getter[*game.MoveResponse, func(context.Context) (*game.MoveResponse, error)]
}

func NewMove(services *Services, move *game.MoveResponse) *Move {
	return &Move{
		services: services,
		move:     move,
		getter: NewGetter(func(ctx context.Context) (*game.MoveResponse, error) {
			return move, nil
		}),
	}
}
