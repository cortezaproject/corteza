# Intent system — outstanding work

Maintained checklist of what the intent system still needs. Update in the same
commit as the work that resolves an item.

## From the 2026-08-13 full-surface audit

First audit covering all five enforced roots at once (260 docs, 874 files).
Every doc was checked against code; doc-wrong findings were fixed in that
commit. What it left open:

- [ ] **`drift-age` does not respond to an audit.** `intent status` reports the
      OLDEST lock timestamp in a root, and `sync` only re-stamps files that
      actually drifted — so after auditing and syncing all five roots, every
      root still reads `2026-07-23`. The SPEC leans on drift-age for the
      "oldest-audited area is always next" rotation (§6, §8 rot defenses), and
      that signal currently cannot move. Either stamp every covered file on
      sync, or record audit dates per root separately from file hashes.
- [ ] **240 of 260 docs still carry `backfilled: true`.** Per SPEC §4 the flag
      is removed "the first time a human-driven change vets it". An audit is
      not a human-driven change, so the flag stayed. Decide whether an audited
      doc should shed it — otherwise the flag stops distinguishing anything.
- [ ] **Rubber-stamp syncs are real and the hash layer cannot see them.** Two
      project dashboard views (`Overview.vue`, `CategoryView.vue`) were gutted
      into thin panel wrappers on 2026-07-28 and their docs were re-synced with
      the prose untouched — both still described the pre-extraction files. Only
      the audit caught it. This is the failure mode SPEC §6 predicts; it now
      has a confirmed instance.
- [ ] **Consumer/`touched-by` lists rot fastest.** The single largest finding
      class was stale consumer lists and reverse-dependency paths that exist on
      disk but no longer describe a real relationship. Phase 6 dependency
      hardening should treat these as derived, not hand-maintained.
- [ ] **`manageNav.js` is missing dashboard.js's insights group** (Backlog,
      Reports) though `config.intent.md` locks them to the same section set.
      DRIFT note recorded; needs either the sections or a revised ruling.
- [ ] **ModuleImporter aborts a batch import on first failure**, contrary to
      `Modules.intent.md`. DRIFT note recorded; the fix is per-item state.

## Rollout phases

- [ ] **Phase 4 — server backfill.** Scoped 2026-08-04; the three open
      questions are ruled and the backfill itself is the remaining work:
  - **`pkg/` granularity: tiered.** 93 packages, median 3 files — 53 have ≤3,
    28 have 4–10, only 12 exceed 10. One grouped `pkg.intent.md` covers the
    small ones by theme; the 12 substantial packages (`envoy/yaml`,
    `envoy/resource`, `dal`, `envoy/store`, `rbac`, `expr`, …) get their own.
    The 28 mid-sized are judged case by case. Ninety-three docs averaging three
    files each would be boilerplate, and a thin intent doc is worse than none
    because it reads as maintained.
  - **Project-only exclusion: `**/project*.go`, not `**/project_*.go`.** The
    originally proposed glob was wrong — it matches only the snake_case files
    and misses `system/rest/project.go`, `system/service/projectTask.go` and
    nine more hand-written camelCase ones. Verified 2026-08-04: `project*.go`
    catches 81 hand-written files across `system/{service,types,rest,rest/request,rest/handlers}`
    and `store/tests`, with no false positives — every match is project-feature
    code. Generated files need no glob; `**/*.gen.*` already excludes them.
  - **Top-level dirs in scope:** system, compose, automation, federation, pkg,
    store, app, auth, cmd, discovery. Out: devdocs, docs, build, var, assets,
    webapp, webconsole — artifacts or generated output.
  - Still to do: the backfill itself, and adding `server` to `enforced` in
    `config.mjs` once it is covered. Adding it before then would fail `check`
    for every uncovered file.
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

- [x] ~~When `feat-run-wf-named-io` merges: add the Run Workflow step
      (`CInputWorkflowInputMap.vue`, Workflow/WorkflowSelector registry
      entries) to `taq/components/builder/builder.intent.md`~~ — merged and
      recorded 2026-08-13, along with `WorkflowInputMap` in the aggregate-type
      list in `taq/composables/composables.intent.md`.
- [x] ~~When `main` merges into this lineage: `lib/vue` field docs must gain the
      PrimeVue `$pcFormField`/`$pcForm` severing contract~~ — merged and
      recorded 2026-08-13: `CFieldEditor.vue` provides both as `undefined`, and
      `field.intent.md` now documents it with its regression spec.
