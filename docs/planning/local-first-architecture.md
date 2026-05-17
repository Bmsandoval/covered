# Local-first architecture (prototype)

Prototype runs **entirely on a developer machine** — no required cloud services. Every external dependency is behind an **interface** so we can swap implementations later (SQLite → MySQL, in-memory cache → Redis, local disk → object storage) without rewriting business logic.

**Issues implement slices of this;** this doc is the agreed shape for Go code in `v0.1.0+`.

Related: [product-phases.md](./product-phases.md), [staged-solution-plan.md](./staged-solution-plan.md), [connectors-and-storage.md](./connectors-and-storage.md).

---

## Principles

1. **Local defaults** — one command (or `make run`) with `local.env`; data under `./data/` (gitignored).
2. **Ports, not packages** — domain and use cases depend on **interfaces** defined in `internal/ports` (or `domain`), not on `database/sql`, Redis, or S3 SDKs.
3. **Constructor injection** — wire concrete adapters in `cmd/covered/main.go` (or a small `internal/wire` package). No global singletons for DB/cache/storage/LLM.
4. **Thin transport** — HTTP handlers parse/validate, call **services**, return DTOs. No SQL or LLM calls in handlers.
5. **Test with fakes** — unit tests pass **in-memory or fake** implementations of each port; integration tests may use **real SQLite + temp dir** locally.
6. **Portable SQL** — schema and queries written so a **MySQL** adapter can be added later; avoid SQLite-only shortcuts in shared migration SQL without a documented exception.

---

## Default implementations (prototype)

| Concern | Prototype (local) | Later swap (examples) |
|---------|-------------------|------------------------|
| **Database** | SQLite file (`./data/covered.db`) | MySQL / Postgres |
| **Cache** | In-process memory (or `ristretto` / `sync.Map` wrapper) | Redis |
| **File uploads** | Local filesystem (`./data/uploads`) | S3-compatible object store |
| **Sessions** | DB table or signed cookie + DB row | Same interface; Redis session store optional later |
| **LLM / chat** | HTTP client to provider API | Same interface; mock in tests |
| **Embeddings** | Provider API or local stub | Same pattern as LLM |
| **OCR / PDF** | Local libraries or provider API behind `DocumentParser` port | Swap provider without changing ingest service |

Select backend via **`ex.env` / `local.env`** (e.g. `DATABASE_DRIVER=sqlite`, `CACHE_BACKEND=memory`, `STORAGE_BACKEND=local`). `main` reads config and constructs the matching adapter.

---

## Suggested Go layout

```
cmd/covered/                 # main: load config, construct adapters, start HTTP
internal/
  config/                    # env → typed Config
  domain/                    # entities, value types (SessionID, DocumentID, Citation)
  ports/                     # interfaces only (no adapter imports)
  app/                       # use cases / services (orchestration)
    ingest/
    documents/
    chat/
    session/
  adapters/
    sqlite/                  # repositories, migrations
    memory/                  # cache
    localfs/                 # blob storage
    session/                 # cookie + repo
    llm/                     # OpenAI/etc. client implementing ChatCompleter
    pdf/                     # extract/chunk (or subpackage per format)
  api/                       # HTTP: routes, middleware, request/response DTOs
migrations/                  # SQL files (portable where possible)
data/                        # gitignored: sqlite db, uploads, cache spill
```

**Dependency rule:** `api` → `app` → `ports` ← `adapters`. `domain` has no imports from `adapters` or `api`.

---

## Ports (interfaces to define early)

Define small, purpose-specific interfaces — prefer **several small ports** over one giant `Repository`.

| Port | Responsibility |
|------|----------------|
| `SessionStore` | Create/load session by id; tie to anonymous session cookie |
| `DocumentRepository` | Metadata CRUD, list by session |
| `ChunkRepository` | Persist/search chunks (or split read vs write if needed) |
| `BlobStore` | Store/retrieve upload bytes by key (`Put`, `Get`, `Delete`, `Open`) |
| `Cache` | Generic TTL cache for hot keys (embeddings, rate limits) |
| `DocumentParser` | PDF/image → text + page boundaries |
| `ChatCompleter` | `Complete(ctx, ChatRequest) (ChatResponse, error)` — **all** LLM request/response shaping here |
| `Embedder` | `Embed(ctx, texts) ([][]float32, error)` if separate from chat |
| `Clock` | `Now()` for tests |

**Chat service** (`app/chat`) depends on `ChatCompleter`, `ChunkRepository`, and policies — not on `http.Client` or vendor SDK types in handlers.

### Chat request/response (keep inside port)

