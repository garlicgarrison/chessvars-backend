package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/pafkiuq/backend/graph/generated"
	"github.com/pafkiuq/backend/graph/model"
)

// UserEdit is the resolver for the userEdit field.
func (r *mutationResolver) UserEdit(ctx context.Context, input model.UserEditInput) (*model.UserMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserDelete is the resolver for the userDelete field.
func (r *mutationResolver) UserDelete(ctx context.Context) (*model.BasicMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// LikeVideo is the resolver for the likeVideo field.
func (r *mutationResolver) LikeVideo(ctx context.Context, id string, like bool) (*model.BasicMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Feed is the resolver for the feed field.
func (r *queryResolver) Feed(ctx context.Context) (*model.Videos, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
