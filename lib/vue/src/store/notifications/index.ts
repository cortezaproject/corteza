import { apiClients, system } from '@cortezaproject/corteza-js'
import { StoreOptions } from 'vuex'

const types = {
  fetchNotifications: 'fetchNotifications',
  setNotifications: 'setNotifications',
  appendNotifications: 'appendNotifications',
  setVisible: 'setVisible',
  markAsRead: 'markAsRead',
  markAllAsRead: 'markAllAsRead',
  addNotification: 'addNotification',
  deleteNotification: 'deleteNotification',
  setPageCursor: 'setPageCursor',
  setMuted: 'setMuted',
  updateReadNotification: 'updateReadNotification',
  updateAllReadNotifications: 'updateAllReadNotifications',
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
    },

    getters: {
      notifications: (state) => state.notifications,
      visible: (state) => state.visible,
      pageCursor: (state) => state.pageCursor,
      hasMorePages: (state) => !!state.pageCursor,
      muted: (state) => state.muted,

      hasUnread: (state) => state.notifications.some(notification => !notification.readAt),
      unreadCount: (state) => state.notifications.filter(notification => !notification.readAt).length,
    },

    actions: {
      toggleVisibility ({ commit, state }) {
        commit(types.setVisible, !state.visible)
      },

      fetchNotifications ({ commit, state }, { unreadOnly = true } = {}) {
        return api.notificationList({
          limit: 25,
          sort: state.pageCursor ? '' : 'createdAt DESC, readAt DESC',
          read: unreadOnly ? 0 : 1,
          pageCursor: state.pageCursor,
        })
          .then((response: KV) => {
            const set = (response.set || []) as Array<system.Notification>
            const filter = (response.filter || {}) as {
              limit: number;
              sort: string;
              read: number;
              nextPage?: string;
            }

            if (state.pageCursor) {
              commit(types.appendNotifications, set)
            } else {
              commit(types.setNotifications, set)
            }

            commit(types.setPageCursor, filter.nextPage)
          })
      },

      markAsRead ({ commit }, notificationID) {
        return api.notificationMarkAsRead({ notificationID })
          .then(() => {
            commit(types.markAsRead, notificationID)
          })
      },

      markAllAsRead ({ commit, state }) {
        if (!state.notifications.length) {
          return Promise.resolve()
        }

        return api.notificationMarkAllAsRead()
          .then(() => {
            commit(types.markAllAsRead)
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

      updateAllReadNotifications ({ commit }, notifications) {
        commit(types.updateAllReadNotifications, notifications)
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

        if (notification) {
          notification.readAt = new Date()
        }
      },

      [types.markAllAsRead] (state) {
        const now = new Date()
        state.notifications.forEach(notification => {
          if (!notification.readAt) {
            notification.readAt = now
          }
        })
      },

      [types.addNotification] (state, notification) {
        // Add to the beginning of the array
        state.notifications.unshift(new system.Notification(notification))
      },

      [types.deleteNotification] (state, notificationID) {
        state.notifications = state.notifications.filter(n =>
          String(n.notificationID) !== String(notificationID),
        )
      },

      [types.setMuted] (state, muted) {
        state.muted = muted
        localStorage.setItem('notificationsMuted', muted)
      },

      [types.updateReadNotification] (state, notification) {
        const existingNotification = state.notifications.find(n =>
          String(n.notificationID) === String(notification.notificationID),
        )

        if (existingNotification) {
          existingNotification.readAt = notification.readAt || new Date()
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
        } else {
          // Otherwise mark all as read
          state.notifications.forEach(notification => {
            if (!notification.readAt) {
              notification.readAt = now
            }
          })
        }
      },
    },
  }
}
