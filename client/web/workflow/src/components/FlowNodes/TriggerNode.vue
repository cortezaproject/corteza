<template>
  <div
    class="trigger-node"
    :class="{
      'trigger-node--selected': selected,
      'trigger-node--disabled': !isEnabled,
      'trigger-node--highlighted': data?.highlighted,
      'trigger-node--trace-success': data?.traceState === 'success',
      'trigger-node--trace-error': data?.traceState === 'error',
    }"
    :style="{ width: '200px' }"
  >
    <!-- Header -->
    <div class="trigger-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="trigger-node__icon" />
      <span class="trigger-node__type">{{ stepTypeLabel }}</span>

      <div class="trigger-node__header-actions">
        <!-- Test button -->
        <img
          v-if="canTest && !dryRunProcessing"
          :src="getIcon('play')"
          class="trigger-node__action-btn trigger-node__hover-show"
          :title="$t('configurator.tooltip.run-workflow')"
          @click.stop="$emit('test', id)"
        />
        <!-- Spinner -->
        <span v-if="dryRunProcessing && dryRunCellID === id" class="trigger-node__spinner" />
        <!-- Cancel -->
        <img
          v-if="dryRunProcessing && dryRunCellID === id && dryRunSessionID"
          :src="getIcon('stop')"
          class="trigger-node__action-btn"
          @click.stop="$emit('cancel')"
        />
        <!-- Issue badge -->
        <img
          v-if="hasIssues"
          :src="getIcon('issue')"
          class="trigger-node__issue"
          @click.stop="$emit('open-issues', id)"
        />
        <!-- ID label -->
        <span v-if="!hasIssues" class="trigger-node__id">{{ id }}</span>
      </div>
    </div>

    <!-- Label row -->
    <div class="trigger-node__label">
      <span class="trigger-node__label-text">{{ data?.label || '/' }}</span>
    </div>

    <!-- Values table (trigger config preview) -->
    <div v-if="valueRows.length > 0" class="trigger-node__values">
      <table>
        <tr v-for="(row, idx) in valueRows" :key="idx" :class="row.class">
          <td v-for="(cell, ci) in row.cells" :key="ci" v-html="cell" />
        </tr>
      </table>
    </div>

    <!-- Trace badge -->
    <div v-if="data?.traceLog" class="trigger-node__trace-badge">
      <img
        :src="getIcon(data.traceState === 'error' ? 'clock-danger' : 'clock-success')"
        class="trigger-node__trace-icon"
      />
    </div>

    <!-- Only source handles (triggers have no inputs) -->
    <Handle type="source" :position="Position.Bottom" id="source-bottom" :style="{ left: '50%' }" />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-left"
      :style="{ left: '25%' }"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-right"
      :style="{ left: '75%' }"
    />
    <Handle type="source" :position="Position.Right" id="source-right" />
  </div>
</template>

<script setup>
import { Handle, Position } from '@vue-flow/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStyleFromKind } from '../../lib/style'
import { getConstraintNameLabel } from '../../lib/constraint'
import { camelToTitle } from '../../lib/string'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
  issues: { type: Object, default: () => ({}) },
  eventTypes: { type: Array, default: () => [] },
  canTest: { type: Boolean, default: false },
  dryRunProcessing: { type: Boolean, default: false },
  dryRunCellID: { type: String, default: '' },
  dryRunSessionID: { type: String, default: '' },
  currentTheme: { type: String, default: 'light' },
})

defineEmits(['test', 'cancel', 'open-issues'])

const { t } = useI18n()

function getIcon(name) {
  if (!name) return ''
  const basePath = `${(document.getElementsByTagName('base')[0] || {}).href || '/'}icons`
  return `${basePath}/${props.currentTheme === 'dark' ? 'dark/' : ''}${name}.svg`
}

const iconSrc = computed(() => {
  const styleInfo = getStyleFromKind({ kind: 'trigger' })
  return styleInfo?.icon ? getIcon(styleInfo.icon) : ''
})

const stepTypeLabel = computed(() => t('steps.trigger.short', 'Trigger'))

const isEnabled = computed(() => {
  return props.data?.triggers?.enabled !== false
})

const hasIssues = computed(() => {
  return !!(props.issues && props.issues[props.id])
})

const encodeHTML = (value = '') => {
  if (!value) return value
  return value.replace(/[\u00A0-\u9999<>&]/gim, i => '&#' + i.charCodeAt(0) + ';')
}

