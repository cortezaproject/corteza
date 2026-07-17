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

      <div v-else class="backlog-list shrink-0">
        <CResourceList
          class="h-full"
          primary-key="id"
          :fields="fields"
          :items="visibleItems"
          :filter="filter"
          @update:filter="Object.assign(filter, $event)"
          :sorting="sorting"
          :pagination="pagination"
          :translations="{
            searchPlaceholder: $t('project.dashboard.backlog.searchPlaceholder'),
            noItems: $t('project.dashboard.list.empty'),
          }"
          clickable
          @sort="onSort"
          @row-click="onRowClick"
        >
          <template #header>
            <Button
              icon="pi pi-plus"
              :label="$t('project.dashboard.backlog.newButton')"
              size="small"
              @click="openCreate"
            />
          </template>

          <template #body-id="{ data }">
            <span class="font-mono text-xs text-muted-color">{{ data.id }}</span>
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
    </div>

    <!-- Create + row-click edit share one dialog — see BacklogItemDialog. -->
    <BacklogItemDialog
      v-model:visible="dialogVisible"
      :record="selectedItem"
      :user-options="eventsStore.ownerOptions"
      :on-save="onSave"
      :on-delete="onDelete"
    />
  </div>
</template>

<script setup>
import BacklogItemDialog from '@/sections/project/components/dashboard/BacklogItemDialog.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { orderIndex } from '@/sections/project/config/chartColors'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CResourceList, CEmptyState } = components

defineOptions({ name: 'BacklogView' })

const { t } = useI18n()
const route = useRoute()
const store = useBacklogItemsStore()
const eventsStore = useEventsStore()
const $toast = inject('$toast')

// Title-bar badge — neutral (this page spans every category, so it doesn't
// borrow any single category's colour), same icon-square shape as
// CategoryView's per-category badge.
const BADGE = { icon: 'pi pi-th-large', bg: 'bg-emphasis', ring: 'ring-surface', text: 'text-color' }

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

const filter = reactive({ query: '' })
const sorting = reactive({ sortBy: 'dateDue', sortDesc: true })
const pagination = reactive({
  limit: 50,
  pageCursor: undefined,
  prevPage: '',
  nextPage: '',
  total: 0,
  page: 1,
})

const fields = computed(() => [
  { key: 'id', header: t('project.dashboard.columns.id'), sortable: true },
  { key: 'title', header: t('project.dashboard.columns.title'), sortable: true },
  { key: 'eventID', header: t('project.dashboard.backlog.columns.linkedEvent'), sortable: false },
  { key: 'assignee', header: t('project.dashboard.backlog.f.assignee'), sortable: true },
  { key: 'priority', header: t('project.dashboard.backlog.f.priority'), sortable: true },
  { key: 'status', header: t('project.dashboard.event.f.status'), sortable: true },
  { key: 'dateDue', header: t('project.dashboard.event.f.dateDue'), sortable: true },
])

// Client-side query filter across the item's string values.
const filteredItems = computed(() => {
  const q = (filter.query || '').trim().toLowerCase()
  if (!q) return store.items
  return store.items.filter(item =>
    Object.values(item).some(v => typeof v === 'string' && v.toLowerCase().includes(q)),
  )
})

// Priority is ordinal (High > Medium > Low) like severity/risk/status, so it
// sorts by rank rather than alphabetically.
const PRIORITY_ORDER = ['High', 'Medium', 'Low']
const priorityIndex = label => {
  const i = PRIORITY_ORDER.indexOf(String(label ?? ''))
  return i === -1 ? Number.MAX_SAFE_INTEGER : i
}
const RANKED_COLUMNS = { status: label => orderIndex('status', label), priority: priorityIndex }

const visibleItems = computed(() => {
  const list = [...filteredItems.value]
  const { sortBy, sortDesc } = sorting
  if (sortBy) {
    const rank = RANKED_COLUMNS[sortBy]
    list.sort((a, b) => {
      if (rank) {
        const ai = rank(a[sortBy])
        const bi = rank(b[sortBy])
        if (ai !== bi) return sortDesc ? bi - ai : ai - bi
        return 0
      }
      const av = a[sortBy] ?? ''
      const bv = b[sortBy] ?? ''
      if (av < bv) return sortDesc ? 1 : -1
      if (av > bv) return sortDesc ? -1 : 1
      return 0
    })
  }
  return list
})

watch(visibleItems, list => { pagination.total = list.length }, { immediate: true })

const onSort = ({ sortField, sortOrder }) => {
  if (!sortField) return
  sorting.sortBy = sortField
  sorting.sortDesc = sortOrder === -1
}

const dialogVisible = ref(false)
const selectedItem = ref(null)

function openCreate() {
  selectedItem.value = null
  dialogVisible.value = true
}

function onRowClick({ data }) {
  selectedItem.value = data
  dialogVisible.value = true
}

// Save handler — passed to BacklogItemDialog's `onSave` prop; branches
// create/update off whether `id` is set (see the dialog's own contract
// comment). Same truthy/falsy-or-throw idiom as the event dialogs.
const onSave = async (id, payload) => {
  try {
    if (id) {
      await store.update(id, payload)
      $toast.toastSuccess(t('project.dashboard.backlog.singular'), t('project.dashboard.backlog.toast.updated'))
    } else {
      await store.add(payload)
      $toast.toastSuccess(t('project.dashboard.backlog.singular'), t('project.dashboard.backlog.toast.created'))
    }
    return true
  } catch (err) {
    console.error('Failed to save backlog item', err)
    const key = id ? 'project.dashboard.backlog.toast.updateFailed' : 'project.dashboard.backlog.toast.createFailed'
    $toast.toastErrorHandler(t(key))(err)
    return false
  }
}

const onDelete = async id => {
  try {
    await store.remove(id)
    $toast.toastSuccess(t('project.dashboard.backlog.singular'), t('project.dashboard.backlog.toast.deleted'))
    return true
  } catch (err) {
    console.error('Failed to delete backlog item', err)
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
    return false
  }
}

// Switching projects resets transient list/dialog state (mirrors
// CategoryView's category-switch watch); store data stays live (reloaded by
// DashboardLayout).
watch(
  () => route.params.projectId,
  () => {
    filter.query = ''
    dialogVisible.value = false
    selectedItem.value = null
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
