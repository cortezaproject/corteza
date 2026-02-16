# Shared Components (`lib/vue`)

All shared components are exported from `@cortezaproject/corteza-vue-next`:

```js
import { components } from '@cortezaproject/corteza-vue-next'
const { CResourceList, CFieldEditor, CSidebar } = components
```

## Directory Structure

```
lib/vue/src/components/
├── field/          # Editable form field system
├── input/          # Selection/search inputs
├── navigation/     # Sidebar, topbar, toolbar
├── loader/         # Loading indicators
└── resource-list/  # Data table with filtering/sorting
```

---

## Field Editors (`field/`)

Dynamic field editing system with a registry that resolves the correct editor based on field type.

### CFieldEditor

Wrapper that auto-resolves the correct editor component via `resolveFieldEditor()`.

| Prop | Type | Description |
|------|------|-------------|
| `field` | Object (required) | Field definition with `kind`, `isMulti`, `options` |
| `modelValue` | string \| Array | Current value(s) |
| `disabled` | boolean | Disable editing |
| `addLabel` | string | Label for "add" button in multi-value mode |

Supports multi-value fields — renders add/remove buttons automatically.

### Available Editors

| Component | Field Kind | Key Options |
|-----------|-----------|-------------|
| CFieldString | `String` | `multiLine`, `maxLength` |
| CFieldNumber | `Number` | `precision`, `min`, `max`, `prefix`, `suffix` |
| CFieldBool | `Bool` | `switch` (toggle vs checkbox) |
| CFieldSelect | `Select` | `options[]` array of `{text, value}` |
| CFieldDateTime | `DateTime` | `onlyDate`, `onlyTime`, `onlyFutureValues`, `onlyPastValues` |

### Registry

```js
import { FIELD_REGISTRY, resolveFieldEditor } from '@cortezaproject/corteza-vue-next'
// resolveFieldEditor(field) returns the component for a given field kind
```

---

## Input Components (`input/`)

### CInputSearch

Search box with clear button and optional submit.

| Prop | Type | Description |
|------|------|-------------|
| `modelValue` | string | Search query |
| `placeholder` | string | Placeholder text |
| `submittable` | boolean | Show search button |

### CInputSelect

Generic autocomplete dropdown built on PrimeVue AutoComplete.

| Prop | Type | Description |
|------|------|-------------|
| `modelValue` | string/number/object | Selected value |
| `options` | Array | Available options |
| `optionLabel` | string \| function | Display property or formatter |
| `showClear` | boolean | Show clear button |

Emits `search` on typing for server-side filtering.

### Resource Pickers

These inputs load data from the API with debounced search and request cancellation.

| Component | Selects | Requires | v-model emits |
|-----------|---------|----------|---------------|
| CInputUser | User | — | userID |
| CInputNamespace | Namespace | — | namespaceID |
| CInputModule | Module | `namespaceID` prop | moduleID |
| CInputRecord | Record | `namespaceID` + `moduleID` props | recordID |

CInputRecord also accepts a `labelField` prop to control which record field is displayed.

---

## Navigation (`navigation/`)

### CSidebar

Collapsible sidebar drawer. Uses PrimeVue Drawer internally.

- **v-model**: `expanded` (boolean)
- Responsive: modal + dismissable on mobile (<1024px)
- Width controlled by CSS variable `--sidebar-width` (320px default)
- Teleport targets for content injection:
  - `#sidebar-header-expanded` — header area (e.g. namespace switcher)
  - `#sidebar-body-expanded` — main nav content (scrollable)
  - `#sidebar-footer-expanded` — footer content

### CTopbar

Application top navigation bar.

| Prop | Type | Description |
|------|------|-------------|
| `sidebarExpanded` | boolean (v-model) | Sidebar toggle state |
| `sidebarDisabled` | boolean | Hide sidebar toggle |
| `hideAppSelector` | boolean | Hide "back to launcher" button |
| `labels` | Object (required) | i18n labels for all UI strings |

Features: hamburger toggle, app selector, help menu, profile menu with theme switcher, avatar.

Teleport targets:
- `#topbar-title` — dynamic page title area
- `#topbar-tools` — custom toolbar buttons

### CToolbar

Simple 3-column layout toolbar with `start`, `center`, `end` slots.

---

## Resource List (`resource-list/`)

### CResourceList

Full-featured data table with search, pagination, sorting, and selection.

| Prop | Type | Description |
|------|------|-------------|
| `primaryKey` | string (required) | DataTable dataKey |
| `fields` | Array | Column definitions: `{key, header, label, sortable, class, frozen}` |
| `items` | Array | Data rows |
| `filter` | Object | Current filters |
| `sorting` | Object | `{sortBy, sortDesc}` |
| `pagination` | Object | `{limit, page, total}` |
| `loading` | boolean | Loading state |
| `selectable` | boolean | Show checkboxes |
| `clickable` | boolean | Pointer cursor on rows |
| `hideSearch` | boolean | Hide search input |
| `hidePagination` | boolean | Hide pagination |

Slots: `header`, `body-{fieldKey}` for custom cell rendering.

---

## Loader (`loader/`)

### CLoaderLogo

Full-screen loading overlay with pulsing logo.

| Prop | Type | Description |
|------|------|-------------|
| `show` | boolean | Show/hide (default true) |
| `logoUrl` | string | Logo image URL |