const valueRows = computed(() => {
  const rows = []
  const trg = props.data?.triggers || {}

  if (!trg.resourceType) return rows

  let resourceLabel = trg.resourceType
    .split(':')
    .map(part =>
      part
        .split('-')
        .map(w => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase())
        .join(' '),
    )
    .join(' - ')

  let eventLabel = trg.eventType ? camelToTitle(trg.eventType.replace('on', '')) : ''

  rows.push({ class: 'title', cells: ['<b>Configuration</b>', '', ''] })
  rows.push({
    class: '',
    cells: ['<var>Resource</var>', '', `<code>${encodeHTML(resourceLabel)}</code>`],
  })
  rows.push({
    class: '',
    cells: ['<var>Event</var>', '', `<code>${encodeHTML(eventLabel)}</code>`],
  })

  // Constraints
  const constraints = trg.constraints || []
  if (constraints.length && trg.eventType !== 'onManual') {
    rows.push({ class: 'title', cells: ['<b>Constraints</b>', '', ''] })
    constraints.forEach(({ name = '', op = '', values = '' }) => {
      const vs = Array.isArray(values) ? values.join(' or ') : values
      rows.push({
        class: '',
        cells: [
          `<samp>${getConstraintNameLabel(name)}</samp>`,
          `<samp>${op}</samp>`,
          `<code>${encodeHTML(vs)}</code>`,
        ],
      })
    })
  }

  // Initial scope properties
  const et = props.eventTypes.find(
    et => trg.resourceType === et.resourceType && trg.eventType === et.eventType,
  )
  if (et?.properties?.length) {
    rows.push({ class: 'title', cells: ['<b>Initial scope</b>', '', ''] })
    et.properties.forEach(({ name = '', type = '' }) => {
      rows.push({ class: '', cells: [`<var>${name}</var>`, '', `<samp>${type || 'Any'}</samp>`] })
    })
  }

  return rows
})
</script>

<style scoped>
.trigger-node {
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-border);
  border-radius: 5px;
  width: 200px;
  box-shadow: var(--p-card-shadow);
  cursor: pointer;
  transition:
    box-shadow 0.2s,
    border-color 0.2s,
    opacity 0.2s;
}

.trigger-node--disabled {
  opacity: 0.7;
}

.trigger-node--selected {
  border-color: var(--p-primary-color);
  border-width: 2px;
  box-shadow: 0 2px 8px color-mix(in srgb, var(--p-primary-color) 25%, transparent);
}

.trigger-node--highlighted {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

.trigger-node--trace-success {
  border-color: var(--p-green-500);
  border-width: 2px;
}

.trigger-node--trace-error {
  border-color: var(--p-red-500);
  border-width: 2px;
}

.trigger-node__header {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  height: 36px;
  color: var(--p-primary-color);
  font-weight: 500;
}

.trigger-node__icon {
  width: 20px;
  height: 20px;
  margin-right: 6px;
}

.trigger-node__type {
  flex: 1;
}

.trigger-node__header-actions {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 4px;
}

.trigger-node__hover-show {
  display: none;
}

.trigger-node:hover .trigger-node__hover-show {
  display: block;
}

.trigger-node__action-btn {
  width: 20px;
  height: 20px;
  cursor: pointer;
}

.trigger-node__spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 0.2em solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: trigger-spin 0.75s linear infinite;
  color: var(--p-text-muted-color);
}

@keyframes trigger-spin {
  to {
    transform: rotate(360deg);
  }
}

.trigger-node__id {
  font-size: 8px;
  opacity: 0.6;
}

.trigger-node:hover .trigger-node__id {
  display: none;
}

.trigger-node__issue {
  width: 20px;
  cursor: pointer;
}

.trigger-node__label {
  border-top: 1px solid var(--p-surface-border);
  padding: 6px 8px;
  min-height: 36px;
  display: flex;
  align-items: center;
  background: var(--p-content-background);
}

.trigger-node__label-text {
  text-align: left;
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
  min-width: 0;
  color: var(--p-text-color);
}

.trigger-node:hover .trigger-node__label-text {
  white-space: normal;
}

.trigger-node__values {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  z-index: 10;
}

.trigger-node:hover .trigger-node__values {
  display: block;
}

.trigger-node__values table {
  width: 100%;
  border-collapse: collapse;
  background: var(--p-content-background);
  border-radius: 0 0 5px 5px;
  box-shadow: var(--p-card-shadow);
}

.trigger-node__values td {
  text-align: left;
  padding: 6px 8px;
  white-space: nowrap;
  font-size: 12px;
}

.trigger-node__values tr.title {
  background-color: var(--p-highlight-background);
}

.trigger-node__trace-badge {
  position: absolute;
  top: -8px;
  right: -8px;
}

.trigger-node__trace-icon {
  width: 16px;
  height: 16px;
}

/* Hide handles by default */
.trigger-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 10px;
  height: 10px;
  transition: opacity 0.2s;
}

.trigger-node:hover :deep(.vue-flow__handle) {
  opacity: 1;
}
</style>
