<template>
  <div id="editor" ref="editor" class="flex w-full h-full">
    <Teleport to="#topbar-title">
      {{ workflow.meta.name || workflow.handle }}
    </Teleport>

    <Teleport to="#topbar-tools">
      <Button
        data-test-id="button-configure-workflow"
        :label="$t('configurator.configuration')"
        icon="pi pi-cog"
        size="small"
        @click="configuratorVisible = true"
      />
    </Teleport>

    <!-- Toolbar strip -->
    <div class="toolbar flex flex-col h-full bg-topbar">
      <div class="flex flex-col items-center mt-1 overflow-auto gap-1 px-1">
        <template v-for="(item, idx) in toolbarItems" :key="idx">
          <hr v-if="item.kind === 'hr'" class="w-full my-1 border-surface" />
          <div
            v-else
            class="toolbar-item"
            draggable="true"
            @dragstart="onToolbarDragStart($event, item)"
            @mouseenter="showToolbarTooltip($event, item)"
            @mouseleave="activeToolbarTooltip = null"
          >
            <img
              v-if="item.iconSrc"
              :src="item.iconSrc"
              class="w-8 h-8 cursor-pointer"
              :alt="item.label"
            />
          </div>
        </template>
      </div>

      <div class="flex flex-grow-1 items-end justify-center py-3">
        <Button
          ref="help"
          text
          rounded
          severity="secondary"
          class="p-2"
          @click="helpVisible = true"
        >
          <i class="pi pi-question-circle text-xl text-primary" />
        </Button>
      </div>
    </div>

    <!-- Main canvas area -->
    <div class="w-full h-full border-0 shadow-sm rounded-none relative">
      <!-- Workflow info overlay -->
      <div v-if="workflow.meta" class="absolute pl-2 pt-2" style="z-index: 1">
        <p
          v-if="workflow.meta.description"
          :class="{ 'mb-2': getRunAs }"
          class="mb-0 truncate"
          style="white-space: pre-line; max-height: 48px"
        >
          {{ workflow.meta.description }}
        </p>

        <p v-if="getRunAs" class="mb-0 truncate">
          <b>{{ $t('editor.run-as') }}</b>
          <samp>{{ getRunAs }}</samp>
        </p>

        <!-- Namespace/Module Labels -->
        <div v-if="workflowLabelsDisplay.length > 0" class="mb-2">
          <div
            v-for="group in workflowLabelsDisplay"
            :key="'group-' + group.namespaceID"
            class="flex items-center flex-wrap gap-1 mb-1"
          >
            <Tag
              v-tooltip.top="$t('general.filter.namespace.label')"
              :value="group.namespaceName"
              severity="primary"
            />

            <Tag
              v-for="mod in group.modules"
              :key="'mod-' + group.namespaceID + '-' + mod.id"
              v-tooltip.top="$t('general.filter.module.label')"
              :value="mod.name"
              severity="secondary"
            />
          </div>
        </div>

        <div class="flex items-center mb-1 gap-1">
          <Tag v-if="workflow.deletedAt" :value="$t('editor.deleted')" severity="danger" />
          <Tag v-if="!workflow.enabled" :value="$t('editor.disabled')" severity="danger" />
          <Tag v-if="hasIssues" :value="$t('editor.detected-issues')" severity="danger" />
          <Tag
            v-if="workflow.meta.subWorkflow"
            :value="$t('general.subworkflow')"
            severity="info"
          />
          <Tag v-if="deferred" :value="$t('editor.deferred')" severity="info" />
          <Tag
            v-if="triggersPathsChanged"
            :value="$t('notification.trigger-paths-changed')"
            severity="warn"
          />
        </div>
      </div>

      <!-- Bottom bar: save + zoom -->
      <div
        class="flex flex-wrap absolute bottom-0 left-0 p-2 gap-2 w-full transition-[padding-right] duration-300 ease-in-out"
        :style="{ zIndex: 1, paddingRight: sidebar.show ? `${drawerWidth + 8}px` : undefined }"
      >
        <Button
          v-if="changeDetected && canUpdateWorkflow"
          data-test-id="button-save-workflow"
          :label="
            $t('editor.detected-changes') + `${canUpdateWorkflow ? $t('editor.click-to-save') : ''}`
          "
          :loading="processingSave"
          severity="primary"
          class="rounded py-2 px-3"
          style="min-width: 20rem"
          @click="saveWorkflow()"
        />

        <div
          class="flex items-center bg-surface border border-surface py-2 px-3 gap-3 rounded ml-auto"
        >
          <Button
            v-tooltip.top="'Undo (Ctrl+Z)'"
            text
            severity="secondary"
            class="p-0"
            :disabled="!canUndo"
            @click="undo(); refreshSidebar()"
          >
            <i class="pi pi-undo" />
          </Button>
          <Button
            v-tooltip.top="'Redo (Ctrl+Shift+Z)'"
            text
            severity="secondary"
            class="p-0"
            :disabled="!canRedo"
            @click="redo(); refreshSidebar()"
          >
            <i class="pi pi-undo" style="transform: scaleX(-1); display: inline-block" />
          </Button>
        </div>

        <div
          class="flex items-center bg-surface border border-surface py-2 px-3 gap-1 rounded"
          style="z-index: 1"
        >
          {{ getZoomPercent }}
          <Button text severity="secondary" class="ml-4 p-0" @click="zoom(false)">
            <i class="pi pi-search-minus" />
          </Button>
          <Button text severity="secondary" class="ml-1 p-0" @click="zoom()">
            <i class="pi pi-search-plus" />
          </Button>
          <Button
            :label="$t('editor.reset')"
            text
            severity="secondary"
            class="ml-2 p-0"
            @click="resetZoom()"
          />
        </div>
      </div>

      <!-- VueFlow canvas -->
      <VueFlow
        :id="vfId"
        v-model:nodes="nodes"
        v-model:edges="edges"
        :snap-to-grid="true"
        :snap-grid="[8, 8]"
        :min-zoom="0.1"
        :max-zoom="3"
        :pan-on-drag="[1, 2]"
        :pan-on-scroll="true"
        :zoom-on-scroll="false"
        :zoom-on-double-click="false"
        :selection-mode="SelectionMode.Partial"
        :selection-key-code="true"
        :multi-selection-key-code="'Shift'"
        :delete-key-code="null"
        :connection-mode="ConnectionMode.Loose"
        :edges-updatable="true"
        :is-valid-connection="isValidConnection"
        :default-edge-options="defaultEdgeOptions"
        connection-line-type="smoothstep"
        class="vueflow-canvas"
        @connect="onConnect"
        @edge-update="onEdgeUpdate"
        @edge-update-start="onEdgeUpdateStart"
        @edge-update-end="onEdgeUpdateEnd"
        @node-click="onNodeClick"
        @edge-click="onEdgeClick"
        @pane-click="onPaneClick"
        @node-drag="onNodeDrag"
        @node-drag-stop="onNodeDragStop"
        @dragover="onCanvasDragOver"
        @drop="onCanvasDrop"
      >
        <Background variant="dots" />
        <template #node-workflow="props">
          <WorkflowNode
            v-bind="props"
            :issues="issues"
            :function-types="functionTypes"
            :event-types="eventTypes"
            :current-theme="currentTheme"
            :used-source-handles="[...nodeConn(props.id).sources]"
            :used-target-handles="[...nodeConn(props.id).targets]"
            :out-count="nodeConn(props.id).outCount"
            @open-issues="openIssuesModal"
          />
        </template>
        <template #node-trigger="props">
          <TriggerNode
            v-bind="props"
            :issues="issues"
            :event-types="eventTypes"
            :can-test="!!workflow.canExecuteWorkflow"
            :dry-run-processing="dryRun.processing"
            :dry-run-cell-i-d="dryRun.cellID"
            :dry-run-session-i-d="dryRun.sessionID"
            :current-theme="currentTheme"
            :has-outgoing-edge="nodeConn(props.id).outCount > 0"
            @test="startTest"
            @cancel="cancelWorkflow"
            @open-issues="openIssuesModal"
          />
        </template>
        <template #node-termination="props">
          <TerminationNode
            v-bind="props"
            :current-theme="currentTheme"
            :used-target-handles="[...nodeConn(props.id).targets]"
            :in-count="nodeConn(props.id).inCount"
          />
        </template>
        <template #node-visual="props">
          <VisualNode v-bind="props" :drop-target="dropTargetSwimlaneId === props.id" />
        </template>
        <template #edge-workflow="props">
          <FlowEdge v-bind="props" @update-label="onEdgeLabelUpdate" />
        </template>
      </VueFlow>
    </div>

    <!-- Sidebar for step configurator -->
    <Transition
      enter-active-class="transition-transform duration-300 ease-in-out"
      enter-from-class="translate-x-full"
      enter-to-class="translate-x-0"
      leave-active-class="transition-transform duration-300 ease-in-out"
      leave-from-class="translate-x-0"
      leave-to-class="translate-x-full"
    >
      <div
        v-show="sidebar.show"
        class="config-drawer right-sidebar shadow-xl z-20 flex rounded-border"
        :style="{ width: `${drawerWidth}px` }"
      >
        <!-- Resize handle -->
        <div
          class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors shrink-0"
          @mousedown="startDrawerResize"
        />
        <!-- Drawer content -->
        <div class="flex-1 flex flex-col overflow-auto px-3 py-2">
          <!-- Header -->
          <div class="flex items-start justify-between pr-1 group/header">
            <div class="flex items-center gap-3 flex-1 min-w-0 pr-2">
              <div class="rounded-border flex items-center justify-center shrink-0 bg-surface">
                <img
                  v-if="getSidebarItemIcon"
                  :src="getSidebarItemIcon"
                  class="h-9 w-9 object-contain"
                />
                <i
                  v-else-if="sidebar.itemType === 'edge'"
                  class="pi pi-arrow-right-arrow-left text-color text-2xl"
                />
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-lg font-semibold text-color truncate m-0">
                  {{ getSidebarItemType }}
                </h3>
                <div v-if="getSelectedItem?.node?.id" class="text-xs text-muted-color truncate">
                  #{{ getSelectedItem.node.id }}
                </div>
              </div>
            </div>
            <Button
              icon="pi pi-times"
              text
              rounded
              size="small"
              class="shrink-0 mt-1"
              @click="sidebarClose()"
            />
          </div>

          <!-- Body -->
          <div class="relative flex-1 mt-3">
            <div v-if="sidebar.item && !isContentNode" class="flex flex-col gap-3">
              <CFormGroup :label="$t('general.label.name')" input-id="sidebar-name">
                <InputText
                  id="sidebar-name"
                  :model-value="sidebar.item.node?.value || ''"
                  class="w-full"
                  @update:model-value="onSidebarNameChange"
                />
              </CFormGroup>
              <CFormGroup :label="$t('general.description')" input-id="sidebar-description">
                <Textarea
                  id="sidebar-description"
                  :model-value="sidebar.item.node?.description || ''"
                  auto-resize
                  rows="2"
                  class="w-full"
                  @update:model-value="onSidebarDescriptionChange"
                />
              </CFormGroup>
            </div>
            <Divider v-if="sidebar.item && !isContentNode" />
            <transition name="component-fade" mode="out-in">
              <configurator
                v-if="sidebar.showItem"
                v-model:item="sidebar.item"
                v-model:edges="sidebarEdges"
                :out-edges="sidebarOutEdges"
                :is-subworkflow="!!workflow.meta.subWorkflow"
                @update-value="setValue($event)"
                @update-default-value="setValue($event, true)"
              />
            </transition>
          </div>

          <!-- Footer -->
          <div class="mt-6 pt-3 border-t border-surface flex flex-col gap-2">
            <CInputDelete
              :label="$t('general.label.delete')"
              :message="$t('notification.delete-confirmation')"
              :header="getSelectedItem?.node?.value || getSidebarItemType"
              outlined
              size="small"
              @confirm="sidebarDelete()"
            />
            <div id="sidebar-footer" />
          </div>
        </div>
      </div>
    </Transition>

    <!-- Workflow Configurator Dialog -->
    <Dialog
      v-model:visible="configuratorVisible"
      :header="
        workflow.workflowID === '0' ? $t('general.new-workflow') : $t('general.edit-workflow')
      "
      modal
      :pt="{ content: 'p-0' }"
      class="w-full max-w-3xl"
    >
      <workflow-configurator
        v-if="workflow.workflowID"
        :workflow="workflow"
        :can-create="canCreate"
        :processing-save="processingSave"
        :processing-delete="processingDelete"
        :import-processing="importProcessing"
        @save="handleWorkflowSave"
        @import="importJSON"
        @delete="$emit('delete')"
        @undelete="$emit('undelete')"
        @close="configuratorVisible = false"
      />
    </Dialog>

    <Dialog
      v-model:visible="helpVisible"
      :header="$t('editor.help')"
      modal
      class="w-full max-w-4xl"
    >
      <help />
    </Dialog>

    <Dialog
      v-model:visible="issuesModal.show"
      :header="$t('editor.issues')"
      modal
      class="w-full max-w-xl"
    >
      <div v-for="(issue, index) in issuesModal.issues" :key="index">
        <p>
          <code>{{ issue[0].toUpperCase() + issue.slice(1) }}</code>
        </p>
      </div>
    </Dialog>

    <!-- Dry Run Dialog -->
    <Dialog
      v-model:visible="dryRun.show"
      :header="$t('editor.initial-scope')"
      modal
      class="w-full max-w-3xl"
    >
      <div class="flex flex-col gap-4">
        <div v-if="dryRun.lookup">
          <small class="text-muted-color">
            {{ $t('editor.input-ids-or-handles') }}
            <br />
            {{ $t('editor.modify-initial-scope-if-no-variables-are-loaded') }}
            <br />
            {{ $t('editor.auto-initialize-empty-variable') }}
            <br />
            <br />
            {{ $t('editor.open-webapp-on-prompt-use') }}
          </small>
          <div v-for="(p, index) in Object.values(dryRun.initialScope)" :key="index" class="mt-4">
            <div v-if="p.lookup" class="flex flex-col gap-1">
              <label class="font-medium text-sm">{{ p.label }}</label>
              <InputText v-model="p.value" class="w-full" />
              <small v-if="p.description" class="text-muted-color">{{ p.description }}</small>
            </div>
          </div>
        </div>
        <div v-else>
          <Textarea
            :model-value="JSON.stringify(dryRun.input, null, 2)"
            rows="15"
            class="w-full font-mono text-sm"
            @update:model-value="onDryRunEdit(JSON.parse($event || '{}'))"
          />
        </div>
      </div>

      <template #footer>
        <div class="flex items-center w-full">
          <Button
            v-if="!dryRun.lookup"
            :label="$t('editor.back')"
            severity="secondary"
            text
            size="small"
            @click="dryRun.lookup = true"
          />
          <div class="flex-1" />
          <Button
            :label="dryRun.lookup ? $t('editor.load-and-configure') : $t('editor.run-workflow')"
            severity="success"
            size="small"
            @click="dryRunOk"
          />
        </div>
      </template>
    </Dialog>

    <!-- Sidebar Rich Tooltip -->
    <div
      v-if="activeToolbarTooltip"
      class="fixed z-50 pointer-events-none transition-opacity"
      :style="toolbarTooltipStyle"
    >
      <div
        class="bg-surface border border-surface shadow-lg rounded-xl overflow-hidden w-64 flex flex-col"
      >
        <div class="p-4 flex flex-col gap-2">
          <h4 class="font-semibold text-lg text-color m-0">{{ activeToolbarTooltip.title }}</h4>
          <p class="text-sm text-color-secondary m-0 leading-snug">
            {{ activeToolbarTooltip.tooltip }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { VueFlow, ConnectionMode, MarkerType, SelectionMode, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'

import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick, inject } from 'vue'
import { debounce } from 'lodash-es'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'

import { decodeWorkflow, encodeWorkflow } from '../lib/codec'
import { getStyleFromKind } from '../lib/style'
import { encodeInput } from '../lib/dry-run'
import toolbarConfig from '../lib/toolbar'
import eventBus from '../lib/eventBus'
import { nextId } from '../lib/id'
import { getIcon } from '../lib/icon'
import { NoID } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
const { CInputDelete } = components

import Configurator from './Configurator/index.vue'
import WorkflowConfigurator from './Configurator/Workflow.vue'
import Help from './Help.vue'
import WorkflowNode from './FlowNodes/WorkflowNode.vue'
import TriggerNode from './FlowNodes/TriggerNode.vue'
import TerminationNode from './FlowNodes/TerminationNode.vue'
import VisualNode from './FlowNodes/VisualNode.vue'
import FlowEdge from './FlowEdge.vue'

import { useWorkflowHistory } from '../composables/useWorkflowHistory'
import { useWorkflowDnD } from '../composables/useWorkflowDnD'
import { useWorkflowHighlight } from '../composables/useWorkflowHighlight'
import { useWorkflowClipboard } from '../composables/useWorkflowClipboard'
import { useLabelsStore } from '../stores/labels'
import { isValidConnection as validateConnection } from '../lib/connectionRules'

const { t } = useI18n()
const toast = useToast()

const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $Auth = inject('$Auth')

/* ─── Props & Emits ─── */
const props = defineProps({
  workflowObject: { type: Object, default: () => ({}) },
  workflowTriggers: { type: Array, default: () => [] },
  changeDetected: { type: Boolean },
  canCreate: { type: Boolean },
  processingSave: { type: Boolean },
  processingDelete: { type: Boolean },
})

const emit = defineEmits(['save', 'change-detected', 'delete', 'undelete'])

/* ─── Reactive State ─── */
const initialized = ref(false)
const helpVisible = ref(false)
const configuratorVisible = ref(false)
const deferred = ref(false)
const triggersPathsChanged = ref(false)
const importProcessing = ref(false)

const workflow = ref({})
const triggers = ref([])
const issues = ref({})
const runAsUser = ref(undefined)

const nodes = ref([])
const edges = ref([])

const eventTypes = ref([])
const functionTypes = ref([])

const deferredKinds = ['delay', 'prompt']
const labelsStore = useLabelsStore()

const activeToolbarTooltip = ref(null)
const toolbarTooltipStyle = ref({ top: '0px', left: '0px' })

const drawerWidth = ref(380)
const isResizingDrawer = ref(false)

const editor = ref(null)

/* ─── Composables ─── */
const vfId = 'workflow-editor-flow'
const {
  fitView,
  zoomIn: vfZoomIn,
  zoomOut: vfZoomOut,
  getViewport,
  project,
  onNodesInitialized,
} = useVueFlow(vfId)

const defaultEdgeOptions = { markerEnd: MarkerType.ArrowClosed }

const { saveToHistory, resetHistory, undo, redo, canUndo, canRedo } = useWorkflowHistory(nodes, edges)
const debouncedSaveToHistory = debounce(saveToHistory, 600)
const {
  onToolbarDragStart,
  onCanvasDragOver,
  onCanvasDrop: dndDrop,
} = useWorkflowDnD(nodes, edges, saveToHistory, project)
const { highlightConnected, clearHighlights } = useWorkflowHighlight(nodes, edges)
const { copySelected, cutSelected, pasteClipboard } = useWorkflowClipboard(
  nodes,
  edges,
  saveToHistory,
)

const zoomLevel = ref(1)

// Fit view once on first load
let hasInitiallyFit = false
onNodesInitialized(() => {
  if (!hasInitiallyFit && nodes.value.length > 0) {
    fitView({ padding: 0.2, maxZoom: 1 })
    hasInitiallyFit = true
    zoomLevel.value = getViewport().zoom || 1
  }
})

/* ─── Sidebar State ─── */
const sidebar = ref({
  item: undefined,
  itemType: undefined,
  show: false,
  showItem: false,
})

// Live out-edge count for the currently focused node. Recomputes as edges change,
// so the Configurator stays in sync when connections are added/removed while open.
const sidebarOutEdges = computed(() => {
  const id = sidebar.value.item?.node?.id
  if (!id) return 0
  return edges.value.filter(e => e.source === id).length
})

// Tracks the edge currently being reconnected so its endpoints can be offered
// again during the drag (otherwise the used-handle masks hide them).
const edgeUpdatingId = ref(null)

// Per-node map of used source/target handle ids — drives the FlowNodes' handle
// visibility so a point that's already wired up doesn't offer a second connection
// at most one edge per exit/entry point.
const nodeConnections = computed(() => {
  const map = {}
  const ensure = id => {
    if (!map[id]) map[id] = { sources: new Set(), targets: new Set(), outCount: 0, inCount: 0 }
    return map[id]
  }
  edges.value.forEach(e => {
    // Pretend the edge being reconnected isn't there yet so both endpoints can
    // light up handles (including on its current node).
    if (e.id === edgeUpdatingId.value) return
    if (e.source) {
      const s = ensure(e.source)
      if (e.sourceHandle) s.sources.add(e.sourceHandle)
      s.outCount++
    }
    if (e.target) {
      const t = ensure(e.target)
      if (e.targetHandle) t.targets.add(e.targetHandle)
      t.inCount++
    }
  })
  return map
})

function nodeConn(id) {
  return (
    nodeConnections.value[id] || { sources: new Set(), targets: new Set(), outCount: 0, inCount: 0 }
  )
}

const issuesModal = ref({
  show: false,
  issues: [],
})

const dryRun = ref({
  show: false,
  processing: false,
  lookup: false,
  cellID: undefined,
  initialScope: {},
  input: {},
  inputEdited: {},
  sessionID: undefined,
})

/* ─── Sidebar Edges Bridge ─── */
// The Configurator expects edges as an object map { [id]: { node, config } }
// We need to bridge that with VueFlow's edges array.
// IMPORTANT: Gateway.vue directly mutates `this.edges[id].config.expr` so we
// must use a reactive ref + deep watcher rather than a computed getter/setter.
const sidebarEdges = ref({})

// Rebuild the bridge map whenever VueFlow edges change
watch(
  edges,
  () => {
    const map = {}
    edges.value.forEach(e => {
      map[e.id] = {
        node: {
          id: e.id,
          value: e.label || '',
          source: {
            id: e.source,
            style: nodes.value.find(n => n.id === e.source)?.data?.kind || '',
          },
          target: {
            id: e.target,
            style: nodes.value.find(n => n.id === e.target)?.data?.kind || '',
          },
        },
        config: {
          parentID: e.source,
          childID: e.target,
          expr: e.data?.expr || '',
        },
      }
    })
    sidebarEdges.value = map
  },
  { immediate: true, deep: true },
)

// Sync mutations from Configurator back to VueFlow edges (deep watch catches
// direct property mutations like Gateway.vue's `this.edges[id].config.expr = expr`)
watch(
  sidebarEdges,
  newMap => {
    Object.entries(newMap).forEach(([edgeId, edgeData]) => {
      const edge = edges.value.find(e => e.id === edgeId)
      if (edge && edgeData.config) {
        const newExpr = edgeData.config.expr || ''
        const newLabel = edgeData.node?.value || ''
        let changed = false
        if (edge.data?.expr !== newExpr) {
          edge.data = { ...edge.data, expr: newExpr }
          changed = true
        }
        if (edge.label !== newLabel) {
          edge.label = newLabel
          changed = true
        }
        if (changed) debouncedSaveToHistory()
      }
    })
  },
  { deep: true },
)

// Sync Configurator sidebar item mutations back to VueFlow nodes.
// Configurators write freely to sidebar.item.config.* (per-kind fields like
// delay duration, prompt payload, iterator iter/key, exec-workflow target,
// gateway paths, expressions, etc.) — merge the full config back so nothing
// is lost at save time. Runtime-only fields (trace, highlight) are excluded
// so a re-render doesn't stomp them.
watch(
  () => sidebar.value.item,
  newItem => {
    if (!newItem?.node?.id) return

    if (sidebar.value.itemType === 'edge') {
      const edge = edges.value.find(e => e.id === newItem.node.id)
      if (!edge) return
      const nextLabel = newItem.node.value || ''
      const nextDesc = newItem.node.description || ''
      let edgeChanged = false
      if (edge.label !== nextLabel) { edge.label = nextLabel; edgeChanged = true }
      if ((edge.data?.description || '') !== nextDesc) {
        edge.data = { ...edge.data, description: nextDesc }
        edgeChanged = true
      }
      if (edgeChanged) debouncedSaveToHistory()
      return
    }

    const node = nodes.value.find(n => n.id === newItem.node.id)
    if (!node) return

    const {
      highlighted: _highlighted,
      traceState: _traceState,
      traceLog: _traceLog,
      ...configData
    } = newItem.config || {}
    const merged = { ...node.data, ...configData }

    if (newItem.triggers) merged.triggers = newItem.triggers
    if (newItem.node.value !== undefined) merged.label = newItem.node.value
    if (newItem.node.description !== undefined) merged.description = newItem.node.description

    // Compare only the keys that the sidebar actually wrote — avoids
    // stringifying the entire node.data (including large arguments/results
    // arrays) on every keystroke. Keys from node.data not in the patch are
    // identical by construction (they came from ...node.data spread).
    const patchKeys = [
      ...Object.keys(configData),
      ...(newItem.triggers ? ['triggers'] : []),
      ...(newItem.node.value !== undefined ? ['label'] : []),
      ...(newItem.node.description !== undefined ? ['description'] : []),
    ]
    const hasChange = patchKeys.some(
      k => JSON.stringify(node.data[k]) !== JSON.stringify(merged[k]),
    )
    if (hasChange) {
      node.data = merged
      debouncedSaveToHistory()
    }
  },
  { deep: true },
)

/* ─── Computed ─── */
const currentTheme = computed(() => $Auth?.user?.meta?.theme || 'light')

const canUpdateWorkflow = computed(() =>
  workflow.value.workflowID === '0' ? props.canCreate : workflow.value.canUpdateWorkflow,
)

const hasIssues = computed(() => (workflow.value.issues || []).length)

const getRunAs = computed(() => {
  if (runAsUser.value) {
    const { userID, name, username, email } = runAsUser.value
    return name || username || email || `<@${userID}>`
  }
  return undefined
})

const getZoomPercent = computed(() => `${Math.floor(zoomLevel.value * 100).toFixed(0)}%`)

/* ─── Sidebar computed ─── */
const isContentNode = computed(
  () =>
    sidebar.value.item?.config?.kind === 'visual' && sidebar.value.item?.config?.ref === 'content',
)

const getSidebarItemType = computed(() => {
  const item = sidebar.value.item
  if (!item) return ''
  if (item.config?.kind === 'edge' || item.node?.edge) {
    return t('steps.path.short')
  }
  if (isContentNode.value) return 'Text'
  const style = getStyleFromKind(item.config || {})?.style || item.config?.kind || ''
  return t(`steps.${style}.short`, style)
})

const getSidebarItemIcon = computed(() => {
  const item = sidebar.value.item
  if (item?.config) {
    const styleInfo = getStyleFromKind(item.config)
    return styleInfo?.icon ? getIcon(styleInfo.icon, currentTheme.value) : undefined
  }
  return undefined
})

const getSelectedItem = computed(() => sidebar.value.item || undefined)

/* ─── Labels display ─── */
const workflowLabelsDisplay = computed(() => {
  const namespaceIDs = []
  const modulesByNamespace = {}
  if (!workflow.value.labels) return []

  if (workflow.value.labels.ref_namespace) {
    const nsValues = Array.isArray(workflow.value.labels.ref_namespace)
      ? workflow.value.labels.ref_namespace
      : [workflow.value.labels.ref_namespace]
    nsValues.forEach(label => {
      const nsID = label.split('/')[1]
      if (nsID && !namespaceIDs.includes(nsID)) {
        namespaceIDs.push(nsID)
        labelsStore.resolveNamespace({ namespaceID: nsID, api: $ComposeAPI })
      }
    })
  }

  if (workflow.value.labels.ref_module) {
    const modValues = Array.isArray(workflow.value.labels.ref_module)
      ? workflow.value.labels.ref_module
      : [workflow.value.labels.ref_module]
    modValues.forEach(label => {
      const parts = label.split('/')
      const nsID = parts[1]
      const modID = parts[2]
      if (!nsID || !modID) return
      if (!namespaceIDs.includes(nsID)) {
        namespaceIDs.push(nsID)
        labelsStore.resolveNamespace({ namespaceID: nsID, api: $ComposeAPI })
      }
      if (!modulesByNamespace[nsID]) modulesByNamespace[nsID] = []
      labelsStore.resolveModule({ moduleID: modID, namespaceID: nsID, api: $ComposeAPI })
      const name = labelsStore.getModule(modID)
      modulesByNamespace[nsID].push({ id: modID, name: name || modID })
    })
  }

  return namespaceIDs.map(nsID => {
    const name = labelsStore.getNamespace(nsID)
    return {
      namespaceID: nsID,
      namespaceName: name || nsID,
      modules: modulesByNamespace[nsID] || [],
    }
  })
})

/* ─── Toolbar Items ─── */
const toolbarItems = computed(() => {
  return toolbarConfig.map(item => {
    if (item.kind === 'hr') return item
    const styleInfo = getStyleFromKind(item) || {}
    const isContent = item.kind === 'visual' && item.ref === 'content'
    const label = isContent ? 'Text' : t(`steps.${styleInfo.style || item.kind}.label`, item.kind)
    const tooltip = t(`steps.${styleInfo.style || item.kind}.tooltip`)
    return {
      ...item,
      ...styleInfo,
      label,
      tooltip,
      iconSrc: styleInfo.icon ? getIcon(styleInfo.icon, currentTheme.value) : '',
    }
  })
})

/* ─── Watchers ─── */
watch(
  () => workflow.value.runAs,
  (runAs = '0') => {
    if (runAs !== '0') {
      $SystemAPI
        .userRead({ userID: runAs })
        .then(user => {
          runAsUser.value = user
        })
        .catch(e => {
          console.warn('Failed to resolve runAs user', runAs, e)
        })
    } else {
      runAsUser.value = undefined
    }
  },
  { immediate: true },
)

// Combined watcher: triggers and workflow must update atomically, otherwise
// render() decodes with a stale triggers.value (missing freshly-assigned
// triggerID) and the next save re-creates triggers server-side → duplicate.
watch(
  [() => props.workflowObject, () => props.workflowTriggers],
  ([wf, t]) => {
    if (wf.workflowID !== workflow.value.workflowID) {
      configuratorVisible.value = false
    }
    workflow.value = wf
    triggers.value = t || []
    if (initialized.value) {
      nextTick(() => render(workflow.value))
    }
  },
  { immediate: true },
)

/* ─── Lifecycle ─── */
onMounted(() => {
  getEventTypes()
  getFunctionTypes()

  eventBus.on('trigger-updated', onTriggerUpdated)
  eventBus.on('change-detected', onEventBusChange)

  render(workflow.value, true)

  if (workflow.value.workflowID && workflow.value.workflowID === '0') {
    configuratorVisible.value = true
  }

  document.addEventListener('keydown', keybinds)

  // Ctrl/Cmd + wheel → zoom. Capture phase + passive:false so we can preventDefault
  // reliably before VueFlow's internal handlers see the event.
  if (editor.value) {
    editor.value.addEventListener('wheel', onWheelZoom, { capture: true, passive: false })
  }

  initialized.value = true
})

onBeforeUnmount(() => {
  eventBus.off('trigger-updated', onTriggerUpdated)
  eventBus.off('change-detected', onEventBusChange)
  document.removeEventListener('keydown', keybinds)
  if (editor.value) {
    editor.value.removeEventListener('wheel', onWheelZoom, { capture: true })
  }
})

function onWheelZoom(event) {
  if (!(event.ctrlKey || event.metaKey)) return
  event.preventDefault()
  event.stopPropagation()
  if (event.deltaY < 0) vfZoomIn()
  else if (event.deltaY > 0) vfZoomOut()
  zoomLevel.value = getViewport().zoom || zoomLevel.value
}

function onTriggerUpdated(node) {
  // Trigger label changed in configurator — refresh the node data
  const nodeId = node?.id || node?.nodeId
  const n = nodes.value.find(n => n.id === nodeId)
  if (n) {
    n.data = { ...n.data }
  }
}

function onEventBusChange() {
  emit('change-detected')
  debouncedSaveToHistory()
}

/* ─── Render ─── */
function render(wf, initial = false) {
  if (sidebar.value.show) sidebarClose()
  clearHighlights()

  if (!wf.steps) wf.steps = []
  if (!wf.paths) wf.paths = []

  // Assemble issues
  issues.value = {}
  if (wf.issues) {
    wf.issues.forEach(({ culprit, description }) => {
      if (culprit) {
        const { step = -1, trigger: triggerIdx = -1 } = culprit
        let stepID = ''
        if (step >= 0) stepID = (wf.steps[step] || {}).stepID
        else if (triggerIdx >= 0)
          stepID = (triggers.value[triggerIdx] || {})?.meta?.visual?.id || ''
        if (stepID) {
          issues.value[stepID]
            ? issues.value[stepID].push(description)
            : (issues.value[stepID] = [description])
        }
      }
    })
  }

  deferred.value = false
  triggersPathsChanged.value = false

  // Decode
  const decoded = decodeWorkflow(wf, triggers.value)
  nodes.value = decoded.nodes
  edges.value = decoded.edges

  // Check deferred
  nodes.value.forEach(n => {
    if (deferredKinds.includes(n.data?.kind)) deferred.value = true
  })

  if (initial) {
    hasInitiallyFit = false
    resetHistory()
  }

  // Check if trigger→step edges match trigger.stepID (#12)
  nextTick(() => checkExistingTriggerPaths())
}

/* ─── VueFlow Event Handlers ─── */

function onConnect(connection) {
  // Belt-and-suspenders: isValidConnection is the primary gate, but guard here
  // too so a stray connect event (e.g. from click-to-connect paths) can't
  // slip past the outbound-cap rule.
  if (!isValidConnection(connection)) return

  const sourceNode = nodes.value.find(n => n.id === connection.source)
  const outCount = edges.value.filter(e => e.source === connection.source).length

  let label = ''
  const kind = sourceNode?.data?.kind
  const ref = sourceNode?.data?.ref

  if (kind === 'gateway') {
    if (ref === 'excl') {
      label = outCount === 0 ? '#1 - If' : `#${outCount + 1} - Else (if)`
    } else if (ref === 'incl') {
      label = 'If'
    }
  } else if (kind === 'error-handler') {
    label = outCount === 0 ? 'Try' : 'Catch'
  } else if (kind === 'iterator') {
    label = outCount === 0 ? 'Body' : 'End'
  }

  const newEdge = {
    id: String(nextId(nodes, edges)),
    source: connection.source,
    target: connection.target,
    sourceHandle: connection.sourceHandle,
    targetHandle: connection.targetHandle,
    type: 'workflow',
    label,
    markerEnd: MarkerType.ArrowClosed,
    data: {
      expr: '',
      parentID: connection.source,
      childID: connection.target,
      highlighted: false,
      traceState: null,
    },
  }

  edges.value = [...edges.value, newEdge]
  saveToHistory()
  emit('change-detected')
}

function onNodeClick({ node, event }) {
  const ctrlDown = event.ctrlKey || event.metaKey

  if (ctrlDown) {
    highlightConnected(node.id)
    return
  }

  clearHighlights()
  highlightConnected(node.id)

  // Build sidebar item
  const sidebarItem = buildSidebarItem(node)
  if (sidebarItem) {
    sidebarReopen(sidebarItem, node.data?.kind || 'workflow')
  }
}

function onEdgeClick({ edge }) {
  clearHighlights()

  // Build sidebar item for edge
  const sidebarItem = {
    node: {
      id: edge.id,
      value: edge.label || '',
      description: edge.data?.description || '',
      edge: true,
      source: {
        id: edge.source,
        style: nodes.value.find(n => n.id === edge.source)?.data?.kind || '',
      },
      target: { id: edge.target },
      edges: [],
    },
    config: {
      kind: 'edge',
      stepID: edge.id,
      parentID: edge.source,
      childID: edge.target,
      expr: edge.data?.expr || '',
    },
  }

  sidebarReopen(sidebarItem, 'edge')
}

function onPaneClick() {
  sidebar.value.show = false
  if (getSelectedItem.value) sidebarClose()
  clearHighlights()
}

// Tracks the swimlane currently under the dragged node so the swimlane can
// paint a "drop target" highlight — cleared on drag stop.
const dropTargetSwimlaneId = ref(null)

function findSwimlaneUnder(node) {
  const current = nodes.value.find(n => n.id === node.id)
  if (!current) return null
  const abs = getAbsolutePosition(current)
  const cx = abs.x + (current.data?.width || 200) / 2
  const cy = abs.y + (current.data?.height || 80) / 2
  for (let i = nodes.value.length - 1; i >= 0; i--) {
    const c = nodes.value[i]
    if (c.id === current.id) continue
    if (c.type !== 'visual' || c.data?.ref !== 'swimlane') continue
    const pAbs = getAbsolutePosition(c)
    const pw = c.data?.width || 400
    const ph = c.data?.height || 240
    if (cx >= pAbs.x && cx <= pAbs.x + pw && cy >= pAbs.y && cy <= pAbs.y + ph) return c
  }
  return null
}

function onNodeDrag({ node }) {
  if (node.type === 'visual') {
    dropTargetSwimlaneId.value = null
    return
  }
  const swim = findSwimlaneUnder(node)
  // Only highlight when the drop would actually change parenting — if the node
  // is already a child of that swimlane, no visual feedback is needed.
  const current = nodes.value.find(n => n.id === node.id)
  dropTargetSwimlaneId.value = swim && current?.parentNode !== swim.id ? swim.id : null
}

function onNodeDragStop({ node }) {
  dropTargetSwimlaneId.value = null
  // Visual swimlanes act as containers — drop detection reparents nodes whose
  // centre falls inside a swimlane and detaches nodes dragged back out.
  if (node.type !== 'visual') {
    reparentInsideSwimlane(node)
  }
  saveToHistory()
  emit('change-detected')
}

function getAbsolutePosition(node) {
  let x = node.position?.x || 0
  let y = node.position?.y || 0
  let cur = node.parentNode ? nodes.value.find(n => n.id === node.parentNode) : null
  const guard = new Set()
  while (cur && !guard.has(cur.id)) {
    guard.add(cur.id)
    x += cur.position?.x || 0
    y += cur.position?.y || 0
    cur = cur.parentNode ? nodes.value.find(n => n.id === cur.parentNode) : null
  }
  return { x, y }
}

function reparentInsideSwimlane(node) {
  const current = nodes.value.find(n => n.id === node.id)
  if (!current) return

  const abs = getAbsolutePosition(current)
  const cx = abs.x + (current.data?.width || 200) / 2
  const cy = abs.y + (current.data?.height || 80) / 2

  let newParent = null
  // Search in reverse so we pick innermost swimlane if nested
  for (let i = nodes.value.length - 1; i >= 0; i--) {
    const candidate = nodes.value[i]
    if (candidate.id === current.id) continue
    if (candidate.type !== 'visual' || candidate.data?.ref !== 'swimlane') continue
    const pAbs = getAbsolutePosition(candidate)
    const pw = candidate.data?.width || 400
    const ph = candidate.data?.height || 240
    if (cx >= pAbs.x && cx <= pAbs.x + pw && cy >= pAbs.y && cy <= pAbs.y + ph) {
      newParent = candidate
      break
    }
  }

  const desiredParent = newParent?.id
  if (current.parentNode === desiredParent) return

  // Rebuild as a new node object so VueFlow picks up the parent change.
  // `extent` is intentionally left unset so the user can later drag the node
  // back out of the swimlane without the child being clipped to parent bounds.
  const updated = {
    ...current,
    parentNode: desiredParent,
    extent: undefined,
  }

  // Recompute local position: absolute → new-parent-relative (or absolute).
  if (desiredParent) {
    const pAbs = getAbsolutePosition(newParent)
    updated.position = { x: abs.x - pAbs.x, y: abs.y - pAbs.y }
  } else {
    updated.position = { x: abs.x, y: abs.y }
  }

  nodes.value = nodes.value.map(n => (n.id === current.id ? updated : n))
}

function onCanvasDrop(event) {
  const newNode = dndDrop(event)
  if (newNode) {
    if (newNode.type !== 'visual') reparentInsideSwimlane(newNode)
    emit('change-detected')
    nextTick(() => {
      const sidebarItem = buildSidebarItem(newNode)
      if (sidebarItem) {
        sidebarReopen(sidebarItem, newNode.data?.kind || 'workflow')
      }
    })
  }
}

function onEdgeLabelUpdate({ id, label }) {
  const edge = edges.value.find(e => e.id === id)
  if (edge) {
    edge.label = label
    // Also update sidebar item if open
    if (sidebar.value.item?.node?.id === id) {
      sidebar.value.item.node.value = label
    }
    emit('change-detected')
  }
}

/* ─── Sidebar ─── */

function buildSidebarItem(node) {
  const { data = {} } = node
  const outEdges = edges.value.filter(e => e.source === node.id)

  const {
    highlighted: _highlighted,
    traceState: _traceState,
    traceLog: _traceLog,
    width: _width,
    height: _height,
    label: _label,
    description: _description,
    triggers: _triggers,
    ...configData
  } = data

  return {
    node: {
      id: node.id,
      value: data.label || '',
      description: data.description || '',
      edges: outEdges.map(e => ({
        id: e.id,
        value: e.label || '',
        source: {
          id: e.source,
          style: nodes.value.find(n => n.id === e.source)?.data?.kind || '',
        },
        target: {
          id: e.target,
          style: nodes.value.find(n => n.id === e.target)?.data?.kind || '',
        },
      })),
      style: getStyleFromKind(data)?.style || data.kind || '',
    },
    config: {
      ...configData,
      stepID: node.id,
      kind: data.kind || '',
      ref: data.ref || '',
      defaultName: data.defaultName || false,
      arguments: JSON.parse(JSON.stringify(data.arguments || [])),
      results: JSON.parse(JSON.stringify(data.results || [])),
    },
    triggers: data.triggers || undefined,
  }
}

function sidebarReopen(item, itemType) {
  if (!sidebar.value.show) {
    sidebar.value.item = item
    sidebar.value.itemType = itemType
    sidebar.value.show = true
    sidebar.value.showItem = true
  } else {
    if (sidebar.value.item && item.node.id === sidebar.value.item.node.id) return
    sidebar.value.showItem = false
    sidebar.value.item = item
    sidebar.value.itemType = itemType
    setTimeout(() => {
      sidebar.value.showItem = true
    }, 100)
  }
}

function refreshSidebar() {
  if (!sidebar.value.show || !sidebar.value.item?.node?.id) return
  if (sidebar.value.itemType === 'edge') {
    const edge = edges.value.find(e => e.id === sidebar.value.item.node.id)
    if (!edge) { sidebarClose(); return }
    sidebar.value.item = {
      node: {
        id: edge.id,
        value: edge.label || '',
        description: edge.data?.description || '',
        edge: true,
        source: { id: edge.source, style: nodes.value.find(n => n.id === edge.source)?.data?.kind || '' },
        target: { id: edge.target },
        edges: [],
      },
      config: { kind: 'edge', stepID: edge.id, parentID: edge.source, childID: edge.target, expr: edge.data?.expr || '' },
    }
    return
  }
  const node = nodes.value.find(n => n.id === sidebar.value.item.node.id)
  if (!node) { sidebarClose(); return }
  sidebar.value.showItem = false
  sidebar.value.item = buildSidebarItem(node)
  nextTick(() => { sidebar.value.showItem = true })
}

function sidebarClose() {
  sidebar.value.show = false
  setTimeout(() => {
    sidebar.value.showItem = false
    sidebar.value.item = undefined
    sidebar.value.itemType = undefined
  }, 300)
}

function sidebarDelete() {
  if (!getSelectedItem.value) return
  const id = getSelectedItem.value.node.id

  if (sidebar.value.itemType === 'edge') {
    // Delete the edge itself (#3)
    const edge = edges.value.find(e => e.id === id)
    edges.value = edges.value.filter(e => e.id !== id)
    // Renumber excl gateway edges if applicable (#1)
    if (edge) renumberGatewayEdges(edge.source)
  } else {
    // Remove node + connected edges
    const sourceIds = new Set([id])
    nodes.value = nodes.value.filter(n => n.id !== id)
    edges.value = edges.value.filter(e => {
      if (e.source === id || e.target === id) {
        // Renumber remaining excl gateway edges for any affected sources
        if (e.target === id) sourceIds.add(e.source)
        return false
      }
      return true
    })
    sourceIds.forEach(sid => renumberGatewayEdges(sid))
  }

  sidebarClose()
  saveToHistory()
  emit('change-detected')
}

function setValue(value, defaultName = false) {
  const item = sidebar.value.item
  if (!item) return

  if (sidebar.value.itemType === 'edge') {
    // Update edge label
    const edge = edges.value.find(e => e.id === item.node.id)
    if (edge) {
      edge.label = value
      item.node.value = value
    }
  } else {
    // Update node label + data
    const node = nodes.value.find(n => n.id === item.node.id)
    if (node) {
      node.data = { ...node.data, label: value, defaultName }
      item.node.value = value
      item.config.defaultName = defaultName
    }
  }

  emit('change-detected')
  debouncedSaveToHistory()
}

/* ─── Sidebar name / description inputs ─── */
function onSidebarNameChange(value) {
  // setValue handles both edges (writes edge.label) and nodes (writes node.data.label
  // + flips defaultName off so per-kind auto-naming stops stomping user input).
  if (value === sidebar.value.item?.node?.value) return
  setValue(value ?? '', false)
}

function onSidebarDescriptionChange(value) {
  const item = sidebar.value.item
  if (!item?.node) return
  const next = value ?? ''
  if (next === item.node.description) return
  item.node.description = next
  emit('change-detected')
  debouncedSaveToHistory()
}

/* ─── Issues ─── */
function openIssuesModal(nodeId) {
  issuesModal.value.issues = issues.value[nodeId] || []
  issuesModal.value.show = true
}

/* ─── Zoom ─── */
function zoom(up = true) {
  if (up) vfZoomIn()
  else vfZoomOut()
  zoomLevel.value = getViewport().zoom || zoomLevel.value
}

function resetZoom() {
  fitView({ padding: 0.2, maxZoom: 1 })
  zoomLevel.value = getViewport().zoom || 1
}

/* ─── Keyboard ─── */
function isEditableTarget(event) {
  const t = event.target
  if (!t) return false
  const tag = t.tagName
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return true
  if (t.isContentEditable) return true
  if (
    typeof t.closest === 'function' &&
    t.closest('[contenteditable="true"], .ProseMirror, .ql-editor, .cm-editor')
  )
    return true
  return false
}

function keybinds(event) {
  const editable = isEditableTarget(event)

  // Ctrl+S
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    if (!document.getElementById('expression-editor')) {
      saveWorkflow()
    }
  }

  // Ctrl+Z
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'z' && !event.shiftKey) {
    if (editable) return
    event.preventDefault()
    undo()
    refreshSidebar()
    nextTick(() => checkExistingTriggerPaths())
  }

  // Ctrl+Shift+Z
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'z' && event.shiftKey) {
    if (editable) return
    event.preventDefault()
    redo()
    refreshSidebar()
    nextTick(() => checkExistingTriggerPaths())
  }

  // Ctrl+C
  if ((event.ctrlKey || event.metaKey) && event.key === 'c') {
    if (editable) return
    copySelected()
  }

  // Ctrl+X
  if ((event.ctrlKey || event.metaKey) && event.key === 'x') {
    if (editable) return
    cutSelected()
  }

  // Ctrl+V
  if ((event.ctrlKey || event.metaKey) && event.key === 'v') {
    if (editable) return
    pasteClipboard()
  }

  // Ctrl+A
  if ((event.ctrlKey || event.metaKey) && event.key === 'a') {
    if (editable) return
    event.preventDefault()
    nodes.value = nodes.value.map(n => ({ ...n, selected: true }))
  }

  // Ctrl+Space
  if ((event.ctrlKey || event.metaKey) && event.key === ' ') {
    event.preventDefault()
    resetZoom()
  }

  // Delete / Backspace
  if (event.key === 'Delete' || event.key === 'Backspace') {
    if (
      event.target.tagName !== 'INPUT' &&
      event.target.tagName !== 'TEXTAREA' &&
      !event.target.closest('[contenteditable]')
    ) {
      deleteSelected()
    }
  }

  // Arrow nudge
  if (['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(event.key)) {
    if (event.target.tagName !== 'INPUT' && event.target.tagName !== 'TEXTAREA') {
      const delta = event.shiftKey ? 8 : 40
      const selected = nodes.value.filter(n => n.selected)
      if (selected.length) {
        event.preventDefault()
        selected.forEach(n => {
          if (event.key === 'ArrowLeft') n.position.x -= delta
          if (event.key === 'ArrowRight') n.position.x += delta
          if (event.key === 'ArrowUp') n.position.y -= delta
          if (event.key === 'ArrowDown') n.position.y += delta
        })
        saveToHistory()
        emit('change-detected')
      }
    }
  }

  // Shift + ?
  if (event.shiftKey && event.key === '?') {
    helpVisible.value = true
  }
}

