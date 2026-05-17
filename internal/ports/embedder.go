package ports

import "context"

// Embedder produces vector embeddings for retrieval.
type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}
