<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
  >
    <template #dimension-options-options="{ dimension, isTemporal }">
      <b-form-checkbox
        v-if="isTemporal && !['WEEK', 'QUARTER'].includes(dimension.modifier)"
        v-model="dimension.timeLabels"
      >
        {{ $t('edit.dimension.timeLabels') }}
      </b-form-checkbox>
    </template>

    <template #dimension-options="{ dimension }">
      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.dimension.rotate.label')"
            :description="$t('edit.dimension.rotate.description')"
            label-class="text-primary"
          >
            <b-input
              v-model="dimension.rotateLabel"
              type="number"
              number
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>

    <template #y-axis="{ report }">
      <div
        class="px-3"
      >
        <h5 class="mb-3">
          {{ $t('edit.yAxis.label') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.yAxis.labelLabel')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input
                  v-model="report.yAxis.label"
                />
                <b-input-group-append>
                  <chart-translator
                    :field.sync="report.yAxis.label"
                    :chart="chart"
                    :disabled="isNew"
                    highlight-key="yAxis.label"
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
              :label="$t('edit.yAxis.labelPosition.label')"
              label-class="text-primary"
            >
              <b-form-select
                v-model="report.yAxis.labelPosition"
                :options="axisLabelPositions"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.yAxis.minLabel')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="report.yAxis.min"
                type="number"
                :placeholder="$t('edit.yAxis.minPlaceholder')"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.yAxis.maxLabel')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="report.yAxis.max"
                type="number"
                :placeholder="$t('edit.yAxis.maxPlaceholder')"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.yAxis.rotate.label')"
              :description="$t('edit.yAxis.rotate.description')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.yAxis.rotateLabel"
                type="number"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.yAxis.options.label')"
              label-class="text-primary"
            >
              <b-form-checkbox
                v-model="report.yAxis.axisType"
                value="logarithmic"
                unchecked-value="linear"
              >
                {{ $t('edit.yAxis.logarithmicScale') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="report.yAxis.axisPosition"
                value="right"
                unchecked-value="left"
              >
                {{ $t('edit.yAxis.axisOnRight') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="report.yAxis.beginAtZero"
              >
                {{ $t('edit.yAxis.axisScaleFromZero') }}
              </b-form-checkbox>
              <b-form-checkbox
                v-model="report.yAxis.horizontal"
              >
                {{ $t('edit.yAxis.horizontal.label') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>
        </b-row>
      </div>
    </template>

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

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.output.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="metric.type"
              :options="chartTypes"
              label="text"
              :reduce="option => option.value"
              :get-option-key="option => option.text"
              :placeholder="$t('edit.metric.output.placeholder')"
              @input="value => chartTypeChanged(metric)"
            />
          </b-form-group>
        </b-col>

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

            <b-form-checkbox
              v-if="hasRelativeDisplay(metric)"
              v-model="metric.relativeValue"
            >
              {{ $t('edit.metric.relative') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-if="metric.type === 'pie'"
              v-model="metric.rose"
            >
              {{ $t('edit.metric.rose') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-if="metric.type === 'line'"
              v-model="metric.fill"
            >
              {{ $t('edit.metric.fillArea') }}
            </b-form-checkbox>
          </b-form-group>

          <b-form-group
            v-if="metric.type === 'line'"
            :label="$t('edit.metric.lineStyle.label')"
            label-class="text-primary"
          >
            <b-form-radio-group
              :checked="getLineStyle(metric)"
              :options="lineStyleOptions"
              @change="setLineStyle($event, metric)"
            />
          </b-form-group>
        </b-col>

        <b-col
          v-if="!hasRelativeDisplay(metric)"
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.stack.label')"
            :description="$t('edit.metric.stack.description')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="metric.stack"
            />
          </b-form-group>
        </b-col>

        <b-col
          v-if="metric.type === 'scatter'"
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.symbol.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="metric.symbol"
              :options="scatterSymbolOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>

    <template #additional-config="{ hasAxis, report }">
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
              :label="$t('edit.additionalConfig.tooltip.formatting.label')"
              :description="$t('edit.additionalConfig.tooltip.formatting.description')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltip.formatting"
                :placeholder="$t('edit.additionalConfig.tooltip.formatting.placeholder')"
              />
            </b-form-group>
          </b-col>
          <b-col
            v-if="!hasAxis"
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.tooltip.labelNextToChart')"
              label-class="text-primary"
            >
              <c-input-checkbox
                :value="!!report.tooltip.labelsNextToPartition"
                switch
                :labels="checkboxLabel"
                @input="$set(report.tooltip, 'labelsNextToPartition', $event)"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <hr>
      <div class="px-3 mb-2">
        <h5 class="mb-3">
          {{ $t('edit.additionalConfig.offset.label') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.offset.default')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="report.offset.isDefault"
                switch
                :labels="checkboxLabel"
                class="mb-3"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row
          v-if="!report.offset.isDefault"
        >
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.offset.position.top')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.offset.top"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.offset.position.right')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.offset.right"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.offset.position.bottom')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.offset.bottom"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.offset.position.left')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.offset.left"
              />
            </b-form-group>
          </b-col>

          <b-col cols="12">
            <small class="text-muted">
              {{ $t('edit.additionalConfig.offset.valueRange') }}
            </small>
          </b-col>
        </b-row>
      </div>
    </template>
  </report-edit>
</template>

<script>
import ReportEdit from './ReportEdit'
import ChartTranslator from 'corteza-webapp-compose/src/components/Chart/ChartTranslator'
import { compose } from '@cortezaproject/corteza-js'
import base from './base'

const ignoredCharts = [
  'funnel',
  'gauge',
  'radar',
]

export default {
  name: 'GenericChart',

  i18nOptions: {
    namespaces: 'chart',
  },

  components: {
    ReportEdit,
    ChartTranslator,
  },

  extends: base,

  data () {
    return {
      chartTypes: Object.values(compose.chartUtil.ChartType)
        .filter(v => !ignoredCharts.includes(v))
        .map(value => ({ value, text: this.$t(`edit.metric.output.${value}`) })),

      legendPositions: [
        { value: 'top', text: this.$t('edit.metric.legend.top') },
        { value: 'left', text: this.$t('edit.metric.legend.left') },
        { value: 'bottom', text: this.$t('edit.metric.legend.bottom') },
        { value: 'right', text: this.$t('edit.metric.legend.right') },
      ],

      axisLabelPositions: [
        { value: 'end', text: this.$t('edit.yAxis.labelPosition.top') },
        { value: 'center', text: this.$t('edit.yAxis.labelPosition.center') },
        { value: 'start', text: this.$t('edit.yAxis.labelPosition.bottom') },
      ],

      tensionSteps: [
        { text: this.$t('edit.metric.lineTension.straight'), value: 0.0 },
        { text: this.$t('edit.metric.lineTension.slight'), value: 0.2 },
        { text: this.$t('edit.metric.lineTension.medium'), value: 0.4 },
        { text: this.$t('edit.metric.lineTension.curvy'), value: 0.6 },
      ],

      lineStyleOptions: [
        { value: '', text: this.$t('edit.metric.lineStyle.default') },
        { value: 'smooth', text: this.$t('edit.metric.lineStyle.smooth') },
        { value: 'step', text: this.$t('edit.metric.lineStyle.step') },
      ],

      scatterSymbolOptions: [
        { value: 'circle', text: this.$t('edit.metric.symbol.circle') },
        { value: 'triangle', text: this.$t('edit.metric.symbol.triangle') },
        { value: 'diamond', text: this.$t('edit.metric.symbol.diamond') },
        { value: 'pin', text: this.$t('edit.metric.symbol.pin') },
        { value: 'arrow', text: this.$t('edit.metric.symbol.arrow') },
        { value: 'rect', text: this.$t('edit.metric.symbol.rect') },
        { value: 'roundRect', text: this.$t('edit.metric.symbol.roundRect') },
      ],
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    hasRelativeDisplay: compose.chartUtil.hasRelativeDisplay,

    getLineStyle (metric) {
      if (metric.smooth) return 'smooth'
      else if (metric.step) return 'step'
      return ''
    },

    setLineStyle (style, metric) {
      this.$set(metric, 'smooth', style === 'smooth')
      this.$set(metric, 'step', style === 'step')
    },

    chartTypeChanged (metric) {
      metric.relativeValue = false
    },

    setDefaultValues () {
      this.chartTypes = []
      this.legendPositions = []
      this.axisLabelPositions = []
      this.tensionSteps = []
      this.lineStyleOptions = []
      this.scatterSymbolOptions = []
    },
  },
}
</script>
<style lang="scss" scoped>
.color-picker {
  max-width: 50px;
}
</style>
