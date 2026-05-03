package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

func NewS3Uploader(client *s3.Client, bucket, endpoint string) *S3Uploader {
	return &S3Uploader{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

// Upload uploads data to S3 and returns the public URL
func (u *S3Uploader) Upload(ctx context.Context, key string, data []byte) (string, error) {
	contentType := mime.TypeByExtension(filepath.Ext(key))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("s3 upload: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", u.endpoint, u.bucket, key), nil
}

// Delete removes an object from S3
func (u *S3Uploader) Delete(ctx context.Context, key string) error {
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})
	return err
}
