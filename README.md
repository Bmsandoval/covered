# Covered

**Understand your health insurance before the bill arrives.**

Covered’s core product is **citation-based Q&A over your own data** — answers grounded in your files with sources, confidence, and honest gaps. We start with **insurance documents** (upload), then expand formats, **cloud connectors** (e.g. Google Drive read-through), and later **opt-in** long-term vault storage.

> **Now (MVP):** “Ask My Insurance Plan” — upload SBCs, EOCs, cards, plan files.  
> **Next:** Any-format cited Q&A → Drive/connectors → cross-file search → optional “pin to vault.”

## Vision

Many products answer from the open internet or offer storage without proof. Covered’s bet: **your data, cited answers.** Insurance is the first domain; the pattern (ingest → chunk → retrieve → cite) generalizes. We do **not** require storing all your files on day one — connectors can read from your cloud; durable vault copy is **opt-in later**.

| Phase | Focus |
|-------|--------|
| **MVP** | Insurance upload — extract plan details, cited Q&A |
| **Post-MVP** | Any-format cited Q&A; Google Drive (read-through) |
| **Later** | Search across all sources; optional pin/save in Covered vault |

**How we build:** prioritized [GitHub issues](https://github.com/Bmsandoval/covered/issues); strategy in [`docs/planning/`](./docs/planning/) (not a substitute for issues).

Staged plan: [`staged-solution-plan.md`](./docs/planning/staged-solution-plan.md) · Connectors: [`connectors-and-storage.md`](./docs/planning/connectors-and-storage.md).

## The problem (MVP)

Insurance plans are dense PDFs, inconsistent terminology, and surprise bills. People need clarity on deductibles, copays, referrals, and in-network rules **before** they choose care — not a chatbot that guesses or a “vault” that does not exist yet.

## What Covered does today (MVP scope)

- **Upload** insurance materials (Summary of Benefits, Evidence of Coverage, card photos, deductible screenshots).
- **Extract** key plan details with **provenance** (document, page, section).
- **Ask questions** in plain English; answers come **only from your uploads**, with citations and confidence.
- **Flag gaps** and suggest what to ask your insurer when documents are silent.

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
> **Note:** Confirm with your insurer before making care decisions; labs, facility fees, and out-of-network care may differ.

## MVP vs later

| MVP (insurance upload) | Post-MVP |
|------------------------|----------|
| SBC, EOC, card, deductible screenshots | More formats (PDF, images, …) |
| Plan-specific extraction | Domain-agnostic cited Q&A |
| Stored uploads + chunks | Drive read-through (files stay in Drive by default) |
| “Ask My Insurance Plan” | Cross-file cited search |
| — | Opt-in **pin to vault** (long-term storage in Covered) |

**Explicitly out of MVP:** Google Drive OAuth, arbitrary formats, mandatory cloud vault, EMR/insurer login, bill prediction (“you will owe X”).

## How it works (MVP)

```mermaid
flowchart LR
  A[Upload insurance docs] --> B[Extract text and tables]
  B --> C[Plan model with provenance]
  C --> D[Dashboard]
  D --> E[Ask a question]
  E --> F[Retrieve relevant chunks]
  F --> G[Cited answer or honest gap]
```

1. **Ingest** — Parse uploads; preserve page numbers, sections, and tables.
2. **Normalize** — Store structured fields each linked to a source (insurance-specific in MVP).
3. **Chat** — Retrieve evidence, answer with citations or refuse when unsupported.

## Product principles

- **Citation-first** — Every answer names the document and location.
- **Interpreter, not oracle** — We explain what your documents say, not what you will definitely pay.
- **Honest uncertainty** — Low confidence and missing data are shown, not hidden.
- **Vault-ready architecture** — Prefer reusable ingestion and retrieval; avoid insurance-only dead ends when a generic pattern is cheap.
- **Trust over hype** — Calm, clear tone; not a replacement for your insurer or care team.

## Development

We build in small, reviewable slices tied to **GitHub issues** (planning docs in [`docs/planning/`](./docs/planning/) inform strategy; **issues are the work queue**).

| Branch | Purpose |
|--------|---------|
| [`develop`](https://github.com/Bmsandoval/covered/tree/develop) | Integration — feature PRs **squash-merge** here |
| `release-X-Y-Z` (e.g. `release-0-0-0`) | Cut from `develop` per minor version; hotfixes land here, then **backmerge** to `develop` |

There is **no `main` / `master`**.

**Workflow (short):**

1. Each minor version has a **parent release issue** (milestone `v0.x.0`) and **sub-issues** for the parts — implement **one sub-issue** at a time.
2. Before the first commit: `gh issue develop <N> --name issue-<N>-<slug> --checkout --base develop` (links branch on the issue).
3. PR to `develop`, title `Issue-<N> - <description>`, first line `Implements https://github.com/Bmsandoval/covered/issues/N` — **squash merge**.
4. Test on `develop`; cut `release-0-1-0` from `develop`; tag on the release branch (e.g. `v0.1.0`).
5. Hotfixes: branch from `release-*` → squash to release → **regular merge** back to `develop`.

Full rules (labels, milestones, no tool branding, env files): [AGENTS.md](./AGENTS.md) · Templates: [`docs/planning/issue-pr-workflow.md`](./docs/planning/issue-pr-workflow.md).

**Local setup:**

```bash
cp ex.env local.env
# Edit local.env with your values
```

## Status

Early stage — MVP (insurance documents + cited Q&A) in active planning and implementation. Post-MVP vault features are direction only until explicitly prioritized in issues.

## Disclaimer

Covered helps you read and question your **own insurance documents** (MVP). It is not medical, legal, or financial advice, and it is not a substitute for your insurer, provider, or licensed professional. Always confirm coverage and costs with your plan before receiving care.

## License

[MIT](./LICENSE) — Copyright (c) 2026 bryan.sandoval
