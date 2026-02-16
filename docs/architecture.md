# Architecture

## Package Hierarchy

```
Web apps  →  @cortezaproject/corteza-vue-next  →  @cortezaproject/corteza-js-next
(one, compose, taq)    (lib/vue)                       (lib/js)
```

Workspace dependencies use `workspace:*` protocol via pnpm.

## Source-Level Imports

Lib packages export TypeScript source directly (`main: "src/index.ts"`). No build step required for libs during development — Vite compiles TypeScript on-the-fly.

## Plugin Setup Order

Every webapp follows the same initialization sequence in `src/plugins/index.js`:

1. **AuthPlugin** — OAuth 2.0 PKCE flow (`rootApp: true`)
2. **API plugins** — SystemAPI, ComposeAPI, AutomationAPI (inject access token from auth)
3. **SettingsPlugin** — Loads system/app settings from API
4. **Pinia** — State management
5. **Vue Router** — Routing
6. **I18nPlugin** — Loads translations for user's preferred language
7. **PrimeVue + theme** — UI framework with theme from settings

## Authentication

OAuth 2.0 with PKCE flow via `AuthPlugin` (`lib/vue/src/plugins/auth.ts`).

- Access tokens: in-memory only (never persisted)
- Refresh tokens: sessionStorage
- Auto-refresh at 75% of token lifetime
- Duplicate tab detection and session management
- Auth server URL from `window.CortezaAPI` config

## Configuration

Each app requires `public/config.js` (copy from `config.example.js`):

```javascript
window.CortezaAPI = 'https://your-corteza-instance/api'
```

## Dependency Injection

All apps share these injected dependencies (via Vue provide/inject):

| Key | Type | Source |
|-----|------|--------|
| `$Auth` | Auth instance | AuthPlugin |
| `$SystemAPI` | API client | SystemAPIPlugin |
| `$ComposeAPI` | API client | ComposeAPIPlugin |
| `$AutomationAPI` | API client | AutomationAPIPlugin |
| `$Settings` | Settings manager | SettingsPlugin |
| `$toast` | Toast service | ToastPlugin |

Access in `<script setup>` with `const api = inject('$ComposeAPI')`.
