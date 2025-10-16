package ports

import (
	"context"
	"io"
)

// FileStorage defines methods for file storage operations
type FileStorage interface {
	Upload(ctx context.Context, bucket, filename string, reader io.Reader, size int64, contentType string) (string, error)
	Download(ctx context.Context, bucket, filename string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, filename string) error
	GetURL(ctx context.Context, bucket, filename string) (string, error)
}
