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
| **Parent** | Release milestone — scope, release checklist, target tag | **No** — close manually when the release ships |
| **Sub-issue** | One deliverable slice or bug fix — what agents implement day to day | **Yes** — one sub-issue per branch/PR |

**When creating a release batch (with maintainer):**

1. Create the **parent issue** first (title e.g. `Release v0.1.0 — Insurance MVP`).
2. Create **sub-issues** for each part; attach them as **sub-issues** of the parent in GitHub.
3. Prioritize and implement **sub-issues only** — one at a time.
4. **Bug fixes** found during testing → always a **new sub-issue** under the same parent (never bundled into an unrelated sub-issue PR).
5. PRs to `develop` use **`Implements`** with the **sub-issue** full URL only (see below) — squash-merge closes the sub-issue.
6. When all sub-issues are closed and the release batch is tested on `develop`, **cut a release branch**, tag (e.g. `v0.1.0`), **close the parent** when release criteria are met.

Do not file flat issues for release work without a parent when that work belongs to a planned minor version.

**There is no `main` / `master`.** Integration happens on `develop`; shipped minors live on **`release-*`** branches.

Templates: [`docs/planning/issue-pr-workflow.md`](./docs/planning/issue-pr-workflow.md) — **Release parent issue** and **Sub-issue**.

### Labels and milestones (required)

Every **parent**, **sub-issue**, and **PR** for a release must have:

| Item | Milestone | Labels |
|------|-----------|--------|
| **Parent release issue** | Minor version (e.g. `v0.0.0`) | `release` |
| **Sub-issue** | Same milestone as parent | Type + stage (e.g. `documentation`, `planning`, `stage:0`) |
| **Pull request** | Same milestone as the sub-issue it closes | Same type labels as sub-issue (omit `stage:*` on PR if you prefer) |

**Milestone = minor version** — create one per release (`v0.0.0`, `v0.1.0`, …) before filing issues.

**`gh` examples:**

```bash
# Create milestone (once per minor version)
gh api repos/Bmsandoval/covered/milestones -f title="v0.1.0" -f description="Insurance MVP — …"

# Parent issue
gh issue create --title "Release v0.1.0 — …" --milestone "v0.1.0" --label "release"

# Sub-issue
gh issue create --title "…" --milestone "v0.1.0" --label "enhancement,stage:1"

# PR (after open)
gh pr edit <number> --milestone "v0.1.0" --add-label "enhancement"
```

**Stage labels:** `stage:0` (foundation), `stage:1` (document core), `stage:2` (insurance), etc. — see [staged-solution-plan.md](./docs/planning/staged-solution-plan.md).

**Type labels:** use existing repo labels (`documentation`, `enhancement`, `planning`, `bug`, …). Add new labels only when needed.

### One issue at a time

- Work **exactly one sub-issue** per branch and PR — no drive-by fixes or bundled unrelated work.
- The **sub-issue is the execution contract**; planning docs do not override its acceptance criteria.
- Issues are created and **prioritized with the maintainer**; do not invent priority or pull in lower-priority work without agreement.
- If scope grows, **split a new sub-issue** under the same release parent instead of expanding the current one.
- **Bugs found while testing** → **new sub-issue** under the parent; do not patch drive-by on another sub-issue’s branch.

### Branches and merges

**No `main` or `master`.** This repo uses **`develop`** plus **release branches** cut from it.

| Branch | Pattern | Purpose |
|--------|---------|---------|
| `develop` | — | Integration line — all feature/fix work lands here first |
| `release-X-Y-Z` | e.g. `release-0-0-0`, `release-0-1-0` | Shipped minor version; receives hotfixes only after cut |

**Feature / sub-issue PRs → `develop`**

- Branch from `develop`: `issue-<number>-<very-short-description>` (e.g. `issue-3-planning-docs`).
- PR title: `Issue-<number> - <slightly longer description>`.
- **Merge method: squash merge** into `develop`.

