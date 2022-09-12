package videos

import (
	"context"
	"errors"
	"io"

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
	fs     firestore.Firestore
	bucket s3.Bucket
	signer signer.Signer
	es     index.Index
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
		fs:     cfg.Firestore,
		bucket: cfg.Bucket,
		signer: cfg.Signer,
		es:     cfg.ElasticSearch,
	}, nil
}

type UploadVideoRequest struct {
	UserID      format.UserID      `json:"user_id"`
	VideoFormat format.VideoFormat `json:"video_format"`
	Video       io.Reader          `json:"video"`
}

type UploadEmbeddedVideoRequest struct {
	OwnerID     format.UserID `json:"owner_id"`
	Link        string        `json:"link"`
	Tags        []string      `json:"tags"`
	Description string        `json:"description"`
}

func (s *Service) UploadVideo(ctx context.Context, request UploadVideoRequest) (*model.Video, error) {
	//videoID := format.NewVideoID()
	return nil, nil
}

func (s *Service) UploadEmbeddedVideo(ctx context.Context, request UploadEmbeddedVideoRequest) (*model.Video, error) {
	return nil, nil
}
