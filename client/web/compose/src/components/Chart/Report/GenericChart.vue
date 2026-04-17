<template>
  <ReportEdit
    :chart="chart"
    :modules="modules"
    :supported-metrics="supportedMetrics"
  >
    <template #dimension-options-options="{ dimension, isTemporal }">
      <div
        v-if="isTemporal && !['WEEK', 'QUARTER'].includes(dimension.modifier)"
        class="flex items-center gap-2"
      >
        <Checkbox v-model="dimension.timeLabels" :binary="true" input-id="timeLabels" />
        <label for="timeLabels">{{ $t('chart.edit.dimension.timeLabels') }}</label>
      </div>
    </template>

    <template #dimension-options="{ dimension }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.dimension.rotate.label') }}
          </label>
          <InputNumber v-model="dimension.rotateLabel" class="w-full" />
          <small class="text-muted-color">
            {{ $t('chart.edit.dimension.rotate.description') }}
          </small>
        </div>
      </div>
    </template>

    <template #y-axis="{ report: r }">
      <div class="px-3">
        <h5 class="mb-3">
          {{ $t('chart.edit.yAxis.label') }}
        </h5>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.labelLabel') }}
            </label>
            <InputText v-model="r.yAxis.label" class="w-full" />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.labelPosition.label') }}
            </label>
            <Select
              v-model="r.yAxis.labelPosition"
              :options="axisLabelPositions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.minLabel') }}
            </label>
            <InputNumber
              v-model="r.yAxis.min"
              :placeholder="$t('chart.edit.yAxis.minPlaceholder')"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.maxLabel') }}
            </label>
            <InputNumber
              v-model="r.yAxis.max"
              :placeholder="$t('chart.edit.yAxis.maxPlaceholder')"
              class="w-full"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.rotate.label') }}
            </label>
            <InputNumber v-model="r.yAxis.rotateLabel" class="w-full" />
            <small class="text-muted-color">{{ $t('chart.edit.yAxis.rotate.description') }}</small>
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.yAxis.options.label') }}
            </label>
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <Checkbox v-model="logScale" :binary="true" input-id="logScale" />
                <label for="logScale">{{ $t('chart.edit.yAxis.logarithmicScale') }}</label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox v-model="axisRight" :binary="true" input-id="axisRight" />
                <label for="axisRight">{{ $t('chart.edit.yAxis.axisOnRight') }}</label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox v-model="r.yAxis.beginAtZero" :binary="true" input-id="beginAtZero" />
                <label for="beginAtZero">{{ $t('chart.edit.yAxis.axisScaleFromZero') }}</label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox v-model="r.yAxis.horizontal" :binary="true" input-id="horizontal" />
                <label for="horizontal">{{ $t('chart.edit.yAxis.horizontal.label') }}</label>
              </div>
            </div>
          </div>
        </div>

        <Divider />

        <!-- Y-axis formatting -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.formatting.prefix.label') }}
            </label>
            <InputText
              v-model="r.yAxis.formatting.prefix"
              :placeholder="$t('chart.edit.formatting.prefix.placeholder')"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.formatting.suffix.label') }}
            </label>
            <InputText
              v-model="r.yAxis.formatting.suffix"
              :placeholder="$t('chart.edit.formatting.suffix.placeholder')"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.formatting.presetFormats.label') }}
            </label>
            <Select
              v-model="r.yAxis.formatting.presetFormat"
              :options="formatOptions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
            <small v-if="r.yAxis.formatting.presetFormat" class="text-muted-color whitespace-pre-line">
              {{ $t(`chart.edit.formatting.presetFormats.description.${r.yAxis.formatting.presetFormat}`) }}
            </small>
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.formatting.format.label') }}
            </label>
            <InputText 
              v-model="r.yAxis.formatting.format" 
              :disabled="r.yAxis.formatting.presetFormat !== 'custom'" 
              :placeholder="$t('chart.edit.formatting.format.placeholder')" 
              class="w-full" 
            />
          </div>
        </div>
      </div>
    </template>

    <template #metric-options="{ metric }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.labelLabel') }}
          </label>
          <InputText v-model="metric.label" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.output.label') }}
          </label>
          <Select
            v-model="metric.type"
            :options="chartTypes"
            option-label="text"
            option-value="value"
            :placeholder="$t('chart.edit.metric.output.placeholder')"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.fx.label') }}
          </label>
          <Textarea v-model="metric.fx" placeholder="n" rows="2" class="w-full" />
          <small class="text-muted-color">{{ $t('chart.edit.metric.fx.description') }}</small>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.options.label') }}
          </label>
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox v-model="metric.fixTooltips" :binary="true" input-id="fixTooltips" />
              <label for="fixTooltips">{{ $t('chart.edit.metric.fixTooltips') }}</label>
            </div>
            <div v-if="hasRelativeDisplay(metric)" class="flex items-center gap-2">
              <Checkbox v-model="metric.relativeValue" :binary="true" input-id="relativeValue" />
              <label for="relativeValue">{{ $t('chart.edit.metric.relative') }}</label>
            </div>
            <div v-if="metric.type === 'pie'" class="flex items-center gap-2">
              <Checkbox v-model="metric.rose" :binary="true" input-id="rose" />
              <label for="rose">{{ $t('chart.edit.metric.rose') }}</label>
            </div>
            <div v-if="metric.type === 'line'" class="flex items-center gap-2">
              <Checkbox v-model="metric.fill" :binary="true" input-id="fill" />
              <label for="fill">{{ $t('chart.edit.metric.fillArea') }}</label>
            </div>
          </div>
        </div>

        <div v-if="!hasRelativeDisplay(metric)" class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.stack.label') }}
          </label>
          <InputText v-model="metric.stack" class="w-full" />
          <small class="text-muted-color">{{ $t('chart.edit.metric.stack.description') }}</small>
        </div>

        <div v-if="metric.type === 'line'" class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.lineStyle.label') }}
          </label>
          <Select
            :model-value="getLineStyle(metric)"
            :options="lineStyleOptions"
            option-label="text"
            option-value="value"
            class="w-full"
            @change="e => setLineStyle(e.value, metric)"
          />
        </div>

        <div v-if="metric.type === 'scatter'" class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.symbol.label') }}
          </label>
          <Select
            v-model="metric.symbol"
            :options="scatterSymbolOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </div>
      </div>

      <Divider />

      <!-- Metric formatting -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.prefix.label') }}
          </label>
          <InputText
            v-model="metric.formatting.prefix"
            :placeholder="$t('chart.edit.formatting.prefix.placeholder')"
            class="w-full"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.suffix.label') }}
          </label>
          <InputText
            v-model="metric.formatting.suffix"
            :placeholder="$t('chart.edit.formatting.suffix.placeholder')"
            class="w-full"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.presetFormats.label') }}
          </label>
          <Select
            v-model="metric.formatting.presetFormat"
            :options="formatOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
          <small v-if="metric.formatting.presetFormat" class="text-muted-color whitespace-pre-line">
            {{ $t(`chart.edit.formatting.presetFormats.description.${metric.formatting.presetFormat}`) }}
          </small>
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.format.label') }}
          </label>
          <InputText 
            v-model="metric.formatting.format" 
            :disabled="metric.formatting.presetFormat !== 'custom'" 
            :placeholder="$t('chart.edit.formatting.format.placeholder')" 
            class="w-full" 
          />
        </div>
      </div>
    </template>

    <template #additional-config="{ report: r, hasAxis }">
      <Divider />
      <div class="px-3">
        <h5 class="mb-3">
          {{ $t('chart.edit.additionalConfig.tooltip.label') }}
        </h5>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.additionalConfig.tooltip.formatting.label') }}
            </label>
            <InputText
              v-model="r.tooltip.formatting"
              :placeholder="$t('chart.edit.additionalConfig.tooltip.formatting.placeholder')"
              class="w-full"
            />
            <small class="text-muted-color">
              {{ $t('chart.edit.additionalConfig.tooltip.formatting.description') }}
            </small>
          </div>

          <div v-if="!hasAxis" class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.additionalConfig.tooltip.labelNextToChart') }}
            </label>
            <div class="flex items-center gap-2 mt-2">
              <ToggleSwitch v-model="r.tooltip.labelsNextToPartition" />
            </div>
          </div>
        </div>
      </div>

      <Divider />

      <div class="px-3 mb-2">
        <h5 class="mb-3">
          {{ $t('chart.edit.additionalConfig.offset.label') }}
        </h5>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.additionalConfig.offset.default') }}
            </label>
            <div class="flex items-center gap-2 mt-2">
              <ToggleSwitch v-model="r.offset.isDefault" />
            </div>
          </div>
        </div>

        <div v-if="!r.offset.isDefault" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('chart.edit.additionalConfig.offset.position.top') }}</label>
            <InputText v-model="r.offset.top" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('chart.edit.additionalConfig.offset.position.right') }}</label>
            <InputText v-model="r.offset.right" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('chart.edit.additionalConfig.offset.position.bottom') }}</label>
            <InputText v-model="r.offset.bottom" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('chart.edit.additionalConfig.offset.position.left') }}</label>
            <InputText v-model="r.offset.left" class="w-full" />
          </div>
          <div class="col-span-1 lg:col-span-2">
            <small class="text-muted-color">{{ $t('chart.edit.additionalConfig.offset.valueRange') }}</small>
          </div>
        </div>
      </div>
    </template>
  </ReportEdit>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import ReportEdit from './ReportEdit.vue'

