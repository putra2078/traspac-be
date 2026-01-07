package storage

import (
	"context"
	"io"
)

type StorageRepository interface {
	Upload(ctx context.Context, bucket string, key string, data io.Reader, contentType string) error
	Delete(ctx context.Context, bucket string, key string) error
	GetURL(bucket string, key string) string
}