- [x] ~~Document `sections/project` when the POC stabilizes~~ — done 2026-07-24
      via interview-first (rulings in `project/project.intent.md`). Remaining
      project-section follow-ups:
  - [ ] Approval **backend persistence** (session-local scaffolding in
        `stores/projects.js` `governanceByProject`; BE planned).
  - [ ] **Govern tab + FRIA content** lock-in — re-interview when governance
        content settles (incl. Govern steps summary/resource-management/
        data-sensitivity). Phasing agreed 2026-07-28: risk scenarios first
        (inherent severity + safeguards, taxonomies as FE config, linked to
        project resources), risk-management/detection rules after.
        **The stated gate is resolved (verified 2026-07-30):** project-scoped
        actionlog landed end to end — columns + indexes
        (`fix_2026_07_13/14/28_*` in `upgrade_fixes.go`), store filter
        (`filters.gen.go`), REST (`rest/actionlog.go`) and JS client. Rolling
        threshold counts already work today via `actionlogReport` with
        from/to + `metrics:['count']`. The REAL remaining gates are different:
    - [ ] **Scenario persistence** — `scenario_id` on any rule has nothing to
          point at while scenarios live in `governanceByProject`.
    - [ ] **Resource-TYPE filtering on the actionlog** — the stored `resource`
          column holds the full ref (`corteza::compose:module/42`) and the
          filter does exact equality, so "all Module events" is not
          expressible. This is ALSO a live bug in the existing ActivityPanel
          resource filter, not just a future blocker.
    - [ ] **A condition evaluator** over action meta/delta — none of the
          mockup's predicates (`meta.*`, `delta.*`) have a query surface.
          `actionlog.RegisterListener` and eventbus's `ConstraintMatcher` are
          the two in-repo substrates.
    - [ ] **RBAC + auth project attribution** — `accessControlAction.ToAction()`
          and the auth actions never populate `ResourceProjectID`, so that
          whole registry domain is invisible per project.
    - [ ] **Project-scoped `action-log.read`** — it is a global component op
          today, so a project risk officer cannot read their own registry.
    - [ ] Incident columns for the risk chain: `ai_system_id`, `scenario_id`,
          `rule_id`, `priority`, `evidence`, `dedupe_key`, `source`. All
          additive to `project_incident.cue`; reuse the existing work items
          rather than inventing a parallel incident type.
    - [ ] Two STALE comments claiming the revision filter is unwired —
          `composables/useEventActivity.js` and
          `components/dashboard/ActivityPanel.vue`. It is wired; fix before
          someone designs around a limit that no longer exists.
    - [ ] `scope.ProjectScopeMiddleware` is defined but never mounted, so
          `Action.ProjectID` (`rel_project`) is always 0 for request traffic —
          the working axis is `resourceProjectID`. Either mount it or delete
          the dead path.
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
  - [x] ~~**Publish tab — BE deployment plan**~~ — done since it was written:
        `ProjectChange` (`system/types/project_revision.go:14`) now carries
        `Op`, `Kind`, `Name`, `Module`, `Detail`, `Risk` and a record count,
        and the diff reports non-compose kinds (verified 2026-07-30: a branched
        revision reported an `agent` change).
  - [ ] **Branch copy — remaining resource kinds** (started 2026-07-30). A
        revision branch cloned ONLY the compose namespace, so agents, chatbots,
        automations, connections and roles were silently dropped and the fresh
        draft's own deployment plan reported them all as deletions. Ruled
        2026-07-30: copy all five, connections keep their credentials, roles
        copy WITH their RBAC rules (refs remapped), resolution via an explicit
        old→new ID map. **Agents are done and verified**
        (`system/service/project_revision_clone.go`); still to do:
    - [ ] `configured_connection` and `ng_automation` — neither has a unique
          handle index, so no migration is needed, unlike agents.
    - [ ] `chatbot` — needs a ruling first: `unique_widget_key` is globally
          unique and a widget key is the PUBLIC embed identifier, so a copy
          cannot reuse it. Proposal: mint a fresh key on copy, so existing
          embeds keep hitting the live revision rather than silently following
          a draft. Its `unique_handle` also needs the same per-project
          treatment agents got (`agent.cue` + a `dropIndexes` fix).
    - [ ] **roles + RBAC rules** — last, because rule resource strings point at
          compose IDs that the namespace clone re-minted, so this is the one
          part needing the ID map actually populated from the clone.
    - [ ] Tests for the whole copy: assert every ref on a copied resource
          resolves INSIDE the draft (`ResourceRefs()` is the oracle).
    - [ ] Envoy is deliberately NOT the mechanism — see the long comment at the
          head of `project_revision_clone.go`. Revisit only if system-resource
          scope support is ever generated.
  - [ ] **Publish tab — FE** (ruled 2026-07-29, design in
        `views/Wizard.intent.md`). Fourth tab, stage components under
        `components/wizard/publish/`, mapping editor, cluster removed from the
        tab row. `wizard.spec.ts` asserts three tabs and must be updated.
  - [ ] **Member rows written against a revision** (pre-2026-07-29 data) — member
        scoping now resolves to the chain root, so any row whose `project_id` is
        a revision rather than its root is invisible. Practically none should
        exist (a branched revision had no members, so nobody could manage it
        from there), but a one-line repoint would make it certain.
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
