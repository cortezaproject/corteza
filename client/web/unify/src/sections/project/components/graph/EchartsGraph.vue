<template>
  <div class="relative w-full h-full">
    <div
      v-if="showLegend"
      class="absolute top-2 left-1/2 -translate-x-1/2 z-10 flex flex-wrap items-center justify-center gap-1 px-2"
    >
      <button
        v-for="k in kindKeys"
        :key="k"
        type="button"
        class="flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium ring-1 transition-colors select-none"
        :class="
          hiddenKinds.has(k)
            ? 'bg-emphasis text-muted-color ring-surface opacity-60'
            : 'bg-surface text-color-emphasis ring-surface hover:ring-primary'
        "
        @mouseenter="isolatedKind = k"
        @mouseleave="isolatedKind = null"
        @click="toggleKind(k)"
      >
        <span class="w-2.5 h-2.5 rounded-sm" :style="{ background: KIND_CONFIG[k].stroke }" />
        {{ KIND_CONFIG[k].label }}
      </button>
    </div>

    <v-chart ref="chartRef" :option="option" autoresize class="w-full h-full" @click="onClick" />
  </div>
</template>

<script setup>
import { KIND_CONFIG, kindConfig } from '@/sections/project/config/kinds'
import { kindIconDataUri } from '@/sections/project/utils/kindIcons'
import { GraphChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, ref } from 'vue'
import VChart from 'vue-echarts'

use([CanvasRenderer, GraphChart, TitleComponent, TooltipComponent])

const props = defineProps({
  // [{ id, name, kind, meta?, fixed?: {x,y}, badge?: string, tooltip?: string }]
  nodes: { type: Array, required: true },
  // [{ source, target, label?, color?, width?, dashed? }]
  edges: { type: Array, default: () => [] },
  // 'force' | 'none' | 'circular'
  layout: { type: String, default: 'force' },
  showLegend: { type: Boolean, default: true },
  emphasizedId: { type: String, default: null },
  directed: { type: Boolean, default: true },
  repulsion: { type: Number, default: 650 },
  edgeLength: { type: Array, default: () => [160, 280] },
  // Visually emphasize nodes of this kind (the current step's kind).
  emphasizeKind: { type: String, default: null },
})

const emit = defineEmits(['node-click'])

const chartRef = ref(null)

const kindKeys = computed(() => Object.keys(KIND_CONFIG))
const categoryIndex = computed(() => {
  const m = {}
  kindKeys.value.forEach((k, i) => (m[k] = i))
  return m
})

const isolatedKind = ref(null)
const hiddenKinds = ref(new Set())

const toggleKind = k => {
  const next = new Set(hiddenKinds.value)
  next.has(k) ? next.delete(k) : next.add(k)
  hiddenKinds.value = next
}

const degreeMap = computed(() => {
  const m = new Map()
  for (const e of props.edges) {
    m.set(e.source, (m.get(e.source) || 0) + 1)
    m.set(e.target, (m.get(e.target) || 0) + 1)
  }
  return m
})

// Chrome colours (node labels, default edge/edge-label colour) follow the
// active PrimeVue theme instead of hardcoded neutrals. Canvas rendering needs
// a literal colour value (not a live CSS variable), so these are resolved via
// getComputedStyle — same idiom as the surface-gap colour in
// CategoryTrendChart/CategoryDonutChart.
function themeColor(varName, fallback) {
  if (typeof document === 'undefined') return fallback
  const val = getComputedStyle(document.documentElement).getPropertyValue(varName).trim()
  return val || fallback
}

