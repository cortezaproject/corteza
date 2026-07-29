# The Intent System

**Status: active since 2026-07-23 (phases 0–3 landed: tooling + all FE packages).
This file is the constitution of the intent system. Every rule the tooling
enforces and every doc the backfill produces derives from here.**

## 1. What this is

Every non-generated, actively-owned part of the repo carries a recorded **intention**:
what it is for, what it touches, what UX/API contract it upholds, and what should be
tested when it changes. Code is compared against intent — during development, in review,
and mechanically in CI. When they disagree, one of them is wrong and the disagreement is
surfaced, never silently ignored.

Development is **intent-first**: for new features the intent docs are written/updated
before implementation, then code is brought into agreement with them.

## 2. File taxonomy

| Artifact | Location | Role |
|---|---|---|
| `<folder>.intent.md` | one per covered folder, named after it (`router/router.intent.md`) | Folder-level intent: purpose, contents map, data touched, contracts |
| `<name>.intent.md` | sidecar next to a load-bearing file (`RunWorkflow.vue` → `RunWorkflow.intent.md`) | File-level intent for files with real contracts |
| `INTENT.md` (repo root) | `/` | Monorepo map: apps, libs, server components, how they relate. The only doc named `INTENT.md` (root dir name varies per clone) |
| `AGENTS.md` (repo root) | `/` | Thin pointer for AGENTS.md-standard tools: describes the intent system and directs any agent to the governing intent docs |
| `.intent/SPEC.md` | this file | The system's own rules |
| `.intent/TODO.md` | outstanding-work checklist | What the system still needs; updated in the same commit as the resolving work |
| `.intent/config.mjs` | coverage manifest | Covered roots, exclusions, file-tier rules, enforced rollout roots |
| `.intent/hooks/` | git + Claude hook scripts | `pre-commit`, Claude `PostToolUse`/`Stop` handlers |
| `.intent/intent.lock.json` | lockfile | Machine-managed: `path → { hash, docs[] }` at last sync |
| `.intent/intent.mjs` | CLI tool | `check` / `sync` / `status` / `coverage` / `affected` |
| `.intent/templates/` | doc templates | Per-file-type skeletons the backfill and new docs must follow |

Naming is fixed: `*.intent.md`, always carrying the name of what it covers — folder
docs are `<foldername>.intent.md` inside the folder, file docs `<filename>.intent.md`
beside the file. The suffix signals — to humans and to AI — a maintained contract,
not notes; the name prefix keeps docs distinguishable in editor tabs and search.
Corner case: a file named like its parent folder (`User/User.vue`) cannot have a
separate sidecar — the folder doc governs it.

### Relation to existing standards (checked 2026-07)
- **AGENTS.md** (Linux Foundation / Agentic AI Foundation): the open standard for
  nested per-directory *agent instructions*. Our folder tier mirrors its nesting
  pattern, but INTENT.md semantics are stricter — a verifiable contract, not
  instructions — so the name stays distinct. A thin root `AGENTS.md` bridges to
  standard-compliant tools.
- **Spec-driven development** (Spec Kit, Kiro, GSD, …): per-feature, ephemeral
  spec→plan→tasks workflows. Complementary: features may be planned with SDD tooling,
  but the durable outcome must land in the touched intent docs.
- **Drift enforcement**: no open standard exists (commercial only: Swimm, Mintlify);
  the `.intent` CLI is deliberately custom.

## 3. Coverage

### Covered (backfill, "everything up front")
- `client/web/unify` — full
- `client/web/chatbot-widget` — full
- `lib/vue`, `lib/js`, `lib/test-utils`, `lib/eslint-client` — full
- `server` — all non-generated, non-vendor packages (~197), folder(=package)-level
- `def/` — cue sources (hand-written, drives codegen: high-value intent)
- `tests/`, `locale/` — folder-level overview docs
- Root-level config (`Makefile`, workspace files) — covered by root `INTENT.md`

### Excluded permanently (never documented, never checked)
- `node_modules/`, `vendor/`, `dist/`, build output
- Generated files: `*.gen.go`, `*.gen.ts`, anything codegen-owned
- Lockfiles, assets (fonts, icons)

### Excluded for now (documented when it stabilizes)
- ~~`client/web/unify/src/sections/project`~~ — covered since 2026-07-24 via
  interview-first documentation: locked contracts + `> **WIP:**` zones recorded
  in `project/project.intent.md`.
