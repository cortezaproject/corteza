<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('project.list.title') }}</span>
  </Teleport>

  <CViewContainer>
    <CResourceList
      ref="resourceListRef"
      primary-key="projectID"
      :fields="fields"
      :items="projects"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="actionItemsFor"
      :translations="{
        searchPlaceholder: $t('project.list.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('project.list.resourceSingle'),
        resourcePlural: $t('project.list.resourcePlural'),
        noItems: $t('project.list.noItems'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="onRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <Button
          icon="pi pi-plus"
          :label="$t('project.list.newProject')"
          size="small"
          @click="newDialogVisible = true"
        />
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span class="font-medium">{{ data.name }}</span>
          <span v-if="data.meta.description" class="text-xs text-muted-color">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-status="{ data }">
        <StatusChip v-if="data.status || data.archivedAt" v-bind="chipFor(data)" small />
        <span v-else>-</span>
      </template>

      <template #body-updatedAt="{ data }">
        <span class="text-sm text-muted-color">{{ formatDate(data.updatedAt) }}</span>
      </template>
      <template #filter>
        <Button
          icon="pi pi-filter"
          severity="secondary"
          size="small"
          text
          @click="toggleFilterMenu"
        />
      </template>
    </CResourceList>
    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('project.list.filterByStatus') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.status" :inputId="'draft'" value="draft" />
            <label for="draft" class="text-sm cursor-pointer">
              {{ $t('project.status.draft') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.status" :inputId="'published'" value="published" />
            <label for="published" class="text-sm cursor-pointer">
              {{ $t('project.status.published') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.status" :inputId="'archived'" value="archived" />
            <label for="archived" class="text-sm cursor-pointer">
              {{ $t('project.status.archived') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
    <NewProjectDialog v-model:visible="newDialogVisible" @created="onCreated" />
    <RenameProjectDialog
      v-model:visible="renameVisible"
      :project="renameTarget"
      @rename="onRename"
    />
  </CViewContainer>
</template>

<script setup>
import NewProjectDialog from '@/sections/project/components/project/NewProjectDialog.vue'
import RenameProjectDialog from '@/sections/project/components/project/RenameProjectDialog.vue'
import StatusChip from '@/sections/project/components/project/StatusChip.vue'
import { chainHasPublished } from '@/sections/project/config/publishState'
import { system } from '@planetcrust/human-js'
import { components, useResourceList } from '@planetcrust/human-vue'
import { useConfirm } from 'primevue/useconfirm'
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
const { CResourceList, CViewContainer } = components

const { t } = useI18n()
const router = useRouter()
const confirm = useConfirm()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// Backend-driven list: fetch/filter/sort/paginate all happen server-side via
// the shared composable (same pattern as every other resource list). Raw
// records are wrapped in the lib Project class so the template reads its
// getters (name, status, …) instead of a hand-rolled unmarshal.
const {
  items: projects,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(
  params => {
    // headsOnly: one row per revision chain (the row a user would act on),
    // not one row per revision — see project.intent.md / ProjectList.intent.md.
    // Sent unconditionally here (not via the reactive `filter`) so it can
    // never be cleared by a filter-menu change and never pollutes the URL.
    const { response, cancel } = $SystemAPI.projectListCancellable({
      ...params,
      ...statusFilterParams(params.status),
      headsOnly: true,
    })
    return {
      cancel,
      response: async () => {
        const result = await response()
        return { ...result, set: (result.set || []).map(p => new system.Project(p)) }
      },
    }
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'updatedAt', sortDesc: true },
    pagination: { limit: 50 },
  },
)

// The filter menu offers one radio group over three things that are not one
// axis on the backend: two lifecycle statuses and the archive shelf. It also
// speaks the user's vocabulary — "published" is `active` in the data model
// (the status a publish actually writes), and asking for a literal "published"
// now matches nothing. Both mappings live here so the menu can stay as it is.
//
// archived: 0 excludes shelved projects, 2 returns only those (pkg/filter
// State). Sent explicitly rather than by omission so the working set is the
// default no matter what the URL carries.
const statusFilterParams = status => {
  if (status === 'archived') return { status: undefined, archived: 2 }
  if (status === 'published') return { status: 'active', archived: 0 }
  return { status: status || undefined, archived: 0 }
}

const resourceListRef = ref()
const newDialogVisible = ref(false)
const renameVisible = ref(false)
const renameTarget = ref(null)

// name (meta.short) is a JSON column with no sort ident, so it's not sortable.
const fields = [
  { key: 'name', sortable: false, header: t('general.label.name') },
  { key: 'status', sortable: true, header: t('general.label.status') },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('project.list.columns.updatedAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]
// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}
// Archived wins over the lifecycle status in this column: for a shelved
// project "archived" is what the reader needs to know, and its status carries
// on saying whatever it said when it went on the shelf.
//
// Everything else defers to StatusChip, which already holds the palette for
// every status the backend writes. This column kept its own three-key table
// for a while, and its keys were the words a user says rather than the ones
// the data model stores — so `active` (what a publish actually writes) and
// `deprecated` (what a superseded revision becomes) matched nothing and every
// live project showed a dash in its status column.
//
// The one thing that stays local is the WORDING for `active`: the filter menu
// just above calls that state "Published", and a list whose filter and rows
// disagree about the name of the same state is worse than either name. The
// colour still comes from the shared table.
const chipFor = project => {
  const status = project.archivedAt ? 'archived' : project.status
  return { status, label: status === 'active' ? t('project.status.published') : '' }
}

const formatDate = date => {
  if (!date) return ''
  try {
    return new Date(date).toLocaleDateString()
  } catch {
    return String(date)
  }
}

// Each row is a chain head (headsOnly, see the list request above). Routing
// favours the dashboard once a chain has ever published — it's the project's
// home, and the wizard is entered from its revision switcher. A chain that has
// never published has no dashboard (see config/publishState.js), so it opens
// the wizard.
const onRowClick = ({ data }) => {
  const name = chainHasPublished(data) ? 'project.overview' : 'project.wizard'
  router.push({ name, params: { projectId: data.projectID } })
}

// Just-created-project flag: a `new=1` query param on the very first
// navigation into the wizard. Members are the first thing to define on a
// fresh project, so the wizard (see views/Wizard.vue / a follow-up members
// dialog) should read `route.query.new === '1'` once on mount to auto-open
// the members dialog, then strip it via router.replace so it never re-opens
// on a later visit, refresh or back-navigation.
const onCreated = project => {
  $toast.toastSuccess(project.name, t('project.list.toast.created'))
  router.push({
    name: 'project.wizard',
    params: { projectId: project.projectID },
    query: { new: '1' },
  })
}

// Run an API mutation and report the outcome — success toast only when the call
// actually went through, error toast with the API message otherwise. Re-lists
// from the backend on success so the row reflects current state.
async function apiCall(fn, success) {
  try {
    await fn()
    filterList()
    if (success) $toast.toastSuccess(success.detail, success.summary)
  } catch (err) {
    $toast.toastErrorHandler(t('project.list.toast.actionFailed'))(err)
  }
}

// Update payload from a Project instance, overriding only the changed fields;
// config round-trips whole.
const updateProject = (project, patch = {}) =>
  $SystemAPI.projectUpdate({
    projectID: project.projectID,
    handle: project.handle,
    status: patch.status ?? project.status,
    config: project.config,
    meta: patch.meta ?? project.meta,
    labels: project.labels,
    updatedAt: project.updatedAt,
  })

const onRename = name => {
  const project = renameTarget.value
  if (!project) return
  apiCall(() => updateProject(project, { meta: { ...project.meta, short: name } }), {
    summary: t('project.list.toast.renamed'),
    detail: name,
  })
}

// Per-row actions
const closeMenu = () => resourceListRef.value?.hideActionsMenu?.()

const openRename = project => {
  closeMenu()
  renameTarget.value = project
  renameVisible.value = true
}

// Archiving has its own endpoints: it is a shelf state, separate from the
// lifecycle status that publish and revision own. Sending `status: 'draft'` to
// unarchive — which is what this did while the two shared a field — could
// silently un-publish a live project and unlock its schema for editing.
const toggleArchive = project => {
  closeMenu()
  const archive = !project.archivedAt
  apiCall(
    () =>
      archive
        ? $SystemAPI.projectArchive({ projectID: project.projectID })
        : $SystemAPI.projectUnarchive({ projectID: project.projectID }),
    {
      summary: archive ? t('project.list.toast.archived') : t('project.list.toast.unarchived'),
      detail: project.name,
    },
  )
}

const confirmDelete = project => {
  closeMenu()
  confirm.require({
    header: t('project.list.confirmDelete.header'),
    message: t('project.list.confirmDelete.message', { name: project.name }),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: {
      label: t('project.list.confirmDelete.accept'),
      severity: 'danger',
      size: 'small',
    },
    accept: () =>
      apiCall(() => $SystemAPI.projectDelete({ projectID: project.projectID }), {
        summary: t('project.list.toast.deleted'),
        detail: project.name,
      }),
  })
}

const actionItemsFor = project => [
  {
    label: t('project.list.actions.rename'),
    icon: 'pi pi-pencil',
    command: () => openRename(project),
  },
  {
    label: project.archivedAt
      ? t('project.list.actions.unarchive')
      : t('project.list.actions.archive'),
    icon: 'pi pi-inbox',
    command: () => toggleArchive(project),
  },
  { separator: true },
  {
    label: t('general.label.delete'),
    icon: 'pi pi-trash',
    class: 'text-red-500',
    command: () => confirmDelete(project),
  },
]
</script>
