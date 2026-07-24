---
kind: file
covers: users.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/stores/events.js
  - client/web/unify/src/sections/project/stores/backlogItems.js
  - client/web/unify/src/sections/project/components
tests: []
---

# project-users store

## Intention

The user directory the project section picks people from (build members,
access users, event owners, backlog assignees) — real system users, loaded
once and resolved by ID everywhere a name/initials is shown.

## State owned

- `users` — normalized directory entries `{ id, name, email }` (name falls
  back through username/email/handle).
- `currentUserID` — from the injected `$Auth` session.
- Resolution helpers: `findUser`, `userName` (falls back to the raw ID),
  `userInitials`.

## API surface consumed

`$SystemAPI.userList` (limit 500, sorted by name). `load()` is
single-flight and cached; `reload()` forces a re-fetch (e.g. after creating
a user) so new users resolve to names instead of bare IDs.

## Consumers

events/backlogItems stores (owner/assignee names), member and user pickers,
dashboard components rendering user names or avatars.

## Invariants

- IDs are strings; `userName` never throws — unknown IDs render as the ID.
- 500-user cap mirrors the platform idiom; a larger directory silently
  truncates (accepted for now).
