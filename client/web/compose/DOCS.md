# Compose (Low-Code Application Builder)

Build CRM, business process, and structured data applications. Manages namespaces, modules, pages, charts, and records.

## Routes

All admin routes are nested under `/namespace/:slug/`.

| Path | Name | View | Description |
|------|------|------|-------------|
| `/` | — | redirect | → `/namespaces` |
| `/namespaces` | `namespace.list` | `Namespace/List.vue` | List namespaces |
| `/namespaces/manage` | `namespace.manage` | `Namespace/Manage.vue` | Manage namespaces |
| `/namespaces/create` | `namespace.create` | `Namespace/Edit.vue` | Create namespace |
| `/namespaces/edit/:slug` | `namespace.edit` | `Namespace/Edit.vue` | Edit namespace |
| `/namespace/:slug` | `namespace.view` | `Namespace/View.vue` | Namespace container (parent) |
| ` ` (empty) | `pages` | `Pages/Index.vue` | Home page redirect / onboarding |
| `.../pages/:pageID` | `page` | `Pages/View.vue` | Public page viewer |
| `.../admin/modules` | `admin.modules` | `Admin/Modules/List.vue` | List modules |
| `.../admin/modules/create` | `admin.modules.create` | `Admin/Modules/Edit.vue` | Create module |
| `.../admin/modules/:moduleID/edit` | `admin.modules.edit` | `Admin/Modules/Edit.vue` | Edit module |
| `.../admin/modules/:moduleID/records` | `admin.modules.record.list` | `Admin/Modules/Records/List.vue` | List records |
| `.../admin/modules/:moduleID/records/create` | `admin.modules.record.create` | `Admin/Modules/Records/Create.vue` | Create record |
| `.../admin/pages` | `admin.pages` | `Admin/Pages/List.vue` | List pages |
| `.../admin/pages/create` | `admin.pages.create` | `Admin/Pages/Edit.vue` | Create page |
| `.../admin/pages/:pageID/edit` | `admin.pages.edit` | `Admin/Pages/Edit.vue` | Edit page |
| `.../admin/pages/:pageID/builder` | `admin.pages.builder` | `Admin/Pages/Builder.vue` | Page builder |
| `.../admin/charts` | `admin.charts` | `Admin/Charts/List.vue` | List charts |
| `.../admin/charts/create` | `admin.charts.create` | `Admin/Charts/Edit.vue` | Create chart |
| `.../admin/charts/:chartID/edit` | `admin.charts.edit` | `Admin/Charts/Edit.vue` | Edit chart |

Sidebar is disabled on top-level namespace routes, auto-expands on desktop for admin routes.

## Stores

All stores use `inject('$ComposeAPI')` or `$SystemAPI`, wrap responses in model classes, and maintain frozen immutable sets.

### useNamespaceStore (`stores/namespace.js`)

- **State:** `loading`, `pending`, `set[]`
- **Getters:** `getByID`, `getByUrlPart`
- **Actions:** `load(force)`, `findByID(id)`, `create(item)`, `clone(item)`, `update(item)`, `delete(item)`

### useModuleStore (`stores/module.js`)

- **State:** `loading`, `pending`, `set[]`, `namespaceID`
- **Getters:** `getByID`, `getByHandle`
- **Actions:** `load({namespace, namespaceID})`, `findByID({namespaceID, moduleID})`, `create`, `update`, `delete`

### usePageStore (`stores/page.js`)

Same pattern as module store. State scoped to `namespaceID`.

### usePageLayoutStore (`stores/page-layout.js`)

- **State:** `loading`, `pending`, `set[]`, `namespaceID`
- **Getters:** `getByID`, `getByPageID` (filters layouts for a given page)
- **Actions:** `load({namespace, namespaceID})`, `findByID({namespaceID, pageID, pageLayoutID})`, `create`, `update`, `delete`
- Uses `$ComposeAPI.pageLayoutListNamespace()` to load all layouts for a namespace
- Preloaded alongside other stores in `Namespace/View.vue`

### useChartStore (`stores/chart.js`)

Same pattern as module store. State scoped to `namespaceID`.

### useUserStore (`stores/user.js`)

- **State:** `pending`, `set[]`
- **Getters:** `findByID`, `findByUsername`
- **Actions:** `load(filter)`, `fetchUsers(userIDs[])`, `resolveUsers(list)` (lazy-loads unknown users only)

## Components

### CNamespaceSidebar (`components/CNamespaceSidebar.vue`)

Teleports content into `CSidebar` slots:
- Header: namespace selector dropdown
- Body: page tree navigation (visible pages as expandable tree, highlights current page)
- Footer: navigation buttons (Modules, Pages, Charts) with active state

**Props:** `namespace` (required), `namespaces` (array)

### PageBlocks System (`components/PageBlocks/`)

Block rendering system for public page viewer:

- **`registry.ts`** — Maps block kind strings to Vue components (like `lib/vue/src/components/field/registry.ts`). Use `resolveBlock(kind)` to get the component. Currently registered: `Content`.
- **`PageBlock.vue`** — Block wrapper that resolves block kind via registry, wraps in Card (if `block.style.wrap.kind === 'card'`) or plain div. Shows fallback message for unknown block types.
- **`Grid.vue`** — Uses `grid-layout-plus` (Vue 3 successor of `vue-grid-layout`) for block positioning. 48 columns, 10px row height, converts `xywh` arrays to layout items. Supports drag-and-drop and resize when `editable` prop is `true` (default `false`). Responsive: collapses to single column on small screens. Emits `update:blocks` with updated xywh values on layout changes.
- **`Blocks/ContentBlock.vue`** — Renders `block.options.body` as HTML via `v-html`.

**Adding new block types:** Add the component to `Blocks/`, then register it in `registry.ts`.

## Public Page Viewer

The public page viewer (`Pages/View.vue`) renders pages with their blocks:

1. Gets page from `pageStore.getByID(pageID)`
2. Gets layouts from `pageLayoutStore.getByPageID(pageID)`
3. Picks the first layout (no multi-layout/visibility expressions yet)
4. Merges layout block positions with page block definitions (layout provides `xywh`, page provides block config)
5. Renders positioned blocks in a CSS grid

**Home page redirect:** `Pages/Index.vue` finds the first visible, first-level, non-record page sorted by weight and redirects to it. Shows onboarding if no pages exist.

## Architecture Notes

- Namespace context comes from URL `:slug` param
- Child stores lazy-load per namespace and cache (skip reload if same namespace)
- `Namespace/View.vue` is the parent container that loads namespace data and renders nested routes
- Public page routes are nested under `namespace.view`, receiving `namespace` as a prop
