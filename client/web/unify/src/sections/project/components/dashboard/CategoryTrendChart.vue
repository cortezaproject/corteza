<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-3">
      <span v-if="accent" class="w-2 h-2 rounded-full shrink-0" :style="{ background: accent }" />
      <div class="text-sm font-medium text-color truncate">{{ $t(titleKey) }}</div>
    </div>

    <!-- Empty state when every bucket is zero. -->
    <div
      v-if="empty"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ height: height + 'px' }"
    >
      —
    </div>

    <v-chart v-else :option="option" autoresize :style="{ height: height + 'px', width: '100%' }" />
  </div>
</template>

<script setup>
import { MUTED } from '@/sections/project/config/chartColors'
import { BarChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent, LegendComponent])

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // X-axis bucket labels, e.g. ['Jul 1', 'Jul 8', …].
  labels: { type: Array, default: () => [] },
  // One or more series: [{ name, color, data: number[] }] aligned to `labels`.
  series: { type: Array, default: () => [] },
  // Optional accent dot next to the title (hex).
  accent: { type: String, default: '' },
  // Chrome-less variant: no card wrapper/title, for embedding in a band.
  bare: { type: Boolean, default: false },
  height: { type: Number, default: 200 },
})

const empty = computed(() => !props.series.some(s => (s.data || []).some(v => v > 0)))

// Show a legend only when there is more than one series to disambiguate.
const showLegend = computed(() => props.series.length > 1)

// The 2px gap between stacked segments is drawn as a border in the chart's
// own surface colour (see the dataviz skill's marks-and-anatomy "surface
// gap" — separation via gap, never a stroke). Same resolution idiom as
// CategoryDonutChart/ResourceGraph: echarts' canvas renderer needs a literal
// colour, so the PrimeVue content-background token is read via
// getComputedStyle rather than passed as a live CSS variable.
function surfaceGapColor() {
  if (typeof document === 'undefined') return '#ffffff'
  const val = getComputedStyle(document.documentElement).getPropertyValue('--p-content-background').trim()
  if (val) return val
  return document.documentElement.classList.contains('dark') ? '#18181b' : '#ffffff'
}

const option = computed(() => ({
  grid: { left: 8, right: 12, top: 8, bottom: showLegend.value ? 28 : 8, containLabel: true },
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: showLegend.value
    ? { bottom: 0, left: 'center', icon: 'circle', itemHeight: 8, itemWidth: 8, textStyle: { color: MUTED, fontSize: 11 } }
    : undefined,
  xAxis: {
    type: 'category',
    data: props.labels,
    axisLabel: { color: MUTED, fontSize: 11 },
    axisTick: { show: false },
    axisLine: { lineStyle: { color: 'rgba(148,163,184,0.25)' } },
  },
  yAxis: {
    type: 'value',
    minInterval: 1,
    axisLabel: { color: MUTED },
    splitLine: { lineStyle: { color: 'rgba(148,163,184,0.15)' } },
  },
  series: props.series.map(s => ({
    name: s.name,
    type: 'bar',
    stack: 'total',
    data: s.data,
    itemStyle: {
      color: s.color,
      borderRadius: props.series.length > 1 ? 0 : [3, 3, 0, 0],
      borderColor: surfaceGapColor(),
      borderWidth: 2,
    },
    barMaxWidth: 28,
  })),
}))
</script>
