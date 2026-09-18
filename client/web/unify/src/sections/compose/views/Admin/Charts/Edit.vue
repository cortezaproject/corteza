<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ isEdit ? chart?.name || $t('chart.edit.title') : $t('chart.edit.title') }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <ButtonGroup v-if="isEdit && chart && namespace" class="gap-1">
      <ChartTranslator :chart="chart" :namespace="namespace" @update:chart="chart = $event" />
    </ButtonGroup>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="chart"
    :resolver="resolver"
    :initial-values="formValues"
    class="flex flex-col h-full"
    @submit="handleSubmit"
  >
    <!-- The scroller is full width; the container only centres what is inside
         it. Putting overflow on the centred element leaves the gutters outside
         it, where a wheel over them scrolls nothing. -->
    <div class="flex-1 overflow-auto">
      <div class="container mx-auto p-4 min-w-0">
        <!-- Actions above cards -->
        <div v-if="isEdit" class="flex justify-end gap-2 mb-4">
          <Button
            v-if="namespace?.canExportCharts"
            :label="$t('general.label.export')"
            icon="pi pi-download"
            size="small"
            severity="secondary"
            outlined
            @click="exportChart"
          />
          <CPermissionsButton
            v-if="chart.canGrant"
            :resource="`corteza::compose:chart/${namespace.namespaceID}/${chart.chartID}`"
            :title="chart.name || chart.handle || chart.chartID"
            :target="chart.name || chart.handle || chart.chartID"
            v-tooltip.bottom="$t('general.label.permissions')"
            outlined
          />
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
          <!-- Left column: Settings -->
          <div class="lg:col-span-7 flex flex-col gap-4">
            <Panel :header="$t('chart.generalSettings')" toggleable>
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <CFormGroup name="name" :label="$t('chart.name')" required>
                  <InputText
                    id="name"
                    name="name"
                    v-model="chart.name"
                    :placeholder="$t('chart.general.placeholder.name')"
                    class="w-full"
                  />
                </CFormGroup>

                <CFormGroup name="handle" :label="$t('chart.handle')">
                  <InputText
                    id="handle"
                    name="handle"
                    v-model="chart.handle"
                    :placeholder="$t('chart.general.placeholder.handle')"
                    class="w-full"
                  />
                </CFormGroup>

                <CFormGroup
                  :label="$t('chart.description')"
                  input-id="description"
                  class="lg:col-span-2"
                >
                  <Textarea
                    id="description"
                    v-model="chart.meta.description"
                    :placeholder="$t('chart.general.placeholder.description')"
                    rows="2"
                    autoResize
                    class="w-full"
                  />
                </CFormGroup>

                <CFormGroup :label="$t('chart.colorScheme.label')" input-id="colorScheme">
                  <div class="flex items-center gap-2">
                    <Select
                      id="colorScheme"
                      v-model="chart.config.colorScheme"
                      :options="colorSchemes"
                      option-label="name"
                      option-value="id"
                      :placeholder="$t('chart.colorScheme.placeholder')"
                      class="flex-1 min-w-0"
                      filter
                      show-clear
                    >
                      <template #value="{ value }">
                        <span v-if="value" class="truncate">{{ getSchemeName(value) }}</span>
                        <span v-else>{{ $t('chart.colorScheme.placeholder') }}</span>
                      </template>
                      <template #option="{ option }">
                        <div class="flex items-center gap-2 min-w-0">
                          <div class="flex gap-0.5 items-center shrink-0">
                            <div
                              v-for="(color, ci) in option.colors"
                              :key="ci"
                              :style="`background: ${color};`"
                              class="w-3.5 h-3.5 rounded-sm"
                            />
                          </div>
                          <span class="truncate">{{ option.name }}</span>
                        </div>
                      </template>
                      <template v-if="canManageColorSchemes" #header>
                        <div class="p-2 border-b border-surface-200 dark:border-surface-700">
                          <Button
                            :label="$t('chart.colorScheme.custom.add')"
                            icon="pi pi-plus"
                            severity="secondary"
                            text
                            size="small"
                            class="w-full"
                            data-testid="chart-color-scheme-add"
                            @click="createColorScheme"
                          />
                        </div>
                      </template>
                    </Select>

                    <Button
                      v-if="showEditColorSchemeButton"
                      v-tooltip.bottom="$t('chart.colorScheme.custom.edit')"
                      :aria-label="$t('chart.colorScheme.custom.edit')"
                      icon="pi pi-pencil"
                      severity="secondary"
                      outlined
                      data-testid="chart-color-scheme-edit"
                      @click="editColorScheme"
                    />
                  </div>

                  <div v-if="currentColorScheme" class="flex flex-wrap gap-0.5 items-center mt-2">
                    <div
                      v-for="(color, ci) in currentColorScheme.colors"
                      :key="ci"
                      :style="`background: ${color};`"
                      class="w-4 h-2 rounded-sm"
                    />
                  </div>
                </CFormGroup>

                <CInputToggleCard
                  v-model="animationEnabled"
                  :label="$t('chart.edit.animation.label')"
                  :description="$t('chart.edit.animation.description')"
                  class="self-start"
                />
              </div>
            </Panel>

            <!-- Report editor: contributes its own panels -->
            <component
              :is="reportEditor"
              v-if="chart && editReport"
              :chart="chart"
              :modules="modules"
            />

            <Panel :header="$t('chart.edit.toolbox.label')" toggleable collapsed>
              <CInputToggleCard
                v-model="saveAsImageEnabled"
                :label="$t('chart.edit.toolbox.saveAsImage.label')"
                :description="$t('chart.edit.toolbox.saveAsImage.description')"
              />
            </Panel>
          </div>

          <!-- Right column: Preview -->
          <div class="lg:col-span-5">
            <div class="sticky top-0">
              <Card>
                <template #content>
                  <div class="relative" style="height: 400px">
                    <div
                      v-if="previewNeedsRecord"
                      class="absolute inset-0 flex items-center justify-center p-3 text-center text-muted-color text-sm"
                    >
                      {{ $t('chart.edit.filter.previewNeedsRecord') }}
                    </div>

                    <ChartRenderer
                      v-else-if="chart"
                      ref="chartPreview"
                      :chart="chart"
                      :reporter="reporter"
                      @updated="onUpdated"
                    />
                  </div>
                </template>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>

    <CEditorActions
      :back-to="true"
      @back="goBack({ name: 'admin.charts', params: { slug: route.params.slug } })"
    >
      <CInputDelete
        v-if="isEdit && chart.canDeleteChart"
        :label="$t('general.label.delete')"
        :message="$t('chart.list.delete')"
        :header="chart.name"
        :disabled="processingDelete"
        @confirm="handleDelete"
      />
      <Button
        v-if="isEdit"
        :label="$t('general.label.saveAsCopy')"
        icon="pi pi-copy"
        severity="secondary"
        :loading="processingClone"
        :disabled="!isValid"
        @click="handleClone"
      />
      <Button
        v-if="!hideSave"
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="processingSave"
      />
    </CEditorActions>
  </Form>

  <Dialog
    v-model:visible="colorSchemeModal.show"
    :header="colorSchemeModalTitle"
    modal
    :style="{ width: '32rem' }"
  >
    <div class="flex flex-col gap-4">
      <CFormGroup
        :label="$t('chart.colorScheme.custom.modal.name.label')"
        input-id="colorSchemeName"
      >
        <InputText
          id="colorSchemeName"
          v-model="colorSchemeModal.colorscheme.name"
          class="w-full"
          autofocus
        />
      </CFormGroup>

      <CFormGroup :label="$t('chart.colorScheme.custom.modal.colors.label')">
        <div class="flex flex-wrap items-center gap-1">
          <div
            v-for="(color, index) in colorSchemeModal.colorscheme.colors"
            :key="index"
            class="flex items-center"
          >
            <CInputColorPicker v-model="colorSchemeModal.colorscheme.colors[index]" />
            <Button
              v-tooltip.bottom="$t('general.label.remove')"
              :aria-label="$t('general.label.remove')"
              :disabled="colorSchemeModal.colorscheme.colors.length < 2"
              icon="pi pi-times"
              severity="danger"
              text
              rounded
              size="small"
              @click="removeColor(index)"
            />
          </div>

          <Button
            v-tooltip.bottom="$t('general.label.add')"
            :aria-label="$t('general.label.add')"
            icon="pi pi-plus"
            severity="secondary"
            outlined
            size="small"
            data-testid="chart-color-scheme-add-color"
            @click="addColor"
          />
        </div>
      </CFormGroup>
    </div>

    <template #footer>
      <CInputDelete
        v-if="colorSchemeModal.edit"
        :label="$t('general.label.delete')"
        :message="$t('chart.colorScheme.custom.modal.delete')"
        :header="colorSchemeModal.colorscheme.name"
        :disabled="colorSchemeModal.processing"
        outlined
        class="mr-auto"
        @confirm="deleteColorScheme"
      />
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        :disabled="colorSchemeModal.processing"
        @click="colorSchemeModal.show = false"
      />
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :disabled="!colorSchemeIsSaveable"
        :loading="colorSchemeModal.processing"
        data-testid="chart-color-scheme-save"
        @click="saveColorScheme"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, computed, watch, inject, provide, toRaw, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { compose, shared } from '@planetcrust/human-js'
