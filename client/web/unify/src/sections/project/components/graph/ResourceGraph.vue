<template>
  <div class="w-full h-full flex flex-col gap-3">
    <!-- Header: what you're looking at + refresh -->
    <div class="shrink-0 flex items-center gap-2">
      <h3 class="text-xs font-semibold uppercase tracking-wide text-muted-color">{{ $t('project.graph.resources') }}</h3>
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
        <span class="text-base font-semibold leading-none">{{ m.count }}</span>
        <span class="text-xs text-muted-color leading-none">{{ $t(m.cfg.labelKey) }}</span>
      </button>
    </div>

    <!-- Relationship graph -->
    <div class="flex-1 min-h-0 rounded-lg border border-surface overflow-hidden bg-emphasis relative">
      <v-chart
        v-if="visibleNodes.length"
        :option="option"
        autoresize
        class="w-full h-full"
        @click="onClick"
      />
      <div v-else class="absolute inset-0 grid place-items-center text-center px-6 text-muted-color">
        <div>
          <i class="pi pi-sitemap text-4xl mb-2" />
          <p class="text-sm">
            {{ graph.nodes.length ? $t('project.graph.allLayersHidden') : $t('project.graph.empty') }}
          </p>
        </div>
      </div>
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
// access overlay is a single toggle covering both role and user chips; user
// nodes join once the Users step lands (until then that chip counts 0).
const accessReady = true

const store = useProjectsStore()
const $toast = inject('$toast')

// Clicking a module node opens its detail editor (meta + fields, provided by
// the wizard); other kinds get their own editors as their steps land.
const inspectResource = inject('inspectResource', null)
const onClick = params => {
  if (params.dataType !== 'node' || props.locked) return
  if (params.data.kind === 'module') inspectResource?.(params.data.id)
}

// --- Data: the backend graph is the single source of truth ------------------
const graph = ref({ nodes: [], edges: [] })
const loading = ref(false)
// A reload requested mid-fetch runs once more after it — back-to-back saves
// coalesce instead of losing the trailing refetch.
let pending = false

async function reload() {
  if (!props.project?.id) return
  if (loading.value) {
    pending = true
    return
  }
  loading.value = true
  try {
    graph.value = await store.graph(props.project.id)
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
  [() => props.project?.id, () => store.graphVersion],
  ([id]) => {
    if (id) reload()
    else graph.value = { nodes: [], edges: [] }
  },
  { immediate: true },
)

// --- Layer selection ---------------------------------------------------------
// Roles/users are only interactive once the access overlay is wired up.
const isAccessKind = kind => ACCESS_KINDS.includes(kind)

function onChipClick(m) {
  if (!m.toggleable) return
  isAccessKind(m.kind) ? store.graphToggleAccess() : store.graphToggleKind(m.kind)
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
const visibleNodes = computed(() =>
  graph.value.nodes.filter(n => store.graphKindVisible(n.kind)),
)
const visibleEdges = computed(() => {
  const ids = new Set(visibleNodes.value.map(n => n.id))
  return graph.value.edges.filter(e => ids.has(e.source) && ids.has(e.target))
})

// --- Rendering ----------------------------------------------------------------

// Human wording for the backend's edge reasons (i18n keys).
const EDGE_REASONS = {
  'module-field-ref': 'project.graph.edgeReason.moduleFieldRef',
  'role-rbac': 'project.graph.edgeReason.roleRbac',
  'user-role': 'project.graph.edgeReason.userRole',
}

const nameById = computed(() => new Map(graph.value.nodes.map(n => [n.id, n.name])))

const degreeMap = computed(() => {
  const m = new Map()
  for (const e of visibleEdges.value) {
    m.set(e.source, (m.get(e.source) || 0) + 1)
    m.set(e.target, (m.get(e.target) || 0) + 1)
  }
  return m
})

const option = computed(() => {
  const dark = document.documentElement.classList.contains('dark')
  const labelColor = dark ? '#cbd5e1' : '#334155'

  const data = visibleNodes.value.map(n => {
    const size = 34 + Math.min(20, (degreeMap.value.get(n.id) || 0) * 3)
    return {
      id: n.id,
      name: n.name,
      symbol: kindIconDataUri(n.kind),
      symbolSize: [size, size],
      symbolKeepAspect: true,
      label: { show: true, position: 'right', fontSize: 11, color: labelColor },
      // Carried for the tooltip only.
      kind: n.kind,
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
        color: '#94a3b8',
        width: 1.5,
        opacity: 0.6,
        curveness: total === 1 ? 0.16 : -0.45 + (0.9 * (i + 0.5)) / total,
      },
    }
  })

  return {
    tooltip: {
      trigger: 'item',
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
          const lines = [
            `<b>${params.data.name}</b>`,
            `${kindLabel} · ${deg} ${refLabel}`,
          ]
          return lines.join('<br/>')
        }
        if (params.dataType === 'edge') {
          const from = nameById.value.get(params.data.source) || params.data.source
          const to = nameById.value.get(params.data.target) || params.data.target
          const reasonKey = EDGE_REASONS[params.data.reason]
          const reason = reasonKey ? t(reasonKey) : params.data.reason || ''
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
