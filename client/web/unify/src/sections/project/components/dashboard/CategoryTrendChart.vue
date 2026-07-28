<template>
  <div
    :class="[bare ? '' : 'rounded-lg border border-surface bg-surface p-4', fill ? 'h-full' : '']"
  >
    <!-- Card header — title left, caller-supplied controls (e.g. the
         TimeRangeSelect that windows this chart) top right, inside the card
         where the thing they act on lives. -->
    <div
      v-if="!bare && (titleKey || $slots.actions)"
      class="flex items-center justify-between gap-2 mb-3"
    >
      <div class="text-xs font-semibold uppercase tracking-wide text-muted-color truncate">
        {{ titleKey ? $t(titleKey) : '' }}
      </div>
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

    <template v-else>
      <v-chart :option="option" autoresize :style="sizeStyle" />
      <!-- height="fill" callers (EventsActivityPanel) are all single-series
           today, so showLegend is always false there in practice and this
           never has to compete with the chart for the parent's fixed height
           — no extra layout machinery needed for the fill case. -->
      <ChartLegend
        v-if="showLegend && !empty"
        :items="legendItems"
        :variant="legendVariant"
        :hidden="hiddenKeys"
        @toggle="onLegendToggle"
      />
    </template>
  </div>
</template>

<script setup>
import ChartLegend from '@/sections/project/components/dashboard/ChartLegend.vue'
import { MUTED } from '@/sections/project/config/chartColors'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent])

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // X-axis bucket labels, e.g. ['Jul 1', 'Jul 8', …].
  labels: { type: Array, default: () => [] },
  // Tooltip-facing bucket range strings aligned to `labels` (e.g.
  // "Jul 13 – 19, 2026" for a week bucket) — see adaptiveWindow's
  // rangeLabels. Optional; tooltips fall back to the short axis label.
  rangeLabels: { type: Array, default: () => [] },
  // One or more series: [{ name, color, data: number[], key? }] aligned to
  // `labels`. `key` is the category key, when the series are per-category.
  series: { type: Array, default: () => [] },
  // Badge family for the HTML legend below the chart ('' → colored-dot rows
  // from each series' own colour). See ChartLegend.
  legendVariant: { type: String, default: '' },
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

// Show a legend only when there is more than one series to disambiguate —
// based on the FULL series list, not the legend-filtered one below, so
// toggling series off never makes the legend itself disappear.
const showLegend = computed(() => props.series.length > 1)

// Legend rows mirror the full series list; every entry stays clickable to
// toggle back on (see ChartLegend).
const legendItems = computed(() =>
  props.series.map(s => ({ label: s.name, key: s.key, color: s.color })),
)

// Same identity ChartLegend's own keyFor uses for a series — `key` when the
// caller set one (per-category series), else the series name itself.
const itemKey = s => s.key ?? s.name

// Legend-toggled-off series keys — presentation-only (see ChartLegend).
const hiddenKeys = ref([])
function onLegendToggle(key) {
  hiddenKeys.value = hiddenKeys.value.includes(key)
    ? hiddenKeys.value.filter(k => k !== key)
    : [...hiddenKeys.value, key]
}

// What actually reaches the stacked bars — legend-hidden series dropped
// entirely, same as unselecting a series in echarts' own canvas legend (the
// stack re-totals over what's left, same as the tooltip's per-bucket sum).
const visibleSeries = computed(() =>
  props.series.filter(s => !hiddenKeys.value.includes(itemKey(s))),
)

// The 2px gap between stacked segments is drawn as a border in the chart's
// own surface colour (see the dataviz skill's marks-and-anatomy "surface
// gap" — separation via gap, never a stroke). Same resolution idiom as
// CategoryDonutChart/ResourceGraph: echarts' canvas renderer needs a literal
// colour, so the PrimeVue content-background token is read via
// getComputedStyle rather than passed as a live CSS variable.
function surfaceGapColor() {
  if (typeof document === 'undefined') return '#ffffff'
  const val = getComputedStyle(document.documentElement)
    .getPropertyValue('--p-content-background')
    .trim()
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
      : getComputedStyle(document.documentElement)
          .getPropertyValue('--p-text-muted-color')
          .trim() || MUTED
  return `color-mix(in srgb, ${resolved} ${alphaPercent}%, transparent)`
}

const option = computed(() => ({
  grid: { left: 8, right: 12, top: 8, bottom: 8, containLabel: true },
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' },
    // Charts sit inside overflow-hidden cards (see OverviewPanel's category
    // cards) and inside scrollable panels — confine (the other echarts
    // tooltip-clipping knob) does the OPPOSITE of what's wanted here (it
    // clamps the tooltip TO the container). appendTo:'body' instead mounts
    // the tooltip's DOM node straight onto <body>, so it escapes every
    // ancestor's overflow/scroll clipping entirely. Confirmed present on the
    // installed echarts@6.0.0 (node_modules/echarts/types/dist/shared.d.ts
    // TooltipOption#appendTo; appendToBody is the same idea but deprecated in
    // this version in its favour).
    appendTo: 'body',
    // Header is the bucket's full range (e.g. "Jul 13 – 19, 2026") rather
    // than the short axis label, so the tooltip disambiguates which days a
    // bar actually sums — everything else replicates echarts' default
    // marker + name + value row per series. `params` already only carries
    // whatever's in `series` below, so a legend-hidden series never shows a
    // row here either — same as echarts' own canvas legend.
    formatter: params => {
      const list = Array.isArray(params) ? params : [params]
      const idx = list[0]?.dataIndex ?? 0
      const head = props.rangeLabels[idx] || list[0]?.axisValue || ''
      const rows = list.map(p => `${p.marker}${p.seriesName}: <b>${p.value}</b>`).join('<br/>')
      return `${head}<br/>${rows}`
    },
  },
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
  series: visibleSeries.value.map(s => ({
    name: s.name,
    type: 'bar',
    stack: 'total',
    data: s.data,
    itemStyle: {
      color: s.color,
      // Rounded top corner once only one series is actually stacked, whether
      // that's because the caller only ever passed one, or because legend
      // toggling has whittled it down to one — same shape either way.
      borderRadius: visibleSeries.value.length > 1 ? 0 : [3, 3, 0, 0],
      borderColor: surfaceGapColor(),
      borderWidth: 2,
    },
    barMaxWidth: 28,
  })),
}))
</script>
