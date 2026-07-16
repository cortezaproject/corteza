<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div v-if="!disabled">
        <Button
          icon="pi pi-plus"
          :label="$t('project.pages.add')"
          size="small"
          @click="createResource?.('page')"
        />
      </div>

      <!-- Pages render as their nav hierarchy (parent/child via selfID). Drag a
           row by its grip: drop on a row's middle to nest it inside, or near a
           row's top/bottom edge to reorder siblings. Each drop persists
           weight/selfID, which drives the compose namespace sidebar. Standalone
           and module-detail pages share one tree. -->
      <div class="flex flex-col gap-2">
        <CEmptyState v-if="!visibleRows.length">
          {{ $t('project.pages.empty') }}
        </CEmptyState>

        <div
          v-for="item in visibleRows"
          :key="item.id"
          class="group flex items-center gap-2 p-3 border bg-surface rounded-border shadow-sm cursor-pointer hover:bg-emphasis transition-colors"
          :class="[
            dragId === item.id ? 'opacity-40' : '',
            dropTarget.id === item.id && dropTarget.pos === 'inside'
              ? 'border-primary ring-1 ring-primary'
              : 'border-surface',
            dropTarget.id === item.id && dropTarget.pos === 'before'
              ? 'border-t-2 !border-t-primary'
              : '',
            dropTarget.id === item.id && dropTarget.pos === 'after'
              ? 'border-b-2 !border-b-primary'
              : '',
          ]"
          :style="{ marginLeft: item.depth * 2 + 'rem' }"
          :draggable="!disabled"
          @dragstart="onDragStart($event, item)"
          @dragover="onDragOver($event, item)"
          @dragleave="onDragLeave($event, item)"
          @drop="onDrop($event, item)"
          @dragend="onDragEnd"
        >
          <!-- Chevron only for pages that have children — no placeholder for leaves. -->
          <button
            v-if="item.hasChildren"
            type="button"
            class="flex items-center shrink-0 -my-3 -ms-3 self-stretch px-2 cursor-pointer text-muted-color hover:text-color transition-colors"
            :aria-label="$t(isOpen(item.id) ? 'project.pages.collapse' : 'project.pages.expand')"
            :title="$t(isOpen(item.id) ? 'project.pages.collapse' : 'project.pages.expand')"
            @click.stop="toggle(item.id)"
          >
            <i
              class="pi pi-chevron-down text-xs transition-transform duration-200"
              :class="{ '-rotate-90': !isOpen(item.id) }"
            />
          </button>

          <!-- Name + description; record (module detail) pages show a module kind
               badge next to the name. min-h-10 keeps every card the same height as
               the other steps whether or not a description is present. Click opens
               the page detail dialog. -->
          <div
            class="flex-1 min-w-0 min-h-10 flex flex-col justify-center"
            @click="inspectResource?.('page', item.id)"
          >
            <div class="flex items-center gap-2 min-w-0">
              <span class="truncate font-medium">{{ item.name }}</span>
              <KindBadge
                v-if="item.isRecordPage && moduleFor(item)"
                kind="module"
                :label="moduleFor(item).name"
                muted
              />
            </div>
            <div v-if="item.description" class="text-xs text-muted-color truncate">
              {{ item.description }}
            </div>
          </div>

          <!-- Right side: the status Tag swaps for the hover actions (open in
               builder, delete) — the actions overlay the Tag's slot so the row
               width stays put. -->
          <div class="relative flex items-center shrink-0">
            <div class="transition-opacity group-hover:opacity-0 group-focus-within:opacity-0">
              <Tag
                :value="item.visible ? $t('project.pages.visible') : $t('project.pages.hidden')"
                :severity="item.visible ? 'success' : 'secondary'"
                class="!text-xs"
              />
            </div>
            <div
              class="absolute inset-y-0 end-0 flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100"
            >
              <CRouterLinkButton
                :to="{
                  name: 'admin.pages.builder',
                  params: { slug: project.namespaceID, pageID: item.id },
                }"
                icon="pi pi-external-link"
                severity="secondary"
                text
                size="small"
                :aria-label="$t('project.pages.openBuilder')"
                :title="$t('project.pages.openBuilder')"
                @click.stop
              />
              <!-- Detail pages are tied to their module; only standalone pages can
                   be removed from here. -->
              <Button
                v-if="!disabled && !item.isRecordPage"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                :aria-label="$t('project.pages.remove')"
                :title="$t('project.pages.remove')"
                @click.stop="onRemove(item)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import KindBadge from '@/sections/project/components/KindBadge.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CEmptyState, CRouterLinkButton } = components

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// Detail/create dialogs are mounted once in the wizard; open them via injection.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

// --- Page hierarchy --------------------------------------------------------------
// Nest pages by selfID (root = null) and order siblings by weight, then name.
// Both standalone and module-detail pages live in this one tree.
const flatPages = computed(() => store.pagesFor(props.project.projectID))

const tree = computed(() => {
  const byId = new Map(flatPages.value.map(p => [p.id, { ...p, children: [] }]))
  const roots = []
  for (const node of byId.values()) {
    // A missing/self parent (deleted parent, cross-namespace, or a page pointing
    // at itself) falls back to root so the page still renders.
    const parent = node.selfID && node.selfID !== node.id ? byId.get(node.selfID) : null
    if (parent) parent.children.push(node)
    else roots.push(node)
  }
  sortSiblings(roots)
  return roots
})

function sortSiblings(list) {
  list.sort((a, b) => a.weight - b.weight || a.name.localeCompare(b.name))
  list.forEach(n => sortSiblings(n.children))
}

