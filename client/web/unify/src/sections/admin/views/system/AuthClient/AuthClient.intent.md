---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# AuthClient

## Intention

Manage OAuth2 clients used to authenticate against the platform: their
credentials, grant behavior, and security constraints.

## Data touched

- `$SystemAPI.authClient*` incl. `authClientExposeSecret` and
  `authClientRegenerateSecret`; `roleRead` / `userGroupList` for security
  fields. Class-based resource `system.AuthClient`.
- API calls identify the client as `clientID`, while the resource field and
  route param are `authClientID`.

## Map

- `List.vue` — client list (route target, own sidecar).
- `Editor.vue` — client create/edit incl. security + developer helpers
  (route target, own sidecar).

## When changing this

- The secret is never part of the read payload — it only exists via the
  expose/regenerate endpoints; never persist or log it client-side.
- Security role IDs are resolved best-effort (missing roles are dropped
  silently); keep that tolerant behavior for deleted roles.
