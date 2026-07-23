---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Email Index view

## Intention

Configure the platform's outgoing SMTP server and verify the configuration
actually works before the platform depends on it for invites, resets, and
notifications.

## UX capabilities

- Edit host:port, username/password, from address, TLS server name, and
  TLS-insecure toggle on a single panel; explicit Save persists.
- Test button runs the SMTP configuration checker against the values
  currently in the form (not the saved ones) and reports success or
  per-field warnings.

## Routes

- `system.email` → `/system/email`; single page, links nowhere else.

## When changing this

- Settings contract is the `smtp.servers` key holding an array of server
  objects; only the first entry is read and written — keep that shape.
- The checker must keep testing live form values, and the `from` address
  doubles as the test recipient.
- The password field is plain settings-backed (round-trips on load); avoid
  logging or echoing it anywhere new.
