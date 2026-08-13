---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useTheme.ts
  - lib/vue/src/composables/useFileUpload.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Theming Index view

## Intention

Brand the whole webapp: upload main/icon logos and tune per-theme colors and
custom CSS (`ui.studio` settings), with changes applied live on save.

## UX capabilities

- Logo drop zones upload straight to the `ui.mainLogo` / `ui.iconLogo` setting endpoints; a logo is clearable (back to default) only when its raw setting value starts with `attachment:`.
- Theme studio tabs — light (default) / dark / general: color variables (light+dark only; primary, success, warning, danger, body/sidebar/topbar backgrounds) with per-variable reset-to-default, and custom CSS per tab.
- Save writes `ui.studio.themes` + `ui.studio.custom-css` and re-applies the theme immediately (`setThemes` + `useTheme`) so the admin sees changes without refreshing.

## Routes

`ui.theming` at `/ui/theming`; no params.

## When changing this

- Colors are stored as `#`-prefixed hex inside a JSON-stringified `values` blob per theme, but edited without the `#` — keep the strip/add round-trip intact.
- The theme runtime and app shell consume these settings — key renames are cross-app contract changes.
- Preserve the `attachment:` check when touching logo upload/clear.
