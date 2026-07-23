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
  - lib/vue/src/composables/useMinDuration.test.ts
  - lib/vue/src/composables/useRBAC.test.ts
  - lib/vue/src/composables/useResourceList.test.ts
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
- useFileUpload.ts — drag-and-drop + FormData upload via an API client's axios instance (auth + base URL).
- useInternalLink.ts — router-aware links: client-side nav for internal shell routes, native nav otherwise.
- useMinDuration.ts — minimum-duration wrapper for async work so spinners don't flash.
- usePermissions.ts — provide/inject context for the app-level permission dialog (`providePermissions` once in App).
- useRBAC.ts — `useRBACStore`: effective permission rules from all APIs; `can(resource, op)` defaults to deny.
- useResourceList.ts — full list-view state machine: filter/sort/cursor pagination synced to route query, abortable requests.
- useRightSidebarResize.ts — mouse-drag resize state for the right sidebar (280–800px clamp).
- useTheme.ts — PrimeVue preset construction from theme variables, light/dark handling.
- useUnsavedGuard.ts — dirty-state route-leave + tab-close confirmation; `markSaved()` to bypass after save.
- useUserResolver.ts — userID → display name via useUserStore (cache-first, single read fallback).

## When changing this

- `useRBAC.ts` defines a Pinia store despite living here — the shell preloads
  it; `can()` must stay deny-by-default.
- `useResourceList` owns the URL query contract of list views (limit, cursor,
  sort, filters); changing serialization breaks bookmarked lists.
- `usePermissions` throws when no provider exists — keep `providePermissions`
  in the app root.
