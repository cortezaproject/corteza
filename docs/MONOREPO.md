# Corteza Vue 3 Monorepo

This repository contains the next-generation Corteza web applications, rebuilt with Vue 3, Vite, and modern tooling.

## 📁 Repository Structure

```
corteza-vue3/
├── client/
│   └── web/                    # Web applications
│       ├── one/                # Main launcher/dashboard app
│       ├── compose/            # Low-code application builder
│       └── taq/                # Task/automation app
├── lib/
│   ├── js/                     # @cortezaproject/corteza-js-next
│   │   └── src/                # TypeScript utilities & API clients
│   └── vue/                    # @cortezaproject/corteza-vue-next
│       └── src/                # Vue plugins, composables, components
├── server/                     # Backend server (Go)
├── locale/                     # Internationalization files
├── extra/                      # Additional resources
├── package.json                # Root shared dependencies
└── pnpm-workspace.yaml         # Workspace configuration
```

## 🏗️ Architecture

### Monorepo with pnpm Workspaces

The repository uses pnpm workspaces for managing multiple packages:

```yaml
# pnpm-workspace.yaml
packages:
  - 'lib/*'
  - 'client/web/*'
```

### Package Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                      Web Applications                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │     one      │  │   compose    │  │     taq      │       │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘       │
│         │                 │                 │                │
│         └─────────────────┼─────────────────┘                │
│                           ▼                                  │
│              ┌────────────────────────┐                      │
│              │  @cortezaproject/      │                      │
│              │  corteza-vue-next      │                      │
│              │  (Vue plugins, comps)  │                      │
│              └───────────┬────────────┘                      │
│                          ▼                                   │
│              ┌────────────────────────┐                      │
│              │  @cortezaproject/      │                      │
│              │  corteza-js-next       │                      │
│              │  (API clients, types)  │                      │
│              └────────────────────────┘                      │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Getting Started

### Prerequisites

- **Node.js** >= 20.0.0
- **pnpm** >= 10.x (`npm install -g pnpm`)
- Access to a running Corteza backend

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd corteza-vue3

# Install all dependencies (from root)
pnpm install

# Or use make
make dev
```

This will install dependencies for all workspaces with shared dependencies at the root.

### Running Applications

Using Make (recommended):

```bash
# From app directory
cd client/web/one && make dev
cd client/web/compose && make dev
cd client/web/taq && make dev

# Server
cd server && make dev
```

Using pnpm:

```bash
# Run the "One" app
pnpm --filter corteza-webapp-one dev

# Or navigate to the app directory
cd client/web/one
pnpm dev
```

### Configuration

Each app requires a `config.js` file in its `public/` directory. Copy the example:

```bash
cd client/web/one/public
cp config.example.js config.js
```

Edit `config.js` to point to your Corteza API:

```javascript
window.CortezaAPI = 'https://your-corteza-instance/api'
```

## 📦 Packages

### `@cortezaproject/corteza-js-next` (`lib/js`)

Core JavaScript/TypeScript library containing:

- **API Clients**: `System`, `Compose`, `Automation`, `Federation`
- **Type Definitions**: TypeScript interfaces and types
- **Utilities**: Shared helper functions
- **Models**: Data models (User, Module, Record, etc.)

```typescript
import { apiClients, system, compose } from '@cortezaproject/corteza-js-next'

const systemAPI = new apiClients.System({ baseURL: '...' })
const user = new system.User({ ... })
```

### `@cortezaproject/corteza-vue-next` (`lib/vue`)

Vue 3 specific library containing:

- **Plugins**: Auth, API, Settings, I18n, Toast
- **Composables**: `useResourceList`, `useTheme`
- **Components**: CTopbar, CSidebar, CInputSearch, etc.
- **Filters**: Date/time formatting, etc.

```typescript
import {
  AuthPlugin,
  SystemAPIPlugin,
  SettingsPlugin,
  components,
} from '@cortezaproject/corteza-vue-next'
```

### Web Applications

| App     | Package Name             | Description                                   |
| ------- | ------------------------ | --------------------------------------------- |
| One     | `corteza-webapp-one`     | Main launcher, app selector, workflow builder |
| Compose | `corteza-webapp-compose` | Low-code application builder                  |
| Taq     | `corteza-webapp-taq`     | Task and automation management                |

## 🛠️ Technology Stack

| Category                 | Technology   | Version |
| ------------------------ | ------------ | ------- |
| **Framework**            | Vue          | 3.5.x   |
| **Build Tool**           | Vite         | 7.x     |
| **Package Manager**      | pnpm         | 10.x    |
| **State Management**     | Pinia        | 3.x     |
| **Routing**              | Vue Router   | 4.x     |
| **Internationalization** | vue-i18n     | 11.x    |
| **UI Components**        | PrimeVue     | 4.x     |
| **Styling**              | Tailwind CSS | 3.x     |
| **Testing**              | Vitest       | 3.x     |
| **Linting**              | ESLint       | 9.x     |

## 🔧 Development

### Dependency Management

Shared dependencies are defined in the root `package.json` and hoisted:

```bash
# Add a dependency to a specific workspace
pnpm --filter corteza-webapp-one add lodash-es

