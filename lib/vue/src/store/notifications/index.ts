import { apiClients, system } from '@cortezaproject/corteza-js'
import { StoreOptions } from 'vuex'

const types = {
  fetchNotifications: 'fetchNotifications',
  setNotifications: 'setNotifications',
  appendNotifications: 'appendNotifications',
  setVisible: 'setVisible',
  markAsRead: 'markAsRead',
  markAsUnread: 'markAsUnread',
  markAllAsRead: 'markAllAsRead',
  markAllAsUnread: 'markAllAsUnread',
  addNotification: 'addNotification',
  deleteNotification: 'deleteNotification',
  setPageCursor: 'setPageCursor',
  setMuted: 'setMuted',
  setTotalUnreadCount: 'setTotalUnreadCount',
  updateReadNotification: 'updateReadNotification',
  updateUnreadNotification: 'updateUnreadNotification',
  updateAllReadNotifications: 'updateAllReadNotifications',
  updateAllUnreadNotifications: 'updateAllUnreadNotifications',
}

interface Options {
  api: apiClients.System;
  ws: WebSocket;
  watchInterval: number;
  webapp: string;
}

interface State {
  notifications: Array<system.Notification>;
  visible: boolean;
  pageCursor: string | null;
  muted: boolean;
  totalUnreadCount: number;
}

interface KV {
  [key: string]: unknown;
}

