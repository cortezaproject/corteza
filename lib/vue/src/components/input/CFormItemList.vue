<template>
  <div v-if="loading" class="flex justify-center p-4">
    <ProgressSpinner />
  </div>

  <div
    v-else-if="!items.length && emptyMessage"
    class="text-muted-color p-4 border rounded-lg bg-emphasis text-center"
  >
    {{ emptyMessage }}
  </div>

  <div v-else-if="items.length" class="flex flex-col gap-4">
    <div
      v-for="(item, index) in items"
      :key="getKey(item, index)"
      :draggable="draggable"
      :class="[
        'group flex items-center gap-2 p-3 border border-surface rounded-border shadow-sm cursor-pointer hover:bg-emphasis transition-colors',
        isSelected(item, index) ? 'bg-highlight' : 'bg-surface',
        draggable && dropTargetIndex === index && draggedIndex !== index
          ? '!border-t-2 !border-t-primary'
          : '',
      ]"
      @click="$emit('select', item, index)"
      @dragstart="draggable && onDragStart(index)"
      @dragover="draggable && onDragOver($event, index)"
      @dragleave="draggable && onDragLeave()"
      @drop.prevent="draggable && onDrop(index)"
      @dragend="draggable && onDragEnd()"
    >
      <i v-if="draggable" class="pi pi-bars text-muted-color cursor-grab shrink-0 mx-2" />
      <div class="flex-1 min-w-0 min-h-10 flex flex-col justify-center">
        <slot :item="item" :index="index" />
      </div>
      <!-- Right-side controls. With `reveal-on-hover` the status (`actions`
           slot) sits at the far right and is swapped for the action buttons
           (`hover-actions` slot + remove) on hover/focus; otherwise both stay
           visible side by side (legacy layout via `display: contents`). -->
      <div :class="revealOnHover ? 'relative flex items-center' : 'contents'">
        <div
          :class="
            revealOnHover
              ? 'transition-opacity group-hover:opacity-0 group-focus-within:opacity-0'
              : 'contents'
          "
        >
          <slot name="actions" :item="item" :index="index" />
        </div>
        <div
          class="flex items-center gap-1"
          :class="
            revealOnHover
              ? 'absolute inset-y-0 right-0 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100'
              : ''
          "
        >
          <slot name="hover-actions" :item="item" :index="index" />
          <Button
            v-if="!hideRemove"
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            :loading="loadingKey !== null && loadingKey === getKey(item, index)"
            :aria-label="removeLabel || undefined"
            :title="removeLabel || undefined"
            @click.stop="$emit('remove', item, index)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  items: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  emptyMessage: { type: String, default: '' },
  hideRemove: { type: Boolean, default: false },
  itemKey: { type: [String, Function], default: 'id' },
  loadingKey: { type: [String, Number, null], default: null },
  removeLabel: { type: String, default: '' },
  draggable: { type: Boolean, default: false },
  // Key of the currently selected item — compared against getKey(item, index).
  // When matched, the row gets bg-highlight styling. Pass null/undefined for no selection.
  selectedKey: { type: [String, Number, null], default: null },
  // When true, the `actions` slot (e.g. a status tag) sits at the far right and
  // is swapped for the action buttons (`hover-actions` slot + remove) when the
  // row is hovered or focused. Default keeps everything visible side by side
  // (legacy behaviour).
  revealOnHover: { type: Boolean, default: false },
})

const emit = defineEmits(['remove', 'reorder', 'select'])

function getKey(item, index) {
  if (typeof props.itemKey === 'function') return props.itemKey(item)
  return item?.[props.itemKey] ?? index
}

function isSelected(item, index) {
  if (props.selectedKey === null || props.selectedKey === undefined) return false
  return getKey(item, index) === props.selectedKey
}

const draggedIndex = ref(null)
const dropTargetIndex = ref(null)

function onDragStart(index) {
  draggedIndex.value = index
}

function onDragOver(e, index) {
  e.preventDefault()
  dropTargetIndex.value = index
}

function onDragLeave() {
  dropTargetIndex.value = null
}

function onDrop(index) {
  const from = draggedIndex.value
  draggedIndex.value = null
  dropTargetIndex.value = null
  if (from === null || from === index) return
  const reordered = [...props.items]
  const [moved] = reordered.splice(from, 1)
  reordered.splice(index, 0, moved)
  emit('reorder', reordered)
}

function onDragEnd() {
  draggedIndex.value = null
  dropTargetIndex.value = null
}
</script>
