<template>
  <div class="w-full h-full flex flex-col gap-3">
    <!-- Header: what you're looking at + refresh -->
    <div class="shrink-0 flex items-center gap-2">
      <h3 class="text-xs font-semibold uppercase tracking-wide text-muted-color">
        {{ $t('project.graph.resources') }}
      </h3>
      <Button
        icon="pi pi-refresh"
        severity="secondary"
        text
        rounded
        size="small"
        class="ml-auto !w-7 !h-7"
        :loading="loading"
        :title="$t('project.graph.reloadTitle')"
        @click="reload"
      />
    </div>

    <!-- Metric cards — the whole-system overview, one card per resource kind,
         counted from the backend graph payload. Each card is also a toggle:
         click to show/hide that kind; hidden kinds dim but keep their count. -->
    <div class="shrink-0 flex flex-wrap gap-2">
      <button
        v-for="m in metrics"
        :key="m.kind"
        type="button"
        class="rounded-lg border border-surface px-2.5 py-1.5 flex items-center gap-2 whitespace-nowrap transition-opacity"
        :class="[
          m.visible ? 'opacity-100' : 'opacity-40',
          m.toggleable ? 'cursor-pointer hover:border-primary' : 'cursor-not-allowed',
        ]"
        :disabled="!m.toggleable"
        :title="m.toggleable ? '' : $t('project.graph.showAccessComingSoon')"
        @click="onChipClick(m)"
      >
        <span
          class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
          :class="[m.cfg.bg, m.cfg.ring]"
        >
          <i :class="[m.cfg.icon, m.cfg.text, 'text-xs']" />
        </span>
        <span class="text-base font-medium leading-none">{{ m.count }}</span>
        <span class="text-xs text-muted-color leading-none">{{ $t(m.cfg.labelKey) }}</span>
      </button>
    </div>

    <!-- Relationship graph -->
    <div
      class="flex-1 min-h-0 rounded-lg border border-surface overflow-hidden bg-emphasis relative"
    >
      <v-chart
        v-if="visibleNodes.length"
        :option="option"
        autoresize
        class="w-full h-full"
        @click="onClick"
      />
      <div
        v-else
        class="absolute inset-0 grid place-items-center text-center px-6 text-muted-color"
      >
        <div>
          <i class="pi pi-sitemap text-4xl mb-2" />
          <p class="text-sm">
            {{
              graph.nodes.length ? $t('project.graph.allLayersHidden') : $t('project.graph.empty')
            }}
          </p>
        </div>
      </div>
    </div>

    <!-- Reference problems — what the backend could NOT draw. A dangling
         reference is invisible on the canvas by nature (it is a line that
         isn't there), so the badge on the referencing node points here and
         this list says what is broken and why. Deliberately NOT filtered by
         the layer toggles: hiding a kind is a way to read the canvas, not a
         way to decide a problem doesn't count. -->
    <div v-if="issues.length" class="shrink-0 rounded-lg border border-surface overflow-hidden">
      <button
        type="button"
        class="w-full flex items-center gap-2 px-3 py-2 text-left hover:bg-emphasis transition-colors"
        :aria-expanded="issuesOpen"
        @click="issuesOpen = !issuesOpen"
      >
        <i class="pi pi-exclamation-triangle text-sm" :class="issueSummary.tone" />
        <span class="text-xs font-semibold uppercase tracking-wide">
          {{ $t('project.graph.issues.title') }}
        </span>
        <span class="text-xs text-muted-color">{{ issueSummary.text }}</span>
        <i
          class="pi ml-auto text-xs text-muted-color"
          :class="issuesOpen ? 'pi-chevron-up' : 'pi-chevron-down'"
        />
      </button>
      <ul v-if="issuesOpen" class="max-h-44 overflow-auto border-t border-surface">
        <li
          v-for="issue in issues"
          :key="issue.key"
          class="flex items-start gap-2 px-3 py-2 text-xs border-b border-surface last:border-b-0"
        >
          <i
            class="pi mt-0.5 shrink-0"
            :class="
              issue.severity === 'missing'
                ? 'pi-times-circle text-red-500'
                : 'pi-exclamation-circle text-amber-500'
            "
            style="font-size: 0.7rem"
          />
          <div class="min-w-0">
            <div class="truncate">
              <span class="font-medium">{{ issue.sourceName }}</span>
              <span class="text-muted-color">· {{ issue.reasonLabel }}</span>
            </div>
            <div class="text-muted-color">{{ issue.detail }}</div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ACCESS_KINDS, OVERVIEW_KINDS, kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { kindIconDataUri } from '@/sections/project/utils/kindIcons'
import { GraphChart } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import VChart from 'vue-echarts'

use([CanvasRenderer, GraphChart, TooltipComponent])

const { t } = useI18n()

const props = defineProps({
  project: { type: Object, default: null },
  // Dimmed while the current section awaits approval.
  locked: { type: Boolean, default: false },
})

// The backend emits project-scoped role nodes with their RBAC edges to the
// resources each role grants on (roles that grant nothing are omitted). The
// role and user chips each toggle their own visibility independently; user
// nodes join once the Users step lands (until then that chip counts 0).
const accessReady = true

const store = useProjectsStore()
const $toast = inject('$toast')

// Clicking a resource node opens its Detail dialog (editable, provided by the
// wizard). Every kind with a detail dialog participates; chart nodes stay inert
// (no editor exists). The `locked` guard still blocks all node clicks.
const inspectResource = inject('inspectResource', null)
const INSPECTABLE_KINDS = new Set([
  'module',
  'connection',
  'automation',
  'agent',
  'chatbot',
  'page',
  'role',
  'user',
])
//
// External nodes are exempt no matter their kind. They are not this project's
// resources — the wizard's dialogs edit the project's own — and for connections
// the two are not even the same kind of id: an in-project connection node is a
// CONFIGURED connection, an external one is the tenant-level DAL connection
// behind it. Handing that id to the connection dialog opens the wrong record,
// or nothing, with no way for the user to tell which happened.
const onClick = params => {
  if (params.dataType !== 'node' || props.locked) return
  const { kind, id, external } = params.data
  if (external) return
  if (INSPECTABLE_KINDS.has(kind)) inspectResource?.(kind, id)
}

// --- Data: the backend graph is the single source of truth ------------------
const emptyGraph = () => ({ nodes: [], edges: [], missing: [], warnings: [] })
const graph = ref(emptyGraph())
const loading = ref(false)
const issuesOpen = ref(true)
// A reload requested mid-fetch runs once more after it — back-to-back saves
// coalesce instead of losing the trailing refetch.
let pending = false

async function reload() {
  if (!props.project?.projectID) return
  if (loading.value) {
    pending = true
    return
  }
  loading.value = true
  try {
    graph.value = await store.graph(props.project.projectID)
  } catch (err) {
    $toast.toastErrorHandler(t('project.graph.toastLoadFailed'))(err)
  } finally {
    loading.value = false
    if (pending) {
      pending = false
      reload()
    }
  }
}

// Refetch after every persisting mutation (the store bumps graphVersion) —
// the backend re-derives the relations from the saved state.
watch(
  [() => props.project?.projectID, () => store.graphVersion],
  ([id]) => {
    if (id) reload()
    else graph.value = emptyGraph()
  },
  { immediate: true },
)

// --- Layer selection ---------------------------------------------------------
// Roles/users are only interactive once the access overlay is wired up.
const isAccessKind = kind => ACCESS_KINDS.includes(kind)

function onChipClick(m) {
  if (!m.toggleable) return
  // Every kind (including role/user) toggles independently; the store routes
  // access kinds to their own visibility set.
  store.graphToggleKind(m.kind)
}

// --- Derived metrics (straight from the payload) ----------------------------
// One card per kind across the whole system, PoC-style; kinds without
// backend-backed steps simply count 0 until they land. Each card is also a
// visibility toggle, so it carries its visible/toggleable state.
const metrics = computed(() => {
  const counts = {}
  for (const n of graph.value.nodes) counts[n.kind] = (counts[n.kind] || 0) + 1
  return OVERVIEW_KINDS.map(kind => ({
    kind,
    cfg: kindConfig(kind),
    count: counts[kind] || 0,
    visible: store.graphKindVisible(kind),
    toggleable: isAccessKind(kind) ? accessReady : true,
  }))
})

// The graph pane shows only the kinds the layer selector has enabled; edges
// stay only between visible nodes (hiding a kind hides its edges too).
const visibleNodes = computed(() => graph.value.nodes.filter(n => store.graphKindVisible(n.kind)))
const visibleEdges = computed(() => {
  const ids = new Set(visibleNodes.value.map(n => n.id))
  return graph.value.edges.filter(e => ids.has(e.source) && ids.has(e.target))
})

// --- Rendering ----------------------------------------------------------------

// Human wording for the backend's reference reasons (i18n keys). The same
// vocabulary labels an edge, a missing reference and a warning — they are all
// the same fact ("this is why A points at B"), only one of them made it onto
// the canvas. Covers every reason server/pkg/resourceref can attach to a
// relation the graph keeps; anything else falls back to its raw string.
const EDGE_REASONS = {
  'module-field-ref': 'project.graph.edgeReason.moduleFieldRef',
  'module-connection': 'project.graph.edgeReason.moduleConnection',
  'page-module': 'project.graph.edgeReason.pageModule',
  'page-chart': 'project.graph.edgeReason.pageChart',
  'page-automation': 'project.graph.edgeReason.pageAutomation',
  'page-agent': 'project.graph.edgeReason.pageAgent',
  'page-chatbot': 'project.graph.edgeReason.pageChatbot',
  'page-navigation': 'project.graph.edgeReason.pageNavigation',
  'chart-module': 'project.graph.edgeReason.chartModule',
  'trigger-module': 'project.graph.edgeReason.triggerModule',
  'agent-module': 'project.graph.edgeReason.agentModule',
  'agent-automation': 'project.graph.edgeReason.agentAutomation',
  'chatbot-agent': 'project.graph.edgeReason.chatbotAgent',
  'chatbot-automation': 'project.graph.edgeReason.chatbotAutomation',
  'step-argument': 'project.graph.edgeReason.stepArgument',
  'step-connection': 'project.graph.edgeReason.stepConnection',
  'role-rbac': 'project.graph.edgeReason.roleRbac',
  'user-role': 'project.graph.edgeReason.userRole',
}

const reasonLabel = reason => (EDGE_REASONS[reason] ? t(EDGE_REASONS[reason]) : reason || '')

const nameById = computed(() => new Map(graph.value.nodes.map(n => [n.id, n.name])))

// --- Reference problems -------------------------------------------------------
// missing: the reference names a resource that is not in the project and could
// not be loaded from anywhere else — a deleted module, a binding a branch copy
// left pointing at the revision it was copied from.
// warning: the reference resolves only when the thing runs (a computed step
// argument), so there is nothing to check now and nothing to draw.
//
// Both are keyed by the REFERENCING node, which is the one the user can act on.
const targetKindLabel = kind =>
  kind ? t(kindConfig(kind).singularKey) : t('project.graph.issues.someResource')

const issues = computed(() => {
  const named = id => nameById.value.get(id) || t('project.graph.issues.unknownSource')

  const missing = (graph.value.missing || []).map((m, i) => ({
    key: `m${i}`,
    severity: 'missing',
    sourceID: m.sourceID,
    sourceName: named(m.sourceID),
    reasonLabel: reasonLabel(m.reason),
    // The identifier is what makes the row actionable — it is the only trace
    // left of the resource that went away, so it is shown when there is one.
    detail: t(
      m.targetIdent || m.targetID
        ? 'project.graph.issues.missingWithTarget'
        : 'project.graph.issues.missing',
      { kind: targetKindLabel(m.kind), target: m.targetIdent || m.targetID },
    ),
  }))

  const warnings = (graph.value.warnings || []).map((w, i) => ({
    key: `w${i}`,
    severity: 'warning',
    sourceID: w.sourceID,
    sourceName: named(w.sourceID),
    reasonLabel: reasonLabel(w.reason),
    detail:
      t('project.graph.issues.warning', { kind: targetKindLabel(w.kind) }) +
      (w.path ? ` ${t('project.graph.issues.atPath', { path: w.path })}` : ''),
  }))

  return [...missing, ...warnings]
})

// Node id → worst severity on it, so the canvas badge matches the list's icon.
const issueSeverityByNode = computed(() => {
  const m = new Map()
  for (const issue of issues.value) {
    if (issue.severity === 'missing' || !m.has(issue.sourceID))
      m.set(issue.sourceID, issue.severity)
  }
  return m
})

const issuesByNode = computed(() => {
  const m = new Map()
  for (const issue of issues.value) {
    if (!m.has(issue.sourceID)) m.set(issue.sourceID, [])
    m.get(issue.sourceID).push(issue)
  }
  return m
})

const issueSummary = computed(() => {
  const broken = issues.value.filter(i => i.severity === 'missing').length
  const deferred = issues.value.length - broken
  const parts = []
  if (broken) parts.push(t('project.graph.issues.summaryBroken', { count: broken }))
  if (deferred) parts.push(t('project.graph.issues.summaryDeferred', { count: deferred }))
  return { text: parts.join(' · '), tone: broken ? 'text-red-500' : 'text-amber-500' }
})

const degreeMap = computed(() => {
  const m = new Map()
  for (const e of visibleEdges.value) {
    m.set(e.source, (m.get(e.source) || 0) + 1)
    m.set(e.target, (m.get(e.target) || 0) + 1)
  }
  return m
})

// Edge lines use the muted-text token so they follow the active theme instead
// of a hardcoded neutral. Canvas rendering needs a literal colour (not a live
// CSS variable), so it's resolved via getComputedStyle — same idiom as the
// surface-gap colour in CategoryTrendChart/CategoryDonutChart.
function mutedEdgeColor() {
  if (typeof document === 'undefined') return '#94a3b8'
  const val = getComputedStyle(document.documentElement)
    .getPropertyValue('--p-text-muted-color')
    .trim()
  return val || '#94a3b8'
}

const option = computed(() => {
  const dark = document.documentElement.classList.contains('dark')
  const labelColor = dark ? '#cbd5e1' : '#334155'

  const data = visibleNodes.value.map(n => {
    const size = 34 + Math.min(20, (degreeMap.value.get(n.id) || 0) * 3)
    const badge = issueSeverityByNode.value.get(n.id) || ''
    return {
      id: n.id,
      name: n.name,
      symbol: kindIconDataUri(n.kind, { external: !!n.external, badge }),
      symbolSize: [size, size],
      symbolKeepAspect: true,
      label: {
        show: true,
        position: 'right',
        fontSize: 11,
        // An external node's own name is not this project's to change, so its
        // label is stated more quietly than the resources the wizard owns.
        color: n.external ? mutedEdgeColor() : labelColor,
        fontStyle: n.external ? 'italic' : 'normal',
      },
      // Nothing opens for an external node (see onClick), so it does not offer
      // the pointer that promises something will.
      cursor: !props.locked && !n.external && INSPECTABLE_KINDS.has(n.kind) ? 'pointer' : 'default',
      // Carried for the tooltip and the click handler.
      kind: n.kind,
      external: !!n.external,
    }
  })

  // Fan out parallel edges so multiple references between the same pair stay visible.
  const pairCounts = new Map()
  const pairIndex = new Map()
  for (const e of visibleEdges.value) {
    const key = [e.source, e.target].sort().join('|')
    pairCounts.set(key, (pairCounts.get(key) || 0) + 1)
  }
  const links = visibleEdges.value.map(e => {
    const key = [e.source, e.target].sort().join('|')
    const total = pairCounts.get(key) || 1
    const i = pairIndex.get(key) || 0
    pairIndex.set(key, i + 1)
    return {
      source: e.source,
      target: e.target,
      reason: e.reason,
      lineStyle: {
        color: mutedEdgeColor(),
        width: 1.5,
        opacity: 0.6,
        curveness: total === 1 ? 0.16 : -0.45 + (0.9 * (i + 0.5)) / total,
      },
    }
  })

  return {
    tooltip: {
      trigger: 'item',
      // Mount on <body> so the tooltip escapes the graph pane's own
      // overflow-hidden (same fix as the dashboard charts, 2026-07-28).
      appendTo: 'body',
      enterable: false,
      backgroundColor: 'transparent',
      borderColor: 'transparent',
      padding: 0,
      extraCssText:
        'background: var(--p-content-background) !important;' +
        'color: var(--p-content-color) !important;' +
        'border: 1px solid var(--p-content-border-color) !important;' +
        'border-radius: 6px;' +
        'padding: 8px 10px;' +
        'max-width: 320px;' +
        'box-shadow: 0 10px 24px rgba(15,23,42,0.15);',
      formatter: params => {
        if (params.dataType === 'node') {
          const deg = degreeMap.value.get(params.data.id) || 0
          const kindLabel = t(kindConfig(params.data.kind).singularKey)
          const refLabel = deg === 1 ? t('project.graph.reference') : t('project.graph.references')
          const lines = [`<b>${params.data.name}</b>`, `${kindLabel} · ${deg} ${refLabel}`]
          if (params.data.external) lines.push(`<i>${t('project.graph.externalHint')}</i>`)
          // The problems on this node, spelled out where the eye already is —
          // the badge only says that there is one.
          for (const issue of issuesByNode.value.get(params.data.id) || []) {
            const color = issue.severity === 'missing' ? '#dc2626' : '#f59e0b'
            lines.push(`<span style="color:${color}">${issue.reasonLabel}: ${issue.detail}</span>`)
          }
          return lines.join('<br/>')
        }
        if (params.dataType === 'edge') {
          const from = nameById.value.get(params.data.source) || params.data.source
          const to = nameById.value.get(params.data.target) || params.data.target
          const reason = reasonLabel(params.data.reason)
          return `<b>${from}</b> → <b>${to}</b>${reason ? `<br/>${reason}` : ''}`
        }
        return ''
      },
    },
    animationDuration: 600,
    animationEasingUpdate: 'cubicOut',
    series: [
      {
        type: 'graph',
        layout: 'force',
        roam: true,
        draggable: true,
        emphasis: { focus: 'adjacency', label: { fontWeight: 'bold' } },
        force: {
          repulsion: 650,
          edgeLength: [140, 240],
          gravity: 0.05,
          friction: 0.6,
          layoutAnimation: true,
        },
        data,
        edges: links,
        // Record references are directed: referencing module → referenced module.
        edgeSymbol: ['none', 'arrow'],
        edgeSymbolSize: [0, 7],
      },
    ],
  }
})
</script>
