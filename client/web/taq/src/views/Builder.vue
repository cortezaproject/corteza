<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ editor.name.value }}</span>
  </Teleport>

  <!-- Loading state -->
  <div v-if="editor.loading.value" class="h-full flex items-center justify-center">
    <ProgressSpinner />
  </div>

  <div v-else class="builder-layout h-full flex flex-col relative overflow-hidden">
    <!-- Enabled toggle -->
    <div
      class="absolute top-3 left-3 z-20 flex items-center gap-3 bg-surface rounded-lg border border-surface px-3 py-2 shadow-sm"
    >
      <div class="flex items-center gap-2">
        <ToggleSwitch
          :model-value="editor.enabled.value"
          @update:model-value="editor.enabled.value = $event"
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
          @click="editor.exec"
        />
      </template>
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
          <TriggerNode v-bind="props" @delete="confirmDeleteNode" />
        </template>
        <template #node-step="props">
          <StepNode v-bind="props" @delete="confirmDeleteNode" />
        </template>
        <template #node-branch="props">
          <BranchNode v-bind="props" @delete="confirmDeleteNode" />
        </template>
        <template #node-end="props">
          <EndNode v-bind="props" />
        </template>

        <!-- Custom edge with + button -->
        <template #edge-addable="props">
          <AddableEdge v-bind="props" @add="onEdgeAdd" />
        </template>
      </VueFlow>
    </div>

    <!-- Bottom Toolbar -->
    <div class="shrink-0 z-10 body-bg">
      <CToolbar>
        <template #start>
          <CRouterLinkButton to="/" icon="pi pi-arrow-left" severity="secondary" />
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
        nodePickerCategory === 'trigger' ? $t('builder.selectTrigger') : $t('builder.addStep')
      "
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
      :pt="{ content: { class: 'p-0' } }"
      :style="{ width: '50vw' }"
    >
      <NodePicker
        :filter-category="nodePickerCategory"
        @select="handleNodeSelect"
        @close="showNodePicker = false"
      />
    </Dialog>

    <!-- Reference Panel (opens on input click, closes via button or sidebar close) -->
    <Transition
      enter-active-class="transition-transform duration-200 ease-out"
      enter-from-class="translate-x-[280px]"
      enter-to-class="translate-x-0"
    >
      <div
        v-if="selectedNode && showReferencePanel"
        class="absolute top-0 bottom-[58px] my-7 bg-surface border-l border-surface shadow z-30 rounded-xl"
        :style="{ right: `calc(${drawerWidth}px + 1rem)`, width: '280px' }"
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
        class="config-drawer flex absolute top-0 right-0 bottom-[58px] m-3 bg-surface border-l border-surface shadow z-30 rounded-xl"
        :style="{ width: `${drawerWidth}px` }"
      >
        <!-- Resize handle -->
        <div
          class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors"
          @mousedown="startDrawerResize"
        />
        <!-- Drawer content -->
        <div class="flex-1 overflow-auto p-4">
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
            @toggle-reference="handleToggleReference"
          />
        </div>
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

const { CToolbar, CRouterLinkButton } = components

import { useConfirmDelete } from '@cortezaproject/corteza-vue-next'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import ConfigSidebar from '@/components/builder/ConfigSidebar.vue'
import NodePicker from '@/components/builder/NodePicker.vue'
import ReferencePanel from '@/components/builder/ReferencePanel.vue'
import AddableEdge from '@/components/flow/AddableEdge.vue'
import BranchNode from '@/components/flow/BranchNode.vue'
import EndNode from '@/components/flow/EndNode.vue'
import StepNode from '@/components/flow/StepNode.vue'
import TriggerNode from '@/components/flow/TriggerNode.vue'

const route = useRoute()
const store = useAutomationStore()
const editor = useFlowEditor()
const { confirmDelete } = useConfirmDelete()
const { t } = useI18n()

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
    fitView({ padding: 0.2, maxZoom: 1.5 })
    hasInitiallyFit = true
  }
})

// Drawer state
const drawerWidth = ref(360)
const isResizingDrawer = ref(false)

// Node picker state
const showNodePicker = ref(false)
const nodePickerCategory = ref(null)
const insertionPoint = ref(null)

// Reference panel state
const showReferencePanel = ref(false)
const activeReferenceArgument = ref(null)
const configSidebarRef = ref(null)

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

// Trigger VueFlow resize when sidebar opens/closes
watch(selectedNode, () => {
  nextTick(() => {
    // Dispatch resize event to trigger VueFlow's internal ResizeObserver
    window.dispatchEvent(new Event('resize'))
  })
})

// Zoom controls
function zoomIn() {
  vfZoomIn()
}

function zoomOut() {
  vfZoomOut()
}

function fitToScreen() {
  fitView({ padding: 0.2, zoom: 1.5 })
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

// Compute current reference (scope + source) for the active argument
const currentReference = computed(() => {
  if (!activeReferenceArgument.value || !selectedNode.value) return null
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
  // argumentInfo is { name: string, types: string[] } from FunctionForm
  activeReferenceArgument.value = argumentInfo
  showReferencePanel.value = true
}

function handleReferenceSelect({ scope, source }) {
  if (!activeReferenceArgument.value) return
  const { name, target } = activeReferenceArgument.value
  configSidebarRef.value?.applyReference(name, { scope, source }, target)
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
  async ([id, ready]) => {
    if (!ready) return
    if (id && id !== 'new') {
      await editor.load(id)
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