function deleteSelected() {
  const selectedNodes = nodes.value.filter(n => n.selected)
  const selectedEdges = edges.value.filter(e => e.selected)

  if (selectedNodes.length === 0 && selectedEdges.length === 0) return

  if (sidebar.value.item && selectedNodes.some(n => n.id === sidebar.value.item.node.id)) {
    sidebarClose()
  }

  // Track gateway sources that need edge renumbering
  const gatewaySourceIds = new Set()

  // Delete selected edges (#10)
  if (selectedEdges.length) {
    const edgeIds = new Set(selectedEdges.map(e => e.id))
    selectedEdges.forEach(e => gatewaySourceIds.add(e.source))
    edges.value = edges.value.filter(e => !edgeIds.has(e.id))
  }

  // Delete selected nodes + their connected edges
  if (selectedNodes.length) {
    const nodeIds = new Set(selectedNodes.map(n => n.id))
    edges.value = edges.value.filter(e => {
      if (nodeIds.has(e.source) || nodeIds.has(e.target)) {
        if (nodeIds.has(e.target)) gatewaySourceIds.add(e.source)
        return false
      }
      return true
    })
    nodes.value = nodes.value.filter(n => !nodeIds.has(n.id))
  }

  // Renumber excl gateway edges (#1)
  gatewaySourceIds.forEach(sid => renumberGatewayEdges(sid))

  clearHighlights()
  saveToHistory()
  emit('change-detected')
}

