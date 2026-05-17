# Agent instructions — Covered

## Sources of truth (read this first)

| Priority | Source | Use it for |
|----------|--------|------------|
| **1 — Work queue** | **GitHub issues** (prioritized with maintainer) | **What to build right now** — scope, acceptance criteria, stage label in issue body if helpful |
| **2 — Strategy** | **`docs/planning/`** | Why, stages, constraints, connector/storage policy |
| **3 — Process** | This file (`AGENTS.md`) | How to branch, PR, env, and behave |

**Rules**

- **Do not implement from planning docs alone.** Every change ties to an **open, agreed issue**. If planning implies new work, ask for or draft an issue — do not expand scope silently.
- **Before coding:** read the **active issue**; skim relevant planning docs only for context.
- **After shipping:** close/link the issue; update planning docs only when strategy materially changes (via issue or maintainer direction).
- Issues win on scope conflicts; planning docs win on long-term direction.

## Product vision

**Primary offering:** **Citation-based Q&A** over the user’s own data (PDF, images, more formats later) — answers grounded in their files with sources, confidence, and caveats.

**North star (later):** **Personal data vault** — but **long-term storage in Covered is opt-in later** (“pin to vault”), not required for early releases. **Connectors** (e.g. Google Drive) may **read through** to the user’s cloud without keeping raw files; see [connectors-and-storage.md](./docs/planning/connectors-and-storage.md).

**MVP (build now):** **Insurance documents via upload** — “Ask My Insurance Plan.” Same core pattern (ingest → chunk → retrieve → cite). Do not implement Drive connectors, arbitrary formats, or vault persistence unless an issue explicitly says so.

**After MVP (issue-driven):** Any-format cited Q&A → cloud connectors → cross-corpus search → opt-in vault storage.

When implementing MVP features, prefer **reusable building blocks** (`DocumentSource`, chunking, retrieval, citation envelope) over one-off insurance shortcuts unless the issue explicitly calls for insurance-specific logic.

---

**MVP positioning:** Covered is an **insurance document interpreter with receipts**, not an AI health advisor or cost predictor. Build and review MVP changes against that positioning.

## Product goal (MVP)

**“Ask My Insurance Plan”** — Users upload or enter insurance information, then ask questions such as:

- “Do I need a referral for a specialist?”
- “What is my ER copay?”
- “How much is therapy?”
- “What’s my deductible?”
- “What happens if I go out of network?”

The app answers **only** from the user’s uploaded documents, with citations. If the documents do not support an answer, say so and suggest what to ask the insurer.

## Non-goals (do not build in MVP)

**Insurance / trust boundaries**

- EMR login or clinical record access
- Insurer portal login or scraping
- Claims / EOB-heavy integrations (EOBs may come later)
- Provider-side or billing-workflow tools
- Definitive “you will owe $X” predictions

**Post-MVP / north star (do not build until prioritized in an issue)**

- Arbitrary non-insurance document types
- Google Drive / Dropbox / OneDrive connectors
- Cross-document “search everything” product surface
- Opt-in **vault** long-term storage (pin/save copy in Covered)
- Multi-domain dashboards beyond what the current issue requires

Keeping MVP to **insurance document upload + cited Q&A** reduces legal and technical risk while we prove the citation-first pattern.

## Hard rules for the chat agent

Every user-facing answer must include:

1. **Answer** — Plain language, grounded in documents only
2. **Source document** — Filename or type (e.g. Summary of Benefits)
3. **Page / section / line reference** — As specific as extraction allows
4. **Confidence** — High / Medium / Low (or equivalent)
5. **Plain-English caveat** — Encourage confirming with the insurer before care decisions

### Do

- Phrase like: “Your uploaded Summary of Benefits lists a $40 specialist copay for in-network care. Source: SBC page 3, Specialist Visit section.”
- Add hedging when appropriate: “I found the likely answer, but confirm with your insurer before making care decisions.”
- Refuse or hedge when evidence is insufficient; offer **questions to ask the insurer**

### Do not

