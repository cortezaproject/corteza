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
    :style="{ width: '180px' }"
  >
    <!-- Header (icon + title + actions) -->
    <div class="trigger-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="trigger-node__icon" />
      <span class="trigger-node__title" :title="data?.label || stepTypeLabel">
        {{ data?.label || stepTypeLabel }}
      </span>

      <div class="trigger-node__header-actions">
        <img
          v-if="canTest && !dryRunProcessing"
          :src="getIcon('play')"
          class="trigger-node__action-btn"
          :title="$t('configurator.tooltip.run-workflow')"
          @click.stop="$emit('test', id)"
        />
        <span v-if="dryRunProcessing && dryRunCellID === id" class="trigger-node__spinner" />
        <img
          v-if="dryRunProcessing && dryRunCellID === id && dryRunSessionID"
          :src="getIcon('stop')"
          class="trigger-node__action-btn"
          @click.stop="$emit('cancel')"
        />
      </div>
    </div>

    <!-- Issue badge (top-right) -->
    <div
      v-if="hasIssues"
      v-tooltip.top="{ value: issueTooltip, pt: { text: 'whitespace-pre-wrap text-xs' } }"
      class="trigger-node__issue-badge"
      :aria-label="$t('editor.issues')"
      @click.stop="$emit('open-issues', id)"
    >
      <img :src="getIcon('issue')" class="trigger-node__issue-icon" alt="" />
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

    <!-- Description (bottom) — falls back to a config-derived label
         (resource · event) so the second row is never empty. -->
    <div v-if="displayDescription" class="trigger-node__description">
      {{ displayDescription }}
    </div>

    <!-- Source handles (triggers have no inputs). `:connectable="false"` once
         the trigger is wired — Human rule: one outbound per trigger. -->
    <Handle
      type="source"
      :position="Position.Top"
      id="source-top"
      :style="{ left: '50%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Top"
      id="source-top-left"
      :style="{ left: '25%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Top"
      id="source-top-right"
      :style="{ left: '75%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom"
      :style="{ left: '50%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-left"
      :style="{ left: '25%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-right"
      :style="{ left: '75%' }"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Right"
      id="source-right"
      :connectable="!hasOutgoingEdge"
    />
    <Handle
      type="source"
      :position="Position.Left"
      id="source-left"
      :connectable="!hasOutgoingEdge"
    />
  </div>
</template>

<script setup>
import { Handle, Position } from '@vue-flow/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStyleFromKind } from '../../lib/style'
import { getConstraintNameLabel } from '../../lib/constraint'
import { camelToTitle } from '../../lib/string'
import { getIcon as resolveIcon } from '../../lib/icon'

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
  hasOutgoingEdge: { type: Boolean, default: false },
})

defineEmits(['test', 'cancel', 'open-issues'])

const { t } = useI18n()

const getIcon = (name) => resolveIcon(name, props.currentTheme)

const iconSrc = computed(() => {
  const styleInfo = getStyleFromKind({ kind: 'trigger' })
  return styleInfo?.icon ? getIcon(styleInfo.icon) : ''
})

const stepTypeLabel = computed(() => t('steps.trigger.short'))
const stepDescription = computed(() => t('steps.trigger.description'))

const configLabel = computed(() => {
  const trg = props.data?.triggers || {}
  if (!trg.resourceType) return ''
  const resourceLabel = trg.resourceType
    .split(':')
    .map(part =>
      part
        .split('-')
        .map(w => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase())
        .join(' '),
    )
    .join(' - ')
  const eventLabel = trg.eventType ? camelToTitle(trg.eventType.replace('on', '')) : ''
  return eventLabel ? `${resourceLabel} · ${eventLabel}` : resourceLabel
})

const displayDescription = computed(
  () => props.data?.description || configLabel.value || stepDescription.value,
)

const isEnabled = computed(() => {
  return props.data?.triggers?.enabled !== false
})

const nodeIssues = computed(() => {
  const list = props.issues?.[props.id]
  return Array.isArray(list) ? list : []
})

const hasIssues = computed(() => nodeIssues.value.length > 0)

const issueTooltip = computed(() => nodeIssues.value.join('\n'))

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
  display: flex;
  flex-direction: column;
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-border);
  border-radius: 5px;
  width: 180px;
  min-height: 64px;
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
  height: 32px;
  color: var(--p-primary-color);
  font-weight: 500;
  font-size: 13px;
}

.trigger-node__icon {
  width: 24px;
  height: 24px;
  margin-right: 6px;
}

.trigger-node__title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--p-text-color);
}

.trigger-node__header-actions {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 4px;
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

.trigger-node__issue-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--p-content-background);
  border-radius: 50%;
  box-shadow: var(--p-card-shadow);
  cursor: pointer;
  z-index: 11;
}

.trigger-node__issue-icon {
  width: 14px;
  height: 14px;
  display: block;
}

.trigger-node__description {
  border-top: 1px solid var(--p-surface-border);
  padding: 6px 8px;
  font-size: 12px;
  line-height: 16px;
  color: var(--p-text-muted-color);
  white-space: pre-wrap;
  word-break: break-word;
}

.trigger-node__values {
  display: none;
  position: absolute;
  top: calc(100% + 14px);
  left: 0;
  width: 180px;
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

/* Hide handles by default. z-index:-1 tucks the circle behind the node so
   only the outer half pokes past the border. */
.trigger-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 12px;
  height: 12px;
  z-index: -1;
  transition: opacity 0.2s;
}

.trigger-node:hover :deep(.vue-flow__handle.connectable) {
  opacity: 1;
}
</style>
