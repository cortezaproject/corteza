<template>
  <div
    class="termination-node"
    :class="{
      'termination-node--selected': selected,
      'termination-node--highlighted': data?.highlighted,
      'termination-node--trace-success': data?.traceState === 'success',
      'termination-node--trace-error': data?.traceState === 'error',
      'termination-node--connecting': isConnecting,
    }"
  >
    <Handle
      type="target"
      :position="Position.Top"
      id="target-top"
      :style="{ left: '50%' }"
      :connectable="!isTargetUsed('target-top')"
    />
    <Handle
      type="target"
      :position="Position.Top"
      id="target-top-left"
      :style="{ left: '25%' }"
      :connectable="!isTargetUsed('target-top-left')"
    />
    <Handle
      type="target"
      :position="Position.Top"
      id="target-top-right"
      :style="{ left: '75%' }"
      :connectable="!isTargetUsed('target-top-right')"
    />
    <Handle
      type="target"
      :position="Position.Left"
      id="target-left"
      :connectable="!isTargetUsed('target-left')"
    />
    <Handle
      type="target"
      :position="Position.Right"
      id="target-right"
      :connectable="!isTargetUsed('target-right')"
    />
    <Handle
      type="target"
      :position="Position.Bottom"
      id="target-bottom"
      :style="{ left: '50%' }"
      :connectable="!isTargetUsed('target-bottom')"
    />
    <Handle
      type="target"
      :position="Position.Bottom"
      id="target-bottom-left"
      :style="{ left: '25%' }"
      :connectable="!isTargetUsed('target-bottom-left')"
    />
    <Handle
      type="target"
      :position="Position.Bottom"
      id="target-bottom-right"
      :style="{ left: '75%' }"
      :connectable="!isTargetUsed('target-bottom-right')"
    />

    <!-- Header (icon + title) -->
    <div class="termination-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="termination-node__icon" />
      <span class="termination-node__title" :title="data?.label || stepTypeLabel">
        {{ data?.label || stepTypeLabel }}
      </span>
    </div>

    <!-- Description (bottom) — falls back to step-type label so the second row
         is never empty. User-edited description overrides. -->
    <div v-if="displayDescription" class="termination-node__description">
      {{ displayDescription }}
    </div>

    <!-- Trace badge -->
    <div v-if="data?.traceLog" class="termination-node__trace-badge">
      <img
        :src="getIcon(data.traceState === 'error' ? 'clock-danger' : 'clock-success')"
        class="termination-node__trace-icon"
      />
    </div>
  </div>
</template>

<script setup>
import { Handle, Position, useVueFlow } from '@vue-flow/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStyleFromKind } from '../../lib/style'
import { getIcon as resolveIcon } from '../../lib/icon'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
  currentTheme: { type: String, default: 'light' },
  usedTargetHandles: { type: Array, default: () => [] },
  inCount: { type: Number, default: 0 },
})

// Terminations are end-points — once a path reaches one, extra inbound arrows
// don't add meaning, so cap inbound at 1 and hide every handle past that.
const inboundFull = computed(() => props.inCount >= 1)
const isTargetUsed = id => inboundFull.value || props.usedTargetHandles.includes(id)

const { connectionStartHandle } = useVueFlow()
const isConnecting = computed(() => !!connectionStartHandle.value)

const { t } = useI18n()

const getIcon = (name) => resolveIcon(name, props.currentTheme)

const iconSrc = computed(() => {
  const styleInfo = getStyleFromKind({ kind: 'termination' })
  return styleInfo?.icon ? getIcon(styleInfo.icon) : ''
})

const stepTypeLabel = computed(() => t('steps.termination.short'))
const stepDescription = computed(() => t('steps.termination.description'))

const displayDescription = computed(() => props.data?.description || stepDescription.value)
</script>

<style scoped>
.termination-node {
  display: flex;
  flex-direction: column;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: 5px;
  width: 180px;
  min-height: 64px;
  box-shadow: var(--p-card-shadow);
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.termination-node--selected {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

.termination-node--highlighted {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

.termination-node--trace-success {
  border-color: var(--p-green-500);
  border-width: 2px;
}

.termination-node--trace-error {
  border-color: var(--p-red-500);
  border-width: 2px;
}

.termination-node__header {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  height: 32px;
  color: var(--p-primary-color);
  font-weight: 500;
  font-size: 13px;
}

.termination-node__title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--p-text-color);
}

.termination-node__icon {
  width: 24px;
  height: 24px;
  margin-right: 6px;
}

.termination-node__description {
  border-top: 1px solid var(--p-content-border-color);
  padding: 6px 8px;
  font-size: 12px;
  line-height: 16px;
  color: var(--p-text-muted-color);
  white-space: pre-wrap;
  word-break: break-word;
}

.termination-node__trace-badge {
  position: absolute;
  top: -8px;
  right: -8px;
}

.termination-node__trace-icon {
  width: 16px;
  height: 16px;
}

/* z-index:-1 tucks the circle behind the node so only the outer half pokes
   past the border. */
.termination-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 12px;
  height: 12px;
  z-index: -1;
  transition: opacity 0.2s;
}

/* Termination only reveals drop zones while another node is dragging a
   connection — never on plain hover. */
.termination-node--connecting :deep(.vue-flow__handle.connectable) {
  opacity: 1;
}
</style>
