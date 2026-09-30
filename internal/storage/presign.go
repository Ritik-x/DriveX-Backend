package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ( s *S3Storage) GenerateUploadUrl(ctx context.Context , key string , contentType string) (string , error){
presignClient := s3.NewPresignClient(s.Client)

request , err := presignClient.PresignPutObject(

	ctx ,
	&s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
	},
	func(options *s3.PresignOptions) {
			options.Expires = 15 * time.Minute
		},

)

	if err != nil {
		return "", err
	}

	return request.URL, nil

}


func ( s *S3Storage) GenerateDownloadUrl(ctx context.Context  , key string ) ( string , error){
	presignClient := s3.NewPresignClient(s.Client)

	request, err := presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
		func(options *s3.PresignOptions) {
			options.Expires = 15 * time.Minute
		},
	)
	if err != nil {
		return "", err
	}

	return request.URL, nil
}