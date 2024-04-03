<template>
  <div
    class="d-flex overflow-auto p-2 w-100"
  >
    <portal
      v-if="!fetchingReport"
      to="topbar-title"
    >
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <b-input-group style="max-width: 300px;">
        <c-input-select
          v-model="scenarios.selected"
          :options="scenarioOptions"
          :get-option-key="getOptionKey"
          :placeholder="$t('builder:pick-scenario')"
          :disabled="processing || fetchingReport"
          size="sm"
          @input="refreshReport()"
        />
        <b-input-group-append>
          <b-button
            v-b-tooltip.noninteractive.hover="{ title: $t('builder:tooltip.configure-scenarios'), container: '#body' }"
            variant="extra-light"
            :disabled="!canUpdate"
            size="sm"
            @click="openScenarioConfigurator"
          >
            <font-awesome-icon
              :icon="['fas', 'cog']"
              class="text-primary"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>

      <b-button
        :disabled="!canUpdate"
        variant="extra-light"
        size="sm"
        @click="openDatasourceConfigurator"
      >
        {{ $t('builder:datasources.label') }}
      </b-button>
      <b-button-group
        size="sm"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center justify-content-center"
          :disabled="!canRead"
          :to="reportViewer"
        >
          {{ $t('builder:report.view') }}
          <font-awesome-icon
            class="ml-2"
            :icon="['far', 'eye']"
          />
        </b-button>
        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('builder:tooltip.edit.report'), container: '#body' }"
          variant="primary"
          class="d-flex align-items-center justify-content-center"
          style="margin-left:2px;"
          :disabled="!canUpdate"
          :to="reportEditor"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      v-if="fetchingReport"
      class="d-flex align-items-center justify-content-center w-100 h-100"
    >
      <b-spinner />
    </div>

    <grid
      v-if="report && canRead && showReport && !fetchingReport"
      :blocks.sync="reportBlocks"
      editable
      @item-updated="onBlockUpdated"
    >
      <template
        slot-scope="{ block, index }"
      >
        <div
          class="h-100"
        >
          <div
            class="toolbox border-0 p-2 m-0 text-light text-center"
          >
            <div
              v-if="unsavedBlocks.has(index)"
              v-b-tooltip.noninteractive.hover="{ title: $t('builder:tooltip.unsavedChanges'), container: '#body' }"
              class="btn border-0"
            >
              <font-awesome-icon
                :icon="['fas', 'exclamation-triangle']"
                class="text-warning"
              />
            </div>
            <b-button-group>
              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('builder:tooltip.add.displayElement'), container: '#body' }"
                variant="outline-light"
                class="border-0"
                @click="openDisplayElementSelector(index)"
              >
                <font-awesome-icon
                  :icon="['fas', 'plus']"
                />
              </b-button>
              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('builder:tooltip.edit.block'), container: '#body' }"
                variant="outline-light"
                class="border-0"
                @click="editBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>
            </b-button-group>
            <c-input-confirm
              :tooltip="$t('builder:tooltip.delete.block')"
              show-icon
              size="md"
              class="ml-1"
              @confirmed="deleteBlock(index)"
            />
          </div>
          <block
            v-if="block"
            :index="index"
            :block="block"
            :scenario="currentSelectedScenario"
            :report-i-d="reportID"
            @item-updated="onBlockUpdated"
          />
        </div>
      </template>
    </grid>
    <b-modal
      :title="$t('builder:block.configuration')"
      :ok-title="$t('builder:save-button')"
      :visible="showEditor"
      ok-variant="primary"
      cancel-variant="light"
      scrollable
      size="xl"
      body-class="p-0 border-top-0"
      header-class="border-bottom-0"
      no-fade
      @hide="hideEditorModal"
      @ok="updateEditorBlock()"
    >
      <b-tabs
        v-if="currentBlock"
        active-tab-class="tab-content h-auto overflow-auto"
        card
      >
        <b-tab
          :title="$t('builder:general')"
          active
        >
          <b-form-group
            :label="$t('builder:title')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="currentBlock.title"
              type="text"
              :placeholder="$t('builder:block.title')"
            />
          </b-form-group>
          <b-form-group
            :label="$t('builder:description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="currentBlock.description"
              :placeholder="$t('builder:block.description')"
            />
          </b-form-group>
          <b-form-group
            :label="$t('builder:layout')"
            label-class="text-primary"
          >
            <b-form-radio-group
              v-model="currentBlock.layout"
              :options="blockLayoutOptions"
              buttons
              button-variant="outline-primary"
            />
          </b-form-group>
        </b-tab>
        <b-tab
          :active="!!currentBlock.elements.length"
          :title="$t('builder:elements')"
        >
          <configurator
            :items="currentDisplayElements"
            :current-index="displayElements.currentIndex"
            draggable
            @select="setCurrentDisplayElement"
            @add="openDisplayElementSelector(editor.currentIndex)"
            @delete="deleteCurrentDisplayElement"
          >
            <template #label="{ item: { kind, name } }">
              {{ name || kind }}
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="text-secondary grab"
              />
            </template>            <template #configurator>
              <display-element-configurator
                v-if="currentDisplayElement"
                :display-element.sync="currentDisplayElement"
                :block="currentBlock"
                :datasources="reportDatasources"
                class="pr-2"
              />
            </template>
          </configurator>
        </b-tab>
      </b-tabs>
    </b-modal>
    <b-modal
      v-model="datasources.showConfigurator"
      :title="$t('builder:datasources.label')"
      cancel-variant="light"
      :cancel-disabled="datasources.processing"
      scrollable
      size="xl"
      body-class="py-3"
      no-fade
      @hide="hideDatasourceConfigurator"
    >
      <configurator
        v-if="report"
        :items="datasources.tempItems"
        :current-index="datasources.currentIndex"
        draggable
        @select="setCurrentDatasource"
        @add="openDatasourceSelector()"
        @delete="deleteCurrentDataSource"
      >
        <template #label="{ item: { step } }">
          <span
            class="d-inline-block text-truncate"
          >
            {{ datasourceLabel(step, datasources.currentIndex) }}
          </span>
        </template>
        <template #configurator>
          <component
            :is="getDatasourceComponent(datasources.tempItems[datasources.currentIndex])"
            v-if="currentDatasourceStep"
            :index="datasources.currentIndex"
            :datasources="datasources.tempItems"
            :step.sync="currentDatasourceStep"
            :creating="datasources.tempItems[datasources.currentIndex].meta.creating"
          />
        </template>
      </configurator>
      <template #modal-footer>
        <c-button-submit
          data-test-id="button-save"
          :disabled="datasourceSaveDisabled"
          :processing="datasources.processing"
          :text="$t('general:label.saveAndClose')"
          @submit="saveDatasources"
        />
      </template>
    </b-modal>
    <b-modal
      v-model="displayElements.showSelector"
      size="lg"
      scrollable
      hide-footer
      :title="$t('builder:add.display-element')"
      body-class="px-0 py-3"
      no-fade
    >
      <selector
        :items="displayElements.types"
        @select="addDisplayElement"
      />
    </b-modal>
    <b-modal
      v-model="datasources.showSelector"
      size="lg"
      scrollable
      hide-footer
      :title="$t('builder:add.datasource')"
      body-class="px-0 py-3"
      no-fade
    >
      <selector
        :items="datasources.types"
        display-mode="text"
        @select="addDatasource"
      />
    </b-modal>
    <b-modal
      v-model="scenarios.showConfigurator"
      size="xl"
      scrollable
      :ok-title="$t('builder:scenarios.save')"
      ok-variant="primary"
      cancel-variant="light"
      :title="$t('builder:scenarios.label')"
      body-class="py-3"
      no-fade
    >
      <configurator
        v-if="report"
        :items="reportScenarios"
        :current-index="scenarios.currentIndex"
        draggable
        @select="setCurrentScenario"
        @add="addScenario()"
        @delete="deleteCurrentScenario()"
      >
        <template #label="{ item: { label } }">
          <span
            class="d-inline-block text-truncate"
          >
            {{ label }}
          </span>
        </template>
        <template #configurator>
          <scenario-configurator
            v-if="currentScenario"
            :current-index="scenarios.currentIndex"
            :datasources="reportDatasources"
            :scenario.sync="currentScenario"
          />
        </template>
      </configurator>
    </b-modal>
    <portal to="report-toolbar">
      <editor-toolbar
        :back-link="{ name: 'report.list' }"
        :delete-disabled="!canDelete"
        :save-disabled="!canUpdate"
        :processing="processing"
        :processing-save="processingSave"
        :processing-delete="processingDelete"
        :processing-clone="processingClone"
        @clone="handleReportCloning"
        @delete="handleDelete"
        @save="handleReportSave"
      >
        <b-button
          variant="light"
          size="lg"
          :disabled="processing"
          @click="createBlock"
        >
          <font-awesome-icon
            :icon="['fas', 'plus']"
            size="sm"
          />
          {{ $t('general:label.add') }}
        </b-button>
      </editor-toolbar>
    </portal>
  </div>
