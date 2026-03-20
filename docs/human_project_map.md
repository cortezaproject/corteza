# Human Project — Mental Map

> **Human** is the Vue 3 rewrite of Corteza's frontend, sharing the same Go backend and `lib/js`. It uses a modern stack (PrimeVue 4, Tailwind 3 + PrimeUI, Vite 7, Pinia 3) and introduces new apps (TAQ, Agentic).

---

## Architecture Overview

```
human/
├── client/web/          # Vue 3 frontend apps (6 apps)
├── lib/js/              # Shared JS library (SAME as Corteza — types, API clients, models)
├── lib/vue/             # NEW Vue 3 reusable library (replaces Corteza's lib/vue)
├── server/              # Go backend (SAME as Corteza)
├── locale/en/           # New locale structure (per-webapp YAML files)
├── old-locale/en/       # Corteza's original locale files (fallback reference)
├── def/protobuf/        # Protobuf definitions
├── docs/                # Documentation
└── extra/               # Extra resources
```

---

## Client Apps (`client/web/`)

| App | Base Path | In Corteza? | Notes |
|-----|-----------|-------------|-------|
| **admin** | `/admin/` | ✅ | System admin panel (users, roles, connections, settings) |
| **compose** | `/compose/` | ✅ | Low-code app builder (namespaces, modules, pages, records) |
| **workflow** | `/workflow/` | ✅ | Visual workflow editor |
| **one** | `/one/` | ✅ | Landing/dashboard app |
| **taq** | `/taq/` | ❌ NEW | TAQ — automation triggers/functions with grouping |
| **agentic** | `/agentic/` | ❌ NEW | AI/Agent capabilities |

### Apps in Corteza NOT in Human
- `discovery` — Search/discovery app
- `privacy` — Privacy/GDPR  
- `reporter` — Reporting/analytics

### Each App Structure
```
app/
├── src/
│   ├── App.vue           # Root component (sidebar, topbar, router-view, toast, confirm)
│   ├── main.js           # Bootstrap: createApp → setupAndAuthenticate → mount
│   ├── config-check.js   # Validates window.CortezaAPI is configured
│   ├── plugins/index.js  # Plugin chain setup
│   ├── router/index.js   # Vue Router routes
│   ├── views/            # Page components
│   ├── components/       # App-specific components
│   └── stores/           # App-specific Pinia stores (some apps)
├── public/               # Static assets + config.js
├── vite.config.js        # Vite config with base path, proxy, version
├── tailwind.config.js    # Extends shared tailwind config
└── package.json          # App-specific deps
```

---

## Bootstrap Chain (per-app `plugins/index.js`)

```
1. AuthPlugin          → OAuth/OIDC auth flow
2. SystemAPI, ComposeAPI, AutomationAPI  → API clients
3. SettingsPlugin      → Load server settings
4. Pinia               → State management
5. Router              → Vue Router
6. I18nPlugin          → Load locale from API
7. PrimeVue            → Theme preset + dark mode + CSS layers
8. PrimeVueComponentsPlugin → Global component registration
```

---

## `lib/vue` — Shared Vue 3 Library

Package: `@cortezaproject/corteza-vue-next`

### Plugins
| Plugin | File | Purpose |
|--------|------|---------|
| `AuthPlugin` | `auth.ts` | OAuth/auth handling |
| `SystemAPIPlugin` | `corteza-api.ts` | System API client |
| `ComposeAPIPlugin` | `corteza-api.ts` | Compose API client |
| `AutomationAPIPlugin` | `corteza-api.ts` | Automation API client |
| `FederationAPIPlugin` | `corteza-api.ts` | Federation API client |
| `I18nPlugin` | `i18n.ts` | vue-i18n setup |
| `SettingsPlugin` | `settings.ts` | Server settings loader |
| `PrimeVueComponentsPlugin` | `primevue-components.ts` | Global PrimeVue component registration |
| `ToastPlugin` | `toast.ts` | Toast notification wrapper |

