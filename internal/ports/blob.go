package ports

import (
	"context"
	"io"

	"github.com/Bmsandoval/covered/internal/domain"
)

// BlobStore stores raw upload bytes.
type BlobStore interface {
	Put(ctx context.Context, key domain.BlobKey, r io.Reader, size int64) error
	Open(ctx context.Context, key domain.BlobKey) (io.ReadCloser, error)
	Delete(ctx context.Context, key domain.BlobKey) error
}
