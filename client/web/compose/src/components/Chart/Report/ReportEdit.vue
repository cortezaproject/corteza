<template>
  <div>
    <!-- Configure source module -->
    <b-form-group
      :label="$t('edit.module.label')"
    >
      <b-form-select
        v-model="moduleID"
        :options="modules"
        text-field="name"
        class="mt-1"
        value-field="moduleID"
      >
        <template slot="first">
          <option
            :value="undefined"
          >
            {{ $t('edit.module.placeholder') }}
          </option>
        </template>
      </b-form-select>
    </b-form-group>

    <!-- Configure report filters -->
    <div
      v-if="!!module"
      class="mt-1"
    >
      <b-form-group
        :label="$t('edit.filter.label')"
      >
        <b-form-select
          v-model="report.filter"
          :disabled="customFilter"
          :options="predefinedFilters"
        >
          <template slot="first">
            <option value="">
              {{ $t('edit.filter.noFilter') }}
            </option>
          </template>
        </b-form-select>

        <b-form-checkbox
          v-model="customFilter"
          class="mt-1"
        >
          {{ $t('edit.filter.customize') }}
        </b-form-checkbox>

        <b-form-textarea
          v-if="customFilter"
          v-model="report.filter"
          placeholder="a = 1 AND b > 2"
        />
      </b-form-group>
    </div>

    <slot
      name="y-axis"
      :report="editReport"
    />

    <!-- Configure report dimensions -->
    <div v-if="!!module">
      <hr>

      <fieldset
        v-for="(d, i) in dimensions"
        :key="i"
      >
        <h4 class="mb-3">
          {{ $t('edit.dimension.label') }}
        </h4>

        <template v-if="usesDimensionsField">
          <b-form-group
            horizontal
            :label-cols="3"
            breakpoint="md"
            :label="$t('edit.dimension.fieldLabel')"
          >
            <b-form-select
              v-model="d.field"
              :options="dimensionFields"
              @change="onDimFieldChange($event, d)"
            >
              <template slot="first">
                <option
                  :value="undefined"
                >
                  {{ $t('edit.dimension.fieldPlaceholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <b-form-group
            horizontal
            :label-cols="3"
            breakpoint="md"
            :label="$t('edit.dimension.function.label')"
          >
            <b-form-select
              v-model="d.modifier"
              :disabled="!d.field || !isTemporalField(d.field)"
              :options="dimensionModifiers"
            >
              <template slot="first">
                <option
                  :value="undefined"
                >
                  {{ $t('edit.dimension.function.placeholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <template v-if="!unSkippable">
            <b-form-group
              horizontal
              :label-cols="3"
              breakpoint="md"
            >
              <b-form-checkbox
                v-model="d.skipMissing"
              >
                {{ $t('edit.dimension.skipMissingValues') }}
              </b-form-checkbox>
            </b-form-group>

            <b-form-group
              v-if="!d.skipMissing"
              horizontal
              :label-cols="3"
              breakpoint="md"
              :label="$t('edit.dimension.defaultValueLabel')"
              :description="$t('edit.dimension.defaultValueFootnote')"
            >
              <b-form-input
                v-model="d.default"
                :type="defaultValueInputType(d)"
              />
            </b-form-group>
          </template>
        </template>

        <slot
          name="dimension-options"
          :index="i"
          :dimension="d"
          :field="getField(d)"
        />
      </fieldset>

      <hr>

      <h4 class="d-flex align-items-center mb-3">
        {{ $t('edit.metric.title') }}
        <b-button
          v-if="canAddMetric"
          variant="link"
          class="text-decoration-none"
          @click.prevent="addMetric"
        >
          + {{ $t('edit.metric.add') }}
        </b-button>
      </h4>

      <!-- Configure report metrics -->
      <draggable
        class="metrics mb-3"
        :list.sync="metrics"
        handle=".metric-handle"
        :options="{ group: `metrics_${moduleID}` }"
      >
        <fieldset
          v-for="(m,i) in metrics"
          :key="i"
          class="main-fieldset mb-3"
        >
          <h5
            v-if="metrics.length > 1"
            class="d-flex align-items-center mb-3"
          >
            <font-awesome-icon
              class="grab metric-handle align-baseline text-secondary mr-2"
              :icon="['fas', 'grip-vertical']"
            />
            {{ $t('edit.metric.label') }} {{ i + 1 }}
            <b-button
              variant="link"
              class="text-danger align-baseline"
              @click.prevent="removeMetric(i)"
            >
              <font-awesome-icon :icon="['far', 'trash-alt']" />
            </b-button>
          </h5>

          <b-form-group
            horizontal
            :label-cols="3"
            breakpoint="md"
            :label="$t('edit.metric.fieldLabel')"
          >
            <b-form-select
              v-model="m.field"
              :options="metricFields"
              @change="onMetricFieldChange($event, m)"
            >
              <template slot="first">
                <option
                  :value="undefined"
                >
                  {{ $t('edit.metric.fieldPlaceholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <b-form-group
            horizontal
            :label-cols="3"
            breakpoint="md"
            :label="$t('edit.metric.function.label')"
          >
            <b-form-select
              v-model="m.aggregate"
              :disabled="!m.field || m.field === 'count'"
              :options="metricAggregates"
            >
              <template slot="first">
                <option
                  :value="undefined"
                >
                  {{ $t('edit.metric.function.placeholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <slot
            name="metric-options"
            :metric="m"
          />
        </fieldset>
      </draggable>
    </div>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import { compose } from '@cortezaproject/corteza-js'
import base from './base'

const aggregateFunctions = [
  {
    value: 'SUM',
    text: 'sum',
  },
  {
    value: 'MAX',
    text: 'max',
  },
  {
    value: 'MIN',
    text: 'min',
  },
  {
    value: 'AVG',
    text: 'avg',
  },
  {
    value: 'STD',
    text: 'std',
  },
]

export default {
  i18nOptions: {
    namespaces: 'chart',
  },

  name: 'ReportEdit',

  components: {
    draggable,
  },

  extends: base,

  data () {
    return {
      customFilter: false,

      metricAggregates: aggregateFunctions.map(af => ({ ...af, text: this.$t(`edit.metric.function.${af.text}`) })),
      dimensionModifiers: compose.chartUtil.dimensionFunctions.map(df => ({ ...df, text: this.$t(`edit.dimension.function.${df.text}`) })),
      predefinedFilters: compose.chartUtil.predefinedFilters.map(pf => ({ ...pf, text: this.$t(`edit.filter.${pf.text}`) })),
    }
  },

  computed: {
    defaultValueInputType () {
      return ({ field }) => (this.module.fields.filter(f => f.name === field)[0] || {}).kind === 'DateTime' ? 'date' : 'text'
    },

    canAddMetric () {
      return (this.supportedMetrics < 0 || this.metrics.length < this.supportedMetrics) && this.moduleID
    },

    module () {
      return this.modules.find(m => m.moduleID === this.moduleID)
    },

    metricFields () {
      return [
        { value: 'count', text: 'Count' },
        ...this.module.fields.filter(f => f.kind === 'Number')
          .map(({ name }) => ({ value: name, text: name }))
          .sort((a, b) => a.text.localeCompare(b.text)),
      ]
    },

    dimensionFields () {
      return [
        ...[...this.module.fields].sort((a, b) => a.label.localeCompare(b.text)),
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(({ kind, options = {} }) => {
        return this.dimensionFieldKind.includes(kind) && !(options.useRichTextEditor || options.multiLine)
      }).map(({ name, label, kind }) => {
        return { value: name, text: `${label} (${kind})`, kind }
      })
    },

    moduleID: {
      get () {
        return this.report.moduleID
      },

      set (v) {
        this.report.moduleID = v
        this.$emit('update:report', { ...this.report, moduleID: v })
      },
    },

    colorScheme: {
      get () {
        return this.report.colorScheme
      },

      set (v) {
        this.report.colorScheme = v
        this.$emit('update:report', { ...this.report, colorScheme: v })
      },
    },

    metrics: {
      get () {
        return this.report.metrics
      },

      set (v) {
        this.report.metrics = v
        this.$emit('update:report', { ...this.report, metrics: v })
      },
    },

    dimensions: {
      get () {
        return this.report.dimensions
      },

      set (v) {
        this.$emit('update:report', { ...this.report, dimensions: v })
      },
    },
  },

  watch: {
    'report.filter': {
      handler: function (v) {
        // !! is required, since :disabled="..." marks the field as disabled if '' is provided
        this.customFilter = (!!v && !!compose.chartUtil.predefinedFilters.find(({ value }) => value === v)) ||
          (!!v)
      },
      immediate: true,
    },
  },

  methods: {
    hasRelativeDisplay: compose.chartUtil.hasRelativeDisplay,

    getField ({ field }) {
      if (!field || !this.module) {
        return undefined
      }

      return this.module.fields.find(({ name }) => name === field)
    },

    addMetric () {
      this.metrics = this.metrics.concat([{}])
    },

    onDimFieldChange (f, d) {
      if (!this.isTemporalField(f)) {
        this.$set(d, 'modifier', this.dimensionModifiers[0].value)
      }
    },

    onMetricFieldChange (field, m) {
      if (field === 'count') {
        this.$set(m, 'aggregate', undefined)
      }
    },

    removeMetric (i) {
      this.metrics.splice(i, 1)
    },

    isTemporalField (name) {
      return this.dimensionFields.some(f => f.value === name && f.kind === 'DateTime')
    },
  },
}
</script>
