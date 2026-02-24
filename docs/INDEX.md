# Developer Documentation Index

Internal developer docs for the Corteza Vue 3 monorepo.

**Before merging any PR, follow the checklist in [TODO.md](TODO.md).**

## Shared Docs (this folder)

| File                                                     | Contents                                                                                                     |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| [architecture.md](architecture.md)                       | Package hierarchy, source imports, auth flow, config, plugin setup                                           |
| [components.md](components.md)                           | Shared `lib/vue` components — field editors, inputs, navigation, resource list, loader                       |
| [composables-and-plugins.md](composables-and-plugins.md) | Shared composables (`useResourceList`, `useRBAC`, `useTheme`) and plugins (auth, API, settings, toast, i18n) |
| [commands.md](commands.md)                               | Install, dev, test, lint, build commands                                                                     |
| [code-style.md](code-style.md)                           | Prettier, ESLint, Tailwind, CSS layer order, dark mode                                                       |
| [i18n.md](i18n.md)                                       | Internationalization patterns and locale file conventions                                                    |
| [TODO.md](TODO.md)                                       | General cross-cutting TODO and PR merge checklist                                                            |

## Per-App Docs

Each webapp has a `DOCS.md` in its root describing app-specific routes, stores, views, and components.

| App     | Docs                                     | TODO                                     | Description                                                           |
| ------- | ---------------------------------------- | ---------------------------------------- | --------------------------------------------------------------------- |
| One     | [DOCS.md](../client/web/one/DOCS.md)     | [TODO.md](../client/web/one/TODO.md)     | Launcher/dashboard — app grid, minimal routing                        |
| Compose | [DOCS.md](../client/web/compose/DOCS.md) | [TODO.md](../client/web/compose/TODO.md) | Low-code app builder — namespaces, modules, pages, charts, records    |
| Taq     | [DOCS.md](../client/web/taq/DOCS.md)     | [TODO.md](../client/web/taq/TODO.md)     | Automation builder — VueFlow visual editor, triggers, steps, branches |
