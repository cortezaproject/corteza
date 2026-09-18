<template>
  <div class="relative">
    <Button
      v-tooltip.bottom="$t('notifications.title')"
      icon="pi pi-bell"
      severity="secondary"
      variant="text"
      rounded
      @click="rightSidebarStore.toggle('notifications')"
    />
    <Badge
      v-if="notifications.badgeCount > 0"
      :value="notifications.badgeCount > 9 ? '9+' : notifications.badgeCount"
      class="!absolute !-top-1 !-right-1 !pointer-events-none"
      severity="danger"
      size="small"
    />
  </div>
</template>

<script setup>
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import { useRightSidebarStore } from '../../stores/useRightSidebarStore'

const notifications = useNotificationsStore()
const rightSidebarStore = useRightSidebarStore()
</script>
