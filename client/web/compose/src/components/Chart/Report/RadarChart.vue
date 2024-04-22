<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
  >
    <template #metric-options="{ metric, presetFormattedOptions }">
      <b-row>
        <b-col
          cols="12"
          lg="6"
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

      <hr>

      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('prefix')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.prefix"
              placeholder="USD/mo"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('suffix')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.suffix"
              placeholder="$"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('numberFormat')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.numberFormat"
              :disabled="report.metricFormatter.presetFormat !== 'noFormat'"
              placeholder="0.00"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.tooltip.formatting.presetFormats.label')"
            label-class="text-primary"
            style="white-space: pre-line;"
          >
            <b-form-select
              v-model="report.metricFormatter.presetFormat"
              :disabled="!!report.metricFormatter.numberFormat"
              :options="presetFormattedOptions.formatOptions"
            />
            <slot
              v-if="report.metricFormatter.presetFormat === 'accountingNumber'"
              name="description"
            >
              <small class="text-muted">{{ presetFormattedOptions.formattedOptionsDescription }}</small>
            </slot>
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
          lg="6"
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
