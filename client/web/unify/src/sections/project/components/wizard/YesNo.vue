<template>
  <i
    v-if="on"
    class="pi pi-check-circle text-emerald-600 dark:text-emerald-400"
    :aria-label="$t('general.label.yes')"
  />
  <i
    v-else
    class="pi pi-minus text-surface-300 dark:text-surface-600"
    :aria-label="$t('general.label.no')"
  />
</template>

<script setup>
defineProps({
  on: { type: Boolean, default: false },
})
</script>
