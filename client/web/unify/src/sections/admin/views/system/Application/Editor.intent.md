---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/utils/appIcons.js
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests:
  - client/web/unify/e2e/sections/admin/application-custom.spec.ts
---

# Application Editor view

## Intention

Create or edit one application: its name, description and enabled state plus
how it presents in the unified shell (display name, URL, listed flag, logo).

## UX capabilities

- Name is required; free-text description (`meta.description`, what the list
  shows under the name); enabled toggle; unify panel with display name, URL,
  and "listed" visibility toggle.
- Logo via drag-and-drop upload to the application attachment endpoint;
  preview resolves through `resolveAppLogoUrl` + the built-in icon map;
  only custom logos are clearable — built-in default icons are not.
- Delete, permissions, unsaved guard; create redirects to edit.
- Kind: "a section of Human or a link" (`unify.kind` empty) or "custom app"
  (`'custom'`, its own HTML page). A custom app's URL is shown locked at
  `app/<applicationID>` — the server stores exactly that on every update, and
  creating one saves twice, since only the second save has an ID to put in it.
- Custom app panel (custom kind, edit only): open the app; the namespace, the
  modules its page may read and those it may change, all set here with the
  namespace and module pickers; the page's size and when it was last replaced;
  the page itself.
- What is stored is names — the namespace's slug and module handles — while the
  pickers hold IDs, so a page carried to another instance is read back by what
  it declared rather than by IDs that mean nothing there. A module dropped from
  the read list is dropped from the writable one with it. Whoever holds `source.manage` edits the page here, in
  `CCodeEditor`, and saves it through the same endpoint and the same refusals
  as `system_application_source_set` — a page the sandbox could not run is
  refused wherever it is written. The declaration travels back untouched:
  sending the page alone would empty what the app may read.

## Routes

- `system.applications.create` → `/system/applications/new`; `.edit` →
  `/system/applications/:applicationID`; back lands on
  `system.applications`.

## When changing this

- All CRUD goes through `useApplicationsStore` (class `system.Application`)
  — keep the store in the loop so the app selector stays fresh.
- Uploaded logo URL is stored absolute (`$SystemAPI.baseURL` + path) with
  `logoID`; clearing resets logo to `''` and logoID to `'0'`.
