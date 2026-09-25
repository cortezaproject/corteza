<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('page.navigation.page') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full flex flex-col overflow-hidden min-w-0">
    <Card
      class="flex-1 overflow-auto min-w-0 w-full max-w-3xl mx-auto"
      :pt="{ body: { class: 'p-0' } }"
    >
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
            class="w-80"
          />
        </div>
      </template>

      <template #content>
        <div
          class="p-3 flex flex-col gap-2 [--page-tree-guide:var(--p-surface-300)] dark:[--page-tree-guide:var(--p-surface-600)]"
        >
          <p v-if="canDrag && treeNodes.length" class="text-xs text-muted-color">
            {{ $t('page.instructions') }}
          </p>

          <Tree
            v-if="shownNodes.length"
            :value="shownNodes"
            v-model:expanded-keys="expandedKeys"
            @node-collapse="expandAll"
            :draggable-nodes="canDrag"
            :droppable-nodes="canDrag"
            :pt="treePT"
            selection-mode="single"
            class="page-tree p-0"
            @node-select="onNodeSelect"
            @node-drop="onNodeDrop"
          >
            <template #default="{ node }">
              <div class="flex items-center gap-2 min-w-0" data-test-id="page-tree-node">
                <span class="whitespace-normal break-words">{{ node.label }}</span>
                <span
                  v-if="node.data.description"
                  class="text-xs text-muted-color truncate min-w-0"
                  :title="node.data.description"
                >
                  {{ node.data.description }}
                </span>

                <Tag
                  v-if="node.data.isRecordPage"
                  :value="recordPageLabel(node.data)"
                  severity="secondary"
                  class="text-xs leading-none py-0.5 shrink-0"
                />
                <Tag
                  v-else-if="!node.data.visible"
                  v-tooltip.bottom="$t('page.list.hiddenTooltip')"
                  :value="$t('page.notVisible')"
                  icon="pi pi-eye-slash"
                  severity="warn"
                  class="text-xs leading-none py-0.5 shrink-0"
                />
                <Button
                  v-if="actionItems(node.data).length"
                  v-tooltip.bottom="$t('general.label.actions')"
                  icon="pi pi-ellipsis-v"
                  text
                  severity="secondary"
                  size="small"
                  class="shrink-0 opacity-40 transition-opacity group-hover:opacity-100 focus-visible:opacity-100"
                  :aria-label="$t('general.label.actions')"
                  data-test-id="page-tree-actions"
                  @click.stop="showActionsMenu($event, node.data)"
                />
              </div>
            </template>
          </Tree>

          <div v-else-if="!loading" class="flex items-center justify-center h-32 text-muted-color">
            {{ treeNodes.length ? $t('general.label.noResults') : $t('page.noPages') }}
          </div>
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

// A drop persists the whole tree it is given, so a filtered tree is never
// draggable: it would reorder the pages the search hides.
const canDrag = computed(() => canReorder.value && !filterValue.value.trim())

// A page stays when its title matches, with everything beneath it; otherwise
// it stays only as the path down to a match.
function filterNodes(nodes, query) {
  return nodes.flatMap(node => {
    if (node.label.toLocaleLowerCase().includes(query)) return [node]
    const children = filterNodes(node.children || [], query)
    return children.length ? [{ ...node, children }] : []
  })
}

const shownNodes = computed(() => {
  const query = filterValue.value.trim().toLocaleLowerCase()
  return query ? filterNodes(treeNodes.value, query) : treeNodes.value
})

const treePT = {
  rootChildren: { class: 'flex flex-col items-start gap-3' },
  nodeChildren: { class: 'tree-branch flex flex-col items-start gap-3 pt-3 ml-3 pl-4' },
  // A row is as wide as its own content, so the hierarchy reads by shape and
  // not only by indent.
  nodeContent: {
    class:
      'group relative inline-flex flex-row items-center gap-1 w-fit max-w-full border rounded-md bg-[var(--p-content-background)] transition-colors hover:bg-emphasis cursor-pointer pl-3 pr-1 py-2',
  },
  // Every branch stays open, so there is nothing to toggle
  nodeToggleButton: { class: 'hidden' },
  nodeLabel: { class: 'min-w-0' },
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

// Collect keys of all parent nodes, which are always expanded
function collectParentKeys(nodes, keys = {}) {
  for (const node of nodes) {
    if (node.children?.length) {
      keys[node.key] = true
      collectParentKeys(node.children, keys)
    }
  }
  return keys
}

// Keyboard navigation can still collapse a node; open it again
function expandAll() {
  expandedKeys.value = collectParentKeys(treeNodes.value)
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
  } catch (e) {
    console.error('Failed to reorder pages:', e)
    $toast.toastErrorHandler(t('page.pageMoveFailed'))(e)
  }

  // The server's order either way: the new one after a save, and after a
  // failure whatever it kept, rather than the drop it rejected
  try {
    await reloadTree()

    // Reload the flat page list so the sidebar reflects the new order
    await pageStore.load({ namespaceID: props.namespace.namespaceID })
  } catch (e) {
    console.error('Failed to reload the page tree:', e)
    $toast.toastDanger(t('notification.page.listFailed'))
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

<style scoped>
/* Tree guides. A parent's trunk drops from under its title to its children,
 * and the elbow hangs off each child row so it meets the row's middle whatever
 * the row holds. Below a top-level page the trunk runs on to the next one, so
 * the top-level pages hang off one spine; it passes behind the opaque rows.
 * Keep each selector on one line: a wrapped :deep() compiles `X ::before`,
 * which blanks every icon glyph inside the row. */
.page-tree :deep(.p-tree-node) {
  position: relative;
}

.page-tree :deep(.tree-branch > .p-tree-node)::before {
  content: '';
  position: absolute;
  left: -15px;
  top: -0.75rem;
  bottom: 0;
  border-left: 1px solid var(--page-tree-guide);
}

.page-tree :deep(.tree-branch > .p-tree-node:last-child)::before {
  bottom: auto;
  height: calc(0.75rem + 24px);
}

.page-tree :deep(.p-tree-root-children > .p-tree-node:not(:last-child))::before {
  content: '';
  position: absolute;
  left: 0.75rem;
  top: 24px;
  bottom: -0.75rem;
  border-left: 1px solid var(--page-tree-guide);
}

.page-tree :deep(.tree-branch > .p-tree-node > .p-tree-node-content)::before {
  content: '';
  position: absolute;
  left: -15px;
  top: 50%;
  width: 11px;
  border-top: 1px solid var(--page-tree-guide);
}
</style>
