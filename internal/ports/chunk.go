package ports

import (
	"context"

	"github.com/Bmsandoval/covered/internal/domain"
)

// ChunkRepository persists and retrieves text chunks.
type ChunkRepository interface {
	Insert(ctx context.Context, chunk domain.Chunk) error
	ListByDocument(ctx context.Context, documentID domain.DocumentID) ([]domain.Chunk, error)
	Search(ctx context.Context, sessionID domain.SessionID, query string, limit int) ([]domain.Chunk, error)
}
