package app

import "github.com/Bmsandoval/covered/internal/ports"

// Deps holds injected ports for use cases. Adapters are wired in cmd/covered.
type Deps struct {
	Sessions  ports.SessionStore
	Documents ports.DocumentRepository
	Chunks    ports.ChunkRepository
	Blobs     ports.BlobStore
	Cache     ports.Cache
	Parser    ports.DocumentParser
	Chat      ports.ChatCompleter
	Embedder  ports.Embedder
	Clock     ports.Clock
}
