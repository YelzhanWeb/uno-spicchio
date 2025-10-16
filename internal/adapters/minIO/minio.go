// internal/adapters/minio/file_storage.go
package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	client *minio.Client
}

func NewFileStorage(endpoint, accessKey, secretKey string, useSSL bool) (*FileStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &FileStorage{client: client}, nil
}

func (fs *FileStorage) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := fs.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = fs.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Set bucket policy to public read
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, bucketName)

	err = fs.client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

func (fs *FileStorage) Upload(ctx context.Context, bucket, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := fs.client.PutObject(
		ctx,
		bucket,
		filename,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the object URL
	url := fmt.Sprintf("%s/%s/%s", fs.client.EndpointURL(), bucket, filename)
	return url, nil
}

func (fs *FileStorage) Download(ctx context.Context, bucket, filename string) (io.ReadCloser, error) {
	object, err := fs.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}

func (fs *FileStorage) Delete(ctx context.Context, bucket, filename string) error {
	err := fs.client.RemoveObject(ctx, bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (fs *FileStorage) GetURL(ctx context.Context, bucket, filename string) (string, error) {
	// For public buckets, just return the direct URL
	url := fmt.Sprintf("%s/%s/%s", fs.client.EndpointURL(), bucket, filename)
	return url, nil
}
