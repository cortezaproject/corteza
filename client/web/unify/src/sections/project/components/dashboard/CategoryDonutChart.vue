<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-2">
      <div class="text-xs font-semibold uppercase tracking-wide text-muted-color truncate">
        {{ $t(titleKey) }}
      </div>
    </div>

    <!-- Empty state when there is nothing to plot -->
    <div
      v-if="!fullTotal"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ height: height + 'px' }"
    >
      —
    </div>

    <div v-else>
      <v-chart :option="option" autoresize :style="{ height: height + 'px', width: '100%' }" />
      <ChartLegend
        :items="legendItems"
        :variant="variant"
        :hidden="hiddenKeys"
        @toggle="onLegendToggle"
      />
    </div>
  </div>
</template>

<script setup>
import ChartLegend from '@/sections/project/components/dashboard/ChartLegend.vue'
import { colorFor, MUTED } from '@/sections/project/config/chartColors'
import { PieChart } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, PieChart, TooltipComponent])

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // [{ label: string, value: number, color?: string, key?: string }]. `key`
  // is the raw category key, used by the 'category' variant's HTML legend
  // (see ChartLegend) to look up its KindIcon config.
  data: { type: Array, default: () => [] },
  // Colour family for the slices: status | type (severity/risk use
  // CategoryRankBar instead — see CategoryView).
  variant: { type: String, default: 'status' },
  // Which category's option list to index into for variant 'type' (see
  // colorFor in config/chartColors) — ignored for every other variant.
  category: { type: String, default: '' },
  // Accepted for caller compatibility; the title accent dot was removed
  // (category color-coding on chart titles carried no information).
  accent: { type: String, default: '' },
  // Chrome-less variant: no card wrapper/title, for embedding in a card.
  bare: { type: Boolean, default: false },
  height: { type: Number, default: 200 },
})

// The donut hole's number, summed over the VISIBLE slices (ruled 2026-07-28):
// filtering the legend re-sums the hole, so hiding "Completed" reads the open
// total straight off the ring. The trade, accepted deliberately: while a filter
// is active this number no longer matches the KPI tiles or board totals beside
// it — it answers "what's showing" rather than "what the category holds".
// `fullTotal` still gates the empty state, so a chart with data never renders
// as empty just because every slice was toggled off.
const fullTotal = computed(() => props.data.reduce((s, d) => s + (d.value || 0), 0))

// Same identity ChartLegend's own keyFor uses for a row — `key` when the data
// carries one (the 'category' variant, or a row with a pinned colour), else
// the label itself.
const itemKey = d => d.key ?? d.label

// Legend-toggled-off slice keys — presentation-only (see ChartLegend); reset
// implicitly whenever a stale key no longer matches anything (a project/
// category switch), since filtering below is a plain `includes` check.
const hiddenKeys = ref([])
function onLegendToggle(key) {
  hiddenKeys.value = hiddenKeys.value.includes(key)
    ? hiddenKeys.value.filter(k => k !== key)
    : [...hiddenKeys.value, key]
}

// Legend rows mirror the FULL slice data (not the filtered-for-the-ring set
// below) — every entry stays clickable to toggle back on. `color` only
// matters for the dot fallback (type/unknown variants) — resolve it the same
// way the slices do.
const legendItems = computed(() =>
  props.data.map(d => ({
    label: d.label,
    key: d.key,
    color: d.color || colorFor(props.variant, d.label, props.category),
  })),
)

// What actually reaches the ring — legend-hidden slices dropped entirely
// (rather than zeroed) so echarts' own d% recomputes over what's left, same
// as toggling off a series in its canvas legend.
const visibleData = computed(() => props.data.filter(d => !hiddenKeys.value.includes(itemKey(d))))
// Declared after visibleData on purpose — it reads it.
const total = computed(() => visibleData.value.reduce((s, d) => s + (d.value || 0), 0))

// The 2px gap between slices is drawn as a border in the chart's own surface
// colour (not a stroke around the data) so it reads as separation, not ink —
// see the dataviz skill's marks-and-anatomy "surface gap". Resolved from the
// PrimeVue content-background token (adapts light/dark automatically) since
// echarts' canvas renderer needs a literal colour, not a live CSS variable.
function surfaceGapColor() {
  if (typeof document === 'undefined') return '#ffffff'
  const val = getComputedStyle(document.documentElement)
    .getPropertyValue('--p-content-background')
    .trim()
  return val || '#ffffff'
}

const option = computed(() => ({
  tooltip: {
    trigger: 'item',
    // A FUNCTION formatter, not a string template — echarts' string-template
    // substitution (`formatTpl` in echarts' own dist/echarts.esm.js) only
    // ever rewrites the fixed {a}/{b}/{c}/{d} aliases off each data param's
    // own `$vars` list; `marker` is never one of them there. `params.marker`
    // (a ready-made <span> dot in that slice's own itemStyle colour) is only
    // populated for the formatter to CONSUME as a function — confirmed via
    // TooltipView.prototype._showSeriesItemTooltip in that same file, which
    // sets `params.marker` immediately before invoking the formatter, with
    // its own comment: "Users can assemble richText text in `formatter`
    // callback and use those markers style." Building the string ourselves
    // like this (rather than hand-rolling a swatch) means the marker and the
    // slice can never disagree on colour.
    formatter: params => `${params.marker}${params.name}: ${params.value} (${params.percent}%)`,
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
  },
  series: [
    {
      type: 'pie',
      radius: ['58%', '74%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: { borderColor: surfaceGapColor(), borderWidth: 2 },
      // The donut hole shows the grand total (see `total` above — always the
      // FULL data, never the legend-filtered set); slices identify via
      // tooltip/legend.
      label: {
        show: true,
        position: 'center',
        formatter: () => String(total.value),
        color: MUTED,
        fontSize: 24,
        fontWeight: 600,
      },
      emphasis: { scale: true, scaleSize: 4, label: { show: true } },
      labelLine: { show: false },
      data: visibleData.value.map(d => ({
        name: d.label,
        value: d.value,
        // A row may pin its own colour (e.g. the backlog's by-category donut,
        // whose labels are translated titles that colorFor can't key on).
        itemStyle: { color: d.color || colorFor(props.variant, d.label, props.category) },
      })),
    },
  ],
}))
</script>
