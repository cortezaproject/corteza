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

Brand the whole webapp per theme: upload a logo and icon for the light and
the dark theme, tune per-theme colors and custom CSS (`ui.studio` settings),
with changes applied live on save.

## UX capabilities

- Theme studio tabs — light (default) / dark / general. The light and dark tabs open with that theme's logo and icon drop zones, previewed on the tab's sidebar colour; uploads go straight to the setting endpoint (`ui.main-logo` / `ui.icon-logo`, dark: `ui.main-logo-dark` / `ui.icon-logo-dark`) and clearing writes the same key; the structured settings payload reads them back JSON-cased (`ui.mainLogo`, `ui.mainLogoDark`, …). A logo is clearable (back to default) only when its raw setting value starts with `attachment:`. The page refetches `$Settings` on load, so the slots show the server's current values rather than the boot snapshot. A dark slot left empty means the light image serves the dark theme too (the fallback lives in `useBrandLogo`); the dark tab previews that inherited image dimmed, labelled as inherited, with nothing to clear.
- Color variables (light+dark only; primary, success, warning, danger, body/sidebar/topbar backgrounds) with per-variable reset-to-default, and custom CSS per tab.
- Save writes `ui.studio.themes` + `ui.studio.custom-css` and re-applies the current theme immediately (`setThemes` + `useTheme(currentTheme)`) so the admin sees changes without refreshing. Logo uploads are not part of Save; they land as they happen.

## Routes

`ui.theming` at `/ui/theming`; no params.

## When changing this

- Colors are stored as `#`-prefixed hex inside a JSON-stringified `values` blob per theme, but edited without the `#` — keep the strip/add round-trip intact.
- The theme runtime and app shell consume these settings — key renames are cross-app contract changes.
- Preserve the `attachment:` check when touching logo upload/clear.
- Logo slots are one `logoSlot(key, read)` per theme and kind; adding a theme means adding its two slots, never a parallel set of refs and handlers.
