<template>
  <g :class="{ 'flow-edge--highlighted': isHighlighted }">
    <BaseEdge :id="id" :path="path" :marker-end="markerEnd" :style="edgeStyle" />
  </g>

  <EdgeLabelRenderer>
    <div
      v-if="label"
      class="flow-edge__label"
      :style="{
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY}px)`,
        pointerEvents: 'all',
      }"
      @dblclick.stop="startEditing"
    >
      <template v-if="editing">
        <input
          ref="editInput"
          v-model="editValue"
          class="flow-edge__label-input"
          @blur="finishEditing"
          @keydown.enter="finishEditing"
          @keydown.escape="cancelEditing"
        >
      </template>
      <template v-else>
        {{ label }}
      </template>
    </div>
  </EdgeLabelRenderer>
</template>

<script setup>
import { BaseEdge, EdgeLabelRenderer, getSmoothStepPath } from '@vue-flow/core'
import { computed, nextTick, ref } from 'vue'

defineOptions({ inheritAttrs: false })

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
  label: { type: String, default: '' },
  data: { type: Object, default: () => ({}) },
  markerEnd: { type: String, default: '' },
  selected: { type: Boolean, default: false },
})

const emit = defineEmits(['update-label'])

// Edit state
const editing = ref(false)
const editValue = ref('')
const editInput = ref(null)

function startEditing () {
  editValue.value = props.label || ''
  editing.value = true
  nextTick(() => {
    editInput.value?.focus()
    editInput.value?.select()
  })
}

function finishEditing () {
  editing.value = false
  emit('update-label', { id: props.id, label: editValue.value })
}

function cancelEditing () {
  editing.value = false
}

// Path computation
const pathResult = computed(() => {
  return getSmoothStepPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    targetX: props.targetX,
    targetY: props.targetY,
    sourcePosition: props.sourcePosition,
    targetPosition: props.targetPosition,
  })
})

const path = computed(() => pathResult.value[0])
const labelX = computed(() => pathResult.value[1])
const labelY = computed(() => pathResult.value[2])

// Highlighting
const isHighlighted = computed(() => props.data?.highlighted === true)

// Edge style
const edgeStyle = computed(() => {
  if (props.data?.traceState === 'success') {
    return { stroke: 'var(--p-green-500)', strokeWidth: '2px' }
  }
  if (props.data?.traceState === 'error') {
    return { stroke: 'var(--p-red-500)', strokeWidth: '2px' }
  }
  if (isHighlighted.value || props.selected) {
    return { stroke: 'var(--p-primary-color)', strokeWidth: '2.5px' }
  }
  return { stroke: 'var(--p-text-muted-color)', strokeWidth: '2px' }
})
</script>

<style scoped>
.flow-edge__label {
  background: var(--p-content-background);
  border: 2px solid var(--p-content-border-color);
  border-radius: 5px;
  padding: 2px 12px;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  color: var(--p-text-color);
}

.flow-edge__label:hover {
  border-color: var(--p-primary-color);
}

.flow-edge__label-input {
  border: none;
  outline: none;
  font-size: 13px;
  min-width: 60px;
  max-width: 200px;
  background: transparent;
  color: var(--p-text-color);
}

/* Highlight edge path */
.flow-edge--highlighted :deep(path) {
  stroke: var(--p-primary-color) !important;
  stroke-width: 2.5px !important;
}
</style>
