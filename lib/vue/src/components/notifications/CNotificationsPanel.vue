<template>
  <div class="flex flex-col min-h-0 flex-1">
    <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
      <i class="pi pi-bell text-primary text-sm" />
      <span class="font-semibold text-base text-color flex-1">{{ $t('notifications.title') }}</span>
      <slot name="actions" />
    </div>

    <div class="flex min-h-0 flex-1 flex-col">
      <div class="min-h-0 flex-1">
        <Notifications />
      </div>
    </div>
  </div>
</template>

<script setup>
import Notifications from './Notifications.vue'
</script>
