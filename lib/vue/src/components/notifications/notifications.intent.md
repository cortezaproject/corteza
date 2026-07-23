---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNotificationsStore.ts
  - lib/vue/src/stores/useRightSidebarStore.ts
touched-by: []
tests: []
---

# notifications/

## Intention

In-app notification center: unread badge in the topbar, right-sidebar panel
with unread/all tabs, mute, mark-read, and per-kind rendering. State and
websocket delivery live in `useNotificationsStore`; this folder is the view.

## Map

- `CNotificationButton.vue` — topbar bell + unread badge; toggles the right
  sidebar via `useRightSidebarStore`.
- `CNotificationSidebar.vue` — right-sidebar host.
- `CNotificationsPanel.vue` — titled panel wrapper (with `actions` slot) around
  Notifications.
- `Notifications.vue` — tabs (unread/all), mark-all-read, mute toggle.
- `NotificationList.vue` / `NotificationItem.vue` — list + row; row emits
  `mark-read`/`mark-unread`/`delete` and resolves record links via `$ComposeAPI`.
- `types/` — per-kind bodies (NotificationSimple, NotificationRecord); a new
  notification kind means a new component here plus its dispatch entry.

## When changing this

Barrel exports only Button/Sidebar/Panel/Notifications — inner list/item/types
are internal and safe to reshape. Store owns unread counts and pagination;
components must not fetch on their own except record-link resolution.
