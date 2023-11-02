<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
  >
    <template #metric-options="{ metric }">
      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.labelLabel')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="metric.label"
              />
              <b-input-group-append>
                <chart-translator
                  :field.sync="metric.label"
                  :chart="chart"
                  :disabled="isNew"
                  :highlight-key="`metrics.${metric.metricID}.label`"
                />
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
        </b-col>
      </b-row>
    </template>

    <template #dimension-options-options="{ dimension }">
      <b-form-checkbox
        v-model="dimension.fixTooltips"
      >
        {{ $t('edit.metric.fixTooltips') }}
      </b-form-checkbox>
    </template>

    <template #dimension-options="{ dimension }">
      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.radar.shape.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="dimension.shape"
              :options="shapeOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>
  </report-edit>
</template>

<script>
import base from './base'
import ReportEdit from './ReportEdit'
import ChartTranslator from 'corteza-webapp-compose/src/components/Chart/ChartTranslator'

export default {
  name: 'RadarChart',

  components: {
    ReportEdit,
    ChartTranslator,
  },

  extends: base,

  data () {
    return {
      shapeOptions: [
        { value: 'polygon', text: this.$t('edit.metric.radar.shape.polygon') },
        { value: 'circle', text: this.$t('edit.metric.radar.shape.circle') },
      ],
    }
  },
}
</script>

<style scoped lang="scss">
  .cursor-pointer {
    cursor: pointer;
  }
</style>
