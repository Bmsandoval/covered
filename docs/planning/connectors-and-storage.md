# Connectors and storage strategy

How Covered accesses user documents without assuming “we store everything forever” on day one.

Related: [product-vision.md](./product-vision.md), [staged-solution-plan.md](./staged-solution-plan.md).

---

## Core product (what we’re really selling)

**Citation-based Q&A over the user’s own data, in any supported format** — PDF, images, and later spreadsheets, exports, etc.

The **personal data vault** is a **container story**, not the first technical requirement. Long-term storage in Covered’s systems should be **opt-in and later**, not the default for every file.

---

## Three ways to access documents

| Mode | When | Where bytes live | Our persistence |
|------|------|------------------|-----------------|
| **Upload (MVP)** | User picks files in app | Copied to app storage (dev: disk; prod: object store) | File + chunks + embeddings until user deletes |
| **Connector read-through** | User links Google Drive (etc.) | Stays in their cloud | Prefer **no full file copy**; see below |
| **Vault pin (later)** | User explicitly “save to Covered” | Covered storage | Durable copy + same chunk/index pipeline |

---

## Connector model (e.g. Google Drive)

### Idea

User OAuth-connects Drive → selects files or folders → asks questions. We **fetch at need**, parse **in memory** (or temp disk), produce chunks and citations, answer with the same envelope as upload.

We **do not** need to mirror their entire Drive to ship value.

### What we might still persist (tradeoffs)

Pure “read into memory every question” is simple for privacy messaging but expensive and slow on every chat turn. Pragmatic **hybrid**:

| Data | Persist? | TTL / rule |
|------|----------|------------|
| OAuth tokens | Yes (encrypted) | Until revoked |
| File metadata (id, name, mime, modified) | Yes | Refresh on sync |
| Extracted **chunks** + embeddings | Optional cache | Short TTL (e.g. 24–72h) or until file `modifiedTime` changes |
| **Raw file bytes** | No by default | Only if user **pins to vault** |

On each Q&A or ingest job:

1. Check cache validity (TTL + Drive `modifiedTime`).
2. If stale → stream/download → extract → chunk in memory → update cache.
3. Answer from chunks; citations point to **Drive file name + page/section** (and link if we have stable file id).

### Why not only in-memory

- Re-parsing a 40-page SBC on every message is slow and costly.
- Embeddings are not “the document” but are derived data — disclose in privacy copy.
- Cache with TTL is a reasonable middle ground until user opts into vault storage.

### Implementation sketch (post-MVP issue)

- `DocumentSource` interface: `list()`, `fetch(ref)`, `metadata(ref)`
- Implementations: `UploadSource`, `GoogleDriveSource` (later: Dropbox, OneDrive, local folder)
- Ingest pipeline takes a `SourceRef`, not only `multipart.File`
- Env: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, redirect URIs in `ex.env`

### Risks / requirements

- OAuth consent screen and Google API verification for production
- Scopes: minimum read-only for selected files (`drive.readonly` or file-picker scope)
- Rate limits and large-file timeouts
- User disconnect = delete tokens + cached chunks for that source

**Do not start connectors in MVP issues** unless explicitly prioritized; design Stage 1 APIs so a source ref is not upload-only.

---

## MVP vs later (summary)

| Phase | Storage story |
|-------|----------------|
| **MVP (insurance upload)** | Simple upload + stored chunks — acceptable to prove ingestion and Q&A |
| **Post-MVP** | Any format + **connectors** with read-through + chunk cache |
| **Later** | **Opt-in vault** — durable Covered copy, export, retention policies |

---

## Privacy narrative (when connectors ship)

- “Your files can stay in Google Drive; we read them to answer questions.”
- “We cache extracted text snippets for speed; you can clear cache or disconnect.”
- “Pin to vault” = explicit choice to keep a copy on Covered.

---

## Open questions (resolve in issues)

- Maximum file size / page count for in-memory extract
- Whether chunk cache requires separate user consent
- Single Google account vs workspace accounts
- Insurance MVP: stay upload-only until Stage 5 release, then connector spike issue
