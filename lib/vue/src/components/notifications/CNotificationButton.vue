<template>
  <div class="relative">
    <Button
      v-tooltip.bottom="
        $t(notifications.muted ? 'notifications.titleMuted' : 'notifications.title')
      "
      :icon="notifications.muted ? 'pi pi-bell-slash' : 'pi pi-bell'"
      :severity="notifications.badgeCount > 0 ? undefined : 'secondary'"
      :class="{ 'notification-bell-ring': ringing }"
      variant="text"
      rounded
      data-testid="notification-bell"
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
import { nextTick, onBeforeUnmount, ref } from 'vue'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import { useRightSidebarStore } from '../../stores/useRightSidebarStore'

const RING_MS = 900

const notifications = useNotificationsStore()
const rightSidebarStore = useRightSidebarStore()

const ringing = ref(false)
let ringTimer

// A short wiggle each time a notification arrives, unless muted. Arrivals come
// through handleRealtime, whose own addNotification call $onAction never sees.
const stopRinging = notifications.$onAction(({ name, args, after }) => {
  const arrived =
    name === 'addNotification' ||
    (name === 'handleRealtime' && args[0]?.['@type'] === 'notification')
  if (!arrived) return

  after(async () => {
    if (notifications.muted) return

    clearTimeout(ringTimer)
    ringing.value = false
    await nextTick()
    ringing.value = true
    ringTimer = setTimeout(() => (ringing.value = false), RING_MS)
  })
})

onBeforeUnmount(() => {
  stopRinging()
  clearTimeout(ringTimer)
})
</script>

<style scoped>
.notification-bell-ring :deep(.p-button-icon) {
  animation: notification-bell-ring 0.8s ease-in-out;
  transform-origin: 50% 0;
}

@keyframes notification-bell-ring {
  0%,
  100% {
    transform: rotate(0);
  }
  15% {
    transform: rotate(14deg);
  }
  30% {
    transform: rotate(-12deg);
  }
  45% {
    transform: rotate(9deg);
  }
  60% {
    transform: rotate(-6deg);
  }
  75% {
    transform: rotate(3deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .notification-bell-ring :deep(.p-button-icon) {
    animation: none;
  }
}
</style>
