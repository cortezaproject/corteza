<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ isEdit ? chart?.name || $t('chart.edit.title') : $t('chart.edit.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="chart" class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 overflow-auto min-w-0">
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
        <div class="lg:col-span-7">
          <Card>
            <template #content>
              <div class="p-3">
                <!-- General settings -->
                <h5 class="mb-3">
                  {{ $t('chart.generalSettings') }}
                </h5>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">{{ $t('chart.name') }}</label>
                    <InputText
                      v-model="chart.name"
                      :placeholder="$t('chart.general.placeholder.name')"
                      :invalid="!chart.name"
                      class="w-full"
                    />
                  </div>

                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">{{ $t('chart.handle') }}</label>
                    <InputText
                      v-model="chart.handle"
                      :placeholder="$t('chart.general.placeholder.handle')"
                      class="w-full"
                    />
                    <small v-if="handleInvalid" class="text-red-500">
                      {{ $t('chart.general.placeholder.invalid-handle-characters') }}
                    </small>
                  </div>

                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('chart.colorScheme.label') }}
                    </label>
                    <Select
                      v-model="chart.config.colorScheme"
                      :options="colorSchemes"
                      option-label="name"
                      option-value="id"
                      :placeholder="$t('chart.colorScheme.placeholder')"
                      class="w-full"
                      filter
                      show-clear
                    >
                      <template #value="{ value }">
                        <div v-if="value" class="flex gap-0.5 items-center">
                          <div
                            v-for="(color, ci) in getSchemeColors(value)"
                            :key="ci"
                            :style="`background: ${color};`"
                            class="w-3.5 h-3.5 rounded-sm"
                          />
                        </div>
                        <span v-else>{{ $t('chart.colorScheme.placeholder') }}</span>
                      </template>
                      <template #option="{ option }">
                        <div class="flex gap-0.5 items-center">
                          <div
                            v-for="(color, ci) in option.colors"
                            :key="ci"
                            :style="`background: ${color};`"
                            class="w-3.5 h-3.5 rounded-sm"
                          />
                        </div>
                      </template>
                    </Select>
                  </div>

                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('chart.edit.animation.label') }}
                    </label>
                    <div class="flex items-center gap-2">
                      <ToggleSwitch v-model="animationEnabled" />
                      <span>{{ $t('chart.edit.animation.enabled') }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <Divider v-if="modules.length" />

              <!-- Report editor -->
              <component
                :is="reportEditor"
                v-if="chart && editReport"
                :report="editReport"
                :chart="chart"
                :modules="modules"
                :supported-metrics="1"
                @update:report="onReportUpdate"
              />

              <Divider />

              <!-- Toolbox settings -->
              <div class="px-3">
                <h5 class="mb-3">
                  {{ $t('chart.edit.toolbox.label') }}
                </h5>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('chart.edit.toolbox.saveAsImage.label') }}
                    </label>
                    <div class="flex items-center gap-2">
                      <ToggleSwitch v-model="saveAsImageEnabled" />
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </Card>
        </div>

        <!-- Right column: Preview -->
        <div class="lg:col-span-5">
          <div class="sticky top-0">
            <Card>
              <template #title>
                <div class="flex items-center justify-between">
                  <span class="text-sm font-medium">{{ $t('chart.edit.loadData') }}</span>
                  <Button
                    icon="pi pi-refresh"
                    text
                    size="small"
                    :disabled="processing || !reportsValid"
                    @click="refreshPreview"
                  />
                </div>
              </template>
              <template #content>
                <div class="relative" style="height: 400px">
                  <ChartRenderer
                    v-if="chart"
                    ref="chartRendererRef"
                    :key="previewKey"
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

    <!-- Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.back()"
        />
        <div class="flex gap-2">
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
            :disabled="disableSave"
            @click="handleClone"
          />
          <Button
            v-if="!hideSave"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="processingSave"
            :disabled="disableSave"
            @click="handleSave"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, inject, toRaw, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { compose, shared } from '@cortezaproject/corteza-js-next'
