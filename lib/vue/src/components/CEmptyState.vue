<template>
  <!-- Small bordered placeholder shown in place of a list/table when there's
       nothing to display. Content-only — callers own the v-if/v-else that
       decides when it's shown. -->
  <div class="text-muted-color text-sm p-4 border rounded-lg bg-emphasis text-center">
    <slot />
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'CEmptyState' })
</script>
