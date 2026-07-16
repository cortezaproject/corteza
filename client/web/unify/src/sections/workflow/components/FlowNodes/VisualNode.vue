<template>
  <div
    class="visual-node"
    :class="{
      'visual-node--selected': selected,
      'visual-node--swimlane': isSwimlane,
      'visual-node--content': isContent,
      'visual-node--drop-target': dropTarget && isSwimlane,
    }"
  >
    <NodeResizer
      :is-visible="selected"
      :min-width="160"
      :min-height="80"
      :line-style="{ borderColor: 'var(--p-primary-color)' }"
      :handle-style="{ backgroundColor: 'var(--p-primary-color)', borderColor: 'transparent', width: '9px', height: '9px' }"
      @resize="onResize"
    />

    <div v-if="isSwimlane" class="visual-node__swimlane-label">
      {{ data?.label || '' }}
    </div>
    <div
      v-else-if="data?.label"
      class="visual-node__content"
      v-html="data.label"
    />
    <div v-else class="visual-node__content visual-node__content--placeholder">
      {{ t('steps.content.placeholder') }}
    </div>
  </div>
</template>

<script setup>
import { NodeResizer } from '@vue-flow/node-resizer'
import '@vue-flow/node-resizer/dist/style.css'
import { useVueFlow } from '@vue-flow/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import eventBus from '../../lib/eventBus'

const { t } = useI18n()

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
  dropTarget: { type: Boolean, default: false },
})

const { findNode } = useVueFlow()

const isSwimlane = computed(() => props.data?.ref === 'swimlane')
const isContent = computed(() => props.data?.ref === 'content')

function onResize ({ params }) {
  const { width, height } = params || {}
  if (width == null || height == null) return
  const node = findNode(props.id)
  if (!node) return
  if (!node.data) node.data = {}
  node.data.width = width
  node.data.height = height
  node.style = { ...(node.style || {}), width: `${width}px`, height: `${height}px` }
  eventBus.emit('change-detected')
}
</script>

<style scoped>
.visual-node {
  border-radius: 8px;
  position: relative;
  box-sizing: border-box;
  width: 100%;
  height: 100%;
}

/* Keep NodeResizer corner/edge handles fully visible (they sit on the node
   border — overflow:hidden would clip them). */
.visual-node :deep(.vue-flow__resize-control) {
  z-index: 5;
}

/* mxGraph-style swimlane: title bar along the left (solid), body transparent
   so nodes dropped onto it read against the canvas background. */
.visual-node--swimlane {
  background: transparent;
  border: 1px solid var(--p-surface-400, var(--p-text-muted-color));
  display: flex;
  flex-direction: row;
}

.visual-node--swimlane .visual-node__swimlane-label {
  writing-mode: vertical-lr;
  text-orientation: mixed;
  transform: rotate(180deg);
  padding: 10px 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--p-text-color);
  background: var(--p-content-background);
  border-right: 1px solid var(--p-content-border-color);
  border-radius: 0 6px 6px 0;
  flex: 0 0 32px;
  text-align: center;
}

.visual-node--content {
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-400, var(--p-text-muted-color));
}

.visual-node--selected {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

/* Drop-target highlight shown while a node is dragged over this swimlane. */
.visual-node--drop-target {
  border-color: var(--p-primary-color);
  border-width: 2px;
  background: color-mix(in srgb, var(--p-primary-color) 8%, transparent);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--p-primary-color) 20%, transparent);
}

.visual-node__content {
  padding: 10px;
  font-size: 13px;
  color: var(--p-text-color);
  overflow: auto;
  word-break: break-word;
}

.visual-node__content--placeholder {
  color: var(--p-text-muted-color);
  font-style: italic;
}

/* Basic typographic styles so rich-text output (from any WYSIWYG) renders
   with sensible defaults. Scoped to .visual-node__content only. */
.visual-node__content :deep(h1) { font-size: 1.4em; font-weight: 700; margin: 0 0 0.4em; }
.visual-node__content :deep(h2) { font-size: 1.2em; font-weight: 700; margin: 0 0 0.35em; }
.visual-node__content :deep(h3) { font-size: 1.05em; font-weight: 700; margin: 0 0 0.3em; }
.visual-node__content :deep(p) { margin: 0 0 0.5em; }
.visual-node__content :deep(ul),
.visual-node__content :deep(ol) { margin: 0 0 0.5em; padding-left: 1.4em; }
.visual-node__content :deep(li) { margin: 0.15em 0; }
.visual-node__content :deep(a) { color: var(--p-primary-color); text-decoration: underline; }
.visual-node__content :deep(strong) { font-weight: 700; }
.visual-node__content :deep(em) { font-style: italic; }
.visual-node__content :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  background: var(--p-surface-100, rgba(0,0,0,0.05));
  padding: 0 0.2em;
  border-radius: 3px;
}
.visual-node__content :deep(blockquote) {
  margin: 0 0 0.5em;
  padding-left: 0.8em;
  border-left: 3px solid var(--p-content-border-color);
  color: var(--p-text-muted-color, var(--p-text-color));
}
</style>
