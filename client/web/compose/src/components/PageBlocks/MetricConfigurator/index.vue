<template>
  <div>
    <b-tab
      :title="$t('metric.edit.tabTitle')"
    >
      <b-row no-gutters>
        <b-col cols="12">
          <div
            v-for="(m, i) in metrics"
            :key="i"
            class="mb-2"
          >
            <b-btn
              variant="light"
              class="mr-1"
              @click="editMetric(m)"
            >
              {{ $t('general.label.edit') }}
            </b-btn>
            <b-btn
              variant="outline-danger"
              class="mr-2"
              @click="removeMetric(i)"
            >
              {{ $t('general.label.remove') }}
            </b-btn>

            <span class="btn">
              {{ m.label || $t('metric.defaultMetricLabel') }}
            </span>
          </div>

          <b-btn
            variant="link"
            class="px-1"
            @click="addMetric"
          >
            + {{ $t('general.label.add') }}
          </b-btn>
        </b-col>
      </b-row>

      <hr>

      <b-row
        class="mt-3"
      >
        <!-- edit metric -->
        <b-col
          v-if="edit"
          cols="12"
          lg="7"
        >
          <b-card
            class="mb-5"
            no-body
          >
            <fieldset>
              <b-form-group
                :label="$t('metric.edit.labelLabel')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="edit.label"
                  :placeholder="$t('metric.edit.labelPlaceholder')"
                  class="mb-1"
                />
              </b-form-group>
            </fieldset>

            <fieldset>
              <h5>
                {{ $t('metric.edit.dimensionLabel') }}
              </h5>

              <b-form-group
                :label="$t('metric.edit.moduleLabel')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="edit.moduleID"
                  :options="modules"
                  label="name"
                  class="mt-1"
                  :reduce="o => o.moduleID"
                  :placeholder="$t('metric.edit.modulePlaceholder')"
                />
              </b-form-group>

              <!-- <b-form-group
                :label="$t('metric.edit.dimensionFieldLabel')"
                label-cols="2"
              >
                <b-form-select
                  v-model="edit.dimensionField"
                  :options="fields"
                  class="mt-1"
                  text-field="label"
                  value-field="name"
                >
                  <template slot="first">
                    <option
                      :value="undefined"
                      disabled
                    >
                      {{ $t('metric.edit.dimensionFieldPlaceholder') }}
                    </option>
                  </template>
                </b-form-select>
              </b-form-group>

              <template>
                <b-form-group
                  :label="$t('metric.edit.dateFormat')"
                  label-cols="2"
                >
                  <b-form-input
                    v-model="edit.dateFormat"
                    :disabled="!edit.dimensionField || !isTemporalField(edit.dimensionField)"
                    placeholder="YYYY-MM-DD"
                    class="mb-1"
                  />
                </b-form-group>

                <b-form-group
                  :label="$t('metric.edit.bucketLabel')"
                  label-cols="2"
                >
                  <b-form-select
                    v-model="edit.bucketSize"
                    :disabled="!edit.dimensionField || !isTemporalField(edit.dimensionField)"
                    :options="dimensionModifiers"
                  >
                    <template slot="first">
                      <option
                        disabled
                        :value="undefined"
                      >
                        {{ $t('metric.edit.bucketPlaceholder') }}
                      </option>
                    </template>
                  </b-form-select>
                </b-form-group>
              </template> -->

              <b-form-group
                :label="$t('metric.edit.filterLabel')"
                label-class="text-primary"
              >
                <c-input-expression
                  v-model="edit.filter"
                  auto-complete
                  lang="javascript"
                  placeholder="(A > B) OR (A < C)"
                  class="mb-1"
                  height="3.448rem"
                  :suggestion-params="recordAutoCompleteParams"
                />

                <i18next
                  path="metric.edit.filterFootnote"
                  tag="small"
                  class="d-block text-muted"
                >
                  <code>${record.values.fieldName}</code>
                  <code>${recordID}</code>
                  <code>${ownerID}</code>
                  <span><code>${userID}</code>, <code>${user.name}</code></span>
                </i18next>
              </b-form-group>
            </fieldset>

            <fieldset v-if="selectedMetricModule">
              <h5>
                {{ $t('metric.edit.metricLabel') }}
              </h5>

              <b-form-group
                :label="$t('metric.edit.metricFieldLabel')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="edit.metricField"
                  :placeholder="$t('metric.edit.metricFieldSelect')"
                  :options="metricFields"
                  :get-option-key="getOptionMetricFieldKey"
                  :get-option-label="getOptionMetricFieldLabel"
                  :reduce="f => f.name"
                  @input="onMetricFieldChange"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.metricAggregateLabel')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="edit.operation"
                  :disabled="edit.metricField === 'count'"
                  :placeholder="$t('metric.edit.metricSelectAggregate')"
                  :options="aggregationOperations"
                  :get-option-key="getOptionAggregationOperationKey"
                  :reduce="a => a.operation"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.transformFunctionLabel')"
                label-class="text-primary"
              >
                <c-input-expression
                  v-model="edit.transformFx"
                  auto-complete
                  lang="javascript"
                  placeholder="v"
                  class="mb-1"
                  height="3.448rem"
                  :suggestion-params="recordAutoCompleteParams"
                />

                <small>{{ $t('metric.edit.transformFunctionDescription') }}</small>
                <i18next
                  path="metric.edit.transformFootnote"
                  tag="small"
                  class="d-block text-muted"
                >
                  <code>${record.values.fieldName}</code>
                  <code>${recordID}</code>
                  <code>${ownerID}</code>
                  <span><code>${userID}</code>, <code>${user.name}</code></span>
                </i18next>
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.numberFormat')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="edit.numberFormat"
                  placeholder="0.00"
                  class="mb-1"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.prefixLabel')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="edit.prefix"
                  placeholder="$"
                  class="mb-1"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.suffixLabel')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model="edit.suffix"
                  placeholder="USD/mo"
                  class="mb-1"
                />
              </b-form-group>

              <b-form-group
                :description="$t('metric.drillDown.description')"
                label-class="d-flex align-items-center text-primary"
                class="mb-1"
              >
                <template #label>
                  {{ $t('metric.drillDown.label') }}

                  <b-form-checkbox
                    v-model="edit.drillDown.enabled"
                    switch
                    class="ml-1 mb-1"
                  />
                </template>

                <b-input-group>
                  <c-input-select
                    v-model="edit.drillDown.blockID"
                    :options="drillDownOptions"
                    :disabled="!edit.drillDown.enabled"
                    :get-option-label="o => o.title || o.kind"
                    :reduce="option => option.blockID"
                    :clearable="true"
                    :placeholder="$t('metric.drillDown.openInModal')"
                    append-to-body
                  />

                  <b-input-group-append>
                    <column-picker
                      :module="selectedMetricModule"
                      :disabled="!!edit.drillDown.blockID || !edit.drillDown.enabled"
                      :fields="selectedDrilldownFields"
                      variant="extra-light"
                      size="md"
                      @updateFields="onUpdateFields"
                    >
                      <font-awesome-icon :icon="['fas', 'wrench']" />
                    </column-picker>
                  </b-input-group-append>
                </b-input-group>
              </b-form-group>
            </fieldset>
          </b-card>

          <m-style
            class="mt-2"
            :options="edit.valueStyle"
          >
            <h5 slot="title">
              {{ $t('metric.editStyle.valueLabel') }}
            </h5>
          </m-style>
        </b-col>

        <b-col
          cols="12"
          lg="5"
        >
          <div
            v-if="metrics.length"
            class="d-flex flex-column position-sticky pt-2"
            style="top: 0;"
          >
            <b-button
              v-b-tooltip.noninteractive.hover="{ title: $t('metric.edit.refreshData'), container: '#body' }"
              variant="outline-light"
              size="lg"
              class="d-flex align-items-center text-primary ml-auto border-0 px-2 mt-2 mr-2"
              @click.prevent="$root.$emit('metric.update')"
            >
              <font-awesome-icon :icon="['fa', 'sync']" />
            </b-button>

            <div
              class="mt-2"
              style="height: 400px;"
            >
              <metric-base
                v-bind="$props"
              />
            </div>
          </div>
        </b-col>
      </b-row>
    </b-tab>
  </div>
