---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/usePermissions.ts
touched-by: []
tests: []
---

# permissions/

## Intention

The shared RBAC editing UX: any screen drops in a button that opens the one
global permissions dialog for a resource. Open-state plumbing lives in
`usePermissions`; the dialog is mounted once per app.

## Map

- `CPermissionsButton.vue` — opens the dialog for its `resource`; PrimeVue
  Button passthrough props for styling.
- `CPermissionsDialog.vue` — role/user column editor: loads effective rules,
  diffs against initial state, saves per role/user.

## Contracts consumers rely on

- Button props: `resource` (RBAC resource string, e.g.
  `corteza::system:role/12345`) is the contract; `title`/`target` feed i18n
  interpolation, `allSpecific` switches to the "all-specific" wording.
- Dialog resolves the owning API client from the resource string — inject
  `$SystemAPI`/`$ComposeAPI`/`$AutomationAPI` (all optional); a component kind
  from a service whose API isn't provided cannot be edited in that app.
- Strings come from `permissions.ui.*` keys in the host app's locale bundle.

## When changing this

One dialog instance serves the whole app via usePermissions state — keep it
resource-agnostic; per-resource behavior belongs in the resource string
parsing, not in call sites.
