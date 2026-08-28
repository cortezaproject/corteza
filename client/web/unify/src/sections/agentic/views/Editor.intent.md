---
kind: file
covers: Editor.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/agentic/components/AiChat.vue
  - client/web/unify/src/sections/agentic/components/AiTrace.vue
  - client/web/unify/src/sections/agentic/components/AgentToolDialog.vue
  - client/web/unify/src/sections/agentic/composables/useEditorSplit.js
  - client/web/unify/src/sections/agentic/toolAccess.js
  - lib/vue/src/stores/useAgentStore.js
touched-by:
  - client/web/unify/src/sections/agentic/index.js
tests: []
---

# Agentic Editor view

## Intention

The agent workbench: configure an agent and converse with it side-by-side,
inspecting every execution decision — so tuning and testing happen in one
screen without saving round-trips elsewhere.

## UX capabilities

- Config tab: general identity + labels; execution (LLM provider → model — both form-validated required, revalidated on selection change — optional temperature toggle, iteration/context/output/timeout limits); behavior (system prompt — form-validated required — and guardrails); Treaty CL (enable toggle auto-merges default + hardwired articles, grouped per treaty, hardwired ones locked); knowledge bases + system-context injection; access (the namespace the agent works in, the tools it holds, and the TAQs and workflows it may set off — the namespace leads because it bounds the rest); TAQ and workflow access lists with per-entry descriptions (names resolved via `$AutomationAPI`); invocation (user chat + sidebar roles, system + service account).
- Exec tab: `AiTrace` for the active conversation; auto-activates on the first trace.
- History tab (edit mode): past `aiConversation`s, row click reopens one as a chat tab (deduped by ID).
- Right: resizable, persisted chat column (`AiChat` + `CConversationTabs`, multi-conversation), placeholder in create mode; hidden on narrow viewports.
- Unsaved-changes guard; delete gated by `canDeleteAgent`; create redirects to the edit route.

## Access panel

- **Works in** — one namespace (`access.allow[0]`), every tool confined to it; empty means whatever the invoker can reach. Changing it runs `confineTo`, which drops tool narrowings naming the old namespace, and a warning says how many lost one.
- **Tools** — the catalog comes from `mcpListTools`; grants are edited in `AgentToolDialog` (a staged `draft`, applied on Apply, never written through) and read back as a readout beside it. Per-tool notes and namespace/module narrowings live in the dialog.
- The readout is a subject × permission cross-tab from `summaryOf`: one row per subject, one column per permission headed by that permission's glyph, and a ruled `All tools` row carrying the whole agent's. The glyphs are the dialog's, and their tooltips are its `mode.*` sentences, so the two surfaces share one notation. A column a subject has nothing in shows a dash. It is guarded on the catalog having loaded — a configured agent must not read as holding nothing while the fetch is in flight, or for good if it fails.
- Access is deny-by-default, so a total of zero is an agent that refuses every question; that state shows the `inheritsHelp` message instead of the table.

## Routes

`agentic.create` at `/agentic/create`, `agentic.edit` at `/agentic/:agentID/edit`; watching the `agentID` param resets conversations, tabs and dialog state before reloading.

## When changing this

- Save serializes with a JSON replacer that strips `_moduleOptions`/`_loadingModules`, the UI-only keys `CInputKnowledgeBase` hangs on a knowledge-base namespace context — any new helper key stored on the agent must be stripped there too.
- `agent.access.tools` is the only tool state; the dialog stages its own copy and emits the whole list back on Apply. What a grant resolves to is `toolAccess.js`'s to say, never recomputed here.
- TCL enablement merges defaults into `tclArticles`, never removes user picks.
- The provider/model `CFormGroup`s carry a `name` only for resolver error display; their inner selects are severed from PrimeVue form binding (`novalidate`, see the shared-input doc), so the resolver must read those values from `agent`, never from form `values`.
