<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
    :dimension-field-kind="['Select']"
    :supported-metrics="1"
    un-skippable
  >
    <template #dimension-options="{ dimension }">
      <b-form-group
        :label="$t('edit.dimension.gaugeSteps')"
        :label-cols="2"
      >
        <b-input-group
          v-for="(step, i) in dimension.meta.steps"
          :key="i"
          class="mb-1"
        >
          <b-form-input
            v-model="step.label"
            plain
            class="w-50"
            :placeholder="$t('general.label.title')"
          />
          <b-input-group-append>
            <chart-translator
              :field.sync="step.label"
              :chart="chart"
              :disabled="isNew"
              :highlight-key="`dimensions.${dimension.dimensionID}.meta.steps.${step.stepID}.label`"
              button-variant="light"
            />
          </b-input-group-append>

          <b-form-input
            v-model="step.value"
            type="number"
            class="text-right w-25"
            :placeholder="$t('general.value')"
          />

          <b-input-group-append>
            <b-button
              variant="link"
              class="border-0 text-danger"
              @click.prevent="dimension.meta.steps.splice(i, 1)"
            >
              <font-awesome-icon :icon="['far', 'trash-alt']" />
            </b-button>
          </b-input-group-append>
        </b-input-group>

        <b-btn
          variant="link"
          class="p-0"
          @click="dimension.meta.steps.push({ label: undefined, color: undefined, value: undefined })"
        >
          + {{ $t('general.label.add') }}
        </b-btn>
      </b-form-group>
    </template>

    <template #metric-options="{ metric }">
      <b-form-group
        :label="$t('edit.metric.fx.label')"
        :description="$t('edit.metric.fx.description')"
      >
        <b-form-textarea
          v-model="metric.fx"
          placeholder="n"
        />
      </b-form-group>

      <b-form-group
        horizontal
        :label="$t('edit.metric.gaugeColor')"
        :label-cols="2"
      >
        <b-form-input
          v-model="metric.backgroundColor"
          type="color"
          class="p-0 border-0 color-picker mb-2"
        />

        <b-form-checkbox
          v-model="metric.showTooltipLabel"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('edit.metric.showTooltipLabel') }}
        </b-form-checkbox>

        <b-form-checkbox
          v-model="metric.fixTooltips"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('edit.metric.fixTooltips') }}
        </b-form-checkbox>
      </b-form-group>
    </template>
  </report-edit>
</template>

<script>
import ReportEdit from './ReportEdit'
import ChartTranslator from 'corteza-webapp-compose/src/components/Chart/ChartTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'chart',
  },

  components: {
    ChartTranslator,
    ReportEdit,
  },

  extends: base,

  props: {
    chart: {
      type: compose.Chart,
      required: true,
    },
  },

  computed: {
    isNew () {
      return this.chart.chartID === NoID
    },
  },

  methods: {
    renderPicker (field) {
      return field && field.kind === 'Select' && field.options.options
    },

    makeOptions (field, dimension) {
      const ff = dimension.meta.fields || []
      return field.options.options.map(({ text, value }) => ({
        text,
        value,
        color: (ff.find(b => value === b.value) || {}).color,
      }))
    },
  },
}
</script>

<style scoped lang="scss">
  .color-picker {
    width: 30px;
  }
</style>