import { components, useHistoryBack, useDraftGuard, useRBACStore } from '@planetcrust/human-vue'
import { chartConstructor } from '../../../lib/charts'
import {
  COLOR_SCHEMES_SETTING,
  colorSchemeOptions,
  isCustomScheme,
  newCustomScheme,
  readColorSchemes,
  removeColorScheme,
  upsertColorScheme,
} from '../../../lib/chart-color-schemes'
import { useChartStore } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import ChartRenderer from '../../../components/Chart/ChartRenderer.vue'
import { evaluatePlacementFilter, usesRecordVariables } from '../../../lib/record-filter'
import ChartTranslator from '../../../components/Admin/Chart/ChartTranslator.vue'
import * as Reports from '../../../components/Chart/Report/index.js'

const { CInputDelete, CInputColorPicker } = components

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const goBack = useHistoryBack()
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings', undefined)
const $toast = inject('$toast')
const $Auth = inject('$Auth', {})

const chartStore = useChartStore()
const moduleStore = useModuleStore()
const rbac = useRBACStore()

const { colorschemes } = shared

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

// State
const chart = ref(null)
const { capture, markSaved } = useDraftGuard({
  draft: chart,
  busy: () => processingSave.value || processingDelete.value || processingClone.value,
})

const loading = ref(false)
const processing = ref(false)
const processingSave = ref(false)
const processingClone = ref(false)
const processingDelete = ref(false)
const editReportIndex = ref(0)
const chartPreview = ref(null)

