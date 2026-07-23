---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/App.vue
  - lib/vue/src/components
tests:
  - client/web/unify/src/sections/compose/stores/module.test.ts
  - client/web/unify/src/sections/compose/stores/record.test.ts
  - client/web/unify/src/sections/compose/stores/namespace.test.ts
---

# Shared stores

## Intention

These Pinia stores are THE app-wide source of truth for cross-section data.
The unify shell instantiates them and preloads the bounded ones once on mount
(applications, notifications, workflow prompts, RBAC, namespaces, first 500
users). Sections and lib components consume these stores and must never
re-create their own copies; per-store contracts live in the sidecar docs.

## Map

- useAgentChatStore.ts — per-agent chat tabs + server-persisted conversation history.
- useAgentStore.js — AI agent list cache.
- useApplicationsStore.js — application list; unify menu visibility + enabled/disabled path checks.
- useAutomationStore.js — TAQ automations list + function/trigger catalog.
- useChartStore.js — compose charts per active namespace, cache-first picker loads.
- useChatbotStore.js — chatbot list cache.
- useModuleStore.js — compose modules, multi-namespace cache with one "active" namespace.
- useNamespaceStore.js — compose namespaces, preloaded by shell.
- useNotificationsStore.ts — notification list, unread counts, realtime message handling.
- usePageLayoutStore.js — compose page layouts per namespace.
- usePageStore.js — compose pages per namespace, tree/reorder/reparent ops.
- useRecordStore.js — compose records + label cache for viewers (see sidecar: findByID vs resolveRecordLabels).
- useRightSidebarStore.ts — which right-sidebar panel is open (single-slot).
- useUserStore.js — user cache; dedup-aware resolveUsers, shell preload.
- useWorkflowPromptsStore.ts — pending workflow user-prompts, resume/cancel.
- useWorkflowStore.js — workflow list cache.

## When changing this

- Stores get API clients via `inject('$SystemAPI'|'$ComposeAPI'|'$AutomationAPI')`
  — they only work inside apps bootstrapped with the API plugins (tests provide
  mocks via Pinia).
- Cache-guarded loads (`loadFor`, preloads) are a contract: pickers rely on not
  refetching per open. Section-local state stays in section store folders.
