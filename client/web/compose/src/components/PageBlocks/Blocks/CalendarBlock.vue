<template>
  <PageBlock :block="block">
    <div class="flex items-center justify-center h-full p-4 text-muted-color">
      <div class="text-center">
        <i class="pi pi-calendar text-3xl mb-2" />
        <div class="text-sm italic">{{ $t('block.calendar.label') }}</div>
        <div class="text-xs mt-1">{{ $t('block.calendar.viewLabel') }}</div>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import PageBlock from './PageBlock.vue'

defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

// Calendar integration requires a full calendar library (FullCalendar)
// which will be implemented as a separate feature
</script>
