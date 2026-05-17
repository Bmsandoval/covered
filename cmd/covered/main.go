package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Bmsandoval/covered/internal/adapters/stub"
	"github.com/Bmsandoval/covered/internal/api"
	"github.com/Bmsandoval/covered/internal/app"
	"github.com/Bmsandoval/covered/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("local.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	deps := app.Deps{
		Sessions:  stub.SessionStore{},
		Documents: stub.DocumentRepository{},
		Chunks:    stub.ChunkRepository{},
		Blobs:     stub.BlobStore{},
		Cache:     stub.Cache{},
		Parser:    stub.DocumentParser{},
		Chat:      stub.ChatCompleter{},
		Embedder:  stub.Embedder{},
		Clock:     stub.Clock{},
	}

	srv := api.NewServer(cfg, deps)
	log.Printf("covered listening on %s (env=%s)", srv.Addr(), cfg.AppEnv)

	if err := http.ListenAndServe(srv.Addr(), srv.Handler()); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