```go
// ports/chat.go — illustrative; refine in code issues.

type ChatRequest struct {
    SessionID   domain.SessionID
    Question    string
    DocumentIDs []domain.DocumentID // scope retrieval
}

type Citation struct {
    DocumentID domain.DocumentID
    Page       int
    Section    string
    Snippet    string
}

type ChatResponse struct {
    Answer     string
    Citations  []Citation
    Confidence string // high | medium | low
    Caveat     string
    Refused    bool
    GapHints   []string // questions for insurer when refused
}
```

HTTP layer maps JSON ↔ `ChatRequest` / `ChatResponse`; **LLM adapter** maps `ChatResponse` envelope ↔ provider API.

---

## Dependency injection

**Prototype:** explicit wiring in `main`:

```go
db := sqlite.Open(cfg.DatabaseDSN)
cache := memory.NewCache()
blobs := localfs.New(cfg.StoragePath)
llm := llm.NewOpenAI(cfg.LLM)
docSvc := documents.NewService(ports.DocumentDeps{Repo: sqlite.NewDocumentRepo(db), Blobs: blobs, ...})
// pass docSvc into api.Router
```

**Testing:** `app` tests receive fakes:

```go
type fakeChatCompleter struct { answer string }
func (f *fakeChatCompleter) Complete(ctx context.Context, req ports.ChatRequest) (ports.ChatResponse, error) { ... }
```

Optional later: `google/wire` or similar — not required for v0.1.0.

---

## SQL: SQLite now, MySQL later

| Practice | Why |
|----------|-----|
| **Repository per aggregate** | Swap `sqlite.DocumentRepo` for `mysql.DocumentRepo` without touching services |
| **Migrations in `migrations/`** | `golang-migrate` or `goose`; one logical schema |
| **Portable types** | `TEXT`, `INTEGER`, `BLOB`, `TIMESTAMP` (UTC); store JSON as `TEXT` |
| **Avoid in shared migrations** | SQLite-only `AUTOINCREMENT` quirks without MySQL equivalent — use sequences/IDs consistently |
| **IDs** | Prefer `TEXT` UUIDs or `BIGINT` — document choice in first migration issue |
| **Queries** | Keep SQL in adapter package, not scattered; optional **sqlc** with two drivers later |

When MySQL is added: new `adapters/mysql` + `DATABASE_DRIVER=mysql` + same interfaces; run parallel integration tests.

**Vector search:** if using sqlite-vec locally, hide behind `ChunkIndex` port so pgvector/MySQL variants do not leak into `app/chat`.

---

## Local data directories

| Path | Contents |
|------|----------|
| `./data/covered.db` | SQLite |
| `./data/uploads/` | Raw uploads |
| `./data/cache/` | Optional spill for large cache entries |

All under `.gitignore`. Document in README/`ex.env`.

---

## Configuration (`ex.env`)

Use **backend switches**, not separate code paths in services:

```bash
DATABASE_DRIVER=sqlite
DATABASE_DSN=file:./data/covered.db
CACHE_BACKEND=memory
STORAGE_BACKEND=local
STORAGE_LOCAL_PATH=./data/uploads
LLM_PROVIDER=openai   # or mock for tests
```

Services read **ports**; only `main` reads these variables.

---

## Testing strategy

| Layer | Approach |
|-------|----------|
| **domain** | Pure unit tests, no mocks |
| **app / services** | Fake ports, table-driven tests |
| **adapters** | Integration tests with real SQLite + temp dir; tag `integration` |
| **api** | `httptest` + services backed by fakes |
| **LLM** | Always mock `ChatCompleter` in CI; optional manual integration with API key locally |

Do not call live LLM APIs in default `go test ./...`.

---

## What not to do in prototype

- Hard-code `*sql.DB`, Redis client, or S3 path in `app` or `api` packages
- Put business rules in HTTP handlers or LLM prompt strings only (extract prompt builders to `app` or `adapters/llm`)
- Assume Postgres-specific features in migrations without a MySQL plan
- Skip interfaces “until we need MySQL” — add ports when the **first** adapter is written

---

## Mapping to v0.1.0 issues

Suggested sub-issue themes under **Prototype: document upload**:

1. Go skeleton + config loader + health check  
2. Port interfaces + fake implementations for tests  
3. SQLite migrations + session + document repositories  
4. Local `BlobStore` + upload API  
5. PDF parse adapter + chunk repository  
6. Wire `main` with DI; integration test upload → chunk row  

Chat/LLM ports land in **v0.2.0** but **define `ChatCompleter` in v0.1.0** if cheap, so ingest does not block the shape.