const option = computed(() => {
  const dimKind = isolatedKind.value
  const hidden = hiddenKinds.value
  const visibleNodes = props.nodes.filter(n => !hidden.has(n.kind))
  const visibleNodeIds = new Set(visibleNodes.map(n => n.id))
  const data = visibleNodes.map(n => {
    const cfg = kindConfig(n.kind)
    const emphasized = props.emphasizeKind && n.kind === props.emphasizeKind
    const baseSize = (emphasized ? 44 : 32) + Math.min(20, (degreeMap.value.get(n.id) || 0) * 2)
    const nodeDim = dimKind && n.kind !== dimKind
    const node = {
      id: n.id,
      name: n.name,
      symbol: kindIconDataUri(n.kind),
      symbolSize: [baseSize, baseSize],
      symbolKeepAspect: true,
      category: categoryIndex.value[n.kind] ?? 0,
      itemStyle: {
        color: cfg.stroke,
        opacity: nodeDim ? 0.12 : 1,
        borderColor: emphasized ? cfg.stroke : 'transparent',
        borderWidth: emphasized ? 2 : 0,
      },
      value: n.kind,
      label: {
        show: !nodeDim,
        position: 'right',
        fontSize: 11,
        color: themeColor('--p-text-color', '#334155'),
        formatter: n.badge ? `{badge|${n.badge}}  ${n.name}` : n.name,
        rich: {
          badge: {
            backgroundColor: cfg.stroke,
            color: '#fff',
            padding: [2, 5, 2, 5],
            borderRadius: 4,
            fontSize: 10,
            fontWeight: 'bold',
          },
        },
      },
    }
    if (n.fixed) {
      node.x = n.fixed.x
      node.y = n.fixed.y
      node.fixed = true
    }
    if (props.emphasizedId === n.id) {
      const big = baseSize + 12
      node.symbolSize = [big, big]
    }
    return node
  })

  // Assign per-edge curveness so parallel/overlapping edges fan out instead of stacking.
  const visibleEdges = props.edges.filter(
    e => visibleNodeIds.has(e.source) && visibleNodeIds.has(e.target),
  )
  const pairCounts = new Map()
  const pairIndex = new Map()
  for (const e of visibleEdges) {
    const key = [e.source, e.target].sort().join('|')
    pairCounts.set(key, (pairCounts.get(key) || 0) + 1)
  }
  const nodeKindById = new Map(visibleNodes.map(n => [n.id, n.kind]))
  const links = visibleEdges.map(e => {
    const key = [e.source, e.target].sort().join('|')
    const total = pairCounts.get(key) || 1
    const i = pairIndex.get(key) || 0
    pairIndex.set(key, i + 1)
    // Spread curveness symmetrically across [-0.45, 0.45] when multiple edges share a pair;
    // single edges get a light bow so they don't fall behind labels.
    const curveness = total === 1 ? 0.18 : -0.45 + (0.9 * (i + 0.5)) / total
    const edgeDim =
      dimKind &&
      nodeKindById.get(e.source) !== dimKind &&
      nodeKindById.get(e.target) !== dimKind
    return {
      source: e.source,
      target: e.target,
      lineStyle: {
        color: e.color || themeColor('--p-text-muted-color', '#94a3b8'),
        width: e.width || 1.5,
        opacity: edgeDim ? 0.05 : 0.55,
        curveness,
        type: e.dashed ? 'dashed' : 'solid',
      },
      label: e.label
        ? {
            show: true,
            formatter: e.label,
            fontSize: 10,
            color: e.color || themeColor('--p-text-muted-color', '#94a3b8'),
          }
        : { show: false },
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
          if (params.data.tooltip) return params.data.tooltip
          const deg = degreeMap.value.get(params.data.id) || 0
          return `<b>${params.data.name}</b><br/>${params.data.value} · ${deg} connections`
        }
        if (params.dataType === 'edge') {
          const label = params.data.label?.formatter || params.data.value || ''
          return `${params.data.source} → ${params.data.target}${label ? `<br/>${label}` : ''}`
        }
        return ''
      },
    },
    animationDuration: 600,
    animationEasingUpdate: 'cubicOut',
    series: [
      {
        type: 'graph',
        layout: props.layout,
        roam: true,
        draggable: true,
        label: { show: true },
        emphasis: { focus: 'adjacency', label: { fontWeight: 'bold' } },
        force: {
          repulsion: props.repulsion,
          edgeLength: props.edgeLength,
          gravity: 0.04,
          friction: 0.6,
          layoutAnimation: true,
        },
        circular: { rotateLabel: true },
        categories: kindKeys.value.map(k => ({
          name: KIND_CONFIG[k].label,
          itemStyle: { color: KIND_CONFIG[k].stroke },
        })),
        data,
        edges: links,
        edgeSymbol: props.directed ? ['none', 'arrow'] : ['none', 'none'],
        edgeSymbolSize: props.directed ? [0, 6] : [0, 0],
        lineStyle: { color: 'source', curveness: 0.15 },
      },
    ],
  }
})

const onClick = params => {
  if (params.dataType === 'node') {
    chartRef.value?.dispatchAction?.({ type: 'focusNodeAdjacency', dataIndex: params.dataIndex })
    emit('node-click', { id: params.data.id })
  }
}
</script>
