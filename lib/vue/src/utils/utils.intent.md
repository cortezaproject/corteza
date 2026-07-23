---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - lib/vue/src/composables/useInternalLink.ts
  - client/web/unify/src
tests: []
---

# Utils

## Intention

Small framework-free helpers backing the unified-shell navigation model and
app branding, kept dependency-light so both lib components and host apps can
import them.

## Map

- appIcons.ts — `resolveAppLogoUrl`: application logo/icon → URL (bundled icon map from host, API-relative paths, absolute URLs; host supplies the iconMap so the lib stays asset-free).
- internalNav.ts — `resolveInternalPath` + `handleInternalAnchorClick`: decide whether a link target is a route the shell's router hosts (client-side push) or external (native navigation).

## When changing this

- internalNav treats a catch-all `:pathMatch` match or a redirect to a
  different path as "not ours" — that fallback is what keeps not-yet-merged
  apps working via full page loads. Modified clicks (ctrl/meta/middle) must
  always stay native.
- resolveAppLogoUrl's API-path branch strips the `/api/system` suffix from the
  base URL — it assumes the system API lives under `<host>/api/system`.
