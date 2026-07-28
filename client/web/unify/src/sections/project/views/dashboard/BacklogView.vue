<template>
  <!-- Backlog: every follow-up work item across every category, as a flat,
       sortable/searchable table (reference demo's Backlog page) — replaces
       the old tag-grouped view now that backlog items are their own linked
       resource (see stores/backlogItems.js) rather than a comma-tag field on
       the category records. Title-bar idiom matches CategoryView. -->
  <div class="flex flex-col h-full min-w-0 overflow-y-auto">
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <KindIcon :config="BADGE" size="xl" />
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.backlog.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.backlog.desc') }}</p>
      </div>
    </header>

    <div class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-4">
      <!-- First load (or a project switch): skeleton rather than a
           zeroed-out page. -->
      <div v-if="store.loading" class="flex flex-col gap-2">
        <span
          v-for="n in 6"
          :key="n"
          class="h-10 rounded-lg bg-emphasis block animate-pulse motion-reduce:animate-none"
        />
      </div>

      <!-- Nothing created yet — a CTA gets the first item created without a
           trip through the (currently empty) table toolbar below. -->
      <CEmptyState v-else-if="!store.items.length">
        <p class="font-medium text-color">{{ $t('project.dashboard.backlog.empty') }}</p>
        <p class="text-xs mt-1">{{ $t('project.dashboard.backlog.emptyHint') }}</p>
        <Button
          icon="pi pi-plus"
          :label="$t('project.dashboard.backlog.newButton')"
          size="small"
          class="mt-3"
          @click="openCreate"
        />
      </CEmptyState>

      <template v-else>
        <CategoryKpiRow :kpis="kpiList" />

        <!-- Chart band — same shape as CategoryView's, computed client-side
             from the loaded items (no report endpoint, same trade as the KPI
             trio above; the created-over-time trend is omitted for the same
             reason — the 200-row load cap would silently truncate it). -->
        <div class="grid grid-cols-2 xl:grid-cols-3 gap-3">
          <CategoryDonutChart
            title-key="project.dashboard.chart.byStatus"
            :data="statusBreakdown"
            variant="status"
            :height="180"
          />
          <CategoryDonutChart
            title-key="project.dashboard.chart.byCategory"
            :data="categoryBreakdown"
            variant="category"
            :height="180"
          />
          <CategoryRankBar
            title-key="project.dashboard.chart.byPriority"
            :data="priorityBreakdown"
            variant="priority"
            :height="180"
          />
        </div>

        <div class="backlog-list shrink-0">
          <CResourceList
            class="h-full"
            primary-key="id"
            :fields="fields"
            :items="listItems"
            :filter="filter"
            @update:filter="Object.assign(filter, $event)"
            :sorting="sorting"
            :pagination="pagination"
            :loading="listLoading"
            :translations="{
              searchPlaceholder: $t('project.dashboard.backlog.searchPlaceholder'),
              noItems: $t('project.dashboard.list.empty'),
            }"
            clickable
            @sort="handleSort"
            @row-click="onRowClick"
            @page-change="handlePageChange"
          >
            <!-- There WAS an "Unassigned" filter chip/checkbox in this
                 toolbar — dropped, see the module doc comment's UNASSIGNED
                 FILTER note: a server page can't honestly claim "N
                 unassigned" when it only means "N unassigned on this page". -->
            <template #header>
              <Button
                icon="pi pi-plus"
                :label="$t('project.dashboard.backlog.newButton')"
                size="small"
                @click="openCreate"
              />
            </template>

            <template #body-title="{ data }">
              <span class="font-medium text-color">{{ data.title || '—' }}</span>
            </template>

            <template #body-eventID="{ data }">
              <div v-if="linkedEventFor(data)" class="flex items-center gap-2 min-w-0">
                <KindIcon :config="CATEGORY_CONFIG[data.category]?.badge" size="sm" />
                <span class="truncate max-w-40 text-sm text-color">
                  {{ linkedEventFor(data).title || $t('project.dashboard.event.untitled') }}
                </span>
              </div>
              <span v-else class="font-mono text-xs text-muted-color">#{{ data.eventID }}</span>
            </template>

            <!-- Chain-wide scope (project.intent.md "Dashboards"): this list
                 spans every revision, so each row says which one it belongs
                 to — unassigned items (no revisionID) render distinctly
                 rather than as a blank cell (see useRevisionLabel). -->
            <template #body-revisionID="{ data }">
              <span
                v-if="revisionInfo(data.revisionID).unassigned"
                class="inline-flex items-center gap-1 text-xs text-muted-color italic"
              >
                <i class="pi pi-question-circle" />
                {{ revisionInfo(data.revisionID).label }}
              </span>
              <Tag
                v-else
                :value="revisionInfo(data.revisionID).label"
                :severity="revisionInfo(data.revisionID).severity"
                class="!text-xs"
              />
            </template>

            <template #body-assignee="{ data }">
              <UserCell :name="data.assignee" />
            </template>

            <template #body-priority="{ data }">
              <EventBadge :value="data.priority" variant="priority" />
            </template>

            <template #body-status="{ data }">
              <EventBadge :value="data.status" variant="status" />
            </template>

            <template #body-dateDue="{ data }">
              <span class="text-sm text-muted-color">{{ formatDate(data.dateDue) }}</span>
            </template>
          </CResourceList>
        </div>
      </template>
    </div>

    <!-- Create + row-click edit share one dialog — see BacklogItemDialog.
         `allow-revision-select` offers the revisionID field on create too —
         this view is always chain-wide (no revisionId/board scope exists
         here), same as CategoryPanel's dashboard-mode "+ New". -->
    <BacklogItemDialog
      v-model:visible="dialogVisible"
      :record="selectedItem"
      :user-options="eventsStore.ownerOptions"
      allow-revision-select
      :on-save="onSave"
      :on-delete="onDelete"
    />

    <!-- Row-click detail drawer — read-only summary; its Edit button opens the
         edit dialog above without closing the drawer, and clicking the linked
         event opens that event's edit dialog below. -->
    <BacklogItemDrawer
      v-model:visible="drawerVisible"
      :record="selectedItem"
      @edit="dialogVisible = true"
      @open-event="openLinkedEvent"
    />

    <!-- Linked event's edit dialog — opened from the drawer's linked-event
         row; saves/deletes go through the events store (same handlers idiom
         as CategoryView's). -->
    <EventDetailDialog
      v-if="selectedItem"
      v-model:visible="eventDialogVisible"
      :category="selectedItem.category"
      :record="linkedEventRecord"
      :user-options="eventsStore.ownerOptions"
      :on-save="onEventSave"
      :on-delete="onEventDelete"
    />
  </div>
</template>

<script setup>
// THE LIST is server-paged directly off the backlog-item resource endpoint
// (GET /project-backlog-items/, see the useResourceList call below) via the same
// useResourceList idiom views/ProjectList.vue uses, with `incTotal` for a
// true count — NOT store.items (capped at 200 rows — that store's own load()
// comment), which stays the source for the KPI trio and chart band below
// (unchanged by this — backlog has no report endpoint, so those were already
// a client-side trade against the store's own cap, not something this pass
// fixes) and for linkedEventFor's lookup into the events store.
//
// UNASSIGNED FILTER — DROPPED, not just hidden (see the old template's
// Popover/Chip, now gone): the generated resource filters treat
// `revisionID = 0` as "no constraint", not "unassigned only" (see
// server/store/adapters/rdbms/filters.gen.go's `if f.RevisionID > 0` guard,
// applied identically across every one of these six resources) — there is no
// server-side way to ask for "only unassigned rows" today. Client-filtering
// one fetched PAGE to fake it would misreport a page's own unassigned count
// as the whole backlog's, which is worse than not offering the filter.
import BacklogItemDialog from '@/sections/project/components/dashboard/BacklogItemDialog.vue'
import BacklogItemDrawer from '@/sections/project/components/dashboard/BacklogItemDrawer.vue'
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryRankBar from '@/sections/project/components/dashboard/CategoryRankBar.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { useRevisionLabel } from '@/sections/project/composables/useRevisionLabel'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import {
  CATEGORY_COLORS,
  PRIORITY_ORDER,
  STATUS_ORDER,
} from '@/sections/project/config/chartColors'
import { mapBacklogRow, useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { toISODate } from '@/sections/project/stores/dateUtils'
import { isOpenStatus, useEventsStore } from '@/sections/project/stores/events'
import { components, useResourceList, useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CResourceList, CEmptyState } = components

defineOptions({ name: 'BacklogView' })

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const route = useRoute()
const store = useBacklogItemsStore()
const eventsStore = useEventsStore()
const { revisionInfo } = useRevisionLabel()
const $toast = inject('$toast')

// Title-bar badge — neutral (this page spans every category, so it doesn't
// borrow any single category's colour), same icon-square shape as
// CategoryView's per-category badge.
const BADGE = {
  icon: 'pi pi-th-large',
  bg: 'bg-emphasis',
  ring: 'ring-surface',
  text: 'text-color',
}

// Resolve a backlog item's linked event via the events store (already loaded
// alongside this store by DashboardLayout) — null when the event no longer
// exists (deleted) or hasn't loaded, in which case the column falls back to
// the raw `#<eventID>` reference.
function linkedEventFor(item) {
  return eventsStore.byCategory(item.category).find(e => e.id === String(item.eventID)) || null
}

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}

