<template>
  <div class="flex flex-col gap-1">
    <CSidebarNavItem
      v-for="node in tree"
      :key="node[idKey]"
      :node="node"
      :id-key="idKey"
      :label-key="labelKey"
      :icon-key="iconKey"
      :icon="icon"
      :divider-key="dividerKey"
      :route-key="routeKey"
      :active-id="activeId"
      :expanded-ids="expandedIds"
      :depth="0"
      @select="onSelect"
      @toggle="onToggle"
    >
      <template v-if="$slots.item" #item="slotProps">
        <slot name="item" v-bind="slotProps" />
      </template>
    </CSidebarNavItem>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import CSidebarNavItem from './CSidebarNavItem.vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  idKey: {
    type: String,
    default: 'pageID',
  },
  parentKey: {
    type: String,
    default: 'selfID',
  },
  labelKey: {
    type: String,
    default: 'title',
  },
  iconKey: {
    type: String,
    default: undefined,
  },
  icon: {
    type: String,
    default: undefined,
  },
  weightKey: {
    type: String,
    default: 'weight',
  },
  activeId: {
    type: String,
    default: null,
  },
  rootId: {
    type: String,
    default: '0',
  },
  filterFn: {
    type: Function,
    default: undefined,
  },
  dividerKey: {
    type: String,
    default: undefined,
  },
  routeKey: {
    type: String,
    default: undefined,
  },
  expandAll: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['select'])

// Collect IDs of all items that have children
function getAllParentIds(items) {
  const childParents = new Set()
  for (const item of items) {
    const parentId = item[props.parentKey] || props.rootId
    if (parentId !== props.rootId) {
      childParents.add(parentId)
    }
  }
  return childParents
}

const expandedIds = ref(props.expandAll ? getAllParentIds(props.items) : new Set())

// Build tree from flat items
const tree = computed(() => {
  let items = props.items
  if (props.filterFn) {
    items = items.filter(props.filterFn)
  }

  const grouped = new Map()
  for (const item of items) {
    const parentId = item[props.parentKey] || props.rootId
    if (!grouped.has(parentId)) {
      grouped.set(parentId, [])
    }
    grouped.get(parentId).push(item)
  }

  function buildChildren(parentId) {
    const children = grouped.get(parentId) || []
    return children
      .sort((a, b) => (a[props.weightKey] || 0) - (b[props.weightKey] || 0))
      .map(item => ({
        ...item,
        _children: buildChildren(item[props.idKey]),
      }))
  }

  return buildChildren(props.rootId)
})

// Build a parent lookup map for auto-expand
const parentMap = computed(() => {
  const map = new Map()
  let items = props.items
  if (props.filterFn) {
    items = items.filter(props.filterFn)
  }
  for (const item of items) {
    map.set(item[props.idKey], item[props.parentKey] || props.rootId)
  }
  return map
})

// Auto-expand ancestors of a given item
function expandAncestors(itemId) {
  if (!itemId) return
  let currentId = parentMap.value.get(itemId)
  while (currentId && currentId !== props.rootId) {
    expandedIds.value.add(currentId)
    currentId = parentMap.value.get(currentId)
  }
}

// For activeId-based nav: watch activeId and expand ancestors
watch(
  () => props.activeId,
  id => expandAncestors(id),
  { immediate: true },
)

// For route-based nav: watch route and expand ancestors of the matching item
if (props.routeKey) {
  const currentRoute = useRoute()
  const router = useRouter()

  watch(
    () => currentRoute.path,
    () => {
      for (const item of props.items) {
        const itemRoute = item[props.routeKey]
        if (!itemRoute) continue
        const resolved = router.resolve(itemRoute)
        if (resolved.path === currentRoute.path) {
          expandAncestors(item[props.idKey])
          break
        }
      }
    },
    { immediate: true },
  )
}

function onSelect(item) {
  emit('select', item)
}

function onToggle(itemId) {
  if (expandedIds.value.has(itemId)) {
    expandedIds.value.delete(itemId)
  } else {
    expandedIds.value.add(itemId)
  }
}
</script>
