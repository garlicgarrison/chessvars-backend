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

func (m *Move) Move(ctx context.Context) (string, error) {
	move, err := m.getter.Call(ctx)
	if err != nil {
		return "", nil
	}

	return move.Move.String(), nil
}

func (m *Move) Timestamp(ctx context.Context) (string, error) {
	move, err := m.getter.Call(ctx)
	if err != nil {
		return "", nil
	}

	return move.Timestamp.String(), nil
}
