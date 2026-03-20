<template>
  <div
    class="inline-flex items-center gap-1 bg-surface border border-surface text-xs px-2 py-1 rounded-border"
    :title="displayLabel"
  >
    <i class="pi pi-link text-primary text-xs" />
    <span class="text-color truncate">{{ displayLabel }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  scope: { type: String, default: '' },
  source: { type: String, default: '' },
  nodes: { type: Array, default: () => [] },
})

const displayLabel = computed(() => {
  return props.source
})
</script>
