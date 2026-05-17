package ports

import (
	"context"
	"time"
)

// Cache is a generic TTL key-value store (memory now, Redis later).
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
