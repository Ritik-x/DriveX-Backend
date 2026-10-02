package storage

import (
	"context"
	"fmt"
	"io"

	"drivex/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client *s3.Client
	Bucket string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					cfg.AWSAccessKeyID,
					cfg.AWSSecretAccessKey,
					"",
				),
			),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Storage{
		Client: client,
		Bucket: cfg.AWSS3Bucket,
	}, nil
}



func (s *S3Storage) GetObjectMetadata(
	ctx context.Context,
	key string,
) (int64, string, error) {

	result, err := s.Client.HeadObject(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)

	if err != nil {
		return 0, "", err
	}

	contentType := ""

	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	var size int64

	if result.ContentLength != nil {
		size = *result.ContentLength
	}

	return size, contentType, nil
}
func (s *S3Storage) PutObject(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
) error {

	_, err := s.Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.Bucket),
			Key:         aws.String(key),
			Body:        body,
			ContentType: aws.String(contentType),
		},
	)

	return err
}


func (s *S3Storage) GetObject(
	ctx context.Context,
	key string,
) (*s3.GetObjectOutput, error) {

	return s.Client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)
}

func (s *S3Storage) DeleteObject(
	ctx context.Context,
	key string,
) error {
	if key == "" {
		return nil
	}

	_, err := s.Client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)

	return err
}