---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/agentic/components/AiChat.vue
  - client/web/unify/src/sections/agentic/components/AiTrace.vue
  - client/web/unify/src/sections/agentic/composables/useEditorSplit.js
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

- Config tab: general identity + labels; execution (LLM provider → model — both form-validated required, revalidated on selection change — optional temperature toggle, iteration/context/output/timeout limits); behavior (system prompt — form-validated required — and guardrails); Treaty CL (enable toggle auto-merges default + hardwired articles, grouped per treaty, hardwired ones locked); knowledge bases + system-context injection; tool access (MCP tools from `mcpListTools`, per-tool hints and optional namespace/module allow-rules edited in a dialog that stages a deep copy until Save); TAQ and workflow access lists with per-entry descriptions (names resolved via `$AutomationAPI`); invocation (user chat + sidebar roles, system + service account).
- Exec tab: `AiTrace` for the active conversation; auto-activates on the first trace.
- History tab (edit mode): past `aiConversation`s, row click reopens one as a chat tab (deduped by ID).
- Right: resizable, persisted chat column (`AiChat` + `CConversationTabs`, multi-conversation), placeholder in create mode; hidden on narrow viewports.
- Unsaved-changes guard; delete gated by `canDeleteAgent`; create redirects to the edit route.

## Routes

`agentic.create` at `/agentic/create`, `agentic.edit` at `/agentic/:agentID/edit`; watching the `agentID` param resets conversations, tabs and dialog state before reloading.

## When changing this

- Save serializes with a JSON replacer that strips `_moduleOptions`/`_loadingModules` — any new UI-only helper key stored on the agent must be stripped there too.
- Tool selection is mirrored state (`selectedTools` UI list vs `agent.access.tools` payload) initialized only after both agent and tool catalog load — keep both sides in sync on add/remove.
- TCL enablement merges defaults into `tclArticles`, never removes user picks.
- The provider/model `CFormGroup`s carry a `name` only for resolver error display; their inner selects are severed from PrimeVue form binding (`novalidate`, see the shared-input doc), so the resolver must read those values from `agent`, never from form `values`.