### Composables
| Composable | Purpose |
|-----------|---------|
| `useRBAC` | Permission checking (RBAC store) |
| `useConfirmDelete` | Confirm dialog for delete actions |
| `useResourceList` | Reusable resource list logic (search, pagination, etc.) |
| `useTheme` | Theme management + `getTheme`/`setThemes` |
| `useMinDuration` | Minimum duration wrapper for async ops (loading states) |
| `useUserResolver` | Resolve user IDs to user objects |

### Stores (Pinia)
| Store | Purpose |
|-------|---------|
| `useApplicationsStore` | Fetch/cache applications list |
| `useComposeResourceStore` | Compose resource management |

### Components
| Component/Dir | Purpose |
|--------------|---------|
| `CEmojiPicker.vue` | Emoji picker |
| `chart/` | Chart components |
| `field/` | Field rendering components |
| `input/` | Input components (`CInputSwitch`, `CInputRole`) |
| `loader/` | Loading indicators |
| `navigation/` | Navigation components (`CTopbar`, `CSidebar`, `CAppListSidebar`) |
| `resource-list/` | Reusable resource list UI |

---

## Styling Stack

| Layer | Technology |
|-------|-----------|
| CSS Framework | **Tailwind CSS 3** |
| Component Library | **PrimeVue 4.5** |
| Integration | **tailwindcss-primeui** plugin |
| CSS Layer Order | `tailwind-base, primevue, tailwind-utilities` |
| Dark Mode | Selector-based: `.dark` class |
| Design Tokens | PrimeVue CSS variables (`--p-*`) |

### Key Style Rules (from gemini.md)
- ❌ No `dark:` classes or `surface-80` etc.
- ✅ Use PrimeVue style classes: `primary`, `primary-contrast`, `border-surface`, `bg-emphasis`, `bg-highlight`, `text-color`, `text-muted-color`, `rounded-border`, etc.

---

## Locale Structure

```
locale/en/                        # Active locale files
  ├── config.yaml
  ├── corteza-server/             # Server-side translations
  ├── corteza-webapp-admin/       # Per-app translations
  ├── corteza-webapp-compose/
  ├── corteza-webapp-taq/
  ├── corteza-webapp-agentic/
  ├── corteza-webapp-one/
  └── corteza-webapp-workflow/

old-locale/en/                    # Corteza originals (fallback reference)
```

### Translation Rules
1. Always use `$t('key')` — no hardcoded text
2. Check `old-locale` first if translations are missing
3. Add new translations to `locale` if not in `old-locale`
4. Never use fallback translations — missing translations must be visible

---

## Corteza vs Human — Key Differences

| Aspect | Corteza | Human |
|--------|---------|-------|
| Vue version | Vue 2 (Options API) | **Vue 3** (Composition API, `<script setup>`) |
| CSS | Bootstrap | **Tailwind 3 + PrimeVue** |
| Components | BootstrapVue | **PrimeVue 4** |
| State | Vuex | **Pinia 3** |
| Build | Webpack | **Vite 7** |
| Routing | vue-router 3 | **vue-router 4** |
| i18n | vue-i18n 8 | **vue-i18n 11** |
| lib/vue | `corteza-vue` (Vue 2) | `corteza-vue-next` (Vue 3) |
| lib/js | Shared | **Same** — shared between both |
| Server | Go backend | **Same** — shared between both |
| Extra apps | discovery, privacy, reporter | **taq, agentic** |
| Package mgr | yarn | **pnpm** |

---

## Workspace Config

- **Root `pnpm-workspace.yaml`** — monorepo workspace
- **Root `package.json`** — shared deps (PrimeVue, Vue, Tailwind, etc.)
- **Root `tailwind.config.shared.js`** — shared Tailwind config (PrimeUI plugin, content paths inc. lib/vue)
- **Root `eslint.config.shared.js`** — shared ESLint config
- **`.agent/workflows/`** — Agent workflow definitions
