<template>
  <div
    class="d-flex overflow-auto px-2 w-100"
  >
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <div
        class="d-inline-block mr-2"
      >
        <b-input-group
          size="sm"
        >
          <vue-select
            v-model="scenarios.selected"
            :options="scenarioOptions"
            :placeholder="$t('builder:pick-scenario')"
            class="h-100 bg-white rounded"
            @input="refreshReport()"
          />

          <b-input-group-append>
            <b-button
              variant="secondary"
              :disabled="!canUpdate"
              class="py-0"
              @click="openScenarioConfigurator"
            >
              <font-awesome-icon
                :icon="['fas', 'cog']"
              />
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </div>

      <b-button
        :disabled="!canUpdate"
        variant="secondary"
        size="sm"
        class="mr-1"
        @click="openDatasourceConfigurator"
      >
        {{ $t('builder:datasources.label') }}
      </b-button>

      <b-button-group
        size="sm"
        class="mr-1"
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
          variant="primary"
          style="margin-left:2px;"
          :disabled="!canUpdate"
          :to="reportEditor"
        >
          <font-awesome-icon
            :icon="['fas', 'pen']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <grid
      v-if="report && canRead &&showReport"
      :blocks.sync="reportBlocks"
      editable
    >
      <template
        slot-scope="{ block, index }"
      >
        <div
          class="h-100 editable-block"
        >
          <div
            class="add-element d-flex align-items-center justify-items-between mr-3 mt-3"
          >
            <b-button
              variant="link"
              class="text-light"
              @click="openDisplayElementSelector(index)"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="h4 mb-0"
              />
            </b-button>
            <b-button
              variant="link"
              class="text-light"
              @click="editBlock(index)"
            >
              <font-awesome-icon
                :icon="['far', 'edit']"
                class="h5 mb-0"
              />
            </b-button>

            <c-input-confirm
              size="md"
              variant="link text-danger"
              @confirmed="deleteBlock(index)"
            />
          </div>

          <block
            v-if="block"
            :index="index"
            :block="block"
            :scenario="currentSelectedScenario"
            :report-i-d="reportID"
          />
        </div>
      </template>
    </grid>

    <b-modal
      v-model="blocks.showConfigurator"
      :title="$t('builder:block.configuration')"
      :ok-title="$t('builder:save-button')"
      ok-variant="primary"
      cancel-variant="link"
      scrollable
      size="xl"
      body-class="p-0 border-top-0"
      header-class="pb-0 px-3 pt-3 border-bottom-0"
      @ok="updateBlock()"
    >
      <b-tabs
        v-if="currentBlock"
        active-nav-item-class="bg-grey"
        nav-wrapper-class="bg-white border-bottom"
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
            @add="openDisplayElementSelector(blocks.currentIndex)"
            @delete="deleteCurrentDisplayElement"
          >
            <template v-slot:label="{ item: { kind, name } }">
              {{ name || kind }}
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="grab text-grey"
              />
            </template>

            <template #configurator>
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
      size="xl"
      scrollable
      :ok-title="$t('builder:datasources.save')"
      ok-variant="primary"
      :title="$t('builder:datasources.label')"
      body-class="py-3"
      @ok="refreshReport()"
    >
      <configurator
        v-if="report"
        :items="reportDatasources"
        :current-index="datasources.currentIndex"
        draggable
        @select="setCurrentDatasource"
        @add="openDatasourceSelector()"
        @delete="deleteCurrentDataSource"
      >
        <template v-slot:label="{ item: { step } }">
          <span
            class="d-inline-block text-truncate"
          >
            {{ datasourceLabel(step, datasources.currentIndex) }}
          </span>
        </template>

        <template #configurator>
          <component
            :is="getDatasourceComponent(reportDatasources[datasources.currentIndex])"
            v-if="currentDatasourceStep"
            :datasources="reportDatasources"
            :step.sync="currentDatasourceStep"
          />
        </template>
      </configurator>
    </b-modal>

    <b-modal
      v-model="displayElements.showSelector"
      size="lg"
      scrollable
      hide-footer
      :title="$t('builder:add.display-element')"
      body-class="px-0 py-3"
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
      :title="$t('builder:scenarios.label')"
      body-class="py-3"
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
        <template v-slot:label="{ item: { label } }">
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
        @delete="handleDelete"
        @save="handleReportSave"
      >
        <b-button
          variant="light"
          size="lg"
          @click="createBlock"
        >
          <font-awesome-icon
            :icon="['fas', 'plus']"
            size="sm"
            class="mr-1"
          />
          {{ $t('general:label.add') }}
        </b-button>
      </editor-toolbar>
    </portal>
  </div>
</template>

