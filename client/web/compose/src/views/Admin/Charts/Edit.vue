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
      <b-row>
        <b-col>
          <b-card
            no-body
            class="shadow-sm"
          >
            <b-card-header class="d-flex py-3 align-items-center border-bottom gap-1">
              <export
                v-if="chart.canExportChart"
                slot="header"
                :list="[chart]"
                type="chart"
              />

              <c-permissions-button
                v-if="namespace.canGrant"
                :title="chart.name || chart.handle || chart.chartID"
                :target="chart.name || chart.handle || chart.chartID"
                :resource="`corteza::compose:chart/${namespace.namespaceID}/${chart.chartID}`"
                :button-label="$t('general.label.permissions')"
                class="btn-lg"
              />
            </b-card-header>

            <b-row no-gutters>
              <b-col
                cols="12"
                lg="7"
                class="border-right"
              >
                <div
                  class="pt-3 px-3"
                >
                  <h5>
                    {{ $t('generalSettings') }}
                  </h5>
                  <b-row
                    v-if="modules"
                  >
                    <b-col
                      cols="12"
                      lg="6"
                    >
                      <b-form-group
                        :label="$t('name')"
                        label-class="text-primary"
                        class="text-primary"
                      >
                        <b-form-input
                          v-model="chart.name"
                          :state="nameState"
                        />
                      </b-form-group>
                    </b-col>

                    <b-col
                      cols="12"
                      lg="6"
                    >
                      <b-form-group
                        :label="$t('handle')"
                        label-class="text-primary"
                        class="text-primary"
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
                    </b-col>

                    <b-col
                      cols="12"
                      lg="6"
                    >
                      <b-form-group
                        :label="$t('colorScheme.label')"
                        label-class="text-primary"
                      >
                        <b-input-group class="d-flex w-100">
                          <c-input-select
                            v-model="chart.config.colorScheme"
                            :options="colorSchemes"
                            :reduce="cs => cs.id"
                            label="name"
                            :get-option-key="o => o.id"
                            :placeholder="$t('colorScheme.placeholder')"
                          >
                            <template #option="option">
                              <p
                                class="mb-1"
                              >
                                {{ option.name }}
                              </p>

                              <div
                                v-for="(color, index) in option.colors"
                                :key="index"
                                :style="`background: ${color};`"
                                class="d-inline-block color-box mr-1 mb-1"
                              />
                            </template>

                            <template
                              v-if="canManageColorSchemes"
                              #list-header
                            >
                              <li class="border-bottom text-center mb-1">
                                <b-button
                                  variant="link"
                                  class="text-decoration-none"
                                  @click="createColorScheme"
                                >
                                  {{ $t('colorScheme.custom.add') }}
                                </b-button>
                              </li>
                            </template>
                          </c-input-select>

                          <b-input-group-append v-if="showEditColorSchemeButton">
                            <b-button
                              v-b-tooltip.noninteractive.hover="{ title: $t('colorScheme.custom.edit'), container: '#body' }"
                              variant="extra-light"
                              class="d-flex align-items-center"
                              @click="editColorScheme()"
                            >
                              <font-awesome-icon :icon="['far', 'edit']" />
                            </b-button>
                          </b-input-group-append>
                        </b-input-group>

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
                      class="mt-2 mt-md-0"
                    >
                      <b-form-group
                        :label="$t('edit.animation.label')"
                        label-class="text-primary"
                        class="text-primary"
                      >
                        <c-input-checkbox
                          v-model="chart.config.noAnimation"
                          :labels="checkboxLabel"
                          switch
                          invert
                        />
                      </b-form-group>
                    </b-col>
                  </b-row>
                </div>
                <hr v-if="modules">

                <!-- Some charts support multiple reports -->
                <!-- <fieldset
                  v-if="supportsMultipleReports"
                  class="form-group"
                >
                  <b-form-group class=" px-3">
                    <h5 class="d-inline-block">
                      {{ $t('configure.reportsLabel') }}
                    </h5>

                    <b-btn
                      v-if="reportsValid"
                      class="float-right p-0"
                      variant="link"
                      @click="onAddReport"
                    >
                      + {{ $t('general.label.add') }}
                    </b-btn>

                    <div>
                      <draggable
                        v-model="reports"
                        handle=".handle"
                        tag="tbody"
                        class="w-100 d-inline-block"
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
                </fieldset> -->

                <!-- General report editing component -->
                <component
                  :is="reportEditor"
                  v-if="chart && editReport"
                  :report.sync="editReport"
                  :chart="chart"
                  :modules="modules"
                  :supported-metrics="1"
                />

                <hr>

                <div
                  class="px-3"
                >
                  <h5 class="mb-3">
                    {{ $t('edit.toolbox.label') }}
                  </h5>

                  <b-row>
                    <b-col
                      cols="12"
                      lg="6"
                    >
                      <b-form-group
                        :label="$t('edit.toolbox.saveAsImage.label')"
                        label-class="text-primary"
                      >
                        <c-input-checkbox
                          :value="!!chart.config.toolbox.saveAsImage"
                          switch
                          :labels="checkboxLabel"
                          @input="$set(chart.config.toolbox, 'saveAsImage', $event)"
                        />
                      </b-form-group>
                    </b-col>

                    <b-col
                      v-if="hasAxis"
                      cols="12"
                      lg="6"
                    >
                      <b-form-group
                        :label="$t('edit.toolbox.timeline.label')"
                        label-class="text-primary"
                      >
                        <b-form-radio-group
                          v-model="chart.config.toolbox.timeline"
                          buttons
                          button-variant="outline-secondary"
                          size="sm"
                          :options="timelineOptions"
                        />
                      </b-form-group>
                    </b-col>
                  </b-row>
                </div>
              </b-col>

              <b-col
                cols="12"
                lg="5"
              >
                <div
                  class="d-flex flex-column position-sticky"
                  style="top: 0;"
                >
                  <b-button
                    v-b-tooltip.noninteractive.hover="{ title: $t('edit.loadData'), container: '#body' }"
                    :disabled="processing || !reportsValid"
                    variant="outline-light"
                    size="lg"
                    class="d-flex align-items-center text-primary ml-auto border-0 px-2 mt-2 mr-2"
                    @click.prevent="update"
                  >
                    <font-awesome-icon :icon="['fa', 'sync']" />
                  </b-button>

                  <chart-component
                    ref="chart"
                    :chart="chart"
                    :reporter="reporter"
                    style="min-height: 400px;"
                    @updated="onUpdated"
                  />
                </div>
              </b-col>
            </b-row>
          </b-card>
        </b-col>
      </b-row>
    </b-container>

    <b-modal
      v-model="colorSchemeModal.show"
      :title="colorSchemeModalTitle"
      :ok-title="$t('general:label.saveAndClose')"
      centered
      size="md"
      cancel-variant="light"
      no-fade
    >
      <b-form-group
        :label="$t('colorScheme.custom.modal.name.label')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="colorSchemeModal.colorscheme.name"
        />
      </b-form-group>

      <b-form-group
        label-class="text-primary"
        class="mb-0"
      >
        <template #label>
          {{ $t('colorScheme.custom.modal.colors.label') }}
          <b-button
            variant="outline-light"
            size="sm"
            class="text-primary border-0"
            @click="addColor"
          >
            <font-awesome-icon :icon="['fa', 'plus']" />
          </b-button>
        </template>

        <c-input-color-picker
          v-for="(color, index) in colorSchemeModal.colorscheme.colors"
          :key="index"
          v-model="colorSchemeModal.colorscheme.colors[index]"
          :show-text="false"
          data-test-id="input-scheme-color"
          :translations="{
            modalTitle: $t('colorScheme.pickAColor'),
            light: $t('general:themes.labels.light'),
            dark: $t('general:themes.labels.dark'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          :theme-settings="themeSettings"
          class="d-inline-flex mr-1"
        >
          <template #footer>
            <c-input-confirm
              variant="danger"
              size="md"
              show-icon
              @confirmed="removeColor(index)"
            />
          </template>
        </c-input-color-picker>
      </b-form-group>

      <template #modal-footer>
        <c-input-confirm
          v-if="colorSchemeModal.edit"
          :disabled="colorSchemeModal.processing"
          variant="danger"
          size="md"
          show-icon
          @confirmed="deleteColorScheme()"
        />

        <b-button
          variant="light"
          class="ml-auto"
          :disabled="colorSchemeModal.processing"
          @click="closeColorSchemeModal"
        >
          {{ $t('general:label.cancel') }}
        </b-button>

        <b-button
          variant="primary"
          :disabled="!colorSchemeModal.colorscheme.name || !colorSchemeModal.colorscheme.colors.length || colorSchemeModal.processing"
          @click="saveColorScheme"
        >
          {{ $t('general:label.saveAndClose') }}
        </b-button>
      </template>
    </b-modal>

    <portal to="admin-toolbar">
      <editor-toolbar
        :processing="processing"
        :processing-save="processingSave"
        :processing-clone="processingClone"
        :processing-save-and-close="processingSaveAndClose"
        :processing-delete="processingDelete"
        :hide-delete="hideDelete"
        :hide-save="hideSave"
        :hide-clone="!isEdit"
        :disable-save="disableSave"
        @delete="handleDelete()"
        @save="handleSave()"
        @clone="handleClone()"
        @saveAndClose="handleSave({ closeOnSuccess: true })"
        @back="$router.push(previousPage || { name: 'admin.charts' })"
      />
    </portal>
  </div>
</template>
<script>
import { isEqual, debounce, cloneDeep } from 'lodash'
import { mapGetters, mapActions } from 'vuex'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import { compose, NoID, shared } from '@cortezaproject/corteza-js'
import Export from 'corteza-webapp-compose/src/components/Admin/Export'
import ChartComponent from 'corteza-webapp-compose/src/components/Chart'
import { handle, components } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import ReportItem from 'corteza-webapp-compose/src/components/Chart/ReportItem'
import Reports from 'corteza-webapp-compose/src/components/Chart/Report'
import { chartConstructor } from 'corteza-webapp-compose/src/lib/charts'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

const { CInputCheckbox, CInputColorPicker } = components
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
    CInputCheckbox,
    CInputColorPicker,
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChart(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChart(next)
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
      initialChartState: undefined,
      processing: false,
      processingSave: false,
      processingClone: false,
      processingSaveAndClose: false,
      processingDelete: false,

      editReportIndex: undefined,

      customColorSchemes: [],
      colorSchemeModal: {
        show: false,
        processing: false,
        edit: false,
        colorscheme: {},
      },

      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      modByID: 'module/getByID',
      previousPage: 'ui/previousPage',
      can: 'rbac/can',
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

      const rr = [...this.customColorSchemes]
      for (const g in colorschemes) {
        for (const sc in colorschemes[g]) {
          const gn = splicer(sc)

          rr.push({
            id: `${g}.${sc}`,
            name: `${capitalize(g)}: ${capitalize(gn.label)} (${this.$t('colorLabel', gn)})`,
            colors: [...colorschemes[g][sc]],
          })
        }
      }

      return rr
    },

    currentColorScheme () {
      return this.colorSchemes.find(({ id }) => id === this.chart.config.colorScheme)
    },

    canManageColorSchemes () {
      return this.can('system/', 'settings.manage')
    },

    colorSchemeModalTitle () {
      return this.$t(`colorScheme.custom.${this.colorSchemeModal.edit ? 'edit' : 'add'}`)
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
      } else if (this.chart instanceof compose.GaugeChart) {
        return Reports.GaugeChart
      } else if (this.chart instanceof compose.RadarChart) {
        return Reports.RadarChart
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

    hasAxis () {
      return this.reports.some(({ metrics = [] }) => metrics.some(m => ['bar', 'line', 'scatter'].includes(m.type)))
    },

    timelineOptions () {
      return [
        { value: '', text: this.$t('edit.toolbox.timeline.options.none') },
        { value: 'x', text: this.$t('edit.toolbox.timeline.options.x') },
        { value: 'y', text: this.$t('edit.toolbox.timeline.options.y') },
        { value: 'xy', text: this.$t('edit.toolbox.timeline.options.xy') },
      ]
    },

    showEditColorSchemeButton () {
      const { config = {} } = this.chart || {}
      return config.colorScheme && config.colorScheme.includes('custom') && this.canManageColorSchemes
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },
  },

  watch: {
    chartID: {
      immediate: true,
      handler (chartID) {
        this.chart = undefined
        this.initialChartState = undefined

        const { namespaceID } = this.namespace

        if (this.canManageColorSchemes) {
          this.fetchCustomColorSchemes()
        }

        if (chartID === NoID) {
          let c = new compose.Chart({ namespaceID: this.namespace.namespaceID })

          switch (this.category) {
            case 'gauge':
              c = new compose.GaugeChart(c)
              break

            case 'funnel':
              c = new compose.FunnelChart(c)
              break

            case 'radar':
              c = new compose.RadarChart(c)
              break
          }
          this.chart = c
          this.initialChartState = cloneDeep(c)
          this.onEditReport(0)
        } else {
          this.findChartByID({ namespaceID, chartID, force: true }).then((chart) => {
            // Make a copy so that we do not change store item by ref
            this.chart = chartConstructor(chart)
            this.initialChartState = cloneDeep(chartConstructor(chart))
            this.onEditReport(0)
          }).catch(this.toastErrorHandler(this.$t('notification:chart.loadFailed')))
        }
      },
    },

    'chart.config': {
      deep: true,
      handler (value, oldValue) {
        if (value && oldValue) {
          this.onConfigUpdate()
        }
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
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      }

      return this.$ComposeAPI.recordReport({ namespaceID: this.namespace.namespaceID, ...nr })
    },

    update () {
      this.$refs.chart.updateChart()
    },

    onConfigUpdate: debounce(function () {
      this.update()
    }, 300),

    onUpdated () {
      this.processing = false
    },

    toggleProcessing ({ closeOnSuccess = false, isClone = false }) {
      this.processing = !this.processing

      if (closeOnSuccess) {
        this.processingSaveAndClose = !this.processingSaveAndClose
      } else if (isClone) {
        this.processingClone = !this.processingClone
      } else {
        this.processingSave = !this.processingSave
      }
    },

    handleSave ({ chart = this.chart, closeOnSuccess = false, isClone = false } = {}) {
      const toggleSaveProcessing = () => {
        this.processing = !this.processing

        if (closeOnSuccess) {
          this.processingSaveAndClose = !this.processingSaveAndClose
        } else if (isClone) {
          this.processingClone = !this.processingClone
        } else {
          this.processingSave = !this.processingSave
        }
      }

      toggleSaveProcessing()

      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage

      const c = Object.assign({}, chart, resourceTranslationLanguage)

      if (chart.chartID === NoID) {
        this.createChart(c).then(({ chartID }) => {
          this.chart = chartConstructor(chart)
          this.initialChartState = cloneDeep(chartConstructor(this.chart))

          this.toastSuccess(this.$t('notification:chart.saved'))
          if (closeOnSuccess) {
            this.redirect()
          } else {
            this.$router.push({ name: 'admin.charts.edit', params: { chartID: chartID } })
          }
        })
          .catch(this.toastErrorHandler(this.$t('notification:chart.saveFailed')))
          .finally(() => {
            toggleSaveProcessing()
          })
      } else {
        this.updateChart(c).then((chart) => {
          this.chart = chartConstructor(chart)
          this.initialChartState = cloneDeep(chartConstructor(chart))

          this.toastSuccess(this.$t('notification:chart.saved'))
          if (closeOnSuccess) {
            this.redirect()
          }
        })
          .catch(this.toastErrorHandler(this.$t('notification:chart.saveFailed')))
          .finally(() => {
            toggleSaveProcessing()
          })
      }
    },

    handleDelete () {
      this.processing = true
      this.processingDelete = true

      this.deleteChart(this.chart).then(() => {
        this.chart.deletedAt = new Date()

        this.toastSuccess(this.$t('notification:chart.deleted'))
        this.$router.push({ name: 'admin.charts' })
      })
        .catch(this.toastErrorHandler(this.$t('notification:chart.deleteFailed')))
        .finally(() => {
          this.processing = false
          this.processingDelete = false
        })
    },

    handleClone () {
      const chart = this.chart.clone()
      chart.chartID = NoID
      chart.name = `${this.chart.name} (copy)`
      chart.handle = this.chart.handle ? `${this.chart.handle}_copy` : ''

      this.handleSave({ chart, isClone: true })
    },

    redirect () {
      this.$router.push(this.previousPage || { name: 'admin.charts' })
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

    async fetchCustomColorSchemes () {
      return this.$SystemAPI.settingsList({ prefix: 'ui.charts.colorSchemes' })
        .then(settings => {
          const { value = [] } = settings[0] || {}
          this.customColorSchemes = value
        })
        .catch(this.toastErrorHandler(this.$t('notification:chart.colorScheme.fetch.failed')))
    },

    async saveColorScheme () {
      this.colorSchemeModal.processing = true

      const action = this.colorSchemeModal.edit ? 'update' : 'create'

      if (this.colorSchemeModal.edit) {
        const index = this.customColorSchemes.findIndex(({ id }) => id === this.colorSchemeModal.colorscheme.id)
        this.customColorSchemes.splice(index, 1, this.colorSchemeModal.colorscheme)
      } else {
        this.customColorSchemes.push(this.colorSchemeModal.colorscheme)
      }

      const values = [
        { name: 'ui.charts.colorSchemes', value: this.customColorSchemes },
      ]

      return this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.chart.config.colorScheme = this.colorSchemeModal.colorscheme.id
          this.closeColorSchemeModal()
          this.toastSuccess(this.$t(`notification:chart.colorScheme.${action}.success`))
          return this.$Settings.fetch().then(() => {
            return this.update()
          })
        })
        .catch(this.toastErrorHandler(this.$t(`notification:chart.colorScheme.${action}.failed`)))
        .finally(() => {
          this.colorSchemeModal.processing = false
        })
    },

    async deleteColorScheme () {
      this.colorSchemeModal.processing = true

      const value = this.customColorSchemes.filter(({ id }) => id !== this.colorSchemeModal.colorscheme.id)
      const values = [
        { name: 'ui.charts.colorSchemes', value },
      ]

      return this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.chart.config.colorScheme = undefined
          this.customColorSchemes = value
          this.closeColorSchemeModal()
          this.toastSuccess(this.$t('notification:chart.colorScheme.delete.success'))
          return this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:chart.colorScheme.delete.success')))
        .finally(() => {
          this.colorSchemeModal.processing = false
        })
    },

    createColorScheme () {
      this.colorSchemeModal.edit = false
      this.colorSchemeModal.colorscheme = {
        id: `custom-${Date.now()}`,
        name: '',
        colors: ['#6C757D', '#000000'],
      }

      this.colorSchemeModal.show = true
    },

    editColorScheme () {
      this.colorSchemeModal.edit = true
      this.colorSchemeModal.colorscheme = {
        id: this.currentColorScheme.id,
        name: this.currentColorScheme.name,
        colors: [...this.currentColorScheme.colors],
      }
      this.colorSchemeModal.show = true
    },

    closeColorSchemeModal () {
      this.colorSchemeModal.show = false
    },

    addColor () {
      this.colorSchemeModal.colorscheme.colors.push('#000000')
    },

    removeColor (index) {
      this.colorSchemeModal.colorscheme.colors.splice(index, 1)
    },

    getOptionKey ({ value }) {
      return value
    },

    checkUnsavedChart (next) {
      if (!this.chart.deletedAt) {
        return next(!isEqual(this.chart, this.initialChartState) ? window.confirm(this.$t('notification.unsavedChanges')) : true)
      }
      next()
    },
  },
}
</script>

<style lang="scss">

.chart-preview {
  max-height: 50%;
}

.color-box {
  width: 18px;
  height: 8px;
}
</style>
