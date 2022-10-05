package users

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type Service interface {
	CreateUser(context.Context, CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, GetUserRequest) (*GetUserResponse, error)
	EditUser(context.Context, EditUserRequest) (*EditUserResponse, error)
}

type User struct {
	UserID   format.UserID `json:"user_id"`
	Username string        `json:"username"`
}

type CreateUserRequest struct {
	UserID format.UserID `json:"user_id"`
}

type CreateUserResponse = User

type GetUserRequest struct {
	UserID format.UserID `json:"user_id"`
}

type GetUserResponse = User

type EditUserRequest struct {
	UserID   format.UserID
	Username *string `json:"username"`
}

type EditUserResponse = User