**Cutting a release (after sub-issues merged and tested on `develop`):**

```bash
git checkout develop && git pull
git checkout -b release-0-1-0
git push -u origin release-0-1-0
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

Use the minor version in the branch name (`release-0-0-0` for `v0.0.0`) and matching git tag.

**Hotfixes / patches on a shipped release:**

1. Branch from the **release branch** (e.g. `release-0-1-0`), not from `develop`.
2. Open PR **into that release branch**.
3. **Squash merge** into the release branch.
4. **Backmerge** release → `develop` with a **regular merge** (merge commit), not squash — so `develop` picks up the patch.
5. Tag a patch release on the release branch if applicable (e.g. `v0.1.1`).

**Release branch must not stay ahead of `develop`.** Patches on `release-*` must be **backmerged** into `develop` promptly. `develop` is the forward line; release branches are cut points plus patch lines that flow back.

**v0.0.0 dry run:** work **squash-merges to `develop` only** first; cutting `release-0-0-0` from `develop` is a later step in the same dry run (after merge and test).

**PR targets by work type:**

| Work | PR base | Merge into base |
|------|---------|-----------------|
| Sub-issue / feature | `develop` | **Squash** |
| Hotfix on shipped release | `release-X-Y-Z` | **Squash** |
| Backmerge after hotfix | `develop` | **Merge** (regular) |

Do not open feature PRs directly into `release-*` unless it is a hotfix for that release.

- **Link the branch on the sub-issue** so it appears under **Development** (see below).

### Link branches and PRs in Development (required)

GitHub’s **Development** sidebar on the **sub-issue** (not the parent release issue) must show the branch and/or PR.

**Order matters — do not push an unlinked branch first.**

1. `git checkout develop && git pull`
2. `gh issue develop <N> --name issue-<N>-<slug> --checkout --base develop` — **before any commit**
3. Commit, push, open PR
4. Verify Development shows branch and/or PR (see below)

**Never** use GraphQL `createLinkedBranch` on an existing branch — GitHub creates a **new auto-named branch** (e.g. `3-add-planning-…`), not `issue-3-planning-docs`, and the real branch stays unlinked.

**If the branch was already pushed without linking:**

- Sub-issue → **Development** → **Link a branch** → select `issue-<N>-<slug>` on the remote.

`gh issue develop` cannot attach a branch that already exists on the remote (API error). Manual link or a new linked branch name is required.

**Parent release issues (#2, etc.)** do not get branches — only **sub-issues** do.

### Issue ↔ PR linking (required)

**On every PR to `develop` — first line (sub-issue only):**

```text
Implements https://github.com/Bmsandoval/covered/issues/<sub-issue-number>
```

Example (sub-issue #3):

```text
Implements https://github.com/Bmsandoval/covered/issues/3
```

- **Sub-issue** — the slice you implemented (branch `issue-<sub>-…`, Development panel shows branch/PR).
- **Parent** — linked as a **sub-issue** in GitHub; do **not** reference it in the PR body. Close the parent manually when the release ships.

Then Summary, Changes, and Test plan (see template). Do **not** add an **Issues** section — Development on the sub-issue is enough.

**Hotfix PRs** (into `release-*`, not `develop`): `Implements` only the **bug-fix sub-issue** URL unless the maintainer says otherwise.

**On every completed release (after all subs merged to `develop`):**

- Cut `release-X-Y-Z`, tag, record the tag on the parent issue, and **close the parent** when release criteria are met.

### Pull request description format

Keep PR descriptions **short but descriptive**:

1. **Summary (required)** — One or two sentences stating the **problem** being solved.
2. **Summary (required)** — One or two sentences on **what we are doing** to solve it (broad approach, not every file).
3. **Changes (required)** — Bullet list of concrete changes (what shipped).
4. **Test plan** — Checklist for how it was or should be verified.

The **first line** must be `Implements https://github.com/Bmsandoval/covered/issues/<N>`. Do not add a separate Issues section — GitHub **Development** links the PR on the sub-issue.

