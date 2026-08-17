<template>
  <ReportEdit :chart="chart" :modules="modules" :supported-metrics="supportedMetrics">
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
        <CFormGroup
          :label="$t('chart.edit.dimension.rotate.label')"
          :description="$t('chart.edit.dimension.rotate.description')"
          input-id="dimensionRotate"
        >
          <InputNumber input-id="dimensionRotate" v-model="dimension.rotateLabel" class="w-full" />
        </CFormGroup>
      </div>
    </template>

    <template #y-axis="{ report: r }">
      <Panel :header="$t('chart.edit.yAxis.label')" toggleable collapsed>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <CFormGroup :label="$t('chart.edit.yAxis.labelLabel')" input-id="yAxisLabel">
            <InputText id="yAxisLabel" v-model="r.yAxis.label" class="w-full" />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.yAxis.labelPosition.label')"
            input-id="yAxisLabelPosition"
          >
            <Select
              id="yAxisLabelPosition"
              v-model="r.yAxis.labelPosition"
              :options="axisLabelPositions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </CFormGroup>

          <CFormGroup :label="$t('chart.edit.yAxis.minLabel')" input-id="yAxisMin">
            <InputNumber
              input-id="yAxisMin"
              v-model="r.yAxis.min"
              :placeholder="$t('chart.edit.yAxis.minPlaceholder')"
              class="w-full"
            />
          </CFormGroup>

          <CFormGroup :label="$t('chart.edit.yAxis.maxLabel')" input-id="yAxisMax">
            <InputNumber
              input-id="yAxisMax"
              v-model="r.yAxis.max"
              :placeholder="$t('chart.edit.yAxis.maxPlaceholder')"
              class="w-full"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.yAxis.rotate.label')"
            :description="$t('chart.edit.yAxis.rotate.description')"
            input-id="yAxisRotate"
          >
            <InputNumber input-id="yAxisRotate" v-model="r.yAxis.rotateLabel" class="w-full" />
          </CFormGroup>

          <CFormGroup :label="$t('chart.edit.yAxis.options.label')">
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
          </CFormGroup>
        </div>

        <NumberFormatting v-model="r.yAxis.formatting" id-prefix="yAxis" />
      </Panel>
    </template>

    <template #metric-options="{ metric, index }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <CFormGroup :label="$t('chart.edit.metric.labelLabel')" :input-id="`metricLabel${index}`">
          <InputText :id="`metricLabel${index}`" v-model="metric.label" class="w-full" />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.output.label')" :input-id="`metricType${index}`">
          <Select
            :id="`metricType${index}`"
            v-model="metric.type"
            :options="chartTypes"
            option-label="text"
            option-value="value"
            :placeholder="$t('chart.edit.metric.output.placeholder')"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('chart.edit.metric.fx.label')"
          :description="$t('chart.edit.metric.fx.description')"
          :input-id="`metricFx${index}`"
        >
          <Textarea
            :id="`metricFx${index}`"
            v-model="metric.fx"
            placeholder="n"
            rows="2"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.options.label')">
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="metric.fixTooltips"
                :binary="true"
                :input-id="`fixTooltips${index}`"
              />
              <label :for="`fixTooltips${index}`">{{ $t('chart.edit.metric.fixTooltips') }}</label>
            </div>
            <div v-if="hasRelativeDisplay(metric)" class="flex items-center gap-2">
              <Checkbox
                v-model="metric.relativeValue"
                :binary="true"
                :input-id="`relativeValue${index}`"
              />
              <label :for="`relativeValue${index}`">{{ $t('chart.edit.metric.relative') }}</label>
            </div>
            <div v-if="metric.type === 'pie'" class="flex items-center gap-2">
              <Checkbox v-model="metric.rose" :binary="true" :input-id="`rose${index}`" />
              <label :for="`rose${index}`">{{ $t('chart.edit.metric.rose') }}</label>
            </div>
            <div v-if="metric.type === 'line'" class="flex items-center gap-2">
              <Checkbox v-model="metric.fill" :binary="true" :input-id="`fill${index}`" />
              <label :for="`fill${index}`">{{ $t('chart.edit.metric.fillArea') }}</label>
            </div>
          </div>
        </CFormGroup>

        <CFormGroup
          v-if="!hasRelativeDisplay(metric)"
          :label="$t('chart.edit.metric.stack.label')"
          :description="$t('chart.edit.metric.stack.description')"
          :input-id="`metricStack${index}`"
        >
          <InputText :id="`metricStack${index}`" v-model="metric.stack" class="w-full" />
        </CFormGroup>

        <CFormGroup
          v-if="metric.type === 'line'"
          :label="$t('chart.edit.metric.lineStyle.label')"
          :input-id="`metricLineStyle${index}`"
        >
          <Select
            :id="`metricLineStyle${index}`"
            :model-value="getLineStyle(metric)"
            :options="lineStyleOptions"
            option-label="text"
            option-value="value"
            class="w-full"
            @change="e => setLineStyle(e.value, metric)"
          />
        </CFormGroup>

        <CFormGroup
          v-if="metric.type === 'scatter'"
          :label="$t('chart.edit.metric.symbol.label')"
          :input-id="`metricSymbol${index}`"
        >
          <Select
            :id="`metricSymbol${index}`"
            v-model="metric.symbol"
            :options="scatterSymbolOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>
      </div>

      <!-- Metric formatting -->
      <NumberFormatting v-model="metric.formatting" :id-prefix="`metric${index}`" />
    </template>

    <template #additional-config="{ report: r, hasAxis }">
      <Panel :header="$t('chart.edit.additionalConfig.tooltip.label')" toggleable collapsed>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <CFormGroup
            :label="$t('chart.edit.additionalConfig.tooltip.formatting.label')"
            :description="$t('chart.edit.additionalConfig.tooltip.formatting.description')"
            input-id="tooltipFormatting"
          >
            <InputText
              id="tooltipFormatting"
              v-model="r.tooltip.formatting"
              :placeholder="$t('chart.edit.additionalConfig.tooltip.formatting.placeholder')"
              class="w-full"
            />
          </CFormGroup>

          <CFormGroup
            v-if="!hasAxis"
            :label="$t('chart.edit.additionalConfig.tooltip.labelNextToChart')"
            input-id="labelsNextToPartition"
          >
            <ToggleSwitch
              input-id="labelsNextToPartition"
              v-model="r.tooltip.labelsNextToPartition"
            />
          </CFormGroup>
        </div>
      </Panel>

      <Panel :header="$t('chart.edit.additionalConfig.offset.label')" toggleable collapsed>
        <CFormGroup
          :label="$t('chart.edit.additionalConfig.offset.default')"
          input-id="offsetDefault"
          class="mb-4"
        >
          <ToggleSwitch input-id="offsetDefault" v-model="r.offset.isDefault" />
        </CFormGroup>

        <div v-if="!r.offset.isDefault" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <CFormGroup
            :label="$t('chart.edit.additionalConfig.offset.position.top')"
            input-id="offsetTop"
          >
            <InputText id="offsetTop" v-model="r.offset.top" class="w-full" />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.additionalConfig.offset.position.right')"
            input-id="offsetRight"
          >
            <InputText id="offsetRight" v-model="r.offset.right" class="w-full" />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.additionalConfig.offset.position.bottom')"
            input-id="offsetBottom"
          >
            <InputText id="offsetBottom" v-model="r.offset.bottom" class="w-full" />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.additionalConfig.offset.position.left')"
            input-id="offsetLeft"
          >
            <InputText id="offsetLeft" v-model="r.offset.left" class="w-full" />
          </CFormGroup>

          <small class="text-muted-color col-span-1 lg:col-span-2">
            {{ $t('chart.edit.additionalConfig.offset.valueRange') }}
          </small>
        </div>
      </Panel>
    </template>
  </ReportEdit>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import ReportEdit from './ReportEdit.vue'
import NumberFormatting from './NumberFormatting.vue'

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
