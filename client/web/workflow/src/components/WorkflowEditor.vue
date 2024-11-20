<template>
  <div
    id="editor"
    ref="editor"
    class="d-flex w-100 h-100"
  >
    <portal to="topbar-title">
      {{ workflow.meta.name || workflow.handle }}
    </portal>

    <portal to="topbar-tools">
      <b-button
        v-b-modal.workflow
        data-test-id="button-configure-workflow"
        variant="primary"
        size="sm"
        class="d-flex align-items-center"
      >
        {{ $t('configurator:configuration') }}
        <font-awesome-icon
          :icon="['fas', 'cog']"
          class="ml-1"
        />
      </b-button>
    </portal>

    <div class="toolbar d-flex flex-column h-100 border-right shadow-lg">
      <div
        id="toolbar"
        ref="toolbar"
        class="d-flex flex-column align-items-center mt-1 overflow-auto"
      />

      <div
        class="d-flex flex-grow-1 align-items-end justify-content-center py-3"
      >
        <b-button
          ref="help"
          v-b-modal.help
          variant="outline-light"
          class="d-flex align-items-center border-0 p-2"
        >
          <font-awesome-icon
            :icon="['far', 'question-circle']"
            class="h4 mb-0 text-primary"
          />
        </b-button>
      </div>
    </div>

    <div
      ref="tooltips"
      class="mh-100"
    />

    <b-card
      no-header
      class="w-100 h-100 border-0 shadow-sm rounded-0"
      body-class="p-0"
    >
      <div
        v-if="workflow.meta"
        class="position-absolute pl-2 pt-2"
        style="z-index: 1;"
      >
        <p
          v-if="workflow.meta.description"
          :class="{ 'mb-2': getRunAs }"
          class="mb-0 text-truncate"
          style="white-space: pre-line; max-height: 48px;"
        >
          {{ workflow.meta.description }}
        </p>

        <p
          v-if="getRunAs"
          class="mb-0 text-truncate"
        >
          <b>{{ $t('editor:run-as') }}</b> <samp>{{ getRunAs }}</samp>
        </p>

        <div
          class="d-flex align-items-center mb-1"
        >
          <h5
            v-if="workflow.deletedAt"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="danger"
            >
              {{ $t('editor:deleted') }}
            </b-badge>
          </h5>

          <h5
            v-if="!workflow.enabled"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="danger"
            >
              {{ $t('editor:disabled') }}
            </b-badge>
          </h5>

          <h5
            v-if="hasIssues"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="danger"
            >
              {{ $t('editor:detected-issues') }}
            </b-badge>
          </h5>

          <h5
            v-if="workflow.meta.subWorkflow"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="info"
            >
              {{ $t('general:subworkflow') }}
            </b-badge>
          </h5>

          <h5
            v-if="deferred"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="info"
            >
              {{ $t('editor:deferred') }}
            </b-badge>
          </h5>

          <h5
            v-if="triggersPathsChanged"
            class="mb-0 mr-1"
          >
            <b-badge
              variant="warning"
            >
              {{ $t('notification:trigger-paths-changed') }}
            </b-badge>
          </h5>
        </div>
      </div>

      <div
        class="bg-white position-absolute m-2 zoom border border-secondary"
        style="z-index: 1; width: fit-content;"
      >
        <div
          class="d-flex align-items-baseline p-2"
        >
          {{ getZoomPercent }}
          <b-button
            variant="link"
            class="ml-4 p-0"
            @click="zoom(false)"
          >
            <font-awesome-icon
              :icon="['fas', 'search-minus']"
            />
          </b-button>
          <b-button
            variant="link"
            class="ml-1 p-0"
            @click="zoom()"
          >
            <font-awesome-icon
              :icon="['fas', 'search-plus']"
              class="pointer"
              @click="zoom()"
            />
          </b-button>
          <b-button
            variant="link"
            class="ml-2 p-0 text-decoration-none"
            @click="resetZoom()"
          >
            {{ $t('editor:reset') }}
          </b-button>
        </div>
      </div>

      <div
        class="d-flex flex-column flex-shrink position-absolute fixed-bottom m-2"
        style="z-index: 1; width: 20vw;"
      >
        <c-button-submit
          v-if="changeDetected && canUpdateWorkflow"
          data-test-id="button-save-workflow"
          variant="primary"
          block
          :processing="processingSave"
          :text="$t('editor:detected-changes') + `${canUpdateWorkflow ? $t('editor:click-to-save') : ''}`"
          :loading-text="$t('editor:saving')"
          class="rounded-0 py-2 px-3"
          @submit="saveWorkflow()"
        />
      </div>

      <div
        id="graph"
        ref="graph"
        class="h-100 p-0"
      />
    </b-card>

    <!--
      no-enforce-focus flag doesn't set focus to sidebar when it is opened.
      Bad for Accessability, since keyboard only users can't use sidebar.
    -->
    <b-sidebar
      v-model="sidebar.show"
      shadow
      right
      lazy
      no-enforce-focus
      width="500px"
      :no-header="!sidebar.item"
      header-class="bg-white border-bottom border-light p-2"
      body-class="bg-white"
      footer-class="bg-white border-top border-light p-1"
    >
      <template
        #header
      >
        <div
          class="d-flex align-items-center w-100 h5 mb-0 p-2"
        >
          <b-img
            v-if="getSidebarItemIcon"
            :src="getSidebarItemIcon"
            class="mr-2"
          />
          <h4
            class="text-primary font-weight-bold mb-0"
          >
            <b>{{ getSidebarItemType }}</b>
          </h4>

          <div
            class="ml-auto"
          >
            {{ $t('editor:id') }} <var>{{ getSelectedItem.node.id }}</var>
          </div>
        </div>
      </template>

      <transition
        name="component-fade"
        mode="out-in"
      >
        <configurator
          v-if="sidebar.showItem"
          :item.sync="sidebar.item"
          :edges.sync="edges"
          :out-edges="sidebar.outEdges"
          :is-subworkflow="!!workflow.meta.subWorkflow"
          @update-value="setValue($event)"
          @update-default-value="setValue($event, true)"
        />
      </transition>

      <template
        #footer
      >
        <div
          class="d-flex m-2"
        >
          <c-input-confirm
            size="md"
            size-confirm="md"
            variant="danger"
            :processing="processingDelete"
            :text="$t('editor:delete')"
            @confirmed="sidebarDelete()"
          />

          <div
            class="ml-auto"
          >
            <portal-target name="sidebar-footer" />
          </div>
        </div>
      </template>
    </b-sidebar>

    <b-modal
      id="workflow"
      :title="$t('editor:workflow-configuration')"
      size="lg"
      :hide-header-close="workflow.workflowID === '0'"
      :no-close-on-backdrop="workflow.workflowID === '0'"
      :no-close-on-esc="workflow.workflowID === '0'"
      no-fade
    >
      <template #modal-title>
        {{ $t('editor:workflow-configuration') }}
      </template>

      <div
        v-if="workflow.workflowID && workflow.workflowID !== '0'"
        class="d-flex mb-3"
      >
        <import
          data-test-id="button-import-workflow"
          :disabled="importProcessing"
          @import="importJSON"
        />

        <export
          data-test-id="button-export-workflow"
          :workflows="[workflow.workflowID]"
          :file-name="workflow.meta.name || workflow.handle"
          size="lg"
          class="ml-1"
        />

        <c-permissions-button
          v-if="workflow.canGrant"
          :title="workflow.meta.name || workflow.handle || workflow.workflowID"
          :target="workflow.meta.name || workflow.handle || workflow.workflowID"
          :resource="`corteza::automation:workflow/${workflow.workflowID}`"
          :button-label="$t('general:permissions')"
          class="btn-lg ml-1"
        />
      </div>

      <workflow-configurator
        v-if="workflow.workflowID"
        :workflow="workflow"
        @delete="$emit('delete')"
      />

      <template #modal-footer>
        <div
          class="d-flex w-100"
        >
          <c-input-confirm
            v-if="workflow.canDeleteWorkflow && !isDeleted"
            size="md"
            size-confirm="md"
            :processing="processingDelete"
            :text="$t('editor:delete')"
            :borderless="false"
            @confirmed="$emit('delete')"
          />

          <c-input-confirm
            v-else-if="isDeleted"
            size="md"
            size-confirm="md"
            :processing="processingDelete"
            :text="$t('editor:undelete')"
            :borderless="false"
            @confirmed="$emit('undelete')"
          />

          <b-button
            v-if="workflow.workflowID === '0'"
            variant="light"
            @click="$router.back()"
          >
            {{ $t('editor:back') }}
          </b-button>

          <c-button-submit
            data-test-id="button-save-workflow"
            :disabled="canSave"
            :processing="processingSave"
            :text="$t('editor:save')"
            class="ml-auto"
            @submit="saveWorkflow()"
          />
        </div>
      </template>
    </b-modal>

    <b-modal
      id="help"
      :title="$t('editor:help')"
      size="lg"
      scrollable
      hide-footer
      no-fade
      body-class="p-0"
    >
      <help />
    </b-modal>

    <b-modal
      id="issues"
      v-model="issuesModal.show"
      :title="$t('editor:issues')"
      hide-footer
      no-fade
    >
      <div
        v-for="(issue, index) in issuesModal.issues"
        :key="index"
      >
        <p>
          <code>{{ issue[0].toUpperCase() + issue.slice(1) }}</code>
        </p>
      </div>
    </b-modal>

    <b-modal
      id="dry-run"
      v-model="dryRun.show"
      size="lg"
      :title="$t('editor:initial-scope')"
      scrollable
      :body-class="dryRun.lookup ? '' : 'p-1'"
      :ok-only="dryRun.lookup"
      :ok-title="`${dryRun.lookup ? $t('editor:load-and-configure') : $t('editor:run-workflow')}`"
      :cancel-title="$t('editor:back')"
      ok-variant="success"
      no-fade
      @cancel.prevent="dryRun.lookup = true"
      @ok="dryRunOk"
    >
      <div
        v-if="dryRun.lookup"
      >
        <small>
          {{ $t('editor:input-ids-or-handles') }}<br>
          {{ $t('editor:modify-initial-scope-if-no-variables-are-loaded') }}<br>
          {{ $t('editor:auto-initialize-empty-variable') }}
          <br><br>
          {{ $t('editor:open-webapp-on-prompt-use') }}
        </small>
        <div
          v-for="(p, index) in Object.values(dryRun.initialScope)"
          :key="index"
          class="mt-4"
        >
          <b-form-group
            v-if="p.lookup"
            :label="p.label"
            :description="p.description"
            label-class="text-primary"
          >
            <b-form-input
              v-model="p.value"
            />
          </b-form-group>
        </div>
      </div>
      <div
        v-else
        class="h-100"
      >
        <vue-json-editor
          :value="dryRun.input"
          :options="{ name: $t('editor:initial-scope') }"
          class="h-100"
          @input="onDryRunEdit"
        />
      </div>
    </b-modal>
  </div>
