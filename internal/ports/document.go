package ports

import (
	"context"

	"github.com/Bmsandoval/covered/internal/domain"
)

// DocumentRepository persists document metadata.
type DocumentRepository interface {
	Create(ctx context.Context, doc domain.Document) error
	Get(ctx context.Context, id domain.DocumentID) (domain.Document, error)
	ListBySession(ctx context.Context, sessionID domain.SessionID) ([]domain.Document, error)
}
