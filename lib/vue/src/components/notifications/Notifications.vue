<template>
  <div class="flex h-full min-h-0 flex-col">
    <Tabs v-model:value="activeTab" class="flex h-full min-h-0 flex-col">
      <div>
        <div class="flex items-center gap-2">
          <TabList>
            <Tab value="unread">{{ $t('notifications.unread') }}</Tab>
            <Tab value="all">{{ $t('notifications.all') }}</Tab>
          </TabList>

          <div class="ml-auto flex items-center gap-1 mr-1">
            <Button
              v-if="notifications.hasUnread"
              v-tooltip.bottom="$t('notifications.markAllAsRead')"
              icon="pi pi-check-circle"
              severity="secondary"
              variant="text"
              rounded
              size="small"
              @click="handleMarkAllAsRead"
            />

            <Button
              v-tooltip.bottom="
                $t(notifications.muted ? 'notifications.unmute' : 'notifications.mute')
              "
              :icon="notifications.muted ? 'pi pi-bell-slash' : 'pi pi-bell'"
              severity="secondary"
              variant="text"
              rounded
              size="small"
              @click="notifications.toggleMuted()"
            />
          </div>
        </div>
      </div>

      <TabPanels class="flex-1 min-h-0 p-0" :pt="{ root: { class: 'h-full' } }">
        <TabPanel value="unread" class="h-full !p-0" :pt="{ root: { class: 'h-full' } }">
          <NotificationList
            :items="notifications.notifications"
            :loading="loading"
            :loading-more="loadingMore"
            :has-more-pages="notifications.hasMorePages"
            @mark-read="handleMarkAsRead"
            @mark-unread="handleMarkAsUnread"
            @delete="handleDelete"
            @load-more="loadMore"
          />
        </TabPanel>

        <TabPanel value="all" class="h-full !p-0" :pt="{ root: { class: 'h-full' } }">
          <NotificationList
            :items="notifications.notifications"
            :loading="loading"
            :loading-more="loadingMore"
            :has-more-pages="notifications.hasMorePages"
            @mark-read="handleMarkAsRead"
            @mark-unread="handleMarkAsUnread"
            @delete="handleDelete"
            @load-more="loadMore"
          />
        </TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import NotificationList from './NotificationList.vue'

const notifications = useNotificationsStore()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const activeTab = ref('unread')
const loading = ref(false)
const loadingMore = ref(false)

async function loadNotifications() {
  if (!$SystemAPI) {
    return
  }

  await notifications.fetchNotifications($SystemAPI, {
    unreadOnly: activeTab.value === 'unread',
  })
}

async function loadMore() {
  if (!$SystemAPI) {
    return
  }

  loadingMore.value = true
  try {
    await loadNotifications()
  } finally {
    loadingMore.value = false
  }
}

async function handleMarkAsRead(notificationID) {
  if (!$SystemAPI) {
    return
  }

  await notifications.markAsRead($SystemAPI, String(notificationID))
}

async function handleMarkAsUnread(notificationID) {
  if (!$SystemAPI) {
    return
  }

  await notifications.markAsUnread($SystemAPI, String(notificationID))
}

async function handleDelete(notificationID) {
  if (!$SystemAPI) {
    return
  }

  try {
    await notifications.deleteNotification($SystemAPI, String(notificationID))
    $toast?.toastSuccess?.(t('notifications.notificationDeleted'))
  } catch (error) {
    $toast?.toastDanger?.(t('notifications.notificationDeletedError'))
  }
}

async function handleMarkAllAsRead() {
  if (!$SystemAPI) {
    return
  }

  try {
    await notifications.markAllAsRead($SystemAPI)
    $toast?.toastSuccess?.(t('notifications.allMarkedAsRead'))
  } catch (error) {
    $toast?.toastDanger?.(t('notifications.markAllAsReadError'))
  }
}

watch(
  activeTab,
  async () => {
    notifications.setPageCursor(null)
    loading.value = true
    try {
      await loadNotifications()
    } finally {
      loading.value = false
    }
  },
  { immediate: true },
)
</script>
