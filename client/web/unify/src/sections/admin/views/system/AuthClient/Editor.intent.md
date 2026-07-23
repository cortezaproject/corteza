---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# AuthClient Editor view

## Intention

Create or edit one OAuth2 client: credentials, grant, scopes, validity,
and security constraints — plus developer helpers for client_credentials.

## UX capabilities

- Name required, handle format-validated; redirect URIs edited as a list
  (persisted space-joined); scope checkboxes backed by the space-separated
  `scope` string; validity window; enabled/trusted toggles.
- Security: permitted/prohibited/forced role lists (role objects shown, IDs
  persisted), impersonated user (client_credentials only), default user
  group (new clients get a best-guess group preselected).
- Secret reveal/regenerate (edit only); developer panel (edit +
  client_credentials) with copyable cURL and in-browser token generation;
  delete, permissions, unsaved guard; create redirects to edit.

## Routes

- `system.authClients.create` → `/system/auth-clients/new`; `.edit` →
  `/system/auth-clients/:authClientID`; back lands on `system.authClients`.

## When changing this

- On save `impersonateUser` is forced to NoID unless grant is
  client_credentials with a user picked; switching to client_credentials
  auto-defaults it to the current user — keep both rules in sync.