const { t } = useI18n()

defineProps({
  chart: {
    type: Object,
    default: () => ({}),
  },
  modules: {
    type: Array,
    required: true,
  },
  supportedMetrics: {
    type: Number,
    default: -1,
  },
})

const report = inject('reportDraft')

const ignoredCharts = ['funnel', 'gauge', 'radar']

const formatOptions = [
  { value: 'custom', text: t('chart.edit.formatting.presetFormats.options.custom') },
  { value: 'accounting', text: t('chart.edit.formatting.presetFormats.options.accounting') },
]

const chartTypes = Object.values(compose.chartUtil.ChartType)
  .filter(v => !ignoredCharts.includes(v))
  .map(value => ({ value, text: t(`chart.edit.metric.output.${value}`) }))

const axisLabelPositions = [
  { value: 'end', text: t('chart.edit.yAxis.labelPosition.top') },
  { value: 'center', text: t('chart.edit.yAxis.labelPosition.center') },
  { value: 'start', text: t('chart.edit.yAxis.labelPosition.bottom') },
]

const lineStyleOptions = [
  { value: '', text: t('chart.edit.metric.lineStyle.default') },
  { value: 'smooth', text: t('chart.edit.metric.lineStyle.smooth') },
  { value: 'step', text: t('chart.edit.metric.lineStyle.step') },
]

