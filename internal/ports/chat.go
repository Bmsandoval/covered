package ports

import (
	"context"

	"github.com/Bmsandoval/covered/internal/domain"
)

// ChatRequest is input to the chat/LLM port.
type ChatRequest struct {
	SessionID   domain.SessionID
	Question    string
	DocumentIDs []domain.DocumentID
}

// Citation points to evidence in a document.
type Citation struct {
	DocumentID domain.DocumentID
	Page       int
	Section    string
	Snippet    string
}

// ChatResponse is the citation envelope returned to callers.
type ChatResponse struct {
	Answer     string
	Citations  []Citation
	Confidence string // high | medium | low
	Caveat     string
	Refused    bool
	GapHints   []string
}

// ChatCompleter runs grounded Q&A (LLM adapter implements this).
type ChatCompleter interface {
	Complete(ctx context.Context, req ChatRequest) (ChatResponse, error)
}
