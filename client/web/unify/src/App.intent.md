---
kind: file
covers: App.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/index.js
  - lib/vue
touched-by:
  - client/web/unify/src/main.js
tests: []
---

# App.vue — shell chrome

## Intention

Render the constant chrome around whatever section is active, and adapt it purely
from the section's declaration — App.vue must know no section specifics.

## UX capabilities

- Topbar (app menu, search when discovery enabled, profile menu incl. section-contributed items, theme switch).
- Per-section left sidebar: only when the section declares one; expand state remembered per section in localStorage, first visit uses `sidebarExpandedByDefault`; routes opt out via `meta.hideSidebar`.
- Global overlays: app-list sidebar, notifications, agent sidebar, workflow prompts, permissions dialog, toasts, confirm dialog.
- Disabled-app gate: section whose registry application is disabled shows the "app disabled" screen instead of content.

## Data touched

- Stores: applications, notifications, rightSidebar, workflowPrompts, rbac, users, namespaces, modules (all from lib/vue).
- `$Settings` (`ui.topbar`, `discovery.enabled`), `$Auth`, websocket, event bus.
- Provides `$appIconMap`; agent context = route context + active section (+ section `agentContext` enrichment).

## When changing this

- Section adaptation goes through the section object (`sidebar`, `topbar`, `profileItems`, `agentContext`) — never branch on section ids here.
- `currentWebapp` (route.meta.section) drives workflow-prompt visibility parity with the old standalone apps.
