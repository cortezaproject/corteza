---
kind: file
covers: App.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/index.js
  - lib/vue
touched-by:
  - client/web/unify/src/main.js
tests:
  - client/web/unify/e2e/sections/admin/application-custom.spec.ts
---

# App.vue — shell chrome

## Intention

Render the constant chrome around whatever section is active, and adapt it purely
from the section's declaration — App.vue must know no section specifics.

## UX capabilities

- Topbar (app menu, search when discovery enabled, profile menu incl. section-contributed items, theme switch).
- Branding follows the theme: the loader logo and the chrome read `useBrandLogo` against the app's current theme; the tab icon and desktop-notification icon read it against the OS colour scheme (`useOsColorScheme`), because the browser draws those on its own chrome.
- Per-section left sidebar: only when the section declares one; expand state remembered per section in localStorage, first visit uses `sidebarExpandedByDefault`; routes opt out via `meta.hideSidebar`.
- Global overlays: app-list sidebar, notifications, agent sidebar, workflow prompts, permissions dialog, toasts, confirm dialog.
- Disabled-app gate: section whose registry application is disabled shows the "app disabled" screen instead of content. One exception: a switched-off custom application stays open, as a preview, to whoever may change its page (`canManageSourceOnApplication`) — the rule is `previewsSwitchedOffApp` in `sections/app/preview.js`.

## Data touched

- Stores: applications, notifications, rightSidebar, workflowPrompts, rbac, users, namespaces (all from lib/vue).
- `$Settings` (`ui.topbar`, `discovery.enabled`), `$Auth`, websocket, event bus.
- Provides `$appIconMap`; agent context = route context + active section (+ section `agentContext` enrichment).

## When changing this

- Section adaptation goes through the section object (`sidebar`, `topbar`, `profileItems`, `agentContext`) — never branch on section ids here.
- `currentWebapp` (route.meta.section) drives workflow-prompt visibility parity with the old standalone apps.
