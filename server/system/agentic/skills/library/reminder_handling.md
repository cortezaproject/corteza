---
name: reminder_handling
description: Creating and managing reminders — the payload is a JSON string, not an object, and an unassigned reminder belongs to nobody.
triggers:
  - system_reminder_create
  - system_reminder_update
  - system_reminder_lookup
  - system_reminder_dismiss
  - system_reminder_undismiss
  - system_reminder_snooze
---

# Reminders

A reminder is a note attached to a resource, surfaced to one user at one time.
Both `resource` and `payload` are free-form — whoever renders the reminder
decides what they mean — so **call `system_reminder_lookup` first and mirror an
existing reminder's shape** rather than inventing one. Existing resources on an
instance look like `namespace:<id>` or `compose:record/<ns>/<module>/<record>`.

## payload is a JSON string

`payload` is declared as a string and its contents are stored verbatim. Passing
an object drops it in **silence** — the reminder is created, the response comes
back with `"payload": {}`, and nothing says why.

```
"payload": "{\"title\":\"Call Ada\"}"    ✅ stored
"payload": {"title": "Call Ada"}          ❌ stored as {}
```

On `system_reminder_update` the same argument **replaces** the stored object
rather than merging into it, and an empty string clears it.

## Assign it, or nobody sees it

`assignedTo` is not filled in for you: leave it out and the reminder is stored
with `assignedTo: "0"` and belongs to no one, which is not the same as
belonging to you. Pass the user ID as a string.

Assigning to anyone other than yourself needs the assign-reminder permission
and fails without it. That same permission is what decides how much you can
read: without it a lookup returns only the reminders assigned to you, and with
it you see everyone's. So the same call answers differently for two callers,
and an empty list is not proof that no reminder exists.

## Ending one

- `system_reminder_dismiss` stops it surfacing and keeps the record;
  `system_reminder_undismiss` reverses it.
- `system_reminder_snooze` pushes `remindAt` out and increments `snoozeCount`.
- `system_reminder_delete` is a soft delete and hides the record.

Dismissing is what a person does with a reminder they have dealt with; deleting
is for one that should never have existed.