// --- Server-paged list -------------------------------------------------------
// Same idiom views/ProjectList.vue uses: useResourceList owns filter/sorting/
// pagination state and the actual fetch, calling straight through to
// projectBacklogItemList (with `incTotal`) rather than a new endpoint. This
// view is always chain-wide (no revisionId scope — see the class-level
// comment), so `projectID` is just `route.params.projectId`, read fresh on
// every call so a project switch (this component's instance is reused across
// one, see the route.params.projectId watcher below) picks it up.
const {
  items: listItems,
  loading: listLoading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  fetchItems: refetchList,
  filterList,
} = useResourceList(
  params => {
    const { response, cancel } = $SystemAPI.projectBacklogItemListCancellable({
      ...params,
      projectID: route.params.projectId,
    })
    return {
      cancel,
      response: async () => {
        const result = await response()
        return { ...result, set: (result.set || []).map(mapBacklogRow) }
      },
    }
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'dateDue', sortDesc: true },
    pagination: { limit: 50 },
  },
)

// KPI trio above the table — computed client-side from the loaded items
// (backlog has no report endpoint; CategoryView's trio comes from one).
// Same rules as everywhere else: open = not Completed (stores/events.js#
// isOpenStatus), overdue = open with a due date before today (dateDue is a
// plain YYYY-MM-DD string, so a string compare against today's ISO date works).
const kpiList = computed(() => {
  const items = store.items
  const todayISO = toISODate(new Date())
  const open = items.filter(i => isOpenStatus(i.status))
  const overdue = open.filter(i => i.dateDue && i.dateDue < todayISO)
  return [
    { labelKey: 'project.dashboard.kpi.total', value: items.length },
    { labelKey: 'project.dashboard.kpi.open', value: open.length },
    { labelKey: 'project.dashboard.kpi.overdue', value: overdue.length },
  ]
})

