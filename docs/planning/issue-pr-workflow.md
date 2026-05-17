# Issue and PR workflow (reference)

Canonical templates for GitHub issues and pull requests. Agents must follow [AGENTS.md](../../AGENTS.md) **Development workflow**; this doc is the detailed model.

## Branching and releases

**No `main` or `master`.** Only **`develop`** and **`release-*`** branches.

| Branch | Purpose |
|--------|---------|
| `develop` | Integration — all sub-issue PRs **squash-merge** here |
| `release-X-Y-Z` | Shipped minor (e.g. `release-0-0-0` for `v0.0.0`) — cut from `develop`; hotfixes land here |

```mermaid
flowchart TB
  subgraph daily [Day to day]
    F[issue-N branch] -->|squash PR| D[develop]
  end
  subgraph release [Release]
    D -->|cut branch| R[release-0-1-0]
    R -->|tag| T[v0.1.0]
  end
  subgraph hotfix [Hotfix]
    H[hotfix branch] -->|squash PR| R
    R -->|regular merge backmerge| D
  end
```

### Release flow

1. Sub-issues merge to **`develop`** via **squash** PRs.
2. Test together on **`develop`**.
3. Cut **`release-X-Y-Z`** from `develop` (e.g. `release-0-0-0`).
4. Tag on the **release branch** (e.g. `v0.0.0`).
5. Update parent issue with tag; close parent when done.

**v0.0.0 dry run:** squash-merge planning work to `develop` first; cut `release-0-0-0` later in the dry run.

### Hotfix flow

1. Branch from **`release-X-Y-Z`** (not `develop`).
2. **Squash merge** PR into the release branch.
3. **Regular merge** (backmerge) `release-X-Y-Z` → `develop` so patches are not lost.
4. Tag patch on release branch if needed (`v0.1.1`).

**Rule:** A release branch must **not** remain ahead of `develop` — backmerge after every hotfix.

### Merge methods (required)

| PR type | Base branch | Merge method |
|---------|-------------|--------------|
| Sub-issue / feature | `develop` | **Squash** |
| Hotfix | `release-X-Y-Z` | **Squash** |
| Backmerge | `develop` | **Merge** (regular merge commit) |

---

## Issue hierarchy (required for each minor version)

Every planned minor release (e.g. `v0.1.0`) uses:

1. **One parent issue** — the release milestone.
2. **Multiple sub-issues** — the parts we implement and merge via PRs.

```mermaid
flowchart TB
  P["Parent: Release v0.1.0"]
  S1[Sub: Upload API]
  S2[Sub: PDF chunking]
  S3[Sub: Cited Q&A API]
  P --> S1
  P --> S2
  P --> S3
  S1 --> PR1[squash PR → develop]
  S2 --> PR2[squash PR → develop]
  S3 --> PR3[squash PR → develop]
```

**Rules**

- Create the **parent first**, then sub-issues; link sub-issues to the parent using GitHub **sub-issues**.
- Agents implement **sub-issues only** — one sub-issue per branch/PR.
- **Link every working branch** on the sub-issue (**Development** sidebar).
- PRs **implement** the **sub-issue** (`Implements` line in body); **squash merge** into `develop`.
- Close the **parent** after: all sub-issues closed, test on `develop`, `release-X-Y-Z` cut, tag pushed, tag on parent.

### Labels and milestones (required)

| Item | Milestone | Labels |
|------|-----------|--------|
| Parent release issue | Minor version (`v0.0.0`, `v0.1.0`, …) | `release` |
| Sub-issue | Same as parent | e.g. `documentation`, `planning`, `enhancement`, `stage:0` |
| Pull request | Same as sub-issue | Match sub-issue type labels |

Create the **milestone first** (one per minor version):

```bash
gh api repos/Bmsandoval/covered/milestones -f title="v0.1.0" -f description="…"
```

Apply when creating or after filing:

```bash
gh issue edit <parent> --milestone "v0.0.0" --add-label "release"
gh issue edit <sub> --milestone "v0.0.0" --add-label "documentation,planning,stage:0"
gh pr edit <pr> --milestone "v0.0.0" --add-label "documentation,planning"
```

