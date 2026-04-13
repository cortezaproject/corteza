<template>
  <ReportEdit
    :report="report"
    :chart="chart"
    :modules="modules"
    :supported-metrics="-1"
    @update:report="$emit('update:report', $event)"
  >
    <template #dimension-options="{ dimension }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.radar.shape.label') }}
          </label>
          <Select
            v-model="dimension.shape"
            :options="radarShapes"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.options.label') }}
          </label>
          <div class="flex items-center gap-2">
            <Checkbox v-model="dimension.fixTooltips" :binary="true" input-id="radarFixTooltips" />
            <label for="radarFixTooltips">{{ $t('chart.edit.metric.fixTooltips') }}</label>
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
  </ReportEdit>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import ReportEdit from './ReportEdit.vue'

const { t } = useI18n()

const formatOptions = [
  { value: 'custom', text: t('chart.edit.formatting.presetFormats.options.custom') },
  { value: 'accounting', text: t('chart.edit.formatting.presetFormats.options.accounting') },
]

defineProps({
  report: {
    type: Object,
    required: true,
  },
  chart: {
    type: Object,
    default: () => ({}),
  },
  modules: {
    type: Array,
    required: true,
  },
})

defineEmits(['update:report'])

const radarShapes = [
  { value: 'polygon', text: t('chart.edit.metric.radar.shape.polygon') },
  { value: 'circle', text: t('chart.edit.metric.radar.shape.circle') },
]
</script>
