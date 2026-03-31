<template>
  <div
    class="workflow-node"
    :class="{
      'workflow-node--selected': selected,
      'workflow-node--highlighted': data?.highlighted,
      'workflow-node--trace-success': data?.traceState === 'success',
      'workflow-node--trace-error': data?.traceState === 'error',
    }"
    :style="{ width: '200px' }"
  >
    <!-- Target handles (top + left) -->
    <Handle type="target" :position="Position.Top" id="target-top" :style="{ left: '50%' }" />
    <Handle type="target" :position="Position.Top" id="target-top-left" :style="{ left: '25%' }" />
    <Handle type="target" :position="Position.Top" id="target-top-right" :style="{ left: '75%' }" />
    <Handle type="target" :position="Position.Left" id="target-left" />
    <Handle type="target" :position="Position.Bottom" id="target-bottom" :style="{ left: '50%' }" />

    <!-- Header -->
    <div class="workflow-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="workflow-node__icon" />
      <span class="workflow-node__type">{{ stepTypeLabel }}</span>

      <div class="workflow-node__header-actions">
        <!-- Test button (trigger only) -->
        <!-- Issue badge -->
        <img
          v-if="hasIssues"
          :src="issueIcon"
          class="workflow-node__issue"
          @click.stop="$emit('open-issues', id)"
        />
        <!-- ID label -->
        <span v-if="!hasIssues" class="workflow-node__id">{{ id }}</span>
      </div>
    </div>

    <!-- Label row -->
    <div class="workflow-node__label">
      <span class="workflow-node__label-text">{{ data?.label || '/' }}</span>
    </div>

    <!-- Values table (arguments/results preview) -->
    <div v-if="valueRows.length > 0" class="workflow-node__values">
      <table>
        <tr v-for="(row, idx) in valueRows" :key="idx" :class="row.class">
          <td v-for="(cell, ci) in row.cells" :key="ci" v-html="cell" />
        </tr>
      </table>
    </div>

    <!-- Trace badge -->
    <div v-if="data?.traceLog" class="workflow-node__trace-badge">
      <img
        :src="getIcon(data.traceState === 'error' ? 'clock-danger' : 'clock-success')"
        class="workflow-node__trace-icon"
      />
    </div>

    <!-- Source handles (bottom + right + left) -->
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
    <Handle type="source" :position="Position.Left" id="source-left" />
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
  issues: { type: Object, default: () => ({}) },
  functionTypes: { type: Array, default: () => [] },
  eventTypes: { type: Array, default: () => [] },
  currentTheme: { type: String, default: 'light' },
})

defineEmits(['open-issues'])

const { t } = useI18n()

function getIcon(name) {
  if (!name) return ''
  const basePath = `${(document.getElementsByTagName('base')[0] || {}).href || '/'}icons`
  return `${basePath}/${props.currentTheme === 'dark' ? 'dark/' : ''}${name}.svg`
}

const iconSrc = computed(() => {
  const styleInfo = getStyleFromKind(props.data)
  return styleInfo?.icon ? getIcon(styleInfo.icon) : ''
})

const issueIcon = computed(() => getIcon('issue'))

const stepTypeLabel = computed(() => {
  const style = getStyleFromKind(props.data)?.style || props.data?.kind || ''
  return t(`steps.${style}.short`, style)
})

const hasIssues = computed(() => {
  return !!(props.issues && props.issues[props.id])
})

/**
 * Build value rows for the values table preview.
 */
