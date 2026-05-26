<template>
  <Teleport to="#topbar-title" defer>
    <span>Projects</span>
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
        searchPlaceholder: 'Search projects…',
        resourceSingle: 'project',
        resourcePlural: 'projects',
        noItems: 'No projects yet.',
      }"
      clickable
      @sort="onSort"
      @row-click="onRowClick"
    >
      <template #header>
        <Button icon="pi pi-plus" label="New Project" size="small" @click="newDialogVisible = true" />
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span class="font-medium">{{ data.name }}</span>
          <span v-if="data.description" class="text-xs text-muted-color">{{ data.description }}</span>
        </div>
      </template>

      <template #body-mode="{ data }">
        <Tag
          :value="capitalize(data.mode)"
          :severity="data.mode === 'gated' ? 'warn' : 'secondary'"
          :icon="data.mode === 'gated' ? 'pi pi-shield' : 'pi pi-unlock'"
        />
      </template>

      <template #body-status="{ data }">
        <Tag :value="capitalize(data.status)" :severity="statusSeverity(data.status)" />
      </template>

      <template #body-version="{ data }">
        <span v-if="store.currentVersion(data)" class="font-medium">v{{ store.currentVersion(data).number }}</span>
        <span v-else class="text-muted-color">—</span>
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
            {{ data.gatesApproved }}/{{ totalGates(data) }} gates
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
import { useRouter } from 'vue-router'

const { CResourceList } = components

const router = useRouter()
const store = useProjectsStore()
const confirm = useConfirm()
const toast = useToast()
const { projects } = storeToRefs(store)

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
  { key: 'name', sortable: true, header: 'Name' },
  { key: 'mode', sortable: true, header: 'Mode' },
  { key: 'status', sortable: true, header: 'Status' },
  { key: 'version', sortable: false, header: 'Version' },
  { key: 'progress', sortable: false, header: 'Build progress' },
  {
    key: 'updatedAt',
    sortable: true,
    header: 'Last modified',
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const capitalize = s => (s ? s[0].toUpperCase() + s.slice(1) : s)

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
  toast.add({ severity: 'success', summary: 'Project created', detail: project.name, life: 2500 })
  router.push({ name: 'project.wizard', params: { projectId: project.id } })
}

const onRename = name => {
  if (renameTarget.value) store.updateProject(renameTarget.value.id, { name })
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
  store.updateProject(project.id, { status: project.status === 'archived' ? 'draft' : 'archived' })
}

const confirmDelete = project => {
  closeMenu()
  confirm.require({
    header: 'Delete project',
    message: `Delete "${project.name}"? This can't be undone.`,
    icon: 'pi pi-exclamation-triangle',
    rejectProps: { label: 'Cancel', severity: 'secondary', text: true },
    acceptProps: { label: 'Delete', severity: 'danger' },
    accept: () => {
      store.removeProject(project.id)
      toast.add({ severity: 'success', summary: 'Project deleted', detail: project.name, life: 2500 })
    },
  })
}

const actionItemsFor = project => [
  { label: 'Rename', icon: 'pi pi-pencil', command: () => openRename(project) },
  {
    label: project.status === 'archived' ? 'Unarchive' : 'Archive',
    icon: 'pi pi-inbox',
    command: () => toggleArchive(project),
  },
  { separator: true },
  { label: 'Delete', icon: 'pi pi-trash', class: 'text-red-500', command: () => confirmDelete(project) },
]
</script>
