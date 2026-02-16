<template>
  <GridLayout
    :layout="layout"
    :col-num="48"
    :row-height="10"
    :margin="[0, 0]"
    :is-draggable="editable"
    :is-resizable="editable"
    :responsive="true"
    :breakpoints="{ lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 }"
    :cols="{ lg: 48, md: 48, sm: 1, xs: 1, xxs: 1 }"
    :vertical-compact="true"
    @layout-updated="onLayoutUpdated"
  >
    <GridItem
      v-for="item in layout"
      :key="item.i"
      :i="item.i"
      :x="item.x"
      :y="item.y"
      :w="item.w"
      :h="item.h"
      :min-w="6"
      :min-h="5"
    >
      <component
        :is="resolveBlock(blockMap.get(item.i)?.kind)"
        v-if="resolveBlock(blockMap.get(item.i)?.kind)"
        :block="blockMap.get(item.i)"
        :namespace="namespace"
        :page="page"
      />
      <div v-else class="p-3 text-muted-color italic">
        {{ $t('block.noConfiguration') }} ({{ blockMap.get(item.i)?.kind }})
      </div>
    </GridItem>
  </GridLayout>
</template>

<script setup>
import { computed } from 'vue'
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
  editable: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:blocks'])

const blockMap = computed(() => {
  const map = new Map()
  for (const block of props.blocks) {
    const id = block.blockID || block.meta?.tempID
    if (id) map.set(String(id), block)
  }
  return map
})

const layout = computed(() => {
  return props.blocks.map(block => {
    const [x, y, w, h] = block.xywh || [0, 0, 48, 15]
    return {
      i: String(block.blockID || block.meta?.tempID),
      x,
      y,
      w,
      h,
    }
  })
})

function onLayoutUpdated(newLayout) {
  if (!props.editable) return

  const updated = props.blocks.map(block => {
    const id = String(block.blockID || block.meta?.tempID)
    const item = newLayout.find(l => l.i === id)
    if (item) {
      return { ...block, xywh: [item.x, item.y, item.w, item.h] }
    }
    return block
  })

  emit('update:blocks', updated)
}
</script>
