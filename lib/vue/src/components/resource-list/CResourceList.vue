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
    <template v-if="$slots.header || $slots.filter || !hideSearch" #header>
      <div class="flex-1 min-w-0">
        <slot name="header" />
      </div>
      <div v-if="!hideSearch || $slots.filter" class="flex-1 flex items-center justify-end gap-2">
        <div class="flex min-w-0">
          <slot v-if="$slots.filter" name="filter" />
        </div>
        <CInputSearch
          v-if="!hideSearch"
          :model-value="filter[queryField]"
          :placeholder="translations.searchPlaceholder"
          size="small"
          class="flex-1 min-w-0 max-w-xl"
          @update:model-value="$emit('update:filter', { ...filter, [queryField]: $event })"
        />
      </div>
    </template>

    <template #content>
      <div class="flex-1 overflow-auto min-h-0 min-w-0 flex flex-col h-full w-full">
        <DataTable
          v-model:selection="selected"
          v-model:expandedRows="expandedRows"
          :dataKey="primaryKey"
          :value="items"
          :loading="loading"
          :sortOrder="sorting.sortDesc ? -1 : 1"
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
          }"
          @sort="$emit('sort', $event)"
          @row-click="$emit('row-click', $event)"
        >
          <template #empty>
            <div class="flex items-center justify-center p-4 text-muted-color">
              {{ translations.noItems || t('general.resourceList.noItems') }}
            </div>
          </template>

          <Column v-if="selectable" selectionMode="multiple" headerStyle="width: 3rem" />
          <Column v-if="expandable" :expander="true" style="width: 3rem" />
          <Column
            v-for="field in computedFields"
            :key="field.key"
            :field="field.key"
            :header="field.header || field.label"
            :sortable="field.sortable"
            :class="field.class"
            :style="field.style"
            :pt="{
              ...field.pt,
              headerCell: {
                ...field.pt?.headerCell,
                class: [
                  'bg-emphasis text-muted-color font-semibold uppercase text-sm tracking-wide',
                  field.pt?.headerCell?.class,
                ],
              },
            }"
            :frozen="field.frozen"
            :alignFrozen="field.alignFrozen"
          >
            <template #body="slotProps">
              <slot :name="`body-${field.key}`" :data="slotProps.data" :field="field">
                {{ slotProps.data[field.key] }}
              </slot>
            </template>
          </Column>

          <!-- Built-in actions column when actionItems prop is provided -->
          <Column
            v-if="actionItems"
            key="__actions"
            header=""
            class="text-right w-12"
            frozen
            alignFrozen="right"
            :pt="{
              headerCell: { class: 'border-l-0' },
              bodyCell: { class: 'px-2 py-1 border-l-0' },
            }"
          >
            <template #body="slotProps">
              <Button
                v-if="actionItems(slotProps.data).length"
                icon="pi pi-ellipsis-v"
                text
                severity="secondary"
                size="small"
                class="row-action-btn w-full"
                @click.stop="showActionsMenu($event, slotProps.data)"
              />
            </template>
          </Column>

          <template v-if="expandable" #expansion="slotProps">
            <slot name="expansion" :data="slotProps.data" />
          </template>

          <template v-if="!hidePagination || $slots.footer" #footer>
            <slot name="footer" />
            <div v-if="!hidePagination" class="flex items-center flex-wrap gap-2 px-3 py-2">
              <div class="flex items-center text-sm">
                <span v-if="!hideTotal" class="whitespace-nowrap">
                  {{ getPagination }}
                </span>

                <Divider layout="vertical" />

                <div v-if="!hidePerPageOption" class="flex items-center gap-2 whitespace-nowrap">
                  <span>
                    {{ translations.recordsPerPage }}
                  </span>
                  <Select
                    :model-value="pagination.limit"
                    :options="perPageOptions"
                    size="small"
                    class="w-24"
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
                  :label="translations.prevPagination"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="!hasPrevPage"
                  @click="goToPage('prevPage')"
                />
                <Button
                  :label="translations.nextPagination"
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

      <!-- One popup shared by every row's actions button -->
      <Menu v-if="actionItems" ref="actionsMenuRef" :model="currentMenuItems" popup>
        <template #item="{ item, props: menuProps }">
          <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
            <a v-ripple :href="href" v-bind="menuProps.action" @click="navigate">
              <span :class="item.icon" />
              <span class="ml-2">{{ item.label }}</span>
            </a>
          </router-link>
          <a v-else v-ripple v-bind="menuProps.action" :class="item.class">
            <span :class="item.icon" />
            <span class="ml-2">{{ item.label }}</span>
          </a>
        </template>
      </Menu>
    </template>
  </Card>
</template>

<script setup>
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
  expandable: {
    type: Boolean,
    default: false,
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
  /**
   * Function that receives row data and returns an array of PrimeVue MenuItem objects.
   * When provided, the component auto-adds an actions column with ellipsis button
   * and one popup Menu shared by every row.
   * Items can have: label, icon, command, class, route (for router-link), separator.
   * Flat only — Menu renders a nested `items` array as a section header with its
   * children inline, and drops anything deeper.
   */
  actionItems: {
    type: Function,
    default: null,
  },
})

const { t } = useI18n()
const selected = ref([])
const expandedRows = ref([])

// -- Actions menu --
const actionsMenuRef = ref()
const currentMenuItems = ref([])

/**
 * Filters out the 'actions' field from the fields array when actionItems is provided,
 * since the component will render its own built-in actions column.
 */
const computedFields = computed(() => {
  if (props.actionItems) {
    return props.fields.filter(f => f.key !== 'actions')
  }
  return props.fields
})

// show() rather than toggle(): one popup serves every row, so clicking a second
// row's button must re-anchor and open there rather than close the first.
function showActionsMenu(event, rowData) {
  currentMenuItems.value = props.actionItems(rowData)
  actionsMenuRef.value.show(event, event.currentTarget)
}

function hideActionsMenu() {
  if (actionsMenuRef.value) {
    actionsMenuRef.value.hide()
  }
}

defineExpose({
  hideActionsMenu,
})

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

/* The actively sorted column reads primary — label and sort arrow both — so
   the sort state is visible at a glance, not just by the small arrow. Done in
   CSS (not a conditional header class) because the Aura theme colours the
   sorted th/icon via its own higher-specificity token rules, which quietly win
   over a single utility class. */
:deep(th.p-datatable-column-sorted),
:deep(th.p-datatable-column-sorted .p-sortable-column-icon) {
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
