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

# Settings

## Intention

Authentication settings for the whole platform, edited as one page. The
`auth/` subfolder holds the external-provider form components and remains
covered by this doc.

## Data touched

- `$SystemAPI.settingsList({ prefix: 'auth.' })` / `settingsUpdate`;
  background image via `settingsSetEndpoint` + access token; `$Settings`
  runtime service for applied values.

## Map

- `Index.vue` — the settings page, flat-key mapping owner (see sidecar).
- `auth/ExternalStd.vue` — standard OAuth provider form (dumb `defineModel`).
- `auth/ExternalOIDC.vue` — OIDC provider form (dynamic, deletable).
- `auth/ExternalSAML.vue` — SAML provider form.
- `auth/ExternalSecurity.vue` — role-mapping rules, embedded by the others;
  the one panel still on `defineProps`/`defineEmits` rather than `defineModel`.

## When changing this

- Sub-components are dumb object forms (bound via `defineModel`, except
  ExternalSecurity); all load/save/mapping logic stays in `Index.vue` —
  don't let API calls leak into `auth/`.
