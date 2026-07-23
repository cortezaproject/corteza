---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/lib/eventBus.js
  - client/web/unify/src/sections/workflow/components/ExpressionTable.vue
  - lib/vue
touched-by:
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
tests: []
---

# Step configurators

## Intention

Per-step-kind configuration panels shown in the editor's right sidebar, plus the
workflow-level config dialog. One panel per step kind; `index.vue` dispatches on
the selected item's `config.kind`/`ref`.

## Map

- `index.vue` — dispatcher: `exec-workflow`→ExecWorkflow, `error-handler`→ErrorHandler, `visual`+`content`→Content, else capitalized kind; forwards `update-value` and gates `update-default-value` behind `defaultName`.
- `base.vue` — shared props contract: `item` ({node, config, triggers}), `edges` (object map id→{node, config}), `outEdges` count, `isSubworkflow`.
- `loader.js` — named-export registry `index.vue` resolves components from.
- `Workflow.vue` — workflow config dialog: name/handle/description, namespace/module labels, runAs, enabled/subWorkflow/keepSessions, import/export/permissions; emits `save`.
- `Trigger.vue` — resource type, event type, constraints, deferred timestamps; announces changes on the eventBus.
- `Function.vue` — function picker with arguments/results mapping (the largest panel); `ExecWorkflow.vue` extends it with a fixed exec-workflow signature and a workflow selector (searches `subWorkflow: 2`, keyed by handle-or-ID).
- `Expressions.vue`, `Gateway.vue`, `Iterator.vue`, `ErrorHandler.vue`, `Error.vue`, `Prompt.vue`, `Delay.vue`, `Edge.vue`, `Content.vue` — remaining per-kind panels.

## When changing this

- Panels write to `item.config.*` (and Gateway directly to `edges[id].config.expr`); WorkflowEditor's deep watchers sync those into the graph — no explicit save emit exists for config fields.
- Auto-naming contract: `update-default-value` only applies while the step still has `defaultName`; a manual rename must permanently stop it.
- New step kind = panel + `loader.js` export + `index.vue` mapping, plus toolbar/style/codec entries in `lib/`.
