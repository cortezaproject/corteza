<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('page.navigation.page') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full flex flex-col overflow-hidden min-w-0">
    <Card class="flex-1 overflow-auto min-w-0" :pt="{ body: { class: 'p-0' } }">
      <template #header>
        <!-- Header: Create button + Search -->
        <div class="flex items-center justify-between gap-3 p-3 border-b">
          <CRouterLinkButton
            v-if="namespace?.canCreatePage"
            :to="{ name: 'admin.pages.create' }"
            :label="$t('page.createLabel')"
            icon="pi pi-plus"
            size="small"
          />
          <div v-else />

          <CInputSearch
            v-model="filterValue"
            :placeholder="$t('page.searchPlaceholder')"
            size="small"
            class="w-80"
          />
        </div>
      </template>

      <template #content>
        <Tree
          v-if="treeNodes.length"
          v-model:value="treeNodes"
          v-model:expanded-keys="expandedKeys"
          :filter="!!filterValue"
          :filter-value="filterValue"
          filter-mode="lenient"
          filter-by="label"
          draggable-nodes
          droppable-nodes
          :pt="treePT"
          selection-mode="single"
          class="p-0"
          @node-select="onNodeSelect"
          @node-drop="onNodeDrop"
        >
          <template #default="{ node }">
            <div class="flex items-center gap-3">
              <span :class="{ 'text-muted-color font-medium': node.data.selfID === '0' }">
                {{ node.label }}
              </span>
            </div>
          </template>
        </Tree>

        <div v-else-if="!loading" class="flex items-center justify-center h-32 text-muted-color">
          {{ $t('page.noPages') }}
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { compose } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { usePageStore } from '@/stores/page'

const { CInputSearch, CRouterLinkButton } = components
const { t } = useI18n()
const router = useRouter()
const $toast = inject('$toast')
const pageStore = usePageStore()

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const treeNodes = ref([])
const expandedKeys = ref({})
const filterValue = ref('')
const loading = ref(false)

const treePT = {
  rootChildren: { class: 'flex flex-col gap-3' },
  nodeChildren: { class: 'flex flex-col gap-2 py-2 ml-5 border-l' },
  wrapper: { class: 'p-2' },
  nodeContent: {
    class:
      'flex flex-row shadow border rounded-lg transition-colors hover:bg-emphasis cursor-pointer px-3 py-2',
  },
  nodeToggleButton: { class: 'order-1 ml-auto' },
  nodeLabel: { class: 'flex-1' },
}

// Convert API page tree (recursive children) to PrimeVue TreeNode format
function toTreeNodes(pages) {
  if (!pages || !Array.isArray(pages)) return []

  return pages
    .sort((a, b) => (a.weight || 0) - (b.weight || 0))
    .map(p => {
      const page = new compose.Page(p)
      return {
        key: page.pageID,
        label: page.title || page.handle || page.pageID,
        data: page,
        children: toTreeNodes(p.children),
      }
    })
}

// Collect keys of all parent nodes so they start expanded
function collectParentKeys(nodes, keys = {}) {
  for (const node of nodes) {
    if (node.children?.length) {
      keys[node.key] = true
      collectParentKeys(node.children, keys)
    }
  }
  return keys
}

// Load page tree on mount
onMounted(async () => {
  loading.value = true
  try {
    const pages = await pageStore.loadTree({
      namespaceID: props.namespace.namespaceID,
    })
    treeNodes.value = toTreeNodes(pages)
    expandedKeys.value = collectParentKeys(treeNodes.value)
  } catch (e) {
    console.error('Failed to load page tree:', e)
    $toast.toastDanger(t('notification.page.listFailed'))
  } finally {
    loading.value = false
  }
})

// Handle node click — navigate to page edit
function onNodeSelect(node) {
  if (node?.key) {
    router.push({
      name: 'admin.pages.edit',
      params: { pageID: node.key },
    })
  }
}

// Handle drag-and-drop
async function onNodeDrop(event) {
  // event.value contains the new tree state after the drop
  const newTree = event.value
  treeNodes.value = newTree

  try {
    await reorderTree(newTree, '0')

    // Refetch tree to stay in sync with server
    const pages = await pageStore.loadTree({
      namespaceID: props.namespace.namespaceID,
    })
    treeNodes.value = toTreeNodes(pages)
    expandedKeys.value = collectParentKeys(treeNodes.value)

    // Reload the flat page list so the sidebar reflects the new order
    await pageStore.load({ namespaceID: props.namespace.namespaceID, force: true })
  } catch (e) {
    console.error('Failed to reorder pages:', e)
    $toast.toastDanger(t('page.pageMoveFailed'))
  }
}

// Walk tree and persist order + reparenting for each level
async function reorderTree(nodes, parentID) {
  if (!nodes?.length) return

  const namespaceID = props.namespace.namespaceID
  const pageIDs = nodes.map(n => n.key)

  // First: update selfID on any reparented nodes (matching old Corteza approach)
  for (const node of nodes) {
    if (node.data?.selfID !== parentID) {
      node.data.selfID = parentID
      node.data.namespaceID = namespaceID
      await pageStore.update(node.data)
    }
  }

  // Then: reorder children under this parent
  await pageStore.reorder({
    namespaceID,
    selfID: parentID,
    pageIDs,
  })

  // Recurse into children
  for (const node of nodes) {
    if (node.children?.length) {
      await reorderTree(node.children, node.key)
    }
  }
}
</script>
