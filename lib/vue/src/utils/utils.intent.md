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
- appUrl.ts — `resolveAppUrl` + `isAppUrlLocal`: a registry application's `unify.url` → the href to link it by, and whether that url is a path this shell serves at all.
- internalNav.ts — `resolveInternalPath` + `handleInternalAnchorClick`: decide whether a link target is a route the shell's router hosts (client-side push) or external (native navigation).

## When changing this

- internalNav treats a catch-all `:pathMatch` match or a redirect to a
  different path as "not ours" — that fallback is what keeps not-yet-merged
  apps working via full page loads. Modified clicks (ctrl/meta/middle) must
  always stay native.
- resolveAppUrl reads a dot in the url's first segment as "this is a host,
  not a path" — that is what tells `www.google.com` from `admin/system/labels`,
  and a bare host without a scheme links as a relative path. `isAppUrlLocal` is
  the same rule, and is what the shell's app menu asks before trying to resolve
  a url to one of its sections.
- resolveAppLogoUrl's API-path branch strips the `/api/system` suffix from the
  base URL — it assumes the system API lives under `<host>/api/system`.
