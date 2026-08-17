<template>
  <ReportEdit :chart="chart" :modules="modules" :supported-metrics="-1">
    <template #dimension-options="{ index, dimension }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <CFormGroup
          :label="$t('chart.edit.metric.radar.shape.label')"
          :input-id="`radarShape${index}`"
        >
          <Select
            :id="`radarShape${index}`"
            v-model="dimension.shape"
            :options="radarShapes"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.options.label')">
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="dimension.fixTooltips"
              :binary="true"
              :input-id="`radarFixTooltips${index}`"
            />
            <label :for="`radarFixTooltips${index}`">
              {{ $t('chart.edit.metric.fixTooltips') }}
            </label>
          </div>
        </CFormGroup>
      </div>
    </template>

    <template #metric-options="{ metric, index }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <CFormGroup :label="$t('chart.edit.metric.labelLabel')" :input-id="`radarLabel${index}`">
          <InputText :id="`radarLabel${index}`" v-model="metric.label" class="w-full" />
        </CFormGroup>
      </div>

      <NumberFormatting v-model="metric.formatting" :id-prefix="`radar${index}`" />
    </template>
  </ReportEdit>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
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
})

const radarShapes = [
  { value: 'polygon', text: t('chart.edit.metric.radar.shape.polygon') },
  { value: 'circle', text: t('chart.edit.metric.radar.shape.circle') },
]
</script>
