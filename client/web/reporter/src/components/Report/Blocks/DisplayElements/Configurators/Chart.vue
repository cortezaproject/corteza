<template>
  <div
    v-if="options"
  >
    <div
      class="mb-3"
    >
      <h5 class="text-primary mb-2">
        {{ $t('display-element:chart.configurator.general') }}
      </h5>

      <b-row>
        <b-col>
          <b-form-group
            :label="$t('display-element:chart.configurator.type')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="options.type"
              :options="chartTypes"
              @change="typeChanged"
            />
          </b-form-group>
        </b-col>
        <b-col>
          <b-form-group
            :label="$t('display-element:chart.configurator.chart-title')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="options.title"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col>
          <b-form-group
            :label="$t('display-element:chart.configurator.color-scheme')"
            label-class="text-primary"
            class="mb-0"
          >
            <vue-select
              v-model="options.colorScheme"
              :options="colorSchemes"
              :reduce="cs => cs.value"
              label="label"
              option-text="label"
              option-value="value"
              clearable
              class="h-100 w-100"
            >
              <template #option="option">
                <div
                  v-for="(color, index) in option.colors"
                  :key="`${option.value}-${index}`"
                  :style="`background: ${color};`"
                  class="d-inline-block color-box mr-1"
                />
              </template>
            </vue-select>
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
          align-self="center"
        >
          <b-form-checkbox
            v-model="options.showTooltips"
            class="pb-2"
          >
            {{ $t('display-element:chart.configurator.show.tooltips') }}
          </b-form-checkbox>
          <b-form-checkbox
            v-model="options.showLegend"
            class="pb-2"
          >
            {{ $t('display-element:chart.configurator.show.legend') }}
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
        {{ $t('display-element:chart.configurator.data') }}
      </h5>

      <b-form-group
        v-if="options.labelColumn !== undefined"
        :label="$t('display-element:chart.configurator.label-column')"
        label-class="text-primary"
      >
        <b-form-select
          v-model="options.labelColumn"
          :options="labelColumns"
          text-field="label"
          value-field="name"
        >
          <template #first>
            <b-form-select-option
              value=""
            >
              {{ $t('display-element:chart.configurator.none') }}
            </b-form-select-option>
          </template>
        </b-form-select>
      </b-form-group>

      <b-form-group
        v-if="options.dataColumns && columns.length"
        :label="$t('display-element:chart.configurator.data-columns')"
        label-class="text-primary"
      >
        <column-picker
          :all-items="dataColumns"
          :selected-items.sync="options.dataColumns"
          class="d-flex flex-column"
        />
      </b-form-group>

      <div
        v-if="['bar', 'line'].includes(options.type)"
      >
        <hr>

        <div
          class="mb-3"
        >
          <h5 class="text-primary mb-2">
            {{ $t('display-element:chart.configurator.x-axis.name') }}
          </h5>
          <b-row>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.x-axis.label')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.xAxis.label"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.x-axis.type')"
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
                      {{ $t('display-element:chart.configurator.default') }}
                    </b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.default-value')"
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
                {{ $t('display-element:chart.configurator.skip-missing-values') }}
              </b-form-checkbox>
            </b-col>
            <b-col>
              <b-form-group
                v-if="options.xAxis.type === 'time'"
                :label="$t('display-element:chart.configurator.time.unit.label')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.xAxis.unit"
                  :options="timeUnits"
                >
                  <template #first>
                    <b-form-select-option
                      :value="undefined"
                    >
                      {{ $t('display-element:chart.configurator.default') }}
                    </b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>
            </b-col>
          </b-row>
        </div>

        <hr>

        <div>
          <h5 class="text-primary mb-2">
            {{ $t('display-element:chart.configurator.y-axis.name') }}
          </h5>
          <b-row>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.y-axis.label')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.label"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.step-size')"
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
              <b-form-group
                :label="$t('display-element:chart.configurator.value.min')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.min"
                  type="number"
                />
              </b-form-group>
            </b-col>
            <b-col>
              <b-form-group
                :label="$t('display-element:chart.configurator.value.max')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="options.yAxis.max"
                  type="number"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group>
                <b-form-checkbox
                  v-model="options.yAxis.beginAtZero"
                >
                  {{ $t('display-element:chart.configurator.begin-axis-at-zero') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.yAxis.type"
                  value="logarithmic"
                  unchecked-value="linear"
                >
                  {{ $t('display-element:chart.configurator.logarithmic-scale') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.yAxis.position"
                  value="right"
                  unchecked-value="left"
                >
                  {{ $t('display-element:chart.configurator.place-axis-on-right-side') }}
                </b-form-checkbox>
              </b-form-group>
            </b-col>
          </b-row>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import base from './base'
import ColumnPicker from 'corteza-webapp-reporter/src/components/Common/ColumnPicker'
import colorschemes from 'chartjs-plugin-colorschemes/src/colorschemes'
import VueSelect from 'vue-select'
import { reporter } from '@cortezaproject/corteza-js'

export default {
  components: {
    VueSelect,
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
    }
  },

  computed: {
    mergedColumns () {
      return [].concat.apply([], this.columns)
    },

    chartTypes () {
      const types = [
        { value: 'bar', text: this.$t('display-element:chart.configurator.types.bar') },
        { value: 'line', text: this.$t('display-element:chart.configurator.types.line') },
        { value: 'pie', text: this.$t('display-element:chart.configurator.types.pie') },
        { value: 'doughnut', text: this.$t('display-element:chart.configurator.types.doughnut') },
      ]

      if (this.datasource && this.datasource.step.group) {
        types.push({ value: 'funnel', text: this.$t('display-element:chart.configurator.types.funnel') })
      }

      return types
    },

    AxisTypes () {
      return [
        { value: 'time', text: this.$t('display-element:chart.configurator.time.label') },
        { value: 'category', text: 'Category' },
      ]
    },

    timeUnits () {
      return [
        { value: 'day', text: this.$t('display-element:chart.configurator.time.unit.types.date') },
        { value: 'week', text: this.$t('display-element:chart.configurator.time.unit.types.week') },
        { value: 'month', text: this.$t('display-element:chart.configurator.time.unit.types.month') },
        { value: 'quarter', text: this.$t('display-element:chart.configurator.time.unit.types.quarter') },
        { value: 'year', text: this.$t('display-element:chart.configurator.time.unit.types.year') },
      ]
    },

    labelColumns () {
      const columns = this.columns.length ? this.columns[0] : []
      return [
        ...columns.filter(({ kind }) => this.allowedLabelKinds.includes(kind)),
      ].sort((a, b) => a.label.localeCompare(b.label))
    },

    dataColumns () {
      return [
        ...this.mergedColumns.filter(({ kind }) => kind === 'Number'),
      ].sort((a, b) => a.label.localeCompare(b.label))
    },

    currentColorScheme () {
      return this.colorSchemes.find(({ value }) => value === this.options.colorScheme)
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
          colors: [...colorschemes[g][sc]].reverse(),
          value: `${g}.${sc}`,
        })
      }
    }

    this.colorSchemes = rr
  },

  methods: {
    setColorscheme (colorscheme) {
      this.options.colorScheme = (colorscheme || {}).value || ''
    },

    typeChanged () {
      this.options = reporter.ChartOptionsMaker(this.options) || {}
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