export default function ({ api }: Options): StoreOptions<State> {
  return {
    strict: true,

    state: {
      notifications: [],
      visible: false,
      pageCursor: null,
      muted: localStorage.getItem('notificationsMuted') === 'true' || false,
      totalUnreadCount: 0,
    },

    getters: {
      notifications: (state) => state.notifications,
      visible: (state) => state.visible,
      pageCursor: (state) => state.pageCursor,
      hasMorePages: (state) => !!state.pageCursor,
      muted: (state) => state.muted,

      hasUnread: (state) => state.totalUnreadCount > 0 || state.notifications.some(notification => !notification.readAt),
      hasRead: (state) => state.notifications.some(notification => !!notification.readAt),
      unreadCount: (state) => state.totalUnreadCount,
    },

    actions: {
      toggleVisibility ({ commit, state }) {
        commit(types.setVisible, !state.visible)
      },

      fetchNotifications ({ commit, state, dispatch }, { unreadOnly = true } = {}) {
        const isInitialPage = !state.pageCursor

        return api.notificationList({
          limit: 25,
          sort: state.pageCursor ? '' : 'createdAt DESC, readAt DESC',
          read: unreadOnly ? 0 : 1,
          pageCursor: state.pageCursor,
          // total count can only be requested on the first page
          incTotal: isInitialPage && unreadOnly,
        })
          .then((response: KV) => {
            const set = (response.set || []) as Array<system.Notification>
            const filter = (response.filter || {}) as {
              limit: number;
              sort: string;
              read: number;
              nextPage?: string;
              total?: number;
            }

            if (state.pageCursor) {
              commit(types.appendNotifications, set)
            } else {
              commit(types.setNotifications, set)
            }

            commit(types.setPageCursor, filter.nextPage)

            if (isInitialPage && unreadOnly && typeof filter.total === 'number') {
              commit(types.setTotalUnreadCount, filter.total)
            } else if (isInitialPage && !unreadOnly) {
              // Refresh the total when loading "all" tab so the badge stays accurate
              dispatch('refreshUnreadCount')
            }
          })
      },

      refreshUnreadCount ({ commit }) {
        return api.notificationList({ limit: 1, read: 0, incTotal: true })
          .then((response: KV) => {
            const filter = (response.filter || {}) as { total?: number }
            if (typeof filter.total === 'number') {
              commit(types.setTotalUnreadCount, filter.total)
            }
          })
      },

      markAsRead ({ commit }, notificationID) {
        return api.notificationMarkAsRead({ notificationID })
          .then(() => {
            commit(types.markAsRead, notificationID)
          })
      },

      markAsUnread ({ commit }, notificationID) {
        return (api as any).notificationMarkAsUnread({ notificationID })
          .then(() => {
            commit(types.markAsUnread, notificationID)
          })
      },

      markAllAsRead ({ commit, state }) {
        if (!state.notifications.length && state.totalUnreadCount === 0) {
          return Promise.resolve()
        }

        return api.notificationMarkAllAsRead()
          .then(() => {
            commit(types.markAllAsRead)
          })
      },

      markAllAsUnread ({ commit, state }) {
        if (!state.notifications.length) {
          return Promise.resolve()
        }

        return (api as any).notificationMarkAllAsUnread()
          .then(() => {
            commit(types.markAllAsUnread)
          })
      },

      addNotification ({ commit }, notification) {
        commit(types.addNotification, notification)
      },

      deleteNotification ({ commit }, notificationID) {
        return api.notificationDelete({ notificationID })
          .then(() => {
            commit(types.deleteNotification, notificationID)
          })
      },

      setPageCursor ({ commit }, pageCursor) {
        commit(types.setPageCursor, pageCursor)
      },

      toggleMuted ({ commit, state }) {
        commit(types.setMuted, !state.muted)
      },

      updateReadNotification ({ commit }, notification) {
        commit(types.updateReadNotification, notification)
      },

      updateUnreadNotification ({ commit }, notification) {
        commit(types.updateReadNotification, { ...notification, readAt: undefined })
      },

      updateAllReadNotifications ({ commit }, notifications) {
        commit(types.updateAllReadNotifications, notifications)
      },

      updateAllUnreadNotifications ({ commit }, notifications) {
        commit(types.updateAllUnreadNotifications, notifications)
      },

      removeNotification ({ commit }, notification) {
        commit(types.deleteNotification, notification.notificationID)
      },
    },

    mutations: {
      [types.setNotifications] (state, notifications = []) {
        state.notifications = notifications.map((n: system.Notification) => new system.Notification(n))
      },

      [types.appendNotifications] (state, notifications = []) {
        const newNotifications = notifications.map((n: system.Notification) => new system.Notification(n))
        state.notifications = [...state.notifications, ...newNotifications]
      },

      [types.setVisible] (state, visible) {
        state.visible = visible
      },

      [types.setPageCursor] (state, pageCursor) {
        state.pageCursor = pageCursor
      },

      [types.markAsRead] (state, notificationID) {
        const notification = state.notifications.find(n =>
          String(n.notificationID) === String(notificationID),
        )

        if (notification && !notification.readAt) {
          notification.readAt = new Date()
          state.totalUnreadCount = Math.max(0, state.totalUnreadCount - 1)
        }
      },

      [types.markAllAsRead] (state) {
        const now = new Date()
        state.notifications.forEach(notification => {
          if (!notification.readAt) {
            notification.readAt = now
          }
        })
        state.totalUnreadCount = 0
      },

      [types.markAllAsUnread] (state) {
        state.notifications.forEach(notification => {
          if (notification.readAt) {
            notification.readAt = undefined
          }
        })
      },

      [types.addNotification] (state, notification) {
        // Add to the beginning of the array
        state.notifications.unshift(new system.Notification(notification))
        if (!notification.readAt) {
          state.totalUnreadCount += 1
        }
      },

      [types.markAsUnread] (state, notificationID) {
        const notification = state.notifications.find(n =>
          String(n.notificationID) === String(notificationID),
        )

        if (notification && notification.readAt) {
          notification.readAt = undefined
          state.totalUnreadCount += 1
        }
      },

      [types.deleteNotification] (state, notificationID) {
        const removed = state.notifications.find(n =>
          String(n.notificationID) === String(notificationID),
        )

        state.notifications = state.notifications.filter(n =>
          String(n.notificationID) !== String(notificationID),
        )

        if (removed && !removed.readAt) {
          state.totalUnreadCount = Math.max(0, state.totalUnreadCount - 1)
        }
      },

      [types.setMuted] (state, muted) {
        state.muted = muted
        localStorage.setItem('notificationsMuted', muted)
      },

      [types.setTotalUnreadCount] (state, count) {
        state.totalUnreadCount = Math.max(0, Number(count) || 0)
      },

      [types.updateReadNotification] (state, notification) {
        const existingNotification = state.notifications.find(n =>
          String(n.notificationID) === String(notification.notificationID),
        )

        const hasReadAt = Object.prototype.hasOwnProperty.call(notification, 'readAt')
        const willBeUnread = hasReadAt ? !notification.readAt : false

        // Determine prior state: prefer local truth, otherwise assume the event
        // reflects a state change (i.e. the prior state was the opposite).
        const wasUnread = existingNotification
          ? !existingNotification.readAt
          : !willBeUnread

        if (existingNotification) {
          if (hasReadAt) {
            existingNotification.readAt = notification.readAt ? new Date(notification.readAt) : undefined
          } else {
            existingNotification.readAt = new Date()
          }
        }

        if (wasUnread && !willBeUnread) {
          state.totalUnreadCount = Math.max(0, state.totalUnreadCount - 1)
        } else if (!wasUnread && willBeUnread) {
          state.totalUnreadCount += 1
        }
      },

      [types.updateAllReadNotifications] (state, notifications) {
        const now = new Date()

        // If notifications array is provided, update only those specific ones
        if (Array.isArray(notifications) && notifications.length > 0) {
          const notificationIds = notifications.map(n => String(n.notificationID))

          state.notifications.forEach(notification => {
            if (notificationIds.includes(String(notification.notificationID))) {
              notification.readAt = now
            }
          })

          state.totalUnreadCount = Math.max(0, state.totalUnreadCount - notifications.length)
        } else {
          // Otherwise mark all as read
          state.notifications.forEach(notification => {
            if (!notification.readAt) {
              notification.readAt = now
            }
          })
          state.totalUnreadCount = 0
        }
      },

      [types.updateAllUnreadNotifications] (state, notifications) {
        // If notifications array is provided, update only those specific ones
        if (Array.isArray(notifications) && notifications.length > 0) {
          const notificationIds = notifications.map(n => String(n.notificationID))

          state.notifications.forEach(notification => {
            if (notificationIds.includes(String(notification.notificationID))) {
              notification.readAt = undefined
            }
          })

          state.totalUnreadCount += notifications.length
        } else {
          // Otherwise mark all as unread
          state.notifications.forEach(notification => {
            notification.readAt = undefined
          })
        }
      },
    },
  }
}
