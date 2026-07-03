<template>
  <!-- Shared item cell for form lists (the CFormItemList `#default` slot and the
       collapsible card headers). One place owns the name/subtitle layout so every
       wizard step's list rows look identical. Pass `title`/`subtitle` for the
       common case; use the `#subtitle` slot for a richer second line and the
       optional `#icon` slot for a leading badge (rarely — the step header already
       carries the kind icon). -->
  <div class="flex items-center gap-3 min-w-0">
    <slot name="icon" />
    <div class="min-w-0">
      <div class="font-medium truncate">{{ title }}</div>
      <slot v-if="$slots.subtitle" name="subtitle" />
      <div v-else-if="subtitle" class="text-xs text-muted-color truncate">{{ subtitle }}</div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>
