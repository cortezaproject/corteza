<template>
  <div
    class="workflow-node"
    :class="{
      'workflow-node--selected': selected,
      'workflow-node--highlighted': data?.highlighted,
      'workflow-node--trace-success': data?.traceState === 'success',
      'workflow-node--trace-error': data?.traceState === 'error',
      'workflow-node--hoverable': outCount === 0,
      'workflow-node--connecting': isConnecting,
    }"
    :style="{ width: '180px' }"
  >
    <!-- Target handles. `:connectable="false"` once the point is used so vue-flow
         strips the `.connectable` class → CSS hover rule leaves it invisible
         AND base styles set pointer-events: none. DOM stays so the edge
         attached to it can still anchor. -->
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
      :position="Position.Bottom"
      id="target-bottom"
      :style="{ left: '50%' }"
      :connectable="!isTargetUsed('target-bottom')"
    />

    <!-- Header (icon + title + id) -->
    <div class="workflow-node__header">
      <img v-if="iconSrc" :src="iconSrc" class="workflow-node__icon" />
      <span class="workflow-node__title" :title="displayTitle">
        {{ displayTitle }}
      </span>

      <div class="workflow-node__header-actions">
        <span class="workflow-node__id">{{ id }}</span>
      </div>
    </div>

    <!-- Issue badge (top-right) -->
    <div
      v-if="hasIssues"
      v-tooltip.top="{ value: issueTooltip, pt: { text: 'whitespace-pre-wrap text-xs' } }"
      class="workflow-node__issue-badge"
      :aria-label="$t('editor.issues')"
      @click.stop="$emit('open-issues', id)"
    >
      <img :src="issueIcon" class="workflow-node__issue-icon" alt="" />
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
    <div
      v-if="data?.traceLog"
      v-tooltip.top="{ value: data.traceLog, pt: { text: 'whitespace-pre-wrap font-mono text-xs' } }"
      class="workflow-node__trace-badge"
    >
      <img
        :src="getIcon(data.traceState === 'error' ? 'clock-danger' : 'clock-success')"
        class="workflow-node__trace-icon"
      />
    </div>

    <!-- Description (bottom) — falls back to the sidebar's step-type label so
         the second row is never empty. User-edited description overrides. -->
    <div v-if="displayDescription" class="workflow-node__description">
      {{ displayDescription }}
    </div>

    <!-- Source handles. `:connectable="false"` once the point is used or the
         node has hit its outbound cap. Vue-flow strips `.connectable`, CSS
         hides the dot on hover and disables pointer events. -->
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom"
      :style="{ left: '50%' }"
      :connectable="!isSourceUsed('source-bottom')"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-left"
      :style="{ left: '25%' }"
      :connectable="!isSourceUsed('source-bottom-left')"
    />
    <Handle
      type="source"
      :position="Position.Bottom"
      id="source-bottom-right"
      :style="{ left: '75%' }"
      :connectable="!isSourceUsed('source-bottom-right')"
    />
    <Handle
      type="source"
      :position="Position.Right"
      id="source-right"
      :connectable="!isSourceUsed('source-right')"
    />
    <Handle
      type="source"
      :position="Position.Left"
      id="source-left"
      :connectable="!isSourceUsed('source-left')"
    />
  </div>
</template>

<script setup>
import { Handle, Position, useVueFlow } from '@vue-flow/core'
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
  usedSourceHandles: { type: Array, default: () => [] },
  usedTargetHandles: { type: Array, default: () => [] },
  outCount: { type: Number, default: 0 },
})

// Outbound cap per step kind (mirrors isValidConnection in WorkflowEditor):
// fork/excl/incl gateways accept unlimited out, iterator & error-handler take
// 2 (Body+End / Try+Catch), everything else is 1. When the cap is reached,
// hide every source handle so hover doesn't even offer a starting point.
const maxOutbound = computed(() => {
  const kind = props.data?.kind
  const ref = props.data?.ref
  if (kind === 'gateway' && ['fork', 'excl', 'incl'].includes(ref)) return Infinity
  if (kind === 'iterator' || kind === 'error-handler') return 2
  return 1
})
const outboundFull = computed(() => props.outCount >= maxOutbound.value)

const isSourceUsed = id => outboundFull.value || props.usedSourceHandles.includes(id)
const isTargetUsed = id => props.usedTargetHandles.includes(id)

const { connectionStartHandle } = useVueFlow()
const isConnecting = computed(() => !!connectionStartHandle.value)

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

const stepTitleLabel = computed(() => {
  const style = getStyleFromKind(props.data)?.style
  if (!style) return stepTypeLabel.value
  const key = `steps.${style}.title`
  const resolved = t(key)
  return resolved === key ? stepTypeLabel.value : resolved
})

const displayTitle = computed(() => props.data?.label || stepTitleLabel.value)

const displayDescription = computed(() => props.data?.description || stepTypeLabel.value)

const nodeIssues = computed(() => {
  const list = props.issues?.[props.id]
  return Array.isArray(list) ? list : []
})

const hasIssues = computed(() => nodeIssues.value.length > 0)

const issueTooltip = computed(() => nodeIssues.value.join('\n'))

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
  height: 32px;
  color: var(--p-primary-color);
  font-weight: 500;
  font-size: 13px;
}

.workflow-node__icon {
  width: 18px;
  height: 18px;
  margin-right: 6px;
  object-fit: contain;
}

.workflow-node__title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--p-text-color);
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

.workflow-node__issue-badge {
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

.workflow-node__issue-icon {
  width: 14px;
  height: 14px;
  display: block;
}

.workflow-node__description {
  border-top: 1px solid var(--p-surface-border);
  padding: 6px 8px;
  font-size: 12px;
  line-height: 16px;
  color: var(--p-text-muted-color);
  white-space: pre-wrap;
  word-break: break-word;
}

.workflow-node__values {
  display: none;
  position: absolute;
  top: calc(100% + 14px);
  left: 0;
  width: 180px;
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

/* Hide handles by default, show on hover. z-index:-1 tucks the circle behind
   the node so only the outer half pokes past the border. */
.workflow-node :deep(.vue-flow__handle) {
  opacity: 0;
  width: 12px;
  height: 12px;
  z-index: -1;
  transition: opacity 0.2s;
}

/* Show handles only when this node can still take a new outbound edge
   (hover), or while a connection drag is in progress (drop target). */
.workflow-node--hoverable:hover :deep(.vue-flow__handle.connectable),
.workflow-node--connecting :deep(.vue-flow__handle.connectable) {
  opacity: 1;
}
</style>
