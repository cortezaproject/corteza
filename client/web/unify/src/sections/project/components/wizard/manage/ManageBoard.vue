<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 p-4 pb-3 border-b border-surface">
      <h2 class="text-lg font-medium mb-1">{{ $t('project.manage.views.board') }}</h2>
      <p class="text-sm text-muted-color">{{ $t('project.manage.board.blurb') }}</p>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <!-- Empty state matters here: every work item defaults to unassigned
         until something assigns it to a revision, so an open revision
         legitimately has nothing to show — explain that instead of leaving
         what looks like a broken/blank board (same "explain, don't leave a
         dead end" idiom as AllEventsView's empty state). -->
    <div
      v-else-if="!boardItems.length"
      class="flex-1 flex flex-col items-center justify-center text-center gap-2 px-6"
    >
      <span
        class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emphasis text-muted-color"
      >
        <i class="pi pi-objects-column text-xl" />
      </span>
      <p class="text-sm font-medium text-color">{{ $t('project.manage.board.empty.title') }}</p>
      <p class="text-sm text-muted-color max-w-md">{{ $t('project.manage.board.empty.hint') }}</p>
    </div>

    <div v-else class="flex-1 min-h-0 overflow-x-auto">
      <div class="h-full flex gap-3 p-4 min-w-max">
        <BoardColumn
          v-for="col in columns"
          :key="col.status"
          class="w-72 shrink-0"
          :status="col.status"
          :items="col.items"
          :disabled="disabled"
          :dragged-key="draggedKey"
          @drop-item="onDropItem(col.status, $event)"
          @card-dragstart="draggedKey = $event"
          @card-dragend="draggedKey = null"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import BoardColumn from './BoardColumn.vue'
import { EVENT_STATUS } from '@/sections/project/config/eventForm'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

// The wizard's Manage & Monitor board — work items for the ONE revision open
// in the wizard (route.params.projectId), not the whole project (that's the
// live dashboard's scope; see project.intent.md's "Dashboards" locked
// contract). Six item types share this board: the five event categories
// (stores/events.js) plus backlog items (stores/backlogItems.js), grouped
// into the shared four-status column set (config/eventForm.js EVENT_STATUS —
// review now carries the same four statuses as everything else).
const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const $toast = inject('$toast')
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

const rootProjectId = computed(() => props.project?.rootProjectID || props.project?.projectID)
const revisionId = computed(() => props.project?.projectID)
const loading = computed(() => eventsStore.loading || backlogStore.loading)

// Scope both dashboard stores to (root project, this revision) — filed
// against the root, assigned to the revision, like an issue against a
// milestone. Re-fetches whenever the open revision changes (e.g. switching
// revisions via the wizard-header revision switcher).
watch(
  [rootProjectId, revisionId],
  ([pid, rid]) => {
    if (!pid || !rid) return
    eventsStore.load(pid, rid)
    backlogStore.load(pid, rid)
  },
  { immediate: true },
)

// Per-category field maps needed to normalize the five event resources into
// one card shape below — which field on the raw (store-mapped) record holds
// the primary owner. Severity/risk are on every category but review (see
// config/eventForm.js's EVENT_FORMS); backlog items carry `priority` instead
// of severity, and their own `assignee` field rather than a category-named
// owner.
const OWNER_KEY = {
  incident: 'issueOwner',
  feature: 'featureOwner',
  privacy: 'requestOwner',
  task: 'owner',
  review: 'reviewer',
}

// One normalized shape for every board item, regardless of which of the six
// underlying resources it came from: { key, id, itemType, linkedCategory?,
// title, status, severity?, priority?, assignee, dueDate }. `key` is the drag
// payload BoardCard/BoardColumn round-trip through dataTransfer.
const boardItems = computed(() => {
  const evItems = eventsStore.events.map(e => ({
    key: `${e.category}:${e.id}`,
    id: e.id,
    itemType: e.category,
    title: e.title,
    status: e.status,
    severity: e.severity || '',
    assignee: e[OWNER_KEY[e.category]] || '',
    dueDate: e.dateDue || '',
  }))
  const blItems = backlogStore.items.map(i => ({
    key: `backlog:${i.id}`,
    id: i.id,
    itemType: 'backlog',
    linkedCategory: i.category,
    title: i.title,
    status: i.status,
    priority: i.priority || '',
    assignee: i.assignee || '',
    dueDate: i.dateDue || '',
  }))
  return [...evItems, ...blItems]
})

// EVENT_STATUS's order is the board's fixed column order (Open / In Progress
// / Ready to Test / Completed) — every item type shares it.
const columns = computed(() =>
  EVENT_STATUS.map(status => ({
    status,
    items: boardItems.value.filter(it => it.status === status),
  })),
)

// --- Drag & drop -------------------------------------------------------------
// `draggedKey` is the one piece of state shared across columns/cards (which
// card, if any, is the current drag source), for the dimmed-source-card
// style; which item got dropped where is read off dataTransfer instead (see
// BoardColumn.vue), so this stays a single flat ref rather than needing
// provide/inject.
const draggedKey = ref(null)

function findItem(key) {
  return boardItems.value.find(it => it.key === key)
}

// Write immediately, optimistically (the store patches the item's status in
// place before the request resolves, so the card visibly moves to the target
// column right away), rolling back on failure — stores/projects.intent.md's
// "Update flows are optimistic with rollback on push failure", applied to
// these two stores' updateStatus() actions.
async function onDropItem(status, key) {
  if (props.disabled) return
  const item = findItem(key)
  if (!item || item.status === status) return
  try {
    if (item.itemType === 'backlog') {
      await backlogStore.updateStatus(item.id, status)
    } else {
      await eventsStore.updateStatus(item.itemType, item.id, status)
    }
  } catch (err) {
    $toast.toastErrorHandler(t('project.manage.board.toast.moveFailed'))(err)
  }
}
</script>
