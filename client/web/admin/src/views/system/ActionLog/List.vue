<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.actionlog.list.title', 'Action Log') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="actionID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.actionlog.list.filter.search', 'Search'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.actionlog.list.title'),
        resourcePlural: $t('system.actionlog.list.title'),
      }"
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-surface-500">
            {{ $t('system.actionlog.list.filter.from', 'Starting from') }}
          </label>
          <DatePicker
            v-model="dateFrom"
            showTime
            hourFormat="24"
            dateFormat="yy-mm-dd"
            size="small"
            @update:modelValue="filterList"
          />
          <label class="text-sm font-medium text-surface-500">
            {{ $t('system.actionlog.list.filter.to', 'Ending at') }}
          </label>
          <DatePicker
            v-model="dateTo"
            showTime
            hourFormat="24"
            dateFormat="yy-mm-dd"
            size="small"
            @update:modelValue="filterList"
          />
        </div>
      </template>

      <template #body-timestamp="{ data }">
        {{ locFullDateTime(data.timestamp) }}
      </template>

      <template #body-severity="{ data }">
        <Tag :value="data.severity || 'info'" :severity="severityMap[data.severity] || 'info'" />
      </template>
    </CResourceList>
  </div>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

const dateFrom = ref(null)
const dateTo = ref(null)

const severityMap = {
  error: 'danger',
  warn: 'warn',
  warning: 'warn',
  info: 'info',
  debug: 'secondary',
  notice: 'success',
}

const fields = [
  {
    key: 'timestamp',
    sortable: true,
    header: t('system.actionlog.list.columns.timestamp', 'Timestamp'),
  },
  {
    key: 'action',
    sortable: true,
    header: t('system.actionlog.list.columns.action', 'Action'),
  },
  {
    key: 'actorID',
    sortable: false,
    header: t('system.actionlog.list.columns.actor', 'Actor'),
  },
  {
    key: 'resource',
    sortable: false,
    header: t('system.actionlog.list.columns.resource', 'Resource'),
  },
  {
    key: 'severity',
    sortable: true,
    header: t('system.actionlog.list.columns.severity', 'Severity'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(
    params =>
      $SystemAPI.actionlogListCancellable({
        ...params,
        from: dateFrom.value ? dateFrom.value.toISOString() : undefined,
        to: dateTo.value ? dateTo.value.toISOString() : undefined,
      }),
    {
      filter: { query: '' },
      sorting: { sortBy: 'timestamp', sortDesc: true },
      pagination: { limit: 50 },
    },
  )
</script>
