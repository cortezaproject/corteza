<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.sessions.list.title', 'Sessions') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="sessionID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('automation.sessions.list.title'),
        resourcePlural: $t('automation.sessions.list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="
        ({ data }) =>
          $router.push({ name: 'automation.sessions.view', params: { sessionID: data.sessionID } })
      "
      @page-change="handlePageChange"
    >
      <template #body-status="{ data }">
        <Tag :value="data.status" :severity="statusSeverity(data.status)" />
      </template>

      <template #body-createdAt="{ data }">
        {{ locFullDateTime(data.createdAt) }}
      </template>
    </CResourceList>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()

const $AutomationAPI = inject('$AutomationAPI')

function statusSeverity(status) {
  switch (status) {
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    case 'canceled':
      return 'secondary'
    case 'started':
    case 'pending':
      return 'info'
    default:
      return 'secondary'
  }
}

const fields = [
  {
    key: 'sessionID',
    sortable: false,
    header: t('automation.sessions.list.columns.sessionID', 'Session ID'),
  },
  {
    key: 'workflowID',
    sortable: false,
    header: t('automation.sessions.list.columns.workflowID', 'Workflow ID'),
  },
  {
    key: 'status',
    sortable: false,
    header: t('automation.sessions.list.columns.status', 'Status'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('automation.sessions.list.columns.createdAt', 'Started'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange } =
  useResourceList(params => $AutomationAPI.sessionListCancellable({ ...params }), {
    filter: { query: '' },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 50 },
  })
</script>
