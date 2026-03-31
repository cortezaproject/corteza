<template>
  <PageBlock :block="block">
    <div class="relative h-full">
      <ChartRenderer v-if="chart" ref="chartRenderer" :chart="chart" :reporter="reporter" :record="record" />
      <div v-else-if="!block.options?.chartID" class="p-3 text-muted-color italic">
        {{ $t('block.chart.noChart') }}
      </div>
      <div v-else class="flex items-center justify-center h-full">
        <ProgressSpinner />
      </div>

      <!-- Live filter button -->
      <Button
        v-if="chart && block.options?.liveFilterEnabled"
        icon="pi pi-filter"
        text
        rounded
        size="small"
        :severity="hasLiveFilter ? 'primary' : 'secondary'"
        class="absolute top-1 right-1 z-10"
        @click="showFilterModal = true"
      />
    </div>

    <!-- Live filter modal -->
    <Dialog
      v-model:visible="showFilterModal"
      :header="$t('chart.filter.modal.title')"
      modal
      :style="{ width: '600px' }"
    >
      <div class="flex flex-col gap-4">
        <div v-if="originalFilter" class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.filter.modal.originalFilter.label') }}
          </label>
          <Textarea :model-value="originalFilter" readonly rows="2" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.filter.modal.liveFilter.label') }}
          </label>
          <Select
            v-model="liveFilterModalValue"
            :options="predefinedFilterOptions"
            option-label="text"
            option-value="value"
            :placeholder="$t('chart.filter.modal.liveFilter.placeholder')"
            class="w-full"
            show-clear
          />
        </div>

        <div v-if="originalFilter && liveFilterModalValue" class="flex flex-col gap-4">
          <Divider />

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.filter.modal.filterPreview.label') }}
            </label>
            <Textarea :model-value="liveFilterPreview" readonly rows="2" class="w-full" />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.filter.modal.options.label') }}
            </label>
            <div class="flex flex-col gap-2">
              <div v-for="opt in filterCombineOptions" :key="opt.value" class="flex items-center gap-2">
                <RadioButton v-model="liveFilterModalOption" :value="opt.value" :input-id="`filter-opt-${opt.value}`" />
                <label :for="`filter-opt-${opt.value}`">{{ $t(`chart.filter.modal.options.${opt.label}`) }}</label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-between w-full">
          <Button
            :label="$t('general.label.reset')"
            severity="secondary"
            text
            @click="resetLiveFilter"
          />
          <div class="flex gap-2">
            <Button
              :label="$t('general.label.cancel')"
              severity="secondary"
              @click="showFilterModal = false"
            />
            <Button
              :label="$t('general.label.save')"
              @click="applyLiveFilter"
            />
          </div>
        </div>
      </template>
    </Dialog>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, inject, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@cortezaproject/corteza-js-next'
import PageBlock from './PageBlock.vue'
import ChartRenderer from '../../Chart/ChartRenderer.vue'
import { useChartStore } from '../../../stores/chart'
import { evaluatePrefilter } from '../../../lib/record-filter'

const { t } = useI18n()

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  page: {
    type: Object,
    default: () => ({}),
  },
  record: {
    type: Object,
    default: undefined,
  },
})

const $ComposeAPI = inject('$ComposeAPI')
const $auth = inject('$auth', {})
const $eventBus = inject('$eventBus', null)
const chartStore = useChartStore()

const chart = ref(null)
const chartRenderer = ref(null)

// Live filter state
const showFilterModal = ref(false)
const originalFilter = ref(undefined)
const liveFilterValue = ref(undefined)
const liveFilterOption = ref('AND')
const liveFilterModalValue = ref(undefined)
const liveFilterModalOption = ref('AND')

const hasLiveFilter = computed(() => !!liveFilterValue.value)

const predefinedFilterOptions = computed(() =>
  compose.chartUtil.predefinedFilters.map(pf => ({
    ...pf,
    text: t(`chart.edit.filter.${pf.text}`),
  })),
)

const filterCombineOptions = [
  { label: 'and', value: 'AND' },
  { label: 'or', value: 'OR' },
  { label: 'overwrite', value: '' },
]

const liveFilterPreview = computed(() =>
  getFilter(liveFilterModalValue.value, liveFilterModalOption.value),
)

function getFilter(liveFilter = liveFilterValue.value, option = liveFilterOption.value) {
  if (liveFilter) {
    return originalFilter.value && option
      ? `(${originalFilter.value}) ${option} (${liveFilter})`
      : liveFilter
  }
  return originalFilter.value
}

function resetLiveFilter() {
  liveFilterModalValue.value = undefined
  liveFilterModalOption.value = 'AND'
}

function applyLiveFilter() {
  liveFilterValue.value = liveFilterModalValue.value
  liveFilterOption.value = liveFilterModalOption.value
  showFilterModal.value = false
  updateChart()
}

function updateChart() {
  if (chart.value) {
    chart.value = { ...chart.value, config: { ...chart.value.config, noAnimation: true } }
  }
}

async function fetchChart() {
  const { chartID } = props.block.options || {}

  if (!chartID) {
    return
  }

  const { namespaceID } = props.namespace

  try {
    chart.value = await chartStore.findByID({ chartID, namespaceID })
  } catch (e) {
    console.error('Failed to load chart:', e)
  }
}

function reporter(r = {}) {
  const { namespaceID } = props.namespace
  let { filter } = r

  // Capture original filter for the live filter modal
  if (originalFilter.value === undefined && filter) {
    originalFilter.value = filter
  }

  // Apply live filter combination
  filter = getFilter()

  // Evaluate prefilter with record context if available
  if (filter) {
    const record = props.record
    const user = $auth?.user || {}

    // If filter uses ${record} or ${ownerID} and there is no record, return empty
    if (!record && (filter.includes('${record') || filter.includes('${ownerID}'))) {
      return Promise.resolve([])
    }

    filter = evaluatePrefilter(filter, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  }

  return $ComposeAPI.recordReport({ namespaceID, ...r, filter })
}

onMounted(() => {
  fetchChart()
})

watch(
  () => props.block.options?.chartID,
  () => {
    chart.value = null
    originalFilter.value = undefined
    liveFilterValue.value = undefined
    fetchChart()
  },
)

const offRefetch = $eventBus?.on('refetch-records', () => fetchChart())

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>
