<template>
  <!-- Single block in view mode: skip grid, use CSS flex to fill -->
  <div v-if="isSingleBlockView" class="single-block-wrapper p-4">
    <div class="block-content">
      <component
        :is="resolveBlock(blocks[0]?.kind)"
        v-if="resolveBlock(blocks[0]?.kind)"
        :block="blocks[0]"
        :blocks="blocks"
        :namespace="namespace"
        :page="page"
        :record="record"
      />
    </div>
  </div>

  <!-- Multiple blocks or builder mode: use grid layout -->
  <GridLayout
    v-else
    v-model:layout="layoutModel"
    :col-num="48"
    :row-height="10"
    :margin="[12, 12]"
    :is-draggable="editable"
    :is-resizable="editable"
    :responsive="true"
    :breakpoints="{ lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 }"
    :cols="{ lg: 48, md: 48, sm: 1, xs: 1, xxs: 1 }"
    :vertical-compact="true"
    :use-css-transforms="true"
    @layout-updated="onLayoutUpdated"
  >
    <GridItem
      v-for="item in layoutModel"
      :key="item.i"
      :i="item.i"
      :x="item.x"
      :y="item.y"
      :w="item.w"
      :h="item.h"
      :min-w="6"
      :min-h="5"
      :class="editable ? 'builder-grid-item' : 'view-grid-item'"
    >
      <!-- Scoped slot for custom per-item overlay (e.g. builder toolbox) -->
      <slot name="item-overlay" :item="item" :block="blockMap.get(item.i)" />

      <!-- Block content — always fills the grid item -->
      <div class="block-content" :class="{ 'pointer-events-none': editable }">
        <component
          :is="resolveBlock(blockMap.get(item.i)?.kind)"
          v-if="resolveBlock(blockMap.get(item.i)?.kind)"
          :block="blockMap.get(item.i)"
          :blocks="blocks"
          :namespace="namespace"
          :page="page"
          :record="record"
        />
        <div v-else class="flex items-center justify-center h-full text-muted-color italic p-2">
          <div class="text-center">
            <i class="pi pi-box text-2xl mb-2" />
            <div>{{ blockMap.get(item.i)?.kind || $t('block.noConfiguration') }}</div>
          </div>
        </div>
      </div>
    </GridItem>
  </GridLayout>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { GridLayout, GridItem } from 'grid-layout-plus'
import { resolveBlock } from './registry'

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
})

const emit = defineEmits(['update:blocks', 'layout-updated'])

// Single block in view mode — bypass grid entirely, let CSS flex handle sizing
const isSingleBlockView = computed(() => !props.editable && props.blocks.length === 1)

// Unique ID for each block — blockID '0' is NoID (unsaved), so fall back to tempID
function getBlockId(block) {
  const bid = block.blockID
  if (bid && bid !== '0') return bid
  return block.meta?.tempID || ''
}

// Block map for grid lookups
const blockMap = computed(() => {
  const map = new Map()
  for (const block of props.blocks) {
    const id = getBlockId(block)
    if (id) map.set(String(id), block)
  }
  return map
})

// Mutable layout ref so grid-layout-plus can update it directly during drag/resize
const layoutModel = ref([])

// Build layout from blocks
function rebuildLayout() {
  layoutModel.value = props.blocks.map(block => {
    const [x, y, w, h] = block.xywh || [0, 0, 24, 18]
    return {
      i: String(getBlockId(block)),
      x,
      y,
      w,
      h,
    }
  })
}

// Rebuild when blocks array reference changes
watch(() => props.blocks, rebuildLayout, { immediate: true })

// Sync grid positions back to blocks
function onLayoutUpdated(newLayout) {
  if (props.editable) {
    for (const item of newLayout) {
      const block = props.blocks.find(b => String(getBlockId(b)) === item.i)
      if (block) {
        block.xywh = [item.x, item.y, item.w, item.h]
      }
    }
  }
  emit('layout-updated', newLayout)
}

// Expose rebuildLayout so parent can call it after add/clone/delete
defineExpose({ rebuildLayout })
</script>

<style>
.vgl-layout {
  .vgl-item--placeholder {
    background-color: var(--p-highlight-focus-background);
    border-radius: var(--p-card-border-radius);
  }
}

/* Disable grid-layout-plus slide animation globally */
.vue-grid-item {
  transition: none !important;
}

.vgl-item__resizer {
  right: 0.25rem;
  bottom: 0.25rem;

  &::before {
    border: 0 solid var(--p-primary-color);
    border-right-width: var(--vgl-resizer-border-width);
    border-bottom-width: var(--vgl-resizer-border-width);
  }
}

.vgl-item--transform {
  right: auto !important;
  left: 0 !important;
  transition-property: none !important;
}
</style>

<style scoped>
/* Single block: fill all available parent space */
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

/* View mode — invisible border to match builder sizing */
.view-grid-item {
  border: 2px solid transparent;
  border-radius: var(--p-card-border-radius);
}

/* Builder mode — dashed border, same 2px as view for layout parity */
.builder-grid-item {
  position: relative;
  border: 2px dashed var(--p-content-border-color);
  border-radius: var(--p-card-border-radius);
}

.builder-grid-item:hover {
  border-color: var(--p-primary-color);
}

.builder-grid-item:hover :deep(.block-toolbox) {
  opacity: 1;
}

.block-content {
  height: 100%;
}
</style>
