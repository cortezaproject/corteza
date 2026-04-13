<template>
  <ReportEdit
    :report="report"
    :chart="chart"
    :modules="modules"
    :supported-metrics="1"
    :uses-dimensions-field="false"
    @update:report="$emit('update:report', $event)"
  >
    <template #dimension-options="{ dimension }">
      <div class="px-0 mt-4">
        <h5 class="mb-3">
          {{ $t('chart.edit.dimension.gaugeSteps') }}
        </h5>

        <div
          v-for="(step, si) in dimension.meta?.steps || []"
          :key="si"
          class="flex items-center gap-2 mb-2"
        >
          <InputText
            v-model="step.label"
            :placeholder="$t('chart.edit.metric.labelLabel')"
            class="flex-1"
          />
          <InputNumber
            v-model="step.value"
            :placeholder="$t('chart.general.value')"
            class="flex-1"
          />
          <InputText
            v-model="step.color"
            :placeholder="$t('chart.edit.metric.gaugeColor')"
            class="w-24"
          />
          <Button
            icon="pi pi-trash"
            text
            severity="danger"
            size="small"
            @click="removeStep(dimension, si)"
          />
        </div>

        <Button
          :label="'+ ' + $t('chart.general.label.add')"
          text
          size="small"
          @click="addStep(dimension)"
        />
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
          <div class="flex items-center gap-2">
            <Checkbox v-model="metric.fixTooltips" :binary="true" input-id="gaugeFixTooltips" />
            <label for="gaugeFixTooltips">{{ $t('chart.edit.metric.fixTooltips') }}</label>
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.angle.start') }}
          </label>
          <InputNumber v-model="metric.startAngle" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.angle.end') }}
          </label>
          <InputNumber v-model="metric.endAngle" class="w-full" />
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

function addStep(dimension) {
  if (!dimension.meta) dimension.meta = {}
  if (!dimension.meta.steps) dimension.meta.steps = []
  dimension.meta.steps.push({ label: '', value: 0, color: '' })
}

function removeStep(dimension, index) {
  if (dimension.meta?.steps) {
    dimension.meta.steps.splice(index, 1)
  }
}
</script>
