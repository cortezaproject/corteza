<template>
  <div class="rounded-lg border border-surface bg-surface-0 dark:bg-surface-900 p-4">
    <div class="text-sm font-medium text-muted-color mb-3">{{ $t(titleKey) }}</div>

    <!-- Empty state when there is nothing to plot -->
    <div
      v-if="!data.length"
      class="flex items-center justify-center text-muted-color text-sm"
      style="height: 120px"
    >
      —
    </div>

    <v-chart
      v-else
      :option="option"
      autoresize
      :style="{ height: chartHeight + 'px', width: '100%' }"
    />
  </div>
</template>

<script setup>
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent])

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, required: true },
  // [{ label: string, value: number }]
  data: { type: Array, default: () => [] },
})

// Grow the chart with the number of bars, but keep a sane minimum.
const chartHeight = computed(() => Math.max(120, props.data.length * 34))

const option = computed(() => ({
  grid: { left: 8, right: 16, top: 8, bottom: 8, containLabel: true },
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  xAxis: {
    type: 'value',
    minInterval: 1,
    axisLabel: { color: '#94a3b8' },
    splitLine: { lineStyle: { color: 'rgba(148,163,184,0.15)' } },
  },
  yAxis: {
    type: 'category',
    data: props.data.map(d => d.label),
    axisLabel: { color: '#94a3b8' },
    axisTick: { show: false },
  },
  series: [
    {
      type: 'bar',
      data: props.data.map(d => d.value),
      itemStyle: { color: '#6366f1', borderRadius: [0, 4, 4, 0] },
      barMaxWidth: 22,
    },
  ],
}))
</script>