const valueRows = computed(() => {
  const { kind, ref, arguments: args = [], results = [] } = props.data || {}
  const rows = []

  const encodeHTML = (value = '') => {
    if (!value) return value
    return value.replace(/[\u00A0-\u9999<>&]/gim, i => '&#' + i.charCodeAt(0) + ';')
  }

  if (kind === 'gateway' && ['excl', 'incl'].includes(ref)) {
    // Gateway edges are shown external to this node
    return rows
  }

  if (
    ['expressions', 'function', 'prompt', 'iterator', 'exec-workflow', 'error-handler'].includes(
      kind,
    )
  ) {
    // Function label
    const fnDef = props.functionTypes.find(f => f.ref === (ref || kind))
    const fnLabel = fnDef?.meta?.short
    if (fnLabel) {
      rows.push({ class: '', cells: [`<b class="text-primary">${encodeHTML(fnLabel)}</b>`, ''] })
    }

    if (kind === 'expressions') {
      args.forEach(({ target, expr, type }) => {
        rows.push({
          class: '',
          cells: [
            `<var>${encodeHTML(target)}</var> <samp>(${type || ''})</samp>`,
            `<code>${encodeHTML(expr)}</code>`,
          ],
        })
      })
    } else {
      // Arguments
      const params = fnDef?.parameters || []
      if (params.length && args.length) {
        rows.push({ class: 'title', cells: ['<b>Arguments</b>', ''] })
        params.forEach(({ name, types = [] }) => {
          const arg = args.find(a => a.target === name)
          const exprType = arg?.type || types[0] || ''
          rows.push({
            class: '',
            cells: [
              `<var>${encodeHTML(name)}</var> <samp>(${exprType})</samp>`,
              `<code>${encodeHTML(arg?.expr || arg?.value || '')}</code>`,
            ],
          })
        })
      }
      // Results
      if (results.length) {
        rows.push({ class: 'title', cells: ['<b>Results</b>', ''] })
        results.forEach(({ target = '', expr = '', value = '' }) => {
          rows.push({
            class: '',
            cells: [
              `<code>${encodeHTML(target)}</code>`,
              `<var>${encodeHTML(expr || value)}</var>`,
            ],
          })
        })
      }
    }
  } else if (kind === 'trigger') {
    // Trigger config preview
    const trg = props.data?.triggers || {}
    if (trg.resourceType) {
      rows.push({ class: 'title', cells: ['<b>Configuration</b>', ''] })
      rows.push({
        class: '',
        cells: ['<var>Resource</var>', `<code>${encodeHTML(trg.resourceType)}</code>`],
      })
      if (trg.eventType) {
        rows.push({
          class: '',
          cells: ['<var>Event</var>', `<code>${encodeHTML(trg.eventType)}</code>`],
        })
      }
    }
  } else if (['error', 'delay'].includes(kind)) {
    const arg = args[0]
    if (arg?.target) {
      rows.push({
        class: '',
        cells: [
          `<var>${encodeHTML(arg.target)}</var>`,
          `<code>${encodeHTML(arg.expr || arg.value || '')}</code>`,
        ],
      })
    }
  }

  return rows
})
</script>

<style scoped>
.workflow-node {
  background: var(--p-content-background);
  border: 1px solid var(--p-surface-border);
  border-radius: 5px;
  width: 200px;
  box-shadow: var(--p-card-shadow);
  cursor: pointer;
  transition:
    box-shadow 0.2s,
    border-color 0.2s;
}

.workflow-node--selected {
  border-color: var(--p-primary-color);
  border-width: 2px;
  box-shadow: 0 2px 8px color-mix(in srgb, var(--p-primary-color) 25%, transparent);
}

.workflow-node--highlighted {
  border-color: var(--p-primary-color);
  border-width: 2px;
}

.workflow-node--trace-success {
  border-color: var(--p-green-500);
  border-width: 2px;
}

.workflow-node--trace-error {
  border-color: var(--p-red-500);
  border-width: 2px;
}

.workflow-node__header {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  height: 36px;
  color: var(--p-primary-color);
  font-weight: 500;
}

.workflow-node__icon {
  width: 20px;
  height: 20px;
  margin-right: 6px;
  object-fit: contain;
}

.workflow-node__type {
  flex: 1;
  truncate: true;
}

.workflow-node__header-actions {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 4px;
}

.workflow-node__id {
  font-size: 8px;
  opacity: 0.6;
}

.workflow-node:hover .workflow-node__id {
  display: none;
}

.workflow-node__issue {
  width: 20px;
  cursor: pointer;
}

.workflow-node__label {
  border-top: 1px solid var(--p-surface-border);
  padding: 6px 8px;
  min-height: 36px;
  display: flex;
  align-items: center;
  background: var(--p-content-background);
}

.workflow-node__label-text {
  text-align: left;
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
  min-width: 0;
  color: var(--p-text-color);
}

.workflow-node:hover .workflow-node__label-text {
  white-space: normal;
  text-overflow: clip;
}

.workflow-node__values {
  display: none;
  position: absolute;
  top: calc(100% + 14px);
  left: 0;
  width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  z-index: 10;
}

.workflow-node:hover .workflow-node__values {
  display: block;
}

.workflow-node__values table {
  width: 100%;
  border-collapse: collapse;
  background: var(--p-content-background);
  border-radius: 0 0 5px 5px;
  box-shadow: var(--p-card-shadow);
}

.workflow-node__values td {
  text-align: left;
  padding: 6px 8px;
  white-space: nowrap;
  font-size: 12px;
}

.workflow-node__values tr.title {
  background-color: var(--p-highlight-background);
}

.workflow-node__trace-badge {
  position: absolute;
  top: -8px;
  right: -8px;
}

.workflow-node__trace-icon {
  width: 16px;
  height: 16px;
}

/* Hide handles by default, show on hover */
.workflow-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 10px;
  height: 10px;
  transition: opacity 0.2s;
}

.workflow-node:hover :deep(.vue-flow__handle) {
  opacity: 1;
}
</style>
