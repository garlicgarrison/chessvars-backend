package elo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	K_FACTOR       = 32
	ELO_DIFFERENCE = 400
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

func (s *service) populateElo(e EloDocument) *Elo {
	return &Elo{
		UserID: e.UserID,
		Game:   e.GameType,
		Elo:    e.Elo,
	}
}

func (s *service) CreateElo(ctx context.Context, request CreateEloRequest) (*CreateEloResponse, error) {
	var elo EloDocument
	err := s.fs.RunTransaction(ctx, func(_ context.Context, t *firestore.Transaction) error {
		eloSnap, err := t.Get(s.getEloRef(request.UserID, request.Game))
		if err != nil {
			if status.Code(err) == codes.NotFound {
				elo = EloDocument{
					UserID:   request.UserID,
					GameType: request.Game,
					Elo:      DEFAULT_ELO,
				}
				return t.Create(s.getEloRef(request.UserID, request.Game), elo)
			}

			return err
		}

		return eloSnap.DataTo(&elo)
	})
	if err != nil {
		return nil, err
	}

	return s.populateElo(elo), nil
}

func (s *service) GetElo(ctx context.Context, request GetEloRequest) (*GetEloResponse, error) {
	eloSnap, err := s.getEloRef(request.UserID, request.Game).Get(ctx)
	if err != nil {
		return nil, err
	}

	var elo EloDocument
	err = eloSnap.DataTo(&elo)
	if err != nil {
		return nil, err
	}

	return s.populateElo(elo), nil
}

func (s *service) GetElos(ctx context.Context, request GetElosRequest) (*Elos, error) {
	eloSnaps, err := s.getElosRef(request.UserID).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	eloArr := make([]*Elo, 0)
	for _, eloSnap := range eloSnaps {
		var elo EloDocument
		err = eloSnap.DataTo(&elo)
		if err != nil {
			return nil, err
		}

		eloArr = append(eloArr, s.populateElo(elo))
	}

	return &Elos{
		Elos: eloArr,
	}, nil
}

func (s *service) UpdateElo(ctx context.Context, request UpdateEloRequest) (*UpdateEloResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	type eloStruct struct {
		elo *Elo
		err error
	}
	eloChan := make(chan eloStruct, 2)
	// get own elo
	go func() {
		defer wg.Done()

		elo, err := s.GetElo(ctx, GetEloRequest{
			UserID: request.UserID,
			Game:   request.Game,
		})
		eloChan <- eloStruct{
			elo: elo,
			err: err,
		}
	}()

	// get other elo
	go func() {
		defer wg.Done()

		elo, err := s.GetElo(ctx, GetEloRequest{
			UserID: request.OtherUserID,
			Game:   request.Game,
		})
		eloChan <- eloStruct{
			elo: elo,
			err: err,
		}
	}()
	wg.Wait()

	var err error
	var myElo Elo
	var otherElo Elo
	for e := range eloChan {
		if e.err != nil {
			err = fmt.Errorf("%sget elo error -- %s --", err, e.err)
			break
		}
		if e.elo.UserID == request.UserID {
			myElo = *e.elo
		} else if e.elo.UserID == request.OtherUserID {
			otherElo = *e.elo
		}
	}
	if err != nil {
		return nil, err
	}

	transformR1 := math.Pow(10, float64(myElo.Elo)/float64(ELO_DIFFERENCE))
	transformR2 := math.Pow(10, float64(otherElo.Elo)/float64(ELO_DIFFERENCE))
	expected := transformR1 / (transformR1 + transformR2)

	var s1 float64
	switch request.Status {
	case game.WIN:
		s1 = 1
	case game.LOSS:
		s1 = 0
	case game.DRAW:
		s1 = 0.5
	default:
		s1 = 0
	}

	newElo := int(math.Round(float64(myElo.Elo) + float64(K_FACTOR)*(s1-expected)))
	_, err = s.getEloRef(request.UserID, request.Game).Set(ctx, map[string]interface{}{
		"elo": newElo,
	}, firestore.MergeAll)
	if err != nil {
		return nil, err
	}

	myElo.Elo = newElo
	return &myElo, nil
}
