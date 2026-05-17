package ports

import (
	"context"

	"github.com/Bmsandoval/covered/internal/domain"
)

// ParsedPage is one page of extracted text.
type ParsedPage struct {
	Page    int
	Section string
	Text    string
}

// DocumentParser extracts text from uploads (PDF, images, etc.).
type DocumentParser interface {
	Parse(ctx context.Context, doc domain.Document, content []byte) ([]ParsedPage, error)
}
