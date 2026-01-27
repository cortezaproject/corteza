<template>
  <Card
    :pt="{
      header: {
        class: 'flex flex-wrap items-center justify-between gap-3 p-3 w-full',
      },
      body: {
        class: 'p-0 overflow-auto h-full',
      },
      content: {
        class: 'overflow-auto h-full',
      },
    }"
    class="overflow-hidden"
  >
    <template v-if="$slots.header || !hideSearch" #header>
      <div class="flex-1">
        <slot name="header" />
      </div>
      <CInputSearch
        v-if="!hideSearch"
        :model-value="filter[queryField]"
        :placeholder="translations.searchPlaceholder || 'Search applications...'"
        submittable
        class="flex-1 max-w-xl"
        @update:model-value="$emit('update:filter', { ...filter, [queryField]: $event })"
        @search="emit('search', $event)"
      />
    </template>

    <template #content>
      <DataTable
        v-model:selection="selected"
        :dataKey="primaryKey"
        :value="items"
        :loading="loading"
        :sortOrder="sorting.sortDesc ? 1 : -1"
        :sortField="sorting.sortBy"
        scrollable
        scrollHeight="flex"
        :paginator="!hidePagination"
        row-hover
        lazy
        resizableColumns
        columnResizeMode="fit"
        :row-class="rowClass"
        :rows="pagination.limit"
        :rowsPerPageOptions="[5, 10, 20, 50]"
        @sort="$emit('sort', $event)"
        @row-click="$emit('row-click', $event)"
      >
        <Column v-if="selectable" selectionMode="multiple" headerStyle="width: 3rem" />
        <Column
          v-for="field in fields"
          :key="field.key"
          :field="field.key"
          :header="field.header || field.label"
          :sortable="field.sortable"
          :class="field.class"
          :pt="field.pt"
        >
          <template #body="slotProps">
            <slot :name="`body-${field.key}`" :data="slotProps.data" :field="field">
              {{ slotProps.data[field.key] }}
            </slot>
          </template>
        </Column>

        <!-- Custom paginator template (uncomment when ready to implement)
        <template #paginatorcontainer="{ first, last, page, pageCount, prevPageCallback, nextPageCallback, totalRecords }">
          <div class="flex items-center flex-wrap gap-4">
            <Button icon="pi pi-chevron-left" rounded text @click="prevPageCallback" :disabled="page === 0" />
            <div class="text-color font-medium">
              <span class="hidden sm:block">Showing {{ first }} to {{ last }} of {{ totalRecords }}</span>
              <span class="block sm:hidden">Page {{ page + 1 }} of {{ pageCount }}</span>
            </div>
            <Button icon="pi pi-chevron-right" rounded text @click="nextPageCallback" :disabled="page === pageCount - 1" />
          </div>
        </template>
        -->
      </DataTable>
    </template>
  </Card>
</template>

<script setup>
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { ref } from 'vue'
import CInputSearch from '../input/CInputSearch.vue'

const emit = defineEmits(['search', 'sort', 'row-click'])

const props = defineProps({
  primaryKey: {
    type: String,
    required: true,
  },
  selectable: {
    type: Boolean,
    default: false,
  },
  fields: {
    type: Array,
    default: () => [],
  },
  items: {
    type: Array,
    default: () => [],
  },
  filter: {
    type: Object,
    default: () => ({}),
  },
  clickable: {
    type: Boolean,
    default: false,
  },
  sorting: {
    type: Object,
    default: () => ({}),
  },
  pagination: {
    type: Object,
    default: () => ({}),
  },
  loading: {
    type: Boolean,
    default: false,
  },
  queryField: {
    type: String,
    default: 'query',
  },
  hideSearch: {
    type: Boolean,
    default: false,
  },
  hidePagination: {
    type: Boolean,
    default: false,
  },

  translations: {
    type: Object,
    default: () => ({}),
  },
})

const selected = ref([])

const rowClass = () => {
  return {
    'cursor-pointer': props.clickable,
  }
}
</script>
