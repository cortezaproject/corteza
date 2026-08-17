<template>
  <ReportEdit :chart="chart" :modules="modules" :supported-metrics="supportedMetrics">
    <template #dimension-options="{ index, dimension, field }">
      <div v-if="showPicker(field)" class="grid grid-cols-1 gap-4 mt-4">
        <CFormGroup
          :label="$t('chart.edit.dimension.options.label')"
          :input-id="`funnelOptions${index}`"
        >
          <MultiSelect
            :id="`funnelOptions${index}`"
            :model-value="getOptions(dimension)"
            :options="field.options.options"
            option-label="value"
            option-value="value"
            class="w-full"
            filter
            @update:model-value="v => setOptions(index, field, v)"
          />
        </CFormGroup>
      </div>
    </template>

    <template #metric-options="{ metric, index }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <CFormGroup :label="$t('chart.edit.metric.labelLabel')" :input-id="`funnelLabel${index}`">
          <InputText :id="`funnelLabel${index}`" v-model="metric.label" class="w-full" />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.options.label')">
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="metric.fixTooltips"
                :binary="true"
                :input-id="`funnelFixTooltips${index}`"
              />
              <label :for="`funnelFixTooltips${index}`">
                {{ $t('chart.edit.metric.fixTooltips') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="metric.relativeValue"
                :binary="true"
                :input-id="`funnelRelativeValue${index}`"
              />
              <label :for="`funnelRelativeValue${index}`">
                {{ $t('chart.edit.metric.relative') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="metric.cumulative"
                :binary="true"
                :input-id="`funnelCumulative${index}`"
              />
              <label :for="`funnelCumulative${index}`">
                {{ $t('chart.edit.metric.cumulative') }}
              </label>
            </div>
          </div>
        </CFormGroup>
      </div>

      <NumberFormatting v-model="metric.formatting" :id-prefix="`funnel${index}`" />
    </template>
  </ReportEdit>
</template>

<script setup>
import { inject } from 'vue'
import ReportEdit from './ReportEdit.vue'
import NumberFormatting from './NumberFormatting.vue'

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

function showPicker(field) {
  return field && field.kind === 'Select' && field.options?.options
}
function getOptions(dimension) {
  const fields = dimension.meta?.fields || []
  return fields.map(f => f.value)
}
function setOptions(index, field, values) {
  if (!report.value.dimensions[index].meta) {
    report.value.dimensions[index].meta = {}
  }
  const options = field.options?.options || []
  report.value.dimensions[index].meta.fields = values
    .map(v => options.find(o => o.value === v))
    .filter(Boolean)
}
</script>
