<template>
  <div
    v-if="visible"
    data-testid="desktop-notifications-prompt"
    class="mx-2 mt-2 flex items-start gap-2 rounded-md border border-surface-200 bg-surface-50 p-2 text-sm dark:border-surface-700 dark:bg-surface-800"
  >
    <i class="pi pi-desktop mt-0.5 text-muted-color" />

    <div class="flex min-w-0 flex-1 flex-col items-start gap-2 break-words">
      <template v-if="permission === 'default'">
        <span>{{ $t('notifications.desktopPrompt') }}</span>
        <Button
          :label="$t('notifications.desktopEnable')"
          size="small"
          data-testid="desktop-notifications-enable"
          @click="request()"
        />
      </template>

      <span v-else>{{ $t('notifications.desktopBlocked') }}</span>
    </div>

    <Button
      v-tooltip.bottom="$t('notifications.desktopDismiss')"
      :aria-label="$t('notifications.desktopDismiss')"
      icon="pi pi-times"
      severity="secondary"
      variant="text"
      rounded
      size="small"
      data-testid="desktop-notifications-dismiss"
      @click="dismiss()"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import { useSystemNotificationPermission } from '../../composables/useSystemNotifications'

const DISMISSED_KEY = 'notificationsDesktopPromptDismissed'

const notifications = useNotificationsStore()
const { permission, request } = useSystemNotificationPermission()

function readDismissed() {
  try {
    return localStorage.getItem(DISMISSED_KEY) === 'true'
  } catch {
    return false
  }
}

const dismissed = ref(readDismissed())

const visible = computed(
  () =>
    !notifications.muted &&
    !dismissed.value &&
    (permission.value === 'default' || permission.value === 'denied'),
)

function dismiss() {
  dismissed.value = true
  try {
    localStorage.setItem(DISMISSED_KEY, 'true')
  } catch {}
}
</script>
