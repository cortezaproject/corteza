---
kind: folder
covers: '.'
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# Unify webapp

## Intention

The single web client for the whole platform. Replaces the old multi-SPA setup with
one shell app whose feature areas are self-contained **sections** (see
`src/sections/intent doc`). One bootstrap, one router, one Pinia, one merged i18n
bundle (`human-webapp`).

## Map

- `src/` — app shell and bootstrap (own intent doc)
- `index.html`, `vite.config.js` — Vite entry; `public/config.js` must define `window.HumanAPI`
- `jsconfig.json` — `@/` alias → `src/`
- `tailwind.config.js`, `postcss.config.js` — extend shared root configs

## When changing this

- Root-level config here affects every section; verify with a full `pnpm build` for this app.
- Anything reusable across apps belongs in `lib/vue`, not here.
