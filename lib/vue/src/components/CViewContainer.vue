<template>
  <div :class="classes">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Top-level page wrapper shared by most section views. Two shapes cover
// nearly every call site: a fixed-height container that owns its own
// internal scrolling (the default — e.g. a CResourceList that scrolls
// itself), or a `scroll` container that grows to fill the page and scrolls
// as a whole (used for edit forms stacked in a column). `gap` only applies
// to the `scroll` shape.
const props = defineProps({
  scroll: { type: Boolean, default: false },
  gap: { type: [Number, String], default: 4 },
})

const classes = computed(() => {
  if (!props.scroll) return 'container mx-auto p-4 h-full overflow-hidden min-w-0'
  return String(props.gap) === '5'
    ? 'container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto'
    : 'container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto'
})
</script>
