<template>
  <div
    class="p-3 border-b surface-border hover:bg-emphasis cursor-pointer"
    :class="{ 'opacity-60': notification.readAt }"
  >
    <div
      class="group relative border surface-border bg-surface rounded-border p-3 pr-2 transition-colors"
      @click="handleClick"
    >
      <div class="flex items-start gap-2">
        <div class="font-semibold text-color min-w-0 flex-1">
          {{ notification.config?.title || '' }}
        </div>

        <Menu ref="menu" :model="menuItems" :popup="true" />
        <Button
          icon="pi pi-ellipsis-v"
          severity="secondary"
          variant="text"
          size="small"
          class="absolute top-0 right-0 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
          @click.stop="toggleMenu"
        />
      </div>

      <div
        v-if="notification.config?.description"
        class="mt-1 text-sm text-muted-color"
        :class="{ 'opacity-60': notification.readAt }"
      >
        {{ notification.config.description }}
      </div>
    </div>
    <div class="mt-3 text-xs text-muted-color text-right">
      {{ formatDateTime(notification.createdAt) }}
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useOpenNotification } from '../../composables/useOpenNotification'

const props = defineProps({
  notification: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['mark-read', 'mark-unread', 'delete'])

const menu = ref()
const openNotification = useOpenNotification()
const { t } = useI18n()

function formatDateTime(value) {
  if (!value) {
    return ''
  }

  const date = value instanceof Date ? value : new Date(value)
  return date.toLocaleString()
}

const menuItems = computed(() => [
  {
    label: props.notification.readAt
      ? t('notifications.markAsUnread')
      : t('notifications.markAsRead'),
    icon: props.notification.readAt ? 'pi pi-envelope' : 'pi pi-check-circle',
    command: () => {
      emit(props.notification.readAt ? 'mark-unread' : 'mark-read')
    },
  },
  {
    label: t('notifications.delete'),
    icon: 'pi pi-trash',
    command: () => emit('delete'),
  },
  ...(props.notification.kind === 'record'
    ? [
        {
          label: t('notifications.open'),
          icon: 'pi pi-external-link',
          command: () => openNotification(props.notification),
        },
      ]
    : []),
])

function handleClick() {
  if (!props.notification.readAt) {
    emit('mark-read')
  }

  if (props.notification.kind === 'record') {
    openNotification(props.notification)
  }
}

function toggleMenu(event) {
  menu.value?.toggle(event)
}
</script>
