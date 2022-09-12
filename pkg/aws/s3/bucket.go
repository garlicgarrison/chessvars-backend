package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type Bucket interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	UploadObject(ctx context.Context, key string, body io.Reader) error
	DeleteObject(ctx context.Context, key string) error
	DeleteFolder(ctx context.Context, key string) error
}

// Returns a new client to the specified bucket
func NewBucket(awsConfig aws.Config, bucketName string, logger *zap.Logger) Bucket {
	client := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)

	return &bucket{
		bucket: bucketName,

		client:     client,
		uploader:   uploader,
		downloader: downloader,

		logger: logger,
	}
}

type bucket struct {
	bucket string

	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader

	logger *zap.Logger
}

func (b *bucket) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (b *bucket) UploadObject(ctx context.Context, key string, body io.Reader) error {
	_, err := b.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (b *bucket) DeleteObject(ctx context.Context, key string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	return err
}

// Deletes all object in the folder
func (b *bucket) DeleteFolder(ctx context.Context, key string) error {
	// TODO: Return more info than just an error

	isTruncated := true
	for isTruncated {
		// Get object list
		resp, err := b.client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(b.bucket),
			Prefix: aws.String(key),
		})
		if err != nil {
			return err
		}

		isTruncated = resp.IsTruncated

		for _, obj := range resp.Contents {
			_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(b.bucket),
				Key:    obj.Key,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
