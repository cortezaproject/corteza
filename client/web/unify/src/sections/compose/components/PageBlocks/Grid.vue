<template>
  <!-- Single block in view mode: skip grid, use CSS flex to fill -->
  <div v-if="isSingleBlockView" class="single-block-wrapper p-4">
    <div class="block-content">
      <component
        :is="resolveBlock(visibleBlocks[0]?.kind)"
        v-if="resolveBlock(visibleBlocks[0]?.kind)"
        :block="visibleBlocks[0]"
        :blocks="blocks"
        :namespace="namespace"
        :page="page"
        :record="record"
      />
      <div v-if="loading" class="block-busy">
        <ProgressSpinner style="width: 28px; height: 28px" />
      </div>
    </div>
  </div>

  <!-- Multiple blocks or builder mode: use gridstack.
       Padding lives on the wrapper, not .grid-stack itself — gridstack items are
       absolutely positioned and their containing block is the padding edge, so
       padding on .grid-stack is ignored visually. -->
  <div v-else class="p-2">
    <div ref="gridEl" class="grid-stack">
      <div
        v-for="block in visibleBlocks"
        :key="getBlockId(block)"
        class="grid-stack-item"
        :class="editable ? 'builder-grid-item' : 'view-grid-item'"
        :gs-id="getBlockId(block)"
        :gs-x="block.xywh?.[0] ?? 0"
        :gs-y="block.xywh?.[1] ?? 0"
        :gs-w="block.xywh?.[2] ?? 24"
        :gs-h="block.xywh?.[3] ?? 18"
        gs-min-w="6"
        gs-min-h="3"
      >
        <!-- Slot lives outside .grid-stack-item-content (which is overflow:auto) so the
             overlay can extend over the dashed border without triggering scrollbars. -->
        <slot name="item-overlay" :item="{ i: getBlockId(block) }" :block="block" />
        <div class="grid-stack-item-content">
          <div class="block-content">
            <component
              :is="resolveBlock(block.kind)"
              v-if="resolveBlock(block.kind)"
              :block="block"
              :blocks="blocks"
              :namespace="namespace"
              :page="page"
              :record="record"
            />
            <div v-else class="flex items-center justify-center h-full text-muted-color italic p-2">
              <div class="text-center">
                <i class="pi pi-box text-2xl mb-2" />
                <div>{{ block.kind || $t('block.noConfiguration') }}</div>
              </div>
            </div>
            <div v-if="loading" class="block-busy">
              <ProgressSpinner style="width: 28px; height: 28px" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { GridStack } from 'gridstack'
import 'gridstack/dist/gridstack.min.css'
import ProgressSpinner from 'primevue/progressspinner'
import { resolveBlock } from './registry'

const COLS = 48

