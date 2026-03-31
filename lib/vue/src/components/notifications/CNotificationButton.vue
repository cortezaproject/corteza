<template>
  <div class="relative">
    <Button
      v-tooltip.bottom="$t('notifications.title')"
      icon="pi pi-bell"
      severity="secondary"
      variant="text"
      rounded
      @click="notifications.toggleVisibility()"
    />
    <Badge
      v-if="notifications.unreadCount > 0 && !notifications.muted"
      :value="notifications.unreadCount > 9 ? '9+' : notifications.unreadCount"
      class="!absolute !-top-1 !-right-1 !pointer-events-none"
      severity="danger"
      size="small"
    />
  </div>
</template>

<script setup>
import { useNotificationsStore } from '../../stores/useNotificationsStore'

const notifications = useNotificationsStore()
</script>