/**
 * Keep parallel gateway `ref` (fork vs join) in sync with actual edge counts.
 * Matches Human's mxGraph behaviour: one style (`gatewayParallel`), role
 * inferred from in/out edge balance — join when in > out, otherwise fork.
 */
function reconcileParallelGateways() {
  let mutated = false
  nodes.value.forEach(n => {
    if (n.data?.kind !== 'gateway') return
    if (n.data?.ref !== 'fork' && n.data?.ref !== 'join') return
    let inCount = 0
    let outCount = 0
    edges.value.forEach(e => {
      if (e.source === n.id) outCount++
      if (e.target === n.id) inCount++
    })
    const expected = inCount > outCount ? 'join' : 'fork'
    if (n.data.ref !== expected) {
      n.data = { ...n.data, ref: expected }
      mutated = true
    }
  })
  return mutated
}

watch(
  edges,
  () => {
    reconcileParallelGateways()
  },
  { deep: false },
)

/**
 * Relabel remaining out-edges of certain source kinds after deletion so their
 * labels match the semantic slot (position-dependent). Covers:
 *   - exclusive gateway:  #1 - If, #N - Else (if)
 *   - iterator:           Body, End
 *   - error-handler:      Try, Catch
 */
function renumberGatewayEdges(sourceId) {
  const sourceNode = nodes.value.find(n => n.id === sourceId)
  if (!sourceNode) return
  const kind = sourceNode.data?.kind
  const ref = sourceNode.data?.ref

  const outEdges = edges.value.filter(e => e.source === sourceId)

  if (kind === 'gateway' && ref === 'excl') {
    outEdges.forEach((e, idx) => {
      e.label = idx === 0 ? '#1 - If' : `#${idx + 1} - Else (if)`
    })
  } else if (kind === 'iterator') {
    outEdges.forEach((e, idx) => {
      e.label = idx === 0 ? 'Body' : 'End'
    })
  } else if (kind === 'error-handler') {
    outEdges.forEach((e, idx) => {
      e.label = idx === 0 ? 'Try' : 'Catch'
    })
  }
}

