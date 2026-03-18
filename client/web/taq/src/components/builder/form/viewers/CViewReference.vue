<template>
  <Tag
    :value="displayLabel"
    severity="primary"
    icon="pi pi-link"
    class="text-sm"
    :title="displayLabel"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  scope: { type: String, default: '' },
  source: { type: String, default: '' },
  nodes: { type: Array, default: () => [] },
})

const displayLabel = computed(() => {
  const step = props.nodes.find(n => n.data?.ref === props.scope)
  const stepLabel = step?.data?.label || props.scope
  return `${stepLabel} → ${props.source}`
})
</script>
