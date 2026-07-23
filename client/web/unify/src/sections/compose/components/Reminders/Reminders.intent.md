---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/stores/reminder.js
touched-by:
  - client/web/unify/src/sections/compose/ComposeHost.vue
tests: []
---

# Reminders

## Intention

The reminder UI for compose: a sidebar to browse/create/edit reminders and
toast popups when a reminder fires. All state and persistence live in the
reminder Pinia store; these components are views over it.

## Data touched

Everything goes through `stores/reminder` (fetch, save, dismiss, delete,
snooze, toast queue). ComposeHost feeds realtime reminder messages from the
websocket into `store.handleRealtimeReminder`. New reminders default
`assignedTo` to the authenticated user and scope `resource` to the current
namespace when one is active.

## Map

- ReminderSidebar.vue — drawer host toggled app-wide; contains ReminderManager
- ReminderManager.vue — switches list ↔ edit based on `store.editing`; seeds new reminders (assignee, resource, payload)
- ReminderList.vue / ReminderEdit.vue — presentation + form; emit intents, never call the API
- ReminderToastHost.vue — fixed-position toasts from `store.toasts`: dismiss, snooze, link to the related record page

## When changing this

Mounted once in ComposeHost — don't mount these per-view. Keep the
list/edit components store-free (props + emits only); the manager is the only
place that binds them to the store. Toast host must survive route changes and
stay above dialogs (z-index) but never intercept clicks outside its cards.
