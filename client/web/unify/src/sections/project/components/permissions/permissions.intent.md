---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/config/kinds.js
  - lib/vue
touched-by:
  - client/web/unify/src/sections/project/components/wizard/steps/PermissionsStep.vue
  - client/web/unify/src/sections/project/components/datamodel/ModuleDetailDialog.vue
  - client/web/unify/src/sections/project/components/pages/PageDetailDialog.vue
  - client/web/unify/src/sections/project/components/automations/AutomationDetailDialog.vue
  - client/web/unify/src/sections/project/components/agents/AgentDetailDialog.vue
  - client/web/unify/src/sections/project/components/chatbots/ChatbotDetailDialog.vue
  - client/web/unify/src/sections/project/components/roles/RoleDetailDialog.vue
  - client/web/unify/src/sections/project/components/users/UserDetailDialog.vue
tests: []
---

# Project permissions (home of the per-resource permissions contract)

## Intention

One matrix component is the entire project access surface, deliberately scoped
to what an END USER of the deployed app can do at runtime — module records,
page viewing, agent/chatbot use, automation running. Build-time ops stay out;
rare ops remain reachable via the per-row ⚙ full permission editor
(`usePermissions`). The locked contract: a shared `ResourcePermissionsSection`
is embedded in each resource detail dialog; the matrix takes a `scope` prop;
automation appears as a runtime **Run** kind (read + execute); **connection is
deliberately omitted**.

## Data touched

- `useProjectsStore`: resource getters per kind; `loadRoles`;
  `loadEffectiveAccess` / `effectiveAccess` (per-role trace),
  `loadUserEffectiveAccess` / `userEffectiveAccess` (resolved per-user trace);
  `setCapabilityAccess` (writes); `graphVersion` (re-trace trigger); `touch`.
- Corteza RBAC resource strings built inline (compose module/record/field/
  page/page-layout under the project namespace; system agent/chatbot;
  automation ng-automation).

## Map

- `ProjectPermissionMatrix.vue` — three views of one capability tree:
  multi-role step view (roles as columns + create/read/update/delete lens),
  single-role detail view (`hideRoleHeader`, all capability chips per row),
  and single-resource scope (`scope: {kind, id}`) for dialog embedding;
  optional read-only "Evaluated" user column (`evalUserId`). Capabilities are
  compound and honest: a toggle may span several ops (Records read =
  record.read + records.search + record.value.read); cells show `partial`
  when ops disagree, a spinner while tracing. Parent toggles offer a cascade
  confirm ("only this" / "include everything nested"). Pages: View cascades
  to layout reads. Agents/chatbots: single "Use". Automations: single "Run".
- `ResourcePermissionsSection.vue` — the embeddable section every resource
  detail dialog uses: loads the project's roles and renders the scoped matrix
  for one resource.

## When changing this

- App access (namespace read) is intentionally absent — every project member
  already reads the namespace; do not add a misleading toggle.

> **WIP:** **permission create-op gaps** — parent-scoped create/list/
> app-access operations (create modules/pages, list within a parent, app
> access) are missing from the matrix; only record.create exists. The
> namespace-vs-component scoping split awaits a ruling — do not extend the
> KINDS/ops table before it lands.
