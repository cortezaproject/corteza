<template>
  <!-- Vertical drag handle for a resizable side panel. Pairs with
       useRightSidebarResize's `startDrawerResize` (or any compatible
       drag-start handler) via the `mousedown` event. -->
  <div
    class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors shrink-0"
    @mousedown="emit('mousedown', $event)"
  />
</template>

<script setup lang="ts">
defineOptions({ name: 'CResizeHandle' })

const emit = defineEmits<{
  (_e: 'mousedown', _event: MouseEvent): void
}>()
</script>
