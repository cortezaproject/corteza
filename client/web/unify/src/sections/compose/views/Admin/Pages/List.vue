<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('page.navigation.page') }}</span>
  </Teleport>

  <div class="max-w-4xl w-full mx-auto p-4 h-full flex flex-col overflow-hidden min-w-0">
    <Card class="flex-1 overflow-auto min-w-0" :pt="{ body: { class: 'p-0' } }">
      <template #header>
        <!-- Header: Create + permissions on the left, search on the right -->
        <div class="flex items-center justify-between gap-3 p-3 border-b">
          <div class="flex items-center gap-2">
            <CRouterLinkButton
              v-if="namespace?.canCreatePage"
              :to="{ name: 'admin.pages.create' }"
              :label="$t('page.createLabel')"
              icon="pi pi-plus"
              size="small"
            />
            <CPermissionsButton
              v-if="namespace?.canGrant"
              v-tooltip.bottom="$t('general.label.permissions')"
              :resource="`corteza::compose:page/${namespace.namespaceID}/*`"
            />
          </div>

          <CInputSearch
            v-model="filterValue"
            :placeholder="$t('page.searchPlaceholder')"
            size="small"
            class="w-64"
          />
        </div>
      </template>

      <template #content>
        <p
          v-if="canReorder && treeNodes.length"
          class="flex items-center gap-2 px-4 pt-3 text-xs text-muted-color"
        >
          <span class="pi pi-arrows-v" />
          {{ $t('page.instructions') }}
        </p>

        <Tree
          v-if="treeNodes.length"
          v-model:value="treeNodes"
          v-model:expanded-keys="expandedKeys"
          :filter="!!filterValue"
          :filter-value="filterValue"
          filter-mode="lenient"
          filter-by="label"
          :draggable-nodes="canReorder"
          :droppable-nodes="canReorder"
          :pt="treePT"
          selection-mode="single"
          class="p-3"
          @node-select="onNodeSelect"
          @node-drop="onNodeDrop"
        >
          <template #default="{ node }">
            <div class="flex items-center gap-2 min-w-0" data-test-id="page-tree-node">
              <span
                :class="node.data.isRecordPage ? 'pi pi-id-card' : 'pi pi-file'"
                class="text-muted-color text-sm shrink-0"
              />
              <span
                class="shrink-0 whitespace-normal break-words"
                :class="{ 'font-semibold': node.data.selfID === '0' }"
              >
                {{ node.label }}
              </span>
              <span
                v-if="node.data.description"
                class="text-xs text-muted-color truncate min-w-0"
                :title="node.data.description"
              >
                {{ node.data.description }}
              </span>

              <div class="ml-auto flex items-center gap-1.5 shrink-0">
                <Tag
                  v-if="node.data.isRecordPage"
                  :value="recordPageLabel(node.data)"
                  severity="secondary"
                  class="text-xs"
                />
                <Tag
                  v-else-if="!node.data.visible"
                  v-tooltip.bottom="$t('page.list.hiddenTooltip')"
                  :value="$t('page.notVisible')"
                  icon="pi pi-eye-slash"
                  severity="warn"
                  class="text-xs"
                />
                <Button
                  v-if="actionItems(node.data).length"
                  v-tooltip.bottom="$t('general.label.actions')"
                  icon="pi pi-ellipsis-v"
                  text
                  rounded
                  severity="secondary"
                  size="small"
                  :aria-label="$t('general.label.actions')"
                  data-test-id="page-tree-actions"
                  @click.stop="showActionsMenu($event, node.data)"
                />
              </div>
            </div>
          </template>
        </Tree>

        <div v-else-if="!loading" class="flex items-center justify-center h-32 text-muted-color">
          {{ $t('page.noPages') }}
        </div>
      </template>
    </Card>
  </div>

  <!-- One popup shared by every node's actions button -->
  <Menu ref="actionsMenuRef" :model="currentMenuItems" popup>
    <template #item="{ item, props: menuProps }">
      <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
        <a v-ripple :href="href" v-bind="menuProps.action" @click="navigate">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </router-link>
      <a v-else v-ripple v-bind="menuProps.action" :class="item.class">
        <span :class="item.icon" />
        <span class="ml-2">{{ item.label }}</span>
      </a>
    </template>
  </Menu>
</template>

<script setup>
import { compose, NoID } from '@planetcrust/human-js'
import {
  components,
  useConfirmDelete,
  useModuleStore,
  usePageStore,
  usePermissions,
} from '@planetcrust/human-vue'
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CInputSearch, CPermissionsButton, CRouterLinkButton } = components
const { t } = useI18n()
const router = useRouter()
const $toast = inject('$toast')
const pageStore = usePageStore()
const moduleStore = useModuleStore()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()

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

// A drop reorders the pages under a parent, which the backend answers with the
// namespace's page-create permission at the root level (onReorder) — without it
// no drop can land, so the tree is not draggable at all.
const canReorder = computed(() => !!props.namespace?.canCreatePage)

const treePT = {
  rootChildren: { class: 'flex flex-col gap-1.5' },
  nodeChildren: { class: 'flex flex-col gap-1.5 pt-1.5 ml-3 pl-4 border-l border-surface' },
  nodeContent: {
    class:
      'flex flex-row items-center gap-1 border rounded-md transition-colors hover:bg-emphasis cursor-pointer px-1.5 py-1',
  },
  nodeToggleButton: { class: 'shrink-0' },
  nodeLabel: { class: 'flex-1 min-w-0' },
}

const actionsMenuRef = ref()
const currentMenuItems = ref([])

