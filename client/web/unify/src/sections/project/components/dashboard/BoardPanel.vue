<template>
  <!-- The project's kanban board: every work item as a card in one of four
       shared status columns, drag-to-move between them. Route-free, so both
       components/wizard/manage/ManageBoard.vue (revision-scoped) and
       views/dashboard/BoardView.vue (chain-wide) mount the same board.
       Mirrors CategoryPanel.vue/ActivityPanel.vue.

       NO background on this container, unlike ActivityPanel.vue: the columns
       below are themselves bg-surface, so a bg-surface container would flatten
       them into one slab and the columns would stop reading as columns. This
       matches CategoryPanel/OverviewPanel, which also sit transparent inside
       DashboardLayout's bordered panel.

       DATA LOADING: unlike CategoryPanel/ActivityPanel, this component DOES
       own its own load — it pages the board endpoint (GET /project-board/,
       server/system/types/project_board.go) directly rather than reading
       eventsStore/backlogStore's already-loaded state, so it stays correct
       past those stores' 200-row-per-category cache (see loadBoard/
       loadMoreColumn below). It still injects both stores for their WRITES
       (create/save/delete/status-change — see the Drag & Drop / Quick-add /
       Detail drawer sections below) and for `ownerOptions`, which is why
       views/dashboard/DashboardLayout.vue and
       components/wizard/manage/ManageBoard.vue load them chain-wide and
       revision-scoped respectively. -->
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

    <!-- A board-fetch failure (network/backend error) gets a visible retry —
         same shape as CategoryPanel's own report-load failure state — rather
         than silently rendering four columns that all read as empty. -->
    <section
      v-else-if="loadFailed"
      class="flex-1 flex flex-col items-center justify-center text-center gap-2 py-8"
    >
      <span
        class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emphasis text-red-500"
      >
        <i class="pi pi-exclamation-triangle text-xl" />
      </span>
      <p class="text-sm font-medium text-color">{{ $t('project.dashboard.overview.error') }}</p>
      <Button
        type="button"
        size="small"
        severity="secondary"
        outlined
        :label="$t('project.dashboard.overview.retry')"
        @click="loadBoard"
      />
    </section>

    <!-- Columns always render, even with zero items — an empty board (chain-
         wide, or a fresh revision) must stay usable, with per-column quick-add
         (see BoardColumn.vue). No whole-board "no work items" note: the
         per-column empty affordance already carries that message. -->
    <div v-else class="flex-1 min-h-0 overflow-x-auto">
      <div class="h-full flex gap-3 p-4 min-w-max">
        <BoardColumn
          v-for="col in columns"
          :key="col.status"
          class="w-80 shrink-0"
          :status="col.status"
          :items="col.items"
          :total="col.total"
          :has-more="!!col.nextPage"
          :loading-more="col.loadingMore"
          :disabled="disabled"
          :dragged-key="draggedKey"
          :add-label="quickAddLabel"
          :show-revision="!isRevisionScoped"
          @drop-item="onDropItem(col.status, $event)"
          @add-item="onQuickAdd(col.status)"
          @card-dragstart="draggedKey = $event"
          @card-dragend="draggedKey = null"
          @card-click="onCardClick"
          @load-more="loadMoreColumn(col.status)"
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
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

// The project's board — six item types share it: the five event categories
// (stores/events.js) plus backlog items (stores/backlogItems.js), grouped
// into the shared four-status column set (config/eventForm.js EVENT_STATUS —
// review now carries the same four statuses as everything else). Rows come
// off the board endpoint (GET /project-board/, see server/system/types/
// project_board.go), NOT stores/events.js#events / stores/backlogItems.js#items
// — those cap at 200 rows per category (their own load()'s comment), so a
// board built on them silently drops cards and undercounts past that; the
// board endpoint pages each column server-side instead (see loadBoard/
// loadMoreColumn below). Both stores are still injected here, but only for
// their WRITE actions (add/update/remove/updateStatus/fetchOne) and
// `ownerOptions` — see each usage below.
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
  // chain-wide either: creating and moving work items is ungated, matching the
  // dashboard's own New-event button (CategoryPanel).
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()
const users = useProjectUsersStore()
const { revisionInfo } = useRevisionLabel()