- Server-side code that exists **only** for the projects feature and cannot work
  standalone (project resource packages, project-only endpoints). Rule of thumb:
  if deleting projects would delete the file, it is excluded; if projects merely
  *uses* it (TAQ, automation, compose), it is covered. Revisit in Phase 4 scoping.
- `extra/server-discovery` — separate service, documented on first touch.

Coverage is declared in `.intent/config.mjs`, not implied. `intent coverage` reports
covered-without-doc and doc-without-covered-path (orphans) — both are CI failures.

### File tier (which files get their own `*.intent.md`)
Automatic (via `fileTier` globs in config): **route-target views**, Pinia stores,
registries (`field/registry.ts`), `App.vue`, server files defining a service or
RBAC surface.
View = one screen = one UX contract: every view a route points at carries its own
sidecar (Intention, UX capabilities, Routes, tests — tests attach per screen).
Non-route sub-components under `views/` are governed by their resource folder doc,
which stays thin: shared data-touched notes + map. The fileTier view globs must
match exactly the filenames referenced as route components — never blanket
`views/**/*.vue` (pilot finding: that forced empty docs onto sub-components).
Everything else is covered by the nearest folder intent doc; a sidecar always
wins over folder docs when present, so any file can be promoted ad hoc.

## 4. Doc format

Every intent doc = YAML frontmatter (machine-readable) + prose sections (human/AI-readable).

```yaml
---
kind: folder | file            # what this doc covers
covers: recursive              # folder docs: "recursive" (default — governs subtree,
                               #   nearest doc wins, like AGENTS.md) or "." (direct
                               #   children only); file docs: the sibling filename
backfilled: true               # present on AI-backfilled docs never human-vetted;
                               #   removed the first time a human-driven change vets it
owner: fe | be | shared
depends-on:                    # best-effort in phase 1, hardened in phase 2
  - lib/vue/src/components/field/registry.ts
touched-by: []                 # known reverse dependencies (consumers)
tests:                         # specs that must pass when covered code changes
  - client/web/unify/src/sections/taq/__tests__/runWorkflow.spec.js
---
```

Prose sections by template (see `.intent/templates/`; abbreviated here):

- **Vue view / section**: Intention · UX capabilities · Data touched (stores, APIs,
  resources) · Routes & navigation · Key interactions · When changing this
- **Vue component**: Intention · Contract (props/emits/slots — the *meaning*, not the
  signatures) · Used by · When changing this
- **Store**: Intention · State owned · API surface consumed · Consumers · Invariants
- **Go package**: Responsibility · Key types & services · RBAC/permissions surface ·
  Store/DAL usage · When changing this
- **Folder overview**: What this area is · Map of contents (one line per child) ·
  Cross-cutting concerns · When changing this

Rules for prose: record **intent and contracts**, never restate what the code plainly
says. A doc that would survive a rewrite of the implementation is correctly written;
a doc that paraphrases the code line-by-line is wrong and will rot.

**Hard length cap: 60 lines of prose** (frontmatter excluded), enforced by
`intent check`. Long docs are unread docs; the cap forces intent over paraphrase.

## 5. Sync & enforcement (the make-or-break layer)

Hashes live in `.intent/intent.lock.json` (machine-managed; docs stay clean prose).

Rollout is gradual: `enforced` roots in `.intent/config.mjs` list the areas where
check/coverage apply. Backfilled areas get added there as each phase lands.

**`intent check`** — for every covered source file under an enforced root, compare
content hash against the lock. A mismatch means the file changed since its governing
doc(s) last acknowledged it. Also validates frontmatter schema, the 60-line cap, and
orphaned lock entries (renames/moves must be re-synced — the CLI reports the old path
so moves never silently orphan docs).
Check passes only when every changed covered file's governing doc was re-synced.
**`intent sync <paths>`** — after updating (or consciously confirming) the governing
doc, re-records hashes. Sync is an *acknowledgment*: the mechanical layer cannot verify
semantics, but it guarantees drift is always deliberate and visible, never accidental.

Enforcement points:
1. **Claude Code hooks** (`.claude/settings.json`):
   - `PostToolUse` on Edit/Write: when a covered file is edited, inject the governing
     intent doc path(s) into context with the instruction to reconcile.
   - `Stop`: run `intent check --changed`; block ending the turn while drift exists.
2. **Pre-commit** (via `core.hooksPath`, installed by `make intent-hooks`, no husky):
   `intent check --staged` — fails the commit on drift. **Opt-in, not part of the
   standard team setup** (decision 2026-07-24): Claude hooks + CI are the
   enforced layers.
