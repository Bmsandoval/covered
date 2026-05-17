# Covered

**Prototype · Citation-backed Q&A over your insurance documents**

Understand your health insurance before the bill arrives — with **sources**, not guesses.

Covered’s core product is **citation-based Q&A over your own data** — answers grounded in your files with document, page, and section references, confidence, and honest gaps. We are in the **prototype** phase (`v0.x`): prove the loop on insurance **upload**, then pursue a full **MVP** later (`v1.x`, much higher bar). See [product phases](docs/planning/product-phases.md).

> **Now (prototype):** Insurance upload → chunks → cited Q&A (Go, anonymous sessions).  
> **Later (MVP):** Shippable “Ask My Insurance Plan” product.  
> **Platform:** Any-format Q&A → Drive/connectors → cross-file search → optional pin-to-vault.

## Vision

Many products answer from the open internet or offer storage without proof. Covered’s bet: **your data, cited answers.** Insurance is the first domain; the pattern (ingest → chunk → retrieve → cite) generalizes.

| Phase | Focus |
|-------|--------|
| **Prototype** (`v0.x`) | Thin releases — upload, cited Q&A, demo UI, hardening on SBCs |
| **MVP** (`v1.x`, planned) | Production-ready insurance product (beyond prototype) |
| **Platform** | Any-format cited Q&A; Google Drive (read-through); cross-corpus search; opt-in vault |

**How we build:** prioritized [GitHub issues](https://github.com/Bmsandoval/covered/issues); strategy in [`docs/planning/`](docs/planning/).

Staged plan: [`staged-solution-plan.md`](docs/planning/staged-solution-plan.md) · Phases: [`product-phases.md`](docs/planning/product-phases.md) · Connectors: [`connectors-and-storage.md`](docs/planning/connectors-and-storage.md).

## The problem (insurance first)

Insurance plans are dense PDFs, inconsistent terminology, and surprise bills. People need clarity on deductibles, copays, referrals, and in-network rules **before** they choose care — not a chatbot that guesses.

## What the prototype aims to demonstrate

- **Upload** insurance materials (SBC, EOC, card photos, deductible screenshots).
- **Extract** key plan details with **provenance** (document, page, section).
- **Ask questions** in plain English; answers from **your uploads only**, with citations.
- **Flag gaps** when documents are silent.

### Example questions

- Do I need a referral for a specialist?
- What is my ER copay?
- How much is therapy?
- What’s my deductible?
- What happens if I go out of network?

### How answers should read

Not: *“You owe $40.”*

Instead:

> Your uploaded Summary of Benefits lists a **$40 specialist copay** for in-network care.  
> **Source:** SBC page 3, “Specialist Visit.”  
> **Confidence:** High.  
> **Note:** Confirm with your insurer before making care decisions.

## Prototype vs later

| Prototype (`v0.x`) | MVP / platform (later) |
|--------------------|-------------------------|
| Go API, **anonymous sessions** | Accounts, production ops (MVP) |
| SBC, EOC, card, screenshots | More formats; Drive connectors |
| Local/demo deploy OK | Hosted, scaled product (MVP) |
| Prove citation loop | Shippable “Ask My Insurance Plan” |

**Out of prototype:** user accounts, Google Drive OAuth, arbitrary formats, mandatory vault, EMR/insurer login, bill prediction.

## How it works (target architecture)

```mermaid
flowchart LR
  A[Upload insurance docs] --> B[Extract text and tables]
  B --> C[Plan model with provenance]
  C --> D[Dashboard]
  D --> E[Ask a question]
  E --> F[Retrieve relevant chunks]
  F --> G[Cited answer or honest gap]
```

## Product principles

- **Citation-first** — Every answer names the document and location.
- **Interpreter, not oracle** — Explain what documents say, not guaranteed final cost.
- **Honest uncertainty** — Show low confidence and missing data.
- **Vault-ready design** — Reusable ingest/retrieval; avoid dead-end shortcuts.
- **Trust over hype** — Calm tone; not a substitute for your insurer or care team.

## Development

We build in small slices tied to **GitHub issues**. Agents read [`AGENTS.md`](AGENTS.md) and [`docs/planning/`](docs/planning/) before coding.

| Branch | Purpose |
|--------|---------|
| [`develop`](https://github.com/Bmsandoval/covered/tree/develop) | Integration — feature PRs **squash-merge** here |
| `release-X-Y-Z` | Cut from `develop` per minor version; hotfixes, then **backmerge** |

There is **no `main` / `master`**.

**Prototype release train (Option B):** `v0.1.0` upload → `v0.2.0` cited Q&A → `v0.3.0` UI → `v0.4.0` hardening. Details in [product-phases.md](docs/planning/product-phases.md).

**Local setup:**

```bash
cp ex.env local.env
# Edit local.env (SESSION_SECRET, keys when wired)
```

## Status

| Release | State |
|---------|--------|
| **v0.0.0** | Shipped — planning, workflow, `AGENTS.md` ([tag](https://github.com/Bmsandoval/covered/releases/tag/v0.0.0)) |
| **v0.1.0** | Next — **Prototype: document upload** (Go skeleton + PDF → chunks) |

**Stack:** Go · anonymous sessions (prototype).

## GitHub About line

Recommended repo description (you can tweak):

`Prototype · Citation-backed Q&A over your insurance documents`

Optional topics: `prototype`, `golang`, `insurance`, `citations`.

## Disclaimer

Covered is a **prototype** tool to help you read and question **your own insurance documents**. It is not medical, legal, or financial advice. Confirm coverage and costs with your plan before receiving care.

## License

[MIT](./LICENSE) — Copyright (c) 2026 bryan.sandoval
