# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Corteza is a low-code platform for building CRM, business process, and structured data applications. This repository contains the Vue 3 monorepo with web applications and shared libraries.

## Repository Structure

```
corteza-vue3/
├── client/web/           # Web applications
│   ├── one/              # Main launcher/dashboard (corteza-webapp-one)
│   ├── compose/          # Low-code application builder (corteza-webapp-compose)
│   └── taq/              # Task/automation management (corteza-webapp-taq)
├── lib/
│   ├── js/               # @cortezaproject/corteza-js-next - API clients, types, utilities
│   └── vue/              # @cortezaproject/corteza-vue-next - Vue plugins, composables, components
├── server/               # Backend server (Go)
└── locale/               # Internationalization files
```

## Commands

### Installation

```bash
pnpm install              # Install all dependencies
make dev                  # Alternative: install via make
make fresh                # Clean install (removes node_modules first)
```

### Development

```bash
# Run individual web apps (from their directories)
cd client/web/one && pnpm dev
cd client/web/compose && pnpm dev
cd client/web/taq && pnpm dev

# Or use make from app directories
cd client/web/one && make dev
```

### Testing

```bash
make test                 # Test all (libs, clients, server)
pnpm --filter corteza-webapp-one test:unit    # Test specific app
cd lib/js && yarn test    # Test js library
```

### Linting

```bash
make lint                 # Lint all libs and clients
pnpm --filter corteza-webapp-one lint         # Lint specific app
cd lib/js && yarn lint    # Lint js library
```

### Building

```bash
cd client/web/one && pnpm build   # Build specific app (outputs to dist/)
```

## Architecture

### Package Hierarchy

Web apps depend on `@cortezaproject/corteza-vue-next`, which depends on `@cortezaproject/corteza-js-next`. Workspace dependencies use `workspace:*` protocol.

### Source-Level Imports

Lib packages export TypeScript source directly (`main: "src/index.ts"`). No build step is required for libs during development - Vite compiles TypeScript on-the-fly.

### Component Imports

- **PrimeVue components are globally registered** via `PrimeVueComponentsPlugin` from `lib/vue` — no per-file imports needed. Use `Button`, `Card`, `DataTable`, etc. directly in templates. To add a new global component, edit `lib/vue/src/plugins/primevue-components.ts`.
- **PrimeVue composables** (`useConfirm`, `useToast`) still need explicit imports: `import { useConfirm } from 'primevue/useconfirm'`
- **lib/vue shared components** are imported via the `components` export: `import { components } from '@cortezaproject/corteza-vue-next'`

### Configuration

Each app requires `public/config.js` (copy from `config.example.js`):

```javascript
window.CortezaAPI = 'https://your-corteza-instance/api'
```

## Technology Stack

- Vue 3.5.x, Vite 7.x, pnpm 10.x workspaces
- Pinia 3.x (state), Vue Router 4.x, vue-i18n 11.x
- PrimeVue 4.x (UI), Tailwind CSS 3.x (styling)
- Vitest 3.x (testing), ESLint 9.x (linting)
- Node.js >= 20.0.0

## Code Style

**Formatting**: Prettier handles all code formatting on save (see `.prettierrc.json`)
**Linting**: ESLint handles code quality only, not formatting (see `eslint.config.shared.js`)

Key style rules (enforced by Prettier):

- 2-space indentation
- Single quotes (JS/TS), double quotes (HTML attributes)
- No semicolons
- Trailing commas in multiline
- 100 character line width

## Styling

PrimeVue and Tailwind are configured to work together with CSS layer order:

```
tailwind-base, primevue, tailwind-utilities
```

Dark mode uses class selector `[class~="dark"]` on root element.

### Tailwind First

- **Always prefer Tailwind utility classes** over custom CSS in `<style>` blocks
- Use Tailwind's arbitrary value syntax when needed: `w-[240px]`, `h-[32px]`
- Only use `<style>` blocks for complex selectors, pseudo-elements, or CSS that can't be expressed with Tailwind
- For colors, use PrimeVue CSS variables with Tailwind arbitrary values: `text-[--p-text-muted-color]`, `bg-[--p-primary-color]`

## Authentication

Uses OAuth 2.0 with PKCE flow via `AuthPlugin`. Access tokens are kept in memory only; refresh tokens in sessionStorage.

## Internationalization (i18n)

- **No hardcoded UI strings** - All user-facing text must use translations
- **Locale files**: `locale/en/corteza-webapp-{appname}/` - organized per view (e.g., `dashboard.yaml`, `builder.yaml`)
- **In templates**: Use `$t('view.key')` - e.g., `{{ $t('dashboard.title') }}`
- **In script setup**: Import `useI18n` and use `t()`:
  ```vue
  import { useI18n } from 'vue-i18n' const { t } = useI18n() toast.add({ summary:
  t('builder.toast.saved.summary') })
  ```
- **Interpolation**: `$t('key', { name: value })` with `{name}` placeholders in YAML
