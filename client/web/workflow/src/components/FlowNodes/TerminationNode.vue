<template>
  <div
    class="termination-node"
    :class="{
      'termination-node--selected': selected,
      'termination-node--highlighted': data?.highlighted,
      'termination-node--trace-success': data?.traceState === 'success',
      'termination-node--trace-error': data?.traceState === 'error',
    }"
  >
    <Handle type="target" :position="Position.Top" id="target-top" :style="{ left: '50%' }" />
    <Handle type="target" :position="Position.Top" id="target-top-left" :style="{ left: '25%' }" />
    <Handle type="target" :position="Position.Top" id="target-top-right" :style="{ left: '75%' }" />
    <Handle type="target" :position="Position.Left" id="target-left" />

    <div class="termination-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="termination-node__icon" />
      <span class="termination-node__type">{{ stepTypeLabel }}</span>
    </div>

    <div v-if="data?.label" class="termination-node__label">
      {{ data.label }}
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
import { Handle, Position } from '@vue-flow/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStyleFromKind } from '../../lib/style'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
  currentTheme: { type: String, default: 'light' },
})

const { t } = useI18n()

function getIcon(name) {
  if (!name) return ''
  const basePath = `${(document.getElementsByTagName('base')[0] || {}).href || '/'}icons`
  return `${basePath}/${props.currentTheme === 'dark' ? 'dark/' : ''}${name}.svg`
}

const iconSrc = computed(() => {
  const styleInfo = getStyleFromKind({ kind: 'termination' })
  return styleInfo?.icon ? getIcon(styleInfo.icon) : ''
})

const stepTypeLabel = computed(() => t('steps.termination.short'))
</script>

<style scoped>
.termination-node {
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-border);
  border-radius: 5px;
  width: 200px;
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
  height: 36px;
  color: var(--p-primary-color);
  font-weight: 500;
}

.termination-node__icon {
  width: 20px;
  height: 20px;
  margin-right: 6px;
}

.termination-node__label {
  border-top: 1px solid var(--p-surface-border);
  padding: 6px 8px;
  font-size: 13px;
  color: var(--p-text-color);
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

/* Only target handle */
.termination-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 10px;
  height: 10px;
  transition: opacity 0.2s;
}

.termination-node:hover :deep(.vue-flow__handle) {
  opacity: 1;
}
</style>
