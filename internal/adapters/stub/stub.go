// Package stub provides no-op adapters so the app compiles before real implementations land.
package stub

import (
	"context"
	"io"
	"time"

	"github.com/Bmsandoval/covered/internal/domain"
	"github.com/Bmsandoval/covered/internal/ports"
)

type SessionStore struct{}

func (SessionStore) Create(context.Context) (domain.SessionID, error) {
	return "", ports.ErrNotImplemented
}

func (SessionStore) Exists(context.Context, domain.SessionID) (bool, error) {
	return false, ports.ErrNotImplemented
}

type DocumentRepository struct{}

func (DocumentRepository) Create(context.Context, domain.Document) error {
	return ports.ErrNotImplemented
}

func (DocumentRepository) Get(context.Context, domain.DocumentID) (domain.Document, error) {
	return domain.Document{}, ports.ErrNotImplemented
}

func (DocumentRepository) ListBySession(context.Context, domain.SessionID) ([]domain.Document, error) {
	return nil, ports.ErrNotImplemented
}

type ChunkRepository struct{}

func (ChunkRepository) Insert(context.Context, domain.Chunk) error {
	return ports.ErrNotImplemented
}

func (ChunkRepository) ListByDocument(context.Context, domain.DocumentID) ([]domain.Chunk, error) {
	return nil, ports.ErrNotImplemented
}

func (ChunkRepository) Search(context.Context, domain.SessionID, string, int) ([]domain.Chunk, error) {
	return nil, ports.ErrNotImplemented
}

type BlobStore struct{}

func (BlobStore) Put(context.Context, domain.BlobKey, io.Reader, int64) error {
	return ports.ErrNotImplemented
}

func (BlobStore) Open(context.Context, domain.BlobKey) (io.ReadCloser, error) {
	return nil, ports.ErrNotImplemented
}

func (BlobStore) Delete(context.Context, domain.BlobKey) error {
	return ports.ErrNotImplemented
}

type Cache struct{}

func (Cache) Get(context.Context, string) ([]byte, bool, error) {
	return nil, false, ports.ErrNotImplemented
}

func (Cache) Set(context.Context, string, []byte, time.Duration) error {
	return ports.ErrNotImplemented
}

func (Cache) Delete(context.Context, string) error {
	return ports.ErrNotImplemented
}

type DocumentParser struct{}

func (DocumentParser) Parse(context.Context, domain.Document, []byte) ([]ports.ParsedPage, error) {
	return nil, ports.ErrNotImplemented
}

type ChatCompleter struct{}

func (ChatCompleter) Complete(context.Context, ports.ChatRequest) (ports.ChatResponse, error) {
	return ports.ChatResponse{}, ports.ErrNotImplemented
}

type Embedder struct{}

func (Embedder) Embed(context.Context, []string) ([][]float32, error) {
	return nil, ports.ErrNotImplemented
}

type Clock struct{}

func (Clock) Now() time.Time {
	return time.Now()
}
