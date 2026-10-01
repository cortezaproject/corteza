<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('page.navigation.page') }}</span>
  </Teleport>

  <CViewContainer class="flex flex-col">
    <Card class="flex-1 overflow-auto min-w-0 w-full" :pt="{ body: { class: 'p-0' } }">
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

          <div v-if="shownNodes.length" ref="treeEl" class="relative">
            <PageTreeBranch :nodes="shownNodes" parent-id="0" :depth="0" root class="page-tree" />
            <!-- Where the carried page will land -->
            <div
              v-if="drag?.lineStyle"
              class="page-drop-line"
              :style="drag.lineStyle"
              data-test-id="page-drop-line"
            />
          </div>

          <div v-else-if="!loading" class="flex items-center justify-center h-32 text-muted-color">
            {{ treeNodes.length ? $t('general.label.noResults') : $t('page.noPages') }}
          </div>
        </div>
      </template>
    </Card>
  </CViewContainer>

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

  <!-- Make sub-page of… -->
  <Dialog
    :visible="!!moveUnder"
    :header="moveUnder ? $t('page.list.makeSubPageOfTitle', { page: moveUnder.page.title }) : ''"
    modal
    :style="{ width: '28rem' }"
    @update:visible="v => !v && (moveUnder = null)"
  >
    <CFormGroup v-if="moveUnder" :label="$t('page.list.parentPage')">
      <Select
        v-model="moveUnder.targetID"
        :options="moveUnder.targets"
        option-label="label"
        option-value="value"
        class="w-full"
        data-test-id="page-move-target"
      />
    </CFormGroup>
    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="moveUnder = null"
      />
      <Button
        :label="$t('page.list.move')"
        icon="pi pi-check"
        size="small"
        data-test-id="page-move-confirm"
        @click="confirmMoveUnder"
      />
    </template>
  </Dialog>

  <!-- The carried page, at the pointer -->
  <Teleport to="body">
    <Transition name="page-carry">
      <div
        v-if="drag"
        class="page-carry fixed flex flex-col gap-1 border rounded-md bg-[var(--p-content-background)] pl-3 pr-3 py-2"
        :class="{ 'page-carry-landing': drag.landing }"
        :style="{
          left: `${drag.x - drag.offsetX}px`,
          top: `${drag.y - drag.offsetY}px`,
          width: `${drag.width}px`,
        }"
        data-test-id="page-carry"
      >
        <div class="flex items-center gap-2 min-w-0">
          <span class="truncate">{{ drag.node.label }}</span>
          <Tag
            v-if="drag.node.children.length"
            :value="$t('page.list.subPages', drag.node.children.length)"
            severity="secondary"
            class="text-xs leading-none py-0.5 shrink-0"
          />
        </div>
        <!-- What travels with it: the subtree, in brief -->
        <ul v-if="carriedRows.length" class="page-carry-kids">
          <li
            v-for="kid in carriedRows"
            :key="kid.key"
            class="truncate text-xs text-muted-color"
            :style="{ paddingLeft: `${kid.depth * 12}px` }"
          >
            {{ kid.label }}
          </li>
          <li v-if="carriedMore" class="text-xs text-muted-color">
            {{ $t('page.list.andMore', carriedMore) }}
          </li>
        </ul>
      </div>
    </Transition>
  </Teleport>
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
import { computed, inject, onBeforeUnmount, onMounted, provide, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import PageTreeBranch from './PageTreeBranch.vue'
import { dropPlan, projectDrop } from './pageTreeDrop'

const { CInputSearch, CPermissionsButton, CRouterLinkButton, CViewContainer } = components
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
const filterValue = ref('')
const loading = ref(false)

// A drop reorders the pages under a parent, which the backend answers with the
// namespace's page-create permission at the root level (onReorder) — without it
// no drop can land, so the tree is not draggable at all.
const canReorder = computed(() => !!props.namespace?.canCreatePage)

// A drop reorders the siblings it lands among, so a filtered tree is never
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

  if (props.namespace?.canCreatePage) {
    items.push({
      label: t('page.list.addSubPage'),
      icon: 'pi pi-plus',
      route: { name: 'admin.pages.create', query: { parent: page.pageID } },
    })
  }

  if (canReorder.value && moveTargets(page).length) {
    items.push({
      label: t('page.list.makeSubPageOf'),
      icon: 'pi pi-sitemap',
      command: () => openMoveUnder(page),
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

// "Make sub-page of…": the pages this one may go under — not itself, nothing
// beneath it, and not the parent it already has. Top level counts as a
// parent.
const moveUnder = ref(null)

function moveTargets(page) {
  const out = []
  if (page.selfID !== NoID) out.push({ value: NoID, label: t('page.edit.noParent') })
  const walk = (nodes, depth) => {
    for (const node of nodes) {
      if (node.key === page.pageID) continue
      if (node.key !== page.selfID) {
        out.push({ value: node.key, label: `${'\u2003'.repeat(depth)}${node.label}` })
      }
      walk(node.children, depth + 1)
    }
  }
  walk(treeNodes.value, 0)
  return out
}

function openMoveUnder(page) {
  const targets = moveTargets(page)
  moveUnder.value = { page, targets, targetID: targets[0]?.value ?? NoID }
}

// The page goes last among its new siblings
async function confirmMoveUnder() {
  const { page, targetID } = moveUnder.value
  moveUnder.value = null
  const namespaceID = props.namespace.namespaceID
  const siblings =
    targetID === NoID ? treeNodes.value : findNode(treeNodes.value, targetID)?.children

  try {
    page.selfID = targetID
    page.namespaceID = namespaceID
    await pageStore.update(page)
    await pageStore.reorder({
      namespaceID,
      selfID: targetID,
      pageIDs: [...(siblings ?? []).map(n => n.key), page.pageID],
    })
    $toast.toastSuccess(t('page.list.moved'))
  } catch (e) {
    console.error('Failed to move page:', e)
    $toast.toastErrorHandler(t('page.pageMoveFailed'))(e)
  }
  await reloadAll()
}

// Convert the API page tree (recursive children) to keyed nodes
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

async function reloadTree() {
  const pages = await pageStore.loadTree({
    namespaceID: props.namespace.namespaceID,
  })
  treeNodes.value = toTreeNodes(pages)
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
function onSelect(node) {
  const page = node?.data
  if (!page) return
  router.push(
    page.canUpdatePage
      ? { name: 'admin.pages.builder', params: { pageID: page.pageID } }
      : viewRoute(page),
  )
}

// One save per drop: the new parent when it changed, then the order of the
// level the page landed in. The level it left keeps its weights, in order
// with one gap. A drop where the page already is saves nothing.
async function applyDrop(node, spot) {
  const namespaceID = props.namespace.namespaceID
  const childrenOf = key =>
    (key === NoID ? treeNodes.value : (findNode(treeNodes.value, key)?.children ?? [])).map(
      n => n.key,
    )
  const plan = dropPlan(spot, node.key, childrenOf)
  const unchanged =
    node.data.selfID === plan.parentKey &&
    JSON.stringify(childrenOf(plan.parentKey)) === JSON.stringify(plan.pageIDs)
  if (unchanged) return

  try {
    if (node.data.selfID !== plan.parentKey) {
      node.data.selfID = plan.parentKey
      node.data.namespaceID = namespaceID
      await pageStore.update(node.data)
    }
    await pageStore.reorder({ namespaceID, selfID: plan.parentKey, pageIDs: plan.pageIDs })
  } catch (e) {
    console.error('Failed to move page:', e)
    $toast.toastErrorHandler(t('page.pageMoveFailed'))(e)
  }
  await reloadAll()
}

// The server's order either way: the new one after a save, and after a
// failure whatever it kept, rather than the drop it rejected. The flat list
// too, so the sidebar follows.
async function reloadAll() {
  try {
    await reloadTree()
    await pageStore.load({ namespaceID: props.namespace.namespaceID })
  } catch (e) {
    console.error('Failed to reload the page tree:', e)
    $toast.toastDanger(t('notification.page.listFailed'))
  }
}

// ─── Carrying a page ─────────────────────────────────────────────────────────
// A page is carried by its grip. The rows stay where they are; a line shows
// where the page will land, and the tree changes only on release. Pointer
// events, so a finger works the same as a mouse.
const DRAG_THRESHOLD = 5
// The child indent: the branch's 1rem margin plus 15px padding; the gap
// between rows
const INDENT = 30
const GAP = 12
const AUTOSCROLL_EDGE = 48

const treeEl = ref()
const drag = ref(null)
const carriedKey = computed(() => drag.value?.node.key ?? null)

// The carried subtree, a few rows of it, in tree order
const CARRIED_ROWS = 4
const carriedAll = computed(() => {
  const out = []
  const walk = (nodes, depth) => {
    for (const node of nodes) {
      out.push({ key: node.key, label: node.label, depth })
      walk(node.children, depth + 1)
    }
  }
  if (drag.value) walk(drag.value.node.children, 1)
  return out
})
const carriedRows = computed(() => carriedAll.value.slice(0, CARRIED_ROWS))
const carriedMore = computed(() => Math.max(0, carriedAll.value.length - CARRIED_ROWS))
const targetKey = computed(() => drag.value?.spot?.intoKey ?? null)

// A press on a grip, until the pointer has travelled far enough to be a drag
let pending = null
let scrollFrame = 0

function onGripDown(event, node) {
  if (event.pointerType === 'mouse' && event.button !== 0) return
  pending = {
    node,
    x: event.clientX,
    y: event.clientY,
    pointerId: event.pointerId,
    grip: event.currentTarget,
  }
  try {
    event.currentTarget.setPointerCapture?.(event.pointerId)
  } catch {
    // an already-released pointer; the drag still works without capture
  }
  document.addEventListener('pointermove', onPointerMove)
  document.addEventListener('pointerup', onPointerUp)
  document.addEventListener('pointercancel', stopDrag)
  document.addEventListener('keydown', onDragKey)
}

function onPointerMove(event) {
  if (pending && !drag.value) {
    if (Math.hypot(event.clientX - pending.x, event.clientY - pending.y) < DRAG_THRESHOLD) return
    startDrag(pending, event)
  }
  if (!drag.value) return
  drag.value.x = event.clientX
  drag.value.y = event.clientY
  project()
  autoscroll()
}

function startDrag(from, event) {
  const rect = from.grip.closest('.page-row').getBoundingClientRect()
  drag.value = {
    node: from.node,
    x: event.clientX,
    y: event.clientY,
    offsetX: from.x - rect.left,
    offsetY: from.y - rect.top,
    width: rect.width,
    spot: null,
    lineStyle: null,
    landing: false,
  }
  // The lib's carry state: no text selection while a page is carried
  document.body.classList.add('c-dragging')
}

// The rows on screen, top to bottom, without the carried page's own block
function visibleRows() {
  const carried = new Set()
  const collect = node => {
    carried.add(node.key)
    node.children.forEach(collect)
  }
  collect(drag.value.node)

  return [...treeEl.value.querySelectorAll('.page-row')]
    .filter(el => !carried.has(el.dataset.key))
    .map(el => {
      const r = el.getBoundingClientRect()
      return {
        key: el.dataset.key,
        parentKey: el.dataset.parent,
        depth: Number(el.dataset.depth),
        top: r.top,
        bottom: r.bottom,
      }
    })
}

// The carried card's left edge says how deep the page goes
function project() {
  const d = drag.value
  const base = treeEl.value.querySelector('.tree-root > .page-node > .page-row')
  const host = treeEl.value.getBoundingClientRect()
  const baseX = base ? base.getBoundingClientRect().left : host.left
  const spot = projectDrop(visibleRows(), d.x - d.offsetX, d.y, { indent: INDENT, baseX, gap: GAP })
  d.spot = spot
  d.lineStyle = {
    top: `${spot.lineY - host.top}px`,
    left: `${baseX + spot.depth * INDENT - host.left}px`,
    width: `${d.width}px`,
  }
}

// Near the top or bottom of the scrolling card, the tree scrolls on its own
function autoscroll() {
  const host = treeEl.value.closest('.overflow-auto')
  globalThis.cancelAnimationFrame?.(scrollFrame)
  if (!host) return
  const r = host.getBoundingClientRect()
  const y = drag.value.y
  const pull =
    y < r.top + AUTOSCROLL_EDGE
      ? y - (r.top + AUTOSCROLL_EDGE)
      : y > r.bottom - AUTOSCROLL_EDGE
        ? y - (r.bottom - AUTOSCROLL_EDGE)
        : 0
  if (!pull) return
  const step = () => {
    if (!drag.value) return
    host.scrollTop += Math.sign(pull) * Math.min(24, Math.abs(pull) / 2)
    project()
    scrollFrame = requestAnimationFrame(step)
  }
  scrollFrame = requestAnimationFrame(step)
}

// On release the card flies to where the line was and stays there while the
// page is saved and the tree reloaded with it in that place; then it fades
// under the real row
async function onPointerUp() {
  const d = drag.value
  if (!d?.spot || !d.lineStyle) return stopDrag()
  unlisten()
  const host = treeEl.value.getBoundingClientRect()
  d.landing = true
  d.x = host.left + parseFloat(d.lineStyle.left) + d.offsetX
  d.y = host.top + parseFloat(d.lineStyle.top) + GAP / 2 + d.offsetY
  d.lineStyle = null
  await applyDrop(d.node, d.spot)
  if (drag.value === d) drag.value = null
}

function onDragKey(event) {
  if (event.key === 'Escape') stopDrag()
}

function unlisten() {
  try {
    pending?.grip.releasePointerCapture?.(pending.pointerId)
  } catch {
    // never captured
  }
  pending = null
  globalThis.cancelAnimationFrame?.(scrollFrame)
  document.body.classList.remove('c-dragging')
  document.removeEventListener('pointermove', onPointerMove)
  document.removeEventListener('pointerup', onPointerUp)
  document.removeEventListener('pointercancel', stopDrag)
  document.removeEventListener('keydown', onDragKey)
}

function stopDrag() {
  unlisten()
  drag.value = null
}

onBeforeUnmount(stopDrag)

// What every level of the tree shares; the branch component reads it
provide('pageTree', {
  canDrag,
  carriedKey,
  targetKey,
  onGripDown,
  onSelect,
  actionItems,
  showActionsMenu,
  recordPageLabel,
})
</script>

<style scoped>
/* The line where the carried page will land: at the depth the card's left
 * edge asks for, with a ring at its start */
.page-drop-line {
  position: absolute;
  height: 2px;
  transition:
    top 120ms ease,
    left 120ms ease,
    width 120ms ease;
  margin-top: -1px;
  border-radius: 1px;
  background: var(--p-primary-color);
  pointer-events: none;
  z-index: 1;
}

.page-drop-line::before {
  content: '';
  position: absolute;
  left: -4px;
  top: -3px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: 2px solid var(--p-primary-color);
  background: var(--p-content-background);
}

.page-carry {
  z-index: 1100;
  pointer-events: none;
  cursor: grabbing;
  opacity: 0.95;
  box-shadow: var(--p-overlay-popover-shadow, 0 8px 24px rgb(0 0 0 / 25%));
}

.page-carry-kids {
  margin: 0;
  padding: 0.125rem 0 0 0.25rem;
  list-style: none;
  border-left: 1px solid var(--page-tree-guide, var(--p-surface-300));
}

.page-carry-enter-active {
  transition:
    opacity 120ms ease,
    transform 120ms ease;
}

.page-carry-enter-from {
  opacity: 0;
  transform: scale(0.97);
}

.page-carry-landing {
  transition:
    left 180ms ease,
    top 180ms ease;
}

.page-carry-leave-active {
  transition: opacity 120ms ease;
}

.page-carry-leave-to {
  opacity: 0;
}
</style>
