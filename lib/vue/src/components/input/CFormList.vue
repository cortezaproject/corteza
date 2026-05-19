<template>
  <div
    v-if="!items.length && emptyMessage"
    class="text-muted-color text-sm p-3 border border-surface rounded-border bg-highlight text-center"
  >
    {{ emptyMessage }}
  </div>

  <div
    v-else-if="items.length"
    class="flex flex-col gap-4"
  >
    <!-- Column headers -->
    <div v-if="hasHeaders" :style="gridStyle" class="grid gap-2 px-3 pt-2">
      <span
        v-for="(col, i) in columns"
        :key="i"
        class="text-xs font-semibold text-muted-color uppercase px-2"
        :class="col.headerClass"
      >
        <span v-if="col.label">{{ col.label }}</span>
        <i
          v-if="col.tooltip"
          v-tooltip.top="col.tooltip"
          class="pi pi-info-circle text-xs cursor-help ml-1"
        />
      </span>
      <span v-if="!hideRemove" class="w-10" />
    </div>

    <!-- Item rows -->
    <div
      v-for="(item, index) in items"
      :key="index"
      class="border border-surface rounded-border p-3 flex flex-col gap-2 shadow-sm hover:bg-emphasis transition-colors"
    >
      <div :style="gridStyle" class="grid gap-2 items-center">
        <slot name="row" :item="item" :index="index" />
        <div v-if="!hideRemove" class="w-10 flex justify-end">
          <Button
            v-if="items.length > minItems"
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
  </div>
</template>

<script setup>
import { computed } from 'vue'

const items = defineModel({ type: Array, required: true })

const props = defineProps({
  columns: { type: Array, default: () => [] },
  emptyMessage: { type: String, default: '' },
  minItems: { type: Number, default: 0 },
  hideRemove: { type: Boolean, default: false },
})

const emit = defineEmits(['change'])

const hasHeaders = computed(() => props.columns.some(c => c.label || c.tooltip))

const gridStyle = computed(() => {
  const cols = props.columns.length
    ? props.columns.map(c => c.width || '1fr')
    : ['1fr']
  return {
    gridTemplateColumns: [
      ...cols,
      ...(props.hideRemove ? [] : ['auto']),
    ].join(' '),
  }
})

function remove(index) {
  items.value.splice(index, 1)
  emit('change')
}
</script>
