<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 p-4 pb-3 border-b border-surface">
      <h2 class="text-lg font-medium mb-1">{{ $t('project.manage.views.board') }}</h2>
      <p class="text-sm text-muted-color">{{ $t('project.manage.board.blurb') }}</p>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <template v-else>
      <!-- Explanatory note only — every work item defaults to unassigned
           until something assigns it to a revision (quick-add below is the
           one path that assigns on create; everything else stays
           unassigned), so an open revision legitimately starts with nothing
           on it. This used to REPLACE the board with a full-panel empty
           state; now it's a slim note ABOVE the always-visible columns
           (each with its own per-column empty affordance — see
           BoardColumn.vue) so the board stays usable from cold: quick-add
           is right there in the column headers. -->
      <div
        v-if="!boardItems.length"
        class="shrink-0 mx-4 mt-3 flex items-start gap-2 rounded-lg border border-surface bg-emphasis/40 px-3 py-2"
      >
        <i class="pi pi-info-circle text-muted-color text-sm mt-0.5 shrink-0" />
        <p class="text-xs text-muted-color">
          <span class="font-medium text-color">{{ $t('project.manage.board.empty.title') }}</span>
          {{ $t('project.manage.board.empty.hint') }}
        </p>
      </div>

      <div class="flex-1 min-h-0 overflow-x-auto">
        <div class="h-full flex gap-3 p-4 min-w-max">
          <BoardColumn
            v-for="col in columns"
            :key="col.status"
            class="w-72 shrink-0"
            :status="col.status"
            :items="col.items"
            :dragged-key="draggedKey"
            :add-label="quickAddLabel"
            @drop-item="onDropItem(col.status, $event)"
            @add-item="onQuickAdd(col.status)"
            @card-dragstart="draggedKey = $event"
            @card-dragend="draggedKey = null"
            @card-click="onCardClick"
          />
        </div>
      </div>
    </template>

    <!-- Quick-add — the EXISTING event-create dialog (NewEventDialog), never
         a bespoke board form. Always creates a 'task' (the simplest, most
         generic of the five categories — see the column header comment in
         BoardColumn.vue); the task schema's own `taskType` field is still
         there for the user to refine, since we pass the category's real
         schema through unmodified bar one thing: a clone with the `status`
         field's default overridden to the clicked column's status, so the
         dialog opens pre-filled without NewEventDialog itself needing to
         know about columns. -->
    <NewEventDialog
      v-model:visible="createDialogVisible"
      category="task"
      :schema="quickAddSchema"
      :user-options="eventsStore.ownerOptions"
      :on-create="onQuickCreate"
    />

    <!-- Row-click detail drawer — mirrors the dashboard exactly (click a
         card -> read-only drawer in the shared right sidebar -> its Edit
         opens the dialog below). Exactly one of the two renders, keyed off
         which of the six item types was clicked; both share the same
         right-sidebar slot (rightSidebar below) since only one is ever
         showing at a time. -->
    <EventDetailDrawer
      v-if="selectedItem && !isBacklogSelected"
      v-model:visible="drawerVisible"
      :category="selectedItem.itemType"
      :record="selectedRecord"
      :user-options="eventsStore.ownerOptions"
      :revision-id="revisionId"
      @edit="editVisible = true"
    />
    <BacklogItemDrawer
      v-if="selectedItem && isBacklogSelected"
      v-model:visible="drawerVisible"
      :record="selectedRecord"
      @edit="editVisible = true"
      @open-event="openLinkedEvent"
    />

    <!-- Edit dialog — same schema/store as the drawer's record, opened from
         the drawer's own Edit button (stacks above it), never the drawer
         itself. -->
    <EventDetailDialog
      v-if="selectedItem && !isBacklogSelected"
      v-model:visible="editVisible"
      :category="selectedItem.itemType"
      :record="selectedRecord"
      :user-options="eventsStore.ownerOptions"
      :on-save="onEventSave"
      :on-delete="onEventDelete"
    />
    <BacklogItemDialog
      v-if="selectedItem && isBacklogSelected"
      v-model:visible="editVisible"
      :record="selectedRecord"
      :user-options="eventsStore.ownerOptions"
      :on-save="onBacklogSave"
      :on-delete="onBacklogDelete"
    />

    <!-- Linked event's own edit dialog — opened from the backlog drawer's
         linked-event row (mirrors BacklogView.vue's same handoff). `category`
         comes off the backlog record itself (always present), not the
         resolved event, which may still be null before a row is clicked. -->
    <EventDetailDialog
      v-if="selectedItem && isBacklogSelected"
      v-model:visible="linkedEventVisible"
      :category="linkedEventCategory"
      :record="linkedEventRecord"
      :user-options="eventsStore.ownerOptions"
      :on-save="onLinkedEventSave"
      :on-delete="onLinkedEventDelete"
    />
  </div>
</template>

<script setup>
import BacklogItemDialog from '@/sections/project/components/dashboard/BacklogItemDialog.vue'
import BacklogItemDrawer from '@/sections/project/components/dashboard/BacklogItemDrawer.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import EventDetailDrawer from '@/sections/project/components/dashboard/EventDetailDrawer.vue'
import NewEventDialog from '@/sections/project/components/dashboard/NewEventDialog.vue'
import BoardColumn from './BoardColumn.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { EVENT_STATUS } from '@/sections/project/config/eventForm'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, ref, watch } from 'vue'
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
  // Passed by Wizard.vue to every M&M section, but deliberately NOT enforced
  // here: creating and moving work items is ungated, matching the dashboard's
  // own New-event button (views/dashboard/CategoryView.vue), which has never
  // been capability-gated. Ruled 2026-07-28 — revisit when member roles and
  // what each may do are refined.
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
// payload BoardCard/BoardColumn round-trip through dataTransfer — also reused
// as the click payload (see onCardClick) so both interactions share one
// lookup (findItem).
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
// / Ready to Test / Completed) — every item type shares it. Columns render
// unconditionally (see template) so the board is usable from cold; this list
// itself never collapses to nothing.
const columns = computed(() =>
  EVENT_STATUS.map(status => ({
    status,
    items: boardItems.value.filter(it => it.status === status),
  })),
)

function findItem(key) {
  return boardItems.value.find(it => it.key === key)
}

// --- Drag & drop -------------------------------------------------------------
// `draggedKey` is the one piece of state shared across columns/cards (which
// card, if any, is the current drag source), for the dimmed-source-card
// style; which item got dropped where is read off dataTransfer instead (see
// BoardColumn.vue), so this stays a single flat ref rather than needing
// provide/inject.
const draggedKey = ref(null)

// Write immediately, optimistically (the store patches the item's status in
// place before the request resolves, so the card visibly moves to the target
// column right away), rolling back on failure — stores/projects.intent.md's
// "Update flows are optimistic with rollback on push failure", applied to
// these two stores' updateStatus() actions.
async function onDropItem(status, key) {
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

// --- Quick-add ---------------------------------------------------------------
// AGREED BEHAVIOUR: an item created from a revision's board is ASSIGNED TO
// THAT REVISION (revisionID = the open revision) — unlike items created
// elsewhere (dashboard's category/backlog "+ New" flows), which stay
// unassigned. Otherwise the new card would vanish the instant it's created:
// the board only ever shows items already assigned to this revision (see the
// scoped load() above). WIP: the backend doesn't persist revisionID on
// create yet (see stores/events.js#add's comment) — this is written to work
// the moment that lands.
const taskCfg = CATEGORY_CONFIG.task
const quickAddLabel = computed(() =>
  t('project.dashboard.newButton', { type: t(taskCfg.singularKey) }),
)

const createDialogVisible = ref(false)
const createStatus = ref(EVENT_STATUS[0])

// Clone the task category's real form schema, overriding only the `status`
// field's default to the clicked column's status — NewEventDialog seeds its
// create-mode model from each field's `default` (see its buildModel), so
// this is enough to pre-fill the dialog without NewEventDialog itself
// needing any board-specific prop. Everything else about the schema
// (including the task-specific `taskType` field) passes through untouched.
const quickAddSchema = computed(() =>
  taskCfg.formSchema.map(section => ({
    ...section,
    fields: section.fields.map(f =>
      f.key === 'status' ? { ...f, default: createStatus.value } : f,
    ),
  })),
)

function onQuickAdd(status) {
  createStatus.value = status
  createDialogVisible.value = true
}

// Passed to NewEventDialog's `on-create` prop — same contract/shape as
// views/dashboard/CategoryView.vue's onCreate (including the inline backlog
// widget's queued titles), plus the revisionID assignment above.
async function onQuickCreate(payload, backlogTitles = []) {
  let event
  try {
    event = await eventsStore.add('task', payload, revisionId.value)
  } catch (err) {
    console.error('Failed to create task', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.createFailed'))(err)
    return false
  }
  $toast.toastSuccess(quickAddLabel.value, t('project.dashboard.event.toast.created'))
  if (backlogTitles.length) {
    const results = await Promise.allSettled(
      backlogTitles.map(title =>
        backlogStore.add(
          { category: 'task', eventID: event.id, title, priority: 'Medium', status: 'Open' },
          revisionId.value,
        ),
      ),
    )
    const failed = results.find(r => r.status === 'rejected')
    if (failed) {
      console.error('Failed to create backlog item', failed.reason)
      $toast.toastErrorHandler(t('project.dashboard.backlog.toast.createFailed'))(failed.reason)
    }
  }
  return true
}

// --- Detail drawer / edit dialog ---------------------------------------------
// The drawer registers with the shared right-sidebar store so it behaves
// like every other right panel: opening one (TAQ config, notifications,
// agent, the dashboard's own event/backlog drawers…) closes the rest, and
// vice versa. One panel name covers both drawer components below — exactly
// one of them is ever mounted at a time (keyed off isBacklogSelected), so
// they can share a slot rather than needing two.
const rightSidebar = useRightSidebarStore()
const DRAWER_PANEL = 'project-wizard-item-detail'
const drawerVisible = computed({
  get: () => rightSidebar.isOpen(DRAWER_PANEL),
  set: v => (v ? rightSidebar.open(DRAWER_PANEL) : rightSidebar.close(DRAWER_PANEL)),
})
// Leaving the tab/section with the drawer open would strand the shared
// store's activePanel on our name — release it (mirrors CategoryView.vue /
// BacklogView.vue).
onUnmounted(() => rightSidebar.close(DRAWER_PANEL))

// The clicked board item (the normalized card shape, not the full record —
// see boardItems above); null until a card is clicked.
const selectedItem = ref(null)
const editVisible = ref(false)
const isBacklogSelected = computed(() => selectedItem.value?.itemType === 'backlog')

// Resolve the full store-mapped record behind the clicked card — the
// drawers/dialogs need every field (description, risk, owners…), not just
// the board's slim projection. A computed (not a captured snapshot) so a
// save elsewhere in the same record (store.update splices a fresh object
// into the list) is reflected live without re-pointing anything by hand.
const selectedRecord = computed(() => {
  const item = selectedItem.value
  if (!item) return null
  if (item.itemType === 'backlog') return backlogStore.items.find(i => i.id === item.id) || null
  return eventsStore.events.find(e => e.category === item.itemType && e.id === item.id) || null
})

function onCardClick(key) {
  const item = findItem(key)
  if (!item) return
  selectedItem.value = item
  editVisible.value = false
  drawerVisible.value = true
}

// --- Event save/delete (drawer's Edit -> EventDetailDialog) -----------------
async function onEventSave(id, payload) {
  const cat = selectedItem.value?.itemType
  try {
    await eventsStore.update(cat, id, payload)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[cat].singularKey),
      t('project.dashboard.event.toast.updated'),
    )
    return true
  } catch (err) {
    console.error('Failed to update event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.updateFailed'))(err)
    return false
  }
}

