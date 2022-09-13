package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/garlicgarrison/chessvars-backend/graph/model"
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Firestore firestore.Firestore
}

type Service struct {
	fs firestore.Firestore
}

const DEFAULT_ELO int = 1200

func NewService(cfg Config) (*Service, error) {
	if cfg.Firestore == nil {
		return nil, errors.New("firestore required")
	}

	return &Service{
		fs: cfg.Firestore,
	}, nil
}

func populateUser(user *UserDocument) *model.User {
	exists := user != nil
	return &model.User{
		ID:       user.UserID.String(),
		Exists:   &exists,
		Username: &user.Username,
	}
}

func (s *Service) CreateUser(ctx context.Context, userID format.UserID) (*model.User, error) {
	user := UserDocument{
		UserID: userID,
		Elo:    DEFAULT_ELO,
	}

	_, err := s.getUserRef(userID).Create(ctx, user)
	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return nil, err
		}
		return s.GetUser(ctx, userID)
	}

	return populateUser(&user), nil
}

func (s *Service) GetUser(ctx context.Context, userID format.UserID) (*model.User, error) {
	userSnap, err := s.getUserRef(userID).Get(ctx)
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

func (s *Service) EditUser(ctx context.Context, userID format.UserID, input model.UserEditInput) (*model.User, error) {
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
		userSnap, err := t.Get(s.getUserRef(userID))
		if err != nil {
			return err
		}

		err = userSnap.DataTo(&user)
		if err != nil {
			return err
		}

		if input.Username != nil {
			username, err := verifyUsername(*input.Username)
			if err != nil {
				return err
			}

			user.Username = username
		}

		return t.Set(s.getUserRef(userID), user)
	})
	if err != nil {
		return nil, err
	}

	return populateUser(&user), nil
}

func (s *Service) DeleteUser(ctx context.Context, userID format.UserID) error {
	return nil
}
