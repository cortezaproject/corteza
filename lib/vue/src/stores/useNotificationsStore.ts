import { system } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

export const useNotificationsStore = defineStore('notifications', () => {
  const $SystemAPI = inject<any>('$SystemAPI')
  const notifications = ref<Array<system.Notification>>([])
  const pageCursor = ref<string | null>(null)
  const muted = ref(localStorage.getItem('notificationsMuted') === 'true')

  const hasMorePages = computed(() => !!pageCursor.value)
  const hasUnread = computed(() => notifications.value.some(notification => !notification.readAt))
  const hasRead = computed(() => notifications.value.some(notification => !!notification.readAt))
  const unreadCount = computed(() => notifications.value.filter(notification => !notification.readAt).length)

  function setNotifications(set: Array<system.Notification> = []) {
    notifications.value = set.map(notification => new system.Notification(notification))
  }

  function appendNotifications(set: Array<system.Notification> = []) {
    notifications.value = [
      ...notifications.value,
      ...set.map(notification => new system.Notification(notification)),
    ]
  }

  async function fetchNotifications({ unreadOnly = true } = {}) {
    const response = await $SystemAPI.notificationList({
      limit: 25,
      sort: pageCursor.value ? '' : 'createdAt DESC, readAt DESC',
      read: unreadOnly ? 0 : 1,
      pageCursor: pageCursor.value ?? undefined,
    })

    const set = (response.set || []) as Array<system.Notification>
    const filter = response.filter || {}

    if (pageCursor.value) {
      appendNotifications(set)
    } else {
      setNotifications(set)
    }

    pageCursor.value = filter.nextPage || null
    return notifications.value
  }

  async function markAsRead(notificationID: string) {
    await $SystemAPI.notificationMarkAsRead({ notificationID })
    updateReadNotification({ notificationID })
  }

  async function markAsUnread(notificationID: string) {
    await $SystemAPI.notificationMarkAsUnread({ notificationID })
    updateUnreadNotification({ notificationID })
  }

  async function markAllAsRead() {
    if (!notifications.value.length) {
      return
    }

    await $SystemAPI.notificationMarkAllAsRead()
    const now = new Date()
    notifications.value.forEach(notification => {
      if (!notification.readAt) {
        notification.readAt = now
      }
    })
  }

  async function markAllAsUnread() {
    if (!notifications.value.length) {
      return
    }

    await $SystemAPI.notificationMarkAllAsUnread()
    notifications.value.forEach(notification => {
      notification.readAt = undefined
    })
  }

  async function deleteNotification(notificationID: string) {
    await $SystemAPI.notificationDelete({ notificationID })
    removeNotification({ notificationID })
  }

  function setPageCursor(value: string | null) {
    pageCursor.value = value
  }

  function toggleMuted() {
    muted.value = !muted.value
    localStorage.setItem('notificationsMuted', String(muted.value))
  }

  function addNotification(notification: system.Notification) {
    notifications.value.unshift(new system.Notification(notification))
  }

  function updateReadNotification(notification: Partial<system.Notification> & { notificationID: string }) {
    const existing = notifications.value.find(n => String(n.notificationID) === String(notification.notificationID))
    if (existing) {
      existing.readAt = notification.readAt ? new Date(notification.readAt) : new Date()
    }
  }

  function updateUnreadNotification(notification: Partial<system.Notification> & { notificationID: string }) {
    const existing = notifications.value.find(n => String(n.notificationID) === String(notification.notificationID))
    if (existing) {
      existing.readAt = undefined
    }
  }

  function updateAllReadNotifications(set: Array<system.Notification> = []) {
    const ids = new Set(set.map(({ notificationID }) => String(notificationID)))
    const now = new Date()
    notifications.value.forEach(notification => {
      if (ids.has(String(notification.notificationID))) {
        notification.readAt = now
      }
    })
  }

  function updateAllUnreadNotifications(set: Array<system.Notification> = []) {
    const ids = new Set(set.map(({ notificationID }) => String(notificationID)))
    notifications.value.forEach(notification => {
      if (ids.has(String(notification.notificationID))) {
        notification.readAt = undefined
      }
    })
  }

  function removeNotification(notification: { notificationID: string }) {
    notifications.value = notifications.value.filter(n => String(n.notificationID) !== String(notification.notificationID))
  }

  return {
    notifications,
    pageCursor,
    muted,
    hasMorePages,
    hasUnread,
    hasRead,
    unreadCount,
    fetchNotifications,
    markAsRead,
    markAsUnread,
    markAllAsRead,
    markAllAsUnread,
    deleteNotification,
    setPageCursor,
    toggleMuted,
    addNotification,
    updateReadNotification,
    updateUnreadNotification,
    updateAllReadNotifications,
    updateAllUnreadNotifications,
    removeNotification,
  }
})