- State certainty about final patient responsibility (“You owe $40.”)
- Answer from general medical or insurance knowledge when the user’s documents do not support it
- Present outputs as medical, legal, or financial advice

## V1 input sources (priority order)

**Support early:**

- Summary of Benefits and Coverage (SBC) PDF
- Evidence of Coverage (EOC) PDF
- Insurance card photo
- Deductible / out-of-pocket screenshots

**Defer:**

- EOB parsing at scale (later phase)

## Core behavior

1. **Ingest** — Upload PDF/image; extract text, page numbers, sections, tables, and key fields (plan year, carrier, deductible, copays, coinsurance, OOP max, referral rules, prior auth, in-network vs out-of-network).
2. **Normalize** — Structured plan model where **every field carries provenance** (value + source document + page + section).
3. **Retrieve & answer** — Classify the question, retrieve relevant chunks, answer only from those chunks, cite sources, refuse when evidence is weak.

Example provenance shape:

```json
{
  "value": 40,
  "source_document": "summary_of_benefits.pdf",
  "page": 3,
  "section": "Specialist Visit"
}
```

Example answer envelope (conceptual):

```
Answer: In-network primary care appears to be $25 per visit.
Source: 2026 Blue Cross Summary of Benefits, page 2, “Primary Care Visit.”
Confidence: High.
Note: This may not include labs, imaging, facility fees, or out-of-network charges.
```

## V1 feature set (target)

- Upload insurance documents
- Extract and surface key plan details (dashboard)
- Ask questions with cited answers
- Compare in-network vs out-of-network when documents allow
- Explain confusing terms (from user docs or controlled glossary tied to citations)
- Flag missing information
- Generate “questions to ask your insurer” when coverage is unclear

## Ideal first user flow

1. User uploads insurance PDF(s).
2. App responds: “I found your deductible, copays, out-of-pocket max, and prescription coverage” (only for fields actually extracted).
3. Dashboard shows structured summary with links to sources.
4. User asks: “How much would urgent care cost?”
5. App answers with citation; if unclear: “The document does not fully answer this. Ask your insurer: ‘Is urgent care subject to deductible or only a copay?’”

## Positioning (copy and UX)

**MVP — use:** “Understand your insurance documents before you get surprised by a bill.”

**MVP — avoid:** “Predict exactly what your care will cost.”

**Long-term (do not lead with in MVP marketing):** “Your personal data vault” / “search all your documents with citations” — accurate as direction, premature as the primary promise until post-MVP.

Brand tone: trustworthy, calm, financially literate — not clinical, not “AI hype,” not insurer-like. The vault vision should feel like **privacy and receipts**, not “upload everything to our cloud” hype.

## Implementation guidance

- Prefer **upload** in MVP; design ingest around a **`DocumentSource`** abstraction so connectors can be added without rewrite (see [connectors-and-storage.md](./docs/planning/connectors-and-storage.md)).
- Do not add insurer portal or EMR integrations.
- Treat extraction errors as first-class: show confidence and missing fields openly.
- Store and display **provenance** alongside normalized values; never show a number without a source path when one exists.
- Retrieval layer: classify question → chunk retrieval → constrained generation → citation validation before response.
- When adding dependencies or external APIs, document data handling and retention expectations.

## Safety and compliance mindset

- Users own their documents; minimize retention where possible; encrypt at rest and in transit.
- Disclaimers on every chat surface: not a substitute for insurer, provider, or licensed advisor.
- Log and test refusal paths (no relevant chunks, conflicting pages, low OCR quality).

## Planning documents

All **LLM- or agent-generated planning artifacts** must be written under **`docs/planning/`**, not in the repo root, `AGENTS.md`, or application source trees.

Examples of what belongs here:

- Architecture and system design drafts
- MVP / milestone breakdowns and roadmaps
- Spike notes, ADRs, and technical decision write-ups
- Extraction or retrieval experiment plans
- Feature specs and user-flow outlines created during planning sessions

**Do not** use `docs/planning/` for:

- End-user documentation (use `README.md` or future `docs/` guides)
- Generated API reference or code-owned docs tied to packages
- Secrets, credentials, or real user insurance documents