# Add a dev dependency
pnpm --filter corteza-webapp-one add -D @types/lodash

# Add to root (shared across all workspaces)
pnpm add -w some-package
```

### Component Imports

#### PrimeVue Components in Web Apps

PrimeVue components are **auto-imported** via `unplugin-vue-components`. No manual imports needed:

```vue
<template>
  <!-- Just use them directly -->
  <Button label="Click me" />
  <DataTable :value="items" />
</template>

<script setup>
// No imports required for PrimeVue components
</script>
```

#### PrimeVue Components in lib/vue

Components in `lib/vue` **must have explicit imports** (not processed by unplugin):

```vue
<script setup>
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
</script>
```

#### lib/vue Components in Web Apps

Import explicitly from the package:

```vue
<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
const { CTopbar, CSidebar } = components
</script>
```

### Source-Level Imports

The lib packages export TypeScript source directly (`main: "src/index.ts"`). This means:

- No build step required for libs during development
- Vite compiles TypeScript on-the-fly
- Hot Module Replacement (HMR) works across packages

### Building for Production

```bash
# Build a specific app
cd client/web/one && pnpm build

# The output will be in client/web/one/dist/
```

### Linting

```bash
# Lint a specific workspace
pnpm --filter corteza-webapp-one lint

# Lint lib packages
pnpm --filter @cortezaproject/corteza-js-next lint
pnpm --filter @cortezaproject/corteza-vue-next lint
```

### Testing

```bash
# Run tests for a specific app
pnpm --filter corteza-webapp-one test:unit
```

## 🎨 Styling Architecture

### CSS Layer Order

PrimeVue and Tailwind CSS are configured to work together:

```javascript
// PrimeVue cssLayer configuration
cssLayer: {
  name: 'primevue',
  order: 'tailwind-base, primevue, tailwind-utilities'
}
```

### Dark Mode

Dark mode is enabled via class selector:

```javascript
// tailwind.config.js
darkMode: ['selector', '[class~="dark"]']
```

Toggle dark mode by adding/removing the `dark` class on the root element.

### Theme System

Themes are managed via the `useTheme` composable:

```typescript
import { useTheme, setThemes } from '@cortezaproject/corteza-vue-next'

// Set available themes
setThemes(themesFromSettings)

// Apply a theme
useTheme('dark')
```

## 🔐 Authentication

Authentication uses OAuth 2.0 with PKCE flow:

1. `AuthPlugin` is installed first during app startup
2. Handles redirect to Corteza Auth
3. Exchanges authorization code for tokens
4. Stores refresh token in sessionStorage
5. Access token kept in memory only

```javascript
// In app startup
app.use(AuthPlugin, { app: 'one', rootApp: true })

await app.config.globalProperties.$Auth.handle()
```

## 📝 Adding a New Web App

1. Create directory: `client/web/your-app/`

2. Copy structure from an existing app (e.g., `one`)

3. Update `package.json`:

   ```json
   {
     "name": "corteza-webapp-your-app",
     "dependencies": {
       "@cortezaproject/corteza-js-next": "workspace:*",
       "@cortezaproject/corteza-vue-next": "workspace:*"
     }
   }
   ```

4. Add a `Makefile`:

   ```makefile
   .PHONY: dev build

   dev:
   	@pnpm dev

   build:
   	@pnpm build

   .DEFAULT_GOAL := dev
   ```

5. Run `pnpm install` from root

6. Create `public/config.js` from example

7. Start development: `make dev`

## 🐛 Troubleshooting

### Dependencies not found

Run `pnpm install` from the repository root to ensure all workspaces are properly linked.

### Changes in lib not reflected

Vite should pick up changes automatically. If not:

1. Restart the dev server
2. Clear Vite cache: `rm -rf node_modules/.vite`

### PrimeVue components not rendering

For web apps: Components are auto-imported, just use them in templates.

For lib/vue: Ensure you have explicit imports for all PrimeVue components used.

### TypeScript errors in IDE

1. Ensure your IDE is using the workspace TypeScript version
2. Run `pnpm install` to ensure all type definitions are available
3. Restart the TypeScript language server

## 📄 License

Apache-2.0 - See [LICENSE](LICENSE) for details.
