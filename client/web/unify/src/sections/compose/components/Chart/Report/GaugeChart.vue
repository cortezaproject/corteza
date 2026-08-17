<template>
  <ReportEdit
    :chart="chart"
    :modules="modules"
    :supported-metrics="1"
    :uses-dimensions-field="false"
  >
    <template #dimension-options="{ dimension }">
      <Fieldset :legend="$t('chart.edit.dimension.gaugeSteps')" class="mt-4">
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
      </Fieldset>
    </template>

    <template #metric-options="{ metric, index }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <CFormGroup :label="$t('chart.edit.metric.labelLabel')" :input-id="`gaugeLabel${index}`">
          <InputText :id="`gaugeLabel${index}`" v-model="metric.label" class="w-full" />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.options.label')">
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="metric.fixTooltips"
              :binary="true"
              :input-id="`gaugeFixTooltips${index}`"
            />
            <label :for="`gaugeFixTooltips${index}`">
              {{ $t('chart.edit.metric.fixTooltips') }}
            </label>
          </div>
        </CFormGroup>

        <CFormGroup
          :label="$t('chart.edit.metric.angle.start')"
          :input-id="`gaugeStartAngle${index}`"
        >
          <InputNumber
            :input-id="`gaugeStartAngle${index}`"
            v-model="metric.startAngle"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.metric.angle.end')" :input-id="`gaugeEndAngle${index}`">
          <InputNumber
            :input-id="`gaugeEndAngle${index}`"
            v-model="metric.endAngle"
            class="w-full"
          />
        </CFormGroup>
      </div>

      <NumberFormatting v-model="metric.formatting" :id-prefix="`gauge${index}`" />
    </template>
  </ReportEdit>
</template>

<script setup>
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
})

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
