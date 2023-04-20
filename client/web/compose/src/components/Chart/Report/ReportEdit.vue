<template>
  <div>
    <!-- Configure source module -->
    <div
      class="px-3"
    >
      <h5 class="mb-3">
        {{ $t('edit.module.title') }}
      </h5>

      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.module.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="moduleID"
              :options="modules"
              text-field="name"
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
        </b-col>

        <b-col
          v-if="!!module"
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.filter.preset')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="report.filter"
              :options="predefinedFilters"
            >
              <template slot="first">
                <b-form-select-option :value="defaultFilterOption">
                  {{ $t('edit.filter.noFilter') }}
                </b-form-select-option>
              </template>
            </b-form-select>
          </b-form-group>
        </b-col>

        <!-- Configure report filters -->
        <b-col
          v-if="!!module"
          cols="12"
          class="mt-1"
        >
          <b-form-group
            :label="$t('edit.filter.label')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="report.filter"
              :placeholder="$t('edit.filter.placeholder')"
            />

            <i18next
              path="edit.filter.footnote"
              tag="small"
              class="text-muted"
            >
              <code>${recordID}</code>
              <code>${ownerID}</code>
              <code>${userID}</code>
            </i18next>
          </b-form-group>
        </b-col>
      </b-row>
    </div>
    <hr v-if="module">

    <!-- Configure report dimensions -->
    <div
      v-if="!!module"
      class="px-3"
    >
      <fieldset
        v-for="(d, i) in dimensions"
        :key="i"
      >
        <h5 class="mb-3">
          {{ $t('edit.dimension.label') }}
        </h5>

        <template v-if="usesDimensionsField">
          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('edit.dimension.fieldLabel')"
                label-class="text-primary"
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
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('edit.dimension.function.label')"
                label-class="text-primary"
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
            </b-col>
          </b-row>

          <template v-if="!unSkippable">
            <b-row>
              <b-col
                v-if="!d.skipMissing"
                cols="12"
                md="6"
              >
                <b-form-group
                  :label="$t('edit.dimension.defaultValueLabel')"
                  :description="$t('edit.dimension.defaultValueFootnote')"
                  label-class="text-primary"
                >
                  <b-form-input
                    v-model="d.default"
                    :type="defaultValueInputType(d)"
                  />
                </b-form-group>
              </b-col>

              <b-col
                cols="12"
                md="6"
              >
                <b-form-group
                  :label="$t('edit.dimension.options.label')"
                  label-class="text-primary"
                >
                  <b-form-checkbox
                    v-model="d.skipMissing"
                  >
                    {{ $t('edit.dimension.skipMissingValues') }}
                  </b-form-checkbox>

                  <slot
                    name="dimension-options-options"
                    :dimension="d"
                  />
                </b-form-group>
              </b-col>
            </b-row>
          </template>
        </template>

        <slot
          name="dimension-options"
          :index="i"
          :dimension="d"
          :field="getField(d)"
        />
      </fieldset>
    </div>
    <hr v-if="!!module">

    <!-- Configure report metrics -->
    <div
      v-if="!!module"
      class="px-3"
    >
      <h5 class="d-flex align-items-center mb-3">
        {{ $t('edit.metric.title') }}
        <b-button
          v-if="canAddMetric"
          variant="link"
          class="text-decoration-none"
          @click.prevent="addMetric"
        >
          + {{ $t('edit.metric.add') }}
        </b-button>
      </h5>

      <draggable
        class="metrics mb-3"
        :list.sync="metrics"
        handle=".grab"
        :group="`metrics_${moduleID}`"
      >
        <div
          v-for="(m,i) in metrics"
          :key="i"
          class="rounded border border-light p-3 mb-3"
          style="background-color: #F9FAFB;"
        >
          <h5
            v-if="metrics.length > 1"
            class="d-flex align-items-center mb-3"
          >
            {{ $t('edit.metric.label') }} {{ i + 1 }}

            <div class="d-flex align-items-center ml-auto">
              <c-input-confirm
                class="mr-2"
                @confirmed="removeMetric(i)"
              />

              <b-button
                variant="link"
                size="sm"
                class="ml-auto px-0"
              >
                <font-awesome-icon
                  class="grab text-secondary"
                  :icon="['fas', 'bars']"
                />
              </b-button>
            </div>
          </h5>

          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('edit.metric.fieldLabel')"
                label-class="text-primary"
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
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('edit.metric.function.label')"
                label-class="text-primary"
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
            </b-col>
          </b-row>

          <slot
            name="metric-options"
            :metric="m"
          />
        </div>
      </draggable>
    </div>

    <hr v-if="!!module && hasAxis">

    <template v-if="hasAxis">
      <slot
        name="y-axis"
        :report="editReport"
      />
    </template>

    <hr v-if="hasLegend">

    <div
      v-if="hasLegend"
      class="px-3"
    >
      <h5 class="mb-3">
        {{ $t('edit.additionalConfig.legend.label') }}
      </h5>

      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.orientation.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="report.legend.orientation"
              :options="orientations"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.show')"
            label-class="text-primary"
          >
            <c-input-checkbox
              :value="!!report.legend.isHidden"
              switch
              invert
              :labels="checkboxLabel"
              @input="$set(report.legend,'isHidden', $event)"
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
            :label="$t('edit.additionalConfig.legend.align.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="report.legend.align"
              :options="alignments"
              :disabled="!report.legend.position.isDefault"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.options.label')"
            label-class="text-primary"
          >
            <b-form-checkbox
              v-model="report.legend.isScrollable"
              :disabled="report.legend.orientation !== 'horizontal'"
            >
              {{ $t('edit.additionalConfig.legend.scrollable') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="report.legend.position.isDefault"
            >
              {{ $t('edit.additionalConfig.legend.position.customize') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
      </b-row>

      <b-row
        v-if="!report.legend.position.isDefault"
      >
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.position.top')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.legend.position.top"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.position.right')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.legend.position.right"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.position.bottom')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.legend.position.bottom"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.legend.position.left')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.legend.position.left"
            />
          </b-form-group>
        </b-col>

        <b-col cols="12">
          <small class="text-muted">
            {{ $t('edit.additionalConfig.legend.valueRange') }}
          </small>
        </b-col>
      </b-row>
    </div>

    <slot
      name="additional-config"
      :report="editReport"
      :metrics="metrics"
      :has-axis="hasAxis"
    />
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
      metricAggregates: aggregateFunctions.map(af => ({ ...af, text: this.$t(`edit.metric.function.${af.text}`) })),
      dimensionModifiers: compose.chartUtil.dimensionFunctions.map(df => ({ ...df, text: this.$t(`edit.dimension.function.${df.text}`) })),
      predefinedFilters: compose.chartUtil.predefinedFilters.map(pf => ({ ...pf, text: this.$t(`edit.filter.${pf.text}`) })),

      alignments: [
        { value: 'left', text: this.$t('edit.additionalConfig.legend.align.left') },
        { value: 'center', text: this.$t('edit.additionalConfig.legend.align.center') },
        { value: 'right', text: this.$t('edit.additionalConfig.legend.align.right') },
      ],

      orientations: [
        { value: 'horizontal', text: this.$t('edit.additionalConfig.legend.orientation.horizontal') },
        { value: 'vertical', text: this.$t('edit.additionalConfig.legend.orientation.vertical') },
      ],
    }
  },

  computed: {
    hasLegend () {
      return !this.metrics.some(({ type }) => ['gauge'].includes(type))
    },

    defaultValueInputType () {
      return ({ field }) => (this.module.fields.filter(f => f.name === field)[0] || {}).kind === 'DateTime' ? 'date' : 'text'
    },

    defaultFilterOption () {
      return this.predefinedFilters.some(({ value }) => value === this.report.filter) ? '' : this.report.filter
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
          .sort((a, b) => (a.label || a.name).localeCompare((b.label || b.name)))
          .map(({ label, name }) => ({ value: name, text: label || name })),
      ]
    },

    hasAxis () {
      return this.metrics.some(({ type }) => ['bar', 'line', 'scatter'].includes(type))
    },

    dimensionFields () {
      return [
        ...[...this.module.fields].sort((a, b) => (a.label || a.name).localeCompare((b.label || b.name))),
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(({ kind, options = {} }) => {
        return this.dimensionFieldKind.includes(kind) && !(options.useRichTextEditor || options.multiLine)
      }).map(({ name, label, kind }) => {
        return { value: name, text: `${label || name} (${kind})`, kind }
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

      this.$set(d.meta, 'fields', [])
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