</template>

<script>
import mxgraph from 'mxgraph'
import Vue from 'vue'
import { encodeGraph } from '../lib/codec'
import { getStyleFromKind, getKindFromStyle } from '../lib/style'
import { encodeInput } from '../lib/dry-run'
import toolbarConfig from '../lib/toolbar'
import Configurator from '../components/Configurator'
import Tooltip from '../components/Tooltip.vue'
import WorkflowConfigurator from '../components/Configurator/Workflow'
import Help from '../components/Help'
import VueJsonEditor from 'v-jsoneditor'
import Import from '../components/Import'
import Export from '../components/Export'
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

const {
  mxClient,
  mxGraph,
  mxEvent,
  mxUtils,
  mxCell,
  mxGeometry,
  mxUndoManager,
  mxGraphHandler,
  mxEdgeHandler,
  mxKeyHandler,
  mxDivResizer,
  mxToolbar,
  mxConstants,
  mxDragSource,
  mxRubberband,
  mxPerimeter,
  mxEdgeStyle,
  mxConnectionHandler,
  mxClipboard,
  mxPoint,
  mxRectangle,
  mxLog,
  mxImage,
  mxConstraintHandler,
  mxConnectionConstraint,
  mxCellState,
  mxEllipse,
  mxCellOverlay,
  mxCellHighlight,
} = mxgraph({
  mxImageBasePath: `${document.getElementsByTagName('base')[0].href}icons`,
})

const originPoint = -2042

