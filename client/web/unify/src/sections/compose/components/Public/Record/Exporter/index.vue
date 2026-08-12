<template>
  <div>
    <Button
      :label="$t('general.label.export')"
      icon="pi pi-download"
      size="small"
      severity="secondary"
      @click="openDialog"
    />

    <Dialog
      v-model:visible="showDialog"
      :header="$t('general.label.export')"
      modal
      :style="{ width: '60rem' }"
      :breakpoints="{ '1200px': '75vw', '768px': '95vw' }"
      :content-style="{ 'max-height': '75vh', overflow: 'auto' }"
    >
      <div class="flex flex-col gap-5">
        <!-- Field selection -->
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.export.selectFields') }}
          </label>
          <CFieldPicker
            v-model="selectedFieldNames"
            :all-fields="allFields"
            list-class="max-h-64"
          />
          <small class="text-muted-color">
            {{ $t('block.recordList.export.limitations') }}
          </small>
        </div>

        <!-- Range type -->
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.export.inRange') }}
          </label>
          <div class="flex flex-wrap gap-x-6 gap-y-2">
            <div class="flex items-center gap-2">
              <RadioButton v-model="rangeType" input-id="range-all" value="all" />
              <label for="range-all" class="cursor-pointer text-sm">
                {{ $t('block.recordList.export.all') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <RadioButton v-model="rangeType" input-id="range-range" value="range" />
              <label for="range-range" class="cursor-pointer text-sm">
                {{ $t('block.recordList.export.inRange') }}
              </label>
            </div>
            <div v-if="hasSelection" class="flex items-center gap-2">
              <RadioButton v-model="rangeType" input-id="range-selection" value="selection" />
              <label for="range-selection" class="cursor-pointer text-sm">
                {{ $t('block.recordList.export.selection') }} ({{ selection.length }})
              </label>
            </div>
          </div>
        </div>

        <!-- Date range (only when rangeType = range) -->
        <div v-if="rangeType === 'range'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.export.rangeBy') }}
            </label>
            <Select
              v-model="rangeBy"
              :options="rangeByOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              @change="fetchCount"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.export.dateRange') }}
            </label>
            <Select
              v-model="rangePreset"
              :options="dateRangeOptions"
              option-label="label"
              option-value="value"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.export.filter.from') }}
            </label>
            <CInputDateTime
              v-model="startDate"
              value-type="date"
              only-date
              :max-date="endDate || undefined"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.export.filter.to') }}
            </label>
            <CInputDateTime
              v-model="endDate"
              value-type="date"
              only-date
              :min-date="startDate || undefined"
            />
          </div>
        </div>

        <!-- Free-form filter (when not exporting a selection) -->
        <div v-if="rangeType !== 'selection'" class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.export.filter.label') }}
          </label>
          <Textarea
            v-model="extraFilter"
            :placeholder="$t('block.recordList.export.filter.placeholder')"
            rows="2"
            class="w-full font-mono text-xs"
            @update:model-value="fetchCount"
          />
          <small class="text-muted-color">
            {{ $t('block.recordList.export.filter.footnote') }}
          </small>
        </div>

        <!-- Timezone -->
        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
            <Checkbox v-model="forTimezone" binary input-id="for-timezone" />
            <label for="for-timezone" class="cursor-pointer text-sm">
              {{ $t('block.recordList.export.specifyTimezone') }}
            </label>
          </div>
          <Select
            v-if="forTimezone"
            v-model="exportTimezone"
            :options="timezones"
            filter
            :placeholder="$t('block.recordList.export.timezonePlaceholder')"
            class="w-full"
          />
        </div>

        <!-- Multi-value delimiter -->
        <div v-if="hasMultiValueField" class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.export.multiValueDelimiter.label') }}
          </label>
          <Select
            v-model="multiValueDelimiter"
            :options="delimiterOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <!-- Resolve refs -->
        <div class="flex flex-col gap-1">
          <div class="flex items-center gap-2">
            <Checkbox v-model="resolveRefs" binary input-id="resolve-refs" />
            <label for="resolve-refs" class="cursor-pointer text-sm">
              {{ $t('block.recordList.export.resolveRefs') }}
            </label>
          </div>
          <small class="text-muted-color">
            {{ $t('block.recordList.export.resolveRefsNote') }}
          </small>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-between gap-2 w-full">
          <div class="flex items-center gap-2 text-sm text-muted-color">
            <ProgressSpinner
              v-if="countLoading"
              style="width: 1rem; height: 1rem"
              stroke-width="6"
            />
            <span v-else>
              {{ $t('block.recordList.export.recordCount', { count: exportableCount }) }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <Button
              :label="$t('block.recordList.export.csv')"
              severity="secondary"
              :disabled="!canExport"
              :loading="exporting === 'csv'"
              @click="doExport('csv')"
            />
            <Button
              :label="$t('block.recordList.export.json')"
              :disabled="!canExport"
              :loading="exporting === 'json'"
              @click="doExport('json')"
            />
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import moment from 'moment'
import { components } from '@planetcrust/human-vue'
import { getFieldFilter } from '../../../../lib/record-filter'

const { CInputDateTime } = components

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const props = defineProps({
  namespace: { type: Object, required: true },
  module: { type: Object, required: true },
  filter: { type: String, default: '' },
  selection: { type: Array, default: () => [] },
})

// --- dialog state ---
const showDialog = ref(false)
const exporting = ref(null)

// --- field selection ---
const systemFieldNames = ['ownedBy', 'createdAt', 'createdBy', 'updatedAt', 'updatedBy']

const allFields = computed(() => {
  const moduleFields = (props.module.fields || [])
    .filter(f => f.kind !== 'File')
    .map(f => ({ name: f.name, label: f.label || f.name, kind: f.kind, isMulti: f.isMulti }))

  const sysFields = systemFieldNames.map(name => ({ name, label: name }))

  return [...moduleFields, ...sysFields]
})

const selectedFieldNames = ref([])

watch(
  allFields,
  fields => {
    if (selectedFieldNames.value.length === 0) {
      selectedFieldNames.value = fields.map(f => f.name)
    }
  },
  { immediate: true },
)

const hasSelection = computed(() => props.selection.length > 0)

const hasMultiValueField = computed(() => {
  const names = new Set(selectedFieldNames.value)
  return (props.module.fields || []).some(f => f.isMulti && names.has(f.name))
})

// --- range / filter ---
const rangeType = ref('all')
const rangeBy = ref('createdAt')
const rangePreset = ref('lastMonth')
const startDate = ref(null)
const endDate = ref(null)
const extraFilter = ref(props.filter || '')

const rangeByOptions = computed(() => [
  { value: 'createdAt', label: t('block.recordList.export.filter.createdAt') },
  { value: 'updatedAt', label: t('block.recordList.export.filter.updatedAt') },
])

const dateRangeOptions = computed(() => [
  { value: 'lastMonth', label: t('block.recordList.export.filter.lastMonth') },
  { value: 'thisMonth', label: t('block.recordList.export.filter.thisMonth') },
  { value: 'lastWeek', label: t('block.recordList.export.filter.lastWeek') },
  { value: 'thisWeek', label: t('block.recordList.export.filter.thisWeek') },
  { value: 'today', label: t('block.recordList.export.filter.today') },
  { value: 'custom', label: t('block.recordList.export.filter.custom') },
])

function applyPreset(preset) {
  const now = moment()
  switch (preset) {
    case 'lastMonth':
      startDate.value = now.clone().subtract(1, 'months').startOf('month').toDate()
      endDate.value = now.clone().subtract(1, 'months').endOf('month').toDate()
      break
    case 'thisMonth':
      startDate.value = now.clone().startOf('month').toDate()
      endDate.value = now.clone().endOf('month').toDate()
      break
    case 'lastWeek':
      startDate.value = now.clone().subtract(1, 'weeks').startOf('week').toDate()
      endDate.value = now.clone().subtract(1, 'weeks').endOf('week').toDate()
      break
    case 'thisWeek':
      startDate.value = now.clone().startOf('week').toDate()
      endDate.value = now.clone().endOf('week').toDate()
      break
    case 'today':
      startDate.value = now.clone().startOf('day').toDate()
      endDate.value = now.clone().endOf('day').toDate()
      break
    case 'custom':
      // leave as-is
      break
  }
}

// init preset to lastMonth
applyPreset('lastMonth')

watch(rangePreset, p => {
  if (p !== 'custom') applyPreset(p)
  fetchCount()
})

let userEditingDate = false
watch([startDate, endDate], () => {
  if (!userEditingDate) {
    userEditingDate = true
    rangePreset.value = 'custom'
    // defer reset so the watcher doesn't loop
    setTimeout(() => {
      userEditingDate = false
    }, 0)
  }
  fetchCount()
})

watch(
  () => props.filter,
  f => {
    extraFilter.value = f || ''
  },
)

// --- timezone ---
const forTimezone = ref(false)
const exportTimezone = ref(null)

const timezones = computed(() => {
  try {
    if (typeof Intl !== 'undefined' && typeof Intl.supportedValuesOf === 'function') {
      return Intl.supportedValuesOf('timeZone')
    }
  } catch {
    /* ignore */
  }
  return [
    'UTC',
    'Europe/London',
    'Europe/Berlin',
    'Europe/Ljubljana',
    'America/New_York',
    'America/Los_Angeles',
    'Asia/Tokyo',
    'Asia/Shanghai',
    'Australia/Sydney',
  ]
})

// --- delimiter / refs ---
const multiValueDelimiter = ref(';')
const resolveRefs = ref(false)

const delimiterOptions = computed(() => [
  { value: ';', label: t('block.recordList.export.multiValueDelimiter.semicolon.label') },
  { value: ',', label: t('block.recordList.export.multiValueDelimiter.comma.label') },
  { value: '|', label: t('block.recordList.export.multiValueDelimiter.pipe.label') },
  { value: '[;]', label: t('block.recordList.export.multiValueDelimiter.semicolonArray.label') },
  { value: '[,]', label: t('block.recordList.export.multiValueDelimiter.commaArray.label') },
  { value: '[|]', label: t('block.recordList.export.multiValueDelimiter.pipeArray.label') },
])

// --- filter assembly ---
function fmt(d) {
  return moment(d).format('YYYY-MM-DD')
}

function makeFilter() {
  if (rangeType.value === 'selection' && props.selection.length > 0) {
    const ids = props.selection.map(r => `'${r.recordID || r}'`).join(',')
    return `recordID IN (${ids})`
  }

  const base = (extraFilter.value || '').trim()

  if (rangeType.value === 'all') {
    return base
  }

  // rangeType === 'range'
  const start = startDate.value ? fmt(startDate.value) : null
  const end = endDate.value ? fmt(endDate.value) : null

  let dateRangeQuery = ''

  // Normalize date-only boundaries to cover the whole selected days: start from the
  // beginning of the start day and up to the end of the end day. Without the
  // end-of-day adjustment the end date resolves to midnight and records logged
  // during the end day itself are excluded from the export.
  const startOfDay = d => moment(d, 'YYYY-MM-DD').utc().format()
  const endOfDay = d => moment(d, 'YYYY-MM-DD').endOf('day').utc().format()

  if (start && end) {
    dateRangeQuery =
      getFieldFilter(
        rangeBy.value,
        'DateTime',
        { start: startOfDay(start), end: endOfDay(end) },
        'BETWEEN',
      ) || ''
  } else if (start) {
    dateRangeQuery = getFieldFilter(rangeBy.value, 'DateTime', startOfDay(start), '>=') || ''
  } else if (end) {
    dateRangeQuery = getFieldFilter(rangeBy.value, 'DateTime', endOfDay(end), '<=') || ''
  }

  return base && dateRangeQuery ? `(${base}) AND ${dateRangeQuery}` : dateRangeQuery || base
}

const canExport = computed(() => selectedFieldNames.value.length > 0 && dateRangeValid.value)

const dateRangeValid = computed(() => {
  if (rangeType.value !== 'range') return true
  if (!startDate.value || !endDate.value) return true
  return moment(endDate.value).isSameOrAfter(moment(startDate.value))
})

// --- record count ---
const exportableCount = ref(0)
const countLoading = ref(false)
let countTimer = null

function fetchCount() {
  clearTimeout(countTimer)
  countTimer = setTimeout(async () => {
    countLoading.value = true
    try {
      const query = makeFilter()
      const result = await $ComposeAPI.recordList({
        namespaceID: props.namespace.namespaceID,
        moduleID: props.module.moduleID,
        query,
        limit: 1,
        incTotal: true,
      })
      exportableCount.value = result?.filter?.total ?? 0
    } catch {
      exportableCount.value = 0
    } finally {
      countLoading.value = false
    }
  }, 300)
}

watch(rangeType, fetchCount)

function openDialog() {
  showDialog.value = true
  rangeType.value = hasSelection.value ? 'selection' : 'all'
  extraFilter.value = props.filter || ''
  fetchCount()
}

async function doExport(ext) {
  const { namespaceID } = props.namespace
  const { moduleID, name } = props.module

  exporting.value = ext
  try {
    const blob = await $ComposeAPI.recordExport(
      {
        namespaceID,
        moduleID,
        filename: name || 'export',
        ext,
        fields: selectedFieldNames.value.join(','),
        filter: makeFilter() || undefined,
        multiValueDelimiter: multiValueDelimiter.value || undefined,
        timezone: forTimezone.value ? exportTimezone.value || undefined : undefined,
        resolveRefs: resolveRefs.value || undefined,
      },
      { responseType: 'blob' },
    )

    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${name || 'export'}.${ext}`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    console.error('Export failed', e)
  } finally {
    exporting.value = null
  }
}
</script>
