<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="route_"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:apigw-route/${route_.routeID}`"
          :title="route_.endpoint || route_.routeID"
          :target="route_.endpoint || route_.routeID"
        />
      </div>
      <!-- Route info panel -->
      <Panel :header="$t('system.apigw.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="endpoint" class="flex flex-col gap-2 md:col-span-2">
            <label for="endpoint" class="font-medium text-primary">
              {{ $t('system.apigw.editor.info.endpoint') }}
              <span class="text-red-500">*</span>
            </label>
            <InputText
              id="endpoint"
              name="endpoint"
              v-model="route_.endpoint"
              placeholder="/api/v1/..."
            />
            <Message
              v-if="$form.endpoint?.invalid"
              severity="error"
              size="small"
              variant="simple"
            >
              {{ $form.endpoint.error?.message }}
            </Message>
          </FormField>

          <div class="flex flex-col gap-2">
            <label for="method" class="font-medium text-primary">
              {{ $t('system.apigw.editor.info.method') }}
            </label>
            <Select
              id="method"
              v-model="route_.method"
              :options="methodOptions"
              option-label="label"
              option-value="value"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label for="group" class="font-medium text-primary">
              {{ $t('system.apigw.editor.info.group') }}
            </label>
            <InputText id="group" v-model="route_.group" />
          </div>

          <div class="flex flex-col gap-2 md:col-span-2">
            <label for="description" class="font-medium text-primary">
              {{ $t('system.apigw.editor.info.description') }}
            </label>
            <Textarea id="description" v-model="route_.meta.description" rows="2" />
          </div>

          <div class="flex items-center gap-3">
            <ToggleSwitch id="enabled" v-model="route_.enabled" />
            <label for="enabled" class="font-medium text-primary cursor-pointer">
              {{ $t('system.apigw.editor.info.enabled') }}
            </label>
          </div>

          <div class="flex items-center gap-3">
            <ToggleSwitch id="async" v-model="route_.meta.async" />
            <label for="async" class="font-medium text-primary cursor-pointer">
              {{ $t('system.apigw.editor.info.async') }}
            </label>
          </div>
        </div>
      </Panel>

      <!-- Filters panel -->
      <Panel v-if="isEdit" :header="$t('system.apigw.editor.filters.title')" toggleable :collapsed="false">
        <!-- Step tabs: Prefilter / Processer / Postfilter -->
        <div class="flex border-b border-surface -mx-5 -mt-2 mb-3">
          <button
            v-for="(step, index) in steps"
            :key="step"
            type="button"
            class="px-4 py-3 text-sm font-medium transition-colors border-b-2"
            :class="activeStep === index
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-color hover:text-color hover:border-surface-300'"
            @click="activeStep = index"
          >
            {{ $t(`system.apigw.editor.filters.step_title.${step}`) }}
          </button>
        </div>

        <!-- Add filter dropdown for current step -->
        <div class="mb-3">
          <Select
            v-model="selectedFilterToAdd"
            :options="availableFiltersForStep"
            option-label="label"
            :placeholder="$t('system.apigw.editor.filters.addFilter')"
            showClear
            class="w-full md:w-72"
            @change="onAddFilter"
          >
            <template #option="{ option }">
              <span :class="{ 'text-muted-color': option.disabled }">
                {{ option.label }}
              </span>
            </template>
          </Select>
        </div>

        <!-- Filter list for current step -->
        <div v-if="filtersLoading" class="flex items-center justify-center py-8">
          <ProgressSpinner style="width: 30px; height: 30px" />
        </div>

        <div v-else-if="filtersForStep.length === 0" class="text-center text-muted-color py-8">
          {{ $t('system.apigw.editor.filters.list.noFilters') }}
        </div>

        <div v-else class="flex flex-col -mx-5">
          <div
            v-for="(filter, index) in filtersForStep"
            :key="filter.ref || filter.filterID || index"
            class="flex items-center gap-3 px-5 py-3 border-b border-surface hover:bg-surface-50 dark:hover:bg-surface-800 transition-colors cursor-pointer group"
            @click="openFilterModal(filter)"
          >
            <i class="pi pi-bars text-muted-color cursor-grab" />
            <span class="flex-1 font-medium text-sm">
              {{ filter.label || filter.ref || filter.kind }}
            </span>
            <Tag
              :value="filter.enabled ? $t('system.apigw.editor.filters.enabled') : $t('system.apigw.editor.filters.disabled')"
              :severity="filter.enabled ? 'success' : 'secondary'"
              class="text-xs"
            />
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              class="opacity-0 group-hover:opacity-100 transition-opacity"
              @click.stop="onRemoveFilter(filter)"
            />
          </div>
        </div>

        <!-- Save filters button -->
        <div class="flex justify-end mt-3">
          <Button
            :label="$t('general.label.save')"
            icon="pi pi-save"
            size="small"
            :loading="savingFilters"
            :disabled="!hasFilterChanges"
            @click="onFiltersSubmit"
          />
        </div>
      </Panel>
    </div>

    <!-- Filter config modal -->
    <Dialog
      v-model:visible="filterModalVisible"
      modal
      :header="modalFilter?.label || $t('system.apigw.editor.filters.edit')"
      :style="{ width: '40vw' }"
      :breakpoints="{ '1199px': '60vw', '575px': '90vw' }"
    >
      <div v-if="modalFilter" class="flex flex-col gap-4">
        <CFilterParamsEditor :filter="modalFilter" />

        <div class="flex items-center gap-3">
          <ToggleSwitch id="modalFilterEnabled" v-model="modalFilter.enabled" />
          <label for="modalFilterEnabled" class="font-medium text-primary cursor-pointer">
            {{ $t('system.apigw.editor.filters.enabled') }}
          </label>
        </div>

        <div class="flex justify-end gap-2 pt-2 border-t border-surface">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            size="small"
            outlined
            @click="filterModalVisible = false"
          />
          <Button
            :label="$t('general.label.save')"
            size="small"
            @click="onFilterModalSave"
          />
        </div>
      </div>
    </Dialog>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.apiGateway' })"
        />
        <div class="flex gap-2">

          <CInputDelete
            v-if="isEdit && route_.canDeleteApigwRoute"
            :label="$t('system.apigw.editor.delete')"
            :message="$t('general.confirm.delete')"
            :header="route_.endpoint || route_.routeID"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useConfirmDelete, useUnsavedGuard } from '@planetcrust/human-vue'
import { NoID } from '@planetcrust/human-js'
import { cloneDeep, isEqual } from 'lodash-es'
import CFilterParamsEditor from '@/components/ApiGateway/CFilterParamsEditor.vue'

const { CInputDelete } = components

const vueRoute = useRoute()
const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const route_ = ref(null)
const initialRoute_ = ref(null)

// Filter state
const activeStep = ref(0)
const filtersLoading = ref(false)
const savingFilters = ref(false)
const filters = ref([])
const availableFilters = ref([])
const filterModalVisible = ref(false)
const modalFilter = ref(null)
const selectedFilterToAdd = ref(null)
const hasFilterChanges = ref(false)

const steps = ['prefilter', 'processer', 'postfilter']

const mapKindToStep = {
  prefilter: 0,
  processer: 1,
  postfilter: 2,
}

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
]

const isEdit = computed(() => !!vueRoute.params.routeID)

const pageTitle = computed(() =>
  isEdit.value ? t('system.apigw.editor.title.edit') : t('system.apigw.editor.title.create'),
)

const initialValues = computed(() => ({
  endpoint: route_.value?.endpoint || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}
  if (!values.endpoint || values.endpoint.trim().length === 0) {
    errors.endpoint = [{ message: t('general.label.required') }]
  }
  return { errors }
})

// Filters for the currently active step
const filtersForStep = computed(() =>
  filters.value
    .filter(f => mapKindToStep[f.kind] === activeStep.value && !f.deleted)
    .sort((a, b) => (a.weight || 0) - (b.weight || 0)),
)

// Available filter definitions for the active step, marking already-used ones as disabled
const availableFiltersForStep = computed(() =>
  availableFilters.value
    .filter(f => mapKindToStep[f.kind] === activeStep.value)
    .map(f => ({
      ...f,
      disabled: filters.value.some(ef => ef.ref === f.ref && !ef.deleted),
    })),
)

// --- Data fetching ---

function newRoute() {
  return {
    endpoint: '',
    method: 'GET',
    enabled: true,
    group: '',
    meta: { description: '', async: false },
  }
}

async function loadRoute() {
  const routeID = vueRoute.params.routeID
  if (!routeID) {
    route_.value = newRoute()
    initialRoute_.value = cloneDeep(route_.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.apigwRouteRead({ routeID })
    route_.value = {
      ...raw,
      meta: { description: raw.meta?.description || '', async: raw.meta?.async || false },
    }
    initialRoute_.value = cloneDeep(route_.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.fetch.error'))(e)
    router.push({ name: 'system.apiGateway' })
  } finally {
    loading.value = false
  }
}

async function fetchAvailableFilters() {
  try {
    const defs = await $SystemAPI.apigwFilterDefFilter({})
    availableFilters.value = (defs || []).map(f => ({
      ...f,
      ref: f.name,
      enabled: true,
    }))
  } catch (e) {
    console.error('Failed to fetch filter definitions:', e)
  }
}

async function fetchFilters() {
  if (!vueRoute.params.routeID) return

  filtersLoading.value = true
  try {
    const result = await $SystemAPI.apigwFilterList({ routeID: vueRoute.params.routeID })
    const routeFilters = result?.set || []

    // Map route filters to their definitions for label + param structure
    filters.value = routeFilters.map(filter => {
      const def = availableFilters.value.find(af => af.ref === filter.ref) || {}
      return {
        ...def,
        ...filter,
        label: def.label || filter.ref || filter.kind,
        params: decodeParams(def, filter.params || {}),
        weight: parseInt(filter.weight) || 0,
        enabled: !!filter.enabled,
      }
    })
    hasFilterChanges.value = false
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.filter.fetch.error'))(e)
  } finally {
    filtersLoading.value = false
  }
}

function decodeParams(filterDef = {}, values = {}) {
  const { params = [] } = filterDef
  return params.map(({ label, type }) => ({
    label,
    type,
    value: values[label],
  }))
}

function encodeParams(params = []) {
  return params.reduce((result, p) => {
    result[p.label] = p.value
    return result
  }, {})
}

// --- Filter CRUD ---

function onAddFilter(event) {
  const filterDef = event.value
  if (!filterDef || filterDef.disabled) {
    selectedFilterToAdd.value = null
    return
  }

  // Check if already exists
  const existing = filters.value.find(f => f.ref === filterDef.ref && !f.deleted)
  if (existing) {
    openFilterModal(existing)
  } else {
    const newFilter = {
      ...filterDef,
      created: true,
      weight: filtersForStep.value.length,
      params: filterDef.params
        ? filterDef.params.map(p => ({ ...p, value: undefined, options: { ...p.options } }))
        : [],
    }
    openFilterModal(newFilter)
  }

  selectedFilterToAdd.value = null
}

function openFilterModal(filter) {
  // Convert response filter params for FE structure
  const modalData = {
    ...filter,
    params: (filter.params || []).map(p => {
      let value = p.value

      if (filter.ref === 'response') {
        if (p.type === 'header' && value) {
          value = Object.entries(value).map(([name, v = []]) => ({ name, expr: Array.isArray(v) ? v.join('') : v }))
        } else if (p.type === 'input') {
          value = { type: 'Any', expr: '', ...value }
        }
      }

      return { ...p, value }
    }),
  }

  modalFilter.value = modalData
  filterModalVisible.value = true
}

function onFilterModalSave() {
  // Convert response filter params back to BE structure
  const filter = {
    ...modalFilter.value,
    params: (modalFilter.value.params || []).map(p => {
      if (modalFilter.value.ref === 'response' && p.type === 'header' && Array.isArray(p.value)) {
        return {
          ...p,
          value: p.value.reduce((obj, { name, expr = '' }) => {
            return { ...obj, [name]: [expr] }
          }, {}),
        }
      }
      return p
    }),
    updated: true,
  }

  // Update or add to filters array
  const existingIndex = filters.value.findIndex(f => f.ref === filter.ref && !f.deleted)
  if (existingIndex >= 0) {
    filters.value.splice(existingIndex, 1, filter)
  } else {
    filters.value.push(filter)
  }

  hasFilterChanges.value = true
  filterModalVisible.value = false
}

function onRemoveFilter(filter) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: filter.label || filter.ref,
    onConfirm: () => {
      if (filter.filterID) {
        // Mark existing filter for deletion
        const idx = filters.value.findIndex(f => f.filterID === filter.filterID)
        if (idx >= 0) {
          filters.value.splice(idx, 1, { ...filter, deleted: true })
        }
      } else {
        // Remove unsaved filter entirely
        const idx = filters.value.findIndex(f => f.ref === filter.ref && !f.deleted)
        if (idx >= 0) {
          filters.value.splice(idx, 1)
        }
      }
      hasFilterChanges.value = true
    },
  })
}

async function onFiltersSubmit() {
  savingFilters.value = true
  try {
    const routeID = vueRoute.params.routeID

    await Promise.all(
      filters.value
        .filter(f => f.created || f.updated || f.deleted)
        .map(async (filter) => {
          const payload = {
            routeID,
            ref: filter.ref,
            kind: filter.kind,
            weight: String(filter.weight || 0),
            enabled: filter.enabled,
            params: encodeParams(filter.params || []),
          }

          if (filter.filterID && filter.filterID !== NoID) {
            if (filter.deleted) {
              return $SystemAPI.apigwFilterDelete({ filterID: filter.filterID })
            }
            payload.filterID = filter.filterID
            return $SystemAPI.apigwFilterUpdate(payload)
          } else if (!filter.deleted) {
            return $SystemAPI.apigwFilterCreate(payload)
          }
        }),
    )

    $toast.toastSuccess(t('notification.gateway.filter.update.success'))
    await fetchFilters()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.filter.update.error'))(e)
  } finally {
    savingFilters.value = false
  }
}

// --- Route CRUD ---

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      endpoint: route_.value.endpoint,
      method: route_.value.method,
      enabled: route_.value.enabled,
      group: route_.value.group || '',
      meta: route_.value.meta,
    }

    if (isEdit.value) {
      payload.routeID = route_.value.routeID
      const raw = await $SystemAPI.apigwRouteUpdate(payload)
      route_.value = {
        ...raw,
        meta: { description: raw.meta?.description || '', async: raw.meta?.async || false },
      }
      initialRoute_.value = cloneDeep(route_.value)
      $toast.toastSuccess(t('notification.gateway.update.success'))
    } else {
      const created = await $SystemAPI.apigwRouteCreate(payload)
      $toast.toastSuccess(t('notification.gateway.create.success'))
      markSaved()
      router.push({ name: 'system.apiGateway.edit', params: { routeID: created.routeID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(t(`notification.gateway.${isEdit.value ? 'update' : 'create'}.error`))(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.apigwRouteDelete({ routeID: route_.value.routeID })
    $toast.toastSuccess(t('notification.gateway.delete.success'))
    router.push({ name: 'system.apiGateway' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!route_.value && !!initialRoute_.value && !isEqual(route_.value, initialRoute_.value),
  messageKey: 'general.editor.unsavedChanges',
})

// --- Lifecycle ---

async function init() {
  await loadRoute()
  if (isEdit.value) {
    await fetchAvailableFilters()
    await fetchFilters()
  }
}

onMounted(() => init())
watch(
  () => vueRoute.params.routeID,
  () => init(),
)
</script>
