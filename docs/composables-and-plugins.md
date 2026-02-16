# Composables and Plugins (`lib/vue`)

## Composables

### useResourceList

Manages list views with pagination, sorting, filtering, and URL query sync.

```js
import { useResourceList } from '@cortezaproject/corteza-vue-next'

const { items, loading, filter, sorting, pagination, fetchItems, filterList, handleSort } =
  useResourceList({ /* options */ })
```

**Returns:**
- `items`, `loading`, `error` — reactive state
- `filter`, `sorting`, `pagination` — reactive config objects
- `fetchItems(updateQuery?)` — fetch with pagination
- `filterList()` — apply filters, reset page
- `handleSort(event)` — handle sort changes
- `abortRequests()` — cancel pending API requests

Features: URL query parameter persistence, cursor pagination, debounced filtering, 300ms loading delay.

### useRBACStore

Pinia store for role-based access control.

```js
import { useRBACStore } from '@cortezaproject/corteza-vue-next'

const rbac = useRBACStore()
rbac.can('compose:namespace', 'read') // boolean
```

**API:**
- `can(resource, operation)` — check permission
- `load(apis[])` — load rules from API clients
- `clear()` — clear cached rules

Rules are prefixed with `corteza::` for resource matching.

### useTheme

Theme switching with PrimeVue integration.

```js
import { useTheme, getTheme, setThemes } from '@cortezaproject/corteza-vue-next'

useTheme('dark')           // Apply dark theme
const config = getTheme()  // Get current theme config
setThemes(customVars)      // Override theme variables
```

Default colors: primary `#FF9661`, success `#43AA8B`, warning `#E27646`, danger `#E54122`.

Generates CSS custom properties: `--topbar-height`, `--sidebar-width`, `--sidebar-bg`, `--body-bg`.

---

## Plugins

### AuthPlugin

OAuth 2.0 PKCE authentication. See [architecture.md](architecture.md) for flow details.

Provides `$Auth` globally with:
- `user` — current user info
- `accessTokenFn()` — returns current access token
- `logout()` — end session

### API Plugins

`SystemAPIPlugin`, `ComposeAPIPlugin`, `AutomationAPIPlugin`, `FederationAPIPlugin`

Each provides an API client (e.g. `$ComposeAPI`) that auto-injects the access token.

```js
const api = inject('$ComposeAPI')
const { set } = await api.moduleList({ namespaceID })
```

### SettingsPlugin

Loads and provides system/app settings as `$Settings`.

```js
const settings = inject('$Settings')
settings.get('ui.topbar.hideProfile')        // dot-notation access
settings.attachment('ui.mainLogo')           // attachment URL
```

### ToastPlugin

Toast notification wrapper provided as `$toast`.

```js
const toast = inject('$toast')
toast.toastSuccess('Saved successfully')
toast.toastDanger('Something went wrong')
toast.toastErrorHandler('prefix')  // returns error handler function
```

### I18nPlugin

Initializes vue-i18n with app-specific translations. Fallback locale: `en`.

### PrimeVueComponentsPlugin

Globally registers PrimeVue components (no per-file imports needed):

Avatar, Badge, Button, Card, Checkbox, Column, ConfirmDialog, DataTable, DatePicker, Dialog, Divider, IconField, InputIcon, InputNumber, InputText, Menu, Message, ProgressSpinner, Select, Tab, TabList, TabPanel, TabPanels, Tabs, Tag, Textarea, TieredMenu, Toast, ToggleSwitch

To add a new global component, edit `lib/vue/src/plugins/primevue-components.ts`.
