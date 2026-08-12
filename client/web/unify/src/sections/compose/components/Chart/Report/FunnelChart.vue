<template>
  <ReportEdit :chart="chart" :modules="modules" :supported-metrics="supportedMetrics">
    <template #dimension-options="{ index, dimension, field }">
      <div v-if="showPicker(field)" class="grid grid-cols-1 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.dimension.options.label') }}
          </label>
          <MultiSelect
            :model-value="getOptions(dimension)"
            :options="field.options.options"
            option-label="value"
            option-value="value"
            class="w-full"
            filter
            @update:model-value="v => setOptions(index, field, v)"
          />
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
            {{ $t('chart.edit.metric.options.label') }}
          </label>
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox v-model="metric.fixTooltips" :binary="true" input-id="funnelFixTooltips" />
              <label for="funnelFixTooltips">{{ $t('chart.edit.metric.fixTooltips') }}</label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="metric.relativeValue"
                :binary="true"
                input-id="funnelRelativeValue"
              />
              <label for="funnelRelativeValue">{{ $t('chart.edit.metric.relative') }}</label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox v-model="metric.cumulative" :binary="true" input-id="funnelCumulative" />
              <label for="funnelCumulative">{{ $t('chart.edit.metric.cumulative') }}</label>
            </div>
          </div>
        </div>
      </div>

      <Divider />

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
            {{
              $t(
                `chart.edit.formatting.presetFormats.description.${metric.formatting.presetFormat}`,
              )
            }}
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
  </ReportEdit>
</template>

<script setup>
import { inject } from 'vue'
import { useI18n } from 'vue-i18n'
import ReportEdit from './ReportEdit.vue'

const { t } = useI18n()

const formatOptions = [
  { value: 'custom', text: t('chart.edit.formatting.presetFormats.options.custom') },
  { value: 'accounting', text: t('chart.edit.formatting.presetFormats.options.accounting') },
]

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