// Chart breakdowns — counts over the loaded items, in each dimension's
// canonical order, zero labels dropped (the donut/rank-bar empty states
// handle the all-zero case).
const countBy = field => {
  const counts = {}
  for (const i of store.items) {
    const key = String(i[field] ?? '')
    if (key) counts[key] = (counts[key] || 0) + 1
  }
  return counts
}

const statusBreakdown = computed(() => {
  const counts = countBy('status')
  return STATUS_ORDER.filter(l => counts[l]).map(label => ({ label, value: counts[label] }))
})

const priorityBreakdown = computed(() => {
  const counts = countBy('priority')
  return PRIORITY_ORDER.filter(l => counts[l]).map(label => ({ label, value: counts[label] }))
})

// Labels are the translated category titles (colorFor can't key on those, so
// each row pins its colour — see CategoryDonutChart's data prop).
const categoryBreakdown = computed(() => {
  const counts = countBy('category')
  return CATEGORY_ORDER.filter(k => counts[k]).map(key => ({
    label: t(CATEGORY_CONFIG[key].titleKey),
    value: counts[key],
    color: CATEGORY_COLORS[key],
    key,
  }))
})

const fields = computed(() => [
  { key: 'title', header: t('project.dashboard.columns.title'), sortable: true },
  { key: 'eventID', header: t('project.dashboard.backlog.columns.linkedEvent'), sortable: false },
  { key: 'revisionID', header: t('project.dashboard.columns.revision'), sortable: false },
  { key: 'assignee', header: t('project.dashboard.backlog.f.assignee'), sortable: true },
  { key: 'priority', header: t('project.dashboard.backlog.f.priority'), sortable: true },
  { key: 'status', header: t('project.dashboard.event.f.status'), sortable: true },
  { key: 'dateDue', header: t('project.dashboard.event.f.dateDue'), sortable: true },
])

const dialogVisible = ref(false)
const selectedItem = ref(null)

