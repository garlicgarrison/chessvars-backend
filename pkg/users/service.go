package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pafkiuq/backend/graph/model"
	"github.com/pafkiuq/backend/pkg/firestore"
	"github.com/pafkiuq/backend/pkg/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Firestore firestore.Firestore
}

type Service struct {
	fs firestore.Firestore
}

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
		Bio:      &user.Bio,
	}
}

func (s *Service) CreateUser(ctx context.Context, userID format.UserID) (*model.User, error) {
	user := UserDocument{
		UserID: userID,
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

		if input.Bio != nil {
			user.Bio = *input.Bio
		}

		if input.Gender != nil {
			switch *input.Gender {
			case model.GenderMale:
				user.Gender = GENDER_MALE
			case model.GenderFemale:
				user.Gender = GENDER_FEMALE
			case model.GenderNonbinary:
				user.Gender = GENDER_NONBINARY
			default:
				user.Gender = GENDER_UNKNOWN
			}
		}

		if input.Preferences != nil {
			user.Preferences = []Preference{}
			for _, pref := range input.Preferences {
				switch pref {
				case model.PreferenceStraight:
					user.Preferences = append(user.Preferences, PREF_STRAIGHT)
				case model.PreferenceGay:
					user.Preferences = append(user.Preferences, PREF_GAY)
				case model.PreferenceTrans:
					user.Preferences = append(user.Preferences, PREF_TRANS)
				}
			}
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
