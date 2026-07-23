---
kind: folder
covers: '.'
owner: fe
depends-on:
  - lib/vue
touched-by: []
tests: []
---

# App shell & bootstrap

## Intention

Boot the app safely (config check → auth → plugins → mount) and render the shell
chrome (topbar, per-section sidebar, global overlays) around the active section.

## Map

- `config-check.js` — hard-fails boot when `window.HumanAPI` is missing (imported first by `main.js`)
- `main.js` — createApp → `setupAndAuthenticate` → mount on `body`; nothing renders before auth resolves
- `App.vue` — shell chrome (own `App.intent.md`)
- `router/`, `plugins/`, `utils/`, `sections/` — own INTENT.md each

## When changing this

- Boot order is a contract: config check, then auth, then API/plugins, then settings/theme, then mount. Reordering breaks auth redirects and theming.
- Unauthenticated users must end up in the auth flow, never on a broken page.
