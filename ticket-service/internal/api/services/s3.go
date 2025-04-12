package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

func NewS3Service(client *s3.Client, bucketName string, region string) *S3Service {
	return &S3Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}
}

func (s *S3Service) UploadObject(ctx context.Context, filePath string, key string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	url := generateObjectUrl(s.bucketName, s.region, key)
	return url, nil
}

func (s *S3Service) DownloadObject(ctx context.Context, key string, destination string) error {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil
	}
	defer result.Body.Close()

	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, result.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Service) UploadBytes(ctx context.Context, data []byte, key string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", nil
	}
	url := generateObjectUrl(s.bucketName, s.region, key)
	return url, nil
}

func generateObjectUrl(bucketName string, region string, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key)
}
