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
- `e2e/` — Playwright suite mirroring the section tree (own intent doc);
  `playwright.config.ts` attaches it to a running dev stack via gitignored
  `.env.e2e` (`.env.e2e.example` documents the shape)
- `index.html`, `vite.config.js` — Vite entry; `public/config.js` must define `window.HumanAPI`.
  The dev-only `human:watch-lib-sources` plugin watches `lib/js/src` and
  `lib/vue/src` as directories: outside vite's `root` they would be watched
  file by file, and a per-file watch does not survive a git write, freezing
  that module's transform until the server restarts.
- `jsconfig.json` — `@/` alias → `src/`
- `tailwind.config.js`, `postcss.config.js` — extend shared root configs

## When changing this

- Root-level config here affects every section; verify with a full `pnpm build` for this app.
- Anything reusable across apps belongs in `lib/vue`, not here.
- Do not drop `human:watch-lib-sources` as redundant — without it a merge or
  branch switch silently freezes a lib module, and a frozen `lib/vue` barrel
  takes the whole app down at import.
