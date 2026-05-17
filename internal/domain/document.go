package domain

import "time"

// Document is upload metadata (bytes live in BlobStore).
type Document struct {
	ID        DocumentID
	SessionID SessionID
	Filename  string
	MimeType  string
	SizeBytes int64
	CreatedAt time.Time
}

// Chunk is searchable text with provenance.
type Chunk struct {
	ID         ChunkID
	DocumentID DocumentID
	SessionID  SessionID
	Page       int
	Section    string
	Text       string
}
