<template>
  <div
    class="visual-node"
    :class="{
      'visual-node--selected': selected,
      'visual-node--swimlane': isSwimlane,
      'visual-node--content': isContent,
    }"
    :style="nodeStyle"
  >
    <NodeResizer
      :is-visible="selected"
      :min-width="200"
      :min-height="100"
      :line-style="{ borderColor: 'var(--p-primary-color)' }"
      :handle-style="{ backgroundColor: 'var(--p-primary-color)', borderColor: 'transparent', width: '9px', height: '9px' }"
      @resize="onResize"
    />

    <div v-if="isSwimlane" class="visual-node__swimlane-label">
      {{ data?.label || '' }}
    </div>
    <div v-else class="visual-node__content">
      {{ data?.label || '' }}
    </div>
  </div>
</template>

<script setup>
import { NodeResizer } from '@vue-flow/node-resizer'
import '@vue-flow/node-resizer/dist/style.css'
import { useVueFlow } from '@vue-flow/core'
import { computed } from 'vue'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})

const { updateNodeData } = useVueFlow()

const isSwimlane = computed(() => props.data?.ref === 'swimlane')
const isContent = computed(() => props.data?.ref === 'content')

const nodeStyle = computed(() => {
  return {
    width: `${props.data?.width || 400}px`,
    height: `${props.data?.height || 240}px`,
  }
})

function onResize ({ width, height }) {
  updateNodeData(props.id, { width, height })
}
</script>

<style scoped>
.visual-node {
  border-radius: 5px;
  position: relative;
}

.visual-node--swimlane {
  background: var(--p-content-background);
  border: 1px solid var(--p-text-color);
}

.visual-node--content {
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-border);
}

.visual-node--selected {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

.visual-node__swimlane-label {
  writing-mode: vertical-lr;
  text-orientation: mixed;
  padding: 8px;
  font-size: 15px;
  font-weight: 500;
  color: var(--p-text-color);
}

.visual-node__content {
  padding: 10px;
  font-size: 13px;
  color: var(--p-text-color);
  overflow: hidden;
  word-break: break-word;
  white-space: pre-wrap;
}
</style>
