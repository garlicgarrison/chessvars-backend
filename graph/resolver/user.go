package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	services *Services
	userID   format.UserID

	getter[*users.User, func(context.Context) (*users.User, error)]
}

func NewUser(services *Services, userID format.UserID) *User {
	return &User{
		services: services,
		userID:   userID,
		getter: NewGetter(func(ctx context.Context) (*users.User, error) {
			user, err := services.Users.GetUser(ctx, users.GetUserRequest{
				UserID: userID,
			})
			if err != nil {
				if status.Code(err) != codes.NotFound {
					return nil, err
				}

				authUserID, ok := GetAuthUserID(ctx)
				if !ok || authUserID != userID {
					return nil, err
				}

				return services.Users.CreateUser(ctx, users.CreateUserRequest{
					UserID: userID,
				})
			}

			return user, nil
		}),
	}
}

func NewUserWithData(services *Services, data *users.User) *User {
	return &User{
		services: services,
		userID:   data.UserID,
		getter: NewGetter(func(ctx context.Context) (*users.User, error) {
			return data, nil
		}),
	}
}

func (u *User) ID(ctx context.Context) (string, error) {
	return u.userID.String(), nil
}

func (u *User) Exists(ctx context.Context) (bool, error) {
	_, err := u.getter.Call(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (u *User) Username(ctx context.Context) (string, error) {
	reply, err := u.getter.Call(ctx)
	if err != nil {
		return "", err
	}

	return reply.Username, nil
}

func (u *User) Elo(ctx context.Context) (int, error) {
	reply, err := u.getter.Call(ctx)
	if err != nil {
		return 1200, err
	}

	return reply.Elo, nil
}