async function onEventDelete(id) {
  const cat = selectedItem.value?.itemType
  try {
    await eventsStore.remove(cat, id)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[cat].singularKey),
      t('project.dashboard.event.toast.deleted'),
    )
    drawerVisible.value = false
    selectedItem.value = null
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}

// --- Backlog save/delete (drawer's Edit -> BacklogItemDialog) ---------------
async function onBacklogSave(id, payload) {
  try {
    await backlogStore.update(id, payload)
    $toast.toastSuccess(
      t('project.dashboard.backlog.singular'),
      t('project.dashboard.backlog.toast.updated'),
    )
    return true
  } catch (err) {
    console.error('Failed to update backlog item', err)
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.updateFailed'))(err)
    return false
  }
}

async function onBacklogDelete(id) {
  try {
    await backlogStore.remove(id)
    $toast.toastSuccess(
      t('project.dashboard.backlog.singular'),
      t('project.dashboard.backlog.toast.deleted'),
    )
    drawerVisible.value = false
    selectedItem.value = null
    return true
  } catch (err) {
    console.error('Failed to delete backlog item', err)
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
    return false
  }
}

// --- Linked event handoff (BacklogItemDrawer's "open-event" row) -----------
// Mirrors views/dashboard/BacklogView.vue's same handoff: the backlog
// drawer's linked-event row opens that event's own edit dialog, stacked
// independently of the backlog edit dialog above. `linkedEventCategory`
// comes off the backlog record's own `category` field (always present) —
// NOT off the resolved event, which stays null until one is found.
const linkedEventVisible = ref(false)
const linkedEventCategory = computed(() =>
  isBacklogSelected.value ? selectedRecord.value?.category || '' : '',
)
const linkedEventRecord = computed(() => {
  const rec = isBacklogSelected.value ? selectedRecord.value : null
  if (!rec) return null
  return eventsStore.byCategory(rec.category).find(e => e.id === String(rec.eventID)) || null
})

function openLinkedEvent() {
  if (linkedEventRecord.value) linkedEventVisible.value = true
}

async function onLinkedEventSave(id, payload) {
  const cat = linkedEventCategory.value
  try {
    await eventsStore.update(cat, id, payload)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[cat].singularKey),
      t('project.dashboard.event.toast.updated'),
    )
    return true
  } catch (err) {
    console.error('Failed to update event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.updateFailed'))(err)
    return false
  }
}

async function onLinkedEventDelete(id) {
  const cat = linkedEventCategory.value
  try {
    await eventsStore.remove(cat, id)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[cat].singularKey),
      t('project.dashboard.event.toast.deleted'),
    )
    linkedEventVisible.value = false
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}
</script>
