<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
    :supported-metrics="1"
    :uses-dimensions-field="false"
    un-skippable
  >
    <template #dimension-options="{ dimension }">
      <b-form-group
        :label="$t('edit.dimension.gaugeSteps')"
        label-class="text-primary"
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
            />
          </b-input-group-append>

          <b-form-input
            v-model="step.value"
            type="number"
            class="text-right w-25"
            :placeholder="$t('general.value')"
          />

          <b-input-group-append>
            <c-input-confirm
              show-icon
              @confirmed="dimension.meta.steps.splice(i, 1)"
            />
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
      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.fx.label')"
            :description="$t('edit.metric.fx.description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="metric.fx"
              placeholder="n"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.options.label')"
            label-class="text-primary"
          >
            <b-form-checkbox
              v-model="metric.fixTooltips"
            >
              {{ $t('edit.metric.fixTooltips') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.angle.start')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="metric.startAngle"
              type="number"
              number
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.angle.end')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="metric.endAngle"
              type="number"
              number
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>
  </report-edit>
</template>

<script>
import ReportEdit from './ReportEdit'
import ChartTranslator from 'corteza-webapp-compose/src/components/Chart/ChartTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  name: 'GaugeChart',

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
}
</script>
