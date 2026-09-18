import { system } from '@planetcrust/human-js'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationsStore } from '../stores/useNotificationsStore'
import { useRightSidebarStore } from '../stores/useRightSidebarStore'
import { useOpenNotification } from './useOpenNotification'

export type SystemNotificationPermission = NotificationPermission | 'unsupported'

// The Human tab that has focus, if any, shared across tabs so an unfocused tab
// can tell that the user is looking at another one. The focused tab renews its
// claim; one that stops being renewed belongs to a tab that is gone.
export const FOCUSED_TAB_KEY = 'notificationsFocusedTab'
export const FOCUS_CLAIM_TTL_MS = 15_000
const FOCUS_CLAIM_RENEW_MS = 5_000

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
    const { id, at } = JSON.parse(localStorage.getItem(FOCUSED_TAB_KEY) || 'null') || {}
    return typeof id === 'string' && Date.now() - at < FOCUS_CLAIM_TTL_MS ? id : null
  } catch {
    return null
  }
}

function trackFocus() {
  if (tracking) return
  tracking = true

  let renewal: ReturnType<typeof setInterval> | undefined

  const write = () => {
    try {
      localStorage.setItem(FOCUSED_TAB_KEY, JSON.stringify({ id: tabID, at: Date.now() }))
    } catch {}
  }
  const claim = () => {
    write()
    clearInterval(renewal)
    renewal = setInterval(write, FOCUS_CLAIM_RENEW_MS)
  }
  const release = () => {
    clearInterval(renewal)
    renewal = undefined
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