const colorSchemeModal = ref({
  show: false,
  processing: false,
  edit: false,
  colorscheme: newCustomScheme(),
})

// Computed
const chartID = computed(() => route.params.chartID)

const isEdit = computed(() => {
  return chart.value && chart.value.chartID && chart.value.chartID !== '0'
})

const modules = computed(() => moduleStore.set || [])

// The instance's own schemes are a live read of the setting, so a save in the
// modal reaches the picker and the preview without refetching the chart.
const customColorSchemes = computed(() => readColorSchemes($Settings))

const colorSchemes = computed(() =>
  colorSchemeOptions(customColorSchemes.value, colorschemes, count =>
    t('chart.colorLabel', { count }),
  ),
)

const currentColorScheme = computed(() =>
  colorSchemes.value.find(({ id }) => id === chart.value?.config?.colorScheme),
)

const canManageColorSchemes = computed(() => rbac.can('system/', 'settings.manage'))

const showEditColorSchemeButton = computed(
  () => canManageColorSchemes.value && isCustomScheme(currentColorScheme.value?.id),
)

const colorSchemeModalTitle = computed(() =>
  t(`chart.colorScheme.custom.${colorSchemeModal.value.edit ? 'edit' : 'add'}`),
)

const colorSchemeIsSaveable = computed(() => {
  const { name, colors } = colorSchemeModal.value.colorscheme
  return !!name?.trim() && !!colors?.length && !colorSchemeModal.value.processing
})

function getSchemeName(id) {
  const scheme = colorSchemes.value.find(s => s.id === id)
  return scheme?.name || id
}

const editReport = computed({
  get: () => {
    if (chart.value && editReportIndex.value !== undefined) {
      return chart.value.config?.reports?.[editReportIndex.value]
    }
    return undefined
  },
  set: v => {
    if (chart.value?.config?.reports) {
      chart.value.config.reports.splice(editReportIndex.value, 1, v)
    }
  },
})

provide('reportDraft', editReport)

const reportEditor = computed(() => {
  if (!chart.value) return undefined

  if (chart.value instanceof compose.FunnelChart) return markRaw(Reports.FunnelChart)
  if (chart.value instanceof compose.GaugeChart) return markRaw(Reports.GaugeChart)
  if (chart.value instanceof compose.RadarChart) return markRaw(Reports.RadarChart)

  return markRaw(Reports.GenericChart)
})

