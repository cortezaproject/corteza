<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.sessions.list.title') }}</span>
  </Teleport>

  <div class="p-4 h-full overflow-hidden min-w-0 flex flex-col">
    <Card
      :pt="{
        root: { class: 'flex-1 flex flex-col min-h-0 overflow-hidden' },
        body: { class: 'p-0 flex flex-col flex-1 min-h-0' },
        content: { class: 'p-0 flex flex-col flex-1 min-h-0' },
      }"
    >
      <template #content>
    <Tabs v-model:value="activeTab" class="flex-1 flex flex-col min-h-0 overflow-hidden">
      <TabList class="rounded-t-lg">
        <Tab value="taq">{{ $t('automation.sessions.list.tabs.taq') }}</Tab>
        <Tab value="workflow">{{ $t('automation.sessions.list.tabs.workflow') }}</Tab>

      </TabList>

      <TabPanels class="flex-1 min-h-0 overflow-hidden p-0">
        <!-- TAQ Executions -->
        <TabPanel value="taq" class="h-full p-0 pt-3">
          <CResourceList
            ref="taqListRef"
            primary-key="executionID"
            :fields="taqFields"
            :items="taqItems"
            :filter="taqFilter"
            :sorting="{}"
            :pagination="{}"
            :loading="taqLoading"
            :translations="{
              resourceSingle: $t('automation.sessions.list.tabs.taq'),
              resourcePlural: $t('automation.sessions.list.tabs.taq'),
            }"
            hide-search
            hide-pagination
            class="h-full"
          >
            <template #body-status="{ data }">
              <Tag :value="data.status" :severity="taqStatusSeverity(data.status)" />
            </template>

            <template #body-startedAt="{ data }">
              {{ locFullDateTime(data.startedAt) }}
            </template>

            <template #filter>
              <Button
                icon="pi pi-filter"
                severity="secondary"
                size="small"
                text
                @click="toggleTaqFilterMenu"
              />
            </template>
          </CResourceList>

          <Popover ref="taqFilterMenu">
            <div class="flex flex-col gap-4 p-2 w-64">
              <div class="flex flex-col gap-2">
                <span class="font-medium text-sm text-primary">
                  {{ $t('automation.sessions.list.columns.status') }}
                </span>
                <div v-for="opt in taqStatusOptions" :key="opt.value" class="flex items-center gap-2">
                  <RadioButton v-model="taqFilter.status" :inputId="`tst-${opt.value}`" :value="opt.value" />
                  <label :for="`tst-${opt.value}`" class="text-sm cursor-pointer">{{ opt.label }}</label>
                </div>
              </div>

              <div class="flex flex-col gap-2">
                <label class="font-medium text-sm text-primary" for="taq-filter-automationID">
                  {{ $t('automation.sessions.list.taqColumns.automationID') }}
                </label>
                <InputText
                  id="taq-filter-automationID"
                  v-model="taqFilter.automationID"
                  size="small"
                />
              </div>
            </div>
          </Popover>
        </TabPanel>

        <!-- Workflow Sessions -->
        <TabPanel value="workflow" class="h-full p-0 pt-3">
          <CResourceList
            primary-key="sessionID"
            :fields="workflowFields"
            :items="workflowItems"
            :filter="workflowFilter"
            :sorting="workflowSorting"
            :pagination="workflowPagination"
            :loading="workflowLoading"
            :translations="{
              showingPagination: 'general.resourceList.pagination.showing',
              singlePluralPagination: 'general.resourceList.pagination.single',
              prevPagination: $t('general.resourceList.pagination.prev'),
              nextPagination: $t('general.resourceList.pagination.next'),
              recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
              resourceSingle: $t('automation.sessions.list.single'),
              resourcePlural: $t('automation.sessions.list.title'),
            }"
            clickable
            hide-search
            class="h-full"
            @sort="handleWorkflowSort"
            @row-click="
              ({ data }) =>
                $router.push({ name: 'automation.sessions.view', params: { sessionID: data.sessionID } })
            "
            @page-change="handleWorkflowPageChange"
          >
            <template #body-status="{ data }">
              <Tag :value="data.status" :severity="workflowStatusSeverity(data.status)" />
            </template>

            <template #body-createdAt="{ data }">
              {{ locFullDateTime(data.createdAt) }}
            </template>

            <template #filter>
              <Button
                icon="pi pi-filter"
                severity="secondary"
                size="small"
                text
                @click="toggleWorkflowFilterMenu"
              />
            </template>
          </CResourceList>

          <Popover ref="workflowFilterMenu">
            <div class="flex flex-col gap-4 p-2 w-64">
              <div class="flex flex-col gap-2">
                <span class="font-medium text-sm text-primary">
                  {{ $t('automation.sessions.list.columns.status') }}
                </span>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-all" value="all" />
                  <label for="st-all" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.all.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-0" value="0" />
                  <label for="st-0" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.started.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-1" value="1" />
                  <label for="st-1" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.prompted.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-2" value="2" />
                  <label for="st-2" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.suspended.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-3" value="3" />
                  <label for="st-3" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.failed.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-4" value="4" />
                  <label for="st-4" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.completed.label') }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioButton v-model="selectedWorkflowStatus" inputId="st-5" value="5" />
                  <label for="st-5" class="text-sm cursor-pointer">
                    {{ $t('automation.sessions.list.filterForm.canceled.label') }}
                  </label>
                </div>
              </div>

              <div class="flex flex-col gap-2">
                <label class="font-medium text-sm text-primary" for="filter-sessionID">
                  {{ $t('automation.sessions.list.columns.sessionID') }}
                </label>
                <InputText
                  id="filter-sessionID"
                  v-model="workflowFilter.sessionID"
                  size="small"
                  @input="workflowFilterList"
                />
              </div>

              <div class="flex flex-col gap-2">
                <label class="font-medium text-sm text-primary" for="filter-workflowID">
                  {{ $t('automation.sessions.list.columns.workflowID') }}
                </label>
                <InputText
                  id="filter-workflowID"
                  v-model="workflowFilter.workflowID"
                  size="small"
                  @input="workflowFilterList"
                />
              </div>
            </div>
          </Popover>
        </TabPanel>
      </TabPanels>
    </Tabs>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const $AutomationAPI = inject('$AutomationAPI')

