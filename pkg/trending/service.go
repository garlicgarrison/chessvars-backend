package trending

import (
	"context"
	"errors"

	"github.com/pafkiuq/backend/graph/model"
	"github.com/pafkiuq/backend/pkg/elasticsearch/index"
	"github.com/pafkiuq/backend/pkg/format"
)

type Config struct {
	ElasticSearch index.Index
}

type Service struct {
	es index.Index
}

func NewService(cfg Config) (*Service, error) {
	if cfg.ElasticSearch == nil {
		return nil, errors.New("elastic search required")
	}

	return &Service{
		es: cfg.ElasticSearch,
	}, nil
}

type GetFeedRequest struct {
	UserID format.UserID `json:"user_id"`
}

func (s *Service) GetFeed(ctx context.Context, request GetFeedRequest) (*model.Videos, error) {
	//videoID := format.NewVideoID()
	return nil, nil
}
