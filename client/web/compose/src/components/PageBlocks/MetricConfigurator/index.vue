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
              >
                <b-form-select
                  v-model="edit.moduleID"
                  :options="modules"
                  text-field="name"
                  class="mt-1"
                  value-field="moduleID"
                >
                  <template slot="first">
                    <option
                      :value="undefined"
                      disabled
                    >
                      {{ $t('metric.edit.modulePlaceholder') }}
                    </option>
                  </template>
                </b-form-select>
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

              <b-form-group :label="$t('metric.edit.filterLabel')">
                <b-form-textarea
                  v-model="edit.filter"
                  placeholder="(A > B) OR (A < C)"
                  class="mb-1"
                />

                <b-form-text>
                  <i18next
                    path="metric.edit.filterFootnote"
                    tag="label"
                  >
                    <code>${recordID}</code>
                    <code>${ownerID}</code>
                    <code>${userID}</code>
                  </i18next>
                </b-form-text>
              </b-form-group>
            </fieldset>

            <fieldset>
              <h5>
                {{ $t('metric.edit.metricLabel') }}
              </h5>

              <b-form-group
                :label="$t('metric.edit.metricFieldLabel')"
              >
                <vue-select
                  v-model="edit.metricField"
                  :placeholder="$t('metric.edit.metricFieldSelect')"
                  :options="metricFields"
                  :reduce="f => f.name"
                  class="bg-white"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.metricAggregateLabel')"
              >
                <vue-select
                  v-model="edit.operation"
                  :disabled="edit.metricField === 'count'"
                  :placeholder="$t('metric.edit.metricSelectAggregate')"
                  :options="aggregationOperations"
                  :reduce="a => a.operation"
                  class="bg-white"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.transformFunctionLabel')"
              >
                <b-form-textarea
                  v-model="edit.transformFx"
                  placeholder="v"
                  class="mb-1"
                />

                <small>{{ $t('metric.edit.transformFunctionDescription') }}</small>
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.numberFormat')"
              >
                <b-form-input
                  v-model="edit.numberFormat"
                  placeholder="0.00"
                  class="mb-1"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.prefixLabel')"
              >
                <b-form-input
                  v-model="edit.prefix"
                  placeholder="$"
                  class="mb-1"
                />
              </b-form-group>

              <b-form-group
                :label="$t('metric.edit.suffixLabel')"
              >
                <b-form-input
                  v-model="edit.suffix"
                  placeholder="USD/mo"
                  class="mb-1"
                />
              </b-form-group>
            </fieldset>
          </b-card>

          <!-- <m-style :options="edit.labelStyle">
            <h5 slot="title">
              {{ $t('metric.editStyle.labelLabel') }}
            </h5>
          </m-style> -->
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
              variant="outline-primary"
              @click.prevent="$root.$emit('metric.update')"
            >
              {{ $t('metric.edit.refreshData') }}
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
import { VueSelect } from 'vue-select'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Metric',
  components: {
    MStyle,
    MetricBase,
    VueSelect,
  },
  extends: base,

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
      moduleByID: 'module/getByID',
    }),

    fields () {
      if (!this.edit || !this.edit.moduleID) {
        return []
      }

      return this.moduleByID(this.edit.moduleID).fields
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
    'edit.metricField': function (mf) {
      if (mf === 'count') {
        this.edit.operation = undefined
      }
    },
  },

  created () {
    if (!this.metrics.length) {
      this.addMetric()
    }

    this.edit = this.metrics[0]
  },

  methods: {
    addMetric () {
      const m = {
        labelStyle: {},
        valueStyle: {
          backgroundColor: '#ffffff',
        },
      }
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
  },
}
</script>
