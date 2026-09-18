import { system } from '@planetcrust/human-js'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationsStore } from '../stores/useNotificationsStore'
import { useRightSidebarStore } from '../stores/useRightSidebarStore'
import { useOpenNotification } from './useOpenNotification'

export type SystemNotificationPermission = NotificationPermission | 'unsupported'

// The Human tab that has focus, if any, shared across tabs so an unfocused tab
// can tell that the user is looking at another one.
export const FOCUSED_TAB_KEY = 'notificationsFocusedTab'

const tabID = Math.random().toString(36).slice(2)

function readPermission(): SystemNotificationPermission {
  if (!('Notification' in window) || !window.isSecureContext) {
    return 'unsupported'
  }
  return Notification.permission
}

// The browser grants per origin, so every caller in the tab shares one state.
const permission = ref<SystemNotificationPermission>(readPermission())
let tracking = false

function focusedTab(): string | null {
  try {
    return localStorage.getItem(FOCUSED_TAB_KEY)
  } catch {
    return null
  }
}

function trackFocus() {
  if (tracking) return
  tracking = true

  const claim = () => {
    try {
      localStorage.setItem(FOCUSED_TAB_KEY, tabID)
    } catch {}
  }
  const release = () => {
    if (focusedTab() !== tabID) return
    try {
      localStorage.removeItem(FOCUSED_TAB_KEY)
    } catch {}
  }

  window.addEventListener('focus', claim)
  window.addEventListener('blur', release)
  window.addEventListener('pagehide', release)
  if (document.hasFocus()) claim()

  navigator.permissions
    ?.query({ name: 'notifications' as PermissionName })
    .then(status => {
      status.onchange = () => {
        permission.value = readPermission()
      }
    })
    .catch(() => {})
}

// A key naming this tab while it has no focus is left over from a missed blur.
function humanTabFocused(): boolean {
  if (document.hasFocus()) return true
  const focused = focusedTab()
  return !!focused && focused !== tabID
}

// Browser permission for OS notifications; `request` only asks while the
// browser has not been asked yet, and must run inside a click.
export function useSystemNotificationPermission() {
  permission.value = readPermission()

  async function request() {
    permission.value = readPermission()
    if (permission.value !== 'default') return permission.value

    permission.value = await Notification.requestPermission()
    return permission.value
  }

  return { permission, request }
}

// Mirrors an incoming notification as an OS notification while no Human tab
// has focus, unless muted or not permitted.
export function useSystemNotifications() {
  const store = useNotificationsStore()
  const rightSidebar = useRightSidebarStore()
  const openNotification = useOpenNotification()
  const { t } = useI18n()

  trackFocus()

  function notify(raw: Partial<system.Notification>) {
    permission.value = readPermission()
    if (store.muted || permission.value !== 'granted' || humanTabFocused()) return

    const notification = new system.Notification(raw as any)
    const { title = '', description = '' } = (notification.config || {}) as any

    let shown: Notification
    try {
      shown = new Notification(title || description || t('notifications.newNotification'), {
        body: title ? description : '',
        tag: notification.resourceID,
      })
    } catch (err) {
      console.warn('OS notification not shown:', err)
      return
    }

    shown.onclick = () => {
      window.focus()
      shown.close()

      const current = store.notifications.find(
        n => String(n.notificationID) === String(notification.notificationID),
      )
      if (current && !current.readAt) {
        store
          .markAsRead(String(notification.notificationID))
          .catch(err => console.warn('notification not marked as read:', err))
      }

      if (notification.kind === 'record') {
        openNotification(notification)
      } else {
        rightSidebar.open('notifications')
      }
    }
  }

  return { notify }
}