Do not write novel-length PR bodies.

### No AI / editor branding in project artifacts

This project may be built with LLM assistance; that is fine. **Do not advertise tools in issues, PRs, commits, or comments.**

**Never add** lines such as:

- “Made with Cursor” / “Generated by Cursor” / `Made with [Cursor](https://cursor.com)`
- “Co-authored-by” trailers or badges for Cursor, Copilot, Claude, ChatGPT, etc.
- “AI-assisted”, “written by AI”, or similar disclaimers in PR descriptions, issue bodies, or release notes
- Footer boilerplate promoting any coding agent or IDE

**GitHub may append a Cursor footer when PRs are created from the IDE.** After `gh pr create` or before merge, **read the PR body and remove** any Cursor line. Use:

```bash
gh pr view <number> --json body --jq .body   # inspect
gh pr edit <number> --body-file pr-body.md   # fix and save without footer
```

Write issues and PRs as **normal engineering artifacts**: problem, approach, changes, test plan — nothing about which tool drafted the text.

**PR body skeleton:**

```markdown
Implements https://github.com/Bmsandoval/covered/issues/<sub>

## Summary

<Problem in 1–2 sentences.>

<Approach in 1–2 sentences.>

## Changes

- ...
- ...

## Test plan

- [ ] ...
```

### GitHub issue format

When creating or drafting issues, use the structure in `docs/planning/issue-pr-workflow.md`:

- **Release parent** — minor version milestone + sub-issue checklist + release criteria (no parent/child/planning/PR links in the body — use GitHub **sub-issues**).
- **Sub-issue** — **Problem**, **Goal**, **Acceptance criteria**, **Out of scope** only. Link to parent via GitHub sub-issues; branch/PR appear under **Development**.

### Agent process (checklist)

Before coding:

- [ ] Confirm the active **sub-issue** number and that it is the current priority
- [ ] Sub-issue has correct **milestone** (minor version) and **labels**
- [ ] Note the **parent release issue** for context (do not implement the whole parent in one PR)
- [ ] Create linked branch: `gh issue develop <N> --name issue-<N>-<slug> --checkout --base develop`
- [ ] PR title: `Issue-<number> - <slightly longer description>` (e.g. `Issue-12 - Add PDF upload endpoint`)

Before opening a PR:

- [ ] Branch appears under the sub-issue **Development** section (via `gh issue develop` or manual link)
- [ ] Changes map only to that issue
- [ ] PR targets `develop`
- [ ] First line: `Implements https://github.com/Bmsandoval/covered/issues/<sub>`
- [ ] PR body has **no** Cursor / “Made with” footer (`gh pr view` to verify)
- [ ] PR has same **milestone** as sub-issue and matching **labels** (`gh pr edit …`)
- [ ] Issue commented with PR link

Before considering a minor version “released”:

- [ ] All sub-issues squash-merged to `develop`
- [ ] Tested on `develop` with maintainer
- [ ] `release-X-Y-Z` cut from `develop` and pushed
- [ ] Tag (e.g. `v0.1.0`) on **release branch**; tag noted on parent issue
- [ ] Any release-branch hotfixes backmerged to `develop` (regular merge)

## Repository conventions

- Keep changes scoped to the task; match existing patterns once code lands.
- Do not commit secrets; use `local.env` for local secrets (see **Environment variables** above).
- Prefer clear module boundaries: `ingestion`, `normalization`, `retrieval`, `chat` (names may evolve with stack).

When unsure whether a feature fits MVP, ask:

1. **Does this help interpret uploaded insurance documents with citations**, without predicting final bills or requiring insurer login?
2. **If we built it generically, would it still help the post-MVP vault?** Prefer designs that satisfy both when cost is similar.

If neither applies, defer it or put it in `docs/planning/` for post-MVP.
