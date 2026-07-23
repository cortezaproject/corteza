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
- [ ] **Phase 6 — dependency/test-map hardening**: densify `depends-on` /
      `touched-by`, populate `tests:` as specs get written, implement
      `intent affected` properly (currently a stub) and wire selective test
      runs into CI.

## Tooling

- [ ] **`/intent-audit` skill** (SPEC §6) — the semantic-drift audit is
      designed but not built. Required for the mandatory audit cadence
      (pre-release + oldest-first rotation via `intent status` drift-age).
- [ ] **First CI run**: `intent.yml` has never executed on GitHub — validated
      locally only. Resolves on the next push of this branch (user pushes).

## Event-driven reconciles (intent check will flag these when they happen)

- [ ] When `feat-run-wf-named-io` merges: add the Run Workflow step
      (`CInputWorkflowInputMap.vue`, Workflow/WorkflowSelector registry
      entries) to `taq/components/builder/builder.intent.md`.
- [ ] When `main` merges into this lineage: `lib/vue` field docs must gain the
      PrimeVue `$pcFormField`/`$pcForm` severing contract (fix `928c3bc50`,
      currently absent here — the multivalue-input bug is live on this branch).
- [ ] When the project POC stabilizes: document `sections/project` (+ its
      server code), remove its exclusions from `.intent/config.mjs`, and
      restore `project` to the section list in `sections/sections.intent.md`
      (currently deliberately omitted).

## Recorded product backlog (from doc reviews)

- (none open — Builder staged-deletes and required-fields blocking landed
  2026-07-23)