function recordPageLabel(page) {
  const module = moduleStore.getByID(page.moduleID)
  return module?.name
    ? t('page.list.recordPageOf', { module: module.name })
    : t('page.list.recordPage')
}

function hasChildren(page) {
  return !!findNode(treeNodes.value, page.pageID)?.children?.length
}

function findNode(nodes, key) {
  for (const node of nodes) {
    if (node.key === key) return node
    const found = findNode(node.children || [], key)
    if (found) return found
  }
  return undefined
}

function viewRoute(page) {
  return page.isRecordPage
    ? { name: 'page.record', params: { pageID: page.pageID, recordID: NoID } }
    : { name: 'page', params: { pageID: page.pageID } }
}

function actionItems(page) {
  const items = []
  const params = { pageID: page.pageID }

  if (page.canUpdatePage) {
    items.push(
      {
        label: t('general.label.pageBuilder'),
        icon: 'pi pi-wrench',
        route: { name: 'admin.pages.builder', params },
      },
      {
        label: t('page.list.settings'),
        icon: 'pi pi-cog',
        route: { name: 'admin.pages.edit', params },
      },
    )
  }

  items.push({ label: t('page.view'), icon: 'pi pi-eye', route: viewRoute(page) })

  if (props.namespace?.canCreatePage && !page.isRecordPage) {
    items.push({
      label: t('page.list.addSubPage'),
      icon: 'pi pi-plus',
      route: { name: 'admin.pages.create', query: { parent: page.pageID } },
    })
  }

  if (props.namespace?.canGrant || page.canGrant) {
    const namespaceID = props.namespace.namespaceID
    const target = page.title || page.handle || page.pageID
    items.push({ separator: true })

    if (props.namespace?.canGrant) {
      items.push({
        label: t('page.list.pagePermissions'),
        icon: 'pi pi-lock',
        command: () =>
          openPermissions({
            title: target,
            target,
            resource: `corteza::compose:page/${namespaceID}/${page.pageID}`,
          }),
      })
    }

    if (page.canGrant) {
      items.push({
        label: t('page.list.layoutPermissions'),
        icon: 'pi pi-lock',
        command: () =>
          openPermissions({
            title: target,
            target,
            resource: `corteza::compose:page-layout/${namespaceID}/${page.pageID}/*`,
            allSpecific: true,
          }),
      })
    }
  }

  if (page.canDeletePage) {
    items.push({ separator: true })

    if (hasChildren(page)) {
      items.push(
        {
          label: t('page.list.deleteKeepSubPages'),
          icon: 'pi pi-trash',
          class: 'text-red-500',
          command: () => onConfirmDelete(page, 'rebase'),
        },
        {
          label: t('page.list.deleteWithSubPages'),
          icon: 'pi pi-trash',
          class: 'text-red-500',
          command: () => onConfirmDelete(page, 'cascade'),
        },
      )
    } else {
      items.push({
        label: t('general.label.delete'),
        icon: 'pi pi-trash',
        class: 'text-red-500',
        command: () => onConfirmDelete(page, 'abort'),
      })
    }
  }

  return items
}

function showActionsMenu(event, page) {
  currentMenuItems.value = actionItems(page)
  actionsMenuRef.value?.show(event, event.currentTarget)
}

function onConfirmDelete(page, strategy) {
  confirmDelete({
    message: t(
      strategy === 'cascade' ? 'page.list.deleteWithSubPagesConfirm' : 'page.edit.deleteConfirm',
    ),
    header: page.title || page.handle,
    onConfirm: () => handleDelete(page, strategy),
  })
}

async function handleDelete(page, strategy) {
  try {
    await pageStore.delete({
      namespaceID: props.namespace.namespaceID,
      pageID: page.pageID,
      strategy,
    })
    $toast.toastSuccess(t('notification.page.deleted'))
  } catch (e) {
    console.error('Failed to delete page:', e)
    $toast.toastErrorHandler(t('notification.page.deleteFailed'))(e)
  }
  await reloadTree()
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

async function reloadTree() {
  const pages = await pageStore.loadTree({
    namespaceID: props.namespace.namespaceID,
  })
  treeNodes.value = toTreeNodes(pages)
  expandedKeys.value = collectParentKeys(treeNodes.value)
}

onMounted(async () => {
  loading.value = true
  try {
    await reloadTree()
  } catch (e) {
    console.error('Failed to load page tree:', e)
    $toast.toastDanger(t('notification.page.listFailed'))
  } finally {
    loading.value = false
  }
})

// A click opens the builder, or the page itself for whoever may not change it
function onNodeSelect(node) {
  const page = node?.data
  if (!page) return
  router.push(
    page.canUpdatePage
      ? { name: 'admin.pages.builder', params: { pageID: page.pageID } }
      : viewRoute(page),
  )
}

// Handle drag-and-drop
async function onNodeDrop(event) {
  // event.value contains the new tree state after the drop
  const newTree = event.value
  treeNodes.value = newTree

  try {
    await reorderTree(newTree, '0')

    // Refetch tree to stay in sync with server
    await reloadTree()

    // Reload the flat page list so the sidebar reflects the new order
    await pageStore.load({ namespaceID: props.namespace.namespaceID })
  } catch (e) {
    console.error('Failed to reorder pages:', e)
    $toast.toastErrorHandler(t('page.pageMoveFailed'))(e)
  }
}

// Walk tree and persist order + reparenting for each level
async function reorderTree(nodes, parentID) {
  if (!nodes?.length) return

  const namespaceID = props.namespace.namespaceID
  const pageIDs = nodes.map(n => n.key)

  // First: update selfID on any reparented nodes (matching old Human approach)
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
