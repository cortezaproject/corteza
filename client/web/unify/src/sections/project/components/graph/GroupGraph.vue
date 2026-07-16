<template>
  <div class="h-full flex flex-col">
    <div class="px-3 py-2 border-b border-surface flex items-center gap-2 text-sm shrink-0">
      <i :class="['pi', headerIcon, 'text-primary']" />
      <span class="font-semibold uppercase tracking-wider text-muted-color">{{ headerLabel }}</span>
      <span v-if="focusedGroup" class="text-muted-color truncate">
        · {{ focusedGroup.name || t('project.groups.untitledGroup') }}
      </span>
      <span class="text-muted-color ml-auto">
        {{ t('project.groups.nodeCount', visibleResourceCount) }}
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
import EchartsGraph from '@/sections/project/components/graph/EchartsGraph.vue'
import { allLinks } from '@/sections/project/utils/links'
import { PERM_COLORS, permissionEdges } from '@/sections/project/utils/rbac'
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

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

  // Module relationships, inferred from record-reference fields (undirected, deduped).
  for (const r of props.project.resources || []) {
    if (r.kind !== 'module' || !includeId(r.id)) continue
    for (const f of r.fields || []) {
      if (f.type !== 'Record' || !f.targetModuleId || !includeId(f.targetModuleId)) continue
      const key = [r.id, f.targetModuleId].sort().join('|')
      if (seenPair.has(key)) continue
      seenPair.add(key)
      out.push({ source: r.id, target: f.targetModuleId })
    }
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
