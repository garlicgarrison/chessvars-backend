package likes

import (
	"context"
	"errors"

	"github.com/pafkiuq/backend/graph/model"
	"github.com/pafkiuq/backend/pkg/aws/s3"
	"github.com/pafkiuq/backend/pkg/aws/signer"
	"github.com/pafkiuq/backend/pkg/elasticsearch/index"
	"github.com/pafkiuq/backend/pkg/firestore"
	"github.com/pafkiuq/backend/pkg/format"
)

type Config struct {
	Firestore     firestore.Firestore
	Bucket        s3.Bucket
	Signer        signer.Signer
	ElasticSearch index.Index
}

type Service struct {
	fs firestore.Firestore
	es index.Index
}

func NewService(cfg Config) (*Service, error) {
	if cfg.Firestore == nil {
		return nil, errors.New("firestore required")
	}

	if cfg.Bucket == nil {
		return nil, errors.New("s3 bucket required")
	}

	if cfg.Signer == nil {
		return nil, errors.New("signer required")
	}

	if cfg.ElasticSearch == nil {
		return nil, errors.New("elastic search required")
	}

	return &Service{
		fs: cfg.Firestore,
		es: cfg.ElasticSearch,
	}, nil
}

type LikeVideoRequest struct {
	UserID  format.UserID  `json:"user_id"`
	VideoID format.VideoID `json:"video_id"`
}

func (s *Service) LikeVideo(ctx context.Context, request LikeVideoRequest) (*model.BasicMutationResponse, error) {
	//videoID := format.NewVideoID()
	return nil, nil
}
