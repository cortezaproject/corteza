import { onBeforeUnmount, ref } from 'vue'

const STORAGE_KEY_MENU = 'home-column-menu-width'
const STORAGE_KEY_NOTIFICATIONS = 'home-column-notifications-width'

const DEFAULT_MENU_WIDTH = 260
const DEFAULT_NOTIFICATIONS_WIDTH = 360
const MIN_MENU_WIDTH = 200
const MAX_MENU_WIDTH = 480
const MIN_NOTIFICATIONS_WIDTH = 280
const MAX_NOTIFICATIONS_WIDTH = 700

function loadWidth(key, defaultValue) {
  try {
    const stored = localStorage.getItem(key)
    if (stored) {
      const parsed = parseInt(stored, 10)
      if (!isNaN(parsed)) return parsed
    }
  } catch {
    // localStorage unavailable
  }
  return defaultValue
}

function saveWidth(key, value) {
  try {
    localStorage.setItem(key, String(value))
  } catch {
    // localStorage unavailable
  }
}

export function useColumnResize() {
  const menuWidth = ref(loadWidth(STORAGE_KEY_MENU, DEFAULT_MENU_WIDTH))
  const notificationsWidth = ref(loadWidth(STORAGE_KEY_NOTIFICATIONS, DEFAULT_NOTIFICATIONS_WIDTH))

  let activeHandle = null

  function onMouseMove(e) {
    if (activeHandle === 'menu') {
      const newWidth = Math.max(MIN_MENU_WIDTH, Math.min(MAX_MENU_WIDTH, e.clientX))
      menuWidth.value = newWidth
    } else if (activeHandle === 'notifications') {
      const newWidth = Math.max(MIN_NOTIFICATIONS_WIDTH, Math.min(MAX_NOTIFICATIONS_WIDTH, window.innerWidth - e.clientX))
      notificationsWidth.value = newWidth
    }
  }

  function onMouseUp() {
    if (activeHandle === 'menu') {
      saveWidth(STORAGE_KEY_MENU, menuWidth.value)
    } else if (activeHandle === 'notifications') {
      saveWidth(STORAGE_KEY_NOTIFICATIONS, notificationsWidth.value)
    }
    activeHandle = null
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  function startMenuResize() {
    activeHandle = 'menu'
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    document.body.style.cursor = 'ew-resize'
    document.body.style.userSelect = 'none'
  }

  function startNotificationsResize() {
    activeHandle = 'notifications'
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    document.body.style.cursor = 'ew-resize'
    document.body.style.userSelect = 'none'
  }

  onBeforeUnmount(() => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  })

  return {
    menuWidth,
    notificationsWidth,
    startMenuResize,
    startNotificationsResize,
  }
}
