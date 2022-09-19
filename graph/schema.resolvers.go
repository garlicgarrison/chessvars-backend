package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/garlicgarrison/chessvars-backend/graph/generated"
	"github.com/garlicgarrison/chessvars-backend/graph/model"
	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"google.golang.org/grpc/codes"
)

// TimeLimit is the resolver for the timeLimit field.
func (r *gameResolver) TimeLimit(ctx context.Context, obj *resolver.Game) (*model.TimeLimit, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserEdit is the resolver for the userEdit field.
func (r *mutationResolver) UserEdit(ctx context.Context, input model.UserEditInput) (*model.UserMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserDelete is the resolver for the userDelete field.
func (r *mutationResolver) UserDelete(ctx context.Context) (*model.BasicMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// GameCreate is the resolver for the gameCreate field.
func (r *mutationResolver) GameCreate(ctx context.Context, typeArg model.GameType, limit model.TimeLimit) (*model.GameMutationResponse, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not validate user")
	}

	var timeLimit game.TimeLimit
	switch limit {
	case model.TimeLimitBullet:
		timeLimit = game.BULLET
	case model.TimeLimitBlitz:
		timeLimit = game.BLITZ
	case model.TimeLimitBlitz2:
		timeLimit = game.BLITZ2
	case model.TimeLimitRapid:
		timeLimit = game.RAPID
	case model.TimeLimitRapid2:
		timeLimit = game.RAPID2
	case model.TimeLimitRapid3:
		timeLimit = game.RAPID3
	case model.TimeLimitRapid4:
		timeLimit = game.RAPID4
	default:
		return nil, fmt.Errorf("time limit not valid")
	}

	var gameType game.GameType
	switch typeArg {
	case model.GameTypeJanggi:
		gameType = game.JANGGI
	default:
		return nil, fmt.Errorf("game not implemented")
	}

	game, err := r.Services.Game.CreateGame(ctx, game.CreateGameRequest{
		UserID:    userID,
		TimeLimit: timeLimit,
		Type:      gameType,
	})
	if err != nil {
		return nil, err
	}

	return &model.GameMutationResponse{
		Code:    int(codes.OK),
		Success: true,
		Message: "game was successfully created",
		Game:    resolver.NewGameWithData(r.Services, game),
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *string) (*resolver.User, error) {
	if id != nil {
		userID, err := format.ParseUserID(*id)
		if err != nil {
			return nil, err
		}

		return resolver.NewUser(r.Services, userID), nil
	}

	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, errors.New("no auth user in context")
	}

	return resolver.NewUser(r.Services, userID), nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type gameResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
