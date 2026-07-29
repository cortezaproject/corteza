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
  - [ ] E2E backfill (waves ruled 2026-07-24; contract-shape assertions only —
        never WIP content; created data uses `e2e-` name prefix + teardown
        delete):
    - [ ] **Wave 1 — project locked shapes**: list CRUD + create flow (FRIA
          placement), wizard step pipeline + graph→editable-dialog, lifecycle
          routing + publish flip + dashboard view set, members dialog `?new=1`.
    - [ ] **Wave 2 — route integrity sweeps**: frozen compose routes resolve;
          admin sidebar → every screen renders.
    - [ ] **Wave 3 — core CRUD flows**: compose record page create/edit +
          RecordModal URL contract; admin user/role editors; workflow list +
          save round-trip; TAQ list + builder save.
    - [ ] **Wave 4 — agentic + chatbot**: agent list/editor CRUD (exec
          skip-if-no-LLM), chatbot list/editor, sessions inbox.
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
        data-sensitivity). Phasing agreed 2026-07-28: risk scenarios first
        (inherent severity + safeguards, taxonomies as FE config, linked to
        project resources), risk-management/detection rules after — the latter
        needs project-scoped actionlog, so it is gated on that work.
  - [ ] **Create-flow governance removal** (ruled 2026-07-28) — delete the
        deployer-category questions, `DEPLOYER_QUESTIONS`, the `create()`
        mapping, BE `ProjectDeployerCategories`/`FriaRequired`/`friaRequired()`
        and the `project.deployer.*` i18n keys. Lands WITH the Govern FRIA
        flow that replaces the determination, never before it.
  - [ ] **Groups** build-out — part of the intended members model; prior broken
        experiment deleted 2026-07-24, build fresh.
  - [ ] **M&M build-time dashboard** — replace the wizard tab placeholder.
        Intent agreed 2026-07-28 (`project.intent.md`, `Wizard.intent.md`):
        revision-as-milestone, board + metrics + activity for one revision,
        rendered from a manage-nav config rather than STEPS. Blocked on the
        BE change repointing the six work-item resources to the root project
        with a nullable revision ref.
  - [ ] **Publish tab — BE deployment plan** (ruled 2026-07-29). `ProjectChange`
        carries only `{path, risk}`, so added and removed modules are
        indistinguishable and a removed field looks like a retyped one; the
        diff also covers modules/fields ONLY. Add an explicit `op`, widen it
        past compose modules, and return the record count behind each change.
  - [ ] **Publish tab — FE** (ruled 2026-07-29, design in
        `views/Wizard.intent.md`). Fourth tab, stage components under
        `components/wizard/publish/`, mapping editor, cluster removed from the
        tab row. `wizard.spec.ts` asserts three tabs and must be updated.
  - [ ] **Federation envoy filters ignore their scope** — `module_exposed.cue`
        and `shared_module.cue` declare `envoy.scoped: true` but not
        `store.extendedFilterBuilder`, so their generated filter builders drop
        the scope exactly as chart/page_layout did (fixed 2026-07-29). Nothing
        scopes them today, so it is latent, not live. Same fix: the cue flag
        plus an `extend<X>Filter`. `module_field` is deliberately NOT in this
        list — its builder is hand-written and CloneFromStore omits it on
        purpose (fields ride along as nested children of their module).
  - [ ] **Publish approval attribution** — the Publish tab's approval stage has
        no "requested by X, approved by Y, at Z": `governanceByProject` holds
        only `{status, note, values}`, with no actor or timestamp, and resets on
        reload. Deliberately left blank rather than stamped session-locally — an
        approval record that forgets itself is worse than none in a governance
        product. Lands with the governance-persistence work above.
  - [ ] **Publish drops every record** — `publishProject` sends `mappings: []`,
        `migrateRecords` no-ops on an empty set, then the old namespace is
        soft-deleted. Any second-or-later publish silently empties the project.
        Closed by the mapping editor above; until then, do not publish a
        revision that has records.
  - [ ] **Permission create-ops** ruling (Bucket A/B split) + matrix support.
  - [ ] Project-only **server code** still excluded — Phase 4 scoping.

## Recorded product backlog (from doc reviews)

- [x] ~~Project list search is dead~~ — ruled BE-side and fixed 2026-07-24:
      rdbms filter override matches handle OR `meta.short`; verified live,
      DRIFT note resolved, e2e specs use the search box again.
- [x] ~~Dev-stack e2e prerequisite~~ — ruled 2026-07-24:
      `AUTH_REQUEST_RATE_LIMIT=0` set in the dev server env; documented in
      `.env.e2e.example` + `e2e.intent.md`. Full project suite green (14/14).
