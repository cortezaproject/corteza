<template>
  <!-- The project's kanban board: every work item as a card in one of four
       shared status columns, drag-to-move between them. Route-free, extracted
       from components/wizard/manage/ManageBoard.vue (which now just resolves
       the wizard's open revision into props below) so the live dashboard can
       mount the exact same board, chain-wide (see views/dashboard/BoardView.vue) —
       ruled 2026-07-28: the project dashboard and the wizard's Manage &
       Monitor tab are ONE surface differing only in scope (see
       views/views.intent.md / DashboardLayout.intent.md). Mirrors
       CategoryPanel.vue/ActivityPanel.vue's own extraction exactly.

       NO background on this container, unlike ActivityPanel.vue: the columns
       below are themselves bg-surface, so a bg-surface container would flatten
       them into one slab and the columns would stop reading as columns. This
       matches CategoryPanel/OverviewPanel, which also sit transparent inside
       DashboardLayout's bordered panel.

       DATA LOADING: like CategoryPanel, this component does NOT call
       eventsStore.load()/backlogStore.load() itself — it only reads their
       already-loaded state (events/items). Whoever mounts it owns the load:
       views/dashboard/DashboardLayout.vue loads both stores chain-wide (no
       revisionId) for every child route, BoardView.vue included; each
       components/wizard/manage/ManageBoard.vue loads them scoped to (root
       project, open revision) itself, watching [rootProjectId, revisionId] —
       the same pattern the Manage<Category>.vue files use. -->
  <div class="h-full flex flex-col min-h-0">
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 shrink-0 bg-emphasis ring-surface"
      >
        <i class="pi pi-objects-column text-primary" />
      </span>
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.board.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.board.desc') }}</p>
      </div>
    </header>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <!-- Columns always render, even with zero items — an empty board (chain-
         wide, or a fresh revision) must stay usable, with per-column quick-add
         (see BoardColumn.vue). The old whole-board "no work items" note that
         used to sit above the columns is gone — the per-column empty
         affordance already carries that message, and stacking a second one
         above it read as noise. -->
    <div v-else class="flex-1 min-h-0 overflow-x-auto">
      <div class="h-full flex gap-3 p-4 min-w-max">
        <BoardColumn
          v-for="col in columns"
          :key="col.status"
          class="w-80 shrink-0"
          :status="col.status"
          :items="col.items"
          :disabled="disabled"
          :dragged-key="draggedKey"
          :add-label="quickAddLabel"
          :show-revision="!isRevisionScoped"
          @drop-item="onDropItem(col.status, $event)"
          @add-item="onQuickAdd(col.status)"
          @card-dragstart="draggedKey = $event"
          @card-dragend="draggedKey = null"
          @card-click="onCardClick"
        />
      </div>
    </div>

    <!-- Quick-add — the EXISTING event-create dialog (NewEventDialog), never
         a bespoke board form. Always creates a 'task' (the simplest, most
         generic of the five categories — see the column header comment in
         BoardColumn.vue); the task schema's own `taskType` field is still
         there for the user to refine, since we pass the category's real
         schema through unmodified bar one thing: a clone with the `status`
         field's default overridden to the clicked column's status, so the
         dialog opens pre-filled without NewEventDialog itself needing to
         know about columns. No `allow-revision-select` here, chain-wide or
         scoped alike (see NewEventDialog's own prop comment, which
         deliberately keeps every board-scoped caller off it): the created
         item's revision assignment is fully owned by `props.revisionId`
         below, not a form field the user could fight with. -->
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
         showing at a time. `revision-id` is null/omitted chain-wide (the
         dashboard's own behaviour) and the open revision when scoped,
         mirroring CategoryPanel's own EventDetailDrawer usage. -->
    <EventDetailDrawer
      v-if="selectedItem && !isBacklogSelected"
      v-model:visible="drawerVisible"
      :category="selectedItem.itemType"
      :record="selectedRecord"
      :user-options="eventsStore.ownerOptions"
      :revision-id="props.revisionId"
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
import BoardColumn from '@/sections/project/components/dashboard/BoardColumn.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import EventDetailDrawer from '@/sections/project/components/dashboard/EventDetailDrawer.vue'
import NewEventDialog from '@/sections/project/components/dashboard/NewEventDialog.vue'
import { useRevisionLabel } from '@/sections/project/composables/useRevisionLabel'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { EVENT_STATUS } from '@/sections/project/config/eventForm'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

// The project's board — six item types share it: the five event categories
// (stores/events.js) plus backlog items (stores/backlogItems.js), grouped
// into the shared four-status column set (config/eventForm.js EVENT_STATUS —
// review now carries the same four statuses as everything else).
const props = defineProps({
  // Chain-root project id. Chain-wide: views/dashboard/BoardView.vue passes
  // route.params.projectId (a chain HEAD, which is the chain ROOT for an
  // original, never-revised project — see system.Project#rootProjectID).
  // Revision-scoped: components/wizard/manage/ManageBoard.vue passes
  // project.rootProjectID.
  projectId: { type: [String, Number], default: '' },
  // The OPEN revision (a lib system.Project instance's projectID) — absent
  // (null, the default) means the dashboard's chain-wide scope: every
  // revision's items, loaded by the mounting dashboard host (see the DATA
  // LOADING note above). Set, it's the wizard's Manage & Monitor
  // single-revision scope (project.intent.md "Dashboards" / wizard.intent.md).
  // Drives: whether each card shows a revision chip (chain-wide only — see
  // isRevisionScoped/BoardColumn's showRevision), whether the drawer's
  // revision-id prop is forwarded, and — the AGREED behaviour — that an item
  // created from this panel's quick-add is assigned to it (mirrors
  // CategoryPanel's onCreate exactly).
  revisionId: { type: [String, Number], default: null },
  // Passed by Wizard.vue to every M&M section, but deliberately NOT enforced
  // chain-wide either: creating and moving work items is ungated, matching
  // the dashboard's own New-event button (CategoryPanel), which has never
  // been capability-gated. Ruled 2026-07-28 — revisit when member roles and
  // what each may do are refined.
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const $toast = inject('$toast')
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()
const { revisionInfo } = useRevisionLabel()

const isRevisionScoped = computed(() => !!props.revisionId)
const loading = computed(() => eventsStore.loading || backlogStore.loading)

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
// title, status, severity?, priority?, assignee, dueDate, revisionLabel }.
// `key` is the drag payload BoardCard/BoardColumn round-trip through
// dataTransfer — also reused as the click payload (see onCardClick) so both
// interactions share one lookup (findItem). `revisionLabel` is resolved HERE,
// once per item via the shared useRevisionLabel composable (same treatment
// CategoryPanel/BacklogView use for their revisionID column, including the
// explicit "Unassigned" state) rather than in BoardCard itself — keeps that
// component a pure presentational leaf, like it already is for typeBadge/
// typeLabel, and avoids every rendered card independently subscribing to the
// revision chain. Only rendered chain-wide (BoardCard's showRevision), but
// always resolved so that stays a pure display decision.
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
    revisionLabel: revisionInfo(e.revisionID),
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
    revisionLabel: revisionInfo(i.revisionID),
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
//
// STATUS-ONLY, chain-wide included: eventsStore.updateStatus/
// backlogStore.updateStatus (stores/events.js / stores/backlogItems.js) both
// build their PUT body off the store's own cached record (`{ ...item, status }`),
// which already carries that record's existing `revisionID` untouched — this
// function never reads or sets revisionID itself, so dropping a card into
// another column changes ONLY its status. A chain-wide drag can never clear
// or reassign the item's revision.
async function onDropItem(status, key) {
  // Clear the drag flag here rather than relying on @card-dragend alone: a
  // card that changes column unmounts from its source list, so the native
  // dragend never reaches it and the dimmed styling would stick until some
  // later interaction re-rendered it. Cleared before the early returns so it
  // resets on a no-op drop too.
  draggedKey.value = null
  const item = findItem(key)
  if (!item || item.status === status) return
  try {
    if (item.itemType === 'backlog') {
      await backlogStore.updateStatus(item.id, status)
    } else {
      await eventsStore.updateStatus(item.itemType, item.id, status)
    }
  } catch (err) {
    $toast.toastErrorHandler(t('project.dashboard.board.toast.moveFailed'))(err)
  }
}

// --- Quick-add ---------------------------------------------------------------
// AGREED BEHAVIOUR (unchanged from the pre-extraction ManageBoard.vue):
// `props.revisionId` is forwarded to the store's add() as its trailing
// revisionId argument — there is no revision field on this dialog to fight
// with (see the template's NewEventDialog comment). Revision-scoped (the
// wizard board), that's always the open revision — otherwise the new card
// would vanish the instant it's created, since the board only ever shows
// items already assigned to this revision (see the DATA LOADING note's
// scoped load). Chain-wide, it's null, so the new item stays unassigned by
// default — same as the dashboard's other "+ New" flows.
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
// CategoryPanel's onCreate (including the inline backlog widget's queued
// titles), plus the revisionID assignment described above.
async function onQuickCreate(payload, backlogTitles = []) {
  let event
  try {
    event = await eventsStore.add('task', payload, props.revisionId)
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
          props.revisionId,
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
const DRAWER_PANEL = 'project-board-item-detail'
const drawerVisible = computed({
  get: () => rightSidebar.isOpen(DRAWER_PANEL),
  set: v => (v ? rightSidebar.open(DRAWER_PANEL) : rightSidebar.close(DRAWER_PANEL)),
})
// Leaving the view/section with the drawer open would strand the shared
// store's activePanel on our name — release it (mirrors CategoryPanel.vue /
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
