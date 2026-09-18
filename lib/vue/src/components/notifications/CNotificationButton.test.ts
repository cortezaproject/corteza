import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import CNotificationButton from './CNotificationButton.vue'

const ButtonStub = {
  props: { icon: String, severity: String },
  template: '<button :data-icon="icon" :data-severity="severity ?? \'primary\'" />',
}

function mountBell() {
  const pinia = createPinia()
  const wrapper = mount(CNotificationButton, {
    global: {
      plugins: [pinia],
      provide: { $SystemAPI: {} },
      stubs: { Button: ButtonStub, Badge: true },
      directives: {
        tooltip: {
          mounted: (el: HTMLElement, b: any) => (el.dataset.tooltip = b.value),
          updated: (el: HTMLElement, b: any) => (el.dataset.tooltip = b.value),
        },
      },
    },
  })
  return { wrapper, store: useNotificationsStore(pinia) }
}

const bell = (wrapper: ReturnType<typeof mount>) =>
  wrapper.find('[data-testid="notification-bell"]')
const arrive = (store: ReturnType<typeof useNotificationsStore>, id: string) =>
  store.handleRealtime({ '@type': 'notification', '@value': { notificationID: id } })

beforeEach(() => {
  localStorage.clear()
  vi.useFakeTimers()
})

afterEach(() => {
  vi.useRealTimers()
})

describe('CNotificationButton', () => {
  it('is a plain bell while nothing is unread', () => {
    const b = bell(mountBell().wrapper)

    expect(b.attributes('data-icon')).toBe('pi pi-bell')
    expect(b.attributes('data-severity')).toBe('secondary')
    expect(b.classes()).not.toContain('notification-bell-ring')
  })

  it('takes the primary colour while something is unread', async () => {
    const { wrapper, store } = mountBell()
    arrive(store, '1')
    await flushPromises()

    expect(bell(wrapper).attributes('data-severity')).toBe('primary')
  })

  it('rings when a notification arrives, and stops', async () => {
    const { wrapper, store } = mountBell()
    arrive(store, '1')
    await flushPromises()

    expect(bell(wrapper).classes()).toContain('notification-bell-ring')

    vi.advanceTimersByTime(1000)
    await flushPromises()

    expect(bell(wrapper).classes()).not.toContain('notification-bell-ring')
  })

  it('does not ring for notifications that were already there', async () => {
    const { wrapper, store } = mountBell()
    store.updateReadNotification({ notificationID: '1' })
    store.removeNotification({ notificationID: '1' })
    await flushPromises()

    expect(bell(wrapper).classes()).not.toContain('notification-bell-ring')
  })

  it('shows a slashed, uncoloured, silent bell while muted', async () => {
    const { wrapper, store } = mountBell()
    store.toggleMuted()
    arrive(store, '1')
    await flushPromises()

    const b = bell(wrapper)
    expect(b.attributes('data-icon')).toBe('pi pi-bell-slash')
    expect(b.attributes('data-severity')).toBe('secondary')
    expect(b.attributes('data-tooltip')).toBe('notifications.titleMuted')
    expect(b.classes()).not.toContain('notification-bell-ring')
  })
})