const isRevisionScoped = computed(() => !!props.revisionId)

// Cards per column page — small enough that a full four-column load (the
// board's initial fetch, no `status`) stays cheap, generous enough that most
// projects never need "load more" at all. Also the `limit` passed to a single
// column's continuation fetch (see loadMoreColumn).
const BOARD_PAGE_LIMIT = 20

// --- Board state (loadBoard/loadMoreColumn own every write to this) --------
// EVENT_STATUS's order is the board's fixed column order (Open / In Progress
// / Ready to Test / Completed) — every item type shares it. Columns render
// unconditionally (see template) so the board is usable from cold; this list
// itself never collapses to nothing. Each column tracks its own items/total/
// nextPage/loadingMore — independent server-paged state, not a client-side
// grouping of some larger loaded set.
const columns = reactive(
  EVENT_STATUS.map(status => ({ status, items: [], total: 0, nextPage: null, loadingMore: false })),
)
const loading = ref(true)
const loadFailed = ref(false)

// One raw board-endpoint item (server/system/types/project_board.go's
// ProjectBoardItem) -> the shape BoardCard.vue renders: { key, id, itemType,
// linkedCategory?, title, status, severity?, priority?, assignee, dueDate,
// revisionLabel }. `owner` comes back as a raw user id (0 = unassigned) — the
// endpoint deliberately leaves name resolution to the frontend (see that
// type's own doc comment), same directory (stores/users.js) every other view
// resolves owners through. `key` is the drag/click payload BoardCard/
// BoardColumn round-trip through dataTransfer (see onDropItem/onCardClick).
function mapBoardItem(raw) {
  return {
    key: raw.itemType === 'backlog' ? `backlog:${raw.id}` : `${raw.itemType}:${raw.id}`,
    id: raw.id,
    itemType: raw.itemType,
    linkedCategory: raw.linkedCategory || '',
    title: raw.title,
    status: raw.status,
    severity: raw.severity || '',
    priority: raw.priority || '',
    assignee: raw.owner ? users.userName(raw.owner) : '',
    dueDate: raw.dueDate || '',
    revisionLabel: revisionInfo(raw.revisionID),
  }
}

