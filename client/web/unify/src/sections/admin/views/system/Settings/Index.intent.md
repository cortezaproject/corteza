---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Settings Index view

## Intention

Edit all authentication settings on one page: internal auth, password
constraints, MFA, mail, auto-logout, external providers, login background.

## UX capabilities

- Panels of switches/inputs bound directly to flat `auth.*` setting keys;
  one Save persists everything at once.
- External providers table (SAML, dynamic OIDC, standard OAuth) edited in a
  dialog that embeds the `auth/` sub-forms (ExternalOIDC/SAML/Security/Std —
  governed by the folder doc); dialog edits a clone, applied on confirm.
- OIDC providers can be marked deleted/undeleted before saving; background
  image uploads immediately; disabling MFA also clears its "enforced" flag.

## Routes

- `system.settings` → `/system/settings`. No sub-routes; no create/edit split.

## When changing this

- The flat-key mapping (`prepareExternal`/`getExternalChanges`) is the real
  contract: deleting an OIDC provider must write `null` for every key, and
  only changed external keys are submitted (diff vs the loaded snapshot).
- Only `auth.`-prefixed settings are loaded — a panel touching another
  prefix must extend the load. No unsaved-changes guard exists here.
