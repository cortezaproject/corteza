---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
touched-by: []
tests: []
---

# Email

## Intention

Configure the platform's outgoing email (SMTP) and verify it works before
relying on it for invites, resets, and notifications.

## Data touched

- `$SystemAPI.settingsList` / `settingsUpdate` — the `smtp.servers` setting
  (array of server objects, no resource class);
  `smtpConfigurationCheckerCheck` for validation.

## Map

- `Index.vue` — the one SMTP settings screen (route target, own sidecar).

## When changing this

- Setting keys are the server contract; the checker should always test the
  values currently in the form, not the last-saved ones.
