<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.labels.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="_rowKey"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.labels.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.labels.list.resourceSingle'),
        resourcePlural: $t('system.labels.list.title'),
      }"
      class="h-full"
      clickable
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @page-change="handlePageChange"
      @row-click="onRowClick"
    >
      <template #body-resourceCount="{ data }">
        <Tag :value="String(data.resourceCount || 0)" severity="secondary" rounded />
      </template>
    </CResourceList>
  </div>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
  components,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'

const { CResourceList } = components

const { t } = useI18n()
const router = useRouter()

const $SystemAPI = inject('$SystemAPI')

const resourceListRef = ref()

const fields = [
  {
    key: 'name',
    sortable: false,
    header: t('system.labels.list.columns.name'),
  },
  {
    key: 'resourceCount',
    sortable: false,
    header: t('system.labels.list.columns.resources'),
  },
]

function labelListCancellable(params) {
  const { response, cancel } = $SystemAPI.labelListCancellable({
    ...params,
    name: params.query || undefined,
  })

  return {
    response: async () => {
      const result = await response()
      const set = Array.isArray(result) ? result : result.set || []

      // Add a synthetic row key since labels have no dedicated ID
      const enriched = set.map((item, idx) => ({
        ...item,
        _rowKey: `${item.kind || ''}_${item.name || ''}_${idx}`,
      }))

      // Return in CResourceList expected shape
      if (Array.isArray(result)) {
        return enriched
      }

      return {
        ...result,
        set: enriched,
      }
    },
    cancel,
  }
}

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange } =
  useResourceList(params => labelListCancellable(params), {
    filter: { query: '' },
    sorting: {},
    pagination: { limit: 100 },
  })

function onRowClick({ data }) {
  if (data?.name) {
    router.push({
      name: 'system.labels.edit',
      params: { labelID: encodeURIComponent(data.name) },
    })
  }
}
</script>
