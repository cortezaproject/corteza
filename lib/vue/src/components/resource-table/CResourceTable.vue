<template>
  <div
    class="c-resource-table flex flex-col"
    :class="{ 'h-full': fillHeight, 'is-reorderable': reorderableRows }"
  >
    <!-- Header slot -->
    <div v-if="$slots.header" class="flex items-center justify-between gap-3 mb-3 shrink-0">
      <slot name="header" />
    </div>

    <!-- DataTable -->
    <DataTable
      ref="dataTableRef"
      :value="items"
      :data-key="primaryKey"
      :loading="loading"
      striped-rows
      scrollable
      :scroll-height="computedScrollHeight"
      row-hover
      :resizable-columns="resizable"
      column-resize-mode="expand"
      :row-class="rowClass"
      class="border border-b-0 border-surface rounded-border"
      :pt="{
        headerCell: { class: 'bg-highlight-emphasis' },
        ...pt,
      }"
      v-bind="$attrs"
      @sort="$emit('sort', $event)"
      @row-click="$emit('row-click', $event)"
      @row-reorder="$emit('row-reorder', $event)"
    >
      <template #empty>
        <div class="flex items-center justify-center p-4 text-muted-color">
          <slot name="empty">
            {{ emptyMessage || '—' }}
          </slot>
        </div>
      </template>

      <!-- Row reorder drag handle column -->
      <Column
        v-if="reorderableRows"
        rowReorder
        headerStyle="width: 2rem"
        :pt="{
          headerCell: { class: 'border-r-0' },
          bodyCell: { class: 'cursor-move border-r-0' },
        }"
      />

      <!-- Dynamic columns from fields prop -->
      <Column
        v-for="field in computedFields"
        :key="field.key"
        :field="field.key"
        :header="field.header"
        :sortable="field.sortable || false"
        :style="field.style"
        :header-style="field.headerStyle"
        :header-class="field.headerClass"
        :body-class="field.bodyClass"
        :frozen="field.frozen"
        :align-frozen="field.alignFrozen"
        :pt="field.pt"
      >
        <template v-if="field.hint" #header>
          <span class="flex items-center gap-1">
            {{ field.header }}
            <i v-tooltip.top="field.hint" class="pi pi-info-circle text-xs cursor-help" />
          </span>
        </template>
        <template #body="slotProps">
          <slot :name="`body-${field.key}`" v-bind="slotProps">
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
        :header-style="actionsHeaderStyle"
        :pt="{
          headerCell: { class: 'border-l-0' },
          bodyCell: { class: 'px-2 py-1 border-l-0' },
        }"
      >
        <template #body="slotProps">
          <Button
            v-if="actionItems(slotProps.data, slotProps.index).length"
            icon="pi pi-ellipsis-v"
            text
            severity="secondary"
            size="small"
            class="row-action-btn w-full"
            @click.stop="showActionsMenu($event, slotProps.data, slotProps.index)"
          />
        </template>
      </Column>

      <!-- Pass through any extra Column components from the parent -->
      <slot />
    </DataTable>

    <!-- One popup shared by every row's actions button -->
    <Menu v-if="actionItems" ref="actionsMenuRef" :model="currentMenuItems" popup append-to="body">
      <template #item="{ item, props }">
        <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
          <a v-ripple :href="href" v-bind="props.action" @click="navigate">
            <span :class="item.icon" />
            <span class="ml-2">{{ item.label }}</span>
          </a>
        </router-link>
        <a v-else v-ripple v-bind="props.action" :class="item.class">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </template>
    </Menu>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  /** Array of items to display */
  items: {
    type: Array,
    default: () => [],
  },

  /** Primary key field for DataTable data-key */
  primaryKey: {
    type: String,
    default: '_dataKey',
  },

  /** Column definitions: { key, header, sortable, style, headerStyle, headerClass, bodyClass, frozen, alignFrozen, pt } */
  fields: {
    type: Array,
    default: () => [],
  },

  /** Loading state */
  loading: {
    type: Boolean,
    default: false,
  },

  /** Empty state message */
  emptyMessage: {
    type: String,
    default: '',
  },

  /**
   * Function that returns an array of menu items for a row.
   * Signature: (rowData, rowIndex) => MenuItem[]
   * When provided, an actions column with an ellipsis button is automatically appended.
   * Flat only — Menu renders a nested `items` array as a section header with its
   * children inline, and drops anything deeper.
   */
  actionItems: {
    type: Function,
    default: null,
  },

  /** Whether columns are resizable */
  resizable: {
    type: Boolean,
    default: false,
  },

  /** Whether the table should fill available height */
  fillHeight: {
    type: Boolean,
    default: false,
  },

  /** Scroll height for DataTable. Defaults to 'flex' when fillHeight, otherwise undefined (auto). */
  scrollHeight: {
    type: String,
    default: '',
  },

  /** Custom row class function / string / object */
  rowClass: {
    type: [Function, String, Object],
    default: undefined,
  },

  /** PrimeVue passthrough for DataTable */
  pt: {
    type: Object,
    default: () => ({}),
  },

  /** Whether rows can be reordered via drag-and-drop */
  reorderableRows: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['sort', 'row-click', 'row-reorder'])

const dataTableRef = ref()

// -- Computed fields: filter out legacy 'actions' key if actionItems is provided --
const computedFields = computed(() => {
  if (props.actionItems) {
    return props.fields.filter(f => f.key !== 'actions')
  }
  return props.fields
})

const computedScrollHeight = computed(() => {
  if (props.scrollHeight) return props.scrollHeight
  return props.fillHeight ? 'flex' : undefined
})

const actionsHeaderStyle = computed(() => 'width: 3rem')

// -- Actions menu --
const actionsMenuRef = ref()
const currentMenuItems = ref([])

// show() rather than toggle(): one popup serves every row, so clicking a second
// row's button must re-anchor and open there rather than close the first.
function showActionsMenu(event, rowData, rowIndex) {
  currentMenuItems.value = props.actionItems(rowData, rowIndex)
  actionsMenuRef.value.show(event, event.currentTarget)
}

function hideActionsMenu() {
  if (actionsMenuRef.value) {
    actionsMenuRef.value.hide()
  }
}

defineExpose({
  hideActionsMenu,
  /** Expose underlying DataTable ref for advanced use cases */
  dataTableRef,
})
</script>

<style scoped>
.c-resource-table :deep(.row-action-btn) {
  opacity: 0;
  transition: opacity 0.15s ease;
}

.c-resource-table :deep(.p-datatable-tbody > tr:hover .row-action-btn) {
  opacity: 1;
}

.c-resource-table.is-reorderable :deep(.p-datatable-tbody) {
  user-select: none;
}

.c-resource-table.is-reorderable :deep(input),
.c-resource-table.is-reorderable :deep(textarea) {
  user-select: auto;
}
</style>
