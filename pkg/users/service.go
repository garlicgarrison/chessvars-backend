package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Firestore firestore.Firestore
}

type service struct {
	fs firestore.Firestore
}

const DEFAULT_ELO int = 1200

func NewService(cfg Config) (Service, error) {
	if cfg.Firestore == nil {
		return nil, errors.New("firestore required")
	}

	return &service{
		fs: cfg.Firestore,
	}, nil
}

func populateUser(user *UserDocument) *GetUserResponse {
	return &GetUserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Elo:      user.Elo,
	}
}

func (s *service) CreateUser(ctx context.Context, request CreateUserRequest) (*GetUserResponse, error) {
	user := UserDocument{
		UserID: request.UserID,
		Elo:    DEFAULT_ELO,
	}

	_, err := s.getUserRef(request.UserID).Create(ctx, user)
	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return nil, err
		}
		return s.GetUser(ctx, GetUserRequest{
			UserID: request.UserID,
		})
	}

	return populateUser(&user), nil
}

func (s *service) GetUser(ctx context.Context, request GetUserRequest) (*GetUserResponse, error) {
	userSnap, err := s.getUserRef(request.UserID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user UserDocument
	err = userSnap.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return populateUser(&user), nil
}

func (s *service) EditUser(ctx context.Context, request EditUserRequest) (*EditUserResponse, error) {
	verifyUsername := func(username string) (string, error) {
		if username == "" {
			return "", errors.New("invalid empty username")
		}

		// lowercase
		username = strings.ToLower(username)

		const allowed = "abcdefghijklmnopqrstuvwxyz0123456789&_"
		for _, c := range username {
			if !strings.ContainsRune(allowed, c) {
				return "", fmt.Errorf("invalid username -- username: %s; invalid: %v", username, c)
			}
		}

		return username, nil
	}

	var user UserDocument
	err := s.fs.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		userSnap, err := t.Get(s.getUserRef(request.UserID))
		if err != nil {
			return err
		}

		err = userSnap.DataTo(&user)
		if err != nil {
			return err
		}

		if request.Username != nil {
			username, err := verifyUsername(*request.Username)
			if err != nil {
				return err
			}

			user.Username = username
		}

		return t.Set(s.getUserRef(request.UserID), user)
	})
	if err != nil {
		return nil, err
	}

	return populateUser(&user), nil
}

func (s *service) DeleteUser(ctx context.Context, userID format.UserID) error {
	return nil
}
