<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.pages.add')"
          size="small"
          @click="createResource?.('page')"
        />
      </div>

      <!-- Pages render as their nav hierarchy (parent/child via selfID). Each row
           is a card indented by its depth; a chevron expands/collapses children.
           Standalone and module-detail pages share one tree. -->
      <div class="mt-1 flex flex-col gap-2">
        <div
          v-if="!visibleRows.length"
          class="text-muted-color p-4 border border-surface rounded-border bg-emphasis text-center"
        >
          {{ $t('project.pages.empty') }}
        </div>

        <div
          v-for="item in visibleRows"
          :key="item.id"
          class="group flex items-stretch border border-surface bg-surface rounded-border shadow-sm hover:bg-emphasis transition-colors overflow-hidden"
          :style="{ marginLeft: item.depth * 2 + 'rem' }"
        >
          <!-- Chevron only for pages that have children — no placeholder for leaves,
               so leaf rows don't carry an empty column. -->
          <button
            v-if="item.hasChildren"
            type="button"
            class="self-stretch flex items-center px-4 shrink-0 cursor-pointer text-muted-color hover:bg-surface-200 dark:hover:bg-surface-700 transition-colors"
            :aria-label="$t(isOpen(item.id) ? 'project.pages.collapse' : 'project.pages.expand')"
            :title="$t(isOpen(item.id) ? 'project.pages.collapse' : 'project.pages.expand')"
            @click.stop="toggle(item.id)"
          >
            <i
              class="pi pi-chevron-down text-xs transition-transform duration-200"
              :class="{ '-rotate-90': !isOpen(item.id) }"
            />
          </button>

          <!-- Page name is the title; record (module detail) pages show a
               bordered module badge (colored icon square + module name) to the
               right, vertically centered. Click opens the page detail dialog. -->
          <div
            class="flex-1 min-w-0 flex items-center gap-2 py-2.5 pr-3 cursor-pointer"
            :class="item.hasChildren ? 'ps-1' : 'ps-3'"
            @click="inspectResource?.('page', item.id)"
          >
            <span class="truncate font-medium">{{ item.name }}</span>
            <!-- Module the record page belongs to, as a kind badge. -->
            <KindBadge
              v-if="item.isRecordPage && moduleFor(item)"
              kind="module"
              :label="moduleFor(item).name"
            />
          </div>

          <!-- Right side: the status Tag swaps for the hover actions (open in
               builder, delete) — the actions overlay the Tag's slot so the row
               width stays put. Mirrors the reveal-on-hover list behaviour. -->
          <div class="relative self-stretch flex items-center shrink-0 ps-2 pe-3">
            <div class="transition-opacity group-hover:opacity-0 group-focus-within:opacity-0">
              <Tag
                :value="item.visible ? $t('project.pages.visible') : $t('project.pages.hidden')"
                :severity="item.visible ? 'success' : 'secondary'"
                class="!text-xs"
              />
            </div>
            <div
              class="absolute inset-y-0 end-0 pe-2 flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100"
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
import { computed, inject, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

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
const flatPages = computed(() => store.pagesFor(props.project.id))

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
  const sortSiblings = list => {
    list.sort((a, b) => a.weight - b.weight || a.name.localeCompare(b.name))
    list.forEach(n => sortSiblings(n.children))
  }
  sortSiblings(roots)
  return roots
})

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
  for (const r of store.resourcesFor(props.project.id)) {
    if (r.kind === 'module') map.set(r.id, r)
  }
  return map
})

// Module a record (module detail) page belongs to; rendered as a badge next to
// the page name (null until resolved).
const moduleFor = item => modulesById.value.get(item.moduleID) || null

function onRemove(p) {
  confirmDelete({
    header: t('project.pages.removeConfirm.header'),
    message: t('project.pages.removeConfirm.message', { name: p.name }),
    onConfirm: async () => {
      try {
        await store.removePage(props.project.id, p.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.pages.toastRemoveFailed'))(err)
      }
    },
  })
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

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
