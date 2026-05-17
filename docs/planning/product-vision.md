# Product vision

## What we’re building (primary offering)

**Citation-based Q&A over the user’s own data** — any supported format (PDF, images, later more). Answers must be grounded in their files with **document, page, and section citations**, plus confidence and honest gaps.

That is the product. Everything else supports getting content into the Q&A pipeline.

## North star (container, not day one)

**Personal data vault** — one place to reason about your documents. Long-term **storage inside Covered** is a **later, opt-in** capability (“pin to vault”), not a prerequisite for the first releases.

Alternative path to content: **connectors** (Google Drive, etc.) that **read through** to the user’s cloud so we often **don’t store raw files**, only short-lived derived chunks for retrieval. See [connectors-and-storage.md](./connectors-and-storage.md).

## Prototype (current focus)

We are in the **prototype** phase (`v0.x` releases) — not MVP. Prototype proves **citation-backed Q&A** on **insurance upload** with thin releases; **MVP** is a later, stricter bar (see [product-phases.md](./product-phases.md)).

- **Domain:** Insurance documents only — “Ask My Insurance Plan” as the **demo narrative**.
- **Input:** Upload SBC, EOC, card, deductible screenshots.
- **Output (by v0.4.0 prototype):** Cited Q&A, plan summary with provenance, honest gaps.
- **Stack:** Go · **anonymous sessions** (no accounts in v0.x).
- **Not in prototype:** Drive connector, arbitrary formats, vault marketing, production MVP polish.

## MVP (future)

**MVP** means a **shippable** insurance product — significantly more than prototype (identity, ops, security, UX, scale). Plan it in issues **after** prototype slices land; use a **`v1.x`** milestone when ready, not `v0.x`.

## Post-prototype / platform trajectory

1. **Any-format cited Q&A** (still upload-first) — generalize ingest beyond insurance.
2. **Cloud connectors** — Google Drive (read-through + chunk cache; raw files stay in Drive unless pinned).
3. **Cross-corpus search** — ask and search across all connected/uploaded sources.
4. **Opt-in vault storage** — user explicitly saves copies in Covered for backup, export, retention.

## Architecture implication

Build **domain-agnostic** primitives early:

- `DocumentSource` abstraction (upload now; Drive later)
- Chunking with page/section metadata
- Retrieval + citation answer envelope
- Refuse when evidence is insufficient

Keep **insurance-specific** logic in extractors, prompts, and dashboard until generalized in issues.

## What we are not claiming today

- A full personal data vault with mandatory cloud storage
- Google Drive or other connectors (planned, not shipped)
- Support for every file format
- Replacement for professional advice in any domain

## Execution vs planning

| Source | Role |
|--------|------|
| **GitHub issues** | What to implement **now** — scope, acceptance criteria, priority |
| **`docs/planning/`** | Why and **which stage** — strategy, not a work queue |

Agents and contributors: **follow the active issue**; use planning docs for alignment, not to expand scope without a new issue.

## Delivery stages

See [product-phases.md](./product-phases.md) (prototype vs MVP), [staged-solution-plan.md](./staged-solution-plan.md) (stages 0–9), and [connectors-and-storage.md](./connectors-and-storage.md) (Drive / storage policy).
