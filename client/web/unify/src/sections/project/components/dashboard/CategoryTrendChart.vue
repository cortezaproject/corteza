<template>
  <div :class="[bare ? '' : 'rounded-lg border border-surface bg-surface p-4', fill ? 'h-full' : '']">
    <!-- Card header — title left, caller-supplied controls (e.g. the
         TimeRangeSelect that windows this chart) top right, inside the card
         where the thing they act on lives. -->
    <div
      v-if="!bare && (titleKey || $slots.actions)"
      class="flex items-center justify-between gap-2 mb-3"
    >
      <div class="text-sm font-medium text-color truncate">{{ titleKey ? $t(titleKey) : '' }}</div>
      <slot name="actions" />
    </div>

    <!-- Empty state when every bucket is zero. -->
    <div
      v-if="empty"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="sizeStyle"
    >
      —
    </div>

    <v-chart v-else :option="option" autoresize :style="sizeStyle" />
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
  // Accepted for caller compatibility; the title accent dot was removed
  // (category color-coding on chart titles carried no information).
  accent: { type: String, default: '' },
  // Chrome-less variant: no card wrapper/title, for embedding in a band.
  bare: { type: Boolean, default: false },
  // Pixel height, or 'fill' to stretch to the parent's height (the parent
  // must then have a concrete height, e.g. a flex-1 min-h-* wrapper —
  // echarts' autoresize follows the container).
  height: { type: [Number, String], default: 200 },
})

const fill = computed(() => props.height === 'fill')
const sizeStyle = computed(() =>
  fill.value ? { height: '100%', width: '100%' } : { height: props.height + 'px', width: '100%' },
)

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
  return val || '#ffffff'
}

// Axis/grid chrome uses the muted-text token so it follows the active theme
// instead of a hardcoded neutral. Echarts' canvas renderer needs a literal
// colour (not a live CSS variable), so the token is resolved via
// getComputedStyle — same idiom as surfaceGapColor above — then blended to
// the desired opacity with a literal-argument color-mix() string (safe for
// canvas, unlike var(), since it needs no cascade to resolve).
function mutedAxisColor(alphaPercent) {
  const resolved =
    typeof document === 'undefined'
      ? MUTED
      : getComputedStyle(document.documentElement).getPropertyValue('--p-text-muted-color').trim() || MUTED
  return `color-mix(in srgb, ${resolved} ${alphaPercent}%, transparent)`
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
    axisLine: { lineStyle: { color: mutedAxisColor(25) } },
  },
  yAxis: {
    type: 'value',
    minInterval: 1,
    axisLabel: { color: MUTED },
    splitLine: { lineStyle: { color: mutedAxisColor(15) } },
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
