# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Detailed developer docs live in [`docs/`](docs/INDEX.md).** Per-app docs are in each webapp's `DOCS.md`. Keep both up to date when making changes.

## Project Overview

Corteza is a low-code platform for building CRM, business process, and structured data applications. This is the Vue 3 monorepo with web applications and shared libraries.

## Repository Structure

```
corteza-vue3/
├── client/web/
│   ├── one/              # Launcher/dashboard (DOCS.md inside)
│   ├── compose/          # Low-code app builder (DOCS.md inside)
│   └── taq/              # Automation builder (DOCS.md inside)
├── lib/
│   ├── js/               # @cortezaproject/corteza-js-next - API clients, types, utilities
│   └── vue/              # @cortezaproject/corteza-vue-next - Vue plugins, composables, components
├── server/               # Backend server (Go)
├── locale/               # i18n files
└── docs/                 # Developer documentation (INDEX.md for navigation)
```

## Quick Reference

See `docs/` for full details. Key points:

- **Commands:** `pnpm install`, `cd client/web/{app} && pnpm dev` — see [docs/commands.md](docs/commands.md)
- **Architecture:** Plugin setup order, auth flow, dependency injection — see [docs/architecture.md](docs/architecture.md)
- **Shared components:** Field editors, inputs, navigation, resource list — see [docs/components.md](docs/components.md)
- **Composables & plugins:** useResourceList, useRBAC, useTheme, auth, API, settings — see [docs/composables-and-plugins.md](docs/composables-and-plugins.md)
- **Code style:** Prettier, ESLint, Tailwind-first, CSS layers — see [docs/code-style.md](docs/code-style.md)
- **i18n:** No hardcoded strings, YAML locale files, `$t()` / `useI18n()` — see [docs/i18n.md](docs/i18n.md)

## Documentation Rule

**Every feature, fix, or refactor that changes behavior must include doc updates.** This applies to both human developers and AI models. See [docs/TODO.md](docs/TODO.md) for the full merge checklist.

## Critical Rules

- **PrimeVue components are globally registered** — no per-file imports. Use `Button`, `Card`, `DataTable`, etc. directly. To add new ones, edit `lib/vue/src/plugins/primevue-components.ts`.
- **PrimeVue composables** still need explicit imports: `import { useConfirm } from 'primevue/useconfirm'`
- **Shared components** are imported via: `import { components } from '@cortezaproject/corteza-vue-next'`
- **Tailwind first** — prefer utility classes over `<style>` blocks. Use PrimeVue CSS vars: `text-primary`
- **No hardcoded UI strings** — use `$t()` with locale YAML files
- **Lib packages export source** (`main: "src/index.ts"`) — no build step needed for libs
- **2-space indent, single quotes, no semicolons, trailing commas** (Prettier enforced)
