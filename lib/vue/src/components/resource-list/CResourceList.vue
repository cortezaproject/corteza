<template>
  <Card
    :pt="{
      root: { class: 'overflow-hidden h-full flex flex-col min-w-0' },
      header: {
        class: 'flex flex-wrap items-center justify-between gap-3 p-3 w-full border-b shrink-0',
      },
      body: {
        class: 'p-0 flex flex-col flex-1 min-h-0 min-w-0',
      },
      content: {
        class: 'flex-1 overflow-auto min-h-0 min-w-0',
      },
    }"
  >
    <template v-if="$slots.header || !hideSearch" #header>
      <div class="flex-1">
        <slot name="header" />
      </div>
      <CInputSearch
        v-if="!hideSearch"
        :model-value="filter[queryField]"
        :placeholder="translations.searchPlaceholder || 'Search applications...'"
        size="small"
        class="flex-1 max-w-xl"
        @update:model-value="$emit('update:filter', { ...filter, [queryField]: $event })"
      />
    </template>

    <template #content>
      <div class="flex-1 overflow-auto min-h-0 min-w-0 flex flex-col h-full w-full">
        <DataTable
          v-model:selection="selected"
          :dataKey="primaryKey"
          :value="items"
          :loading="loading"
          :sortOrder="sorting.sortDesc ? 1 : -1"
          :sortField="sorting.sortBy"
          scrollable
          scrollHeight="flex"
          row-hover
          lazy
          resizableColumns
          columnResizeMode="expand"
          tableStyle="min-width: 50rem"
          :row-class="rowClass"
          :pt="{
            root: { class: 'flex-1 flex flex-col min-h-0 max-w-full' },
            tableContainer: { class: 'flex-1 overflow-auto max-w-full' },
            emptyMessageCell: { class: 'h-full' },
            footer: { class: 'p-0 border-0' },
            headerCell: { class: 'bg-highlight-emphasis' },
          }"
          @sort="$emit('sort', $event)"
          @row-click="$emit('row-click', $event)"
        >
          <template #empty>
            <div class="flex items-center justify-center p-4 text-muted">
              {{ translations.noItems || t('general.resourceList.noItems') }}
            </div>
          </template>

          <Column v-if="selectable" selectionMode="multiple" headerStyle="width: 3rem" />
          <Column
            v-for="field in fields"
            :key="field.key"
            :field="field.key"
            :header="field.header || field.label"
            :sortable="field.sortable"
            :class="field.class"
            :style="field.style"
            :pt="field.pt"
            :frozen="field.frozen"
            :alignFrozen="field.alignFrozen"
          >
            <template #body="slotProps">
              <slot :name="`body-${field.key}`" :data="slotProps.data" :field="field">
                {{ slotProps.data[field.key] }}
              </slot>
            </template>
          </Column>

          <template #footer>
            <div class="flex items-center flex-wrap gap-2 px-3 py-2">
              <div class="flex items-center text-sm">
                <span v-if="!hideTotal" class="whitespace-nowrap">
                  {{ getPagination }}
                </span>

                <Divider layout="vertical" />

                <div v-if="!hidePerPageOption" class="flex items-center gap-2 whitespace-nowrap">
                  <span>
                    {{ translations.recordsPerPage || 'Per Page' }}
                  </span>
                  <Select
                    :model-value="pagination.limit"
                    :options="perPageOptions"
                    size="small"
                    class="w-20"
                    @update:model-value="handlePerPageChange"
                  />
                </div>
              </div>

              <div class="flex items-center ml-auto gap-1">
                <Button
                  icon="pi pi-angle-double-left"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="!hasPrevPage"
                  @click="goToPage()"
                />
                <Button
                  icon="pi pi-angle-left"
                  :label="translations.prevPagination || 'Previous'"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="!hasPrevPage"
                  @click="goToPage('prevPage')"
                />
                <Button
                  :label="translations.nextPagination || 'Next'"
                  icon="pi pi-angle-right"
                  iconPos="right"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="!hasNextPage"
                  @click="goToPage('nextPage')"
                />
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </template>
  </Card>
</template>

<script setup>
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputSearch from '../input/CInputSearch.vue'

const emit = defineEmits(['search', 'sort', 'row-click', 'update:filter', 'page-change'])

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
  hideTotal: {
    type: Boolean,
    default: false,
  },
  hidePerPageOption: {
    type: Boolean,
    default: false,
  },
  perPageOptions: {
    type: Array,
    default: () => [10, 20, 50, 100],
  },
  translations: {
    type: Object,
    default: () => ({}),
  },
})

const { t } = useI18n()
const selected = ref([])

const hasPrevPage = computed(() => !!props.pagination.prevPage)
const hasNextPage = computed(() => !!props.pagination.nextPage)

const getPagination = computed(() => {
  let { total = 0, limit = 10, page = 1 } = props.pagination
  total = isNaN(total) ? 0 : total

  const from = (page - 1) * limit + 1
  const to = limit > 0 ? Math.min(page * limit, total) : total
  const data = total === 1 ? props.translations.resourceSingle : props.translations.resourcePlural

  if (total > limit && props.translations.showingPagination) {
    return t(props.translations.showingPagination, { from, to, count: total, data })
  }

  if (total <= limit && props.translations.singlePluralPagination) {
    return t(props.translations.singlePluralPagination, { count: total, data }, total)
  }

  return `${from} - ${to} of ${total}`
})

const goToPage = direction => {
  if (!direction) {
    // First page
    emit('page-change', { pageCursor: '', page: 1 })
  } else if (direction === 'prevPage') {
    emit('page-change', { pageCursor: props.pagination.prevPage, page: props.pagination.page - 1 })
  } else if (direction === 'nextPage') {
    emit('page-change', { pageCursor: props.pagination.nextPage, page: props.pagination.page + 1 })
  }
}

const handlePerPageChange = value => {
  emit('page-change', { pageCursor: '', page: 1, limit: value })
}

const rowClass = () => {
  return {
    'cursor-pointer': props.clickable,
  }
}
</script>

<style scoped>
:deep(.p-datatable-gridlines .p-datatable-paginator-bottom) {
  border-width: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :first-child) {
  border-left: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :last-child) {
  border-right: 0;
}

:deep(.p-datatable-gridlines .p-datatable-thead > tr > th) {
  border-top: 0;
}

:deep(.p-datatable-mask) {
  background: color-mix(in srgb, var(--p-content-background) 80%, transparent) !important;
}

:deep(.p-datatable-mask .p-icon-spin) {
  color: var(--p-primary-color);
}

:deep(.row-action-btn) {
  opacity: 0;
  transition: opacity 0.15s ease;
}

:deep(.p-datatable-tbody > tr:hover .row-action-btn) {
  opacity: 1;
}
</style>