const handlePattern = /^[a-zA-Z][a-zA-Z0-9_.]*[a-zA-Z0-9]$/

// The one statement of what makes a chart saveable. The resolver runs it on
// submit and marks the offending field; save-as-copy, which never submits,
// reads the same verdict rather than keeping a second copy of the rules.
function chartErrors({ name, handle, moduleID }) {
  const errors = {}

  if (!name || !name.trim()) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (handle && !handlePattern.test(handle)) {
    errors.handle = [{ message: t('chart.general.placeholder.invalid-handle-characters') }]
  }

  // A report with no module returns nothing to plot, so the chart saves and
  // then renders empty wherever it is placed.
  if (!moduleID) {
    errors.moduleID = [{ message: t('general.label.required') }]
  }

  return errors
}

const formValues = computed(() => ({
  name: chart.value?.name || '',
  handle: chart.value?.handle || '',
  moduleID: chart.value?.config?.reports?.[0]?.moduleID || '',
}))

const resolver = ({ values }) => ({ errors: chartErrors(values) })

const isValid = computed(() => Object.keys(chartErrors(formValues.value)).length === 0)

const hideSave = computed(() => {
  return isEdit.value && !chart.value?.canUpdateChart
})

const animationEnabled = computed({
  get: () => !chart.value?.config?.noAnimation,
  set: v => {
    if (chart.value?.config) {
      chart.value.config.noAnimation = !v
    }
  },
})

const saveAsImageEnabled = computed({
  get: () => !!chart.value?.config?.toolbox?.saveAsImage,
  set: v => {
    if (!chart.value.config.toolbox) {
      chart.value.config.toolbox = {}
    }
    chart.value.config.toolbox.saveAsImage = v
  },
})

// Watch chartID to fetch on route change
watch(
  chartID,
  () => {
    fetchChart()
  },
  { immediate: true },
)

// Methods
async function fetchChart() {
  const cID = chartID.value
  const { namespaceID } = props.namespace

  if (!cID) {
    // New chart
    const category = route.query?.category || ''
    let c = new compose.Chart({ namespaceID })

    switch (category) {
      case 'gauge':
        c = new compose.GaugeChart(c)
        break
      case 'funnel':
        c = new compose.FunnelChart(c)
        break
      case 'radar':
        c = new compose.RadarChart(c)
        break
    }

    chart.value = c
    capture()
    editReportIndex.value = 0
  } else {
    loading.value = true
    processing.value = true

    try {
      const raw = await chartStore.findByID({ namespaceID, chartID: cID, force: true })
      chart.value = chartConstructor(raw)
      capture()
      editReportIndex.value = 0
    } catch (e) {
      console.error('Failed to load chart:', e)
      $toast.toastDanger(t('notification.chart.loadFailed'))
      router.push({ name: 'admin.charts' })
    } finally {
      setTimeout(() => {
        loading.value = false
        processing.value = false
      }, 300)
    }
  }
}

function reporter(r = {}) {
  const { namespaceID } = props.namespace
  let { filter } = r

  // The preview has no record, so a filter reading one cannot be evaluated
  // here; previewNeedsRecord says as much in place of the chart.
  filter = evaluatePlacementFilter(filter, { user: $Auth?.user || {} })
  if (filter === undefined && r.filter) return Promise.resolve([])

  return $ComposeAPI.recordReport({ namespaceID, ...r, filter })
}

// The preview has no record to read, and the report editor's hint already
// tells the author these resolve at placement; saying so here beats a blank
// chart that looks like an empty result set.
const previewNeedsRecord = computed(() =>
  (chart.value?.config?.reports || []).some(({ filter }) => usesRecordVariables(filter)),
)

function onUpdated() {
  processing.value = false
}

function createColorScheme() {
  colorSchemeModal.value = {
    show: true,
    processing: false,
    edit: false,
    colorscheme: newCustomScheme(),
  }
}

function editColorScheme() {
  const { id, name, colors = [] } = currentColorScheme.value

  colorSchemeModal.value = {
    show: true,
    processing: false,
    edit: true,
    colorscheme: { id, name, colors: [...colors] },
  }
}

function addColor() {
  colorSchemeModal.value.colorscheme.colors.push('#000000')
}

