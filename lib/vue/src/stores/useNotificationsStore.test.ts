import { describe, it, expect, beforeEach } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { useNotificationsStore } from './useNotificationsStore'

function mountStore() {
  const pinia = createPinia()
  mount(
    defineComponent({
      setup() {
        useNotificationsStore()
        return () => null
      },
    }),
    { global: { plugins: [pinia], provide: { $SystemAPI: {} } } },
  )
  return useNotificationsStore(pinia)
}

beforeEach(() => {
  localStorage.clear()
})

describe('useNotificationsStore badgeCount', () => {
  it('counts the unread notifications, and none while muted', () => {
    const store = mountStore()
    store.handleRealtime({ '@type': 'notification', '@value': { notificationID: '1' } })
    store.handleRealtime({ '@type': 'notification', '@value': { notificationID: '2' } })

    expect(store.badgeCount).toBe(2)

    store.toggleMuted()

    expect(store.badgeCount).toBe(0)
    expect(store.unreadCount).toBe(2)
  })
})