</template>

<script>
import base from '../base'
import MStyle from './MStyle'
import { mapGetters } from 'vuex'
import MetricBase from '../MetricBase'
import ColumnPicker from 'corteza-webapp-compose/src/components/Admin/Module/Records/ColumnPicker'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Metric',
  components: {
    MStyle,
    MetricBase,
    ColumnPicker,
    CInputExpression,
  },
  extends: base,

  mixins: [autocomplete],

  data () {
    return {
      edit: undefined,
      dimensionModifiers: compose.chartUtil.dimensionFunctions.map(df => ({ ...df, text: this.$t(`chart.edit.dimension.function.${df.text}`) })),
      aggregationOperations: [
        {
          label: this.$t('metric.edit.operationSum'),
          operation: 'sum',
        },
        {
          label: this.$t('metric.edit.operationMax'),
          operation: 'max',
        },
        {
          label: this.$t('metric.edit.operationMin'),
          operation: 'min',
        },
        {
          label: this.$t('metric.edit.operationAvg'),
          operation: 'avg',
        },
      ],
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      getModuleByID: 'module/getByID',
    }),

    fields () {
      if (!this.edit || !this.edit.moduleID) {
        return []
      }

      return this.getModuleByID(this.edit.moduleID).fields
    },

    selectedDrilldownFields () {
      if (!this.edit || !this.edit.drillDown.recordListOptions.fields) return []

      return this.edit.drillDown.recordListOptions.fields
    },

    metricFields () {
      return [
        { name: 'count', label: 'Count' },
        ...this.fields.filter(f => f.kind === 'Number')
          .sort((a, b) => a.label.localeCompare(b.label)),
      ]
    },

    metrics: {
      get () {
        return this.options.metrics
      },
      set (m) {
        this.options.metrics = m
      },
    },

    drillDownOptions () {
      return this.page.blocks.filter(({ blockID, kind, options = {} }) => kind === 'RecordList' && blockID !== NoID && options.moduleID === this.edit.moduleID)
    },

    selectedMetricModule () {
      if (!this.edit.moduleID) return undefined

      return this.getModuleByID(this.edit.moduleID)
    },

    recordAutoCompleteParams () {
      return this.processRecordAutoCompleteParams({ module: this.selectedMetricModule })
    },
  },

  watch: {
    'edit.dimensionField': function (df) {
      if (!this.isTemporalField(df)) {
        this.edit.bucketSize = undefined
        this.edit.dateFormat = undefined
      } else {
        this.edit.dateFormat = this.edit.dateFormat || 'YYYY-MM-DD'
      }
    },
  },

  created () {
    if (!this.metrics.length) {
      this.addMetric()
    }

    this.edit = this.metrics[0]
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    addMetric () {
      const m = this.block.makeMetric()
      this.metrics.push(m)
      this.editMetric(m)
    },

    editMetric (m) {
      this.edit = m
    },

    removeMetric (i) {
      this.metrics.splice(i, 1)
      this.edit = undefined
    },

    isTemporalField (name) {
      return !!this.fields.find(f => f.name === name && f.kind === 'DateTime')
    },

    getOptionMetricFieldKey ({ name }) {
      return name
    },

    getOptionMetricFieldLabel ({ name, label }) {
      return label || name
    },

    getOptionAggregationOperationKey ({ operation }) {
      return operation
    },

    onMetricFieldChange (field) {
      if (field === 'count') {
        this.edit.operation = undefined
      } else if (!this.edit.operation) {
        this.edit.operation = this.aggregationOperations[0].operation
      }
    },

    onUpdateFields (fields) {
      this.edit.drillDown.recordListOptions.fields = fields
    },

    setDefaultValues () {
      this.edit = undefined
      this.dimensionModifiers = []
      this.aggregationOperations = []
    },
  },
}
</script>