function isValidConnection(connection) {
  return validateConnection(connection, {
    nodes: nodes.value,
    edges: edges.value,
    edgeUpdatingId: edgeUpdatingId.value,
  })
}

function onEdgeUpdateStart({ edge }) {
  edgeUpdatingId.value = edge?.id || null
}

function onEdgeUpdateEnd() {
  edgeUpdatingId.value = null
}

/**
 * Reconnect an existing edge to a new source/target/handle. Rewires in place,
 * preserves edge id, label and style; the codec re-derives mxStyle on save.
 */
function onEdgeUpdate({ edge, connection }) {
  edges.value = edges.value.map(e => {
    if (e.id !== edge.id) return e
    return {
      ...e,
      source: connection.source,
      target: connection.target,
      sourceHandle: connection.sourceHandle ?? null,
      targetHandle: connection.targetHandle ?? null,
    }
  })
  saveToHistory()
}

/* ─── Toolbar tooltip ─── */
function showToolbarTooltip(event, item) {
  const rect = event.target.getBoundingClientRect()
  activeToolbarTooltip.value = { title: item.label, icon: item.iconSrc, tooltip: item.tooltip }
  toolbarTooltipStyle.value = {
    top: `${rect.top}px`,
    left: `${rect.right + 10}px`,
  }
}

