---
kind: file
covers: toolAccess.js
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/agentic/views/Editor.vue
  - client/web/unify/src/sections/agentic/components/AgentToolDialog.vue
tests:
  - client/web/unify/src/sections/agentic/toolAccess.test.ts
---

# Agent tool access model

## Intention

One definition of what an agent's tool grants add up to, so the panel that
summarises them and the dialog that edits them cannot disagree about what the
agent will actually be allowed to do.

## Contract

- **Two grant shapes.** `splitGrants` separates named entries (`{name, permission}`)
  from family entries (`{group, maxRisk}`). Everything else reads those two.
- **Deny by default.** A tool no grant reaches resolves to `deny`. Blocked and
  absent are one state: an agent that may not use a tool and one that was never
  given it come to the same thing.
- **A grant naming no mode resolves by risk** — `read` → `always`, anything else
  → `ask` (`defaultModeFor`). This mirrors the runtime's own resolution
  (`withResolvedPermissions`); the editor resolves the same way so an agent is
  shown back as it will run.
- **A family ceiling admits its own risk level and everything below it**
  (`RISK_ORDER`: read &lt; write &lt; destructive). A destructive tool under a write
  ceiling is not covered and stays available to grant by name. A named entry
  overrides the family covering it, and a named `deny` always wins.
- **`DOMAINS` is the subject taxonomy** — eight subjects addressed by area
  (`areaOf`), plus `other` for anything unplaced. A subject holds one subject and
  not two. `HIDDEN_AREAS` drops `system_skill`: a skill is attached to the tool
  that triggers it rather than chosen.
- **Scoping applies to compose resources only.** `canScope` is true for the
  `compose_record_`, `compose_module_` and `compose_namespace_` prefixes; an
  allow entry is ignored elsewhere, and on `automation_taq_exec` /
  `automation_workflow_exec` (`scopeBlocked`) it makes the runtime refuse the
  call outright. `scopesModules` is false for the namespace tools, whose resource
  has no module segment to narrow.
- **`sectionsOf` reports the subjects a tool set falls into**, each with its mode
  tally, leaving out any subject holding nothing the agent has.
- **`summaryOf` returns `{ total, totals, subjects }`** — one row per subject
  carrying both `always` and `ask`, and the same two counts summed. The total is
  summed from the subjects, never from the tool list.
- **`confineTo` drops narrowings** naming a namespace the agent no longer works
  in, and reports how many _tools_ lost one — not how many entries.
- **`hasSettings`** says whether a grant carries anything beyond its mode: a note
  the model reads, or a narrowing.

## Used by

- `views/Editor.vue` — the Access panel's readout (`summaryOf`, `confineTo`,
  `MODE_ICONS`, `MODE_COLOURS`, `SUMMARY_MODES`).
- `components/AgentToolDialog.vue` — row grouping, per-section mode, and the
  per-tool scoping controls.

## When changing this

- `defaultModeFor` mirrors the server's resolution. Change it alone and the
  editor shows an agent back differently from how it will run.
- `DOMAINS` is read by both consumers; a subject added here appears in the dialog
  and the readout together, and a tool whose area is in no domain falls into
  `other` rather than disappearing.
- `summaryOf` totals from the subjects. Counting the whole tool list instead
  includes tools no subject holds — a granted skill — and puts the headline above
  the rows it heads.
- `SUMMARY_MODES` is the readout's column set as well as its vocabulary: adding a
  mode adds a column.
