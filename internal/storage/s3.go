package storage

import (
	"context"
	"fmt"

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