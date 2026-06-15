<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('project.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="id"
      :fields="fields"
      :items="visibleProjects"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="false"
      :action-items="actionItemsFor"
      :translations="{
        searchPlaceholder: $t('project.list.searchPlaceholder'),
        resourceSingle: $t('project.list.resourceSingle'),
        resourcePlural: $t('project.list.resourcePlural'),
        noItems: $t('project.list.noItems'),
      }"
      clickable
      @sort="onSort"
      @row-click="onRowClick"
    >
      <template #header>
        <Button icon="pi pi-plus" :label="$t('project.list.newProject')" size="small" @click="newDialogVisible = true" />
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span class="font-medium">{{ data.name }}</span>
          <span v-if="data.description" class="text-xs text-muted-color">{{ data.description }}</span>
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

      <template #body-progress="{ data }">
        <div v-if="showProgress(data)" class="flex items-center gap-2">
          <ProgressBar
            :value="progressPct(data)"
            :show-value="false"
            class="w-24"
            :pt="{ root: { style: 'height: 6px' } }"
          />
          <span class="text-xs text-muted-color whitespace-nowrap">
            {{ $t('project.list.gates', { approved: data.gatesApproved, total: totalGates(data) }) }}
          </span>
        </div>
        <span v-else class="text-muted-color">—</span>
      </template>

      <template #body-updatedAt="{ data }">
        <span class="text-sm text-muted-color">{{ formatDate(data.updatedAt) }}</span>
      </template>
    </CResourceList>

    <NewProjectDialog v-model:visible="newDialogVisible" @created="onCreated" />
    <RenameProjectDialog v-model:visible="renameVisible" :project="renameTarget" @rename="onRename" />
  </div>
</template>

<script setup>
import NewProjectDialog from '@/sections/project/components/project/NewProjectDialog.vue'
import RenameProjectDialog from '@/sections/project/components/project/RenameProjectDialog.vue'
import { gateCount } from '@/sections/project/config/pipeline'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { storeToRefs } from 'pinia'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList } = components

const { t } = useI18n()
const router = useRouter()
const store = useProjectsStore()
const confirm = useConfirm()
const toast = useToast()
const { projects } = storeToRefs(store)

store.load()

const resourceListRef = ref()
const newDialogVisible = ref(false)
const renameVisible = ref(false)
const renameTarget = ref(null)

const filter = reactive({ query: '' })
const sorting = reactive({ sortBy: 'updatedAt', sortDesc: true })
const pagination = reactive({
  limit: 100,
  pageCursor: undefined,
  prevPage: '',
  nextPage: '',
  total: 0,
  page: 1,
})

const fields = [
  { key: 'name', sortable: true, header: t('general.label.name') },
  { key: 'mode', sortable: true, header: t('project.list.columns.mode') },
  { key: 'status', sortable: true, header: t('general.label.status') },
  { key: 'progress', sortable: false, header: t('project.list.columns.progress') },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('project.list.columns.updatedAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const statusSeverity = status =>
  ({ published: 'success', draft: 'info', archived: 'secondary' })[status] ?? null

const totalGates = data => gateCount(data.mode)
const showProgress = data => data.mode === 'gated' && data.status === 'draft' && totalGates(data) > 0
const progressPct = data => Math.round((data.gatesApproved / totalGates(data)) * 100)

const formatDate = iso => {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString()
  } catch {
    return iso
  }
}

const filteredProjects = computed(() => {
  const q = (filter.query || '').trim().toLowerCase()
  if (!q) return projects.value
  return projects.value.filter(
    p => p.name?.toLowerCase().includes(q) || p.description?.toLowerCase().includes(q),
  )
})

const visibleProjects = computed(() => {
  const list = [...filteredProjects.value]
  const { sortBy, sortDesc } = sorting
  if (sortBy) {
    list.sort((a, b) => {
      const av = a[sortBy] ?? ''
      const bv = b[sortBy] ?? ''
      if (av < bv) return sortDesc ? 1 : -1
      if (av > bv) return sortDesc ? -1 : 1
      return 0
    })
  }
  return list
})

watch(visibleProjects, list => { pagination.total = list.length }, { immediate: true })

const onSort = ({ sortField, sortOrder }) => {
  if (!sortField) return
  sorting.sortBy = sortField
  sorting.sortDesc = sortOrder === -1
}

const onRowClick = ({ data }) => {
  const name = data.status === 'published' ? 'project.overview' : 'project.wizard'
  router.push({ name, params: { projectId: data.id } })
}

const onCreated = project => {
  toast.add({ severity: 'success', summary: t('project.list.toast.created'), detail: project.name, life: 2500 })
  router.push({ name: 'project.wizard', params: { projectId: project.id } })
}

// Run a store mutation and report the outcome — success toast only when the
// call actually went through, error toast with the API message otherwise.
async function apiCall(fn, success) {
  try {
    await fn()
    if (success) toast.add({ severity: 'success', life: 2500, ...success })
  } catch (err) {
    toast.add({ severity: 'error', summary: t('project.list.toast.actionFailed'), detail: err.message, life: 4000 })
  }
}

const onRename = name => {
  if (!renameTarget.value) return
  apiCall(() => store.updateProject(renameTarget.value.id, { name }), {
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
  apiCall(
    () => store.updateProject(project.id, { status: archive ? 'archived' : 'draft' }),
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
    rejectProps: { label: t('general.label.cancel'), severity: 'secondary', text: true },
    acceptProps: { label: t('project.list.confirmDelete.accept'), severity: 'danger' },
    accept: () =>
      apiCall(() => store.removeProject(project.id), {
        summary: t('project.list.toast.deleted'),
        detail: project.name,
      }),
  })
}

const actionItemsFor = project => [
  { label: t('project.list.actions.rename'), icon: 'pi pi-pencil', command: () => openRename(project) },
  {
    label: project.status === 'archived' ? t('project.list.actions.unarchive') : t('project.list.actions.archive'),
    icon: 'pi pi-inbox',
    command: () => toggleArchive(project),
  },
  { separator: true },
  { label: t('general.label.delete'), icon: 'pi pi-trash', class: 'text-red-500', command: () => confirmDelete(project) },
]
</script>
