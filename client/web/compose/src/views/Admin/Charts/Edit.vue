<template>
  <div class="py-3">
    <portal to="topbar-title">
      {{ $t('edit.title') }}
    </portal>

    <div
      v-if="!chart"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <b-container
      v-else
      tag="form"
      fluid="xl"
      @submit.prevent="handleSave"
    >
      <b-row no-gutters>
        <b-col>
          <b-card
            no-body
            class="shadow-sm"
          >
            <b-card-header
              header-bg-variant="white"
              class="d-flex py-3 align-items-center border-bottom"
            >
              <export
                slot="header"
                :list="[chart]"
                type="chart"
                class="float-right"
              />
            </b-card-header>
            <b-container
              fluid
              class="px-4 py-3"
            >
              <b-row>
                <b-col
                  md="6"
                  sm="12"
                >
                  <div v-if="modules">
                    <b-form-group
                      :label="$t('name')"
                    >
                      <b-form-input
                        v-model="chart.name"
                        :state="nameState"
                      />
                    </b-form-group>

                    <b-form-group
                      :label="$t('handle')"
                    >
                      <b-form-input
                        v-model="chart.handle"
                        :placeholder="$t('general.placeholder.handle')"
                        :state="handleState"
                        class="mb-1"
                      />
                      <b-form-invalid-feedback :state="handleState">
                        {{ $t('general.placeholder.invalid-handle-characters') }}
                      </b-form-invalid-feedback>
                    </b-form-group>

                    <b-form-group
                      :label="$t('colorScheme.label')"
                    >
                      <vue-select
                        v-model="chart.config.colorScheme"
                        :options="colorSchemes"
                        :reduce="cs => cs.value"
                        label="label"
                        option-text="label"
                        option-value="value"
                        :placeholder="$t('colorScheme.placeholder')"
                        :clearable="true"
                        class="bg-white h-100 w-100"
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
                  </div>

                  <!-- Some charts support multiple reports -->
                  <fieldset
                    v-if="supportsMultipleReports"
                    class="form-group mt-2"
                  >
                    <b-form-group class="mb-2">
                      <h4 class="d-inline-block">
                        {{ $t('configure.reportsLabel') }}
                      </h4>
                      <b-btn
                        v-if="reportsValid"
                        class="float-right p-0"
                        variant="link"
                        @click="onAddReport"
                      >
                        + {{ $t('general.label.add') }}
                      </b-btn>
                      <div class="ml-1">
                        <draggable
                          v-model="reports"
                          :options="{ handle:'.handle' }"
                          class="w-100 d-inline-block"
                          tag="tbody"
                        >
                          <report-item
                            v-for="(r, i) in reports"
                            :key="i"
                            :report="r"
                            :fixed="reports.length === 1"
                            @edit="onEditReport(i)"
                            @remove="onRemoveReport(i)"
                          >
                            <template #report-label>
                              <template v-if="r.moduleID">
                                {{ moduleName(r.moduleID) }}
                              </template>
                              <template v-else>
                                {{ $t('edit.unconfiguredReport') }}
                              </template>
                            </template>
                          </report-item>
                        </draggable>
                      </div>
                    </b-form-group>
                  </fieldset>

                  <!-- Generic report editing component -->
                  <component
                    :is="reportEditor"
                    v-if="editReport"
                    :report.sync="editReport"
                    :chart="chart"
                    :modules="modules"
                    :dimension-field-kind="['Select']"
                    :supported-metrics="1"
                  />
                </b-col>

                <b-col
                  md="6"
                  sm="12"
                >
                  <div
                    class="d-flex flex-column position-sticky"
                    style="top: 0;"
                  >
                    <b-button
                      :disabled="processing || !reportsValid"
                      variant="outline-primary"
                      @click.prevent="update"
                    >
                      {{ $t('edit.loadData') }}
                    </b-button>

                    <div
                      class="chart-preview mt-5"
                    >
                      <chart-component
                        ref="chart"
                        :chart="chart"
                        :reporter="reporter"
                        style="min-height: 400px;"
                        @updated="onUpdated"
                      />
                    </div>
                  </div>
                </b-col>
              </b-row>
            </b-container>
          </b-card>
        </b-col>
      </b-row>
    </b-container>

    <portal to="admin-toolbar">
      <editor-toolbar
        :processing="processing"
        :back-link="{ name: 'admin.charts' }"
        :hide-delete="hideDelete"
        :hide-save="hideSave"
        hide-clone
        :disable-save="disableSave"
        @delete="handleDelete"
        @save="handleSave()"
        @saveAndClose="handleSave({ closeOnSuccess: true })"
      />
    </portal>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import { compose, NoID, shared } from '@cortezaproject/corteza-js'
import Export from 'corteza-webapp-compose/src/components/Admin/Export'
import ChartComponent from 'corteza-webapp-compose/src/components/Chart'
import { handle } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import ReportItem from 'corteza-webapp-compose/src/components/Chart/ReportItem'
import Reports from 'corteza-webapp-compose/src/components/Chart/Report'
import { chartConstructor } from 'corteza-webapp-compose/src/lib/charts'
import VueSelect from 'vue-select'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { debounce } from 'lodash'
const { colorschemes } = shared

const defaultReport = {
  moduleID: undefined,
  metrics: [{ field: 'count' }],
  dimensions: [{ field: 'createdAt', modifier: 'MONTH' }],
}