// The drawer registers with the shared right-sidebar store so it behaves
// like the app's other right panels (exclusive with TAQ/notifications/the
// event drawer — opening one closes the others).
const rightSidebar = useRightSidebarStore()
const drawerVisible = computed({
  get: () => rightSidebar.isOpen('project-backlog-item-detail'),
  set: v =>
    v
      ? rightSidebar.open('project-backlog-item-detail')
      : rightSidebar.close('project-backlog-item-detail'),
})

// Leaving the view with the drawer open would strand the shared store's
// active panel — close it on unmount.
onUnmounted(() => rightSidebar.close('project-backlog-item-detail'))

function openCreate() {
  selectedItem.value = null
  dialogVisible.value = true
}

function onRowClick({ data }) {
  selectedItem.value = data
  drawerVisible.value = true
}

// Save handler — passed to BacklogItemDialog's `onSave` prop; branches
// create/update off whether `id` is set (see the dialog's own contract
// comment). Same truthy/falsy-or-throw idiom as the event dialogs. Refreshes
// the server-paged list on success — it's no longer store-driven (see the
// module doc comment), so a create/update wouldn't otherwise show up:
// refetchList() (keep the current page) for an in-place edit, filterList()
// (back to page 1) for a create, mirroring CategoryPanel's own onUpdate/
// onCreate split.
const onSave = async (id, payload) => {
  try {
    if (id) {
      await store.update(id, payload)
      if (selectedItem.value?.id === String(id)) {
        selectedItem.value = store.items.find(i => i.id === String(id)) || selectedItem.value
      }
      $toast.toastSuccess(
        t('project.dashboard.backlog.singular'),
        t('project.dashboard.backlog.toast.updated'),
      )
      refetchList()
    } else {
      await store.add(payload)
      $toast.toastSuccess(
        t('project.dashboard.backlog.singular'),
        t('project.dashboard.backlog.toast.created'),
      )
      filterList()
    }
    return true
  } catch (err) {
    console.error('Failed to save backlog item', err)
    const key = id
      ? 'project.dashboard.backlog.toast.updateFailed'
      : 'project.dashboard.backlog.toast.createFailed'
    $toast.toastErrorHandler(t(key))(err)
    return false
  }
}

const onDelete = async id => {
  try {
    await store.remove(id)
    if (selectedItem.value?.id === String(id)) {
      drawerVisible.value = false
      selectedItem.value = null
    }
    $toast.toastSuccess(
      t('project.dashboard.backlog.singular'),
      t('project.dashboard.backlog.toast.deleted'),
    )
    refetchList()
    return true
  } catch (err) {
    console.error('Failed to delete backlog item', err)
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
    return false
  }
}

// Linked event's edit dialog — opened from the drawer. `linkedEventRecord` is
// resolved live from the events store so a save repoints it automatically
// (store.update splices a fresh object into the list).
const eventDialogVisible = ref(false)
const linkedEventRecord = computed(() =>
  selectedItem.value ? linkedEventFor(selectedItem.value) : null,
)

function openLinkedEvent() {
  if (!linkedEventRecord.value) return
  eventDialogVisible.value = true
}

// Save/delete for the linked event — events store + toasts, same contract as
// CategoryView's onUpdate/onDelete (no metrics reload: this page's charts
// count backlog items, not events).
const onEventSave = async (id, payload) => {
  const category = selectedItem.value?.category
  try {
    await eventsStore.update(category, id, payload)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[category].singularKey),
      t('project.dashboard.event.toast.updated'),
    )
    return true
  } catch (err) {
    console.error('Failed to update event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.updateFailed'))(err)
    return false
  }
}

const onEventDelete = async id => {
  const category = selectedItem.value?.category
  try {
    await eventsStore.remove(category, id)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[category].singularKey),
      t('project.dashboard.event.toast.deleted'),
    )
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}

// Switching projects resets transient list/dialog state (mirrors
// CategoryPanel's category-switch watch) and re-runs the server-paged list
// from its first page (filterList) — store data stays live (reloaded by
// DashboardLayout).
watch(
  () => route.params.projectId,
  () => {
    filter.query = ''
    dialogVisible.value = false
    drawerVisible.value = false
    eventDialogVisible.value = false
    selectedItem.value = null
    sorting.sortBy = 'dateDue'
    sorting.sortDesc = true
    filterList()
  },
)
</script>

<style scoped>
/* Same fixed-height idiom as CategoryView's `.category-list` — Tailwind's
   scale has no 70vh utility and arbitrary bracket values are off the table,
   so this lives in plain CSS. */
.backlog-list {
  height: 70vh;
}
</style>