// Initial (or full-refresh) load — every column, first page each. Sequence-
// guarded like OverviewPanel's loadAll: a project/revision switch (or a rapid
// retry click) mid-flight must not let a stale response overwrite a newer
// one's columns.
let loadSeq = 0
async function loadBoard() {
  if (!props.projectId) {
    loading.value = false
    return
  }
  const mySeq = ++loadSeq
  loading.value = true
  loadFailed.value = false
  try {
    const { columns: cols = [] } = await $SystemAPI.projectBoardBoard({
      projectID: props.projectId,
      revisionID: props.revisionId || undefined,
      limit: BOARD_PAGE_LIMIT,
    })
    if (mySeq !== loadSeq) return // stale — a newer load/retry is in flight
    for (const col of cols) {
      const target = columns.find(c => c.status === col.status)
      if (!target) continue
      target.items = (col.items || []).map(mapBoardItem)
      // Fresh load (no cursor): the endpoint omits `total` (json
      // `omitempty`) both for a genuinely empty column AND would for a
      // zero total in general — either way, missing here means 0, never
      // undefined/NaN in the column header's count pill.
      target.total = Number(col.total || 0)
      target.nextPage = col.nextPage || null
      target.loadingMore = false
    }
  } catch (err) {
    if (mySeq !== loadSeq) return
    console.error('Failed to load board', err)
    loadFailed.value = true
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

// One column's next page — status + pageCursor, per the endpoint's paging
// contract (server/system/service/project_board.go's Board/loadColumn).
async function loadMoreColumn(status) {
  const col = columns.find(c => c.status === status)
  if (!col || !col.nextPage || col.loadingMore) return
  col.loadingMore = true
  try {
    const { columns: cols = [] } = await $SystemAPI.projectBoardBoard({
      projectID: props.projectId,
      revisionID: props.revisionId || undefined,
      status,
      pageCursor: col.nextPage,
      limit: BOARD_PAGE_LIMIT,
    })
    const page = cols[0] // Board() with `status` set returns exactly that one column.
    if (page) {
      col.items.push(...(page.items || []).map(mapBoardItem))
      col.nextPage = page.nextPage || null
      // Continuation page: the endpoint NEVER sends `total` here by design
      // (loadColumn's own doc comment — total can't be combined with a page
      // cursor), so `page.total` is always undefined at this point. Do NOT
      // coalesce that to 0 the way loadBoard does for a fresh load — this
      // branch only exists in case that contract ever changes server-side,
      // and even then only overwrites when a value is actually present.
      if (page.total !== undefined) col.total = Number(page.total)
    }
  } catch (err) {
    console.error('Failed to load more board items', err)
    $toast.toastErrorHandler(t('project.dashboard.board.toast.loadMoreFailed'))(err)
  } finally {
    col.loadingMore = false
  }
}

watch([() => props.projectId, () => props.revisionId], () => loadBoard(), { immediate: true })

function findItem(key) {
  for (const col of columns) {
    const item = col.items.find(it => it.key === key)
    if (item) return item
  }
  return null
}

// --- Drag & drop -------------------------------------------------------------
// `draggedKey` is the one piece of state shared across columns/cards (which
// card, if any, is the current drag source), for the dimmed-source-card
// style; which item got dropped where is read off dataTransfer instead (see
// BoardColumn.vue), so this stays a single flat ref rather than needing
// provide/inject.
const draggedKey = ref(null)

// Move the card between this component's own column state immediately
// (optimistic), persist via the store, and roll the move back on failure —
// stores/projects.intent.md's "Update flows are optimistic with rollback on
// push failure", now implemented at the board's own column-state level
// rather than relying on a store mutation to ripple into this component's
// rendering (see the DATA LOADING note above: this component doesn't read
// the stores' cached lists for its rows any more, so their own optimistic
// patch — still there, see stores/events.js#updateStatus — is invisible to
// it either way).
//
// WHY STILL CALL eventsStore.updateStatus/backlogStore.updateStatus rather
// than writing to $SystemAPI directly here: the update endpoint is a full PUT
// (server/system/service/project_incident.go#Update replaces every field,
// not a merge), and the board's own card shape (ProjectBoardItem) is a slim
// projection — no description/risk/dates-other-than-due/etc — so this
// component can't safely build that PUT body itself. The stores already own
// the correct full-record body construction (including remapping owner
// fields back to `<key>Id`); reusing them keeps that logic in one place.
// Both stores' updateStatus was hardened alongside this change (see
// stores/events.js/backlogItems.js) to fetch the record fresh when it isn't
// already cached, rather than silently no-op'ing — needed here specifically,
// since the board can now show (and so drag) cards well past either store's
// own 200-row cache.
//
// STATUS-ONLY, chain-wide included: both stores' updateStatus build their PUT
// body off the full record (cached or freshly fetched), which already
// carries that record's existing `revisionID` untouched — this function
// never reads or sets revisionID itself, so dropping a card into another
// column changes ONLY its status. A chain-wide drag can never clear or
// reassign the item's revision.
async function onDropItem(status, key) {
  // Clear the drag flag here rather than relying on @card-dragend alone: a
  // card that changes column unmounts from its source list, so the native
  // dragend never reaches it and the dimmed styling would stick until some
  // later interaction re-rendered it. Cleared before the early returns so it
  // resets on a no-op drop too.
  draggedKey.value = null
  const item = findItem(key)
  if (!item || item.status === status) return
  const sourceCol = columns.find(c => c.status === item.status)
  const targetCol = columns.find(c => c.status === status)
  const idx = sourceCol?.items.findIndex(it => it.key === key) ?? -1
  if (!sourceCol || !targetCol || idx === -1) return

  const prevStatus = item.status
  const [moved] = sourceCol.items.splice(idx, 1)
  moved.status = status
  targetCol.items.unshift(moved)
  sourceCol.total = Math.max(0, sourceCol.total - 1)
  targetCol.total += 1

  try {
    if (item.itemType === 'backlog') {
      await backlogStore.updateStatus(item.id, status)
    } else {
      await eventsStore.updateStatus(item.itemType, item.id, status)
    }
  } catch (err) {
    // Roll back the optimistic move — best-effort re-insert at the source's
    // original index; exact position among same-status cards doesn't matter
    // (the board endpoint's own ordering is "all of one source before the
    // next", not chronological — see project_board.go's doc comment).
    const revertIdx = targetCol.items.findIndex(it => it.key === key)
    if (revertIdx !== -1) targetCol.items.splice(revertIdx, 1)
    moved.status = prevStatus
    sourceCol.items.splice(Math.min(idx, sourceCol.items.length), 0, moved)
    sourceCol.total += 1
    targetCol.total = Math.max(0, targetCol.total - 1)
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
// titles), plus the revisionID assignment described above. Refreshes the
// whole board on success (loadBoard) rather than splicing the new card into
// local state by hand — creation is comparatively rare (a deliberate,
// dialog-driven action, not a hot interaction like drag), so the simplicity
// of one authoritative refetch outweighs the extra round trip.
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
  loadBoard()
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
// see mapBoardItem above); null until a card is clicked.
const selectedItem = ref(null)
const editVisible = ref(false)
const isBacklogSelected = computed(() => selectedItem.value?.itemType === 'backlog')

// Fallback full record for a clicked card that ISN'T in either store's own
// cached list (see selectedRecord below) — fetched fresh on click (see
// onCardClick) via the stores' fetchOne, reset on every new click so a stale
// fallback never lingers behind a different selection.
const selectedRecordFallback = ref(null)

// Resolve the full record behind the clicked card — the drawers/dialogs need
// every field (description, risk, owners…), not just the board's slim
// projection. A computed (not a captured snapshot) so a save elsewhere in
// the same record (store.update splices a fresh object into the list) is
// reflected live without re-pointing anything by hand.
//
// Checks the store's own cache FIRST (cheap, synchronous, and already the
// freshest copy after a save/delete via this drawer): most boards are small
// enough that every card is inside the store's 200-row-per-category load.
// Falls back to selectedRecordFallback for a card the board endpoint reaches
// but the store's own capped load never cached (see onCardClick) — the exact
// scenario this whole slice exists to stop silently breaking.
const selectedRecord = computed(() => {
  const item = selectedItem.value
  if (!item) return null
  if (item.itemType === 'backlog') {
    return backlogStore.items.find(i => i.id === item.id) || selectedRecordFallback.value
  }
  return (
    eventsStore.events.find(e => e.category === item.itemType && e.id === item.id) ||
    selectedRecordFallback.value
  )
})

async function onCardClick(key) {
  const item = findItem(key)
  if (!item) return
  selectedItem.value = item
  selectedRecordFallback.value = null
  editVisible.value = false
  drawerVisible.value = true

  const cached =
    item.itemType === 'backlog'
      ? backlogStore.items.find(i => i.id === item.id)
      : eventsStore.events.find(e => e.category === item.itemType && e.id === item.id)
  if (cached) return
  try {
    selectedRecordFallback.value =
      item.itemType === 'backlog'
        ? await backlogStore.fetchOne(item.id)
        : await eventsStore.fetchOne(item.itemType, item.id)
  } catch (err) {
    // The drawer stays open with an empty record (its own title/summary
    // fields read blank) rather than a dedicated failure state — a secondary
    // detail fetch failing is rare and the user can just close and re-click.
    console.error('Failed to load board item details', err)
  }
}

// --- Event save/delete (drawer's Edit -> EventDetailDialog) -----------------
// Every save/delete below also refreshes the board (loadBoard) — an edit can
// change the card's status (moving it to another column), title, severity or
// owner, none of which this component would otherwise see since it no longer
// reads the stores' cached lists for its rows (see the DATA LOADING note).
async function onEventSave(id, payload) {
  const cat = selectedItem.value?.itemType
  try {
    await eventsStore.update(cat, id, payload)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[cat].singularKey),
      t('project.dashboard.event.toast.updated'),
    )
    loadBoard()
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
    loadBoard()
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
    loadBoard()
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
    loadBoard()
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
    loadBoard()
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
    loadBoard()
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}
</script>
