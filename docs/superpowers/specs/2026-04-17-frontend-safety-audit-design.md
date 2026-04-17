# Frontend Safety Audit — Design

**Date:** 2026-04-17
**Topic:** Audit methodology & deliverable spec for a frontend-wide safety audit
**Scope:** All seven frontend webapps (`admin`, `agentic`, `compose`, `home`, `one`, `taq`, `workflow`) plus shared libs (`lib/js/`, `lib/vue/`)

## Goals

Produce a single prioritized report covering:

1. **Security** — XSS/injection risks, unsafe HTML/URL handling, auth token handling, CSP.
2. **Runtime safety** — null/undefined crashes, unhandled promise rejections, defensive rendering gaps.
3. **Dependency safety** — `pnpm audit` findings (direct vs transitive), outdated critical libs.
4. **Lint-debt shape** — summary-level only, to quantify technical debt without enumerating every warning.

Out of scope: authorization/RBAC logic review, type-safety coverage, accessibility, performance.

## Non-goals

- Fixing anything. The audit surfaces findings; remediation happens in follow-up work.
- Exhaustive enumeration of every lint warning or Medium-severity pattern instance.

## Audit dimensions & tooling

For each dimension: tooling sweep + targeted manual review of hotspots.

| Dimension | Tooling sweep | Manual hotspots |
|---|---|---|
| Security (XSS/injection) | Grep: `v-html`, `innerHTML`, `eval`, `new Function`, `target="_blank"` missing `rel`, `href="javascript:"`, `document.write`, raw `$sanitize` bypasses | Record/field viewers (user-supplied content), markdown renderers, iframe embedders, URL/link builders, auth token storage (`localStorage` vs cookie), login/auth templates, CSP in `nginx.conf` |
| Runtime safety | Grep: unhandled `.then(` without `.catch`, bare `JSON.parse`, `await` without try/catch, optional-chain gaps on known-optional fields | Pinia store injection fallbacks, router guards, API response shape assumptions, `recordStore.findByID` throw sites, `userStore.findByID` computed access, `resolveRecordLabels` callers |
| Dependency safety | `pnpm audit --prod` per workspace + root; `pnpm outdated` for critical libs (vue, vuex/pinia, axios, vue-router, primevue) | Flag transitive-only vulns separately; note upstream-blocked items |

## Severity rubric

- **Critical** — exploitable security (XSS sink reached by user content, auth token leak, known-exploited dep CVE with PoC). Every one enumerated with file:line + fix.
- **High** — likely-exploitable or crash-prone in normal use (unsafe `v-html` on user content; unhandled promise in hot code path; high-severity dep CVE with a fix available). Every one enumerated.
- **Medium** — conditional risk or repeated pattern. Aggregated as one finding per pattern with count + representative examples (3-5 file:line samples). Not fully enumerated.
- **Low / lint debt** — style, unused vars, deprecated APIs, formatting. Summarized only: total warning count per app, top rule categories. No file list.

Rationale: the user has acknowledged significant lint debt; filtering matters more than completeness. Critical + High are the actionable layer; Medium shows systemic patterns; Low quantifies debt shape.

## Execution plan

Four phases:

1. **Tooling sweep.** Run `pnpm audit --prod` for the monorepo and each workspace. Run existing lint configs if present (`pnpm -r lint` or per-app). Run the security + runtime grep pattern list across all seven apps plus `lib/`. Collect raw findings into a working bucket.
2. **Hotspot manual review.** Walk the following surfaces in each app (where present):
   - Auth/token handling (login flows, token refresh, storage layer)
   - Record/field viewers (user content rendered to DOM — `lib/vue/src/components/field/`)
   - Markdown/HTML renderers (anywhere raw HTML is composed or displayed)
   - URL/link builders (anywhere user input feeds an `href` or `window.open`)
   - Router guards (`src/router/`)
   - API response handling (unchecked `.data.response.foo` chains)
   Apply the CLAUDE.md gotchas: `userStore.findByID` is a computed that IS the finder; `recordStore.findByID` throws if module not loaded.
3. **Triage & rank.** Apply severity rubric. Dedupe patterns. Aggregate Mediums. Summarize Lows. Drop false positives.
4. **Write report.** Single doc at `docs/superpowers/audits/2026-04-17-frontend-safety-audit.md`.

## Report structure

```
# Frontend Safety Audit — 2026-04-17

## Executive summary
- Counts per severity, per app
- Top 5 recommended fixes this week

## Critical findings (enumerated)
### C1. <title> — <app>/<file>:<line>
  Risk | Reproduction | Recommended fix

## High findings (enumerated)
(same shape as Critical)

## Medium patterns (aggregated)
### M1. <pattern title>
  Count: N | Examples: <3-5 file:line> | Recommended fix pattern

## Low / lint debt (summary)
- <app>: <count> lint warnings — top categories: <rule IDs>
- ...

## Dependency findings
- Direct-dep CVEs (actionable) — enumerated
- Transitive-only / upstream-blocked — listed separately
- Outdated critical libs — flagged

## Appendix
- Methodology
- Tools & commands run
- Grep patterns used
- Known limitations / things not covered
```

## Output locations

- **This design doc:** `docs/superpowers/specs/2026-04-17-frontend-safety-audit-design.md` (current file).
- **Audit report output:** `docs/superpowers/audits/2026-04-17-frontend-safety-audit.md` (new directory).

## Success criteria

- Report is actionable: a reader can pick the top 5 fixes and start work without further investigation.
- Every Critical and High finding has file:line and a recommended fix.
- Medium patterns include enough examples to verify the pattern exists, not every instance.
- Low summary communicates debt shape in one short table.
- Dependency section distinguishes actionable from upstream-blocked.
- Methodology appendix lets the audit be reproduced later (same greps, same tools).

## Known limitations

- Tooling sweep will miss logic-level flaws that pattern-matching can't detect.
- Manual review is hotspot-only, not exhaustive; low-traffic surfaces may harbor untriaged issues.
- Dep audit reflects advisory DB at audit time; new CVEs can land the next day.
- Authz/RBAC review is deliberately out of scope and would require its own pass.
