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
        <!-- The endpoint's TRUE per-column total (server/system/types/
             project_board.go's ProjectBoardColumn.Total), not just how many
             cards this page has loaded — a column past its first page still
             reads its real count, not "20+". -->
        <span class="text-xs text-muted-color rounded-full bg-emphasis px-1.5 py-0.5 shrink-0">
          {{ total }}
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
    <div class="flex-1 min-h-0 overflow-y-auto p-2 flex flex-col gap-2" @scroll="onScroll">
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

      <!-- This column's next page — server-paged off the board endpoint's
           own cursor (BoardPanel.vue's loadMoreColumn), never a client-side
           slice. Scrolling near the bottom auto-triggers it (onScroll
           below); the button underneath is the same action, for anyone who'd
           rather click than scroll. -->
      <div v-if="hasMore" class="shrink-0 flex items-center justify-center py-1">
        <ProgressSpinner v-if="loadingMore" style="width: 1rem; height: 1rem" stroke-width="6" />
        <button
          v-else
          type="button"
          class="text-xs text-muted-color hover:text-color underline-offset-2 hover:underline"
          @click="$emit('load-more')"
        >
          {{ $t('project.dashboard.board.loadMore') }}
        </button>
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
  // The endpoint's TRUE per-column total (see the header pill above) —
  // defaults to 0 so this column still reads sensibly before its first load
  // resolves.
  total: { type: Number, default: 0 },
  // Whether this column has a further page to fetch (BoardPanel's
  // column.nextPage, threaded through as a plain boolean).
  hasMore: { type: Boolean, default: false },
  // A load-more request for this column is in flight.
  loadingMore: { type: Boolean, default: false },
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
const emit = defineEmits([
  'drop-item',
  'card-dragstart',
  'card-dragend',
  'add-item',
  'card-click',
  'load-more',
])

// Infinite-scroll trigger — fires once per approach to the bottom (loadingMore
// guards re-entrancy; BoardPanel also no-ops a load-more call with no
// nextPage, so a stray extra emit near the boundary is harmless).
function onScroll(e) {
  if (!props.hasMore || props.loadingMore) return
  const el = e.target
  if (el.scrollHeight - el.scrollTop - el.clientHeight < 80) {
    emit('load-more')
  }
}

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
