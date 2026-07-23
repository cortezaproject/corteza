---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/workflow/components
  - client/web/unify/src/sections/workflow/composables
tests:
  - client/web/unify/src/sections/workflow/lib/codec.test.js
  - client/web/unify/src/sections/workflow/lib/connectionRules.test.js
---

# Workflow editor lib

## Intention

Pure, component-free logic of the editor. Its core duty is legacy
compatibility: the codec keeps the VueFlow canvas interchangeable with the old
mxGraph editor, so workflows created in either render and save correctly in both.

## Map

- `codec.js` — `decodeWorkflow(workflow, triggers)` → `{nodes, edges}`; `encodeWorkflow(nodes, edges)` → `{steps, paths, triggers}`. Canonical `HANDLE_TO_MX` map ties named handles to mx exit/entry coords (round-trip stable by construction); trigger out-edges are stored in `trigger.meta.visual.edges`, not `paths`; legacy kind `workflow` normalizes to `exec-workflow`; non-handle mx style properties from legacy edges are preserved.
- `connectionRules.js` — shared legality: outbound caps (default 1; iterator/error-handler 2; excl/incl/fork gateways ∞; trigger 1), no inbound to triggers, no outbound from termination, ≤1 inbound to termination, one edge per handle, no self-loops, visuals unconnectable.
- `dry-run.js` — `encodeInput`: resolves handles/IDs typed in the test dialog into full resources (namespace/module/page/record/user/role/application, old\* variants) to build the event-args payload.
- `toolbar.js` — ordered drag-palette definition (kind/ref entries + separators).
- `style.js` — kind/ref → {width, height, icon, style}; `style` is the mxGraph style name persisted server-side.
- `editor-auto-complete.js` — expression-language completion list.
- `eventBus.js` — mitt singleton replacing Vue 2 `$root` events.
- `id.js` — incrementing integer node/edge IDs (mxGraph parity; 1 reserved for root).
- `icon.js` — themed step-icon URLs (dark variants, with exceptions); `string.js`, `constraint.js`, `version.js` — label and docs-URL helpers.

## When changing this

- Codec round-trip stability is the hard invariant: decode→encode→decode must be identity for anchors, geometry, and trigger wiring. Run both test files.
- `connectionRules` is the single source of truth for WorkflowEditor, node components, and tests — never fork the rules into a component.
