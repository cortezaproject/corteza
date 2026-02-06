# Corteza Vue 3 Development Guide

## Migration Status

### Tech Stack Migration (complete)

| Old | New |
|-----|-----|
| Vue 2.7 | Vue 3.5 |
| Vue CLI / Webpack | Vite 7 |
| Vuex | Pinia 3 |
| Bootstrap-Vue | PrimeVue 4 + Tailwind CSS 3 |
| vue-i18next | vue-i18n 11 |
| Jest | Vitest 3 |
| Built lib packages (dist/) | Source-level imports (no lib build step) |

### Web Apps

| App | Status | Notes |
|-----|--------|-------|
| **compose** | ~25-30% | Namespace management done, admin stubs, no public pages/builder/records |
| **one** | In progress | Basic shell |
| **taq** | In progress | Replaces old `workflow` app, uses `@vue-flow/core` |
| admin | Not started | |
| discovery | Not started | |
| privacy | Not started | |
| reporter | Not started | |

### Compose App Progress

**Done:**
- App shell (topbar, sidebar navigation)
- Authentication flow
- Namespace management (List, Edit, Create, Manage, View)
- Admin - Module List/Edit (partial — 352 lines vs 1065 in old)
- Admin - Page List/Edit (basic)
- Admin - Chart List/Edit (Edit is a stub)
- Pinia stores (namespace, module, page, chart, user)
- Routing infrastructure
- i18n integration

**Not started:**
- Public page viewing/rendering (Pages/Index.vue is a stub)
- Page Builder (was 1260 lines — core feature)
- Page blocks system
- Record CRUD operations and list views
- Module field editing UI
- Chart rendering
- Drafts system

### Shared Library Progress (lib/vue)

**Done:**
- Navigation: CSidebar, CTopbar, CToolbar
- Resource list: CResourceList
- Basic inputs: CInputText, CInputSearch, CInputSelect, CInputUser
- Composables: useResourceList, useTheme
- Plugins: auth, corteza-api, i18n, settings, toast
- Loader: CLoaderLogo

**Not started (existed in old lib/vue):**
- Rich text editor
- File preview / upload
- DateTime, ColorPicker, Expression inputs
- Chart components
- Permissions UI
- Reminders
- Privacy settings
- Map / webcam components
- Lightbox
- Search

---

## Important Patterns

### Toast Notifications

**Use the toast plugin, NOT PrimeVue's `useToast` directly:**

```typescript
// ❌ WRONG - Don't use PrimeVue's useToast directly
import { useToast } from 'primevue/usetoast'
const toast = useToast()
toast.add({ severity: 'success', summary: 'Saved', life: 3000 })

// ✅ CORRECT - Use the toast plugin from lib/vue
import { inject } from 'vue'
const $toast = inject('$toast')
$toast.success('Saved')
$toast.error('Failed to save')
```

The toast plugin provides a consistent API across all apps and handles common patterns.

---

## Component Architecture

### Two-Layer Field System (lib/vue) — planned

> **Status:** Only basic Layer 1 inputs exist (Text, Search, Select, User). Layer 2 has not been started.

**Layer 1: Generic Inputs** (`lib/vue/src/components/input/`)
- Pure Vue input components without Corteza business logic
- Reusable in any Vue 3 application
- Currently implemented: `CInputText`, `CInputSearch`, `CInputSelect`, `CInputUser`
- Not yet ported: `CInputNumber`, `CInputDateTime`, `CInputColorPicker`, `CInputRichText`, etc.

**Layer 2: Record-Aware Fields** (`lib/vue/src/components/field/`) — not started
- Wrappers that understand `Record` and `ModuleField` context
- Will use Layer 1 components internally
- Examples: `CFieldString`, `CFieldNumber`, `CFieldDateTime`

```
lib/vue/src/components/
├── input/          # Layer 1: Generic inputs (4 components done)
│   ├── CInputText.vue
│   ├── CInputSearch.vue
│   ├── CInputSelect.vue
│   ├── CInputUser.vue
│   └── ...              # Many more needed
└── field/          # Layer 2: Record-aware fields (not started)
    ├── CFieldString.vue
    ├── CFieldNumber.vue
    └── ...
```

---

## Internationalization (i18n)

### Use Existing Translations

**Always check existing translations before adding new keys.** The locale files already have comprehensive translations.

**Translation files:** `locale/en/corteza-webapp-compose/`
- `general.yaml` - Common labels (save, delete, back, loading, etc.)
- `namespace.yaml` - Namespace-specific translations
- `module.yaml` - Module-specific translations
- `notification.yaml` - Toast/notification messages
- `field.yaml` - Field-related translations

### Key Patterns

```vue
<!-- Use general labels for common actions -->
{{ $t('general.label.save') }}
{{ $t('general.label.delete') }}
{{ $t('general.label.back') }}
{{ $t('general.label.loading') }}

<!-- Use domain-specific prefixes -->
{{ $t('namespace.edit') }}
{{ $t('module.createLabel') }}

<!-- Notifications are in notification.yaml -->
$toast.toastSuccess(t('notification.namespace.saved'))
$toast.toastDanger(t('notification.module.deleteFailed'))
```

### Before Adding New Keys

1. Search existing yaml files for similar keys
2. Check `general.yaml` for common labels
3. Check `notification.yaml` for toast messages
4. Only add new keys if truly unique to your feature
