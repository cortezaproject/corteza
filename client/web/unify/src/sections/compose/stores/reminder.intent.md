---
kind: file
covers: reminder.js
backfilled: true
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/compose/ComposeHost.vue
  - client/web/unify/src/sections/compose/index.js
  - client/web/unify/src/sections/compose/components/Reminders
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/RecordListBlock.vue
tests: []
---

# reminder store

## Intention

Single source of truth for the current user's reminders inside compose: the
right-sidebar list, due-reminder toasts, the create/edit form state, and the
topbar profile badge count. Bridges three inputs — REST fetch, realtime
`reminder` messages (forwarded by ComposeHost from the shell socket), and a
local due-time timer — into one consistent list.

## State owned

Sorted reminder list (active by soonest `remindAt` first, then dismissed by
most recent), active toasts, the reminder being edited (null = closed form),
`processing` flag, shown-version map (dedupes toasts), and the single
`setTimeout` handle for the next due reminder.

## API surface consumed

`$SystemAPI` reminderList/Create/Update/Dismiss/Undismiss/Delete/Snooze;
`$Auth.user.userID` scopes the list; lib/vue `useRightSidebarStore` to open
the `reminders` panel on create/edit/save.

## Consumers

ComposeHost (fetch on mount, realtime forwarding, `dispose()` on unmount),
Reminders components (sidebar, toast host, manager), section `profileItems`
(`activeCount` badge), RecordListBlock (per-record "add reminder" row action).

## Invariants

- Mutations go through the API then refetch (`fetchReminders`) — the list is
  never optimistically edited except for delete/realtime upsert.
- A toast re-fires only when the reminder's "version" (remindAt, dismissal,
  snooze count, payload text) changed since last shown; dismissing clears the
  toast and, on undismiss, resets the version so it can fire again.
- Exactly one timer is scheduled, for the earliest pending `remindAt`;
  `dispose()` must be called on section unmount or the timer leaks across
  sections.
- With no API or an anonymous user, `fetchReminders` resolves to an empty
  list instead of throwing.