<script>
import { reporter } from '@cortezaproject/corteza-js'
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
import VueSelect from 'vue-select'
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
    VueSelect,
    Prefilter,
  },

  mixins: [
    report,
  ],

  data () {
    return {
      processing: false,
      showReport: true,

      report: undefined,

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

        currentIndex: undefined,

        types: [
          {
            label: this.$t('builder:datasource.types.load.label'),
            kind: 'Load',
            value: this.$t('builder:datasource.types.load.loads-data-from-specified-resource-such-as-compose-records'),
          },
          {
            label: this.$t('builder:datasource.types.join.label'),
            kind: 'Join',
            value: this.$t('builder:datasource.types.join.joins-two-load-datasources-such-as-compose-record-selector'),
          },
          {
            label: this.$t('builder:datasource.types.group.label'),
            kind: 'Group',
            value: this.$t('builder:datasource.types.group.groups-data-from-load-datasource-like-counting-number-of-accounts-with-same-status'),
          },
        ],
      },

      scenarios: {
        showConfigurator: false,

        currentIndex: undefined,

        selected: undefined,
      },
    }
  },

  computed: {
    reportID () {
      return this.$route.params.reportID
    },

    pageTitle () {
      const title = this.report ? (this.report.meta.name || this.report.handle) : ''
      return `${this.$t('builder:report.builder')} - '${title}'` || this.$t('builder:report.builder')
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
        return this.datasources.currentIndex !== undefined ? this.reportDatasources[this.datasources.currentIndex].step : undefined
      },

      set (step) {
        if (this.datasources.currentIndex !== undefined) {
          this.reportDatasources[this.datasources.currentIndex].step = step
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
        return this.blocks.currentIndex !== undefined ? this.reportBlocks[this.blocks.currentIndex] : undefined
      },

      set (block) {
        if (this.blocks.currentIndex !== undefined) {
          this.reportBlocks[this.blocks.currentIndex] = block
        }
      },
    },

    currentDisplayElement: {
      get () {
        return this.displayElements.currentIndex !== undefined ? this.currentDisplayElements[this.displayElements.currentIndex] : undefined
      },

      set (element) {
        if (this.displayElements.currentIndex !== undefined) {
          this.currentDisplayElements[this.displayElements.currentIndex] = element
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
  },

  watch: {
    reportID: {
      immediate: true,
      handler (reportID) {
        this.scenarios.selected = undefined

        if (reportID) {
          this.processing = true

          this.fetchReport(this.reportID)
            .then(() => {
              this.mapBlocks()
            }).finally(() => {
              this.processing = false
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

      return `${currentIndex}`
    },

    openDatasourceSelector () {
      this.datasources.showSelector = true
      this.datasources.currentIndex = this.reportDatasources.length ? 0 : undefined
    },

    openDatasourceConfigurator () {
      this.datasources.showConfigurator = true
      this.datasources.currentIndex = this.reportDatasources.length ? 0 : undefined
    },

    setCurrentDatasource (index) {
      this.datasources.currentIndex = index
    },

    deleteCurrentDataSource () {
      this.reportDatasources.splice(this.datasources.currentIndex, 1)
      this.datasources.currentIndex = this.reportDatasources.length ? 0 : undefined
    },

    addDatasource (kind = '') {
      if (kind) {
        let step

        switch (kind) {
          case 'Group':
            step = reporter.StepFactory({
              group: {
                name: 'Group',
                keys: [],
                columns: [],
                filter: {},
                sort: '',
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

        this.reportDatasources.push({
          step,
          meta: {},
        })
      }

      // Select newly added datasource in configurator
      this.datasources.currentIndex = this.reportDatasources.length - 1

      // Close selector, open configurator
      this.datasources.showSelector = false
      this.datasources.showConfigurator = true
    },

    // Blocks
    handleReportSave () {
      this.report.blocks = this.reportBlocks.map(({ moved, x, y, w, h, i, ...p }) => {
        return { ...p, key: `${i}`, xywh: [x, y, w, h] }
      })

      this.handleSave()
        .then(() => {
          this.mapBlocks()
          this.refreshReport()
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
      if (this.currentBlock) {
        const elements = this.currentBlock.elements

        this.reportBlocks.splice(this.blocks.currentIndex, 1, { ...this.currentBlock, elements: [] })
        setTimeout(() => {
          this.reportBlocks.splice(this.blocks.currentIndex, 1, { ...this.currentBlock, elements })
        }, 50)
      }
    },

    editBlock (index = undefined) {
      this.blocks.currentIndex = index
      this.setCurrentDisplayElement(this.reportBlocks[this.blocks.currentIndex].elements.length ? 0 : undefined)
      this.blocks.showConfigurator = true
    },

    deleteBlock (index = undefined) {
      this.reindexBlocks(this.reportBlocks.filter((p, i) => index !== i))
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
  },
}
</script>

<style lang="scss">
.add-element {
  position: absolute;
  background-color: #1e2224;
  bottom: 0;
  left: 0;
  z-index: 1021;
  opacity: 0.5;
  border-top-right-radius: 10px;

  &:hover {
    opacity: 1;
  }
}
</style>
