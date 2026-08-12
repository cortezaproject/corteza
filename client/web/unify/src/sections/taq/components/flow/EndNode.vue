<template>
  <div
    class="end-node-wrapper flex justify-center"
    :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
  >
    <Handle type="target" :position="Position.Top" />

    <div
      class="flex items-center justify-center py-2 px-6 bg-highlight text-sm font-semibold text-color rounded-border"
    >
      {{ $t('builder.nodes.end') }}
    </div>
  </div>
</template>

<script setup>
import { NODE_DIMENSIONS } from '@/sections/taq/utils/flow-constants'
import { Handle, Position } from '@vue-flow/core'

defineProps({
  id: { type: String, required: true },
  data: { type: Object, default: () => ({}) },
})
</script>