const props = defineProps({
  blocks: {
    type: Array,
    default: () => [],
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  page: {
    type: Object,
    default: () => ({}),
  },
  record: {
    type: Object,
    default: undefined,
  },
  editable: {
    type: Boolean,
    default: false,
  },
  // Covers each block where it stands, rather than the page blanking. The
  // blocks keep their geometry and their current render underneath, so what is
  // on screen is only ever a settled state.
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['layout-updated'])

const visibleBlocks = computed(() => props.blocks.filter(block => !block.meta?.hidden))

const isSingleBlockView = computed(() => !props.editable && visibleBlocks.value.length === 1)

function getBlockId(block) {
  const bid = block.blockID
  if (bid && bid !== '0') return String(bid)
  return String(block.meta?.tempID || '')
}

// Gridstack ships CSS only for 12 columns. Inject percentage rules for 48 cols once.
function ensureColumnStyles() {
  const styleId = `gs-cols-${COLS}-style`
  if (document.getElementById(styleId)) return
  const rules = [`.gs-${COLS} > .grid-stack-item { width: ${100 / COLS}%; }`]
  for (let i = 1; i < COLS; i++) {
    rules.push(`.gs-${COLS} > .grid-stack-item[gs-x="${i}"] { left: ${(i * 100) / COLS}%; }`)
  }
  for (let i = 2; i <= COLS; i++) {
    rules.push(`.gs-${COLS} > .grid-stack-item[gs-w="${i}"] { width: ${(i * 100) / COLS}%; }`)
  }
  const style = document.createElement('style')
  style.id = styleId
  style.textContent = rules.join('\n')
  document.head.appendChild(style)
}

const gridEl = ref(null)
let grid = null

function initGrid() {
  if (!gridEl.value) return
  ensureColumnStyles()

  const opts = {
    column: COLS,
    cellHeight: 10,
    // Gridstack applies margin to all 4 sides of each item, so the gap between adjacent
    // items is 2 × margin. Half the old [12, 12] inter-item gap to preserve spacing.
    margin: 6,
    float: true,
    animate: true,
    disableDrag: !props.editable,
    disableResize: !props.editable,
    draggable: { handle: '.block-drag-handle' },
    resizable: { handles: 'all', autoHide: true },
    alwaysShowResizeHandle: false,
  }

  // Mobile reflow only in view mode — builder stays full 48-col so authors design for desktop.
  // columnMax must be set explicitly; otherwise gridstack defaults it to 12 and collapses the
  // grid on any width above the 1-col breakpoint.
  if (!props.editable) {
    opts.columnOpts = {
      breakpointForWindow: true,
      columnMax: COLS,
      breakpoints: [{ w: 768, c: 1, layout: 'list' }],
    }
  }

  grid = GridStack.init(opts, gridEl.value)

  grid.on('change', onGridChange)
}

function onGridChange(_event, items) {
  if (props.editable && Array.isArray(items)) {
    for (const item of items) {
      const block = props.blocks.find(b => getBlockId(b) === item.id)
      if (block) {
        block.xywh = [item.x ?? 0, item.y ?? 0, item.w ?? 1, item.h ?? 1]
      }
    }
  }
  emit('layout-updated', items)
}

function rebuildLayout() {
  if (!grid || !gridEl.value) return
  grid.batchUpdate()
  // Detach all from gridstack tracking; Vue still owns the DOM nodes.
  grid.removeAll(false)
  for (const block of visibleBlocks.value) {
    const id = getBlockId(block)
    if (!id) continue
    const el = gridEl.value.querySelector(`.grid-stack-item[gs-id="${CSS.escape(id)}"]`)
    if (el) grid.makeWidget(el)
  }
  grid.commit()
}

onMounted(() => {
  // Defer init until v-for children are in the DOM so gridstack picks them up.
  nextTick(initGrid)
})

onBeforeUnmount(() => {
  if (grid) {
    grid.off('change', onGridChange)
    grid.destroy(false)
    grid = null
  }
})

// Rebuild only when the set of visible block IDs changes (add/remove/visibility),
// not when positions inside xywh change — gridstack already owns those.
watch(
  () => visibleBlocks.value.map(getBlockId).join('|'),
  () => {
    nextTick(rebuildLayout)
  },
)

defineExpose({ rebuildLayout })
</script>

<style>
/* Drag-target placeholder: subtle dashed outline instead of a solid fill — a filled
 * placeholder flashes hard each time it jumps to a new cell as you drag. */
.grid-stack > .grid-stack-placeholder > .placeholder-content {
  border: 2px dashed var(--p-primary-color);
  border-radius: var(--p-card-border-radius);
  opacity: 0.5;
}

/* Item border (view mode keeps it invisible to match builder layout) */
.grid-stack > .view-grid-item > .grid-stack-item-content {
  border: 2px solid transparent;
  border-radius: var(--p-card-border-radius);
}

.grid-stack > .builder-grid-item > .grid-stack-item-content {
  border: 2px dashed var(--p-content-border-color);
  border-radius: var(--p-card-border-radius);
}

/* View mode: let the Card's box-shadow render past the item bounds so blocks look
 * like proper cards. Builder mode keeps the default overflow — otherwise the Card
 * shadow renders through the dashed border's gaps and looks like a second border. */
.grid-stack > .view-grid-item > .grid-stack-item-content {
  overflow: visible;
}

.grid-stack > .builder-grid-item:hover > .grid-stack-item-content {
  /* Use the full shorthand so border-style stays dashed — without this the cascade
   * flips the style to solid on hover. */
  border: 2px dashed var(--p-primary-color);
}

/* Hide the per-block Card border (PageBlock.vue adds `border border-surface` when
 * block.style.border.enabled is true) while editing — the dashed builder indicator
 * is the relevant visual in this mode; the Card border just doubles up underneath.
 * View mode is untouched so a user-configured border still shows there. */
.grid-stack > .builder-grid-item .p-card {
  border: none;
}

/* The drag handle is a PrimeVue <Button> (a real <button>). Gridstack's dd-draggable
 * skips mousedown when e.target.closest('button') matches but e.target isn't the dragEl
 * itself. Disabling pointer-events on descendants forces e.target === button, so
 * gridstack always starts the drag regardless of which pixel inside the button is hit. */
.block-drag-handle * {
  pointer-events: none;
}

/* Resize handles: invisible by design. Corners are shrunk to 10×10 (vs gridstack's
 * default 20×20) and the chevron icon is stripped so they don't crowd the toolbox.
 * Edge handles keep gridstack's default dimensions but have no background.
 * Visual feedback comes from the resize cursor; z-index keeps them under the
 * toolbox (z-index 5) so the drag handle stays clickable. */
.grid-stack > .grid-stack-item > .ui-resizable-handle {
  z-index: 1;
}

.grid-stack > .grid-stack-item > .ui-resizable-se,
.grid-stack > .grid-stack-item > .ui-resizable-sw,
.grid-stack > .grid-stack-item > .ui-resizable-ne,
.grid-stack > .grid-stack-item > .ui-resizable-nw {
  background-image: none;
  width: 10px;
  height: 10px;
}
</style>

<style scoped>
.single-block-wrapper {
  height: 100%;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

.single-block-wrapper > .block-content {
  height: 100%;
  overflow: hidden;
}

.block-content {
  position: relative;
  height: 100%;
  isolation: isolate;
}

/* Sits over the block's own render, which stays mounted: remounting it would
 * cost every block its state and reload what it had already fetched. */
.block-busy {
  position: absolute;
  inset: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--p-content-background);
}
</style>
