<script setup lang="ts">
import type { App as McpApp } from '@modelcontextprotocol/ext-apps'
import CChart from '@planetcrust/human-vue/src/components/chart/CChart.vue'
import CResourceTable from '@planetcrust/human-vue/src/components/resource-table/CResourceTable.vue'
import Button from 'primevue/button'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { hostTheme, mcpAppKey } from '../shared/mount'
import {
  type Metric,
  type Report,
  dimensionKey,
  formOf,
  formatValue,
  groupLabel,
  metricsOf,
  palette,
  parseReport,
  unitOf,
} from './recordReport'

const { t } = useI18n()
const bridge = inject<McpApp>(mcpAppKey)!

const report = ref<Report>()
const mode = ref<'chart' | 'table'>('chart')

bridge.ontoolresult = r => {
  report.value = parseReport(r)
}

const metrics = computed(() => (report.value ? metricsOf(report.value) : []))
const labels = computed(() =>
  (report.value?.rows ?? []).map(r =>
    groupLabel(r[dimensionKey], report.value!, t('mcpApp.recordReport.empty')),
  ),
)
const form = computed(() => (report.value ? formOf(report.value, labels.value) : 'bar'))
const dimensionName = computed(() => {
  const d = report.value?.view?.dimension
  return d ? d.label || d.name : t('mcpApp.recordReport.group')
})

function metricTitle(m: Metric) {
  if (m.key === 'count') return t('mcpApp.recordReport.count')
  return m.key
}

function format(value: unknown, m: Metric) {
  return formatValue(value, unitOf(m, report.value!))
}

// One chart per metric: metrics differ in scale and a chart has one axis.
function optionFor(m: Metric) {
  const values = report.value!.rows.map(r => r[m.key] as number)
  const horizontal = form.value === 'barHorizontal'
  const category = {
    type: 'category',
    data: labels.value,
    axisTick: { show: false },
    inverse: horizontal,
  }
  const value = {
    type: 'value',
    splitLine: { lineStyle: { opacity: 0.35 } },
    axisLabel: { formatter: (v: number) => format(v, m) },
  }
  const directLabels = labels.value.length <= 12

  // Dates sit on a time axis, so uneven gaps between groups read as uneven.
  const line = form.value === 'line'
  const times = report.value!.rows.map(r => r[dimensionKey] as string)

  const series = line
    ? {
        type: 'line',
        data: values.map((v, i) => [times[i], v]),
        lineStyle: { width: 2 },
        symbolSize: 8,
        showSymbol: values.length <= 30,
      }
    : {
        type: 'bar',
        data: values,
        barMaxWidth: 32,
        itemStyle: { borderRadius: horizontal ? [0, 4, 4, 0] : [4, 4, 0, 0] },
        label: {
          show: directLabels,
          position: horizontal ? 'right' : 'top',
          formatter: (p: { value: number }) => format(p.value, m),
        },
      }

  return {
    darkMode: hostTheme.value === 'dark',
    backgroundColor: 'transparent',
    color: palette,
    grid: { left: 8, right: 24, top: 16, bottom: 8, containLabel: true },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: form.value === 'line' ? 'line' : 'shadow' },
      valueFormatter: (v: number) => format(v, m),
    },
    xAxis: horizontal ? value : line ? { type: 'time', splitLine: { show: false } } : category,
    yAxis: horizontal ? category : value,
    series: [{ ...series, name: metricTitle(m) }],
  }
}

function chartHeight() {
  return form.value === 'barHorizontal'
    ? `${Math.max(160, labels.value.length * 28 + 40)}px`
    : '220px'
}

const tableFields = computed(() => [
  { key: 'group', header: dimensionName.value },
  ...metrics.value.map(m => ({ key: m.key, header: metricTitle(m) })),
])
const tableRows = computed(() =>
  (report.value?.rows ?? []).map((r, i) => ({
    _dataKey: String(i),
    group: labels.value[i],
    ...Object.fromEntries(metrics.value.map(m => [m.key, format(r[m.key], m)])),
  })),
)
</script>

<template>
  <div class="p-3 text-sm">
    <div v-if="!report" id="status" class="text-muted-color">
      {{ t('mcpApp.recordReport.waiting') }}
    </div>

    <template v-else>
      <div class="flex items-center justify-between gap-3 mb-2">
        <div class="font-semibold truncate">
          {{ report.view?.module.name }}
          <span v-if="form !== 'total'" class="text-muted-color font-normal">
            · {{ dimensionName }}
          </span>
        </div>
        <div v-if="form !== 'total'" class="flex gap-1 shrink-0">
          <Button
            :label="t('mcpApp.recordReport.chart')"
            size="small"
            :severity="mode === 'chart' ? 'primary' : 'secondary'"
            :text="mode !== 'chart'"
            @click="mode = 'chart'"
          />
          <Button
            :label="t('mcpApp.recordReport.table')"
            size="small"
            :severity="mode === 'table' ? 'primary' : 'secondary'"
            :text="mode !== 'table'"
            @click="mode = 'table'"
          />
        </div>
      </div>

      <div v-if="form === 'total'" id="status" class="flex flex-wrap gap-6">
        <div v-for="m in metrics" :key="m.key">
          <div class="text-muted-color text-xs">{{ metricTitle(m) }}</div>
          <div class="text-3xl font-semibold tabular-nums">
            {{ format(report.rows[0]?.[m.key], m) }}
          </div>
        </div>
      </div>

      <CResourceTable
        v-else-if="mode === 'table'"
        :items="tableRows"
        :fields="tableFields"
        size="small"
      />

      <div v-else id="status" class="flex flex-col gap-4">
        <div v-for="m in metrics" :key="m.key">
          <div v-if="metrics.length > 1 || m.key !== 'count'" class="text-muted-color text-xs mb-1">
            {{ metricTitle(m) }}
          </div>
          <div :style="{ height: chartHeight() }">
            <CChart :chart="optionFor(m)" />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
