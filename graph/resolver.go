package graph

import (
	"context"
	"errors"

	"github.com/pafkiuq/backend/graph/model"
	"github.com/pafkiuq/backend/middleware"
	"github.com/pafkiuq/backend/pkg/format"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func GetAuthUserID(ctx context.Context) (format.UserID, bool) {
	userID, ok := ctx.Value(middleware.AUTH_USER_CONTEXT_KEY).(format.UserID)
	return userID, ok
}

type UsersService interface {
	CreateUser(context.Context, format.UserID) (*model.User, error)
	GetUser(context.Context, format.UserID) (*model.User, error)
	EditUser(context.Context, format.UserID, model.UserEditInput) (*model.User, error)
	DeleteUser(context.Context, format.UserID) error
}

type Config struct {
	UsersService
}

type Resolver struct {
	UsersService
}

func NewResolver(cfg Config) (*Resolver, error) {
	if cfg.UsersService == nil {
		return nil, errors.New("users service required")
	}

	return &Resolver{
		UsersService: cfg.UsersService,
	}, nil
}
