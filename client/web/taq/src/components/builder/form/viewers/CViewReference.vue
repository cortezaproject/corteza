<template>
  <div
    class="inline-flex items-center gap-1.5 bg-surface border border-surface text-sm"
    :title="displayLabel"
    :style="{ padding: 'var(--p-form-field-sm-padding-y) var(--p-form-field-sm-padding-x)', borderRadius: 'var(--p-form-field-border-radius)' }"
  >
    <i class="pi pi-link text-primary text-sm" />
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
