package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/garlicgarrison/chessvars-backend/graph/generated"
	"github.com/garlicgarrison/chessvars-backend/graph/model"
	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
)

// TimeLimit is the resolver for the timeLimit field.
func (r *gameResolver) TimeLimit(ctx context.Context, obj *resolver.Game) (*model.TimeLimit, error) {
	panic(fmt.Errorf("not implemented"))
}

// Move is the resolver for the move field.
func (r *moveResolver) Move(ctx context.Context, obj *resolver.Move) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Timestamp is the resolver for the timestamp field.
func (r *moveResolver) Timestamp(ctx context.Context, obj *resolver.Move) (*string, error) {
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
func (r *mutationResolver) GameCreate(ctx context.Context, typeArg model.GameType) (*model.GameMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *string) (*resolver.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Move returns generated.MoveResolver implementation.
func (r *Resolver) Move() generated.MoveResolver { return &moveResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type gameResolver struct{ *Resolver }
type moveResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
