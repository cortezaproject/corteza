<template>
  <div class="flex flex-col gap-2 min-w-0">
    <!-- An empty range takes a line, not the chart's height -->
    <div v-if="empty" class="flex items-center text-base text-muted-color py-2">
      {{ $t('dashboard.empty') }}
    </div>
    <div v-else :style="{ height: height + 'px' }" class="min-w-0 cursor-pointer" @click="onClick">
      <v-chart
        ref="chartRef"
        :option="option"
        :update-options="updateOptions"
        :theme="isDark ? 'dark' : 'light'"
        autoresize
        class="w-full h-full overflow-hidden"
      />
    </div>

    <!-- Identity never rides on colour alone: two or more series get a legend -->
    <div
      v-if="series.length > 1 && !empty"
      class="flex flex-wrap gap-x-4 gap-y-1 text-sm text-muted-color"
    >
      <span v-for="s in series" :key="s.key" class="inline-flex items-center gap-1.5">
        <span class="inline-block h-2 w-2 rounded-full" :style="{ backgroundColor: s.color }" />
        {{ s.name }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { BarChart, LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { chromeColors, useIsDark } from './chartTheme'

use([CanvasRenderer, BarChart, LineChart, GridComponent, TooltipComponent])

const chartRef = ref(null)
const updateOptions = { replaceMerge: ['series', 'xAxis', 'yAxis'] }

const props = defineProps({
  labels: { type: Array, default: () => [] },
  rangeLabels: { type: Array, default: () => [] },
  // [{ key, name, color, data: number[], type?: 'bar' | 'line' }]
  series: { type: Array, default: () => [] },
  stacked: { type: Boolean, default: false },
  height: { type: Number, default: 220 },
})

const emit = defineEmits(['select'])

const isDark = useIsDark()

// A click anywhere in a bucket's column opens that bucket, bar or no bar.
// Bound on the wrapper rather than through the chart's own event path, and
// resolved against the chart's pixel space.
function onClick(e) {
  const chart = chartRef.value
  if (!chart || !e) return
  const rect = chart.getDom().getBoundingClientRect()
  // the grid finder answers in data space: [category index, value]
  const found = chart.convertFromPixel({ gridIndex: 0 }, [
    e.clientX - rect.left,
    e.clientY - rect.top,
  ])
  const index = Array.isArray(found) ? Math.round(found[0]) : NaN
  if (Number.isInteger(index) && index >= 0 && index < props.labels.length) emit('select', index)
}

const empty = computed(() => !props.series.some(s => (s.data || []).some(v => v > 0)))

// Only the topmost non-zero segment of a stacked column is rounded, so the
// column reads as one bar.
function topmost(bucketIndex) {
  for (let i = props.series.length - 1; i >= 0; i--) {
    if ((props.series[i].data[bucketIndex] || 0) > 0) return i
  }
  return -1
}

const option = computed(() => {
  const c = chromeColors(isDark.value)
  const n = props.labels.length
  const dense = n > 16

  const series = props.series.map((s, idx) => {
    const type = s.type || 'bar'
    if (type === 'line') {
      return {
        name: s.name,
        type: 'line',
        data: s.data,
        smooth: 0.3,
        showSymbol: false,
        symbolSize: 8,
        lineStyle: { color: s.color, width: 2 },
        itemStyle: { color: s.color, borderColor: c.surface, borderWidth: 2 },
        areaStyle: { color: s.color, opacity: 0.08 },
        emphasis: { focus: 'none' },
        z: 3,
      }
    }

    return {
      name: s.name,
      type: 'bar',
      stack: props.stacked ? 'all' : undefined,
      barMaxWidth: dense ? 14 : 28,
      barCategoryGap: '35%',
      itemStyle: { color: s.color, borderColor: c.surface, borderWidth: props.stacked ? 1 : 0 },
      emphasis: { itemStyle: { opacity: 0.85 } },
      data: s.data.map((v, i) => ({
        value: v,
        itemStyle: {
          borderRadius: !props.stacked || topmost(i) === idx ? [3, 3, 0, 0] : 0,
        },
      })),
    }
  })

  return {
    darkMode: isDark.value,
    backgroundColor: 'transparent',
    animationDuration: 300,
    grid: { left: 8, right: 8, top: 12, bottom: 4, containLabel: true },
    tooltip: {
      trigger: 'axis',
      // rendered on the body so a card's overflow never clips it
      appendTo: 'body',
      axisPointer: { type: props.series.some(s => s.type === 'line') ? 'line' : 'shadow' },
      backgroundColor: c.surface,
      borderColor: c.grid,
      textStyle: { color: c.text, fontSize: 13 },
      formatter: params => {
        const i = params[0]?.dataIndex ?? 0
        const head = props.rangeLabels[i] || props.labels[i] || ''
        const rows = params
          .map(
            p =>
              `<div style="display:flex;justify-content:space-between;gap:16px"><span>${p.marker}${p.seriesName}</span><b>${Number(p.value).toLocaleString()}</b></div>`,
          )
          .join('')
        return `<div style="margin-bottom:4px;color:${c.muted}">${head}</div>${rows}`
      },
    },
    xAxis: {
      type: 'category',
      data: props.labels,
      axisLabel: { color: c.muted, fontSize: 12, interval: 'auto', hideOverlap: true },
      axisLine: { lineStyle: { color: c.grid } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: c.muted, fontSize: 12 },
      splitLine: { lineStyle: { color: c.grid, type: 'solid' } },
      splitNumber: 3,
    },
    series,
  }
})
</script>
