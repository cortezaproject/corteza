<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
  >
    <template #metric-options="{ metric, report }">
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
              placeholder="0.00"
            />
          </b-form-group>
        </b-col>
      </b-row>

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

    <template #additional-config="{ report, presetFormattedOptions }">
      <hr>
      <div class="px-3">
        <h5 class="mb-3">
          {{ $t('edit.additionalConfig.tooltip.label') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('numberFormat')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.numberFormat"
                placeholder="0.00"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('prefix')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.prefix"
                placeholder="USD/mo"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('suffix')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.suffix"
                placeholder="$"
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
              :description="presetFormattedOptions.formattedOptionsDescription"
            >
              <b-form-select
                v-model="report.tooltipFormatter.presetFormat"
                :options="presetFormattedOptions.formatOptions"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>
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