export default {
  name: 'WorkflowEditor',

  components: {
    Configurator,
    WorkflowConfigurator,
    Help,
    VueJsonEditor,
    Import,
    Export,
  },

  props: {
    workflowObject: {
      type: Object,
      default: () => {},
    },

    workflowTriggers: {
      type: Array,
      default: () => [],
    },

    changeDetected: {
      type: Boolean,
    },

    canCreate: {
      type: Boolean,
    },

    processingSave: {
      type: Boolean,
    },

    processingDelete: {
      type: Boolean,
    },
  },

  data () {
    return {
      initialized: false,

      deferred: false,
      triggersPathsChanged: false,

      graph: undefined,
      keyHandler: undefined,
      undoManager: undefined,

      workflow: {},
      triggers: [],
      vertices: {},
      edges: {},
      issues: {},

      highlights: {
        success: undefined,
        error: undefined,
      },

      runAsUser: undefined,

      toolbar: undefined,

      edgeConnected: false,

      rendering: false,

      sidebar: {
        item: undefined,
        itemType: undefined,
        outEdges: 0,
        show: false,
        showItem: false,
      },

      issuesModal: {
        show: false,
        issues: [],
      },

      dryRun: {
        show: false,
        processing: false,
        lookup: false,
        cellID: undefined,
        initialScope: {},
        input: {},
        inputEdited: {},
        sessionID: undefined,
      },

      selection: [],

      importProcessing: false,

      zoomLevel: 1,

      currentLabel: undefined,

      eventTypes: [],
      functionTypes: [],

      deferredKinds: ['delay', 'prompt'],
    }
  },

  computed: {
    getSidebarItemType () {
      const { item = {} } = this.sidebar || {}
      const { style, edge } = item.node || {}

      if (edge) {
        return this.$t('steps:path.short')
      }

      return this.$t(`steps:${style}.short`) || style
    },

    getSidebarItemIcon () {
      const { item } = this.sidebar

      if (item && item.config) {
        return this.getIcon(getStyleFromKind(item.config).icon, this.currentTheme)
      }

      return undefined
    },

    getSelectedItem () {
      return this.sidebar.item ? this.sidebar.item : undefined
    },

    getZoomPercent () {
      return `${Math.floor(this.zoomLevel * 100).toFixed(0)}%`
    },

    canUpdateWorkflow () {
      return this.workflow.workflowID === '0' ? this.canCreate : this.workflow.canUpdateWorkflow
    },

    nameState () {
      return this.workflow.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.workflow.handle)
    },

    canSave () {
      return !this.canUpdateWorkflow || [this.nameState, this.handleState].includes(false)
    },

    isDeleted () {
      return this.workflow.deletedAt
    },

    hasIssues () {
      return (this.workflow.issues || []).length
    },

    getRunAs () {
      if (this.runAsUser) {
        const { userID, name, username, email } = this.runAsUser
        return name || username || email || `<@${userID}>`
      }
      return undefined
    },

    currentTheme () {
      return this.$auth.user ? this.$auth.user.meta.theme : 'light'
    },
  },

  watch: {
    'workflow.runAs': {
      immediate: true,
      handler (runAs = '0') {
        if (runAs !== '0') {
          this.$SystemAPI.userRead({ userID: runAs })
            .then(user => {
              this.runAsUser = user
            })
        } else {
          this.runAsUser = undefined
        }
      },
    },

    workflowObject: {
      immediate: true,
      handler (workflow) {
        // If first save was successful, close workflow configurator modal
        if (workflow.workflowID !== this.workflow.workflowID) {
          this.$bvModal.hide('workflow')
        }

        this.workflow = workflow

        // Every change to workflowObject from parent component after initial render triggers rerender
        if (this.initialized) {
          this.render(this.workflow)
        }
      },
    },

    workflowTriggers: {
      immediate: true,
      handler (triggers) {
        this.triggers = triggers
      },
    },
  },

  mounted () {
    try {
      if (!mxClient.isBrowserSupported()) {
        throw new Error(mxUtils.error(this.$t('editor:unsupported-browser'), 200, false))
      }

      mxEvent.disableContextMenu(this.$refs.graph)
      this.graph = new mxGraph(this.$refs.graph, null, mxConstants.DIALECT_STRICTHTML)
      this.keyHandler = new mxKeyHandler(this.graph)

      this.setup()

      this.initToolbar()
      this.initUndoManager()
      this.initClipboard()

      this.keys()
      this.events()
      this.cellOverlay()

      this.styling()
      this.connectionHandler()

      this.getEventTypes()
      this.getFunctionTypes()

      this.$root.$on('trigger-updated', ({ mxObjectId }) => {
        this.redrawLabel(mxObjectId)
      })

      this.render(this.workflow, true)

      // Open workflow configurator if workflow is new
      if (this.workflow.workflowID && this.workflow.workflowID === '0') {
        this.$bvModal.show('workflow')
      }

      this.initialized = true
    } catch (e) {
      console.error(e)
    }
  },

  beforeDestroy () {
    // Destroy mxgraph singletons
    this.graph.destroy()
    this.keyHandler.destroy()
    this.toolbar.destroy()
    document.removeEventListener('keydown', this.keybinds)
  },

  methods: {
    deleteSelectedCells () {
      if (this.sidebar.item && this.graph.isCellSelected(this.sidebar.item.node)) {
        this.sidebarClose()
      }
      this.graph.removeCells()
    },

    sidebarClose () {
      this.sidebar.show = false

      setTimeout(() => {
        const mxObjectId = this.sidebar.item.node.mxObjectId
        this.sidebar.showItem = false
        this.sidebar.item = undefined
        this.sidebar.itemType = undefined
        this.redrawLabel(mxObjectId)
      }, 300)
    },

    sidebarDelete () {
      if (this.getSelectedItem) {
        this.graph.removeCells([this.getSelectedItem.node])
        this.sidebarClose()
      }
    },

    sidebarReopen (item, itemType) {
      this.sidebar.outEdges = (item.node.edges || []).length

      // If not open, just open sidebar
      if (!this.sidebar.show) {
        this.sidebar.item = item
        this.sidebar.itemType = itemType
        this.sidebar.show = true
        this.sidebar.showItem = true
        this.redrawLabel(item.node.mxObjectId)
      } else {
        // If item already opened in sidebar, keep open
        if (this.sidebar.item && item.node.id === this.sidebar.item.node.id) {
          return
        }

        // Otherwise fade item in and out
        const oldMxObjectId = ((this.getSelectedItem || {}).node || {}).mxObjectId
        this.sidebar.showItem = false
        this.sidebar.item = item
        this.sidebar.itemType = itemType
        this.redrawLabel(oldMxObjectId)
        this.redrawLabel(item.node.mxObjectId)
        setTimeout(() => {
          this.sidebar.showItem = true
        }, 100)
      }
    },

    setup () {
      this.graph.zoomFactor = 1.2

      // Sets a background image and restricts child movement to its bounds
      this.graph.setBackgroundImage(new mxImage(this.getIcon('grid', this.currentTheme), 8192, 8192))
      this.graph.maximumGraphBounds = new mxRectangle(0, 0, 8192, 8192)
      this.graph.gridSize = 8

      this.graph.setPanning(true)
      this.graph.setConnectable(true)
      this.graph.setAllowDanglingEdges(false)
      this.graph.setTooltips(true)

      /* eslint-disable no-new */
      new mxRubberband(this.graph) // Enables multiple selection
      this.graph.edgeLabelsMovable = false

      // Prevent showing tooltips on regular cells, just show overlay
      this.graph.getTooltipForCell = () => {}

      // Enables guides
      mxGraphHandler.prototype.guidesEnabled = true

      // Prevent cloning with ctrl + drag

      // Alt disables guides
      mxGraphHandler.prototype.useGuidesForEvent = (evt) => {
        return !mxEvent.isAltDown(evt.getEvent())
      }

      const mxGraphHandlerIsValidDropTarget = mxGraphHandler.prototype.isValidDropTarget
      mxGraphHandler.prototype.isValidDropTarget = function (target, me) {
        return mxGraphHandlerIsValidDropTarget.apply(this, arguments) && !target.edge
      }

      mxEdgeHandler.prototype.snapToTerminals = true

      mxGraph.prototype.minFitScale = 1
      mxGraph.prototype.maxFitScale = 1

      this.graph.isHtmlLabel = cell => {
        return true
      }

      this.graph.isWrapping = cell => {
        return true
      }

      this.graph.getLabel = cell => {
        let label = mxGraph.prototype.getLabel.apply(this, arguments)

        // Used to encode html labels to prevent security issues
        const encodeHTML = (value = '') => {
          if (value) {
            return value.replace(/[\u00A0-\u9999<>&]/gim, i => {
              return '&#' + i.charCodeAt(0) + ';'
            })
          }

          return value
        }

        if (cell.edge) {
          if (cell.value) {
            label = `<div id="openSidebar" class="text-nowrap py-1 px-3 mb-0 rounded bg-white pointer" style="border: 2px solid #A7D0E3; border-radius: 5px; color: var(--dark);">${encodeHTML(cell.value)}</div>`
          }
        } else if (this.vertices[cell.id]) {
          const vertex = this.vertices[cell.id]
          const { kind } = vertex.config
          const { style } = vertex.node

          if (vertex && kind !== 'visual') {
            const icon = this.getIcon(getStyleFromKind(vertex.config).icon, this.currentTheme)
            const type = this.$t(`steps:${style}.short`)
            const isSelected = this.selection.includes(cell.mxObjectId)
            const shadow = isSelected ? 'shadow' : 'shadow-sm'
            const cog = this.getIcon('cog')
            const issue = this.getIcon('issue')
            const playIcon = this.getIcon('play')
            const stopIcon = this.getIcon('stop')
            const opacity = kind === 'trigger' && !vertex.triggers.enabled ? 'opacity: 0.7;' : ''

            let test = ''
            let issues = ''
            let id = ''
            if (this.issues[cell.id]) {
              issues = `<img id="openIssues" src="${issue}" class="ml-2 pointer" style="width: 20px;"/>`
            } else {
              id = `<span class="show id-label">${cell.id}</span>`
            }

            let values = []

            if (kind === 'gateway' && cell.edges && cell.style !== 'gatewayParallel') {
              values = cell.edges
                .filter(({ source }) => cell.id === source.id)
                .map(({ id }) => this.edges[id])
                .map(({ node, config }) => `<tr><td><var>${encodeHTML(node.value)}</var></td><td><code>${encodeHTML(config.expr || '')}</code></td></tr>`)
                .join('')
            } else if (['expressions', 'function', 'prompt', 'iterator', 'exec-workflow'].includes(kind)) {
              let { arguments: args = [], results = [], ref } = vertex.config || {}

              const { meta = {}, results: functionResults = [] } = this.functionTypes.find(f => f.ref === ref) || {}

              const functionLabel = meta.short

              if (functionLabel) {
                values.push(`<tr><td><b class="text-primary">${functionLabel}</b></td><td/></tr>`)
              }

              if (args.length && kind !== 'expressions') {
                values.push('<tr class="title"><td><b>Arguments</b></td><td/></tr>')
              }
              args = args.map(({ target = '', type = 'Any', expr = '', value = '' }) => `<tr><td><var>${encodeHTML(target)}</var> <samp>(${type})</samp></td><td><code>${encodeHTML(expr || value)}</code></td></tr>`)

              if (results.length) {
                args.push('<tr class="title border-top"><td><b>Results</b></td><td /></tr>')
              }

              results = results.map(({ target = '', expr = '', value = '' }) => {
                const { types = [] } = functionResults.find(({ name }) => name === expr || name === value) || {}
                const type = types.length ? `(${types[0]})` : ''
                return `<tr><td><code>${encodeHTML(target)}</code> <samp>${type}</samp></td><td><var>${encodeHTML(expr || value)}</var></td></tr>`
              })

              values = [...values, ...args, ...results].join('')
            } else if (kind === 'trigger') {
              let { resourceType = '', eventType = '', constraints = [] } = vertex.triggers || {}
              let { properties = [] } = this.eventTypes.find(et => resourceType === et.resourceType && eventType === et.eventType) || {}

              if (resourceType) {
                resourceType = resourceType.split(':').map(s => {
                  return s[0].toUpperCase() + s.slice(1).toLowerCase()
                }).join(' ')
              }

              values.push('<tr class="title"><td><b>Configuration</b></td><td/><td/></tr>')
              values.push(`<tr><td><var>Resource</var></td><td/><td><code>${resourceType || ''}</code></td></tr>`)
              values.push(`<tr><td><var>Event</var></td><td/><td><code>${eventType || ''}</code></td></tr>`)

              if (constraints.length && eventType && eventType !== 'onManual') {
                constraints = [
                  '<tr class="title"><td><b>Constraints</b></td><td/><td/></tr>',
                  ...constraints.map(({ name = '', op = '', values = '' }) => {
                    return `<tr><td><samp>${name || eventType.includes('on') ? eventType.replace('on', '') : ''}</var></td><td><samp>${op}</samp></td><td><code>${encodeHTML(values.join(' or '))}</code></td></tr>`
                  }),
                ]
              } else {
                constraints = []
              }

              if (properties.length) {
                properties = [
                  '<tr class="title"><td><b>Initial scope</b></td><td/><td/></tr>',
                  ...properties.map(({ name = '', type = '' }) => {
                    return `<tr><td><var>${name}</var></td><td/><td><samp>${type || 'Any'}</samp></td></tr>`
                  }),
                ]
              }

              values = [
                ...values,
                ...constraints,
                ...properties,
              ].join('')
            } else if (['error', 'delay'].includes(kind)) {
              const { arguments: args = [] } = vertex.config || {}
              const { target, expr, value } = args[0] || {}

              if (target) {
                values = `<tr><td><var>${target}</var></td><td><code>${encodeHTML(expr || value)}</code></td></tr>`
              }
            } else {
              values = ''
            }

            if (values) {
              values = values
                ? '<div class="step-values rounded hide-label">' +
                    '<table class="table bg-white shadow mb-0">' +
                      values +
                    '</table>' +
                  '</div>'
                : ''
            }

            if (this.workflow.canExecuteWorkflow && vertex.triggers && (cell.edges || []).length) {
              if (!this.dryRun.processing) {
                test = `<img id="testWorkflow" title="${this.$t('configurator:tooltip.run-workflow')}" src="${playIcon}" class="hide pointer" style="width: 20px;"/>`
              } else if (this.dryRun.cellID === cell.id) {
                // If this is the trigger that is currently running
                test = `<span class="spinner-border text-success" data-toggle="tooltip" data-placement="top" style="width: 20px; height: 20px; cursor: default;" title="Testing in progress. If your workflow includes Prompt or Delay steps, it may be waiting for them to complete">
                          <span class="sr-only">
                            Spinning
                          </span>
                        </span>
                      `
                if (this.dryRun.sessionID) {
                  test = test + `<img id="cancelWorkflow" src="${stopIcon}" class="ml-2 hide pointer" style="width: 20px; height: 20px;"/>`
                }
              }
            }

            label = `<div class="d-flex flex-column bg-white border rounded step position-relative ${shadow}" style="min-width: 200px; border-radius: 5px;${opacity}">` +
                      '<div class=label-container">' +
                        '<div class="d-flex flex-row align-items-center text-primary px-2 my-1 h6 mb-0" style="width: 200px; height: 36px;">' +
                          `<img src="${icon}" class="mr-2"/>${type}` +
                          '<div class="d-flex h-100 ml-auto align-items-center">' +
                            test +
                            `<img id="openSidebar" title="${this.$t('steps:tooltip.configure-step')}" src="${cog}" class="hide pointer ml-2" style="width: 20px;"/>` +
                            id +
                            issues +
                          '</div>' +
                        '</div>' +
                        `<div class="label d-flex flex-grow-1 align-items-stretch bg-white border-top ${values ? 'wide-label' : ''}" style="max-width: 200px; min-height: 36px;">` +
                          `<span class="d-inline-block hover-untruncate p-2 bg-white">${encodeHTML(cell.value || '/')}</span>` +
                        '</div>' +
                      '</div>' +
                      values +
                    '</div>'
          } else {
            label = cell.value
          }
        }

        return label
      }

      this.graph.isCellEditable = () => {
        return false
      }

      // Disables mxGraph console window
      mxLog.setVisible = () => {}
      mxLog.DEBUG = false
      mxLog.TRACE = false

      mxGraph.prototype.expandedImage = undefined

      if (mxClient.IS_QUIRKS) {
        document.body.style.overflow = 'hidden'
        /* eslint-disable no-new */
        new mxDivResizer(this.graph.container)
      }

      if (mxClient.IS_NS) {
        mxEvent.addListener(this.graph.container, 'mousedown', () => {
          if (!this.graph.isEditing()) {
            this.graph.container.setAttribute('tabindex', '-1')
          }
        })
      }
    },

    initToolbar () {
      this.toolbar = new mxToolbar(this.$refs.toolbar)
      this.graph.dropEnabled = true

      // Matches DnD inside the this.editor.
      mxDragSource.prototype.getDropTarget = (graph, x, y) => {
        let cell = graph.getCellAt(x, y)

        if (!graph.isValidDropTarget(cell)) {
          cell = null
        }

        return cell
      }

      const addCell = ({ icon, width = 60, height = 60, style }) => {
        const { label, tooltip } = this.translateCell(style)
        let value = tooltip

        if (['break', 'continue'].includes(style)) {
          value = style === 'break' ? 'Stop iterator execution' : 'Skip current iteration'
        } else if (style.includes('gateway')) {
          value = style.split('gateway')[1]
        } else if (style === 'expressions') {
          value = 'Define and mutate scope variables'
        } else if (style === 'content') {
          value = 'Text here'
        }

        const cell = new mxCell(
          value,
          new mxGeometry(0, 0, width, height),
          style,
        )
        cell.setVertex(true)

        this.addToolbarItem(label, this.graph, this.toolbar, cell, icon, tooltip)
      }

      toolbarConfig.forEach(cell => {
        if (cell.kind === 'hr') {
          this.toolbar.addLine()
        } else if (cell.kind === 'nl') {
          this.toolbar.addBreak()
        } else {
          const cellStyle = getStyleFromKind(cell)
          if (cellStyle) {
            addCell({
              ...cell,
              ...cellStyle,
            })
          }
        }
      })
    },

    initUndoManager () {
      this.undoManager = new mxUndoManager()
      // Register UNDO and REDO
      const listener = (sender, evt) => {
        if (!this.rendering) {
          this.undoManager.undoableEditHappened(evt.getProperty('edit'))
        }
      }

      this.graph.getModel().addListener(mxEvent.UNDO, listener)
      this.graph.getView().addListener(mxEvent.UNDO, listener)

      this.graph.getModel().addListener(mxEvent.REDO, listener)
      this.graph.getView().addListener(mxEvent.REDO, listener)
    },

    makeCellCopy ({ edge, id }) {
      const cell = edge ? this.edges[id] : this.vertices[id]
      const node = this.graph.model.cloneCell(cell.node, false)
      node.id = cell.node.id
      node.parent = cell.node.parent.id

      if (edge) {
        node.source = cell.node.source.id
        node.target = cell.node.target.id
      }

      const cellCopy = {
        node,
      }

      // Need to use JSON.parse to remove references
      if (cell.config) {
        cellCopy.config = JSON.parse(JSON.stringify(cell.config))
      }

      if (cell.triggers) {
        const triggers = {
          enabled: cell.triggers.enabled,
          constraints: cell.triggers.constraints,
          eventType: cell.triggers.eventType,
          resourceType: cell.triggers.resourceType,
        }

        cellCopy.triggers = JSON.parse(JSON.stringify(triggers))
      }

      return cellCopy
    },

    initClipboard () {
      const absoluteGeometry = (cell) => {
        // If parent not container cell
        if (!cell.parent.geometry) {
          return {
            x: cell.geometry.x,
            y: cell.geometry.y,
          }
        }

        // Get absoluteGeometry recursively
        const { x, y } = absoluteGeometry(cell.parent)
        return {
          x: cell.geometry.x + x,
          y: cell.geometry.y + y,
        }
      }

      const copyCells = (cells, parentID) => {
        let copiedCells = {}
        let copiedEdges = []

        cells.forEach(cell => {
          const newCell = this.makeCellCopy(cell)

          if (cell.edge) {
            copiedEdges.push(newCell)
          } else {
            if (!copiedCells[cell.parent.id]) {
              copiedCells[cell.parent.id] = []
            }

            // Cell parent is container cell, and copyCells wasn't called by parent of curent cell
            if (cell.parent.geometry && cell.parent.id !== parentID) {
              // Get relative x,y recursively - Needed because of relative x,y in swimlanes. And paste not respecting relative x,y
              const { x, y } = absoluteGeometry(cell)
              newCell.node.geometry.x = x
              newCell.node.geometry.y = y
            }

            copiedCells[cell.parent.id].push(newCell)

            // Handle children
            if (cell.children) {
              if (!copiedCells[cell.id]) {
                copiedCells[cell.id] = []
              }

              // Recursively copy children
              const childrenCells = copyCells(cell.children, cell.id)
              copiedCells = { ...copiedCells, ...childrenCells.cells }
              copiedEdges = [...copiedEdges, ...childrenCells.edges]
            }
          }
        })

        return {
          cells: copiedCells,
          edges: copiedEdges,
        }
      }

      const pasteCells = (evt) => {
        if (evt.clipboardData.getData('text').includes('"cells":')) {
          // Get cells from actual system clipboard
          const { cells = {}, edges = [] } = JSON.parse(evt.clipboardData.getData('text')) || {}

          const delta = mxClipboard.insertCount * this.graph.gridSize
          const defaultParent = this.graph.getDefaultParent()
          const newCellIDs = {}
          const allCells = []

          this.graph.getModel().beginUpdate()

          // Handle cells
          Object.entries(cells).forEach(([parentID, children]) => {
            children.forEach(({ node, ...rest }) => {
              const parent = newCellIDs[parentID] ? this.graph.model.getCell(newCellIDs[parentID]) : defaultParent
              const { id, geometry, value, style } = node

              // Offset only the topmost cell. Children are relative and will be offset automaticially
              if (!newCellIDs[parentID]) {
                geometry.x += delta
                geometry.y += delta
              }

              const newVertex = this.graph.insertVertex(parent, null, value, geometry.x, geometry.y, geometry.width, geometry.height, style)
              allCells.push(newVertex)

              newCellIDs[id] = newVertex.id
              rest.config.stepID = newVertex.id

              this.vertices[newVertex.id] = { node: newVertex, ...rest }
            })
          })

          // Handle edges
          edges.forEach(({ node, ...rest }) => {
            const parent = newCellIDs[node.id] ? this.graph.model.getCell(newCellIDs[node.id]) : defaultParent
            const { id, geometry, value, style } = node

            const source = (this.vertices[newCellIDs[rest.config.parentID || node.source]] || {}).node
            const target = (this.vertices[newCellIDs[rest.config.childID || node.target]] || {}).node
            if (!source || !target) {
              return
            }

            node.source = source
            node.target = target

            const newEdge = this.graph.insertEdge(parent, null, value, node.source, node.target, style)
            newEdge.geometry.points = (geometry.points || []).map(({ x, y }) => {
              return new mxPoint(x, y)
            })
            allCells.push(newEdge)

            newCellIDs[id] = newEdge.id
            rest.config.parentID = node.source.id
            rest.config.childID = node.target.id

            this.edges[newEdge.id] = { node: newEdge, ...rest }
          })

          Object.keys(this.vertices).forEach(vID => this.updateVertexConfig(vID))

          mxClipboard.insertCount++
          this.graph.setSelectionCells(allCells)
          this.graph.getModel().endUpdate() // Updates the display
        }
      }

      mxClipboard.copy = (graph, cells) => {
        const exportableCells = graph.getExportableCells(graph.model.getTopmostCells(cells || graph.getSelectionCells()))
        const copiedCells = copyCells(exportableCells)

        // Copy to actual browser(system) clipboard
        const editor = this.$refs.editor
        const tempInput = document.createElement('input')
        editor.appendChild(tempInput)
        tempInput.setAttribute('value', JSON.stringify(copiedCells))
        tempInput.select()
        document.execCommand('copy')
        editor.removeChild(tempInput)

        mxClipboard.insertCount = 1

        return copiedCells
      }

      mxClipboard.cut = (graph, cells) => {
        const copiedCells = mxClipboard.copy(graph, cells)

        const cutCells = []
        Object.entries(copiedCells.cells).forEach(([parentID, children]) => {
          children.forEach(({ node }) => {
            cutCells.push(this.graph.model.getCell(node.id))
          })
        })

        copiedCells.edges.forEach(({ node }) => {
          cutCells.push(this.graph.model.getCell(node.id))
        })

        mxClipboard.insertCount = 0
        this.graph.removeCells(cutCells)

        return cells
      }

      // Register paste handler
      document.querySelector('body').addEventListener('paste', pasteCells)
    },

    // Works on all of editor (mostly)
    keybinds (event) {
      // Ctrl + S
      if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault()
        // Prevent the workflow from being saved if expressions editor is open
        if (!document.getElementById('expression-editor')) {
          this.saveWorkflow()
        }
      }
    },

    // Only works when canvas is focused
    keys () {
      // Register general keydown event if we need it (we destroy it in beforeDestroy)
      document.addEventListener('keydown', this.keybinds)

      // Register control and meta key if Mac
      this.keyHandler.getFunction = (evt) => {
        if (evt != null) {
          // If CTRL or META key is pressed
          if (evt.ctrlKey || (mxClient.IS_MAC && evt.metaKey)) {
            // If SHIFT key is pressed
            if (evt.shiftKey) {
              return this.keyHandler.controlShiftKeys[evt.keyCode]
            }
            return this.keyHandler.controlKeys[evt.keyCode]
          }

          // If only normal keys are pressed
          return this.keyHandler.normalKeys[evt.keyCode]
        }

        return null
      }

      // Ctrl + X
      this.keyHandler.controlKeys[88] = () => {
        mxClipboard.cut(this.graph, this.graph.getSelectionCells())
      }

      // Ctrl + C
      this.keyHandler.controlKeys[67] = () => {
        mxClipboard.copy(this.graph, this.graph.getSelectionCells())
      }

      // Ctrl + A
      this.keyHandler.controlKeys[65] = () => {
        this.graph.selectAll()
      }

      // Ctrl + Z
      this.keyHandler.controlKeys[90] = () => {
        this.undoManager.undo()
        this.checkExistingTriggerPaths()
      }

      // Ctrl + Shift + Z
      this.keyHandler.controlShiftKeys[90] = () => {
        this.undoManager.redo()
        this.checkExistingTriggerPaths()
      }

      // Backspace
      this.keyHandler.normalKeys[8] = () => {
        this.deleteSelectedCells()
      }

      // Delete
      this.keyHandler.normalKeys[46] = () => {
        this.deleteSelectedCells()
      }

      // Ctrl + Space, Resets view to original state (zoom = 1, x = 0, y = 0)
      this.keyHandler.controlKeys[32] = () => {
        if (this.graph.model.getChildCount(this.graph.getDefaultParent())) {
          this.graph.fit()
          this.graph.view.setTranslate(this.graph.view.translate.x + 79, this.graph.view.translate.y + 220)
          this.zoomLevel = this.graph.view.scale
        } else {
          this.resetZoom()
          this.graph.view.setTranslate(originPoint, originPoint)
        }
      }

      // Shift + ?
      this.keyHandler.bindKey(191, (event) => {
        if (event.shiftKey && event.key === '?') {
          this.$refs.help.click()
        }
      })

      // Nudge
      const nudge = (keyCode, evt) => {
        if (!this.graph.isSelectionEmpty()) {
          let dx = 0
          let dy = 0

          // If shift is not pressed move cell by whole grid block
          const delta = evt.shiftKey ? this.graph.gridSize : this.graph.gridSize * 5

          if (keyCode === 37) {
            dx = -delta
          } else if (keyCode === 38) {
            dy = -delta
          } else if (keyCode === 39) {
            dx = delta
          } else if (keyCode === 40) {
            dy = delta
          }

          this.graph.moveCells(this.graph.getSelectionCells(), dx, dy)
        }
      }

      // Move cells with arrow keys
      this.keyHandler.bindKey(37, (evt) => {
        nudge(37, evt)
      })

      this.keyHandler.bindKey(38, (evt) => {
        nudge(38, evt)
      })

      this.keyHandler.bindKey(39, (evt) => {
        nudge(39, evt)
      })

      this.keyHandler.bindKey(40, (evt) => {
        nudge(40, evt)
      })
    },

    checkExistingTriggerPaths () {
      // If trigger was reconnected, check if all triggers are still connected to the previous stepID
      this.triggersPathsChanged = [...this.triggers].some(({ stepID = '0', meta = {} }) => {
        if (stepID !== NoID) {
          let [triggerEdge] = this.vertices[meta.visual.id].node.edges || []

          if (triggerEdge) {
            triggerEdge = this.graph.model.getCell(triggerEdge.id)
            return triggerEdge.target && triggerEdge.target.id !== stepID
          } else {
            return true
          }
        }

        return false
      })
    },

    events () {
      // Redraw selected border of cells that have been newly added/removed
      this.graph.getSelectionModel().addListener(mxEvent.CHANGE, (sender, evt) => {
        const cells = [...(evt.getProperty('added') || []), ...(evt.getProperty('removed') || [])]
        this.selection = this.graph.getSelectionCells().map(({ mxObjectId }) => mxObjectId)
        cells.forEach(({ mxObjectId }) => {
          this.redrawLabel(mxObjectId)
        })
      })

      // Edge connect event
      this.graph.connectionHandler.addListener(mxEvent.CONNECT, (sender, evt) => {
        const node = evt.getProperty('cell')

        this.edges[node.id] = {
          node,
          config: {
            parentID: node.source.id,
            childID: node.source.id,
          },
        }

        const source = this.vertices[node.source.id]
        const target = this.vertices[node.target.id]
        const outPaths = source.node.edges.filter(e => e.source.id === source.node.id) || []

        if (target.config.kind === 'gateway') {
          if (['join', 'fork'].includes(target.config.ref)) {
            this.updateVertexConfig(target.node.id)
          }
        }

        if (source.config.kind === 'gateway') {
          if (['join', 'fork'].includes(source.config.ref)) {
            this.updateVertexConfig(source.node.id)
          }

          if (source.config.ref === 'excl') {
            this.edges[node.id].node.value = `#${outPaths.length} - ${outPaths.length === 1 ? 'If' : 'Else (if)'}`
          } else if (source.config.ref === 'incl') {
            this.edges[node.id].node.value = 'If'
          }

          this.sidebar.outEdges = (source.node.edges || []).length
        } else if (source.config.kind === 'error-handler') {
          this.edges[node.id].node.value = `${outPaths.length === 1 ? 'Try' : 'Catch'}`
        } else if (source.config.kind === 'iterator') {
          this.edges[node.id].node.value = `${outPaths.length === 1 ? 'Body' : 'End'}`
        }

        this.edgeConnected = true
      })

      this.graph.addListener(mxEvent.CELL_CONNECTED, (sender, evt) => {
        if (!this.rendering) {
          const edge = evt.getProperty('edge')
          const source = this.vertices[edge.source.id]
          if (source.config.kind === 'trigger') {
            this.checkExistingTriggerPaths()
          }
        }
      })

      this.graph.addListener(mxEvent.CELLS_ADDED, (sender, evt) => {
        if (!this.rendering) {
          const cells = evt.getProperty('cells')
          let lastVertexID = null
          cells.forEach(cell => {
            if (cell && cell.vertex) {
              if (!this.rendering) {
                cell.defaultName = true
                this.addCellToVertices(cell)
                this.graph.setSelectionCells([cell])
                lastVertexID = cell.id
              }
            }
          })

          if (lastVertexID) {
            this.$nextTick(() => {
              const vertex = this.vertices[lastVertexID]
              this.sidebarReopen(vertex, vertex.config.kind)
            })
          }
        }
      })

      this.graph.addListener(mxEvent.CELLS_REMOVED, (sender, evt) => {
        const cells = evt.getProperty('cells') || []
        cells.forEach(cell => {
          if (cell.edge) {
            const source = this.vertices[cell.source.id]
            const target = this.vertices[cell.target.id]

            // If exlusive gateway, update edge indexes (#n)
            if (source.config.kind === 'gateway') {
              if (source.config.ref === 'excl') {
                source.node.edges.filter(e => e.source.id === source.node.id).forEach((edge, index) => {
                  /* eslint-disable no-unused-vars */
                  const [edgeID, ...rest] = edge.value.split(' - ')

                  this.edges[edge.id].node.value = `#${index + 1} - ${rest.join(' - ')}`
                  this.redrawLabel(edge.mxObjectId)
                })
              }

              if (['join', 'fork'].includes(target.config.ref)) {
                this.updateVertexConfig(source.node.id)
              }
            } else if (source.config.kind === 'iterator' || source.config.kind === 'error-handler') {
              // Since later placed edges will have greater id than the ones placed before, we can filter by id
              // Remove all edges that were placed after the one that was just deleted.
              // This needs to be done to preserve edge order
              this.graph.removeCells(source.node.edges.filter(e => e.source.id === source.node.id && e.id > cell.id))
            } else if (source.config.kind === 'trigger') {
              this.checkExistingTriggerPaths()
            }

            if (target.config.kind === 'gateway') {
              if (['join', 'fork'].includes(target.config.ref)) {
                this.updateVertexConfig(target.node.id)
              }
            }
          }
        })
      })

      // Zoom event
      mxEvent.addMouseWheelListener((event, up) => {
        if (mxEvent.isConsumed(event)) {
          return
        }

        if (mxEvent.isControlDown(event) || (mxClient.IS_MAC && mxEvent.isMetaDown(event))) {
          return
        }

        this.zoom(up)
        mxEvent.consume(event)
      }, this.graph.container)

      // On hover, bring cell to foreground
      this.graph.addMouseListener({
        mouseMove: (sender, evt) => {
          if (this.currentLabel !== null && evt.getState() === this.currentLabel) {
            return
          }

          let tmp = sender.view.getState(evt.getCell())

          // Ignores everything but vertices
          if (tmp !== null && !sender.getModel().isVertex(tmp.cell)) {
            tmp = null
          }

          if (tmp !== this.currentLabel) {
            this.currentLabel = tmp
            if (this.currentLabel?.cell) {
              this.rendering = true
              sender.orderCells(false, [this.currentLabel.cell])
              this.rendering = false
            }
          }
        },
        mouseUp: (sender, evt) => {
          evt.consume()
        },
        mouseDown: (sender, evt) => {
          // Prevent click event handling if edge was just connected
          if (this.edgeConnected) {
            this.edgeConnected = false
            return
          }

          const event = evt.evt
          const cell = evt.state?.cell

          if (event) {
            if (mxEvent.isControlDown(event) || (mxClient.IS_MAC && mxEvent.isMetaDown(event))) {
              // Prevent sidebar opening/closing when CTRL(CMD) is pressed while clicking
            } else if (cell) {
              // If clicked on Cog icon
              const item = cell.edge ? this.edges[cell.id] : this.vertices[cell.id]
              const itemType = cell.edge ? 'edge' : item.config.kind

              if (event.target.id === 'openSidebar' || item.config.kind === 'visual') {
                this.sidebarReopen(item, itemType)
              } else if (event.target.id === 'openIssues') {
                this.issuesModal.issues = this.issues[cell.id]
                this.issuesModal.show = true
              } else if (event.target.id === 'testWorkflow') {
                this.dryRun.cellID = cell.id
                this.loadTestScope()
              } else if (event.target.id === 'cancelWorkflow') {
                this.cancelWorkflow()
              }
            } else if (!event.defaultPrevented) {
              // If click is on background and is not multiple selection, deselect all selected cells
              this.graph.getSelectionModel().clear()
              this.sidebar.show = false
              if (this.getSelectedItem) {
                this.sidebarClose()
              }
            }
          }

          evt.consume()
        },
      })

      this.graph.model.addListener(mxEvent.CHANGE, (sender, evt) => {
        if (!this.rendering) {
          this.removeDryRunOverlay()
          this.$root.$emit('change-detected')
        }
      })
    },

    styling () {
      // General
      mxConstants.VERTEX_SELECTION_COLOR = '#A7D0E3'
      mxConstants.VERTEX_SELECTION_STROKEWIDTH = 2
      mxConstants.EDGE_SELECTION_COLOR = '#A7D0E3'
      mxConstants.EDGE_SELECTION_STROKEWIDTH = 2
      mxConstants.DEFAULT_FONTFAMILY = 'Poppins-Regular'
      mxConstants.DEFAULT_FONTSIZE = 13

      mxConstants.HANDLE_FILLCOLOR = '#A7D0E3'
      mxConstants.HANDLE_STROKECOLOR = 'none'
      mxConstants.CONNECT_HANDLE_FILLCOLOR = '#A7D0E3'
      mxConstants.OUTLINE_HIGHLIGHT_COLOR = '#A7D0E3'
      mxConstants.TARGET_HIGHLIGHT_COLOR = '#A7D0E3'
      mxConstants.DROP_TARGET_COLOR = '#A7D0E3'
      mxConstants.DEFAULT_VALID_COLOR = '#A7D0E3'
      mxConstants.VALID_COLOR = '#A7D0E3'
      mxGraphHandler.prototype.previewColor = '#A7D0E3'

      mxConstants.STYLE_PERIMETER = mxPerimeter.RectanglePerimeter

      mxConstants.GUIDE_COLOR = 'var(--dark)'
      mxConstants.GUIDE_STROKEWIDTH = 1

      // Creates the default style for vertices

      let style = this.graph.getStylesheet().getDefaultVertexStyle()
      style[mxConstants.STYLE_SHAPE] = mxConstants.SHAPE_RECTANGLE
      style[mxConstants.STYLE_PERIMETER] = mxPerimeter.RectanglePerimeter
      style[mxConstants.STYLE_STROKECOLOR] = 'none'
      style[mxConstants.STYLE_STROKEWIDTH] = 0
      style[mxConstants.STYLE_ROUNDED] = true
      style[mxConstants.STYLE_ARCSIZE] = 5
      style[mxConstants.STYLE_RESIZABLE] = false
      style[mxConstants.STYLE_FILLCOLOR] = 'none'
      style[mxConstants.STYLE_FONTCOLOR] = 'var(--dark)'
      style[mxConstants.STYLE_FONTSIZE] = 13
      this.graph.getStylesheet().putDefaultVertexStyle(style)

      // Creates the default style for edges
      style = this.graph.getStylesheet().getDefaultEdgeStyle()
      style[mxConstants.STYLE_STROKECOLOR] = '#A7D0E3'
      style[mxConstants.STYLE_EDGE] = mxEdgeStyle.OrthConnector
      style[mxConstants.STYLE_ROUNDED] = true
      style[mxConstants.STYLE_ORTHOGONAL] = true
      style[mxConstants.STYLE_MOVABLE] = false
      style[mxConstants.STYLE_FONTCOLOR] = 'var(--dark)'
      style[mxConstants.STYLE_STROKEWIDTH] = 2
      style[mxConstants.STYLE_ENDSIZE] = 15
      style[mxConstants.STYLE_STARTSIZE] = 15
      style[mxConstants.STYLE_SOURCE_JETTY_SIZE] = 40
      style[mxConstants.STYLE_TARGET_JETTY_SIZE] = 40
      this.graph.getStylesheet().putDefaultEdgeStyle(style)

      // Swimlane
      style = {}
      style[mxConstants.STYLE_ROUNDED] = true
      style[mxConstants.STYLE_ARCSIZE] = 5
      style[mxConstants.STYLE_RESIZABLE] = true
      style[mxConstants.STYLE_SHAPE] = mxConstants.SHAPE_SWIMLANE
      style[mxConstants.STYLE_FONTSIZE] = 15
      style[mxConstants.STYLE_HORIZONTAL] = false
      style[mxConstants.STYLE_VERTICAL_LABEL_POSITION] = mxConstants.ALIGN_MIDDLE
      style[mxConstants.STYLE_VERTICAL_ALIGN] = mxConstants.ALIGN_MIDDLE
      style[mxConstants.STYLE_FILLCOLOR] = 'var(--white)'
      style[mxConstants.STYLE_STROKECOLOR] = 'var(--dark)'
      style[mxConstants.STYLE_STROKEWIDTH] = 1
      this.graph.getStylesheet().putCellStyle('swimlane', style)

      // Content
      style = {}
      style[mxConstants.STYLE_RESIZABLE] = true
      style[mxConstants.STYLE_CONNECTABLE] = false
      style[mxConstants.STYLE_FILLCOLOR] = 'var(--white)'
      style[mxConstants.STYLE_STROKECOLOR] = 'var(--extra-light)'
      style[mxConstants.STYLE_STROKEWIDTH] = 1
      style[mxConstants.STYLE_VERTICAL_ALIGN] = mxConstants.ALIGN_TOP
      style[mxConstants.STYLE_ALIGN] = mxConstants.ALIGN_LEFT
      style[mxConstants.STYLE_SPACING_TOP] = 10
      style[mxConstants.STYLE_SPACING_LEFT] = 10
      style[mxConstants.STYLE_WHITE_SPACE] = 'wrap'
      style[mxConstants.STYLE_OVERFLOW] = 'hidden'
      this.graph.getStylesheet().putCellStyle('content', style)
    },

    translateCell (style) {
      return {
        label: this.$t(`steps:${style}.label`),
        tooltip: this.$t(`steps:${style}.tooltip`),
      }
    },

    cellOverlay () {
      mxCellOverlay.prototype.defaultOverlap = 1.2
    },

    connectionHandler () {
      mxConstraintHandler.prototype.intersects = function (icon, point, source, existingEdge) {
        return (!source || existingEdge) || mxUtils.intersects(icon.bounds, point)
      }

      // Removes default connect logic (from center of cell)
      if (this.graph.connectionHandler.connectImage === null) {
        this.graph.connectionHandler.isConnectableCell = () => {
          return false
        }
        mxEdgeHandler.prototype.isConnectableCell = cell => {
          return this.graph.connectionHandler.isConnectableCell(cell)
        }
      }

      this.graph.getAllConnectionConstraints = function (terminal, source = false) {
        if (!terminal) {
          return null
        }

        const { cell } = terminal

        let isConnectable = this.model.isVertex(cell) && !['swimlane', 'content'].includes(cell.style)

        // Only one outbound connection per trigger
        if (cell.style.includes('trigger') && cell.edges) {
          isConnectable = isConnectable && !cell.edges.length
        }

        if (isConnectable) {
          let possibleConnections = [
            [0, 0],
            [0.25, 0],
            [0.5, 0],
            [0.75, 0],
            [1, 0],
            [1, 0.25],
            [1, 0.5],
            [1, 0.75],
            [1, 1],
            [0.75, 1],
            [0.5, 1],
            [0.25, 1],
            [0, 1],
            [0, 0.75],
            [0, 0.5],
            [0, 0.25],
          ]

          // Allows for multiple inbound edges on the same point, but not outbound from the same point
          if (source) {
            const edges = cell.edges || []
            edges.forEach(({ source, target, style }) => {
              const points = {}
              if (style) {
                style.split(';').forEach(point => {
                  const [key, value] = point.split('=')
                  if (key && value) {
                    points[key] = parseFloat(value)
                  }
                })

                possibleConnections = possibleConnections.filter(([x, y]) => {
                  // Outgoing edge, check exitX/Y
                  if (source.id === cell.id) {
                    // Filter out exit point
                    return !(x === points.exitX && y === points.exitY)
                  } else if (target.id === cell.id) {
                    // Incoming edge
                    return !(x === points.entryX && y === points.entryY)
                  }
                  return true
                })
              }
            })
          } else {
            // Prevent triggers from being connected to
            if (cell.style.includes('trigger')) {
              possibleConnections = []
            }
          }

          return possibleConnections.map(([x, y]) => {
            return new mxConnectionConstraint(new mxPoint(x, y), true)
          })
        }

        return null
      }

      // Connect preview
      mxConnectionHandler.prototype.createEdgeState = function (me) {
        const edge = this.graph.createEdge(null, null, null, null, null)
        return new mxCellState(this.graph.view, edge, this.graph.getStylesheet().getDefaultEdgeStyle())
      }

      // Resets control points when related cells are moved
      this.graph.resetEdgesOnMove = true
      mxGraph.prototype.resetEdges = function (cells) {
        if (cells != null) {
          this.model.beginUpdate()
          try {
            cells.forEach(cell => {
              const edges = this.model.getEdges(cell)
              if (edges != null) {
                edges.forEach(edge => {
                  this.resetEdge(edge)
                })
              }

              this.resetEdges(this.model.getChildren(cell))
            })
          } finally {
            this.model.endUpdate()
          }
        }
      }

      // Image for fixed point
      mxConstraintHandler.prototype.pointImage = new mxImage(this.getIcon('connection-point'), 16, 16)

      // On hover outline for fixed point
      mxConstraintHandler.prototype.createHighlightShape = function () {
        return new mxEllipse(null, '#A7D0E3', '#A7D0E3', 1)
      }
    },

    addToolbarItem (title, graph, toolbar, prototype, icon, tooltip) {
      const funct = (graph, evt, cell) => {
        graph.stopEditing(false)

        const pt = graph.getPointForEvent(evt)
        const vertex = graph.getModel().cloneCell(prototype)
        vertex.geometry.x = pt.x
        vertex.geometry.y = pt.y

        graph.importCells([vertex], 0, 0, cell)
      }

      const dragElt = document.createElement('div')
      dragElt.style.border = 'dashed #A7D0E3 2px'
      dragElt.style.width = `${prototype.geometry.width}px`
      dragElt.style.height = `${prototype.geometry.height}px`

      icon = this.getIcon(icon, this.currentTheme)

      const img = toolbar.addMode(title, icon, funct)

      const ds = mxUtils.makeDraggable(img, graph, funct, dragElt, null, null, this.graph.autoscroll, true)

      // Init step tooltip
      img.id = prototype.style.split(';')[0]
      const TooltipComponent = Vue.extend(Tooltip)
      const instance = new TooltipComponent({
        propsData: { title, kind: img.id, img: icon, text: tooltip },
      })
      instance.$mount()
      this.$refs.tooltips.appendChild(instance.$el)

      // When dragged over toolbar it shows as img otherwise show border
      ds.createDragElement = mxDragSource.prototype.createDragElement
    },

    addCellToVertices (cell) {
      const triggers = this.triggers.find(({ meta }) => {
        return ((meta || {}).visual || {}).id === cell.id
      })

      const {
        kind = '',
        ref = '',
        defaultName = false,
        arguments: args,
        results = [],
        meta = {},
      } = (this.workflow.steps || []).find(({ stepID }) => {
        return stepID === cell.id
      }) || {}

      this.vertices[cell.id] = {
        node: cell,
        config: {
          stepID: cell.id,
          kind: kind || '',
          ref: ref || '',
          defaultName: defaultName || meta.visual?.defaultName || cell.defaultName || false,
          ...(this.rendering ? {} : getKindFromStyle(cell)),
        },
      }

      if (args) {
        this.vertices[cell.id].config.arguments = args
      }

      if (results) {
        this.vertices[cell.id].config.results = results
      }

      if (triggers || cell.style === 'trigger') {
        this.vertices[cell.id].triggers = triggers || {
          resourceType: null,
          eventType: null,
          constraints: [],
          enabled: true,
        }
      }
    },

    updateVertexConfig (vID) {
      const { node, config } = this.vertices[vID]
      this.vertices[vID].config = { ...config, ...(this.rendering ? {} : getKindFromStyle(node)) }
    },

    setValue (value, defaultName = false) {
      this.graph.model.setValue(this.sidebar.item.node, value)

      if (this.sidebar.itemType !== 'edge') {
        this.vertices[this.sidebar.item.node.id].config.defaultName = defaultName
      }
    },

    zoom (up = true) {
      if (up && this.graph.view.scale < 3) {
        this.graph.zoomIn()
      } else if (!up && this.graph.view.scale > 0.1) {
        this.graph.zoomOut()
      }
      this.zoomLevel = this.graph.view.scale
    },

    resetZoom () {
      this.graph.zoomTo(1)
      this.zoomLevel = this.graph.view.scale
    },

    redrawLabel (id = '') {
      if (id) {
        const state = this.graph.view.states.map[id]
        if (state) {
          this.graph.cellRenderer.redrawLabel(state)
        }
      }
    },

    removeDryRunOverlay () {
      if (this.highlights.length > 0) {
        this.highlights.forEach(h => {
          h.destroy()
        })
        this.highlights = []
        this.graph.clearCellOverlays()
      }
    },

    async loadTestScope () {
      // Can only test saved workflow
      if (this.changeDetected) {
        this.toastWarning(this.$t('notification:save-workflow'), this.$t('notification:failed-test'))
        return
      }

      // Can only test valid workflow
      if (this.hasIssues) {
        this.toastWarning(this.$t('notification:resolve-issues'), this.$t('notification:failed-test'))
        return
      }

      const lookupableTypes = [
        'record',
        'oldRecord',
        'module',
        'oldModule',
        'page',
        'oldPage',
        'namespace',
        'oldNamespace',
        'user',
        'oldUser',
        'role',
        'oldRole',
        'application',
        'oldApplication',
      ]

      // Assume trigger is valid since workflow is saved and has no issues
      const { resourceType, eventType } = this.vertices[this.dryRun.cellID].triggers
      const et = (this.eventTypes.find(et => resourceType === et.resourceType && eventType === et.eventType) || {}).properties
      if (et) {
        // Flag to check if lookup should be opened, or JSON editor
        let lookup = false
        if (et.length) {
          this.dryRun.initialScope = et.reduce((initialScope, p) => {
            let label = `${p.name}${lookupableTypes.includes(p.name) ? this.$t('editor:id-parenthesis') : ''}`
            if (p.type === 'ComposeNamespace' || p.type === 'ComposeModule') {
              label = `${p.name} ${this.$t('editor:handle')}`
            }

            let description = ''
            if (p.type === 'ComposeRecord') {
              description = this.$t('editor:required-namespace-and-module')
            } else if (p.type === 'ComposeModule' || p.name === 'page' || p.name === 'oldPage') {
              description = this.$t('editor:required-namespace')
            }

            initialScope[p.name] = ({
              label,
              value: (this.dryRun.initialScope[p.name] || {}).value,
              lookup: lookupableTypes.includes(p.name),
              description,
            })

            lookup = lookup ? true : lookupableTypes.includes(p.name)
            return initialScope
          }, {})

          // Set initial values for unlookable types
          encodeInput(this.dryRun.initialScope, this.$ComposeAPI, this.$SystemAPI)
            .then(input => {
              this.dryRun.input = input
              this.dryRun.lookup = lookup
              this.dryRun.show = true
            })
            .catch(this.toastErrorHandler(this.$t('notification:initial-scope-load-failed')))
        } else {
          // If no constraints, just run
          this.dryRun.initialScope = {}
          this.testWorkflow()
        }
      } else {
        this.toastWarning(this.$t('notification:event-type-not-found'), this.$t('notification:failed-test'))
      }
    },

    async dryRunOk (e) {
      if (this.dryRun.lookup) {
        e.preventDefault()
        // Lookup based on provided ids
        encodeInput(this.dryRun.initialScope, this.$ComposeAPI, this.$SystemAPI)
          .then(input => {
            this.dryRun.input = input
            this.dryRun.inputEdited = input
            this.dryRun.lookup = false
          })
          .catch(this.toastErrorHandler(this.$t('notification:initial-scope-load-failed')))
      } else {
        this.testWorkflow(this.dryRun.inputEdited)
      }
    },

    onDryRunEdit (e) {
      this.dryRun.inputEdited = e
    },

    async testWorkflow (input = {}) {
      this.removeDryRunOverlay()
      this.dryRun.processing = true
      this.redrawLabel(this.graph.model.getCell(this.dryRun.cellID).mxObjectId)

      const testParams = {
        workflowID: this.workflow.workflowID,
        stepID: this.vertices[this.dryRun.cellID].triggers.stepID,
        trace: this.workflow.canManageWorkflowSessions || false,
        wait: false,
        async: true,
        input,
      }

      this.toastInfo(this.$t('notification:started-test'), this.$t('notification:test-in-progress'))

      await this.$AutomationAPI.workflowExec(testParams)
        .then(({ sessionID, error: wfExecErr }) => {
          this.dryRun.sessionID = sessionID
          this.redrawLabel(this.graph.model.getCell(this.dryRun.cellID).mxObjectId)

          const sessionHandler = ({ completedAt, status, stacktrace, error = false }) => {
            if (completedAt) {
              // If stacktrace exists, render it
              if (stacktrace) {
                this.renderTrace(testParams.stepID, stacktrace)

                if (status === 'completed') {
                  this.toastSuccess(this.$t('notification:workflow-test-completed'), this.$t('notification:test-completed'))
                }
              }

              // If error or no stacktrace, raise an error/warning
              if (error) {
                throw new Error(error)
              } else if (!stacktrace) {
                this.toastWarning(this.$t('notification:trace-unavailable'), this.$t('notification:test-completed'))
              }
            } else {
              setTimeout(sessionReader, 1000)
            }
          }

          // Check if session is completed/failed every second
          const sessionReader = () => {
            this.$AutomationAPI.sessionRead({ sessionID })
              .then(sessionHandler)
              .catch(err => {
                // In case of a workflow step crashing, the session may not always be available
                //
                // In this case, if the wf exec raises an error and the session is not found,
                // make a dummy session so the UI is able to recover without needing to
                // refresh the page.
                if (wfExecErr && err.meta && err.meta.resource === 'automation:session' && err.meta.type === 'notFound') {
                  sessionHandler({ completedAt: new Date(), status: 'failed', error: wfExecErr })
                  return
                }
                throw err
              }).catch(this.toastErrorHandler(this.$t('notification:failed-test')))
          }

          setTimeout(sessionReader, 1000)
        }).catch(this.toastErrorHandler(this.$t('notification:failed-test')))
        .finally(() => {
          // Reset state and refresh the trigger label so spinner disappears
          this.dryRun.lookup = true
          this.dryRun.processing = false
          this.dryRun.sessionID = undefined
          this.redrawLabel(this.graph.model.getCell(this.dryRun.cellID).mxObjectId)
        })
    },

    cancelWorkflow () {
      const { sessionID, processing } = this.dryRun
      if (processing && sessionID) {
        this.dryRun.sessionID = undefined
        this.redrawLabel(this.graph.model.getCell(this.dryRun.cellID).mxObjectId)

        this.$AutomationAPI.sessionCancel({ sessionID })
          .then(() => {
            this.toastInfo('Workflow test canceled', 'Stopping test')
          })
          .catch(e => {
            this.dryRun.sessionID = sessionID
            this.toastErrorHandler('Test cancel failed')(e)
          })
      }
    },

    render (workflow, initial = false) {
      this.rendering = true

      if (this.sidebar.show) {
        this.sidebarClose()
      }

      const { x = originPoint, y = originPoint } = this.graph.view.translate
      const { scale } = this.graph.view

      if (!this.workflow.steps) {
        this.workflow.steps = []
      }

      if (!this.workflow.paths) {
        this.workflow.paths = []
      }

      // Add triggers to steps/paths
      this.triggers.forEach(({ meta, ...config }) => {
        this.workflow.steps.push({
          stepID: meta.visual.id,
          kind: 'trigger',
          defaultName: meta.visual.defaultName || false,
          meta,
        })

        meta.visual.edges.forEach(edge => {
          this.workflow.paths.push(edge)
        })
      })

      // Assemble issues
      this.issues = {}
      if (this.workflow.issues) {
        this.workflow.issues.forEach(({ culprit, description }) => {
          if (culprit) {
            const { step = -1, trigger = -1 } = culprit
            let stepID = ''

            if (step >= 0) {
              stepID = (this.workflow.steps[step] || {}).stepID
            } else if (trigger >= 0) {
              stepID = (this.triggers[trigger] || {}).meta?.visual?.id || ''
            }

            if (stepID) {
              this.issues[stepID] ? this.issues[stepID].push(description) : this.issues[stepID] = [description]
            }
          }
        })
      }

      this.deferred = false
      this.triggersPathsChanged = false

      const steps = workflow.steps || []
      const paths = workflow.paths || []
      const root = this.graph.getDefaultParent()

      this.vertices = {}
      this.edges = {}

      if (initial) {
        this.graph.view.rendering = false
      }

      this.graph.getModel().clear()

      this.graph.getModel().beginUpdate() // Adds cells to the model in a single step

      try {
        // Add vertices
        steps.sort((a, b) => a.meta.visual.parent - b.meta.visual.parent)
          .forEach(({ meta = {}, ...config }) => {
            const node = (meta || {}).visual
            if (node) {
              node.parent = this.graph.model.getCell(node.parent) || root

              const { width, height, style } = getStyleFromKind(config)

              const newCell = this.graph.insertVertex(node.parent, node.id, node.value, node.xywh[0], node.xywh[1], node.xywh[2] || width, node.xywh[3] || height, style)
              this.addCellToVertices(newCell)

              // Only set if not yet true
              this.deferred = this.deferred || this.deferredKinds.includes(config.kind)
            }
          })

        // Add edges
        paths.forEach(({ meta, ...config }) => {
          const edge = (meta || {}).visual
          if (edge) {
            edge.parent = this.graph.model.getCell(edge.parent) || root
            edge.source = config.parentID || edge.source
            edge.target = config.childID || edge.target

            const newEdge = this.graph.insertEdge(edge.parent, edge.id, edge.value, this.vertices[edge.source].node, this.vertices[edge.target].node, edge.style)
            newEdge.geometry.points = (edge.points || []).map(({ x, y }) => {
              return new mxPoint(x, y)
            })

            this.edges[edge.id] = {
              node: newEdge,
              config,
            }
          }
        })

        // Updates vertices now that edges are present
        Object.keys(this.vertices).forEach(vID => this.updateVertexConfig(vID))
      } finally {
        this.graph.view.scale = scale
        this.graph.view.setTranslate(x || originPoint, y || originPoint)

        if (this.undoManager && initial) {
          this.undoManager.clear()
        }

        this.graph.getModel().endUpdate() // Updates the display

        // Resolves problems with same id's being reused
        this.graph.getModel().nextId = this.graph.getModel().nextId + 1

        if (initial) {
          this.graph.fit()
          this.graph.view.rendering = true
          this.graph.refresh()

          this.graph.view.setTranslate(this.graph.view.translate.x + 79, this.graph.view.translate.y + 220)
          this.zoomLevel = this.graph.view.scale
        }

        this.rendering = false
      }
    },

    renderTrace (firstStepID, trace = []) {
      const cells = {}

      // Build cells object for easier drawing of overlay
      trace.filter(t => t).forEach(({ stepID, parentID, stepTime, error = false }, index) => {
        const cell = {
          index,
          stepID,
          parentID,
          stepTime,
          error,
        }

        if (cells[stepID]) {
          cells[stepID].push(cell)
        } else {
          cells[stepID] = [cell]
        }
      })

      this.highlights = []

      // Handle first cell & edge
      this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--success)', 2)) - 1].highlight(this.graph.view.getState(this.graph.model.getCell(this.dryRun.cellID)))
      const firstEdge = this.graph.model.getEdgesBetween(this.graph.model.getCell(this.dryRun.cellID), this.graph.model.getCell(firstStepID), true)[0]
      if (firstEdge) {
        this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--success)', 2)) - 1].highlight(this.graph.view.getState(firstEdge))
      }

      // Handle others
      Object.entries(cells).forEach(([stepID, frames]) => {
        if (stepID !== '0') {
          let error = frames[0].error
          let log = `#${frames[0].index + 1} - ${frames[0].stepTime}ms${error ? this.$t('notification:error') + error : ''}`
          if (frames.length < 2) {
            const [cell] = frames
            // If first cell, dont paint parent edge
            if (cell && cell.index !== 0) {
              this.graph.model.getEdgesBetween(this.graph.model.getCell(cell.parentID), this.graph.model.getCell(stepID), true)
                .forEach(edge => {
                  this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--success)', 2)) - 1].highlight(this.graph.view.getState(edge))
                })
            }
          } else {
            // If step is visited multiple times, keep track of execution info
            const time = {
              min: frames[0].stepTime,
              max: frames[0].stepTime,
              avg: 0,
              sum: 0.0,
            }

            error = ''

            frames.forEach(({ index, parentID, stepTime, error }, i) => {
              if (i !== 0) {
                if (stepTime < time.min) {
                  time.min = stepTime
                }

                if (stepTime > time.max) {
                  time.max = stepTime
                }

                log = `${log}<br>#${index + 1} - ${stepTime}ms${error ? this.$t('notification:error') + error : ''}`
              }

              time.sum += stepTime
              this.graph.model.getEdgesBetween(this.graph.model.getCell(parentID), this.graph.model.getCell(stepID), true)
                .forEach(edge => {
                  this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--success)', 2)) - 1].highlight(this.graph.view.getState(edge))
                })
            })

            time.avg = time.sum ? (time.sum / frames.length).toFixed(2) : time.sum
            log = `${log}<br><br>MIN: ${time.min}<br>MAX: ${time.max}<br>AVG: ${time.avg}<br>SUM: ${time.sum}`
          }

          // Set info overlay
          const time = new mxCellOverlay(new mxImage(this.getIcon(`clock-${error ? 'danger' : 'success'}`), 16, 16), `<span>${log}</span>`)
          this.graph.addCellOverlay(this.graph.model.getCell(stepID), time)

          // Highlight cell based on error
          if (error) {
            this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--danger)', 2)) - 1].highlight(this.graph.view.getState(this.graph.model.getCell(stepID)))
          } else {
            this.highlights[this.highlights.push(new mxCellHighlight(this.graph, 'var(--success)', 2)) - 1].highlight(this.graph.view.getState(this.graph.model.getCell(stepID)))
          }
        }
      })
    },

    getJsonModel () {
      return encodeGraph(this.graph.getModel(), this.vertices, this.edges)
    },

    importJSON (workflows = []) {
      try {
        this.importProcessing = true

        // Only render first workflow
        const [workflow] = workflows

        // Replace triggers
        this.triggers = workflow.triggers || []

        // Replace workflow steps and paths
        this.workflow = {
          ...this.workflow,
          steps: workflow.steps || [],
          paths: workflow.paths || [],
        }

        // Fresh render
        this.render(this.workflow)

        this.importProcessing = false
        this.$root.$emit('change-detected')
        this.$bvModal.hide('import')
        this.toastSuccess(this.$t('notification:imported-workflow'))
      } catch (e) {
        this.toastErrorHandler(this.$t('notification:import-failed'))(e)
      }
    },

    saveWorkflow () {
      // Just emit, let parent component take care of permission checks
      this.$emit('save', { ...this.workflow, ...this.getJsonModel() })
    },

    async getFunctionTypes () {
      return this.$AutomationAPI.functionList()
        .then(({ set }) => {
          this.functionTypes = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-functions')))
    },

    async getEventTypes () {
      return this.$AutomationAPI.eventTypesList()
        .then(({ set }) => {
          this.eventTypes = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:event-type-fetch-failed')))
    },

    getIcon (icon, mode = 'light') {
      if (!icon) {
        return ''
      }

      return `${mxClient.imageBasePath}/${mode === 'dark' ? 'dark/' : ''}${icon}.svg`
    },
  },
}
</script>

<style lang="scss" scoped>
#workflow-editor {
  color: var(--dark);
}

#graph {
  outline: none;
}