function removeColor(index) {
  colorSchemeModal.value.colorscheme.colors.splice(index, 1)
}

// The setting is global and the chart is not: saving a scheme persists the
// palette and selects it, and the chart itself stays unsaved until the author
// presses Save, the way every other edit on this screen behaves.
async function writeColorSchemes(value, { selects, notification }) {
  colorSchemeModal.value.processing = true

  try {
    await $SystemAPI.settingsUpdate({ values: [{ name: COLOR_SCHEMES_SETTING, value }] })
    await $Settings?.fetch()

    // Editing a scheme in place changes no part of the chart, so the preview's
    // own watcher has nothing to fire on; where the selection does change it
    // fires on its own, and asking for a second render races the first.
    const reselected = chart.value.config.colorScheme !== selects
    chart.value.config.colorScheme = selects
    colorSchemeModal.value.show = false
    $toast.toastSuccess(t(`notification.chart.colorScheme.${notification}.success`))

    if (!reselected) {
      chartPreview.value?.updateChart()
    }
  } catch (e) {
    console.error('Failed to save color scheme:', e)
    $toast.toastErrorHandler(t(`notification.chart.colorScheme.${notification}.failed`))(e)
  } finally {
    colorSchemeModal.value.processing = false
  }
}

async function saveColorScheme() {
  if (!colorSchemeIsSaveable.value) return

  const scheme = colorSchemeModal.value.colorscheme

  await writeColorSchemes(upsertColorScheme(customColorSchemes.value, scheme), {
    selects: scheme.id,
    notification: colorSchemeModal.value.edit ? 'update' : 'create',
  })
}

async function deleteColorScheme() {
  const { id } = colorSchemeModal.value.colorscheme

  await writeColorSchemes(removeColorScheme(customColorSchemes.value, id), {
    selects: chart.value?.config?.colorScheme === id ? undefined : chart.value?.config?.colorScheme,
    notification: 'delete',
  })
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    return
  }

  processingSave.value = true
  processing.value = true

  const c = toRaw(chart.value)

  try {
    if (!isEdit.value) {
      const created = await chartStore.create(c)
      chart.value = chartConstructor(created)
      capture()
      $toast.toastSuccess(t('notification.chart.created'))
      markSaved()
      router.push({ name: 'admin.charts.edit', params: { chartID: created.chartID } })
    } else {
      const updated = await chartStore.update(c)
      chart.value = chartConstructor(updated)
      capture()
      $toast.toastSuccess(t('notification.chart.updated'))
    }
  } catch (e) {
    console.error('Failed to save chart:', e)
    $toast.toastErrorHandler(t('notification.chart.updateFailed'))(e)
  } finally {
    processing.value = false
    processingSave.value = false
  }
}

async function handleClone() {
  if (!isValid.value) return

  processingClone.value = true
  processing.value = true

  try {
    const c = toRaw(chart.value)
    const cloned = {
      ...c,
      chartID: '0',
      name: `${c.name} (${t('general.label.clone').toLowerCase()})`,
      handle: '',
    }

    const created = await chartStore.create(cloned)
    chart.value = chartConstructor(created)
    capture()
    $toast.toastSuccess(t('notification.chart.created'))
    markSaved()
    router.push({ name: 'admin.charts.edit', params: { chartID: created.chartID } })
  } catch (e) {
    console.error('Failed to clone chart:', e)
    $toast.toastErrorHandler(t('notification.chart.createFailed'))(e)
  } finally {
    processing.value = false
    processingClone.value = false
  }
}

async function handleDelete() {
  processingDelete.value = true
  processing.value = true

  try {
    await chartStore.delete(toRaw(chart.value))
    capture()
    $toast.toastSuccess(t('notification.chart.deleted'))
    router.push({ name: 'admin.charts' })
  } catch (e) {
    console.error('Failed to delete chart:', e)
    $toast.toastErrorHandler(t('notification.chart.deleteFailed'))(e)
  } finally {
    processing.value = false
    processingDelete.value = false
  }
}

function exportChart() {
  if (!chart.value) return
  const blob = new Blob([JSON.stringify({ type: 'chart', list: [chart.value] }, null, 2)], {
    type: 'application/json',
  })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${chart.value.handle || chart.value.name || 'chart'}-export.json`
  a.click()
  URL.revokeObjectURL(url)
}
</script>
