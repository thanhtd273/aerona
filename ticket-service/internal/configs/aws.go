package configs

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func NewAWSConfig(logger *zap.Logger) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION_ID")),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"), SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
		}))
	if err != nil {
		logger.Fatal("Failed to load AWS configuration", zap.Error(err))
	}
	return cfg
}

func NewS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}
