package ports

import (
	"context"

	"github.com/Bmsandoval/covered/internal/domain"
)

// SessionStore manages anonymous sessions.
type SessionStore interface {
	Create(ctx context.Context) (domain.SessionID, error)
	Exists(ctx context.Context, id domain.SessionID) (bool, error)
}
