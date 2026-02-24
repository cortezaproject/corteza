# Permissions System

## How It Works

Corteza uses a **two-level** permission system in the frontend:

### Level 1: Resource-Level Permissions

The API returns `can*` boolean flags alongside every resource. The JS type classes (`Module`, `Namespace`, `Page`, `Chart`, `Record`) already have these properties defined and populated via `Apply()`.

**Example API response:**
```json
{
  "moduleID": "123",
  "name": "Leads",
  "canUpdateModule": true,
  "canDeleteModule": false,
  "canCreateRecord": true,
  "canGrant": false
}
```

**Usage in templates:**
```vue
<Button v-if="module.canDeleteModule" @click="handleDelete" />
```

### Level 2: Global RBAC Store

The `useRBACStore` Pinia store loads **effective permissions** from the API and exposes a `can(resource, operation)` method. Used for permission checks not tied to a specific resource instance (e.g., "can the user create applications system-wide?").

**Location:** `lib/vue/src/composables/useRBAC.ts`

**Loaded in App.vue** during init alongside namespaces and users:
```typescript
import { useRBACStore } from '@cortezaproject/corteza-vue-next'
const rbacStore = useRBACStore()
rbacStore.load([$ComposeAPI, $SystemAPI])
```

**Usage in components:**
```vue
<script setup>
import { useRBACStore } from '@cortezaproject/corteza-vue-next'
const rbac = useRBACStore()
const canCreateApp = computed(() => rbac.can('system/', 'application.create'))
</script>
```

---

## Resource Permission Properties

| Resource | Properties | File |
|----------|-----------|------|
| Namespace | `canCreateChart`, `canCreateModule`, `canCreatePage`, `canDeleteNamespace`, `canUpdateNamespace`, `canManageNamespace`, `canCloneNamespace`, `canExportNamespace`, `canGrant`, `canExportCharts`, `canExportModules` | `lib/js/src/compose/types/namespace.ts` |
| Module | `canUpdateModule`, `canDeleteModule`, `canCreateRecord`, `canCreateOwnedRecord`, `canGrant` | `lib/js/src/compose/types/module.ts` |
| Page | `canUpdatePage`, `canDeletePage`, `canGrant` | `lib/js/src/compose/types/page.ts` |
| Chart | `canUpdateChart`, `canDeleteChart`, `canGrant` | `lib/js/src/compose/types/chart/base.ts` |
| Record | `canUpdateRecord`, `canReadRecord`, `canDeleteRecord`, `canUndeleteRecord`, `canManageOwnerOnRecord`, `canSearchRevision`, `canGrant` | `lib/js/src/compose/types/record.ts` |
| ModuleField | `canUpdateRecordValue`, `canReadRecordValue` | `lib/js/src/compose/types/module-field/base.ts` |

---

## Views With Permission Checks

All admin list views follow a consistent pattern:
- Create button guarded by `namespace.canCreate*`
- Row click guarded by `canUpdate* || canDelete*`
- Actions menu items guarded by `canUpdate*` and `canDelete*`
- Delete via confirm dialog with API call

| View | Checks | Status |
|------|--------|--------|
| `Admin/Modules/List.vue` | `canCreateModule`, `canUpdateModule`, `canDeleteModule` | Done |
| `Admin/Modules/Edit.vue` | `canDeleteModule` | Done |
| `Admin/Pages/List.vue` | `canCreatePage`, `canUpdatePage`, `canDeletePage` | Done |
| `Admin/Pages/Edit.vue` | Stub — will need checks when built out | Pending |
| `Admin/Charts/List.vue` | `canCreateChart`, `canUpdateChart`, `canDeleteChart` | Done |
| `Admin/Charts/Edit.vue` | Stub — will need checks when built out | Pending |
| `Namespace/Manage.vue` | `canDeleteNamespace` | Done |
| `Namespace/Edit.vue` | `canDeleteNamespace` | Done |

---

## Still TODO

### Permissions UI Components

Port from old codebase (`lib/vue/src/components/permissions/`):

1. **CPermissionsButton** — button shown when `resource.canGrant` is true
   - Takes a `resource` string prop (e.g., `corteza::compose:module/123/456`)
   - Opens the permissions modal

2. **CPermissionsModal** — full RBAC rule editor
   - Allow/deny/inherit toggles per operation per role
   - Uses system API to fetch roles and permission rules

### Field-Level Permissions

When record views are implemented, respect `ModuleField.canReadRecordValue` and `canUpdateRecordValue`:

```vue
<template v-if="field.canReadRecordValue">
  <CFieldViewer :field="field" :record="record" />
</template>
<span v-else class="text-muted-color">{{ $t('field.noPermission') }}</span>

<CFieldEditor
  :field="field"
  :record="record"
  :disabled="!field.canUpdateRecordValue"
/>
```

---

## API Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /api/compose/permissions/effective` | Effective permissions for current user (Compose) |
| `GET /api/system/permissions/effective` | Effective permissions for current user (System) |
| `GET /api/automation/permissions/effective` | Effective permissions for current user (Automation) |
| `GET /api/{service}/permissions/{roleID}/rules` | Permission rules for a role |
| `PATCH /api/{service}/permissions/{roleID}/rules` | Update permission rules for a role |
