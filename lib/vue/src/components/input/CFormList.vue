<template>
  <div class="overflow-x-auto">
    <div
      v-if="!items.length && !$slots.footer && emptyMessage"
      class="text-muted-color text-sm p-3 border border-surface rounded-border bg-emphasis text-center"
    >
      {{ emptyMessage }}
    </div>

    <!-- Bounded table: header bar + divided rows inside a single border. The
         header and footer live inside the sortable too, and stay put: only
         [data-drag-item] rows are addressed. -->
    <CDraggableList
      v-else
      v-model="items"
      class="flex flex-col min-w-max rounded-border border border-surface bg-surface overflow-hidden"
      handle=".c-drag-handle"
      :disabled="!draggable"
      @update="() => emit('reorder', items)"
    >
      <!-- Column headers (the `cform-list-header` class is a stable hook for
           consumers that restyle the list). -->
      <div
        v-if="hasHeaders && (items.length || $slots.footer)"
        :style="gridStyle"
        class="cform-list-header grid gap-2 py-2 px-3 bg-emphasis border-b border-surface"
      >
        <span v-if="draggable" class="w-10" />
        <span
          v-for="(col, i) in columns"
          :key="i"
          class="text-sm font-semibold text-muted-color uppercase"
          :class="col.headerClass"
        >
          <!-- Same required marker as CFormGroup labels -->
          <span v-if="col.label">
            {{ col.label }}
            <span v-if="col.required" class="text-red-500">*</span>
          </span>
          <i
            v-if="col.tooltip"
            v-tooltip.top="col.tooltip"
            class="pi pi-info-circle text-sm cursor-help ml-1"
          />
        </span>
        <span v-if="!hideRemove" class="w-10" />
      </div>

      <!-- Item rows -->
      <div
        v-for="(item, index) in items"
        :key="index"
        data-drag-item
        class="border-t border-surface first:border-t-0 p-3 flex flex-col gap-2 hover:bg-emphasis transition-colors"
      >
        <div :style="gridStyle" class="grid gap-2 items-center">
          <i
            v-if="draggable"
            class="c-drag-handle pi pi-bars text-muted-color cursor-grab text-center"
          />
          <slot name="row" :item="item" :index="index" />
          <div v-if="!hideRemove" class="w-10 flex justify-end">
            <CInputDelete
              v-if="items.length > minItems && confirmRemove"
              icon="pi pi-trash"
              text
              size="small"
              :message="confirmRemove"
              @confirm="remove(index)"
            />
            <Button
              v-else-if="items.length > minItems"
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="remove(index)"
            />
          </div>
        </div>
        <slot name="extra" :item="item" :index="index" />
      </div>

      <slot
        name="footer"
        :grid-style="gridStyle"
        :draggable="draggable"
        :hide-remove="hideRemove"
      />
    </CDraggableList>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import CInputDelete from './CInputDelete.vue'
import CDraggableList from '../drag/CDraggableList.vue'

const items = defineModel({ type: Array, required: true })

const props = defineProps({
  columns: { type: Array, default: () => [] },
  emptyMessage: { type: String, default: '' },
  minItems: { type: Number, default: 0 },
  hideRemove: { type: Boolean, default: false },
  draggable: { type: Boolean, default: false },
  // When set, removing a row asks for confirmation with this message
  confirmRemove: { type: String, default: '' },
})

const emit = defineEmits(['change', 'reorder'])

const hasHeaders = computed(() => props.columns.some(c => c.label || c.tooltip))

const gridStyle = computed(() => {
  const cols = props.columns.length ? props.columns.map(c => c.width || '1fr') : ['1fr']
  return {
    gridTemplateColumns: [
      ...(props.draggable ? ['2.5rem'] : []),
      ...cols,
      ...(props.hideRemove ? [] : ['2.5rem']),
    ].join(' '),
  }
})

function remove(index) {
  items.value.splice(index, 1)
  emit('change')
}
</script>
