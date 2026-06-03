<script setup>
import { BaseEdge, EdgeLabelRenderer, getSmoothStepPath, useVueFlow } from '@vue-flow/core';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineOptions({
  inheritAttrs: false,
})

const props = defineProps({
  id: { type: String, required: true },
  source: { type: String, required: true },
  target: { type: String, required: true },
  sourceX: { type: Number, required: true },
  sourceY: { type: Number, required: true },
  targetX: { type: Number, required: true },
  targetY: { type: Number, required: true },
  sourcePosition: { type: String, required: true },
  targetPosition: { type: String, required: true },
  sourceHandleId: { type: String, default: null },
  data: { type: Object, default: () => ({}) },
  markerEnd: { type: String, default: '' },
  traceActive: { type: Boolean, default: false },
  traceTraversed: { type: Boolean, default: false },
})

const emit = defineEmits(['add'])

// Get nodes to check source type
const { getNodes, getEdges } = useVueFlow()

// Constants for branch edge layout
const FORK_DROP = 30      // How far down before splitting

// Determine if this is a branch edge (source is a branch node)
const isBranchEdge = computed(() => {
  const sourceNode = getNodes.value.find(n => n.id === props.source)
  return sourceNode?.type === 'branch'
})

// Determine if this is an iterator edge (source is an iterator node)
const isIteratorEdge = computed(() => {
  const sourceNode = getNodes.value.find(n => n.id === props.source)
  return sourceNode?.type === 'iterator'
})

// Whether to use forked path rendering (branch or iterator)
const isForkedEdge = computed(() => isBranchEdge.value || isIteratorEdge.value)

// Get the index of this edge among siblings from same source (for labeling)
const edgeIndex = computed(() => {
  if (!isForkedEdge.value) return -1
  const siblingEdges = getEdges.value.filter(e => e.source === props.source)
  return siblingEdges.findIndex(e => e.id === props.id)
})

// Calculate the path and center position
const path = computed(() => {
  if (isForkedEdge.value) {
    // Custom forked path for branch/iterator edges
    // Shape: down from center -> horizontal to target X -> down to target
    const startX = props.sourceX
    const startY = props.sourceY
    const forkY = startY + FORK_DROP
    const endX = props.targetX
    const endY = props.targetY

    // Path: start -> down to fork -> horizontal to endX -> down to target
    const edgePath = `M ${startX} ${startY}
                      L ${startX} ${forkY}
                      L ${endX} ${forkY}
                      L ${endX} ${endY}`

    // Label near top of vertical segment, button centered
    const labelX = endX
    const labelY = forkY + 20  // 20px from top of vertical segment
    const buttonY = forkY + (endY - forkY) / 2  // Centered on vertical segment

    return { edgePath, labelX, labelY, buttonY }
  } else {
    // Regular smooth step path for non-forked edges
    const [edgePath, labelX, labelY] = getSmoothStepPath({
      sourceX: props.sourceX,
      sourceY: props.sourceY,
      targetX: props.targetX,
      targetY: props.targetY,
      sourcePosition: props.sourcePosition,
      targetPosition: props.targetPosition,
    })
    return { edgePath, labelX, labelY, buttonY: labelY }
  }
})

// Edge label for branches: If (first), Else (second/last)
// Edge label for iterators: Body (first), Done (second/last)
const edgeLabel = computed(() => {
  if (isBranchEdge.value) {
    const siblingEdges = getEdges.value.filter(e => e.source === props.source)
    if (edgeIndex.value === 0) return t('builder.branch.if')
    if (edgeIndex.value === siblingEdges.length - 1) return t('builder.branch.else')
    return t('builder.branch.elseIf')
  }
  if (isIteratorEdge.value) {
    const siblingEdges = getEdges.value.filter(e => e.source === props.source)
    if (edgeIndex.value === 0) return ''
    if (edgeIndex.value === siblingEdges.length - 1) return t('builder.iterator.done')
    return ''
  }
  return null
})

// Check if this edge is on the highlighted path
const isHighlighted = computed(() => props.data?.highlighted === true)

// Compute edge style for trace mode
const edgeStyle = computed(() => {
  if (props.traceActive) {
    if (props.traceTraversed) {
      return { stroke: 'var(--p-green-500)', strokeWidth: '2.5px' }
    }
  }
  if (isHighlighted.value) {
    return { stroke: 'var(--p-primary-color)', strokeWidth: '2.5px' }
  }
  return {}
})

function handleAdd() {
  emit('add', { edgeId: props.id, source: props.source, target: props.target })
}
</script>

<template>
  <!-- The edge path wrapped for highlighting -->
  <g :class="{ 'highlighted-edge': isHighlighted && !traceActive }">
    <BaseEdge
      :id="id"
      :path="path.edgePath"
      :marker-end="markerEnd"
      :style="edgeStyle"
    />
  </g>

  <!-- Edge label and + button -->
  <EdgeLabelRenderer>
    <!-- Label for branch paths (positioned separately) -->
    <span
      v-if="edgeLabel"
      class="edge-label"
      :style="{
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${path.labelX}px, ${path.labelY}px)`,
        pointerEvents: 'all',
      }"
    >
      {{ edgeLabel }}
    </span>

    <!-- Add button (centered, hidden during trace) -->
    <div
      v-if="!traceActive"
      class="edge-button"
      :style="{
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${path.labelX}px, ${path.buttonY}px)`,
        pointerEvents: 'all',
      }"
    >
      <Button
        icon="pi pi-plus"
        size="small"
        severity="secondary"
        class="!p-1 !min-w-0 !w-7 !h-7 bg-surface"
        @click.stop="handleAdd"
      />
    </div>
  </EdgeLabelRenderer>
</template>

<style scoped>
.edge-content {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.edge-label {
  font-size: 12px;
  color: var(--p-text-muted-color);
  font-weight: 500;
  background: var(--p-content-background);
  padding: 2px 8px;
  border-radius: 4px;
}

/* Highlighted path styling */
.highlighted-edge :deep(path) {
  stroke: var(--p-primary-color) !important;
  stroke-width: 2.5 !important;
}
</style>
