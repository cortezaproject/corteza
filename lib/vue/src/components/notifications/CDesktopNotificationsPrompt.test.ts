import { describe, it, expect, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import CDesktopNotificationsPrompt from './CDesktopNotificationsPrompt.vue'

const FakeNotification = {
  permission: 'default' as NotificationPermission,
  requestPermission: vi.fn(async () => 'granted' as NotificationPermission),
}

Object.defineProperty(window, 'Notification', { value: FakeNotification, configurable: true })
Object.defineProperty(window, 'isSecureContext', {
  value: true,
  configurable: true,
  writable: true,
})

function mountPrompt() {
  const pinia = createPinia()
  const wrapper = mount(CDesktopNotificationsPrompt, {
    global: { plugins: [pinia], provide: { $SystemAPI: {} } },
  })
  return { wrapper, store: useNotificationsStore(pinia) }
}

const prompt = (wrapper: ReturnType<typeof mount>) =>
  wrapper.find('[data-testid="desktop-notifications-prompt"]')
const enable = (wrapper: ReturnType<typeof mount>) =>
  wrapper.find('[data-testid="desktop-notifications-enable"]')

beforeEach(() => {
  FakeNotification.permission = 'default'
  FakeNotification.requestPermission.mockClear()
  // @ts-expect-error overriding the read-only stub set above
  window.isSecureContext = true
  localStorage.clear()
})

describe('CDesktopNotificationsPrompt', () => {
  it('offers to enable desktop notifications while the browser has not been asked', async () => {
    const { wrapper } = mountPrompt()

    expect(prompt(wrapper).text()).toContain('notifications.desktopPrompt')
    await enable(wrapper).trigger('click')
    await flushPromises()

    expect(FakeNotification.requestPermission).toHaveBeenCalledTimes(1)
    expect(prompt(wrapper).exists()).toBe(false)
  })

  it('says how to unblock once the browser has blocked them', () => {
    FakeNotification.permission = 'denied'

    const { wrapper } = mountPrompt()

    expect(prompt(wrapper).text()).toContain('notifications.desktopBlocked')
    expect(enable(wrapper).exists()).toBe(false)
  })

  it('stays away once they are allowed', () => {
    FakeNotification.permission = 'granted'

    expect(prompt(mountPrompt().wrapper).exists()).toBe(false)
  })

  it('stays away while notifications are muted', async () => {
    const { wrapper, store } = mountPrompt()
    store.toggleMuted()
    await wrapper.vm.$nextTick()

    expect(prompt(wrapper).exists()).toBe(false)
  })

  it('stays away where the browser cannot show them', () => {
    // @ts-expect-error overriding the read-only stub set above
    window.isSecureContext = false

    expect(prompt(mountPrompt().wrapper).exists()).toBe(false)
  })

  it('stays dismissed in this browser', async () => {
    const { wrapper } = mountPrompt()

    await wrapper.find('[data-testid="desktop-notifications-dismiss"]').trigger('click')

    expect(prompt(wrapper).exists()).toBe(false)
    expect(prompt(mountPrompt().wrapper).exists()).toBe(false)
  })
})
