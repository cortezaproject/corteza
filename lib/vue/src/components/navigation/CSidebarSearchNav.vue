<template>
  <div class="flex flex-col min-h-0">
    <CInputSearch
      v-model="query"
      :placeholder="placeholder"
      size="small"
      class="my-2"
      data-testid="sidebar-nav-search"
    />

    <div class="flex-1 overflow-auto">
      <CSidebarNav
        v-if="filteredItems.length"
        :items="filteredItems"
        :id-key="idKey"
        :parent-key="parentKey"
        :label-key="labelKey"
        :icon-key="iconKey"
        :icon="icon"
        :weight-key="weightKey"
        :active-id="activeId"
        :root-id="rootId"
        :filter-fn="filterFn"
        :divider-key="dividerKey"
        :route-key="routeKey"
        :badge-key="badgeKey"
        :expand-all="expandAll || hasQuery"
        :match-type="matchType"
        @select="$emit('select', $event)"
      >
        <template v-if="$slots.item" #item="slotProps">
          <slot name="item" v-bind="slotProps" />
        </template>
        <template v-if="$slots.badge" #badge="slotProps">
          <slot name="badge" v-bind="slotProps" />
        </template>
      </CSidebarNav>

      <div
        v-else-if="hasQuery"
        class="flex items-center justify-center py-8 text-muted-color text-sm"
        data-testid="sidebar-nav-no-results"
      >
        {{ noResultsLabel }}
      </div>
    </div>
  </div>
</template>

<script setup>
// CSidebarNav with a search box pinned above it: the box stays put, the tree
// scrolls under it. The column fits whatever height it is given, so the caller
// hands it `flex-1 min-h-0` inside a full-height sidebar body.
import { computed, ref } from 'vue'
import CInputSearch from '../input/CInputSearch.vue'
import CSidebarNav from './CSidebarNav.vue'

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
  // Extra text on an item that the query matches besides its label — a handle,
  // say, which is what the resource lists search on and not always what the
  // tree shows.
  searchKey: {
    type: String,
    default: '_search',
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
  badgeKey: {
    type: String,
    default: undefined,
  },
  expandAll: {
    type: Boolean,
    default: false,
  },
  matchType: {
    type: String,
    default: 'exact',
  },
  placeholder: {
    type: String,
    default: '',
  },
  noResultsLabel: {
    type: String,
    default: '',
  },
})

defineEmits(['select'])

const query = ref('')

const hasQuery = computed(() => query.value.trim().length > 0)

// What a search leaves behind is still a tree: a match keeps its ancestors, so
// it renders where it lives, and a group that matches keeps its whole subtree,
// so searching for the group name shows what is in it.
const filteredItems = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return props.items

  const items = props.items
  const byId = new Map(items.map(i => [i[props.idKey], i]))

  const matched = new Set()
  for (const item of items) {
    const text = `${item[props.labelKey] ?? ''} ${item[props.searchKey] ?? ''}`.toLowerCase()
    if (text.includes(q)) matched.add(item[props.idKey])
  }

  const keep = new Set(matched)

  for (const id of matched) {
    const seen = new Set()
    let parentId = byId.get(id)?.[props.parentKey]
    while (parentId && parentId !== props.rootId && !seen.has(parentId)) {
      seen.add(parentId)
      keep.add(parentId)
      parentId = byId.get(parentId)?.[props.parentKey]
    }
  }

  const childrenOf = new Map()
  for (const item of items) {
    const parentId = item[props.parentKey] || props.rootId
    if (!childrenOf.has(parentId)) childrenOf.set(parentId, [])
    childrenOf.get(parentId).push(item)
  }

  const pending = [...matched]
  while (pending.length) {
    for (const child of childrenOf.get(pending.pop()) || []) {
      const childId = child[props.idKey]
      if (keep.has(childId)) continue
      keep.add(childId)
      pending.push(childId)
    }
  }

  return items.filter(i => keep.has(i[props.idKey]))
})
</script>