/* ─── Drawer resize ─── */
function startDrawerResize() {
  isResizingDrawer.value = true
  document.addEventListener('mousemove', resizeDrawer)
  document.addEventListener('mouseup', stopDrawerResize)
}

function resizeDrawer(e) {
  if (!isResizingDrawer.value) return
  const newWidth = window.innerWidth - e.clientX
  drawerWidth.value = Math.max(280, Math.min(800, newWidth))
}

function stopDrawerResize() {
  isResizingDrawer.value = false
  document.removeEventListener('mousemove', resizeDrawer)
  document.removeEventListener('mouseup', stopDrawerResize)
}

/* ─── Trigger path check ─── */
function checkExistingTriggerPaths() {
  triggersPathsChanged.value = [...triggers.value].some(({ stepID = '0', meta = {} }) => {
    if (stepID !== NoID) {
      const triggerNodeId = meta?.visual?.id
      const outEdge = edges.value.find(e => e.source === String(triggerNodeId))
      if (outEdge) {
        return outEdge.target !== String(stepID)
      }
      return true
    }
    return false
  })
}

/* ─── Dry Run / Test ─── */
function startTest(cellID) {
  dryRun.value.cellID = cellID
  loadTestScope()
}

async function loadTestScope() {
  if (props.changeDetected) {
    toast.add({
      severity: 'warn',
      summary: t('notification.failed-test'),
      detail: t('notification.save-workflow'),
      life: 5000,
    })
    return
  }
  if (hasIssues.value) {
    toast.add({
      severity: 'warn',
      summary: t('notification.failed-test'),
      detail: t('notification.resolve-issues'),
      life: 5000,
    })
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

  const triggerNode = nodes.value.find(n => n.id === String(dryRun.value.cellID))
  const trg = triggerNode?.data?.triggers
  if (!trg) return

  const { resourceType, eventType } = trg
  const et = (
    eventTypes.value.find(et => resourceType === et.resourceType && eventType === et.eventType) ||
    {}
  ).properties

  if (et) {
    let lookup = false
    if (et.length) {
      dryRun.value.initialScope = et.reduce((scope, p) => {
        let label = `${p.name}${lookupableTypes.includes(p.name) ? t('editor.id-parenthesis') : ''}`
        if (p.type === 'ComposeNamespace' || p.type === 'ComposeModule') {
          label = `${p.name} ${t('editor.handle')}`
        }

        let description = ''
        if (p.type === 'ComposeRecord') description = t('editor.required-namespace-and-module')
        else if (p.type === 'ComposeModule' || p.name === 'page' || p.name === 'oldPage')
          description = t('editor.required-namespace')

        scope[p.name] = {
          label,
          value: (dryRun.value.initialScope[p.name] || {}).value,
          lookup: lookupableTypes.includes(p.name),
          description,
        }
        lookup = lookup || lookupableTypes.includes(p.name)
        return scope
      }, {})

      encodeInput(dryRun.value.initialScope, $ComposeAPI, $SystemAPI)
        .then(input => {
          dryRun.value.input = input
          dryRun.value.lookup = lookup
          dryRun.value.show = true
        })
        .catch(e =>
          toast.add({
            severity: 'error',
            summary: t('notification.initial-scope-load-failed'),
            detail: e?.message,
            life: 5000,
          }),
        )
    } else {
      dryRun.value.initialScope = {}
      testWorkflow()
    }
  } else {
    toast.add({
      severity: 'warn',
      summary: t('notification.failed-test'),
      detail: t('notification.event-type-not-found'),
      life: 5000,
    })
  }
}

async function dryRunOk(e) {
  if (dryRun.value.lookup) {
    e.preventDefault()
    encodeInput(dryRun.value.initialScope, $ComposeAPI, $SystemAPI)
      .then(input => {
        dryRun.value.input = input
        dryRun.value.inputEdited = input
        dryRun.value.lookup = false
      })
      .catch(e =>
        toast.add({
          severity: 'error',
          summary: t('notification.initial-scope-load-failed'),
          detail: e?.message,
          life: 5000,
        }),
      )
  } else {
    dryRun.value.show = false
    testWorkflow(dryRun.value.inputEdited)
  }
}

function onDryRunEdit(e) {
  dryRun.value.inputEdited = e
}

async function testWorkflow(input = {}) {
  clearHighlights()
  dryRun.value.processing = true

  const triggerNode = nodes.value.find(n => n.id === String(dryRun.value.cellID))
  const trg = triggerNode?.data?.triggers

  const testParams = {
    workflowID: workflow.value.workflowID,
    stepID: trg?.stepID || '0',
    trace: workflow.value.canManageWorkflowSessions || false,
    wait: false,
    async: true,
    input,
  }

  toast.add({
    severity: 'info',
    summary: t('notification.test-in-progress'),
    detail: t('notification.started-test'),
    life: 3000,
  })

  $AutomationAPI
    .workflowExec(testParams)
    .then(({ sessionID }) => {
      dryRun.value.sessionID = sessionID

      const pollSession = () => {
        return new Promise((resolve, reject) => {
          const checkSession = () => {
            $AutomationAPI
              .sessionRead({ sessionID })
              .then(session => {
                const { completedAt, status, stacktrace, error = false } = session
                setTimeout(() => {
                  if (completedAt) {
                    if (stacktrace) {
                      renderTrace(testParams.stepID, stacktrace)
                      if (status === 'completed') {
                        toast.add({
                          severity: 'success',
                          summary: t('notification.test-completed'),
                          detail: t('notification.workflow-test-completed'),
                          life: 3000,
                        })
                      }
                    } else {
                      toast.add({
                        severity: 'warn',
                        summary: t('notification.test-completed'),
                        detail: t('notification.trace-unavailable'),
                        life: 5000,
                      })
                    }
                    if (error) reject(new Error(error))
                    else resolve()
                  } else {
                    checkSession()
                  }
                }, 1000)
              })
              .catch(reject)
          }
          checkSession()
        })
      }

      return pollSession()
    })
    .catch(e =>
      toast.add({
        severity: 'error',
        summary: t('notification.failed-test'),
        detail: e?.message,
        life: 5000,
      }),
    )
    .finally(() => {
      dryRun.value.lookup = true
      dryRun.value.processing = false
      dryRun.value.sessionID = undefined
    })
}

function cancelWorkflow() {
  const { sessionID, processing } = dryRun.value
  if (processing && sessionID) {
    dryRun.value.sessionID = undefined
    dryRun.value.processing = false
    $AutomationAPI
      .sessionCancel({ sessionID })
      .then(() =>
        toast.add({
          severity: 'info',
          summary: 'Stopping test',
          detail: 'Workflow test canceled',
          life: 3000,
        }),
      )
      .catch(e =>
        toast.add({
          severity: 'error',
          summary: 'Test cancel failed',
          detail: e?.message,
          life: 5000,
        }),
      )
  }
}

/* ─── Trace rendering ─── */
function renderTrace(firstStepID, trace = []) {
  clearHighlights()

  const cells = {}
  trace
    .filter(t => t)
    .forEach(({ stepID, parentID, stepTime, error = false }, index) => {
      const cell = { index, stepID, parentID, stepTime, error }
      if (cells[stepID]) cells[stepID].push(cell)
      else cells[stepID] = [cell]
    })

  // Highlight first edge (trigger → first step)
  const firstEdge = edges.value.find(
    e => e.source === String(dryRun.value.cellID) && e.target === String(firstStepID),
  )
  if (firstEdge) {
    firstEdge.data = { ...firstEdge.data, traceState: 'success' }
  }

  Object.entries(cells).forEach(([stepID, frames]) => {
    if (stepID !== '0') {
      const error = frames.some(f => f.error)
      const logParts = frames.map(
        f => `#${f.index + 1} - ${f.stepTime}ms${f.error ? ` (Error: ${f.error})` : ''}`,
      )
      const times = frames.map(f => Number(f.stepTime) || 0)
      const sum = times.reduce((a, b) => a + b, 0)
      const min = times.length ? Math.min(...times) : 0
      const max = times.length ? Math.max(...times) : 0
      const avg = times.length ? Math.round(sum / times.length) : 0
      const header = `N=${times.length}  MIN=${min}ms  MAX=${max}ms  AVG=${avg}ms  SUM=${sum}ms`
      const log = `${header}\n${logParts.join('\n')}`

      // Set trace state on node
      const node = nodes.value.find(n => n.data?.stepID === stepID || n.id === stepID)
      if (node) {
        node.data = { ...node.data, traceState: error ? 'error' : 'success', traceLog: log }
      }

      // Highlight connected edges
      frames.forEach(({ parentID }) => {
        const edge = edges.value.find(
          e => e.source === String(parentID) && e.target === String(stepID),
        )
        if (edge) {
          edge.data = { ...edge.data, traceState: error ? 'error' : 'success' }
        }
      })
    }
  })
}

/* ─── Save / Import ─── */
function getJsonModel() {
  return encodeWorkflow(nodes.value, edges.value)
}

function handleWorkflowSave(wf) {
  workflow.value = wf
  saveWorkflow()
}

function saveWorkflow() {
  emit('save', { ...workflow.value, ...getJsonModel() })
}

function importJSON(workflows = []) {
  try {
    importProcessing.value = true
    const [wf] = workflows
    triggers.value = wf.triggers || []
    workflow.value = {
      ...workflow.value,
      steps: wf.steps || [],
      paths: wf.paths || [],
    }
    render(workflow.value)
    importProcessing.value = false
    emit('change-detected')
    configuratorVisible.value = false
    toast.add({ severity: 'success', summary: t('notification.imported-workflow'), life: 3000 })
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('notification.import-failed'),
      detail: e?.message,
      life: 5000,
    })
  }
}

