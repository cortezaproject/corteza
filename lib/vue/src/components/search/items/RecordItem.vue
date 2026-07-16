<template>
  <div
    class="flex items-center justify-between px-4 py-2 cursor-pointer hover:bg-emphasis group"
    @click="$emit('click')"
  >
    <div class="flex flex-col min-w-0">
      <span class="text-sm font-medium text-color truncate">{{ hit.highlight?.label }}</span>
      <span class="text-xs text-muted-color truncate">{{ hit.highlight?.value }}</span>
    </div>

    <Button
      v-tooltip.left="labels.openInNewTab"
      icon="pi pi-external-link"
      severity="secondary"
      variant="text"
      size="small"
      rounded
      class="ml-2 opacity-0 group-hover:opacity-100 shrink-0"
      @click.stop="$emit('open-new-tab')"
    />
  </div>
</template>

<script setup>
defineProps({
  hit: {
    type: Object,
    required: true,
  },
  labels: {
    type: Object,
    default: () => ({}),
  },
})

defineEmits(['click', 'open-new-tab'])
</script>