export default {
  i18nOptions: {
    namespaces: 'chart',
  },

  components: {
    EditorToolbar,
    Export,
    ChartComponent,
    draggable,
    ReportItem,
    VueSelect,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    chartID: {
      type: String,
      required: false,
      default: NoID,
    },

    category: {
      type: String,
      required: false,
      default: '',
    },
  },

  data () {
    return {
      chart: undefined,
      processing: false,

      editReportIndex: undefined,
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      modByID: 'module/getByID',
    }),

    colorSchemes () {
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
            label: `${capitalize(g)}: ${capitalize(gn.label)} (${this.$t('colorLabel', gn)})`,
            colors: [...colorschemes[g][sc]],
            value: `${g}.${sc}`,
          })
        }
      }

      return rr
    },

    currentColorScheme () {
      return this.colorSchemes.find(({ value }) => value === this.chart.config.colorScheme)
    },

    defaultReport () {
      return Object.assign({}, defaultReport)
    },

    nameState () {
      return this.chart.name.length > 0 ? null : false
    },

    handleState () {
      return handle.handleState(this.chart.handle)
    },

    supportsMultipleReports () {
      if (!this.chart) {
        return false
      }

      if (this.chart instanceof compose.FunnelChart) {
        return true
      }
      return false
    },

    reportsValid () {
      if (!this.reports) {
        return false
      }

      return !this.reports.find(({ moduleID }) => !moduleID)
    },

    reportEditor () {
      if (!this.chart) {
        return undefined
      }

      if (this.chart instanceof compose.FunnelChart) {
        return Reports.FunnelChart
      }
      if (this.chart instanceof compose.GaugeChart) {
        return Reports.GaugeChart
      }
      return Reports.GenericChart
    },

    reports: {
      get () {
        return this.chart.config.reports
      },
      set (r) {
        this.chart.config.reports = r
      },
    },

    editReport: {
      get () {
        if (this.editReportIndex !== undefined) {
          return this.reports[this.editReportIndex]
        }
        return undefined
      },
      set (v) {
        this.reports.splice(this.editReportIndex, 1, v)
      },
    },

    disableSave () {
      return !this.chart || [this.nameState, this.handleState].includes(false)
    },

    hideDelete () {
      return !this.isEdit || !this.chart.canDeleteChart || !!this.chart.deletedAt
    },

    hideSave () {
      return this.isEdit && !this.chart.canUpdateChart
    },

    isEdit () {
      return this.chart && this.chart.chartID !== NoID
    },
  },

  watch: {
    chartID: {
      immediate: true,
      handler (chartID) {
        this.chart = undefined
        const { namespaceID } = this.namespace

        if (chartID === NoID) {
          let c = new compose.Chart({ namespaceID: this.namespace.namespaceID })
          switch (this.category) {
            case 'gauge':
              c = new compose.GaugeChart(c)
              break

            case 'funnel':
              c = new compose.FunnelChart(c)
              break
          }
          this.chart = c
          this.onEditReport(0)
        } else {
          this.findChartByID({ namespaceID, chartID, force: true }).then((chart) => {
            // Make a copy so that we do not change store item by ref
            this.chart = chartConstructor(chart)
            this.onEditReport(0)
          }).catch(this.toastErrorHandler(this.$t('notification:chart.loadFailed')))
        }
      },
    },

    'chart.config': {
      deep: true,
      handler () {
        debounce(function () {
          this.update()
        }, 500)
      },
    },
  },

  methods: {
    ...mapActions({
      findChartByID: 'chart/findByID',
      createChart: 'chart/create',
      updateChart: 'chart/update',
      deleteChart: 'chart/delete',
    }),

    moduleName (moduleID) {
      return this.modByID(moduleID).name
    },

    reporter (r) {
      const nr = { ...r }
      if (nr.filter) {
        nr.filter = evaluatePrefilter(nr.filter, {
          record: this.record,
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      }
      return this.$ComposeAPI.recordReport({ namespaceID: this.namespace.namespaceID, ...nr })
    },

    update () {
      this.processing = true
      this.$refs.chart.updateChart()
    },

    onUpdated () {
      this.processing = false
    },

    handleSave ({ closeOnSuccess = false } = {}) {
      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage

      const c = Object.assign({}, this.chart, resourceTranslationLanguage)

      if (this.chart.chartID === NoID) {
        this.createChart(c).then(({ chartID }) => {
          this.toastSuccess(this.$t('notification:chart.saved'))
          if (closeOnSuccess) {
            this.redirect()
          } else {
            this.$router.push({ name: 'admin.charts.edit', params: { chartID: chartID } })
          }
        }).catch(this.toastErrorHandler(this.$t('notification:chart.saveFailed')))
      } else {
        this.updateChart(c).then((chart) => {
          this.chart = chartConstructor(chart)
          this.update()
          this.toastSuccess(this.$t('notification:chart.saved'))
          if (closeOnSuccess) {
            this.redirect()
          }
        }).catch(this.toastErrorHandler(this.$t('notification:chart.saveFailed')))
      }
    },

    handleDelete () {
      this.deleteChart(this.chart).then(() => {
        this.toastSuccess(this.$t('notification:chart.deleted'))
        this.$router.push({ name: 'admin.charts' })
      }).catch(this.toastErrorHandler(this.$t('notification:chart.deleteFailed')))
    },

    redirect () {
      this.$router.push({ name: 'admin.charts' })
    },

    onEditReport (i) {
      this.editReportIndex = i
    },

    onRemoveReport (i) {
      this.reports.splice(i, 1)
      if (this.editReportIndex === i) {
        this.editReportIndex = undefined
      }
    },

    onAddReport () {
      this.reports.push(this.chart.defReport())
    },
  },
}
</script>
<style lang="scss" scoped>
.chart-preview {
  max-height: 50vh;
}

.color-box {
  width: 28px;
  height: 12px;
}
</style>