3. **CI** (GitHub workflow): `intent check` + `intent coverage` on every PR. Hard fail.

Escape hatch: none in CI. If code and intent genuinely must diverge temporarily, the
doc gets a `> **DRIFT:** …` note and is synced — the divergence is then itself recorded
intent.

WIP zones: areas whose intent is deliberately not locked yet carry a
`> **WIP:** …` note stating what must not be relied on. Locked contracts and WIP
notes may coexist in one doc; reviewers vet only the locked parts. A WIP note is a
promise that the area will be re-interviewed and locked later, tracked in TODO.md.

## 6. Workflows

**New feature (intent-first):** write/extend the intent docs for every area the feature
touches → review the intent (this is the design review) → implement → `intent check`
green → tests from the `tests:` fields of touched docs pass.

**Maintenance/bugfix:** hook surfaces governing docs → fix code → if behavior/contract
changed, update doc; if not, sync-acknowledge.

**Audit (semantic, AI-driven):** `/intent-audit <path>` (skill:
`.claude/skills/intent-audit/`) — read intent docs + code under a path, report
semantic drift the hash layer can't see, propose edits to whichever side is wrong.
Cadence is mandatory, not optional: before every release, and rotating
per-area so the oldest-audited area is always next (`intent status` shows drift-age —
time since each area's last sync/audit — to make rot visible).
The audit is the only defense against rubber-stamp syncing (hashes green, docs stale);
the hash layer cannot detect that failure mode.

**Change workflow (opt-in):** `/intent-task` (`.claude/skills/intent-task/`)
standardizes the orient → confirm → implement → verify loop. It applies only when the
human explicitly invokes the skill — ordinary changes are not required to run through
it (see root `CLAUDE.md`).

**Test selection (phase 2):** `intent affected <changed files>` walks `depends-on` /
`touched-by` / `tests` to output the spec list for selective Playwright/vitest runs.

## 7. Rollout phases

- **Phase 0 — tooling first.** `.intent/` CLI (`check/sync/status/coverage`), config,
  lockfile format, templates, Claude hooks, pre-commit, CI workflow, root `INTENT.md`.
  Nothing else starts until this is merged: backfill must produce conforming docs.
- **Phase 1 — friction pilot: unify infra + admin/home**: router, plugins, utils,
  `App.vue`, sections `home` + `admin`. *Double review checkpoint: (a) format proven
  on real code, adjust templates; (b) friction verdict — if upkeep feels heavy here,
  fix the tooling before any further backfill.*
- **Phase 2 — remaining unify sections**: compose, workflow, taq, chatbot, agentic.
- **Phase 3 — libs**: `lib/vue`, `lib/js`, `lib/test-utils`, `chatbot-widget`.
- **Phase 4 — server** (3–4 sub-stages by component: system → compose → automation →
  pkg/store/other), excluding project-only packages. Folder(=package)-level.
- **Phase 5 — periphery**: `def/`, `tests/`, `locale/`, root `INTENT.md` finalization.
- **Phase 6 — dependency hardening**: densify `depends-on`/`touched-by`/`tests`,
  ship `intent affected`, wire selective test runs into CI.

Each phase: multi-agent backfill within the stage → docs land trusted as-is (per
decision) → user skims before the next phase starts (catching format drift, not
auditing content).

## 8. Decisions log

| Decision | Choice |
|---|---|
| Granularity | Tiered: folder intent doc default + `*.intent.md` for load-bearing files |
| Scope | Everything up front, minus project-POC-only code (FE `sections/project` + server project-only) |
| Enforcement | Claude hooks + pre-commit + CI hash check (lockfile), no escape hatch in CI |
| Placement | Sidecar, next to the code |
| Name | `<foldername>.intent.md` / `<filename>.intent.md` (root-only `INTENT.md`); thin root `AGENTS.md` bridges to the AGENTS.md standard |
| Backfill | Staged by area, user review checkpoint between stages |
| Trust | Backfilled docs trusted as-is, but marked `backfilled: true` until first human-driven vetting |
| Doc size | Hard 60-line prose cap, enforced by `intent check` |
| Rot defenses | Mandatory audit cadence (pre-release + oldest-first rotation), drift-age in `intent status`, move/rename detection |
| Test/dep mapping | Schema from day one, best-effort data; hardened in phase 6 |
| Workflow | Intent-first development; hash sync = mechanical layer, `/intent-audit` = semantic layer |