/* ─── API Calls ─── */
async function getFunctionTypes() {
  return $AutomationAPI
    .functionList()
    .then(({ set = [] }) => {
      functionTypes.value = [
        ...set,
        ...(components.promptDefinitions || []),
        ...[
          {
            ref: 'error-handler',
            kind: 'error-handler',
            meta: { short: 'Handle error' },
            parameters: [],
            results: [
              { name: 'error', types: ['Any'] },
              { name: 'errorMessage', types: ['String'] },
              { name: 'errorStepID', types: ['Integer'] },
            ],
          },
          {
            ref: 'exec-workflow',
            kind: 'error-handler',
            meta: { short: 'Execute a workflow' },
            parameters: [
              { name: 'workflow', types: ['ID', 'Handle'], required: true },
              { name: 'scope', types: ['Vars'], required: false },
            ],
            results: [],
          },
        ],
      ]
    })
    .catch(e =>
      toast.add({
        severity: 'error',
        summary: t('notification.failed-fetch-functions'),
        detail: e?.message,
        life: 5000,
      }),
    )
}

async function getEventTypes() {
  return $AutomationAPI
    .eventTypesList()
    .then(({ set = [] }) => {
      eventTypes.value = set
    })
    .catch(e =>
      toast.add({
        severity: 'error',
        summary: t('notification.event-type-fetch-failed'),
        detail: e?.message,
        life: 5000,
      }),
    )
}

