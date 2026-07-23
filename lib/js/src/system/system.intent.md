---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/cast.ts
  - lib/js/src/guards.ts
touched-by:
  - lib/vue
  - client/web/unify
tests: []
---

# System resource classes

## Intention

Client-side model layer for the system subsystem: one class per resource
(`User`, `Role`, `UserGroup`, `Application`, `AuthClient`, `Template`,
`Reminder`, `DalConnection`, `Connection`, `Notification`, `Sink`,
`LlmProvider`, `Agent`, `Chatbot`/session types, …) plus system event
constructors (`events.ts`) for the eventbus.

## Contract (what apps may rely on)

- Every class constructs from a partial/raw API payload via `Apply`: IDs cast to strings (`NoID = '0'`), timestamps to `Date`, and unknown keys ignored — apps can pass API responses straight to constructors.
- `User.meta` carries UI-relevant preferences (preferredLanguage, avatar fields, theme, security policy) — the shape is shared with the auth/profile screens.
- Classes are plain data holders + `toJSON`; no API calls happen in here.

## When changing this

- These classes are the deserialization boundary for every system API response in the apps; adding a field means casting it here, not patching consumers.
- Some admin views intentionally use raw API objects instead (SensitivityLevel, Queue, ApiGateway) — don't assume every system resource has a class here.
