package users

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	USERNAME_REGEX = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))$`
	EMAIL_REGEX    = `/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/`
)

type Config struct {
	Firestore firestore.Firestore
}

type service struct {
	fs firestore.Firestore
}

func NewService(cfg Config) (Service, error) {
	if cfg.Firestore == nil {
		return nil, errors.New("firestore required")
	}

	return &service{
		fs: cfg.Firestore,
	}, nil
}

func populateUser(user *UserDocument) *User {
	return &GetUserResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func (s *service) verifyUsername(ctx context.Context, username string) (bool, error) {
	re := regexp.MustCompile(USERNAME_REGEX)
	ok := re.Match([]byte(username))
	if !ok {
		return false, fmt.Errorf("invalid username")
	}

	userSnaps, err := s.getUsersRef().Where("username", "==", username).
		Documents(ctx).
		GetAll()
	if err != nil || len(userSnaps) != 0 {
		return false, err
	}

	return true, nil
}

func (s *service) getUsernameFromEmail(ctx context.Context, email string) (string, error) {
	re := regexp.MustCompile(EMAIL_REGEX)
	matches := re.FindAllStringSubmatch(email, -1)

	if len(matches) < 1 || len(matches[0]) < 1 {
		return "", fmt.Errorf("could not extract username")
	}

	username := matches[0][1]
	newUsername := username
	identifier := 1
	for {
		ok, err := s.verifyUsername(ctx, newUsername)
		if err != nil {
			return "", err
		}
		if ok {
			break
		}

		newUsername = username +
			strconv.FormatInt(int64(identifier), 10)
		identifier++
	}

	return newUsername, nil
}

func (s *service) CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, error) {
	user := UserDocument{
		UserID:    request.UserID,
		Email:     request.Email,
		CreatedAt: time.Now(),
	}

	username, err := s.getUsernameFromEmail(ctx, request.Email)
	if err != nil {
		fmt.Printf("[getUsernameFromEmail] error -- %s", err)
		username = ""
	}
	user.Username = username

	_, err = s.getUserRef(request.UserID).Create(ctx, user)
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
	var user UserDocument
	err := s.fs.RunTransaction(ctx, func(_ context.Context, t *firestore.Transaction) error {
		userSnap, err := t.Get(s.getUserRef(request.UserID))
		if err != nil {
			return err
		}

		err = userSnap.DataTo(&user)
		if err != nil {
			return err
		}

		if request.Username != nil {
			ok, err := s.verifyUsername(ctx, *request.Username)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("invalid username")
			}

			user.Username = *request.Username
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
