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

      <template #body-mode="{ data }">
        <Tag
          :value="$t(`project.mode.${data.mode}`)"
          :severity="data.mode === 'gated' ? 'warn' : 'secondary'"
          :icon="data.mode === 'gated' ? 'pi pi-shield' : 'pi pi-unlock'"
        />
      </template>

      <template #body-status="{ data }">
        <Tag :value="$t(`project.status.${data.status}`)" :severity="statusSeverity(data.status)" />
      </template>

      <template #body-updatedAt="{ data }">
        <span class="text-sm text-muted-color">{{ formatDate(data.updatedAt) }}</span>
      </template>
    </CResourceList>

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
// getters (name, mode, …) instead of a hand-rolled unmarshal.
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
    const { response, cancel } = $SystemAPI.projectListCancellable(params)
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

const resourceListRef = ref()
const newDialogVisible = ref(false)
const renameVisible = ref(false)
const renameTarget = ref(null)

// name (meta.short) is a JSON column with no sort ident, so it's not sortable;
// mode is a top-level column, sortable server-side.
const fields = [
  { key: 'name', sortable: false, header: t('general.label.name') },
  { key: 'mode', sortable: true, header: t('project.list.columns.mode') },
  { key: 'status', sortable: true, header: t('general.label.status') },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('project.list.columns.updatedAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const statusSeverity = status =>
  ({ active: 'success', published: 'success', draft: 'info', archived: 'secondary' })[status] ??
  null

const formatDate = date => {
  if (!date) return ''
  try {
    return new Date(date).toLocaleDateString()
  } catch {
    return String(date)
  }
}

// A live project (active/published) opens its dashboard; anything still in the
// build/draft lifecycle opens the wizard.
const onRowClick = ({ data }) => {
  const live = ['active', 'published'].includes(data.status)
  const name = live ? 'project.overview' : 'project.wizard'
  router.push({ name, params: { projectId: data.projectID } })
}

const onCreated = project => {
  $toast.toastSuccess(project.name, t('project.list.toast.created'))
  router.push({ name: 'project.wizard', params: { projectId: project.projectID } })
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

// Update payload from a Project instance, overriding only the changed fields.
// Mode is immutable (also enforced server-side), so config round-trips whole.
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

const toggleArchive = project => {
  closeMenu()
  const archive = project.status !== 'archived'
  apiCall(() => updateProject(project, { status: archive ? 'archived' : 'draft' }), {
    summary: archive ? t('project.list.toast.archived') : t('project.list.toast.unarchived'),
    detail: project.name,
  })
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
    label:
      project.status === 'archived'
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
