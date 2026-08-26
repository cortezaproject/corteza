---
kind: folder
covers: '.'
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Compose section

## Intention

The low-code Compose app (namespaces, modules, pages, charts, records) mounted
as a section of the unified webapp under `/compose`. Ports the legacy standalone
compose SPA route-for-route while delegating shell concerns (auth, topbar,
realtime socket, shared stores) to the unify shell.

## Section contract (index.js)

- `id: 'compose'`, all routes prefixed `/compose`, every route gets
  `meta.section: 'compose'`.
- Layout route `/compose` renders `ComposeHost.vue` and is intentionally
  unnamed; the empty-path child carries the name `compose` and redirects to
  `namespace.list` (Vue Router cannot render an unnamed empty-path child of a
  named parent).
- `routes.js` holds the raw legacy route tree (names unchanged: `namespace.*`,
  `pages`/`page`, `admin.modules|pages|charts.*`, `page.record`); index.js
  prefixes its absolute paths/redirects. Route meta picks the sidebar: the
  namespace list sets `meta.hideSidebar`, the create/edit screens set
  `meta.sidebar: 'namespaces'`.
- `sidebar: ComposeSidebar`, expanded by default.
- `agentContext(ctx)`: enriches AI agent context from route params — resolves
  `slug`→namespace, `pageID`→page, `moduleID`→module, `recordID`→record via the
  shared stores (record checked in both `records` and `labelCache` maps; array
  or object `values` normalized to an object).
- `profileItems(t)`: one "Reminders" topbar profile entry showing
  `reminderStore.activeCount`, toggling the shared right sidebar's `reminders`
  panel. Reminders are a compose feature by design — they do not move to the
  shell.

## Namespace administration

A namespace's modules, pages and charts are configured in its **admin panel** —
the `admin.*` routes and the sidebar section that leads to them. The server has
one op for it, `compose:namespace` `manage` ("Access to namespace admin panel"),
and it is the only thing that opens the panel. Every way in is gated on it: the
sidebar's admin section, the `admin.*` routes themselves, the topbar shortcuts
on public pages, and the onboarding buttons on an empty namespace.

Per-resource rights (`canUpdatePage`, `canUpdateModule`) decide **which** of
those shortcuts is offered, never **whether** the panel is reachable — a user
who may edit a page but holds no `manage` has no way into the builder. The gate
withholds rather than disables: everything behind it refuses on arrival, and an
offer that refuses is worse than one never made.

Hiding the entrances is not the gate. `Namespace/View.vue` guards the routes
themselves, because a bookmark, a pasted link or a namespace switch all arrive
without ever touching the sidebar.

## Map

- `ComposeHost.vue` — section layout: RouterView + compose-only overlays
  (ReminderSidebar, ReminderToastHost, CTranslatorDialog). Fetches reminders on
  mount, forwards `reminder` messages from the shell's `$eventBus` `realtime`
  events to the reminder store, disposes timers on unmount.
- `routes.js` — raw legacy compose route tree (see contract above).
- `sidebar/` — section sidebar (own doc).
- `views/`, `components/` — route targets and their building blocks (own docs).
- `stores/` — compose-specific Pinia stores + shared-store tests (own doc).
- `lib/` — pure helpers: filters, charts, resource translations (own doc).
- `composables/` — page visibility + resource-translation settings (own doc).

## When changing this

- Shared module/namespace/record/user/page stores live in lib/vue and are
  provided/preloaded by the shell — never re-provide them here.
- Routes are frozen: paths and names are a stability promise to bookmarks,
  deep links and integrations. Renames are forbidden unless the old path keeps
  working via a redirect, and all `router.push({ name })` callers move in the
  same change.
- The shell owns the single realtime socket; compose only subscribes via the
  event bus and must unsubscribe on unmount.
