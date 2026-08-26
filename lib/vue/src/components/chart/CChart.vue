<template>
  <v-chart
    ref="chartRef"
    :option="chart"
    :update-options="updateOptions"
    :theme="theme"
    autoresize
    class="w-full h-full overflow-hidden"
  />
</template>

<script setup>
import { computed, ref, onBeforeUnmount } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  BarChart,
  LineChart,
  PieChart,
  ScatterChart,
  RadarChart,
  FunnelChart,
  GaugeChart,
} from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DatasetComponent,
  ToolboxComponent,
  DataZoomComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  BarChart,
  LineChart,
  PieChart,
  ScatterChart,
  RadarChart,
  FunnelChart,
  GaugeChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DatasetComponent,
  ToolboxComponent,
  DataZoomComponent,
])

const props = defineProps({
  chart: {
    type: Object,
    required: true,
  },
})

const chartRef = ref(null)

// Left to itself, vue-echarts diffs the option and puts every top-level array
// that lost an entry into replaceMerge. `color` is one — the palette, not a
// list of components — and echarts refuses it, then reports the update after it
// as arriving mid-process. Naming the component arrays here keeps the replace a
// rebuilt option needs and keeps the palette out of it.
const updateOptions = { replaceMerge: ['series', 'xAxis', 'yAxis'] }

const theme = computed(() => {
  const { darkMode } = props.chart || {}
  return darkMode ? 'dark' : 'light'
})

onBeforeUnmount(() => {
  if (chartRef.value) {
    chartRef.value.dispose()
  }
})
</script>
