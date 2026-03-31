<template>
  <div class="flex h-full min-h-0 flex-col">
    <div v-if="loading" class="flex flex-1 items-center justify-center p-8">
      <ProgressSpinner style="width: 2rem; height: 2rem" />
    </div>

    <div v-else-if="items.length" class="flex flex-col flex-1 min-h-0 overflow-auto">
      <NotificationItem
        v-for="notification in items"
        :key="notification.notificationID"
        :notification="notification"
        @mark-read="$emit('mark-read', notification.notificationID)"
        @mark-unread="$emit('mark-unread', notification.notificationID)"
        @delete="$emit('delete', notification.notificationID)"
      />

      <div v-if="hasMorePages" class="p-4 text-center">
        <Button
          :label="loadingMore ? '' : $t('notifications.loadMore')"
          severity="secondary"
          outlined
          :loading="loadingMore"
          @click="$emit('load-more')"
        />
      </div>
    </div>

    <div
      v-else
      class="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center text-muted-color"
    >
      <i class="pi pi-bell text-4xl" />
      <div>{{ $t('notifications.empty') }}</div>
    </div>
  </div>
</template>

<script setup>
import NotificationItem from './NotificationItem.vue'

defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  loadingMore: {
    type: Boolean,
    default: false,
  },
  hasMorePages: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['mark-read', 'mark-unread', 'delete', 'load-more'])
</script>
