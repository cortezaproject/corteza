<template>
  <div class="h-full flex flex-col">
    <div class="px-3 py-2 border-b border-surface flex items-center gap-2 text-sm shrink-0">
      <i :class="['pi', headerIcon, 'text-primary']" />
      <span class="font-semibold uppercase tracking-wider text-muted-color">{{ headerLabel }}</span>
      <span v-if="focusedGroup" class="text-muted-color truncate">
        · {{ focusedGroup.name || 'Untitled group' }}
      </span>
      <span class="text-muted-color ml-auto">
        {{ visibleResourceCount }} {{ visibleResourceCount === 1 ? 'node' : 'nodes' }}
      </span>
    </div>

    <div class="flex-1 min-h-0 relative">
      <EchartsGraph
        v-if="graphNodes.length"
        :nodes="graphNodes"
        :edges="graphEdges"
        layout="force"
        :show-legend="false"
        :directed="false"
        :emphasize-kind="emphasizeKind"
        @node-click="onNodeClick"
      />
      <div
        v-else
        class="absolute inset-0 grid place-items-center text-center px-4 text-muted-color"
      >
        <div>
          <i class="pi pi-sparkles text-4xl mb-2" />
          <p class="text-sm">
            {{
              focusedGroup
                ? 'Add resources to this group to populate the preview.'
                : 'Nodes will appear here as you add groups and resources.'
            }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import EchartsGraph from '@/sections/project-poc/components/graph/EchartsGraph.vue'
import { useProjectsStore } from '@/sections/project-poc/stores/projects'
import { allLinks } from '@/sections/project-poc/utils/links'
import { PERM_COLORS, permissionEdges } from '@/sections/project-poc/utils/rbac'
import { computed, inject, ref, watch } from 'vue'

// Provided by the wizard: open the link/permission dialog, or focus a group.
const configureResource = inject('configureResource', null)
const openGroup = inject('openGroup', null)

// Resource/role nodes open the config dialog; group nodes open that group.
const onNodeClick = ({ id }) => {
  if (typeof id !== 'string') return
  if (id.startsWith('group:')) openGroup?.(id.slice('group:'.length))
  else configureResource?.(id)
}

const props = defineProps({
  project: { type: Object, default: null },
  // '__all__' / null = whole-project overview; otherwise focus a single group.
  selectedGroupId: { type: String, default: '__all__' },
  // Resource kind to visually emphasize (the current step's kind).
  emphasizeKind: { type: String, default: null },
})

const groups = computed(() => props.project?.groups || [])

// Backend-composed graph of real resources and their derived relations
// (module-field-ref, …). Refetched whenever the project's real resources
// change shape. FE-only overlays (group membership, permissions, manual
// planning links) are merged on top below.
const store = useProjectsStore()
const beGraph = ref({ nodes: [], edges: [] })
const graphSignature = computed(() =>
  JSON.stringify(
    (props.project?.resources || [])
      .filter(r => r.kind === 'module')
      .map(r => [r.id, r.sensitivity, (r.fields || []).map(f => [f.type, f.targetModuleId])]),
  ),
)
watch(
  [() => props.project?.id, graphSignature],
  ([id]) => {
    if (id) store.graph(id).then(g => (beGraph.value = g))
    else beGraph.value = { nodes: [], edges: [] }
  },
  { immediate: true },
)

const focusedGroup = computed(() => {
  const id = props.selectedGroupId
  if (!id || id === '__all__' || !props.project) return null
  return groups.value.find(g => g.id === id) || null
})

const headerIcon = computed(() => (focusedGroup.value ? 'pi-folder' : 'pi-globe'))
const headerLabel = computed(() => (focusedGroup.value ? 'Group preview' : 'Live overview'))

const inFocusResourceIds = computed(() => {
  if (!focusedGroup.value) return null
  const ids = new Set()
  for (const id of focusedGroup.value.roleIds || []) ids.add(id)
  for (const id of focusedGroup.value.resourceIds || []) ids.add(id)
  return ids
})

const visibleResources = computed(() => {
  const all = props.project?.resources || []
  if (!inFocusResourceIds.value) return all
  return all.filter(r => inFocusResourceIds.value.has(r.id))
})

const visibleResourceCount = computed(() => visibleResources.value.length)

const groupNodes = computed(() => {
  if (!props.project || focusedGroup.value) return []
  return groups.value
    .filter(g => (g.resourceIds || []).length || (g.roleIds || []).length)
    .map(g => ({
      id: `group:${g.id}`,
      name: g.name,
      kind: 'group',
      meta: { description: g.description },
    }))
})

const graphNodes = computed(() => {
  const base = visibleResources.value.map(r => ({
    id: r.id,
    name: r.name,
    kind: r.kind,
    meta: r.meta,
  }))
  return [...base, ...groupNodes.value]
})

const graphEdges = computed(() => {
  if (!props.project) return []

  const includeId = id => inFocusResourceIds.value === null || inFocusResourceIds.value.has(id)

  const out = []

  // Resource-to-resource links — undirected, deduped across both directions.
  const seenPair = new Set()
  for (const e of allLinks(props.project)) {
    if (!includeId(e.source) || !includeId(e.target)) continue
    const key = [e.source, e.target].sort().join('|')
    if (seenPair.has(key)) continue
    seenPair.add(key)
    out.push({ source: e.source, target: e.target })
  }

  // Relations the backend derives from the real resources (module-field-ref, …),
  // undirected and deduped against the manual links above.
  for (const e of beGraph.value.edges || []) {
    const s = String(e.sourceID)
    const t = String(e.targetID)
    if (!includeId(s) || !includeId(t)) continue
    const key = [s, t].sort().join('|')
    if (seenPair.has(key)) continue
    seenPair.add(key)
    out.push({ source: s, target: t })
  }

  // Group → role/resource membership edges — only in the overview.
  if (!focusedGroup.value) {
    for (const g of groups.value) {
      for (const rid of [...(g.roleIds || []), ...(g.resourceIds || [])]) {
        if (!includeId(rid)) continue
        out.push({ source: `group:${g.id}`, target: rid, color: '#e11d48', width: 1.25, dashed: true })
      }
    }
  }

  // Role → resource permission edges, colored by effective level.
  for (const e of permissionEdges(props.project, includeId)) {
    out.push({ source: e.source, target: e.target, color: PERM_COLORS[e.level].stroke, width: 1.5 })
  }

  return out
})
</script>
