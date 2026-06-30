<template>
  <Dialog
    v-model:visible="open"
    modal
    :header="header"
    :style="{ width: '46rem' }"
    :pt="{ footer: { class: 'flex justify-end gap-2' } }"
  >
    <div v-if="loading" class="py-10 grid place-items-center">
      <ProgressSpinner style="width: 2.5rem; height: 2.5rem" />
    </div>

    <div v-else-if="!ops.length" class="py-8 text-center text-muted-color text-sm">
      {{ $t('project.permissions.dialog.noOps') }}
    </div>

    <div v-else class="overflow-auto -mx-1 px-1">
      <table class="text-sm border-collapse w-full">
        <thead>
          <tr class="bg-emphasis">
            <th class="text-left font-medium px-3 py-2 sticky left-0 bg-emphasis z-10 min-w-48">
              {{ $t('project.permissions.dialog.operation') }}
            </th>
            <th
              v-for="role in roles"
              :key="role.id"
              class="font-medium px-3 py-2 border-l border-surface text-center whitespace-nowrap"
              :class="{ 'bg-primary/10 text-primary': role.id === preselectId }"
            >
              {{ role.name }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="op in ops" :key="op.op" class="border-t border-surface hover:bg-emphasis/40">
            <td class="px-3 py-2 sticky left-0 bg-surface" :title="op.description">
              {{ op.label }}
            </td>
            <td
              v-for="role in roles"
              :key="role.id"
              class="border-l border-surface text-center px-3 py-2 text-base cursor-pointer hover:bg-emphasis"
              :class="{ 'bg-primary/5': role.id === preselectId }"
              @click="cycle(role.id, op.op)"
            >
              <i
                v-if="access(role.id, op.op) === 'allow'"
                class="pi pi-check text-primary"
              />
              <i
                v-else-if="access(role.id, op.op) === 'deny'"
                class="pi pi-times text-rose-500"
              />
              <i v-else class="pi pi-minus text-muted-color/40" />
            </td>
          </tr>
        </tbody>
      </table>
      <p class="text-xs text-muted-color mt-3">{{ $t('project.permissions.dialog.legend') }}</p>
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="open = false"
      />
      <Button
        :label="$t('project.permissions.dialog.save')"
        size="small"
        :loading="saving"
        :disabled="!dirty.size"
        @click="save"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  // Full RBAC resource string, e.g. corteza::compose:module/123/456 or .../123/*
  resource: { type: String, default: '' },
  title: { type: String, default: '' },
  // Role to highlight on open.
  preselectId: { type: String, default: '' },
  // Project roles to show as columns: [{ id, name }].
  roles: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue', 'saved'])

const { t, te } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $toast = inject('$toast')

const open = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const header = computed(() =>
  t('project.permissions.dialog.header', { target: props.title || t('project.permissions.dialog.resource') }),
)

// Route to the API that owns this resource's component.
const api = computed(() => {
  const m = props.resource.match(/^corteza::(\w+):/)
  switch (m?.[1]) {
    case 'compose':
      return $ComposeAPI
    case 'automation':
      return $AutomationAPI
    default:
      return $SystemAPI
  }
})
const resourceType = computed(() => props.resource.split('/', 1)[0])

// i18n prefix for operation labels, e.g. "compose.module".
const i18nPrefix = computed(() => {
  const parts = resourceType.value.split(':').filter(Boolean)
  return `${parts[1] || ''}.${parts[2] || 'component'}`
})
function humanize(s) {
  return s.replace(/[-.]/g, ' ').replace(/^./, c => c.toUpperCase())
}
function opLabel(op) {
  const k = `permissions.resources.${i18nPrefix.value}.operations.${op}.title`
  return te(k) ? t(k) : humanize(op)
}
function opDesc(op) {
  const k = `permissions.resources.${i18nPrefix.value}.operations.${op}.description`
  return te(k) ? t(k) : ''
}

const loading = ref(false)
const saving = ref(false)
const ops = ref([]) // [{ op, label, description }]
// roleId -> { op: 'allow'|'deny'|'inherit' }
const grants = ref({})
const dirty = ref(new Set())

function access(roleId, op) {
  return grants.value[roleId]?.[op] || 'inherit'
}

// Click cycles allow → inherit (two-state; deny is shown if the server already
// has it but isn't set by clicking).
function cycle(roleId, op) {
  const cur = access(roleId, op)
  const next = cur === 'allow' ? 'inherit' : 'allow'
  grants.value = {
    ...grants.value,
    [roleId]: { ...(grants.value[roleId] || {}), [op]: next },
  }
  dirty.value = new Set(dirty.value).add(roleId)
}

async function load() {
  if (!props.resource || !api.value || !props.roles.length) return
  loading.value = true
  dirty.value = new Set()
  try {
    const catalog = await api.value.permissionsList().catch(() => [])
    ops.value = (catalog || [])
      .filter(p => p.type === resourceType.value)
      .map(p => ({ op: p.op, label: opLabel(p.op), description: opDesc(p.op) }))

    const next = {}
    await Promise.all(
      props.roles.map(async role => {
        const rr = await api.value
          .permissionsRead({ roleID: role.id, resource: [props.resource] })
          .catch(() => [])
        const map = {}
        for (const r of rr || []) map[r.operation] = r.access
        next[role.id] = map
      }),
    )
    grants.value = next
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!dirty.value.size) return
  saving.value = true
  try {
    await Promise.all(
      [...dirty.value].map(roleId => {
        const map = grants.value[roleId] || {}
        const rules = ops.value.map(o => ({
          resource: props.resource,
          operation: o.op,
          access: map[o.op] || 'inherit',
        }))
        return api.value.permissionsUpdate({ roleID: roleId, rules })
      }),
    )
    dirty.value = new Set()
    $toast.toastSuccess(t('project.permissions.toastSaved'))
    emit('saved')
    open.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.permissions.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}

watch(
  () => [props.modelValue, props.resource],
  ([visible]) => {
    if (visible) load()
  },
  { immediate: true },
)
</script>
