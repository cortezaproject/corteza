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
          <CFormGroup
            name="endpoint"
            :label="$t('system.apigw.editor.info.endpoint')"
            required
            class="md:col-span-2"
          >
            <InputText
              id="endpoint"
              name="endpoint"
              v-model="route_.endpoint"
              placeholder="/api/v1/..."
            />
          </CFormGroup>

          <CFormGroup :label="$t('system.apigw.editor.info.method')" input-id="method">
            <Select
              id="method"
              v-model="route_.method"
              :options="methodOptions"
              option-label="label"
              option-value="value"
            />
          </CFormGroup>

          <CInputToggleCard
            v-model="route_.enabled"
            :label="$t('system.apigw.editor.info.enabled')"
            :description="$t('system.apigw.editor.info.enabledDescription')"
            class="self-end"
          />

          <CFormGroup
            :label="$t('system.apigw.editor.info.description')"
            input-id="description"
            class="md:col-span-2"
          >
            <Textarea id="description" v-model="route_.meta.description" rows="2" />
          </CFormGroup>

          <CInputToggleCard
            v-model="route_.meta.async"
            :label="$t('system.apigw.editor.info.async')"
            :description="$t('system.apigw.editor.info.asyncDescription')"
            class="self-start"
          />
        </div>
      </Panel>

      <!-- Filters panel -->
      <Panel
        v-if="isEdit"
        :header="$t('system.apigw.editor.filters.title')"
        toggleable
        :collapsed="false"
        :pt="{ content: { style: 'padding: 0 !important' } }"
      >
        <!-- Step tabs: Prefilter / Processer / Postfilter -->
        <Tabs v-model:value="activeStep">
          <TabList>
            <Tab v-for="step in steps" :key="step" :value="step">
              {{ $t(`system.apigw.editor.filters.step_title.${step}`) }}
            </Tab>
          </TabList>
        </Tabs>

        <div class="p-3">
          <!-- Add filter dropdown for current step -->
          <div class="mb-3">
            <Button
              :label="$t('system.apigw.editor.filters.addFilter')"
              icon="pi pi-plus"
              icon-pos="right"
              size="small"
              :disabled="!addFilterMenuItems.length"
              @click="toggleAddFilterMenu"
            />
            <Menu ref="addFilterMenu" :model="addFilterMenuItems" popup />
          </div>

          <!-- Filter list for current step -->
          <CFormItemList
            :items="filtersForStep"
            :loading="filtersLoading"
            :empty-message="$t('system.apigw.editor.filters.list.noFilters')"
            item-key="ref"
            draggable
            @remove="filter => onRemoveFilter(filter)"
            @reorder="onReorderFilters"
          >
            <template #default="{ item: filter }">
              <button
                type="button"
                class="flex items-center gap-3 w-full text-left bg-transparent border-0 p-0 cursor-pointer"
                @click="openFilterModal(filter)"
              >
                <span>
                  {{ filter.label || filter.ref || filter.kind }}
                </span>
                <Tag
                  :value="
                    filter.enabled
                      ? $t('system.apigw.editor.filters.enabled')
                      : $t('system.apigw.editor.filters.disabled')
                  "
                  :severity="filter.enabled ? 'success' : 'secondary'"
                  class="text-xs"
                />
              </button>
            </template>
          </CFormItemList>
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
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            outlined
            size="small"
            @click="filterModalVisible = false"
          />
          <Button
            :label="$t('general.label.save')"
            size="small"
            @click="onFilterModalSave"
          />
        </div>
      </template>
    </Dialog>

    <CEditorActions :back-to="{ name: 'system.apiGateway' }">
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
    </CEditorActions>
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

const { CInputDelete, CInputToggleCard } = components

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
const steps = ['prefilter', 'processer', 'postfilter']
const activeStep = ref(steps[0])
const filtersLoading = ref(false)
const filters = ref([])
const availableFilters = ref([])
const filterModalVisible = ref(false)
const modalFilter = ref(null)
const hasFilterChanges = ref(false)
const addFilterMenu = ref(null)

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
    .filter(f => f.kind === activeStep.value && !f.deleted)
    .sort((a, b) => (a.weight || 0) - (b.weight || 0)),
)

// Available filter definitions for the active step, marking already-used ones as disabled
const availableFiltersForStep = computed(() =>
  availableFilters.value
    .filter(f => f.kind === activeStep.value)
    .map(f => ({
      ...f,
      disabled: filters.value.some(ef => ef.ref === f.ref && !ef.deleted),
    })),
)

// Menu items for the "Add filter" dropdown button
const addFilterMenuItems = computed(() =>
  availableFiltersForStep.value.map(f => ({
    label: f.label,
    disabled: f.disabled,
    command: () => onAddFilter(f),
  })),
)

function toggleAddFilterMenu(event) {
  addFilterMenu.value?.toggle(event)
}

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

    // Map route filters to their definitions: definition wins for kind/label/params shape,
    // only stored fields (params values, weight, filterID, enabled) carry over from the route filter.
    filters.value = routeFilters.map(filter => {
      const def = availableFilters.value.find(af => af.ref === filter.ref) || {}
      return {
        ...def,
        params: decodeParams(def, filter.params || {}),
        weight: parseInt(filter.weight) || 0,
        filterID: filter.filterID,
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

function onAddFilter(filterDef) {
  if (!filterDef || filterDef.disabled) return

  // Check if already exists
  const existing = filters.value.find(f => f.ref === filterDef.ref && !f.deleted)
  if (existing) {
    openFilterModal(existing)
    return
  }

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

function openFilterModal(filter) {
  // Convert response filter params for FE structure
  const modalData = {
    ...filter,
    params: (filter.params || []).map(p => {
      let value = p.value

      if (filter.ref === 'response') {
        if (p.type === 'header' && value) {
          value = Object.entries(value).map(([name, v = []]) => ({
            name,
            expr: Array.isArray(v) ? v.join('') : v,
          }))
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

function onReorderFilters(reordered) {
  reordered.forEach((filter, idx) => {
    const original = filters.value.find(f => f.ref === filter.ref && !f.deleted)
    if (original && original.weight !== idx) {
      original.weight = idx
      original.updated = true
    }
  })
  hasFilterChanges.value = true
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

async function saveFilters(routeID) {
  await Promise.all(
    filters.value
      .filter(f => f.created || f.updated || f.deleted)
      .map(async filter => {
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

  await fetchFilters()
}

// --- Route CRUD ---

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      endpoint: route_.value.endpoint,
      method: route_.value.method,
      enabled: route_.value.enabled,
      group: route_.value.group || NoID,
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
      if (hasFilterChanges.value) {
        await saveFilters(route_.value.routeID)
      }
      $toast.toastSuccess(t('notification.gateway.update.success'))
    } else {
      const created = await $SystemAPI.apigwRouteCreate(payload)
      $toast.toastSuccess(t('notification.gateway.create.success'))
      markSaved()
      router.push({ name: 'system.apiGateway.edit', params: { routeID: created.routeID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(t(`notification.gateway.${isEdit.value ? 'update' : 'create'}.error`))(
      e,
    )
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
  isDirty: () =>
    !saving.value &&
    !deleting.value &&
    !!route_.value &&
    !!initialRoute_.value &&
    (!isEqual(route_.value, initialRoute_.value) || hasFilterChanges.value),
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
