<template>
  <div
    class="end-node-wrapper flex justify-center"
    :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
  >
    <Handle type="target" :position="Position.Top" />

    <div
      class="flex items-center justify-center p-2 rounded-full border-2 border-surface bg-highlight shadow-sm"
    >
      <i class="pi pi-refresh text-muted-color-emphasis text-sm font-bold"></i>
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