/* ─── Expose for parent ─── */
defineExpose({
  getJsonModel,
})
</script>

<style scoped>
#editor {
  color: var(--p-text-color, #1e293b);
}

.toolbar {
  width: 55px;
}

.toolbar-item {
  cursor: grab;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s;
}

.toolbar-item:hover {
  background: var(--p-highlight-background);
}

.toolbar-item:active {
  cursor: grabbing;
}

.component-fade-enter-active,
.component-fade-leave-active {
  transition: opacity 0.3s ease;
}

.component-fade-enter,
.component-fade-leave-to {
  opacity: 0;
}

.vueflow-canvas {
  width: 100%;
  height: 100%;
}

.vueflow-canvas :deep(.vue-flow__background) {
  background-color: transparent;
}

.vueflow-canvas :deep(.vue-flow__background pattern circle) {
  fill: color-mix(in srgb, var(--p-text-muted-color) 30%, transparent);
}

.vueflow-canvas :deep(.vue-flow__edge-path) {
  stroke: var(--p-text-muted-color);
  stroke-width: 2;
}

.vueflow-canvas :deep(.vue-flow__edge.selected .vue-flow__edge-path),
.vueflow-canvas :deep(.vue-flow__edge:hover .vue-flow__edge-path) {
  stroke: var(--p-primary-color);
}

.vueflow-canvas :deep(.vue-flow__minimap) {
  display: none;
}

/* https://stackoverflow.com/a/40991531/17926309 */
.saving::after {
  display: inline-block;
  animation: saving steps(1, end) 1s infinite;
  content: '';
}

@keyframes saving {
  0% {
    content: '';
  }
  25% {
    content: '.';
  }
  50% {
    content: '..';
  }
  75% {
    content: '...';
  }
  100% {
    content: '';
  }
}
</style>

<style>
/* VueFlow global overrides (unscoped) */
.vue-flow__connection-line {
  stroke: var(--p-text-muted-color);
  stroke-width: 2;
}

.vue-flow__edge-path {
  stroke: var(--p-text-muted-color);
}

/* Selection rectangle */
.vue-flow__selection-pane {
  cursor: grab;
}
</style>