Use clear filenames (e.g. `mvp-scope.md`, `ingestion-spike.md`, `adr-001-document-storage.md`). Prefer updating an existing planning doc over creating duplicates when the topic is the same.

**Key planning docs:** [product-vision.md](./docs/planning/product-vision.md), [staged-solution-plan.md](./docs/planning/staged-solution-plan.md), [connectors-and-storage.md](./docs/planning/connectors-and-storage.md).

Planning informs **issue writing and design**; **issues drive implementation**. Tag issues with stage when useful (e.g. `stage:1`). Do not implement post-MVP stages without a prioritized issue.

## Environment variables

All configuration and secrets belong in env files at the repo root:

| File | Committed | Purpose |
|------|-----------|---------|
| `ex.env` | Yes | Template — every variable name, safe defaults, and comments. **Add new keys here first.** |
| `local.env` | No (gitignored) | Developer machine — real API keys, DB URLs, and secrets. Never commit. |

**Rules for agents and contributors:**

- Do **not** introduce `.env`, `.env.local`, or ad-hoc env files elsewhere; use `ex.env` / `local.env` only.
- When a feature needs a new setting, add the key to **`ex.env`** (with a placeholder or empty value) and document it; developers copy or update **`local.env`** locally.
- Application code should load **`local.env`** in development (and production should use the host’s secret manager or injected env — not `local.env`).
- Never commit secrets, tokens, or real user data in `ex.env`; keep placeholders only.

Setup for a new clone:

```bash
cp ex.env local.env
# Edit local.env with real values
```

`.gitignore` ignores `*.env` and explicitly allows `ex.env`.

## Development workflow

Follow this process on every change. Full issue/PR examples live in [`docs/planning/issue-pr-workflow.md`](./docs/planning/issue-pr-workflow.md).

### Issue hierarchy: release parent + sub-issues

Each **minor version** (e.g. `v0.1.0`) gets one **top-level parent issue**. All implementation work for that release is tracked as **sub-issues** under that parent.

| Level | Purpose | PRs? |
|-------|---------|------|
| **Parent** | Release milestone — scope, release checklist, target tag | Rarely (e.g. `develop` → `main` promotion may `Refs` parent) |
| **Sub-issue** | One deliverable slice — what agents implement day to day | **Yes** — one sub-issue per branch/PR |

**When creating a release batch (with maintainer):**

1. Create the **parent issue** first (title e.g. `Release v0.1.0 — Insurance MVP`).
2. Create **sub-issues** for each part; attach them as **sub-issues** of the parent in GitHub.
3. Prioritize and implement **sub-issues only** — one at a time.
4. PRs use `Closes #N` on the **sub-issue**, not the parent (unless the PR is explicitly release-wide).
5. When all sub-issues are done, we test on `develop`, promote to `main`, tag (e.g. `v0.1.0`), record the tag on the **parent**, then **close the parent**.

Do not file flat issues for release work without a parent when that work belongs to a planned minor version.

Templates: [`docs/planning/issue-pr-workflow.md`](./docs/planning/issue-pr-workflow.md) — **Release parent issue** and **Sub-issue**.

### One issue at a time

- Work **exactly one sub-issue** per branch and PR — no drive-by fixes or bundled unrelated work.
- The **sub-issue is the execution contract**; planning docs do not override its acceptance criteria.
- Issues are created and **prioritized with the maintainer**; do not invent priority or pull in lower-priority work without agreement.
- If scope grows, **split a new sub-issue** under the same release parent instead of expanding the current one.

### Branches and merges

| Branch | Use |
|--------|-----|
| `develop` | Default target for all PRs |
| `main` | Minor releases only, after joint testing on `develop` |

- **Branch name:** `issue-<number>-<very-short-description>` (e.g. `issue-3-planning-docs`).
- **PR title:** `Issue-<number> - <slightly longer description>` (e.g. `Issue-3 - Add planning documents and workflow`).
- Open PRs **into `develop`**, never into `main`, unless explicitly instructed for a release promotion PR.
- **Minor releases:** merge `develop` → `main` only after we have tested the release batch together on `develop`.
- **Tag** each release commit on `main` (e.g. `v0.2.0`) and record that tag on the related issue(s).

