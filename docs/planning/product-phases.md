# Product phases

How **prototype**, **MVP**, and **platform** relate. GitHub **milestones** (`v0.x.0`) track **prototype** thin releases; a future **MVP** line (e.g. `v1.x.0`) is a separate, larger bar.

## Phase overview

| Phase | Version line (planned) | What “done” means |
|-------|------------------------|-------------------|
| **Prototype** | `v0.x.0` | Prove the **citation loop** on insurance upload — thin releases, local-friendly, **anonymous sessions**, not production-hardened |
| **MVP** | `v1.x.0` (TBD) | **Shippable** “Ask My Insurance Plan” — security, reliability, UX, ops, and scope well beyond prototype |
| **Platform** | Post-MVP issues | Any-format Q&A, connectors, cross-corpus search, opt-in vault (stages 6–9) |

**Now:** **Prototype** only. Do not scope MVP work into `v0.x` issues unless the maintainer explicitly promotes a slice.

## Prototype (current)

**Goal:** Learn and demo — upload insurance docs → chunks → (optional structure) → **cited Q&A** with honest gaps.

**Stack (agreed):** **Go** for application code (API and services).

**Auth (agreed):** **Anonymous sessions** — browser/session cookie, no user accounts, no OAuth in prototype. Isolate data by session id server-side.

**Delivery:** **Option B — thin releases** (one user-visible slice per minor version):

| Release | Theme | Stages (see [staged-solution-plan.md](./staged-solution-plan.md)) |
|---------|--------|---------------------------------------------------------------------|
| **v0.1.0** | Go app + upload PDF → chunks / doc list | 0 (app skeleton), 1 (document core) |
| **v0.2.0** | Cited Q&A API on uploaded docs | 3 |
| **v0.3.0** | Plan summary + browser shell | 2, 4 |
| **v0.4.0** | Prototype hardening on real SBCs | 5 |

Each release = one **parent issue** + sub-issues; squash-merge to `develop`; tag on `release-0-x-0`.

**Out of prototype:** production hosting SLAs, full accounts, Drive connectors, arbitrary formats, vault, insurer login, bill prediction.

## MVP (future)

MVP is **not** the current milestone name. Treat it as the **next major product gate** after prototype proves the pattern.

Likely includes (exact list via issues later):

- Accounts or stronger identity than anon session
- Production deployment, monitoring, backup story
- Stronger security / privacy review for real user data
- Polished UX, onboarding, error handling at scale
- Legal/compliance copy and processes appropriate for public use
- Performance and cost controls for LLM/OCR

Prototype code should stay **vault-ready** (`DocumentSource`, provenance, cite-or-refuse) so MVP builds on it rather than replacing it.

## Platform (north star)

See [product-vision.md](./product-vision.md) and stages 6–9 in the staged plan — connectors, cross-corpus search, opt-in vault. Issue-driven only.

## GitHub visibility

| Surface | Convention |
|---------|------------|
| **Repo About** | Lead with phase: `Prototype · Citation-backed Q&A over your insurance documents` (optional: `· Go`) |
| **Topics** | e.g. `prototype`, `golang`, `insurance`, `citations` |
| **Milestone `v0.x.0`** | Description prefix: `Prototype:` + slice name |
| **Parent release issue** | Title: `Release v0.1.0 — Prototype: <theme>` |
| **Label `prototype`** | Optional on v0.x work (alongside `stage:N`) |

Do not label prototype work as `MVP` in titles, milestones, or marketing copy.
