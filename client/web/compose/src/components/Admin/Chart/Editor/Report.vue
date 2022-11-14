<template>
  <div>
    <fieldset v-if="modules">
      <b-form-group>
        <b-form-select
          v-model="moduleID"
          :options="modules"
          text-field="name"
          class="mt-1"
          value-field="moduleID"
        >
          <template slot="first">
            <option
              :value="null"
              disabled
            >
              {{ $t('edit.modulePick') }}
            </option>
          </template>
        </b-form-select>
      </b-form-group>
    </fieldset>

    <div
      v-if="!!module"
      class="mt-1"
    >
      <div class="mb-2">
        <h5 class="mb-3">
          {{ $t('edit.filter.label') }}
        </h5>
        <b-form-group>
          <b-form-select
            v-model="report.filter"
            :disabled="customFilter"
            :options="predefinedFilters"
          >
            <template slot="first">
              <option :value="null">
                {{ $t('edit.filter.noFilter') }}
              </option>
            </template>
          </b-form-select>
          <b-form-checkbox v-model="customFilter">
            {{ $t('edit.filter.customize') }}
          </b-form-checkbox>
          <b-form-textarea
            v-if="customFilter"
            v-model="report.filter"
            placeholder="a = 1 AND b > 2"
          />
        </b-form-group>
      </div>
    </div>

    <div v-if="!!module">
      <div class="px-3 py-2 mb-2">
        <fieldset
          v-for="(d,i) in dimensions"
          :key="'d'+i"
        >
          <h5 class="mb-3">
            {{ $t('edit.dimension.label') }}
          </h5>
          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            :label="$t('edit.dimension.fieldLabel')"
          >
            <b-form-select
              v-model="d.field"
              :options="dimensionFields"
              text-field="name"
              value-field="name"
              @change="onDimFieldChange($event, d)"
            >
              <template slot="first">
                <option
                  disabled
                  :value="undefined"
                >
                  {{ $t('edit.dimension.fieldPlaceholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>
          <b-form-group
            horizontal
            :label-cols="2"
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
                  disabled
                  :value="undefined"
                >
                  {{ $t('edit.dimension.function.placeholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>
          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
          >
            <b-form-checkbox v-model="d.skipMissing">
              {{ $t('edit.dimension.skipMissingValues') }}
            </b-form-checkbox>
          </b-form-group>
          <b-form-group
            v-if="!d.skipMissing"
            horizontal
            :label-cols="2"
            breakpoint="md"
            :label="$t('edit.dimension.defaultValueLabel')"
            :description="$t('edit.dimension.defaultValueFootnote')"
          >
            <b-form-input
              v-model="d.default"
              :type="defaultValueInputType(d)"
            />
          </b-form-group>
          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            label=""
          >
            <b-form-checkbox
              v-model="d.autoSkip"
              :value="true"
              :unchecked-value="false"
            >
              {{ $t('edit.dimension.calculateLabelCount') }}
            </b-form-checkbox><br>
          </b-form-group>

          <!--<b-form-group horizontal v-if="d.field && isTemporalField(d.field)">-->
          <!--<b-form-input type="date" v-model="d.conditions.min"></b-form-input>-->
          <!--<b-form-input type="date" v-model="d.conditions.max"></b-form-input>-->
          <!--</b-form-group>-->
        </fieldset>
      </div>
      <draggable
        class="metrics px-3 py-2"
        :list.sync="metrics"
        :options="{ group: 'metrics_'+moduleID, sort: true }"
      >
        <fieldset
          v-for="(m,i) in metrics"
          :key="'m'+i"
          class="main-fieldset"
        >
          <font-awesome-icon
            v-if="metrics.length>1"
            class="align-baseline text-secondary mr-2"
            :icon="['fas', 'grip-vertical']"
          />
          <h5 class="mb-3 d-inline-block">
            {{ $t('edit.metric.label') }}
          </h5>
          <b-button
            v-if="metrics.length>1"
            variant="link"
            class="text-danger align-baseline"
            @click.prevent="removeMetric(i)"
          >
            <font-awesome-icon :icon="['far', 'trash-alt']" />
          </b-button>
          <b-form
            horizontal
            class="w-25 d-inline-block float-right"
          >
            <b-form-input
              v-model="m.backgroundColor"
              type="color"
            />
          </b-form>

          <b-form-group
            horizontal
            :label-cols="2"
            class="mt-1"
            breakpoint="md"
            :label="$t('edit.metric.labelLabel')"
          >
            <b-form-input
              v-model="m.label"
              :placeholder="$t('edit.metric.labelPlaceholder')"
            />
          </b-form-group>
          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            :label="$t('edit.metric.fieldLabel')"
          >
            <b-form-select
              v-model="m.field"
              :options="metricFields"
              text-field="name"
              value-field="name"
            >
              <template slot="first">
                <option
                  disabled
                  :value="undefined"
                >
                  {{ $t('edit.metric.fieldPlaceholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <b-form-group
            horizontal
            :label-cols="2"
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
                  disabled
                  :value="undefined"
                >
                  {{ $t('edit.metric.function.placeholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>

          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            :label="$t('edit.metric.fx.label')"
            :description="$t('edit.metric.fx.description')"
          >
            <b-form-textarea
              v-model="m.fx"
              placeholder="n"
            />
          </b-form-group>

          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            :label="$t('edit.metric.output.label')"
          >
            <b-form-select
              v-model="m.type"
              :disabled="!m.field"
              :options="chartTypes"
            >
              <template slot="first">
                <option
                  disabled
                  :value="undefined"
                >
                  {{ $t('edit.metric.output.placeholder') }}
                </option>
              </template>
            </b-form-select>
          </b-form-group>
          <b-form-group
            horizontal
            :label-cols="2"
            breakpoint="md"
            label=""
          >
            <template v-if="hasRelativeDisplay(m)">
              <b-form-checkbox
                v-model="m.relativeValue"
                :value="true"
                :unchecked-value="false"
              >
                {{ $t('edit.metric.relative') }}
              </b-form-checkbox>

              <template v-if="m.relativeValue">
                <b-form-group
                  horizontal
                  breakpoint="md"
                  :label="$t('edit.metric.relativePrecision')"
                >
                  <b-form-input
                    v-model="m.relativePrecision"
                    type="number"
                    placeholder="2"
                  />
                </b-form-group>
              </template>
            </template>

            <template v-else>
              <b-form-checkbox
                v-model="m.axisType"
                value="logarithmic"
                unchecked-value="linear"
              >
                {{ $t('edit.metric.logarithmicScale') }}
              </b-form-checkbox>
              <b-form-checkbox
                v-model="m.axisPosition"
                value="right"
                unchecked-value="left"
              >
                {{ $t('edit.metric.axisOnRight') }}
              </b-form-checkbox>
              <b-form-checkbox
                v-model="m.beginAtZero"
                :value="true"
                :unchecked-value="false"
                checked
              >
                {{ $t('edit.metric.axisScaleFromZero') }}
              </b-form-checkbox>
            </template>

            <b-form-checkbox
              v-show="m.type === 'line'"
              v-model="m.fill"
              :value="true"
              :unchecked-value="false"
            >
              {{ $t('edit.metric.fillArea') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="m.fixTooltips"
              :value="true"
              :unchecked-value="false"
            >
              {{ $t('edit.metric.fixTooltips') }}
            </b-form-checkbox>
          </b-form-group>
        </fieldset>
      </draggable>
    </div>
    <b-button
      variant="link"
      class="float-right"
      @click.prevent="metrics.push({})"
    >
      + {{ $t('edit.metric.add') }}
    </b-button>
  </div>
</template>
<script>
import draggable from 'vuedraggable'
import { compose } from '@cortezaproject/corteza-js'

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

  name: 'Report',

  components: {
    draggable,
  },

  props: {
    report: {
      type: [Object, undefined],
      required: false,
      default: () => ({}),
    },

    modules: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      customFilter: false,

      metricAggregates: aggregateFunctions.map(af => ({ ...af, text: this.$t(`edit.metric.function.${af.text}`) })),
      dimensionModifiers: compose.chartUtil.dimensionFunctions.map(df => ({ ...df, text: this.$t(`edit.dimension.function.${df.text}`) })),
      predefinedFilters: compose.chartUtil.predefinedFilters.map(pf => ({ ...pf, text: this.$t(`edit.filter.${pf.text}`) })),
      chartTypes: Object.values(compose.chartUtil.ChartType).map(value => ({ value, text: this.$t(`edit.metric.output.${value}`) })),
    }
  },

  computed: {
    defaultValueInputType () {
      return ({ field }) => field === 'created_at' || (this.module.fields.filter(f => f.name === field)[0] || {}).kind === 'DateTime' ? 'date' : 'text'
    },

    module () {
      return this.modules.find(m => m.moduleID === this.moduleID)
    },

    metricFields () {
      return [{ name: 'count' }, ...this.module.fields.filter(f => f.kind === 'Number')]
    },

    dimensionFields () {
      return [
        { name: 'created_at' },
        ...this.module.fields.map(f => {
          const { name, label, kind, options } = f
          let disabled = true
          switch (kind) {
            case 'DateTime':
            case 'Select':
            case 'Number':
            case 'Bool':
              disabled = false
              break
            case 'String':
              disabled = options.useRichTextEditor || options.multiLine
              break
          }

          return { name, label, disabled }
        })]
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
        // this.report.dimensions = v
        this.$emit('update:report', { ...this.report, dimensions: v })
      },
    },
  },

  watch: {
    'report.filter' (v) {
      this.customFilter = !compose.chartUtil.predefinedFilters.includes(v)
    },
  },

  methods: {
    hasRelativeDisplay: compose.chartUtil.hasRelativeDisplay,

    onDimFieldChange (f, d) {
      if (!this.isTemporalField(f)) {
        this.$set(d, 'modifier', this.dimensionModifiers[0].value)
      }
    },

    removeMetric (i) {
      this.metrics.splice(i, 1)
    },

    isTemporalField (name) {
      if (name === 'created_at') {
        return true
      }

      return !!this.module.fields.find(f => f.name === name && f.kind === 'DateTime')
    },
  },
}
</script>
