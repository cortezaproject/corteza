# Intent system — outstanding work

Maintained checklist of what the intent system still needs. Update in the same
commit as the work that resolves an item.

## Rollout phases

- [ ] **Phase 4 — server backfill.** Not scoped yet; open questions to settle
      with the user first:
  - `pkg/` granularity: 93 mostly-small packages — tiered (grouped root doc +
    own docs for substantial packages) vs per-package everywhere.
  - Project-only exclusion mechanics: proposed glob `server/**/project_*.go`
    (project files live inside shared packages, e.g. `system/types/project_*.go`);
    verify the glob catches only project-feature files.
  - Which top-level server dirs are in scope (proposal: system, compose,
    automation, federation, pkg, store, app, auth, cmd, discovery; devdocs,
    docs, build, var, assets, webapp, webconsole stay out).
- [ ] **Phase 5 — periphery**: `def/`, `tests/`, `locale/`, `server/codegen`
      templates (hand-written, drives generation), `provision/`.
      `extra/server-discovery` stays on-touch. `lib/eslint-client` has no
      matching source files — decide: cover config file or drop from `covered`.
      Known locale cleanup for when coverage lands: `locale/en/human-webapp/home.yaml`
      keys `assistant.noAgents`, `column.assistant`, `column.notifications` have
      no consumers (found by 2026-07-24 sections/home audit — the agent column
      uses the `agent.sidebar.*` namespace instead).
- [ ] **Phase 6 — dependency/test-map hardening**: densify `depends-on` /
      `touched-by`, populate `tests:` as specs get written. Started 2026-07-24:
      Playwright e2e layer at `client/web/unify/e2e/` (3 seed specs wired into
      `tests:` fields), `intent affected` implemented (governing doc +
      depends-on consumers). Remaining:
  - [ ] E2E backfill: smoke specs per section (agent fan-out) with `tests:`
        wiring; grow beyond smoke as screens stabilize.
  - [x] ~~First e2e run triage~~ — done 2026-07-24: both failures were
        spec-wrong (users.spec strict-mode multi-match assertion; wizard.spec
        false-skip — looked for anchors in a row-click DataTable). Code and
        auth.setup.ts were correct; suite green (4/4) after spec fixes.
  - [ ] CI story for e2e: boot + seed a test instance (currently attach-only,
        excluded from CI).
  - [ ] Wire `intent affected` into CI selective runs once e2e runs in CI.

## Tooling

- [x] ~~`/intent-audit` skill~~ — built 2026-07-24 (`.claude/skills/intent-audit/`),
      alongside `/intent-task` (standard change workflow) and a root CLAUDE.md.
      First real audit run: `sections/home` 2026-07-24 — both docs verified
      accurate, zero drift, no fixes needed.
- [ ] **First CI run**: `intent.yml` has never executed on GitHub — validated
      locally only. Resolves on the next push of this branch (user pushes).

## Event-driven reconciles (intent check will flag these when they happen)

- [ ] When `feat-run-wf-named-io` merges: add the Run Workflow step
      (`CInputWorkflowInputMap.vue`, Workflow/WorkflowSelector registry
      entries) to `taq/components/builder/builder.intent.md`.
- [ ] When `main` merges into this lineage: `lib/vue` field docs must gain the
      PrimeVue `$pcFormField`/`$pcForm` severing contract (fix `928c3bc50`,
      currently absent here — the multivalue-input bug is live on this branch).
- [x] ~~Document `sections/project` when the POC stabilizes~~ — done 2026-07-24
      via interview-first (rulings in `project/project.intent.md`). Remaining
      project-section follow-ups:
  - [ ] Approval **backend persistence** (session-local scaffolding in
        `stores/projects.js` `governanceByProject`; BE planned).
  - [ ] **Govern tab + FRIA content** lock-in — re-interview when governance
        content settles (incl. Govern steps summary/resource-management/
        data-sensitivity).
  - [ ] **Groups** build-out — part of the intended members model; prior broken
        experiment deleted 2026-07-24, build fresh.
  - [ ] **M&M build-time dashboard** — replace the wizard tab placeholder
        (current-version categories + backlog, version-shared with live).
  - [ ] **Permission create-ops** ruling (Bucket A/B split) + matrix support.
  - [ ] Project-only **server code** still excluded — Phase 4 scoping.

## Recorded product backlog (from doc reviews)

- (none open — Builder staged-deletes and required-fields blocking landed
  2026-07-23)