**Stage labels:** `stage:0` … `stage:9` aligned with [staged-solution-plan.md](./staged-solution-plan.md).

### Link branch and PR (Development section)

**Before the first commit** (required):

```bash
git checkout develop && git pull
gh issue develop <issue-number> --name issue-<issue-number>-<short-slug> --checkout --base develop
# commit, push, then open PR
```

**Branch already on remote without link:** sub-issue → **Development** → **Link a branch** → `issue-<N>-<slug>`. Do not use GraphQL `createLinkedBranch` (wrong branch name).

**PR must start with** `Implements https://github.com/Bmsandoval/covered/issues/N`. Link the PR on the sub-issue **Development** panel if it does not appear automatically (GitHub does not treat `Implements` as a closing keyword).

**Parent release issues** do not show branches — link work on **sub-issues** only.

---

## Release parent issue template

**Title:** `Release v0.1.0 — <short milestone name>`

```markdown
## Summary

One paragraph: what this minor version delivers.

## Target tag

`v0.1.0`

## Release branch

`release-0-1-0` (cut from `develop` after sub-issues merge and test)

## Sub-issues

- [ ] #__ — <title>

## Release acceptance criteria

- [ ] All sub-issues closed (squash-merged to `develop`)
- [ ] Tested on `develop`
- [ ] `release-0-1-0` cut from `develop` and pushed
- [ ] Tag `v0.1.0` on release branch and pushed
- [ ] Tag recorded on this issue

## Test plan (release batch)

- [ ] …

## Links

- Planning: `docs/planning/staged-solution-plan.md`
- Tag: _(fill after release)_
```

---

## Sub-issue template

**Title:** Short, imperative — e.g. `Add PDF upload endpoint`

```markdown
## Parent release

Part of **Release v0.1.0** — #<parent>

## Problem

…

## Goal

…

## Acceptance criteria

- [ ] …
- [ ] PR **squash-merged** to `develop`

## Out of scope

- …

## Links

- Parent: https://github.com/Bmsandoval/covered/issues/<parent>
- PR: …
```

---

## Pull request template

**Branch name:** `issue-<number>-<very-short-description>`

**PR title:** `Issue-<number> - <slightly longer description>`

**Base branch:** `develop` (features) or `release-X-Y-Z` (hotfixes only)

**Merge:** **Squash** (features and hotfixes); **regular merge** for backmerge PRs only.

```markdown
Implements https://github.com/Bmsandoval/covered/issues/N

## Summary

<Problem in 1–2 sentences.>

<Approach in 1–2 sentences.>

## Changes

- …

## Test plan

- [ ] …

## Issue

- https://github.com/Bmsandoval/covered/issues/N
```

---

## No tool branding

Do **not** include “Made with Cursor”, “AI-generated”, or similar in issues, PRs, commits, or release notes. See [AGENTS.md](../../AGENTS.md).

---

## Linking checklist

**When creating a minor version**

- [ ] Milestone created (e.g. `v0.0.0`)
- [ ] Parent issue + milestone + `release` label
- [ ] Sub-issues linked as sub-issues of parent; same milestone + labels

**When starting work (sub-issue)**

- [ ] Linked branch on **Development** (`gh issue develop` or manual link)
- [ ] Branch from `develop`

**When opening a PR (sub-issue)**

- [ ] Base: **`develop`**
- [ ] Merge method: **squash**
- [ ] `Implements` line, issue URL, backlink comment; PR linked under **Development**
- [ ] Milestone and labels on PR

**When releasing**

- [ ] All sub-issues squash-merged to `develop`
- [ ] Test on `develop`
- [ ] Cut `release-X-Y-Z` from `develop`
- [ ] Tag on release branch; note on parent; close parent

**When hotfixing a release**

- [ ] PR into `release-X-Y-Z` — squash merge
- [ ] Backmerge `release-X-Y-Z` → `develop` — **regular merge**
- [ ] Release branch not left ahead of `develop`