const activeTab = ref('taq')

// ── TAQ Executions ──────────────────────────────────────────────────────────

const taqListRef = ref()
const taqFilterMenu = ref()

function toggleTaqFilterMenu(event) {
  taqFilterMenu.value.toggle(event)
}

const taqStatusOptions = [
  { value: '', label: t('automation.sessions.list.filterForm.all.label') },
  { value: 'created', label: t('automation.sessions.list.taqStatus.created') },
  { value: 'running', label: t('automation.sessions.list.taqStatus.running') },
  { value: 'paused', label: t('automation.sessions.list.taqStatus.paused') },
  { value: 'completed', label: t('automation.sessions.list.taqStatus.completed') },
  { value: 'failed', label: t('automation.sessions.list.taqStatus.failed') },
  { value: 'cancelled', label: t('automation.sessions.list.taqStatus.cancelled') },
]

function taqStatusSeverity(status) {
  switch (status) {
    case 'completed': return 'success'
    case 'failed': return 'danger'
    case 'cancelled': return 'secondary'
    case 'running': return 'info'
    case 'paused': return 'warn'
    default: return 'secondary'
  }
}

const taqFields = [
  { key: 'executionID', sortable: false, header: t('automation.sessions.list.taqColumns.executionID') },
  { key: 'executableID', sortable: false, header: t('automation.sessions.list.taqColumns.automationID') },
  { key: 'status', sortable: false, header: t('automation.sessions.list.columns.status') },
  { key: 'startedAt', sortable: false, header: t('automation.sessions.list.taqColumns.startedAt'), class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'duration', sortable: false, header: t('automation.sessions.list.taqColumns.duration'), class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
]

const taqFilter = ref({ automationID: '', status: '' })
const taqItems = ref([])
const taqLoading = ref(false)

async function fetchTaqItems() {
  taqLoading.value = true
  try {
    const params = {}
    if (taqFilter.value.automationID) params.automationID = taqFilter.value.automationID
    if (taqFilter.value.status) params.status = taqFilter.value.status
    const result = await $AutomationAPI.ngAutomationAllExecutions(params)
    taqItems.value = result?.set || (Array.isArray(result) ? result : [])
  } catch (e) {
    console.error('Failed to fetch TAQ executions:', e)
    taqItems.value = []
  } finally {
    taqLoading.value = false
  }
}

watch(taqFilter, () => fetchTaqItems(), { deep: true })
watch(activeTab, (tab) => { if (tab === 'taq') fetchTaqItems() }, { immediate: true })

// ── Workflow Sessions ────────────────────────────────────────────────────────

const workflowFilterMenu = ref()

function toggleWorkflowFilterMenu(event) {
  workflowFilterMenu.value.toggle(event)
}

function workflowStatusSeverity(status) {
  switch (status) {
    case 'completed': return 'success'
    case 'failed': return 'danger'
    case 'canceled': return 'secondary'
    case 'started':
    case 'pending': return 'info'
    default: return 'secondary'
  }
}

const workflowFields = [
  { key: 'sessionID', sortable: false, header: t('automation.sessions.list.columns.sessionID') },
  { key: 'workflowID', sortable: false, header: t('automation.sessions.list.columns.workflowID') },
  { key: 'eventType', sortable: false, header: t('automation.sessions.list.columns.eventType') },
  { key: 'status', sortable: false, header: t('automation.sessions.list.columns.status') },
  {
    key: 'createdAt',
    sortable: true,
    header: t('automation.sessions.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const {
  items: workflowItems,
  loading: workflowLoading,
  filter: workflowFilter,
  sorting: workflowSorting,
  pagination: workflowPagination,
  handleSort: handleWorkflowSort,
  handlePageChange: handleWorkflowPageChange,
  filterList: workflowFilterList,
} = useResourceList(params => $AutomationAPI.sessionListCancellable({ ...params }), {
  filter: { sessionID: null, workflowID: null, status: null },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
  immediate: false,
})

watch(activeTab, (tab) => { if (tab === 'workflow') workflowFilterList() })

const selectedWorkflowStatus = computed({
  get: () => (workflowFilter.status === null ? 'all' : String(workflowFilter.status)),
  set: (val) => {
    workflowFilter.status = val === 'all' ? null : Number(val)
    workflowFilterList()
  },
})
</script>
