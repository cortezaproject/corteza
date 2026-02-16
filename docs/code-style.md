# Code Style

## Formatting (Prettier)

Configured in `.prettierrc.json`. Runs on save.

- 2-space indentation
- Single quotes (JS/TS), double quotes (HTML attributes)
- No semicolons
- Trailing commas in multiline
- 100 character line width

## Linting (ESLint)

Configured in `eslint.config.shared.js`. Handles code quality only, not formatting.

## Styling

### CSS Layer Order

PrimeVue and Tailwind work together with this layer order:

```
tailwind-base, primevue, tailwind-utilities
```

This ensures Tailwind utilities can override PrimeVue defaults.

### Tailwind First

- **Always prefer Tailwind utility classes** over custom CSS in `<style>` blocks
- Use arbitrary value syntax when needed: `w-[240px]`, `h-[32px]`
- Only use `<style>` blocks for complex selectors, pseudo-elements, or CSS that can't be expressed with Tailwind
- For PrimeVue colors, use arbitrary values with CSS variables: `text-[--p-text-muted-color]`, `bg-[--p-primary-color]`

### Dark Mode

Uses class selector `[class~="dark"]` on the root element. Toggled via `useTheme('dark')` from the profile menu.

## Naming Conventions

### Component prefixes

- **`C` prefix** — shared components in `lib/vue` that are reusable across all apps (e.g. `CFieldEditor`, `CInputSearch`, `CSidebar`, `CResourceList`). The `C` stands for Corteza.
- **App prefix** — app-specific components use the app name as prefix (e.g. `TaqIcon` in taq, `CNamespaceSidebar` in compose). This avoids collisions when components are only used within one app.

Rule of thumb: if it lives in `lib/vue/src/components/`, use `C`. If it lives in `client/web/{app}/src/components/`, use the app name or a descriptive prefix.

### Files and directories

- **Views:** PascalCase (`Dashboard.vue`, `Builder.vue`), organized in directories matching route structure
- **Stores:** camelCase (`namespace.js`, `automation.ts`), one store per entity type
- **Composables:** camelCase with `use` prefix (`useFlowEditor.ts`, `useResourceList.ts`)
- **Utilities:** camelCase (`taq-parser.ts`, `flow-constants.ts`)

### Route names

Dot-separated hierarchy: `namespace.view`, `admin.modules.edit`, `admin.pages.builder`

## Technology Stack

- Vue 3.5.x, Vite 7.x, pnpm 10.x workspaces
- Pinia 3.x (state), Vue Router 4.x, vue-i18n 11.x
- PrimeVue 4.x (UI), Tailwind CSS 3.x (styling)
- Vitest 3.x (testing), ESLint 9.x (linting)
- Node.js >= 20.0.0