.toolbar {
  background-color: var(--sidebar-bg) !important;
  width: 66px;
}

.zoom {
  right: 0;
  bottom: 0;
}

.component-fade-enter-active, .component-fade-leave-active {
  transition: opacity 0.3s ease;
}

.component-fade-enter, .component-fade-leave-to {
  opacity: 0;
}

// https://stackoverflow.com/a/40991531/17926309
.saving::after {
  display: inline-block;
  animation: saving steps(1, end) 1s infinite;
  content: '';
}

@keyframes saving {
  0% { content: ''; }
  25% { content: '.'; }
  50% { content: '..'; }
  75% { content: '...'; }
  100% { content: ''; }
}
</style>

<style>
.hide {
  display: none;
}

.step:hover .hide {
  display: flex;
}

.show {
  display: flex;
}

.step:hover .show {
  display: none;
}

.hide-label {
  display: none;
}

.step:hover .hide-label {
  text-align: justify;
  display: flex;
}

.id-label {
  position: absolute;
  font-size: 8px;
  top: 4px;
  right: 4px;
}

.hover-untruncate {
  text-align: left;
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.step:hover .hover-untruncate {
  overflow: visible;
}

.label-container {
  overflow: hidden;
}

.step:hover .label-container {
  overflow-x: visible;
}

.step-values {
  position: absolute;
  min-width: 200px;
  top: 80px;
  border-top: 0;
}

.step-values td, th {
  text-align: left;
  padding: 8px;
  white-space: nowrap;
}

.step-values tr.title {
  background-color: var(--light) !important;
}

.step-values tr.title th {
  border-top: none;
}

#toolbar > hr {
  margin: 0.5rem 0 0.5rem 0 !important;
  align-self: stretch;
}
</style>
