<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface-0 dark:bg-surface-900 p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-2">
      <span v-if="accent" class="w-2 h-2 rounded-full shrink-0" :style="{ background: accent }" />
      <div class="text-sm font-medium text-color truncate">{{ $t(titleKey) }}</div>
    </div>

    <!-- Empty state when there is nothing to plot -->
    <div
      v-if="!total"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ height: height + 'px' }"
    >
      —
    </div>

    <v-chart v-else :option="option" autoresize :style="{ height: height + 'px', width: '100%' }" />
  </div>
</template>

<script setup>
import { colorFor, MUTED } from '@/sections/project/config/chartColors'
import { PieChart } from 'echarts/charts'
import { LegendComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // [{ label: string, value: number }]
  data: { type: Array, default: () => [] },
  // Colour family for the slices: status | type (severity/risk use
  // CategoryRankBar instead — see CategoryView).
  variant: { type: String, default: 'status' },
  // Which category's option list to index into for variant 'type' (see
  // colorFor in config/chartColors) — ignored for every other variant.
  category: { type: String, default: '' },
  // Optional accent dot next to the title (hex).
  accent: { type: String, default: '' },
  // Chrome-less variant: no card wrapper/title, for embedding in a card.
  bare: { type: Boolean, default: false },
  height: { type: Number, default: 200 },
})

const total = computed(() => props.data.reduce((s, d) => s + (d.value || 0), 0))

// The 2px gap between slices is drawn as a border in the chart's own surface
// colour (not a stroke around the data) so it reads as separation, not ink —
// see the dataviz skill's marks-and-anatomy "surface gap". Resolved from the
// PrimeVue content-background token (adapts light/dark automatically; same
// dark-detection idiom as ResourceGraph.vue) since echarts' canvas renderer
// needs a literal colour, not a live CSS variable.
function surfaceGapColor() {
  if (typeof document === 'undefined') return '#ffffff'
  const val = getComputedStyle(document.documentElement).getPropertyValue('--p-content-background').trim()
  if (val) return val
  return document.documentElement.classList.contains('dark') ? '#18181b' : '#ffffff'
}

const option = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: {
    bottom: 0,
    left: 'center',
    icon: 'circle',
    itemHeight: 8,
    itemWidth: 8,
    itemGap: 12,
    textStyle: { color: MUTED, fontSize: 11 },
  },
  series: [
    {
      type: 'pie',
      radius: ['58%', '74%'],
      center: ['50%', '46%'],
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
        itemStyle: { color: colorFor(props.variant, d.label, props.category) },
      })),
    },
  ],
}))
</script>
