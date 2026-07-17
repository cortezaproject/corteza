<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-2">
      <div class="text-xs font-semibold uppercase tracking-wide text-muted-color truncate">{{ $t(titleKey) }}</div>
    </div>

    <!-- Empty state when there is nothing to plot -->
    <div
      v-if="!total"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ height: height + 'px' }"
    >
      —
    </div>

    <div v-else>
      <v-chart :option="option" autoresize :style="{ height: height + 'px', width: '100%' }" />
      <ChartLegend :items="legendItems" :variant="variant" />
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
import { computed } from 'vue'
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

const total = computed(() => props.data.reduce((s, d) => s + (d.value || 0), 0))

// Legend rows mirror the slice data; `color` only matters for the dot
// fallback (type/unknown variants) — resolve it the same way the slices do.
const legendItems = computed(() =>
  props.data.map(d => ({
    label: d.label,
    key: d.key,
    color: d.color || colorFor(props.variant, d.label, props.category),
  })),
)

// The 2px gap between slices is drawn as a border in the chart's own surface
// colour (not a stroke around the data) so it reads as separation, not ink —
// see the dataviz skill's marks-and-anatomy "surface gap". Resolved from the
// PrimeVue content-background token (adapts light/dark automatically) since
// echarts' canvas renderer needs a literal colour, not a live CSS variable.
function surfaceGapColor() {
  if (typeof document === 'undefined') return '#ffffff'
  const val = getComputedStyle(document.documentElement).getPropertyValue('--p-content-background').trim()
  return val || '#ffffff'
}

const option = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  series: [
    {
      type: 'pie',
      radius: ['58%', '74%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: { borderColor: surfaceGapColor(), borderWidth: 2 },
      // The donut hole shows the grand total; slices identify via tooltip/legend.
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
      data: props.data.map(d => ({
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
