<template>
  <!-- One status column. Columns share a drag group, so a card crosses between
       them through the shared draggable rather than any state of this
       component's own; the card lands in `items` and BoardPanel.vue is told
       which one moved, so it can set the status and undo the move if the
       server refuses it. -->
  <div class="flex flex-col min-h-0 rounded-lg border border-surface bg-surface transition-colors">
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
    <CDraggableList
      v-model="cards"
      class="flex-1 min-h-0 overflow-y-auto p-2 flex flex-col gap-2"
      :drag-key="status"
      :group="dragGroup"
      :disabled="disabled"
      @scroll="onScroll"
      @add="$emit('card-added', $event.item?.key)"
    >
      <!-- The v-for is never wrapped in a v-if on `items.length`: dragging the
           last card out would tear the whole branch down while the sortable is
           putting the dragged node back, and the node is left behind, untracked
           and visible in a column the data says is empty. The empty state is a
           sibling for the same reason. -->
      <div v-for="item in items" :key="item.key" data-drag-item>
        <BoardCard
          :item="item"
          :disabled="disabled"
          :show-revision="showRevision"
          @click="$emit('card-click', $event)"
        />
      </div>
      <!-- Per-column empty affordance — an empty status still reads as a
           real (drop-targetable, quick-addable) column rather than a blank
           gap, which is the whole point of always rendering every column
           (see BoardPanel.vue). This is now the ONLY empty-state affordance
           the board carries — the old whole-board note above the columns
           was dropped as redundant noise on top of this. -->
      <div v-if="!items.length" class="flex-1 flex items-center justify-center px-2">
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
    </CDraggableList>
  </div>
</template>

<script setup>
import BoardCard from './BoardCard.vue'
import { colorFor } from '@/sections/project/config/chartColors'
import { components } from '@planetcrust/human-vue'
import { computed } from 'vue'

const { CDraggableList } = components

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
const emit = defineEmits(['card-added', 'update:items', 'add-item', 'card-click', 'load-more'])

// Every column of one board shares a group, so a card may be pulled from any of
// them and put into any other.
const dragGroup = { name: 'project-board-items', pull: true, put: true }

// The column's cards are a prop; the sortable writes the new list back up so
// BoardPanel, which owns the board's state, stays the only thing that mutates
// it.
const cards = computed({
  get: () => props.items,
  set: next => emit('update:items', next),
})

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
</script>