// Expand state, default expanded (page trees are shallow).
const open = reactive(new Map())
const isOpen = id => (open.has(id) ? open.get(id) : true)
const toggle = id => open.set(id, !isOpen(id))

// Depth-first flatten of the currently-visible rows (one card each).
const visibleRows = computed(() => {
  const rows = []
  const walk = (nodes, depth) => {
    for (const node of nodes) {
      const hasChildren = node.children.length > 0
      rows.push({ ...node, depth, hasChildren })
      if (hasChildren && isOpen(node.id)) walk(node.children, depth + 1)
    }
  }
  walk(tree.value, 0)
  return rows
})

// --- Row presentation ------------------------------------------------------------
// Record (module detail) pages show the module they belong to as a badge.
const modulesById = computed(() => {
  const map = new Map()
  for (const r of store.resourcesFor(props.project.projectID)) {
    if (r.kind === 'module') map.set(r.id, r)
  }
  return map
})

const moduleFor = item => modulesById.value.get(item.moduleID) || null

function onRemove(p) {
  confirmDelete({
    header: t('project.pages.removeConfirm.header'),
    message: t('project.pages.removeConfirm.message', { name: p.name }),
    onConfirm: async () => {
      try {
        await store.removePage(props.project.projectID, p.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.pages.toastRemoveFailed'))(err)
      }
    },
  })
}

// --- Drag & drop -----------------------------------------------------------------
// Native HTML5 DnD: drop on a row's middle band nests the dragged page inside it;
// dropping near the top/bottom edge reorders it as a sibling before/after. The
// drop is persisted (reparent + sibling order) then the store refetches.
const dragId = ref(null)
const dropTarget = reactive({ id: null, pos: null })

function clearDrag() {
  dragId.value = null
  dropTarget.id = null
  dropTarget.pos = null
}

function onDragStart(e, row) {
  if (props.disabled) return
  dragId.value = row.id
  e.dataTransfer.effectAllowed = 'move'
  // Firefox only fires drag events once dataTransfer carries something.
  e.dataTransfer.setData('text/plain', row.id)
}

function onDragOver(e, row) {
  if (dragId.value == null || row.id === dragId.value) return
  e.preventDefault()
  e.dataTransfer.dropEffect = 'move'
  const rect = e.currentTarget.getBoundingClientRect()
  const y = e.clientY - rect.top
  // Outer 30% bands reorder as a sibling; the middle nests as a child.
  const pos = y < rect.height * 0.3 ? 'before' : y > rect.height * 0.7 ? 'after' : 'inside'
  dropTarget.id = row.id
  dropTarget.pos = pos
}

function onDragLeave(e, row) {
  if (dropTarget.id === row.id) {
    dropTarget.id = null
    dropTarget.pos = null
  }
}

function onDrop(e, row) {
  e.preventDefault()
  const draggedId = dragId.value
  const targetId = row.id
  const pos = dropTarget.pos
  clearDrag()
  if (!draggedId || draggedId === targetId || !pos) return
  movePage(draggedId, targetId, pos)
}

function onDragEnd() {
  clearDrag()
}

// Is `nodeId` inside the subtree rooted at `ancestorId` (walk up parent links)?
// Guards against dropping a page into one of its own descendants.
function isWithinSubtree(byId, nodeId, ancestorId) {
  let cur = nodeId
  const seen = new Set()
  while (cur && !seen.has(cur)) {
    if (cur === ancestorId) return true
    seen.add(cur)
    cur = byId.get(cur)?.selfID || null
  }
  return false
}

// Ordered child page ids under a parent (null = root), by weight then name.
function orderedChildren(flat, parentId) {
  return flat
    .filter(p => (p.selfID || null) === (parentId || null))
    .sort((a, b) => a.weight - b.weight || a.name.localeCompare(b.name))
    .map(p => p.id)
}

async function movePage(draggedId, targetId, pos) {
  const flat = flatPages.value
  const byId = new Map(flat.map(p => [p.id, p]))
  const dragged = byId.get(draggedId)
  const target = byId.get(targetId)
  if (!dragged || !target) return

  // New parent + the full sibling order under it (with the dragged page inserted).
  let newParent
  let siblings
  if (pos === 'inside') {
    if (isWithinSubtree(byId, targetId, draggedId)) return // no dropping into own subtree
    newParent = targetId
    siblings = orderedChildren(flat, targetId).filter(id => id !== draggedId)
    siblings.push(draggedId)
  } else {
    newParent = target.selfID || null
    if (isWithinSubtree(byId, newParent, draggedId)) return
    siblings = orderedChildren(flat, newParent).filter(id => id !== draggedId)
    const idx = siblings.indexOf(targetId)
    siblings.splice(pos === 'before' ? idx : idx + 1, 0, draggedId)
  }

  try {
    if ((dragged.selfID || null) !== (newParent || null)) {
      await store.updatePage(props.project.projectID, draggedId, { selfID: newParent || '0' })
    }
    await store.reorderPages(props.project.projectID, newParent, siblings)
    if (pos === 'inside') open.set(targetId, true) // reveal the newly nested child
  } catch (err) {
    $toast.toastErrorHandler(t('project.pages.toastMoveFailed'))(err)
  } finally {
    // Weights are stamped server-side; refetch so the tree matches.
    try {
      await store.loadPages(props.project.projectID)
    } catch {
      /* toast above already surfaced the failure */
    }
  }
}

async function refresh(id) {
  if (!id) return
  try {
    // Load resources too so module-detail rows can resolve their module name.
    await Promise.all([store.loadResources(id), store.loadPages(id)])
  } catch (err) {
    $toast.toastErrorHandler(t('project.pages.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.projectID))
watch(
  () => props.project.projectID,
  id => id && refresh(id),
)
</script>
