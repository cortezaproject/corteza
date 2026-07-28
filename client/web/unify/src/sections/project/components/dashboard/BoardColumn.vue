<template>
  <!-- One status column. Drop target for native HTML5 DnD — reads the
       dragged item's key off dataTransfer (set by BoardCard's dragstart) so
       it doesn't need a shared ref with its sibling columns; BoardPanel.vue
       still keeps a single `draggedKey` for the dimmed-source-card styling,
       threaded down as a prop. -->
  <div
    class="flex flex-col min-h-0 rounded-lg border bg-surface transition-colors"
    :class="over ? 'border-primary ring-1 ring-primary' : 'border-surface'"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <!-- Status accent — the SAME hex the status donut and the EventBadge
         pills use (config/chartColors STATUS_COLORS), so a column, a chart
         slice and a badge for the same status always read as one colour.
         Inline style, not a Tailwind class: these are validated hexes owned
         by chartColors, and there is no class equivalent for them. -->
    <div
      class="shrink-0 h-1 rounded-t-lg"
      :style="{ backgroundColor: statusColor }"
      aria-hidden="true"
    />
    <div class="shrink-0 px-3 py-2 flex items-center justify-between gap-2 border-b border-surface">
      <span class="flex items-center gap-2 min-w-0">
        <span
          class="w-2 h-2 rounded-full shrink-0"
          :style="{ backgroundColor: statusColor }"
          aria-hidden="true"
        />
        <span class="text-sm font-medium truncate">{{ status }}</span>
        <!-- Count sits with the title, not in the right-hand group: it
             describes the column, whereas the right side is for actions.
             bg-emphasis, not bg-surface — the column itself is bg-surface now,
             so a surface-toned pill would be invisible against it. -->
        <span class="text-xs text-muted-color rounded-full bg-emphasis px-1.5 py-0.5 shrink-0">
          {{ items.length }}
        </span>
      </span>
      <span class="flex items-center shrink-0">
        <!-- Quick-add — opens BoardPanel's create dialog pre-filled with this
             column's status: every column gets one so an empty board can
             still be filled in from cold. Unlike BoardCard's click
             (view-only, never gated), creating is a write action, so this
             respects `disabled` same as drag. -->
        <Button
          type="button"
          icon="pi pi-plus"
          severity="secondary"
          text
          rounded
          size="small"
          :disabled="disabled"
          :aria-label="addLabel"
          v-tooltip.top="addLabel"
          @click="$emit('add-item')"
        />
      </span>
    </div>
    <div class="flex-1 min-h-0 overflow-y-auto p-2 flex flex-col gap-2">
      <template v-if="items.length">
        <BoardCard
          v-for="item in items"
          :key="item.key"
          :item="item"
          :disabled="disabled"
          :dragging="item.key === draggedKey"
          :show-revision="showRevision"
          @dragstart="$emit('card-dragstart', $event)"
          @dragend="$emit('card-dragend')"
          @click="$emit('card-click', $event)"
        />
      </template>
      <!-- Per-column empty affordance — an empty status still reads as a
           real (drop-targetable, quick-addable) column rather than a blank
           gap, which is the whole point of always rendering every column
           (see BoardPanel.vue). This is now the ONLY empty-state affordance
           the board carries — the old whole-board note above the columns
           was dropped as redundant noise on top of this. -->
      <div v-else class="flex-1 flex items-center justify-center px-2">
        <p class="text-xs text-muted-color text-center">{{ $t('project.dashboard.list.empty') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import BoardCard from './BoardCard.vue'
import { colorFor } from '@/sections/project/config/chartColors'
import { computed, ref } from 'vue'

const props = defineProps({
  status: { type: String, required: true },
  items: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
  // The key of the card currently being dragged (shared across every
  // column so only the true source card dims, not just the one under it).
  draggedKey: { type: String, default: null },
  // Tooltip/aria-label for the quick-add button — owned by BoardPanel.vue
  // (it knows the default created type), kept out of this presentational
  // component.
  addLabel: { type: String, default: '' },
  // Chain-wide only (BoardPanel's `!isRevisionScoped`) — passed straight
  // through to every card so each shows which revision it belongs to. False
  // revision-scoped: every card already shares the one open revision, so the
  // chip would be redundant on every card (same reasoning CategoryPanel uses
  // to drop its revision column when scoped).
  showRevision: { type: Boolean, default: false },
})
const emit = defineEmits(['drop-item', 'card-dragstart', 'card-dragend', 'add-item', 'card-click'])

// Resolved through colorFor rather than reading STATUS_COLORS directly, so an
// unrecognised status degrades to the shared MUTED grey instead of undefined.
const statusColor = computed(() => colorFor('status', props.status))

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
