<template>
  <div class="flex flex-col gap-3">
    <!-- Chart selection -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.chart.display') }}</label>
      <div class="flex gap-2">
        <Select
          v-model="chartID"
          :options="charts"
          option-label="name"
          option-value="chartID"
          :placeholder="$t('block.chart.pick')"
          class="flex-1"
          filter
        />
        <Button
          v-if="selectedChart"
          icon="pi pi-external-link"
          severity="secondary"
          :title="$t('block.chart.openInBuilder')"
          @click="goToChart"
        />
      </div>
    </div>

    <template v-if="selectedChart">
      <!-- Live filter -->
      <CInputSwitch v-model="liveFilterEnabled" :label="$t('block.chart.enableLiveFilter')" />

      <Divider />

      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.chart.drillDown.label') }}</h5>

        <CInputSwitch v-model="drillDownEnabled" :label="$t('block.chart.drillDown.enabled')" />

        <template v-if="drillDownEnabled">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.chart.drillDown.blockID') }}</label>
            <InputText
              v-model="drillDownBlockID"
              :placeholder="$t('block.chart.drillDown.blockIDPlaceholder')"
              class="w-full"
            />
            <small class="text-muted-color">{{ $t('block.chart.drillDown.blockIDFootnote') }}</small>
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useChartStore } from '@/stores/chart'

const router = useRouter()
const chartStore = useChartStore()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const charts = computed(() => chartStore.set || [])

const selectedChart = computed(() => {
  if (!chartID.value) return null
  return charts.value.find(c => c.chartID === chartID.value) || null
})

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const chartID = computed({
  get: () => props.block.options?.chartID,
  set: v => updateOptions('chartID', v),
})

const liveFilterEnabled = computed({
  get: () => !!props.block.options?.liveFilterEnabled,
  set: v => updateOptions('liveFilterEnabled', v),
})

const drillDownEnabled = computed({
  get: () => !!props.block.options?.drillDown?.enabled,
  set: v => updateOptions('drillDown', { ...props.block.options?.drillDown, enabled: v }),
})

const drillDownBlockID = computed({
  get: () => props.block.options?.drillDown?.blockID || '',
  set: v => updateOptions('drillDown', { ...props.block.options?.drillDown, blockID: v }),
})

function goToChart() {
  if (!selectedChart.value) return
  router.push({
    name: 'admin.charts.edit',
    params: { chartID: selectedChart.value.chartID },
  })
}
</script>
