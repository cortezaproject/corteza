<template>
  <div>
    <div
      class="mb-3"
    >
      <h5 class="text-primary mb-2">
        {{ $t('chart.configurator.general') }}
      </h5>

      <b-row
        align-v="stretch"
      >
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('chart.configurator.type')"
            label-class="text-primary"
          >
            <b-form-select
              :value="displayElementOptions.type"
              :options="chartTypes"
              @change="typeChanged"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('chart.configurator.chart-title')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="options.title"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('chart.configurator.color-scheme')"
            label-class="text-primary"
            class="mb-0"
          >
            <c-input-select
              v-model="options.colorScheme"
              :options="colorSchemes"
              :get-option-key="getOptionKey"
              :reduce="cs => cs.value"
              :placeholder="$t('general:label.default')"
            >
              <template #option="option">
                <p
                  class="mb-1"
                >
                  {{ option.label }}
                </p>

                <div
                  v-for="(color, index) in option.colors"
                  :key="`${option.value}-${index}`"
                  :style="`background: ${color};`"
                  class="d-inline-block color-box mr-1 mb-1"
                />
              </template>
            </c-input-select>

            <template
              v-if="currentColorScheme"
            >
              <div
                v-for="(color, index) in currentColorScheme.colors"
                :key="`${currentColorScheme.value}-${index}`"
                :style="`background: ${color};`"
                class="d-inline-block color-box mr-1"
              />
            </template>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
          class="d-flex flex-column justify-content-center"
        >
          <b-form-checkbox
            v-model="options.noAnimation"
            :value="false || undefined"
            :unchecked-value="true"
            switch
            class="mt-3 pt-2"
          >
            {{ $t('chart.configurator.animation.enabled') }}
          </b-form-checkbox>
        </b-col>
      </b-row>
      <hr>
    </div>

    <div
      v-if="options.source"
      class="mb-3"
    >
      <h5 class="text-primary mb-2">
        {{ $t('chart.configurator.data') }}
      </h5>

      <b-form-group
        v-if="options.labelColumn !== undefined"
        :label="$t('chart.configurator.label-column')"
        label-class="text-primary"
      >
        <column-selector
          v-model="options.labelColumn"
          :columns="labelColumns"
          style="min-width: 100% !important;"
        />
      </b-form-group>

      <b-form-group
        v-if="options.dataColumns && columns.length"
        :label="$t('chart.configurator.data-columns')"
        label-class="text-primary"
      >
        <column-picker
          :all-items="dataColumns"
          :selected-items="options.dataColumns"
          class="d-flex flex-column"
          @update:selected-items="updateDataColumns"
        />
      </b-form-group>

      <div
        v-if="['bar', 'line'].includes(options.type)"
      >
        <hr>

        <h5 class="text-primary mb-2">
          {{ $t('chart.configurator.metrics.label') }}
        </h5>

        <b-form-group
          v-if="options.dataColumns && options.dataColumns.length"
          label-class="text-primary"
        >
          <c-form-table-wrapper
            v-for="(column, index) in metricColumns"
            :key="index"
            hide-add-button
            class="mb-3"
          >
            <h6 class="mb-3">
              {{ column.label }}
            </h6>

            <b-row>
              <b-col
                cols="12"
                md="6"
              >
                <b-form-group
                  :label="$t('chart.configurator.metrics.label-field.label')"
                  label-class="text-primary"
                  class="mb-0"
                >
                  <b-form-input
                    v-model="column.label"
                  />
                </b-form-group>
              </b-col>

              <b-col
                cols="12"
                md="6"
              >
                <b-form-group
                  :label="$t('chart.configurator.metrics.stack.label')"
                  :description="$t('chart.configurator.metrics.stack.description')"
                  label-class="text-primary"
                  class="mb-0"
                >
                  <b-form-input
                    v-model="column.stack"
                  />
                </b-form-group>
              </b-col>
            </b-row>
          </c-form-table-wrapper>
        </b-form-group>
      </div>

      <div
        v-if="['bar', 'line'].includes(options.type)"
      >
        <hr>

        <div
          class="mb-3"
        >
          <h5 class="text-primary mb-2">
            {{ $t('chart.configurator.x-axis.name') }}
          </h5>
          <b-row>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.x-axis.label')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.xAxis.label"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.x-axis.type')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.xAxis.type"
                  :options="AxisTypes"
                >
                  <template #first>
                    <b-form-select-option
                      value=""
                    >
                      {{ $t('chart.configurator.default') }}
                    </b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.default-value')"
                label-class="text-primary"
                class="mb-1"
              >
                <b-form-input
                  v-model="options.xAxis.defaultValue"
                  :disabled="options.xAxis.skipMissing"
                  :type="options.xAxis.type === 'time' ? 'date' : 'text'"
                />
              </b-form-group>

              <b-form-checkbox
                v-model="options.xAxis.skipMissing"
                class="mb-3"
              >
                {{ $t('chart.configurator.skip-missing-values') }}
              </b-form-checkbox>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.x-axis.labelRotation.label')"
                label-class="text-primary"
                class="mb-1"
              >
                <b-input
                  v-model="options.xAxis.labelRotation"
                  type="number"
                  number
                />
              </b-form-group>
            </b-col>
          </b-row>
        </div>

        <hr>

        <div>
          <h5 class="text-primary mb-2">
            {{ $t('chart.configurator.y-axis.name') }}
          </h5>
          <b-row>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.y-axis.label')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.label"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.y-axis.labelPosition.label')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.yAxis.labelPosition"
                  :options="axisLabelPositions"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.y-axis.labelRotation.label')"
                label-class="text-primary"
                class="mb-1"
              >
                <b-input
                  v-model="options.yAxis.labelRotation"
                  type="number"
                  number
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.value.min')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.min"
                  type="number"
                  number
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.value.max')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.max"
                  type="number"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('chart.configurator.step-size')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.stepSize"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group>
                <b-form-checkbox
                  v-model="options.yAxis.type"
                  value="logarithmic"
                  unchecked-value="linear"
                >
                  {{ $t('chart.configurator.logarithmic-scale') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.yAxis.beginAtZero"
                >
                  {{ $t('chart.configurator.begin-axis-at-zero') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.yAxis.position"
                  value="right"
                  unchecked-value="left"
                >
                  {{ $t('chart.configurator.place-axis-on-right-side') }}
                </b-form-checkbox>
              </b-form-group>
            </b-col>
          </b-row>
        </div>
      </div>

      <hr>

      <div>
        <h5 class="text-primary mb-2">
          {{ $t('chart.configurator.legend.name') }}
        </h5>

        <b-form-checkbox
          v-model="options.legend.hide"
          class="mb-3"
        >
          {{ $t('chart.configurator.legend.hide') }}
        </b-form-checkbox>

        <b-row v-if="!options.legend.hide">
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('chart.configurator.legend.align.label')"
              label-class="text-primary"
              class="mb-1"
            >
              <b-form-select
                v-model="options.legend.align"
                :options="alignments"
                :disabled="!options.legend.position.default"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('chart.configurator.legend.orientation.label')"
              label-class="text-primary"
              class="mb-1"
            >
              <b-form-select
                v-model="options.legend.orientation"
                :options="orientations"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-checkbox
              v-model="options.legend.position.default"
              :class="{ 'mb-3': !options.legend.position.default }"
            >
              {{ $t('chart.configurator.legend.position.default') }}
            </b-form-checkbox>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-checkbox
              v-model="options.legend.scrollable"
              :disabled="options.legend.orientation !== 'horizontal'"
            >
              {{ $t('chart.configurator.legend.scrollable') }}
            </b-form-checkbox>
          </b-col>

          <template
            v-if="!options.legend.position.default"
          >
            <b-col
              cols="12"
              lg="6"
              xl="3"
            >
              <b-form-group
                :label="$t('chart.configurator.legend.position.top')"
                label-class="text-primary"
              >
                <b-input
                  v-model="options.legend.position.top"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
              xl="3"
            >
              <b-form-group
                :label="$t('chart.configurator.legend.position.right')"
                label-class="text-primary"
              >
                <b-input
                  v-model="options.legend.position.right"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
              xl="3"
            >
              <b-form-group
                :label="$t('chart.configurator.legend.position.bottom')"
                label-class="text-primary"
              >
                <b-input
                  v-model="options.legend.position.bottom"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
              xl="3"
            >
              <b-form-group
                :label="$t('chart.configurator.legend.position.left')"
                label-class="text-primary"
              >
                <b-input
                  v-model="options.legend.position.left"
                />
              </b-form-group>
            </b-col>

            <b-col>
              <small>{{ $t('chart.configurator.position-description') }}</small>
            </b-col>
          </template>
        </b-row>
      </div>

      <hr>

      <div>
        <h5 class="text-primary mb-2">
          {{ $t('chart.configurator.tooltips.name') }}
        </h5>

        <b-row>
          <b-col>
            <b-form-checkbox
              v-model="options.tooltips.showAlways"
            >
              {{ $t('chart.configurator.tooltips.show.always') }}
            </b-form-checkbox>
          </b-col>
        </b-row>
      </div>

      <hr>

      <div class="mb-2">
        <h5 class="text-primary mb-2">
          {{ $t('chart.configurator.offset.name') }}
        </h5>

        <b-form-checkbox
          v-model="options.offset.default"
          class="mb-3"
        >
          {{ $t('chart.configurator.offset.default') }}
        </b-form-checkbox>

        <b-row v-if="!options.offset.default">
          <b-col
            cols="12"
            lg="6"
            xl="3"
          >
            <b-form-group
              :label="$t('chart.configurator.offset.top')"
              label-class="text-primary"
            >
              <b-input
                v-model="options.offset.top"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
            xl="3"
          >
            <b-form-group
              :label="$t('chart.configurator.offset.right')"
              label-class="text-primary"
            >
              <b-input
                v-model="options.offset.right"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
            xl="3"
          >
            <b-form-group
              :label="$t('chart.configurator.offset.bottom')"
              label-class="text-primary"
            >
              <b-input
                v-model="options.offset.bottom"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
            xl="3"
          >
            <b-form-group
              :label="$t('chart.configurator.offset.left')"
              label-class="text-primary"
            >
              <b-input
                v-model="options.offset.left"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
          >
            <small>{{ $t('chart.configurator.position-description') }}</small>
          </b-col>
        </b-row>
      </div>
    </div>
  </div>
</template>

<script>
import base from './base'
import ColumnSelector from 'corteza-webapp-reporter/src/components/Common/ColumnSelector.vue'
import ColumnPicker from 'corteza-webapp-reporter/src/components/Common/ColumnPicker'
import { reporter, shared } from '@cortezaproject/corteza-js'
const { colorschemes } = shared

export default {
  i18nOptions: {
    namespaces: 'display-element',
  },

  name: 'ChartConfigurator',

  components: {
    ColumnSelector,
    ColumnPicker,
  },

  extends: base,

  data () {
    return {
      colorSchemes: [],

      allowedLabelKinds: [
        'Date',
        'DateTime',
        'Select',
        'Number',
        'Bool',
        'String',
        'Record',
        'User',
      ],

      axisLabelPositions: [
        { value: 'end', text: this.$t('chart.configurator.y-axis.labelPosition.top') },
        { value: 'center', text: this.$t('chart.configurator.y-axis.labelPosition.center') },
        { value: 'start', text: this.$t('chart.configurator.y-axis.labelPosition.bottom') },
      ],

      orientations: [
        { value: 'horizontal', text: this.$t('chart.configurator.legend.orientation.horizontal') },
        { value: 'vertical', text: this.$t('chart.configurator.legend.orientation.vertical') },
      ],

      alignments: [
        { value: 'left', text: this.$t('chart.configurator.legend.align.left') },
        { value: 'center', text: this.$t('chart.configurator.legend.align.center') },
        { value: 'right', text: this.$t('chart.configurator.legend.align.right') },
      ],
    }
  },

  computed: {
    mergedColumns () {
      return [].concat.apply([], this.columns)
    },

    chartTypes () {
      const types = [
        { value: 'bar', text: this.$t('chart.configurator.types.bar') },
        { value: 'line', text: this.$t('chart.configurator.types.line') },
        { value: 'pie', text: this.$t('chart.configurator.types.pie') },
        { value: 'doughnut', text: this.$t('chart.configurator.types.doughnut') },
      ]

      if (this.datasource && this.datasource.step.aggregate) {
        types.push({ value: 'funnel', text: this.$t('chart.configurator.types.funnel') })
      }

      return types
    },

    AxisTypes () {
      return [
        { value: 'time', text: this.$t('chart.configurator.time.label') },
        { value: 'category', text: 'Category' },
      ]
    },

    timeUnits () {
      return [
        { value: 'day', text: this.$t('chart.configurator.time.unit.types.date') },
        { value: 'week', text: this.$t('chart.configurator.time.unit.types.week') },
        { value: 'month', text: this.$t('chart.configurator.time.unit.types.month') },
        { value: 'quarter', text: this.$t('chart.configurator.time.unit.types.quarter') },
        { value: 'year', text: this.$t('chart.configurator.time.unit.types.year') },
      ]
    },

    labelColumns () {
      const columns = this.columns.length ? this.columns[0] : []
      return columns.filter(({ kind }) => this.allowedLabelKinds.includes(kind))
    },

    dataColumns () {
      return [
        ...this.mergedColumns.filter(({ kind }) => kind === 'Number'),
      ].sort((a, b) => a.label.localeCompare(b.label))
    },

    currentColorScheme () {
      return this.colorSchemes.find(({ value }) => value === this.options.colorScheme)
    },

    metricColumns () {
      return this.options.dataColumns
    },
  },

  created () {
    const capitalize = w => `${w[0].toUpperCase()}${w.slice(1)}`
    const splicer = sc => {
      const rr = (/(\D+)(\d+)$/gi).exec(sc)
      return {
        label: rr[1],
        count: rr[2],
      }
    }

    const rr = []
    for (const g in colorschemes) {
      for (const sc in colorschemes[g]) {
        const gn = splicer(sc)

        rr.push({
          label: `${capitalize(g)}: ${capitalize(gn.label)}`,
          colors: [...colorschemes[g][sc]],
          value: `${g}.${sc}`,
        })
      }
    }

    this.colorSchemes = rr
  },

  methods: {
    updateDataColumns (columns) {
      let result = columns.map((c) => this.dataColumns.find(({ name }) => name === c))

      result.forEach((c, idx) => {
        const initialCol = this.options.dataColumns.find(({ name }) => name === c.name)
        if (initialCol) result[idx] = initialCol
      })

      this.options.dataColumns = result
    },

    setColorscheme (colorscheme) {
      this.options.colorScheme = (colorscheme || {}).value || ''
    },

    typeChanged (type) {
      this.options = reporter.ChartOptionsMaker({ ...this.options, type })
    },

    getOptionKey ({ value }) {
      return value
    },
  },
}
</script>

<style lang="scss" scoped>
.color-box {
  width: 28px;
  height: 12px;
}
</style>