</template>

<script>
import { cloneDeep } from 'lodash'
import { system, reporter } from '@cortezaproject/corteza-js'
import report from 'corteza-webapp-reporter/src/mixins/report'
import Grid from 'corteza-webapp-reporter/src/components/Report/Grid'
import Block from 'corteza-webapp-reporter/src/components/Report/Blocks'
import datasources from 'corteza-webapp-reporter/src/components/Report/Datasources/loader'
import Configurator from 'corteza-webapp-reporter/src/components/Common/Configurator'
import Selector from 'corteza-webapp-reporter/src/components/Common/Selector'
import EditorToolbar from 'corteza-webapp-reporter/src/components/EditorToolbar'
import DisplayElementConfigurator from 'corteza-webapp-reporter/src/components/Report/Blocks/DisplayElements/Configurators'
import ScenarioConfigurator from 'corteza-webapp-reporter/src/components/Report/Scenarios'
import * as displayElementThumbnails from 'corteza-webapp-reporter/src/assets/DisplayElements'
import Prefilter from 'corteza-webapp-reporter/src/components/Common/Prefilter'

export default {
  name: 'ReportBuilder',

  components: {
    Grid,
    Selector,
    Configurator,
    Block,
    DisplayElementConfigurator,
    ScenarioConfigurator,
    EditorToolbar,
    Prefilter,
  },

  mixins: [
    report,
  ],

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedBlocks(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedBlocks(next)
  },

  data () {
    return {
      processing: false,
      processingSave: false,
      processingDelete: false,
      processingClone: false,
      fetchingReport: false,

      showReport: true,

      report: undefined,

      unsavedBlocks: new Set(),

      dataframes: [],

      blocks: {
        showConfigurator: false,

        currentIndex: undefined,

        items: [],
      },

      displayElements: {
        showSelector: false,

        currentIndex: undefined,

        types: [
          {
            label: this.$t('builder:display-elements.types.text'),
            kind: 'Text',
            value: displayElementThumbnails.Text,
          },
          {
            label: this.$t('builder:display-elements.types.metric'),
            kind: 'Metric',
            value: displayElementThumbnails.Metric,
          },
          {
            label: this.$t('builder:display-elements.types.table'),
            kind: 'Table',
            value: displayElementThumbnails.Table,
          },
          {
            label: this.$t('builder:display-elements.types.chart'),
            kind: 'Chart',
            value: displayElementThumbnails.Chart,
          },
        ],
      },

      datasources: {
        showSelector: false,
        showConfigurator: false,

        processing: false,
        currentIndex: undefined,
        tempItems: [],

        types: [
          {
            label: this.$t('builder:datasource.types.load.label'),
            kind: 'Load',
            value: this.$t('builder:datasource.types.load.data-from-resource'),
          },
          {
            label: this.$t('builder:datasource.types.link.label'),
            kind: 'Link',
            value: this.$t('builder:datasource.types.link.load-datasources'),
          },
          {
            label: this.$t('builder:datasource.types.join.label'),
            kind: 'Join',
            value: this.$t('builder:datasource.types.join.load-datasources'),
          },
          {
            label: this.$t('builder:datasource.types.aggregate.label'),
            kind: 'Aggregate',
            value: this.$t('builder:datasource.types.aggregate.load-datasource'),
          },
        ],
      },

      scenarios: {
        showConfigurator: false,

        currentIndex: undefined,

        selected: undefined,
      },

      editor: undefined,
    }
  },

  computed: {
    reportID () {
      return this.$route.params.reportID
    },

    pageTitle () {
      const title = this.report ? (this.report.meta.name || this.report.handle) : ''
      return `${this.$t('builder:report.builder')} - "${title}"` || this.$t('builder:report.builder')
    },

    canRead () {
      return this.report ? this.report.canReadReport : false
    },

    canDelete () {
      return this.report ? this.report.canDeleteReport : false
    },

    canUpdate () {
      return this.report ? this.report.canUpdateReport : false
    },

    currentDisplayElements () {
      return this.currentBlock ? this.currentBlock.elements : []
    },

    reportDatasources: {
      get () {
        return this.report ? this.report.sources : []
      },

      set (sources) {
        this.report.sources = sources
      },
    },

    currentDatasourceStep: {
      get () {
        return this.datasources.currentIndex !== undefined ? this.datasources.tempItems[this.datasources.currentIndex].step : undefined
      },

      set (step) {
        if (this.datasources.currentIndex !== undefined) {
          this.datasources.tempItems[this.datasources.currentIndex].step = step
        }
      },
    },

    reportBlocks: {
      get () {
        return this.blocks.items || []
      },

      set (blocks) {
        this.blocks.items = blocks
      },
    },

    currentBlock: {
      get () {
        return this.editor ? this.editor.block : undefined
      },

      set (block) {
        if (this.editor && this.editor.currentIndex !== undefined) {
          this.editor.block = block
        }
      },
    },

    currentDisplayElement: {
      get () {
        return this.displayElements.currentIndex !== undefined ? this.currentBlock.elements[this.displayElements.currentIndex] : undefined
      },

      set (element) {
        if (this.displayElements.currentIndex !== undefined) {
          this.currentBlock.elements.splice(this.displayElements.currentIndex, 1, element)
        }
      },
    },

    reportScenarios: {
      get () {
        return this.report ? this.report.scenarios : []
      },

      set (scenarios) {
        this.report.scenarios = scenarios
      },
    },

    currentScenario: {
      get () {
        return this.scenarios.currentIndex !== undefined ? this.reportScenarios[this.scenarios.currentIndex] : undefined
      },

      set (scenario) {
        if (this.scenarios.currentIndex !== undefined) {
          this.reportScenarios[this.scenarios.currentIndex] = scenario
        }
      },
    },

    currentSelectedScenario () {
      return this.scenarios.selected ? this.reportScenarios.find(({ label }) => label === this.scenarios.selected) : undefined
    },

    scenarioOptions () {
      return this.reportScenarios.map(({ label }) => label)
    },

    reportViewer () {
      return this.report ? { name: 'report.view', params: { reportID: this.report.reportID } } : undefined
    },

    reportEditor () {
      return this.report ? { name: 'report.edit', params: { reportID: this.report.reportID } } : undefined
    },

    blockLayoutOptions () {
      return [
        { text: this.$t('builder:layout-options.horizontal'), value: 'horizontal' },
        { text: this.$t('builder:layout-options.vertical'), value: 'vertical' },
      ]
    },

    datasourceSaveDisabled () {
      const uniqueDatasources = new Set()
      const hasDuplicates = this.datasources.tempItems.some(({ step }) => {
        const name = step[Object.keys(step)].name
        return !name || uniqueDatasources.size === uniqueDatasources.add(name).size
      })

      return this.datasources.processing || hasDuplicates
    },

    showEditor () {
      return this.editor && this.editor.currentIndex !== undefined
    },
  },

  watch: {
    reportID: {
      immediate: true,
      handler (reportID) {
        this.unsavedBlocks.clear()
        this.scenarios.selected = undefined
        this.reportBlocks = []
        this.report = undefined
        if (reportID) {
          this.processing = true
          this.fetchingReport = true

          this.fetchReport(this.reportID)
            .then(() => {
              this.mapBlocks()
            }).catch(() => {
              this.toastErrorHandler(this.$t('notification:report.loadFailed'))
            })
            .finally(() => {
              setTimeout(() => {
                this.fetchingReport = false
                this.processing = false
              }, 400)
            })
        }
      },
    },
  },

  methods: {
    refreshReport () {
      this.showReport = false
      return setTimeout(() => {
        this.showReport = true
      }, 50)
    },

    // If block is added/reordered or deleted, vue-grid-layout needs fresh indexes to work properly
    reindexBlocks (blocks = this.reportBlocks || []) {
      this.reportBlocks = blocks.map((block, i) => {
        return { ...block, i }
      })
    },

    // Datasources
    getDatasourceComponent ({ step }) {
      if (step) {
        for (const s in step) {
          return datasources(s)
        }
      }

      return undefined
    },

    datasourceLabel (datasource, currentIndex) {
      for (const v of Object.values(datasource)) {
        if (v && v.name) {
          return v.name
        }
      }

      return `${this.$t('datasources:source')} ${currentIndex}`
    },

    openDatasourceSelector () {
      this.datasources.showSelector = true
      this.datasources.currentIndex = this.datasources.tempItems.length ? 0 : undefined
    },

    openDatasourceConfigurator () {
      this.datasources.showConfigurator = true
      this.datasources.tempItems = cloneDeep(this.reportDatasources).map(ds => {
        ds.meta.creating = false
        return ds
      })
      this.datasources.currentIndex = this.datasources.tempItems.length ? 0 : undefined
    },

    hideDatasourceConfigurator () {
      this.datasources.showConfigurator = false
      this.datasources.tempItems = []
      this.datasources.currentIndex = undefined
    },

    setCurrentDatasource (index) {
      this.datasources.currentIndex = index
    },

    deleteCurrentDataSource () {
      this.datasources.tempItems.splice(this.datasources.currentIndex, 1)
      this.datasources.currentIndex = this.datasources.tempItems.length ? 0 : undefined
    },

    addDatasource (kind = '') {
      if (kind) {
        let step

        switch (kind) {
          case 'Aggregate':
            step = reporter.StepFactory({
              aggregate: {
                name: 'Aggregate',
                keys: [],
                columns: [],
                filter: {},
                sort: '',
              },
            })
            break

          case 'Link':
            step = reporter.StepFactory({
              link: {
                name: 'Link',
                foreignColumn: '',
                foreignSource: '',
                localColumn: '',
                localSource: '',
              },
            })
            break

          case 'Join':
            step = reporter.StepFactory({
              join: {
                name: 'Join',
                foreignColumn: '',
                foreignSource: '',
                localColumn: '',
                localSource: '',
              },
            })
            break

          default:
            step = reporter.StepFactory({
              load: {
                name: 'Load',
                source: 'composeRecords',
                definition: {},
                filter: {},
                sort: '',
              },
            })
        }

        this.datasources.tempItems.push({
          step,
          meta: {},
        })
      }

      // Select newly added datasource in configurator
      this.datasources.currentIndex = this.datasources.tempItems.length - 1

      // Close selector, open configurator
      this.datasources.showSelector = false
      this.datasources.showConfigurator = true
    },

    saveDatasources () {
      // Prevent closing of modal and manually close it when request is complete
      this.datasources.processing = true

      const sources = this.datasources.tempItems
      const { reportID } = this.report

      // Fetch saved report and merge with datasources
      return this.$SystemAPI.reportRead({ reportID }).then(report => {
        return this.$SystemAPI.reportUpdate(new system.Report({ ...report, sources }))
      }).then(report => {
        report.scenarios = this.report.scenarios
        this.report = new system.Report(report)
        this.refreshReport()
        this.hideDatasourceConfigurator()
        this.toastSuccess(this.$t('notification:report.datasources.updated'))
      }).catch(this.toastErrorHandler(this.$t('notification:report.datasources.updateFailed')))
        .finally(() => {
          this.datasources.processing = false
        })
    },

    // Blocks
    handleReportSave () {
      this.processingSave = true

      this.report.blocks = this.reportBlocks.map(({ moved, x, y, w, h, i, ...p }) => {
        return { ...p, key: `${i}`, xywh: [x, y, w, h] }
      })

      this.handleSave()
        .then(() => {
          this.mapBlocks()
          this.refreshReport()
          this.unsavedBlocks.clear()
        })
        .finally(() => {
          this.processingSave = false
        })
    },

    handleReportCloning () {
      this.handleClone(this.report).then(({ reportID }) => {
        this.$root.$emit('refetch:reports')
        this.$router.push({ name: 'report.builder', params: { reportID } })
      })
    },

    mapBlocks () {
      this.reportBlocks = this.report.blocks.map(({ xywh, ...p }, i) => {
        const [x, y, w, h] = xywh
        return { ...p, x, y, w, h, i }
      })
    },

    createBlock () {
      let newBlock = {
        ...new reporter.Block(),
      }

      const [x, y, w, h] = newBlock.xywh
      newBlock = {
        ...newBlock,
        x,
        y,
        w,
        h,
      }

      this.reindexBlocks([...this.reportBlocks, newBlock])
    },

    updateBlock () {
      this.unsavedBlocks.add(this.blocks.currentIndex)

      if (this.currentBlock) {
        const elements = this.currentBlock.elements

        this.reportBlocks.splice(this.blocks.currentIndex, 1, { ...this.currentBlock, elements: [] })
        setTimeout(() => {
          this.reportBlocks.splice(this.blocks.currentIndex, 1, { ...this.currentBlock, elements })
        }, 50)
      }
    },

    updateEditorBlock (block = this.editor.block) {
      const { currentIndex } = this.editor
      this.reportBlocks[currentIndex] = block
      this.editor = undefined
      this.onBlockUpdated(currentIndex)
      this.refreshReport()
    },

    editBlock (index = undefined) {
      const { x, y, w, h, i } = this.reportBlocks[index]
      const block = new reporter.Block(this.reportBlocks[index])

      block.x = x
      block.y = y
      block.w = w
      block.h = h
      block.i = i

      this.editor = {
        currentIndex: index,
        block,
      }
      this.setCurrentDisplayElement(this.editor.block.elements.length ? 0 : undefined)
    },

    deleteBlock (index = undefined) {
      this.reindexBlocks(this.reportBlocks.filter((p, i) => index !== i))
      this.unsavedBlocks.add(index)
    },

    hideEditorModal () {
      this.editor = undefined
      this.displayElements.currentIndex = undefined
    },

    // Display elements
    openDisplayElementSelector (index) {
      this.blocks.currentIndex = index
      this.displayElements.showSelector = true
    },

    setCurrentDisplayElement (index) {
      this.displayElements.currentIndex = index
    },

    deleteCurrentDisplayElement () {
      this.currentBlock.elements.splice(this.displayElements.currentIndex, 1)
      this.displayElements.currentIndex = this.currentBlock.elements.length ? 0 : undefined
      this.setCurrentDisplayElement(this.displayElements.currentIndex)
    },

    addDisplayElement (kind) {
      const newDisplayElement = reporter.DisplayElementMaker({ kind })

      this.reportBlocks[this.blocks.currentIndex].elements.push(newDisplayElement)

      this.displayElements.showSelector = false

      this.editBlock(this.blocks.currentIndex)
      this.setCurrentDisplayElement(this.currentBlock.elements.length - 1)

      this.updateBlock()
    },

    // Scenarios
    openScenarioConfigurator () {
      this.scenarios.showConfigurator = true

      if (this.reportScenarios.length) {
        this.setCurrentScenario(0)
      }
    },

    setCurrentScenario (index = -1) {
      this.scenarios.currentIndex = this.reportScenarios.length && index >= 0 ? index : undefined
    },

    addScenario () {
      if (!this.reportScenarios) {
        this.reportScenarios = []
      }

      this.reportScenarios.push({
        label: 'Scenario Name',
        filters: {},
      })

      this.setCurrentScenario(this.reportScenarios.length - 1)
    },

    deleteCurrentScenario () {
      this.reportScenarios.splice(this.scenarios.currentIndex, 1)
      this.scenarios.currentIndex = this.reportScenarios.length ? 0 : undefined
      this.setCurrentScenario(this.scenarios.currentIndex)
    },

    getOptionKey (scenario) {
      return scenario
    },

    // Trigger browser dialog on page leave to prevent unsaved changes
    checkUnsavedBlocks (next) {
      if (this.report.deletedAt) {
        return next(true)
      }

      next(!this.unsavedBlocks.size || window.confirm(this.$t('builder:unsaved-changes')))
    },

    onBlockUpdated (index) {
      this.unsavedBlocks.add(index)
    },
  },
}
</script>

<style lang="scss">
div.toolbox {
  position: absolute;
  background-color: var(--secondary);
  bottom: 0;
  left: 0;
  z-index: 1001;
  border-top-right-radius: 10px;
  opacity: 0.5;
  pointer-events: none;

  &:hover {
    opacity: 1;
  }

  & * {
    pointer-events: auto;
  }
}

[dir="rtl"] {
  div.toolbox {
    left: 0;
    right: auto;
  }
}
</style>
