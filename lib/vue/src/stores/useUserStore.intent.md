---
kind: file
covers: useUserStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/App.vue
  - lib/vue/src/composables/useUserResolver.ts
tests:
  - lib/vue/src/composables/useUserResolver.test.ts
---

# useUserStore

## Intention

App-wide user cache so any viewer/label can turn a userID into a display name
without per-component fetching.

## State owned

`set` — `system.User` instances, upserted by userID; `pending` flag.

## API surface consumed

`$SystemAPI.userList`.

## Consumers

Unify shell preloads the first 500 users on mount (`load({ limit: 500 })`);
`useUserResolver`, user field viewers/editors, admin user views.

## Invariants

- `findByID` is a computed whose VALUE IS the finder function — call
  `userStore.findByID(id)` directly, never `.value` first.
- `resolveUsers(ids)` is dedup-aware: fetches only IDs missing from the cache,
  in one batched `userList` call; no-op when everything is cached.
- `storeUsers(users)` caches pre-fetched users (upsert) without an API call —
  use it when another fetch already returned user objects.