import { components, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { cloneDeep, isEqual } from 'lodash-es'
import { chartConstructor } from '../../../lib/charts'
import { useChartStore } from '../../../stores/chart'
import { useModuleStore } from '../../../stores/module'
import ChartRenderer from '../../../components/Chart/ChartRenderer.vue'
import * as Reports from '../../../components/Chart/Report/index.js'

const { CInputDelete } = components

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast')

const chartStore = useChartStore()
const moduleStore = useModuleStore()

const { colorschemes } = shared

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

// State
const chart = ref(null)
const initialChart = ref(null)

const { markSaved } = useUnsavedGuard({
  isDirty: () => !processingSave.value && !processingDelete.value && !processingClone.value && !!chart.value && !!initialChart.value && !isEqual(toRaw(chart.value), initialChart.value),
  messageKey: 'general.editor.unsavedChanges',
})
const loading = ref(false)
const processing = ref(false)
const processingSave = ref(false)
const processingClone = ref(false)
const processingDelete = ref(false)
const editReportIndex = ref(0)
const previewKey = ref(0)
const chartRendererRef = ref(null)

// Computed
const chartID = computed(() => route.params.chartID)

const isEdit = computed(() => {
  return chart.value && chart.value.chartID && chart.value.chartID !== '0'
})

const modules = computed(() => moduleStore.set || [])

const colorSchemes = computed(() => {
  const capitalize = w => `${w[0].toUpperCase()}${w.slice(1)}`
  const splicer = sc => {
    const rr = /(\D+)(\d+)$/gi.exec(sc)
    return { label: rr?.[1] || sc, count: rr?.[2] || '' }
  }

  const rr = []
  for (const g in colorschemes) {
    for (const sc in colorschemes[g]) {
      const gn = splicer(sc)
      rr.push({
        id: `${g}.${sc}`,
        name: `${capitalize(g)}: ${capitalize(gn.label)} (${gn.count} colors)`,
        colors: [...colorschemes[g][sc]],
      })
    }
  }
  return rr
})

function getSchemeColors(id) {
  const scheme = colorSchemes.value.find(s => s.id === id)
  return scheme?.colors || []
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

const reportsValid = computed(() => {
  if (!chart.value?.config?.reports) return false
  return !chart.value.config.reports.find(({ moduleID }) => !moduleID)
})

const reportEditor = computed(() => {
  if (!chart.value) return undefined

  if (chart.value instanceof compose.FunnelChart) return markRaw(Reports.FunnelChart)
  if (chart.value instanceof compose.GaugeChart) return markRaw(Reports.GaugeChart)
  if (chart.value instanceof compose.RadarChart) return markRaw(Reports.RadarChart)

  return markRaw(Reports.GenericChart)
})

const handleInvalid = computed(() => {
  if (!chart.value?.handle) return false
  return !/^[a-zA-Z][a-zA-Z0-9_.]*[a-zA-Z0-9]$/.test(chart.value.handle)
})

const disableSave = computed(() => {
  return !chart.value || !chart.value.name || handleInvalid.value
})

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
    initialChart.value = cloneDeep(toRaw(c))
    editReportIndex.value = 0
  } else {
    loading.value = true
    processing.value = true

    try {
      const raw = await chartStore.findByID({ namespaceID, chartID: cID, force: true })
      chart.value = chartConstructor(raw)
      initialChart.value = cloneDeep(toRaw(chart.value))
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
  return $ComposeAPI.recordReport({ namespaceID, ...r })
}

function onReportUpdate(report) {
  if (chart.value?.config?.reports) {
    chart.value.config.reports.splice(editReportIndex.value, 1, report)
  }
}

function refreshPreview() {
  chart.value.config.noAnimation = true
  previewKey.value++
}

function onUpdated() {
  processing.value = false
}

async function handleSave() {
  if (disableSave.value) return

  processingSave.value = true
  processing.value = true

  const c = toRaw(chart.value)

  try {
    if (!isEdit.value) {
      const created = await chartStore.create(c)
      chart.value = chartConstructor(created)
      initialChart.value = cloneDeep(toRaw(chart.value))
      $toast.toastSuccess(t('notification.chart.created'))
      markSaved()
      router.push({ name: 'admin.charts.edit', params: { chartID: created.chartID } })
    } else {
      const updated = await chartStore.update(c)
      chart.value = chartConstructor(updated)
      initialChart.value = cloneDeep(toRaw(chart.value))
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
  if (disableSave.value) return

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
    initialChart.value = cloneDeep(toRaw(chart.value))
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
    initialChart.value = cloneDeep(toRaw(chart.value))
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
  const blob = new Blob(
    [JSON.stringify({ type: 'chart', list: [chart.value] }, null, 2)],
    { type: 'application/json' },
  )
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${chart.value.handle || chart.value.name || 'chart'}-export.json`
  a.click()
  URL.revokeObjectURL(url)
}
</script>
