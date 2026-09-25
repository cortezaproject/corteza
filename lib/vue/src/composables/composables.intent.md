---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores
  - lib/vue/src/utils/internalNav.ts
touched-by:
  - client/web/unify/src
tests:
  - lib/vue/src/composables/useHistoryBack.test.ts
  - lib/vue/src/composables/useMinDuration.test.ts
  - lib/vue/src/composables/useRBAC.test.ts
  - lib/vue/src/composables/useResourceList.test.ts
  - lib/vue/src/composables/useChangedAt.test.ts
  - lib/vue/src/composables/useDraftGuard.test.ts
  - lib/vue/src/composables/useFileUpload.test.ts
  - lib/vue/src/composables/useUnsavedGuard.test.ts
  - lib/vue/src/composables/useUserResolver.test.ts
---

# Shared composables

## Intention

Reusable Composition-API building blocks for behavior every section needs the
same way: list views, permission checks, navigation guards, uploads, theming.
Sections use these instead of re-implementing per-app variants.

## Map

- useAgentRouteContextProvider.ts — default route-snapshot context provider for CAgentSidebar; apps layer richer context on top.
- useConfirmDelete.ts — standardized PrimeVue delete-confirmation dialog (labels, severities).
- useFileUpload.ts — drag-and-drop + FormData upload, either through an API client's axios instance (auth + base URL) or by raw fetch against an endpoint URL + token. Both treat an `error` payload as a failure whatever the HTTP status, because the API reports one under 200; the raw path asks for JSON and falls back to reading a plain-text error's first line.
- useHistoryBack.ts — a Back that always lands somewhere: the previous screen when the history holds one, otherwise a fallback route the caller names.
- useInternalLink.ts — router-aware links: client-side nav for internal shell routes, native nav otherwise.
- useMinDuration.ts — minimum-duration wrapper for async work so spinners don't flash.
- usePermissions.ts — provide/inject context for the app-level permission dialog (`providePermissions` once in App).
- useRBAC.ts — `useRBACStore`: effective permission rules from all APIs; `can(resource, op)` defaults to deny.
- useChangedAt.ts — the one "last change" column every resource list ends with: `changedAtField(header)` for the field, `changedAt`/`changedAtText` for the cell, `resourceState` for the row's lifecycle tag (deleted, then suspended, then archived). The value is the most recent of deletedAt/updatedAt/createdAt, read from either a wire row (strings) or a lib/js model (Dates), and Go's zero time counts as no timestamp.
- useResourceList.ts — full list-view state machine: filter/sort/cursor pagination synced to route query, abortable requests; exposes `filterDefaults` and resets the other state filters when one is set to Only.
- useRightSidebarResize.ts — mouse-drag resize state for the right sidebar (280–800px clamp).
- useTheme.ts — PrimeVue preset construction from theme variables, light/dark handling.
- useDraftGuard.ts — the unsaved-changes guard for a screen editing one resource: owns the baseline, the deep comparison and the busy suppression; `capture()` marks the current state saved, `extra` covers state held beside the draft.
- useUnsavedGuard.ts — dirty-state route-leave + tab-close confirmation; `markSaved()` to bypass after save. The primitive under `useDraftGuard`, used directly when dirtiness is already a flag rather than a comparison.
- useUserResolver.ts — userID → display name via useUserStore (cache-first, single read fallback; caching is best-effort, so a user the store cannot model is still returned to the caller, just uncached).

## When changing this

- `useRBAC.ts` defines a Pinia store despite living here — the shell preloads
  it; `can()` must stay deny-by-default.
- `useResourceList` owns the URL query contract of list views (limit, cursor,
  sort, filters); changing serialization breaks bookmarked lists.
- `changedAt` is the one column key `useResourceList` rewrites: it sends
  `coalesce(deletedAt, updatedAt, createdAt)` so the column orders by the value
  it displays. Sorting such a column on `updatedAt` alone gives every
  never-updated row a NULL key, which the database groups at one end.
- `changedAtText` returns '' for a resource with no timestamp. The date filters
  parse through moment, which reads a missing value as *now* — an unguarded
  call renders today for a record that was never written.
- `usePermissions` throws when no provider exists — keep `providePermissions`
  in the app root.
- `useHistoryBack` reads the router's `state.back`, never `history.length`:
  that counts entries from before the app and never shrinks, so it calls a
  dead end a live one.
- `useDraftGuard.capture()` is the caller's to invoke, and deliberately not
  automatic: a screen's data arrives asynchronously and its field editors
  resolve presets on their own schedule, so any moment the composable picked
  for itself would sometimes land before the form the user was shown and
  report an untouched editor as dirty. Only the screen knows when its draft is
  what the user saw.
- `useDraftGuard` reads `isDirty` as a getter, never a computed: a draft
  mutated through a raw (non-reactive) reference registers no dependency, and a
  cached "clean" would stand while the user's edits went unguarded.
