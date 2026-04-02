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
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps({
  notification: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['mark-read', 'mark-unread', 'delete'])

const menu = ref()
const router = useRouter()
const $ComposeAPI = inject('$ComposeAPI', null)
const $toast = inject('$toast', null)
const { t } = useI18n()

function formatDateTime(value) {
  if (!value) {
    return ''
  }

  const date = value instanceof Date ? value : new Date(value)
  return date.toLocaleString()
}

async function openRecordNotification() {
  const { config = {}, kind } = props.notification
  if (kind !== 'record' || !$ComposeAPI) {
    return
  }

  const { namespaceID, moduleID, recordID, openMode, edit } = config

  try {
    const namespace = await $ComposeAPI.namespaceRead({ namespaceID })
    if (!namespace) {
      $toast?.toastDanger?.(t('notifications.namespaceNotFound'))
      return
    }

    const { set: recordPages = [] } = await $ComposeAPI.pageList({ namespaceID, moduleID })
    if (!recordPages.length) {
      $toast?.toastDanger?.(t('notifications.pageNotFound'))
      return
    }

    if (recordID && recordID !== '0') {
      try {
        const record = await $ComposeAPI.recordRead({ namespaceID, moduleID, recordID })
        if (!record) {
          $toast?.toastDanger?.(t('notifications.recordNotFound'))
          return
        }
      } catch {
        $toast?.toastDanger?.(t('notifications.recordNotFound'))
        return
      }
    }

    const slug = namespace.slug || namespace.namespaceID
    const pageID = recordPages[0].pageID
    const recID = !recordID || recordID === '0' ? '0' : recordID

    if (!router.hasRoute('page.record')) {
      const u = new URL(window.location)
      let url = `${u.origin}/compose/namespace/${slug}/pages/${pageID}/records/${recID}`
      if (edit) {
        url += '?edit=1'
      }

      if (openMode === 'newTab') {
        window.open(url, '_blank', 'noopener')
      } else {
        window.location = url
      }
      return
    }

    const routeLocation = {
      name: 'page.record',
      params: {
        slug,
        pageID,
        recordID: recID,
      },
      query: edit ? { edit: '1' } : {},
    }

    if (openMode === 'newTab') {
      const resolved = router.resolve(routeLocation)
      window.open(resolved.href, '_blank', 'noopener')
    } else {
      router.push(routeLocation)
    }
  } catch {
    $toast?.toastDanger?.(t('notifications.recordRedirectError'))
  }
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
          command: () => openRecordNotification(),
        },
      ]
    : []),
])

function handleClick() {
  if (!props.notification.readAt) {
    emit('mark-read')
  }

  if (props.notification.kind === 'record') {
    openRecordNotification()
  }
}

function toggleMenu(event) {
  menu.value?.toggle(event)
}
</script>