### Issue ↔ PR linking (required)

**On every PR:**

1. Include `Closes #N` (or `Refs #N` if the issue stays open) in the PR body.
2. Add a full issue URL in the PR body (see template below).
3. **Backlink:** comment on the issue with the PR URL when the PR is opened or updated.

**On every completed release:**

- Add the **release tag** (e.g. `v0.2.0`) to the **parent release issue** Links section.
- Close the **parent** when all sub-issues are closed and acceptance criteria for the minor version are met.
- Sub-issues already closed via their own PRs; do not re-close them at release time.

### Pull request description format

Keep PR descriptions **short but descriptive**:

1. **Summary (required)** — One or two sentences stating the **problem** being solved.
2. **Summary (required)** — One or two sentences on **what we are doing** to solve it (broad approach, not every file).
3. **Changes (required)** — Bullet list of concrete changes (what shipped).
4. **Test plan** — Checklist for how it was or should be verified.
5. **Issue** — Full link to the GitHub issue.

Do not write novel-length PR bodies. Do not omit the issue link.

### No AI / editor branding in project artifacts

This project may be built with LLM assistance; that is fine. **Do not advertise tools in issues, PRs, commits, or comments.**

**Never add** lines such as:

- “Made with Cursor” / “Generated by Cursor”
- “Co-authored-by” trailers or badges for Cursor, Copilot, Claude, ChatGPT, etc.
- “AI-assisted”, “written by AI”, or similar disclaimers in PR descriptions, issue bodies, or release notes
- Footer boilerplate promoting any coding agent or IDE

Write issues and PRs as **normal engineering artifacts**: problem, approach, changes, test plan — nothing about which tool drafted the text.

If the user or maintainer wants to mention tooling elsewhere (e.g. personal blog), that is their call — **not** in this repo’s GitHub surface by default.

**PR body skeleton:**

```markdown
Closes #N

## Summary

<Problem in 1–2 sentences.>

<Approach in 1–2 sentences.>

## Changes

- ...
- ...

## Test plan

- [ ] ...

## Issue

- https://github.com/Bmsandoval/covered/issues/N
```

### GitHub issue format

When creating or drafting issues, use the structure in `docs/planning/issue-pr-workflow.md`:

- **Release parent** — minor version milestone + sub-issue checklist + release criteria.
- **Sub-issue** — **Problem**, **Goal**, **Acceptance criteria**, **Out of scope**, **Links** (parent `#N`, PR, optional `stage:` in body).

### Agent process (checklist)

Before coding:

- [ ] Confirm the active **sub-issue** number and that it is the current priority
- [ ] Note the **parent release issue** for context (do not implement the whole parent in one PR)
- [ ] Branch from latest `develop`: `issue-<number>-<very-short-description>` (e.g. `issue-12-pdf-upload`)
- [ ] PR title: `Issue-<number> - <slightly longer description>` (e.g. `Issue-12 - Add PDF upload endpoint`)

Before opening a PR:

- [ ] Changes map only to that issue
- [ ] PR targets `develop`
- [ ] PR body follows the format above with `Closes #N` and issue URL
- [ ] Issue commented with PR link

Before considering work “released”:

- [ ] Tested on `develop` with maintainer per release batch
- [ ] Release tagged on `main`; tag noted on issue

## Repository conventions

- Keep changes scoped to the task; match existing patterns once code lands.
- Do not commit secrets; use `local.env` for local secrets (see **Environment variables** above).
- Prefer clear module boundaries: `ingestion`, `normalization`, `retrieval`, `chat` (names may evolve with stack).

When unsure whether a feature fits MVP, ask:

1. **Does this help interpret uploaded insurance documents with citations**, without predicting final bills or requiring insurer login?
2. **If we built it generically, would it still help the post-MVP vault?** Prefer designs that satisfy both when cost is similar.

If neither applies, defer it or put it in `docs/planning/` for post-MVP.
