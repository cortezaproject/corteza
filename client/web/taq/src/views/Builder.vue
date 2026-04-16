<template>
  <Teleport to="#topbar-title" defer>
    <div
      class="flex items-center gap-2 group/title hover:bg-surface-hover px-2 py-1 rounded-border cursor-pointer -ml-2 transition-colors"
      @click="showConfigDialog = true"
    >
      {{ editor.name.value }}
      <Button icon="pi pi-pencil" text rounded size="small" class="!w-5 !h-5 !p-0 shrink-0" />
    </div>
  </Teleport>

  <!-- Loading state -->
  <div v-if="editor.loading.value" class="h-full flex items-center justify-center">
    <ProgressSpinner />
  </div>

  <div v-else class="builder-layout h-full flex flex-col relative overflow-hidden">
    <!-- Top-left overlay actions -->
    <div class="absolute top-1 left-3 z-20 flex flex-col gap-2 max-w-screen-lg">
      <!-- Optional Description -->
      <div
        v-if="editor.automation.value.meta?.description"
        class="text-sm text-muted-color whitespace-pre-wrap px-1"
      >
        {{ editor.automation.value.meta.description }}
      </div>

      <!-- Run As -->
      <div v-if="runAsUser" class="flex items-center gap-2 px-1">
        <span class="text-sm text-muted-color">{{ $t('builder.canvas.runAs') }}</span>
        <Tag
          :value="runAsUser"
          severity="secondary"
          icon="pi pi-user"
          class="border border-surface"
        />
      </div>

      <div class="flex items-center gap-2 mt-2">
        <!-- Main toggle/run card -->
        <div
          class="flex items-center gap-3 bg-surface rounded-lg border border-surface px-3 py-2 shadow-sm"
        >
          <div class="flex items-center gap-2">
            <ToggleSwitch
              :model-value="editor.enabled.value"
              @update:model-value="toggleEnabled"
              input-id="taq-enabled"
            />
            <label for="taq-enabled" class="text-sm">{{ $t('builder.enabled') }}</label>
          </div>
          <template v-if="editor.automationId.value && !editor.isEmpty.value">
            <Divider layout="vertical" class="!m-0" />
            <Button
              :label="$t('builder.run')"
              icon="pi pi-play"
              severity="success"
              outlined
              size="small"
              :loading="editor.running.value"
              :disabled="!editor.enabled.value || editor.running.value"
              @click="onRunClick"
            />
          </template>
          <Divider layout="vertical" class="!m-0" />
          <Button
            v-tooltip.bottom="
              showAllPreviews
                ? $t('builder.preview.hideConfigurations')
                : $t('builder.preview.showConfigurations')
            "
            :icon="showAllPreviews ? 'pi pi-eye' : 'pi pi-eye-slash'"
            :severity="showAllPreviews ? 'primary' : 'secondary'"
            :outlined="!showAllPreviews"
            size="small"
            @click="togglePreviews"
          />
        </div>

        <!-- Standalone permissions button -->
        <CPermissionsButton
          v-if="editor.automationId.value"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::automation:ng-automation/${editor.automationId.value}`"
          :title="editor.name.value"
          :target="editor.name.value"
          severity="secondary"
          class="bg-surface shadow-sm !border-surface"
        />
      </div>

      <!-- Execution result card (below the toolbar card) -->
      <Transition
        enter-active-class="transition-all duration-200 ease-out"
        enter-from-class="-translate-y-2 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition-all duration-150 ease-in"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="-translate-y-2 opacity-0"
      >
        <div
          v-if="isTraceActive"
          class="bg-surface rounded-lg border border-surface px-3 py-2 shadow-sm flex items-center gap-2 w-fit"
        >
          <Tag
            :severity="editor.traceStatus.value === 'failed' ? 'danger' : 'success'"
            :icon="
              editor.traceStatus.value === 'failed' ? 'pi pi-times-circle' : 'pi pi-check-circle'
            "
            :value="
              editor.traceStatus.value === 'failed'
                ? $t('builder.trace.failed')
                : $t('builder.trace.completed')
            "
            class="text-sm"
          />
          <span
            v-if="editor.traceExecution.value?.duration"
            class="text-xs text-muted-color flex items-center gap-1"
          >
            <i class="pi pi-clock text-xs" />
            {{ editor.traceExecution.value.duration }}
          </span>
          <Divider layout="vertical" class="!m-0 !mx-1" />
          <Button
            :label="$t('builder.trace.clear')"
            icon="pi pi-times"
            severity="secondary"
            text
            size="small"
            @click="editor.clearTrace"
          />
        </div>
      </Transition>

      <!-- Trace Step Panel (shows I/O for the selected node, underneath execution banner) -->
      <Transition
        enter-active-class="transition-all duration-200 ease-out"
        enter-from-class="-translate-y-2 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition-all duration-150 ease-in"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="-translate-y-2 opacity-0"
      >
        <div
          v-if="isTraceActive && selectedNode && selectedTraceFrame && showTracePanel"
          class="bg-surface border border-surface rounded-lg shadow-sm flex flex-col overflow-hidden pointer-events-auto"
          :style="{ width: '300px', maxHeight: 'calc(100vh - 10rem)' }"
        >
          <TracePanel
            :frame="selectedTraceFrame"
            :execution-error="editor.traceExecution.value?.error"
            :step-name="selectedNode?.data?.label || $t('builder.configSidebar.node')"
            @close="showTracePanel = false"
          />
        </div>
      </Transition>
    </div>

    <!-- VueFlow Canvas -->
    <div class="flex-1 min-h-0" @auxclick="onMiddleMouseClick">
      <VueFlow
        v-model:nodes="editor.nodes.value"
        v-model:edges="editor.edges.value"
        :default-viewport="{ zoom: 1.5 }"
        :min-zoom="0.5"
        :max-zoom="2"
        :nodes-draggable="false"
        :nodes-connectable="false"
        :edges-updatable="false"
        :pan-on-drag="[0, 1, 2]"
        :pan-on-scroll="true"
        :zoom-on-scroll="false"
        :zoom-on-double-click="false"
        :selection-key-code="null"
        @node-click="onNodeClick"
        @pane-click="onPaneClick"
      >
        <!-- Background with dots pattern -->
        <Background :gap="27" :size="1" />

        <!-- Custom node types -->
        <template #node-trigger="props">
          <TriggerNode
            v-bind="props"
            :triggers="store.triggers"
            :nodes="editor.nodes.value"
            :always-show-preview="showAllPreviews"
            :trace-active="isTraceActive"
            @delete="confirmDeleteNode"
            @replace="startReplace"
          />
        </template>
        <template #node-step="props">
          <StepNode
            v-bind="props"
            :functions="store.functions"
            :nodes="editor.nodes.value"
            :always-show-preview="showAllPreviews"
            :trace-frame="getTraceFrame(props)"
            :trace-active="isTraceActive"
            @delete="confirmDeleteNode"
            @replace="startReplace"
          />
        </template>
        <template #node-branch="props">
          <BranchNode
            v-bind="props"
            :functions="store.functions"
            :nodes="editor.nodes.value"
            :edges="editor.edges.value"
            :always-show-preview="showAllPreviews"
            :trace-frame="getTraceFrame(props)"
            :trace-active="isTraceActive"
            @delete="confirmDeleteNode"
            @replace="startReplace"
          />
        </template>
        <template #node-iterator="props">
          <IteratorNode
            v-bind="props"
            :functions="store.functions"
            :nodes="editor.nodes.value"
            :edges="editor.edges.value"
            :always-show-preview="showAllPreviews"
            :trace-frame="getTraceFrame(props)"
            :trace-active="isTraceActive"
            @delete="confirmDeleteNode"
            @replace="startReplace"
          />
        </template>
        <template #node-end="props">
          <EndNode v-bind="props" />
        </template>
        <template #node-loop="props">
          <LoopNode v-bind="props" />
        </template>

        <!-- Custom edge with + button -->
        <template #edge-addable="props">
          <AddableEdge
            v-bind="props"
            :trace-active="isTraceActive"
            :trace-traversed="isEdgeTraversed(props)"
            @add="onEdgeAdd"
          />
        </template>
      </VueFlow>
    </div>

    <!-- Bottom Toolbar -->
    <div class="shrink-0 z-10 body-bg">
      <CToolbar>
        <template #start>
          <Button
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="goBack"
          />
        </template>
        <template #center>
          <!-- Zoom -->
          <div v-if="!editor.isEmpty.value" class="flex items-center gap-1">
            <Button icon="pi pi-minus" severity="secondary" @click="zoomOut" />
            <Button icon="pi pi-arrows-alt" severity="secondary" @click="fitToScreen" />
            <Button icon="pi pi-plus" severity="secondary" @click="zoomIn" />
          </div>
        </template>
        <template #end>
          <div class="flex items-center gap-3">
            <!-- Undo/Redo -->
            <div class="flex gap-1">
              <Button
                icon="pi pi-undo"
                severity="secondary"
                :disabled="!editor.canUndo.value"
                @click="editor.undo"
              />
              <Button
                icon="pi pi-refresh"
                severity="secondary"
                :disabled="!editor.canRedo.value"
                @click="editor.redo"
              />
            </div>
            <Divider layout="vertical" class="!m-0" />
            <Button
              :label="$t('builder.save')"
              icon="pi pi-save"
              :loading="editor.saving.value"
              :disabled="editor.saving.value"
              @click="editor.save"
            />
          </div>
        </template>
      </CToolbar>
    </div>

    <!-- Node Picker Dialog -->
    <Dialog
      v-model:visible="showNodePicker"
      modal
      :header="
        replaceNodeId
          ? nodePickerCategory === 'trigger'
            ? $t('builder.replaceTrigger')
            : $t('builder.replaceStep')
          : nodePickerCategory === 'trigger'
            ? $t('builder.selectTrigger')
            : $t('builder.addStep')
      "
      :pt="{ content: { class: 'p-0' } }"
      :style="{ width: '50rem' }"
    >
      <NodePicker
        :filter-category="nodePickerCategory"
        @select="handleNodeSelect"
        @close="closeNodePicker"
      />
    </Dialog>

    <!-- Reference Panel (opens on input click, closes via button or sidebar close) -->
    <Transition
      enter-active-class="transition-transform duration-200 ease-out"
      enter-from-class="translate-x-[300px]"
      enter-to-class="translate-x-0"
    >
      <div
        v-if="selectedNode && showReferencePanel"
        class="right-sidebar"
        :style="{ right: `calc(${drawerWidth}px + 1rem)`, width: '300px' }"
      >
        <ReferencePanel
          :upstream-results="upstreamResults"
          :active-argument="activeReferenceArgument"
          :current-reference="currentReference"
          @select="handleReferenceSelect"
          @close="closeReferencePanel"
        />
      </div>
    </Transition>

    <!-- Config Sidebar (resizable drawer) -->
    <Transition
      enter-active-class="transition-transform duration-300 ease-in-out"
      enter-from-class="translate-x-full"
      enter-to-class="translate-x-0"
      leave-active-class="transition-transform duration-300 ease-in-out"
      leave-from-class="translate-x-0"
      leave-to-class="translate-x-full"
    >
      <div
        v-if="selectedNode"
        class="config-drawer right-sidebar flex"
        :style="{ width: `${drawerWidth}px` }"
      >
        <!-- Resize handle -->
        <div
          class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors"
          @mousedown="startDrawerResize"
        />
        <!-- Drawer content -->
        <ConfigSidebar
          ref="configSidebarRef"
          :node="selectedNode"
          :edges="editor.edges.value"
          :nodes="editor.nodes.value"
          :functions="store.functions"
          :triggers="store.triggers"
          :upstream-results="upstreamResults"
          @close="clearSelection"
          @delete="handleDeleteSelected"
          @add-branch="handleAddBranch"
          @reorder-branches="handleReorderBranches"
          @update-arguments="handleUpdateArguments"
          @update-constraints="handleUpdateConstraints"
          @update-input-schema="handleUpdateInputSchema"
          @toggle-reference="handleToggleReference"
          @update-gateway-type="handleUpdateGatewayType"
          @update-branch-expr="handleUpdateBranchExpr"
          @update-metadata="handleUpdateMetadata"
        />
      </div>
    </Transition>

    <!-- Empty state: Add first trigger -->
    <div v-if="editor.isEmpty.value" class="absolute inset-0 flex items-center justify-center">
      <div class="text-center">
        <div class="text-muted-color text-lg mb-4">{{ $t('builder.emptyState') }}</div>
        <Button
          :label="$t('builder.addTrigger')"
          icon="pi pi-plus"
          @click="openNodePicker({ first: true })"
        />
      </div>
    </div>

    <!-- Unsaved Changes Dialog -->
    <Dialog
      v-model:visible="showUnsavedDialog"
      modal
      :header="$t('builder.unsavedChanges.header')"
      :style="{ width: '40rem' }"
    >
      <div class="mb-4 text-color whitespace-pre-line">
        {{ $t('builder.unsavedChanges.description') }}
      </div>

      <template #footer>
        <div class="flex justify-end w-full h-full items-center gap-2">
          <Button
            :label="$t('general.label.cancel')"
            text
            size="small"
            severity="secondary"
            @click="showUnsavedDialog = false"
          />
          <Button
            :label="$t('builder.saveAndRun')"
            icon="pi pi-save"
            severity="success"
            size="small"
            @click="proceedRun(true)"
          />
        </div>
      </template>
    </Dialog>

    <!-- Run Modal -->
    <RunModal v-model:visible="showRunModal" :properties="runProperties" @run="onRunConfirm" />

    <!-- General Config Modal -->
    <TaqConfigModal
      v-model:visible="showConfigDialog"
      mode="edit"
      :initial-name="editor.name.value"
      :initial-description="editor.automation.value.meta?.description"
      :initial-labels="editor.automation.value.labels"
      :initial-run-as="editor.automation.value.runAs"
      @saved="handleConfigSave"
    />
  </div>
</template>

<script setup>
import { Background } from '@vue-flow/background'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'

import { useFlowEditor } from '@/composables/useFlowEditor'
import { useAutomationStore } from '@/stores/automation'
import { components } from '@cortezaproject/corteza-vue-next'

const { CToolbar } = components

import { useConfirmDelete, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { computed, inject, nextTick, onMounted, onUnmounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import ConfigSidebar from '@/components/builder/ConfigSidebar.vue'
import NodePicker from '@/components/builder/NodePicker.vue'
import ReferencePanel from '@/components/builder/ReferencePanel.vue'
import TracePanel from '@/components/builder/TracePanel.vue'
import AddableEdge from '@/components/flow/AddableEdge.vue'
import BranchNode from '@/components/flow/BranchNode.vue'
import EndNode from '@/components/flow/EndNode.vue'
import LoopNode from '@/components/flow/LoopNode.vue'
import IteratorNode from '@/components/flow/IteratorNode.vue'
import StepNode from '@/components/flow/StepNode.vue'
import TriggerNode from '@/components/flow/TriggerNode.vue'
import RunModal from '@/components/builder/RunModal.vue'
import TaqConfigModal from '@/components/common/TaqConfigModal.vue'

const route = useRoute()
const router = useRouter()
const store = useAutomationStore()
const editor = useFlowEditor()
const { confirmDelete } = useConfirmDelete()
const { t } = useI18n()

useUnsavedGuard({
  isDirty: editor.isDirty,
  messageKey: 'general.editor.unsavedChanges',
})

const $SystemAPI = inject('$SystemAPI')
const runAsUser = ref(null)

watch(
  () => editor.automation.value?.runAs,
  async runAsID => {
    if (!runAsID || runAsID === '0' || runAsID === 0) {
      runAsUser.value = null
      return
    }
    try {
      const user = await $SystemAPI.userRead({ userID: runAsID })
      runAsUser.value = user?.name || user?.email || user?.handle || null
    } catch {
      runAsUser.value = null
    }
  },
  { immediate: true },
)

// VueFlow instance for viewport control and selection
const {
  fitView,
  getViewport,
  setCenter,
  zoomIn: vfZoomIn,
  zoomOut: vfZoomOut,
  getSelectedNodes,
  removeSelectedNodes,
  addSelectedNodes,
  onNodesInitialized,
} = useVueFlow()

// Fit view only once when nodes are initially loaded
let hasInitiallyFit = false
onNodesInitialized(() => {
  if (!hasInitiallyFit) {
    fitView({ padding: 0.2, maxZoom: 1.2 })
    hasInitiallyFit = true
  }
})

// Drawer state
const drawerWidth = ref(360)
const isResizingDrawer = ref(false)

// Node picker state
const showAllPreviews = ref(false)
const showNodePicker = ref(false)
const nodePickerCategory = ref(null)
const insertionPoint = ref(null)
const replaceNodeId = ref(null)

// Reference panel state
const showReferencePanel = ref(false)
const activeReferenceArgument = ref(null)
const configSidebarRef = ref(null)
const showTracePanel = ref(true)

// Run Modal State
const showRunModal = ref(false)
const runProperties = ref([])
const showUnsavedDialog = ref(false)

// Config Modal State
const showConfigDialog = ref(false)

function handleConfigSave({ name, description, labels, runAs }) {
  editor.name.value = name
  editor.automation.value.meta = {
    ...editor.automation.value.meta,
    description: description,
  }
  editor.automation.value.labels = labels || {}
  editor.automation.value.runAs = runAs || '0'
  editor.save()
}

// Provide active reference argument to descendant components (DynamicInput, CInputFieldValueMap)
provide('activeReferenceArgument', activeReferenceArgument)

// Trace helpers
const isTraceActive = computed(() => editor.traceStatus.value !== 'idle')

// Look up trace frame for a step node by matching its data.ref (handle) to the traceByHandle map
// Falls back to matching by stepID if handle-based lookup fails
function getTraceFrame(nodeProps) {
  if (!isTraceActive.value) return null
  const handle = nodeProps.data?.ref
  if (handle) {
    const frame = editor.traceByHandle.value.get(handle)
    if (frame) return frame
  }
  // Fallback: match by stepID
  const stepID = nodeProps.data?.stepID
  if (stepID) {
    return editor.traceFrames.value.find(f => f.stepID === stepID) || null
  }
  return null
}

// Check if an edge was traversed: both source and target nodes must have been executed
// Helper to check if a node has a matching trace frame
function nodeHasTrace(node) {
  if (!node) return false
  if (node.type === 'trigger') return true // triggers always count as traced
  const handle = node.data?.ref
  if (handle && editor.traceByHandle.value.get(handle)) return true
  const stepID = node.data?.stepID
  if (stepID && editor.traceFrames.value.some(f => f.stepID === stepID)) return true
  return false
}

function isEdgeTraversed(edgeProps) {
  if (!isTraceActive.value) return false
  const sourceNode = editor.nodes.value.find(n => n.id === edgeProps.source)
  const targetNode = editor.nodes.value.find(n => n.id === edgeProps.target)
  if (!sourceNode || !targetNode) return false
  return nodeHasTrace(sourceNode) && nodeHasTrace(targetNode)
}

// Get the trace frame for the currently selected node
const selectedTraceFrame = computed(() => {
  if (!selectedNode.value || !isTraceActive.value) return null
  return getTraceFrame(selectedNode.value) || null
})

// Get selected node - look up from nodes array to get latest version after updates
const selectedNode = computed(() => {
  const selected = getSelectedNodes.value
  if (!selected || selected.length === 0) return null
  // VueFlow's selection state may hold stale reference, look up fresh node by ID
  const node = editor.nodes.value.find(n => n.id === selected[0].id)
  // Don't open sidebar for end nodes
  if (node && node.type !== 'end') return node
  return null
})

// Compute upstream results for the selected node
const upstreamResults = computed(() => {
  if (!selectedNode.value) return []
  return editor.getUpstreamResults(selectedNode.value.id)
})

// Close reference panel when selected node changes (not on data updates)
watch(
  () => selectedNode.value?.id,
  () => {
    showReferencePanel.value = false
    activeReferenceArgument.value = null
    showTracePanel.value = true
  },
)

// Update edge highlighting when selection changes
watch(
  () => getSelectedNodes.value,
  selected => {
    const targetNode = selected?.[0]

    // First, clear all highlights
    editor.edges.value.forEach(edge => {
      edge.data = { ...edge.data, highlighted: false }
    })

    if (!targetNode) return

    // Walk backwards from selected node to find all ancestor edges
    const visited = new Set()
    function walkBackwards(nodeId) {
      if (visited.has(nodeId)) return
      visited.add(nodeId)

      const incomingEdges = editor.edges.value.filter(e => e.target === nodeId)
      incomingEdges.forEach(edge => {
        edge.data = { ...edge.data, highlighted: true }
        walkBackwards(edge.source)
      })
    }

    walkBackwards(targetNode.id)
  },
  { immediate: true, deep: true },
)

// Trigger VueFlow resize when sidebar opens/closes (watch ID only,
// not the full object reference which changes on every data update)
watch(
  () => selectedNode.value?.id,
  () => {
    nextTick(() => {
      // Dispatch resize event to trigger VueFlow's internal ResizeObserver
      window.dispatchEvent(new Event('resize'))
    })
  },
)

// Zoom controls
function zoomIn() {
  vfZoomIn()
}

function zoomOut() {
  vfZoomOut()
}

function fitToScreen() {
  fitView({ padding: 0.2, maxZoom: 1.2 })
}

function toggleEnabled(val) {
  editor.enabled.value = val
}

function togglePreviews() {
  showAllPreviews.value = !showAllPreviews.value
}

function goBack() {
  router.push('/')
}

function closeNodePicker() {
  showNodePicker.value = false
  replaceNodeId.value = null
}

function onRunClick() {
  if (editor.isDirty.value) {
    showUnsavedDialog.value = true
  } else {
    proceedRun(false)
  }
}

async function proceedRun(saveFirst) {
  showUnsavedDialog.value = false
  if (saveFirst) {
    await editor.save()
  }

  const props = editor.getTriggerProperties()
  if (props && props.length > 0) {
    runProperties.value = props
    showRunModal.value = true
  } else {
    // If no complex scopes, just execute directly
    editor.exec()
  }
}

function onRunConfirm(scope) {
  editor.exec(scope)
}

// Drawer resize handlers
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

// Center viewport on a node (preserving current zoom)
import { getNodeCenterOffset } from '@/utils/flow-constants'

function centerOnNode(node) {
  if (!node) return
  const { zoom } = getViewport()
  const offset = getNodeCenterOffset(node.type)
  setCenter(node.position.x + offset.x, node.position.y + offset.y, { zoom, duration: 300 })
}

// Event handlers
function onNodeClick({ node = {} } = {}) {
  // Ignore clicks on end nodes
  if (node.type === 'end') return

  const foundNode = editor.nodes.value.find(n => n.id === node.id)
  if (foundNode) {
    removeSelectedNodes(getSelectedNodes.value)
    addSelectedNodes([foundNode])
    centerOnNode(foundNode)
  }
}

function onPaneClick() {
  clearSelection()
}

function clearSelection() {
  removeSelectedNodes(getSelectedNodes.value)
  showReferencePanel.value = false
  activeReferenceArgument.value = null
}

// Track middle mouse click timing for double-click detection
let lastMiddleClickTime = 0

function onMiddleMouseClick(event) {
  if (event.button !== 1) return
  event.preventDefault()

  const now = Date.now()
  if (now - lastMiddleClickTime < 300) {
    fitToScreen()
    lastMiddleClickTime = 0
  } else {
    lastMiddleClickTime = now
  }
}

function onEdgeAdd(payload) {
  insertionPoint.value = { edgeId: payload.edgeId, source: payload.source, target: payload.target }
  nodePickerCategory.value = null
  showNodePicker.value = true
}

// Open node picker
function openNodePicker(payload) {
  insertionPoint.value = payload
  nodePickerCategory.value = payload.first ? 'trigger' : null
  showNodePicker.value = true
}

// Handle node selection from picker
function handleNodeSelect(nodeType) {
  if (replaceNodeId.value) {
    // Replace mode
    editor.replaceNode(replaceNodeId.value, nodeType)
    replaceNodeId.value = null
    showNodePicker.value = false
    return
  }

  const newNode = editor.addNode(nodeType, insertionPoint.value)

  // Center and select the new node
  nextTick(() => {
    if (newNode) {
      addSelectedNodes([newNode])
      centerOnNode(newNode)
    }
  })

  showNodePicker.value = false
}

// Handle replace from node context menu
function startReplace(nodeId) {
  const node = editor.nodes.value.find(n => n.id === nodeId)
  if (!node) return

  replaceNodeId.value = nodeId
  nodePickerCategory.value = node.type === 'trigger' ? 'trigger' : null
  showNodePicker.value = true
}

// Handle delete from node context menu (with confirmation)
function confirmDeleteNode(nodeId) {
  const node = editor.nodes.value.find(n => n.id === nodeId)
  if (!node) return

  confirmDelete({
    message: t('builder.confirmDelete.message'),
    header: t('builder.confirmDelete.header'),
    icon: 'pi pi-trash',
    onConfirm: () => {
      editor.deleteNode(node)
      clearSelection()
    },
  })
}

// Handle delete from sidebar or keyboard
function handleDeleteSelected() {
  const selectedNodes = getSelectedNodes.value
  if (selectedNodes.length === 0) return

  selectedNodes.forEach(node => {
    // End steps cannot be deleted directly
    if (node.type !== 'end') {
      editor.deleteNode(node)
    }
  })

  clearSelection()
}

// Handle adding a new branch output (Else If)
function handleAddBranch() {
  const selected = getSelectedNodes.value?.[0]
  if (selected && selected.type === 'branch') {
    editor.addBranchOutput(selected)
  }
}

// Handle reordering branch edges
function handleReorderBranches(edgeIds) {
  const selected = getSelectedNodes.value?.[0]
  if (selected && selected.type === 'branch') {
    editor.reorderBranchEdges(selected.id, edgeIds)
  }
}

// Handle argument updates from ConfigSidebar
function handleUpdateArguments(args) {
  const selected = getSelectedNodes.value?.[0]
  if (!selected) return

  const node = editor.nodes.value.find(n => n.id === selected.id)
  if (node) {
    editor.updateNodeData(node.id, { arguments: args })
  }
}

// Handle constraint updates from ConfigSidebar (trigger nodes)
function handleUpdateConstraints(constraints) {
  const selected = getSelectedNodes.value?.[0]
  if (!selected) return

  const node = editor.nodes.value.find(n => n.id === selected.id)
  if (node) {
    editor.updateNodeData(node.id, { constraints })
  }
}

// Handle inputSchema updates from ConfigSidebar (agent trigger only)
function handleUpdateInputSchema(inputSchema) {
  const selected = getSelectedNodes.value?.[0]
  if (!selected) return

  const node = editor.nodes.value.find(n => n.id === selected.id)
  if (node) {
    editor.updateNodeData(node.id, { inputSchema })
  }
}

// Handle gateway type change (exclusive/inclusive) from ConfigSidebar
function handleUpdateGatewayType(gatewayRef) {
  const selected = getSelectedNodes.value?.[0]
  if (!selected) return
  editor.updateGatewayType(selected.id, gatewayRef)
}

// Handle branch condition change from ConfigSidebar
function handleUpdateBranchExpr({ edgeId, condition }) {
  editor.updateEdgeCondition(edgeId, condition)
}

// Handle node metadata updates (label, description) from ConfigSidebar
function handleUpdateMetadata(meta) {
  const selected = getSelectedNodes.value?.[0]
  if (!selected) return

  const node = editor.nodes.value.find(n => n.id === selected.id)
  if (node) {
    editor.updateNodeData(node.id, meta)
  }
}

// Compute current reference (scope + source) for the active argument
const currentReference = computed(() => {
  if (!activeReferenceArgument.value) return null

  if (activeReferenceArgument.value.isCondition) {
    const { edgeId, side, rowIndex } = activeReferenceArgument.value
    const edge = editor.edges.value.find(e => e.id === edgeId)
    if (!edge?.data?.condition) return null

    const condition = edge.data.condition
    let refNode = null

    if (condition.ref === 'and' || condition.ref === 'or') {
      const row = condition.args?.[rowIndex]
      if (!row) return null
      refNode = side === 'variable' ? row.args[0] : row.args[1]
    } else {
      refNode = side === 'variable' ? condition.args[0] : condition.args[1]
    }

    if (refNode?.meta?.scope && refNode.symbol) {
      return { scope: refNode.meta.scope, source: refNode.symbol }
    }
    return null
  }

  if (!selectedNode.value) return null
  const args = selectedNode.value.data?.arguments || []
  const { name, target } = activeReferenceArgument.value

  const match = target
    ? args.find(a => a.argumentName === name && a.target === target)
    : args.find(a => a.argumentName === name)

  if (match?.scope && (match?.source || match?.expr)) {
    return { scope: match.scope, source: match.source || match.expr }
  }
  return null
})

// Reference panel handlers
function handleToggleReference(argumentInfo) {
  // Condition references come as "condition:edgeId:side:rowIndex"
  if (typeof argumentInfo === 'string' && argumentInfo.startsWith('condition:')) {
    const [, edgeId, side, rowIndex] = argumentInfo.split(':')
    activeReferenceArgument.value = {
      isCondition: true,
      edgeId,
      side,
      rowIndex: parseInt(rowIndex, 10),
    }
  } else {
    // Normal function argument reference
    activeReferenceArgument.value = argumentInfo
  }
  showReferencePanel.value = true
}

function handleReferenceSelect({ scope, source }) {
  if (!activeReferenceArgument.value) return

  if (activeReferenceArgument.value.isCondition) {
    // Apply reference to condition AST node
    const { edgeId, side, rowIndex } = activeReferenceArgument.value
    const edge = editor.edges.value.find(e => e.id === edgeId)
    if (!edge?.data?.condition) return

    const condition = JSON.parse(JSON.stringify(edge.data.condition))
    const refNode = { symbol: source, meta: { scope } }

    // If it's a group (and/or), update the specific child row
    if (condition.ref === 'and' || condition.ref === 'or') {
      const row = condition.args?.[rowIndex]
      if (!row) return
      if (side === 'variable') {
        row.args[0] = refNode
      } else {
        row.args[1] = refNode
      }
    } else {
      // Single condition (bare comparison node)
      if (side === 'variable') {
        condition.args[0] = refNode
      } else {
        condition.args[1] = refNode
      }
    }

    editor.updateEdgeCondition(edgeId, condition)
  } else {
    // Normal function argument reference
    const { name, target } = activeReferenceArgument.value
    configSidebarRef.value?.applyReference(name, { scope, source }, target)
  }
}

function closeReferencePanel() {
  showReferencePanel.value = false
  activeReferenceArgument.value = null
}

// Handle keyboard deletion
function onKeyDown(event) {
  if (event.key === 'Delete' || event.key === 'Backspace') {
    if (event.target.tagName !== 'INPUT' && event.target.tagName !== 'TEXTAREA') {
      handleDeleteSelected()
    }
  }
}

// Catalog is pre-loaded by App.vue; just register keyboard listener
onMounted(() => {
  window.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
})

// Watch route for changes (waits for catalog to be ready)
watch(
  [() => route.params.id, () => store.catalogReady],
  async ([id, ready], oldVals) => {
    if (!ready) return

    const oldId = oldVals?.[0]
    if (id !== oldId) {
      hasInitiallyFit = false
    }

    if (id && id !== 'new') {
      await editor.load(id)

      // Ensure we fit view when navigating to an already mounted component
      nextTick(() => {
        setTimeout(() => {
          if (!hasInitiallyFit) {
            fitView({ padding: 0.2, maxZoom: 1.2 })
            hasInitiallyFit = true
          }
        }, 100)
      })
    } else {
      editor.reset()
    }
  },
  { immediate: true },
)
</script>

<style>
.builder-layout {
  height: calc(100vh - var(--topbar-height));
}

/* Override VueFlow styles to match our theme */
.vue-flow {
  background-color: var(--p-surface-ground);
}

.vue-flow__background {
  background-color: transparent;
}

.vue-flow__background pattern circle {
  fill: color-mix(in srgb, var(--p-text-muted-color) 40%, transparent);
}

/* Edge styling */
.vue-flow__edge-path {
  stroke: var(--p-text-muted-color);
  stroke-width: 2;
}

.vue-flow__edge.selected .vue-flow__edge-path,
.vue-flow__edge:hover .vue-flow__edge-path {
  stroke: var(--p-primary-color);
  stroke-width: 2;
}

/* Hide node handles - edges connect invisibly */
.vue-flow__handle {
  opacity: 0;
  pointer-events: none;
}
</style>
