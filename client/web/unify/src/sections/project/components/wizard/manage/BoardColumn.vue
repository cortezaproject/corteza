<template>
  <!-- One status column. Drop target for native HTML5 DnD — reads the
       dragged item's key off dataTransfer (set by BoardCard's dragstart) so
       it doesn't need a shared ref with its sibling columns; ManageBoard.vue
       still keeps a single `draggedKey` for the dimmed-source-card styling,
       threaded down as a prop. -->
  <div
    class="flex flex-col min-h-0 rounded-lg border bg-emphasis/40 transition-colors"
    :class="over ? 'border-primary ring-1 ring-primary' : 'border-surface'"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <div class="shrink-0 px-3 py-2 flex items-center justify-between gap-2 border-b border-surface">
      <span class="text-sm font-medium truncate">{{ status }}</span>
      <span class="text-xs text-muted-color rounded-full bg-surface px-1.5 py-0.5 shrink-0">
        {{ items.length }}
      </span>
    </div>
    <div class="flex-1 min-h-0 overflow-y-auto p-2 flex flex-col gap-2">
      <BoardCard
        v-for="item in items"
        :key="item.key"
        :item="item"
        :disabled="disabled"
        :dragging="item.key === draggedKey"
        @dragstart="$emit('card-dragstart', $event)"
        @dragend="$emit('card-dragend')"
      />
    </div>
  </div>
</template>

<script setup>
import BoardCard from './BoardCard.vue'
import { ref } from 'vue'

const props = defineProps({
  status: { type: String, required: true },
  items: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
  // The key of the card currently being dragged (shared across every
  // column so only the true source card dims, not just the one under it).
  draggedKey: { type: String, default: null },
})
const emit = defineEmits(['drop-item', 'card-dragstart', 'card-dragend'])

const over = ref(false)

function onDragOver(e) {
  if (props.disabled) return
  e.preventDefault()
  e.dataTransfer.dropEffect = 'move'
  over.value = true
}
function onDragLeave() {
  over.value = false
}
function onDrop(e) {
  if (props.disabled) return
  e.preventDefault()
  over.value = false
  const key = e.dataTransfer.getData('text/plain')
  if (key) emit('drop-item', key)
}
</script>
