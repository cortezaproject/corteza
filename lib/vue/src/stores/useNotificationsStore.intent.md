---
kind: file
covers: useNotificationsStore.ts
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/App.vue
  - lib/vue/src/components/notifications/Notifications.vue
tests: []
---

# useNotificationsStore

## Intention

Single notification inbox for the app: list, unread badge counts, and the
receiving end of realtime notification events.

## State owned

`notifications` (`system.Notification` instances), `pageCursor` for
cursor-based load-more, persisted `muted` flag (localStorage).

## API surface consumed

`$SystemAPI.notificationList/MarkAsRead/MarkAsUnread/MarkAllAsRead/MarkAllAsUnread/Delete`.

## Consumers

Unify shell (preload, unread count in document title, websocket dispatch),
notifications sidebar panel.

## Invariants

- `handleRealtime(msg)` applies `notification.*` websocket messages and
  returns true when handled — the shell uses this to short-circuit its own
  dispatch. New realtime notification types belong here.
- `fetchNotifications` appends when a `pageCursor` is set, replaces otherwise;
  mark-as-read mutations patch local state after the API call resolves, never
  optimistically before it, and never by refetching.
