<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('general.workflow-list') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="workflowID"
      :fields="workflowFields"
      :items="workflowList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('general.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.workflow.single'),
        resourcePlural: $t('general.workflow.plural'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <CRouterLinkButton
            v-if="canCreate"
            data-test-id="button-create-workflow"
            :to="{ name: 'workflow.create' }"
            :label="$t('general.new-workflow')"
            icon="pi pi-plus"
            size="small"
          />
          <Import
            v-if="canCreate"
            data-test-id="button-import-workflow"
            :disabled="importProcessing"
            @import="importJSON"
          />
          <Export :workflows="workflowIDs" size="small" severity="secondary" />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::automation:workflow/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col gap-1">
          <div class="flex items-center gap-2 flex-wrap">
            <span>{{ data.meta?.name || data.handle || '-' }}</span>
            <Tag
              v-if="data.meta?.subWorkflow"
              :value="$t('general.subworkflow')"
              severity="info"
              class="text-xs"
            />
          </div>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
          <div
            v-for="group in getWorkflowLabels(data)"
            :key="'group-' + group.namespaceID"
            class="flex items-center flex-wrap gap-1"
          >
            <Tag
              v-tooltip.top="$t('general.filter.namespace.label')"
              :value="group.namespaceName"
              severity="primary"
              class="text-xs"
            />
            <Tag
              v-for="mod in group.modules"
              :key="mod.id"
              v-tooltip.top="$t('general.filter.module.label')"
              :value="mod.name"
              severity="secondary"
              class="text-xs"
            />
          </div>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('general.enabled') : $t('general.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
          rounded
        />
      </template>

      <template #body-steps="{ data }">
        {{ (data.steps || []).length }}
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
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
      <div class="flex flex-col gap-4 p-2 w-72">
        <!-- SubWorkflow filter -->
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">{{ $t('general.subworkflows') }}</span>
          <div v-for="opt in radioOptions" :key="'sw-' + opt.value" class="flex items-center gap-2">
            <RadioButton
              v-model="filter.subWorkflow"
              :inputId="'sw' + opt.value"
              :value="opt.value"
              @change="filterList"
            />
            <label :for="'sw' + opt.value" class="text-sm cursor-pointer">{{ opt.label }}</label>
          </div>
        </div>

        <!-- Disabled filter -->
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">{{ $t('general.disabled') }}</span>
          <div
            v-for="opt in radioOptions"
            :key="'dis-' + opt.value"
            class="flex items-center gap-2"
          >
            <RadioButton
              v-model="filter.disabled"
              :inputId="'dis' + opt.value"
              :value="opt.value"
              @change="filterList"
            />
            <label :for="'dis' + opt.value" class="text-sm cursor-pointer">{{ opt.label }}</label>
          </div>
        </div>

        <!-- Deleted filter -->
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">{{ $t('general.deleted') }}</span>
          <div
            v-for="opt in radioOptions"
            :key="'del-' + opt.value"
            class="flex items-center gap-2"
          >
            <RadioButton
              v-model="filter.deleted"
              :inputId="'del' + opt.value"
              :value="opt.value"
              @change="filterList"
            />
            <label :for="'del' + opt.value" class="text-sm cursor-pointer">{{ opt.label }}</label>
          </div>
        </div>

        <Divider class="my-0" />

        <!-- Namespace / module label filter -->
        <NamespaceModuleSelector
          :namespace-labels="selectedNamespaceLabels"
          :module-labels="selectedModuleLabels"
          @change="handleLabelFilterChange"
        />
      </div>
    </Popover>
  </div>
</template>

<script setup>
import {
  components,
  filters,
  useConfirmDelete,
  useRBACStore,
  useResourceList,
} from '@planetcrust/human-vue'
import { inject, ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { saveAs } from 'file-saver'
import { useLabelsStore } from '@/stores/labels'
import { useWorkflowStore } from '@/stores/workflow'
import Import from '@/components/Import.vue'
import Export from '@/components/Export.vue'
import NamespaceModuleSelector from '@/components/NamespaceModuleSelector.vue'

const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const toast = useToast()
const $AutomationAPI = inject('$AutomationAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $Auth = inject('$Auth')
const labelsStore = useLabelsStore()
const workflowStore = useWorkflowStore()

// Dialog / popover visibility
const filterMenu = ref()
const importProcessing = ref(false)

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Label filter state
const selectedNamespaceLabels = ref([])
const selectedModuleLabels = ref([])
const labelsFilter = ref([])

const radioOptions = [
  { label: t('general.without'), value: '0' },
  { label: t('general.including'), value: '1' },
  { label: t('general.only'), value: '2' },
]

const userID = computed(() => $Auth?.user?.userID)

const resourceListRef = ref()

// RBAC
const rbacStore = useRBACStore()
const canCreate = computed(() => rbacStore.can('automation/', 'workflow.create'))
const canGrant = computed(() => rbacStore.can('automation/', 'grant'))

// Column definitions
const workflowFields = [
  {
    key: 'name',
    sortable: true,
    header: t('general.columns.name'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('general.columns.enabled'),
    class: 'text-center w-28',
  },
  {
    key: 'steps',
    sortable: false,
    header: t('general.columns.steps'),
    class: 'text-center w-20',
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('general.columns.changedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items: workflowList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(
  params => $AutomationAPI.workflowListCancellable({ ...params, labels: labelsFilter.value }),
  {
    filter: { query: '', subWorkflow: '1', disabled: '1', deleted: '0' },
    sorting: { sortBy: 'name', sortDesc: false },
    pagination: { limit: 50 },
  },
)

const workflowIDs = computed(() => workflowList.value.map(w => w.workflowID))

// Resolve namespace/module names whenever the list changes
watch(workflowList, workflows => {
  if (!workflows?.length) return

  const namespaceIDs = new Set()
  const modules = []

  workflows.forEach(wf => {
    if (!wf.labels) return
    const ns = [wf.labels.ref_namespace].flat().filter(Boolean)
    const mod = [wf.labels.ref_module].flat().filter(Boolean)

    ns.forEach(l => {
      const id = l.split('/')[1]
      if (id) namespaceIDs.add(id)
    })
    mod.forEach(l => {
      const [, nsID, modID] = l.split('/')
      if (nsID) namespaceIDs.add(nsID)
      if (modID) modules.push({ moduleID: modID, namespaceID: nsID })
    })
  })

  if (namespaceIDs.size) {
    labelsStore.resolveMultipleNamespaces({ namespaceIDs: [...namespaceIDs], api: $ComposeAPI })
  }
  if (modules.length) {
    labelsStore.resolveMultipleModules({ modules, api: $ComposeAPI })
  }
})

// Methods
function getWorkflowLabels(workflow) {
  if (!workflow.labels) return []

  const nsIDs = []
  const modsByNs = {}

  ;[workflow.labels.ref_namespace]
    .flat()
    .filter(Boolean)
    .forEach(l => {
      const id = l.split('/')[1]
      if (id && !nsIDs.includes(id)) nsIDs.push(id)
    })
  ;[workflow.labels.ref_module]
    .flat()
    .filter(Boolean)
    .forEach(l => {
      const [, nsID, modID] = l.split('/')
      if (!nsID || !modID) return
      if (!nsIDs.includes(nsID)) nsIDs.push(nsID)
      if (!modsByNs[nsID]) modsByNs[nsID] = []
      modsByNs[nsID].push({ id: modID, name: labelsStore.getModule(modID) || modID })
    })

  return nsIDs.map(id => ({
    namespaceID: id,
    namespaceName: labelsStore.getNamespace(id) || id,
    modules: modsByNs[id] || [],
  }))
}

function handleLabelFilterChange({ namespaceLabels, moduleLabels }) {
  selectedNamespaceLabels.value = namespaceLabels || []
  selectedModuleLabels.value = moduleLabels || []

  const labels = []
  if (selectedNamespaceLabels.value.length) {
    labels.push(`ref_namespace=${JSON.stringify(selectedNamespaceLabels.value)}`)
  }
  if (selectedModuleLabels.value.length) {
    labels.push(`ref_module=${JSON.stringify(selectedModuleLabels.value)}`)
  }
  labelsFilter.value = labels
  filterList()
}

async function handleStatusChange(workflow) {
  const enabled = !workflow.enabled
  const key = enabled ? 'enable' : 'disable'
  try {
    const w = await $AutomationAPI.workflowRead({ workflowID: workflow.workflowID })
    const updated = await $AutomationAPI.workflowUpdate({ ...w, enabled })
    workflowStore.updateInList(updated)
    toast.add({ severity: 'success', summary: t(`notification.list.${key}.success`), life: 3000 })
    filterList()
  } catch {
    toast.add({ severity: 'error', summary: t(`notification.list.${key}.failed`), life: 5000 })
  }
}

async function handleExportWorkflow(workflow) {
  try {
    const { set: tSet = [] } = await $AutomationAPI.triggerList({
      workflowID: [workflow.workflowID],
      disabled: 1,
    })
    const triggers = {}
    tSet.forEach(tr => {
      if (!triggers[tr.workflowID]) triggers[tr.workflowID] = []
      triggers[tr.workflowID].push({
        resourceType: tr.resourceType,
        eventType: tr.eventType,
        constraints: tr.constraints,
        enabled: tr.enabled,
        stepID: tr.stepID,
        meta: tr.meta,
      })
    })

    const { set: wSet = [] } = await $AutomationAPI.workflowList({
      workflowID: [workflow.workflowID],
      disabled: 1,
      subWorkflow: 1,
    })
    const workflows = wSet.map(w => ({
      handle: w.handle,
      enabled: w.enabled,
      meta: w.meta,
      keepSessions: w.keepSessions,
      steps: w.steps,
      paths: w.paths,
      triggers: triggers[w.workflowID],
    }))

    const blob = new Blob([JSON.stringify({ workflows }, null, 2)], { type: 'application/json' })
    const filename = (workflow.meta?.name || workflow.handle || 'workflow').replace(
      /[/\\?%*:|"<>]/g,
      '',
    )
    saveAs(blob, `${filename}.json`)
  } catch {
    toast.add({ severity: 'error', summary: t('notification.failed-fetch-workflows'), life: 5000 })
  }
}

async function importJSON(workflows = []) {
  importProcessing.value = true
  const skipped = []

  await Promise.all(
    workflows.map(({ triggers = [], ...wf }) =>
      $AutomationAPI
        .workflowCreate({ ownedBy: userID.value, runAs: '0', ...wf })
        .then(({ workflowID }) =>
          Promise.all(
            triggers.map(tr =>
              $AutomationAPI.triggerCreate({
                ...tr,
                workflowID,
                workflowStepID: tr.stepID,
                ownedBy: userID.value,
              }),
            ),
          ),
        )
        .catch(({ message }) => {
          if (wf.handle) skipped.push(`${wf.handle}${message ? ' - ' + message : ''}`)
        }),
    ),
  )

  if (skipped.length) {
    toast.add({
      severity: 'warn',
      summary: t('notification.import.skipped-workflows'),
      detail: skipped.join('; '),
      life: 5000,
    })
  } else {
    toast.add({
      severity: 'success',
      summary: t('notification.import.imported-workflows'),
      life: 3000,
    })
  }

  importProcessing.value = false
  filterList()
}

function handleRowClick({ data }) {
  router.push({
    name: 'workflow.edit',
    params: { workflowID: data.workflowID },
  })
}

function getActionsMenuItems(workflow) {
  const items = []

  items.push({
    label: t('general.label.edit'),
    icon: 'pi pi-pencil',
    route: {
      name: 'workflow.edit',
      params: { workflowID: workflow.workflowID },
    },
  })

  if (!workflow.deletedAt) {
    items.push({
      label: workflow.enabled ? t('general.disable') : t('general.enable'),
      icon: workflow.enabled ? 'pi pi-power-off' : 'pi pi-check-circle',
      command: () => handleStatusChange(workflow),
    })
  }

  items.push({
    label: t('general.export'),
    icon: 'pi pi-download',
    command: () => handleExportWorkflow(workflow),
  })

  items.push({ separator: true })

  if (workflow.deletedAt) {
    items.push({
      label: t('general.undelete'),
      icon: 'pi pi-undo',
      command: () => handleUndelete(workflow),
    })
  } else {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(workflow),
    })
  }

  return items
}

function onConfirmDelete(workflow) {
  confirmDelete({
    message: t('notification.delete.confirm'),
    header: workflow.meta?.name || workflow.handle || t('general.label.delete'),
    onConfirm: () => handleDelete(workflow),
  })
}

async function handleDelete(workflow) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $AutomationAPI.workflowDelete({
      workflowID: workflow.workflowID,
    })
    workflowStore.removeFromList(workflow.workflowID)
    toast.add({ severity: 'success', summary: t('notification.delete.success'), life: 3000 })
    filterList()
  } catch (e) {
    console.error('Failed to delete workflow:', e)
    toast.add({ severity: 'error', summary: t('notification.delete.failed'), life: 5000 })
  }
}

async function handleUndelete(workflow) {
  resourceListRef.value.hideActionsMenu()
  try {
    const restored = await $AutomationAPI.workflowUndelete({
      workflowID: workflow.workflowID,
    })
    workflowStore.updateInList(restored)
    toast.add({ severity: 'success', summary: t('notification.undelete.success'), life: 3000 })
    filterList()
  } catch (e) {
    console.error('Failed to restore workflow:', e)
    toast.add({ severity: 'error', summary: t('notification.undelete.failed'), life: 5000 })
  }
}
</script>