const scatterSymbolOptions = [
  { value: 'circle', text: t('chart.edit.metric.symbol.circle') },
  { value: 'triangle', text: t('chart.edit.metric.symbol.triangle') },
  { value: 'diamond', text: t('chart.edit.metric.symbol.diamond') },
  { value: 'pin', text: t('chart.edit.metric.symbol.pin') },
  { value: 'arrow', text: t('chart.edit.metric.symbol.arrow') },
  { value: 'rect', text: t('chart.edit.metric.symbol.rect') },
  { value: 'roundRect', text: t('chart.edit.metric.symbol.roundRect') },
]

const { hasRelativeDisplay } = compose.chartUtil

const logScale = computed({
  get: () => report.value.yAxis?.axisType === 'logarithmic',
  set: v => {
    if (report.value.yAxis) {
      report.value.yAxis.axisType = v ? 'logarithmic' : 'linear'
    }
  },
})

const axisRight = computed({
  get: () => report.value.yAxis?.axisPosition === 'right',
  set: v => {
    if (report.value.yAxis) {
      report.value.yAxis.axisPosition = v ? 'right' : 'left'
    }
  },
})

function getLineStyle(metric) {
  if (metric.smooth) return 'smooth'
  else if (metric.step) return 'step'
  return ''
}

function setLineStyle(style, metric) {
  metric.smooth = style === 'smooth'
  metric.step = style === 'step'
}
</script>
