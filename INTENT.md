---
kind: folder
covers: '.'
owner: shared
depends-on: []
touched-by: []
tests: []
---

# Human monorepo

## Intention

Low-code platform (Corteza fork): Go server + Vue 3 web clients. This doc maps the
monorepo; every covered area carries its own `INTENT.md` (see `.intent/SPEC.md` for
the intent system's rules).

## Map

- `client/web/unify` — the single unified webapp (all sections: admin, compose, taq, workflow, …)
- `client/web/chatbot-widget` — embeddable chatbot widget (separate Vite build)
- `lib/js` — API clients + shared JS (codegen from server rest.yaml)
- `lib/vue` — shared Vue component library (`@planetcrust/human-vue`); one-way dependency: apps import it, never the reverse
- `server` — Go backend (system, compose, automation, federation components; cue-driven codegen)
- `def` — cue definitions driving server codegen
- `locale` — translation sources; unify uses the single merged `human-webapp` bundle
- `tests` — e2e fixtures and tests
- `extra` — auxiliary services (server-discovery), not covered by the intent system

## When changing this

- Root config (Makefile, pnpm workspace, shared eslint/tailwind/prettier) affects every package — run `make intent-check` plus the affected app builds.
- Codegen order matters: server codegen before lib codegen (`make codegen`).
