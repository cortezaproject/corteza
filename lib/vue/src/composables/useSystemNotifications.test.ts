import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { useNotificationsStore } from '../stores/useNotificationsStore'
import { useRightSidebarStore } from '../stores/useRightSidebarStore'
import {
  FOCUS_CLAIM_TTL_MS,
  FOCUSED_TAB_KEY,
  useSystemNotificationPermission,
  useSystemNotifications,
} from './useSystemNotifications'

const openNotification = vi.hoisted(() => vi.fn())

vi.mock('./useOpenNotification', () => ({
  useOpenNotification: () => openNotification,
}))

class FakeNotification {
  static permission: NotificationPermission = 'granted'
  static requestPermission = vi.fn(async () => FakeNotification.permission)
  static shown: FakeNotification[] = []

  onclick: (() => void) | null = null
  close = vi.fn()

  constructor(
    public title: string,
    public options: NotificationOptions = {},
  ) {
    FakeNotification.shown.push(this)
  }
}

Object.defineProperty(window, 'Notification', { value: FakeNotification, configurable: true })
Object.defineProperty(window, 'isSecureContext', { value: true, configurable: true })

const RECORD = {
  notificationID: '11',
  kind: 'record',
  config: {
    record: {
      title: 'Invoice overdue',
      description: 'Invoice 42 is 30 days late',
      namespaceID: '1',
      moduleID: '2',
      recordID: '3',
    },
  },
}

const SIMPLE = {
  notificationID: '12',
  kind: 'simple',
  config: { simple: { title: 'Export ready', description: '' } },
}

let hasFocus: ReturnType<typeof vi.spyOn>

const claim = (id: string, at: number) =>
  localStorage.setItem(FOCUSED_TAB_KEY, JSON.stringify({ id, at }))

function setup(options = {}) {
  const pinia = createPinia()
  const api = { notificationMarkAsRead: vi.fn(async () => ({})) }
  let notify!: (raw: any) => Promise<void>

  mount(
    defineComponent({
      setup() {
        ;({ notify } = useSystemNotifications(options))
        return () => null
      },
    }),
    { global: { plugins: [pinia], provide: { $SystemAPI: api } } },
  )

  const store = useNotificationsStore(pinia)
  // The shell hands every realtime message to the store before the OS sees it.
  const receive = (raw: any) => {
    store.handleRealtime({ '@type': 'notification', '@value': raw })
    return notify(raw)
  }

  return { receive, store, api, rightSidebar: useRightSidebarStore(pinia) }
}

beforeEach(() => {
  FakeNotification.permission = 'granted'
  FakeNotification.requestPermission.mockClear()
  FakeNotification.shown = []
  openNotification.mockClear()
  localStorage.clear()
  hasFocus = vi.spyOn(document, 'hasFocus').mockReturnValue(false)
})

afterEach(() => {
  hasFocus.mockRestore()
})

describe('useSystemNotifications', () => {
  it('shows an OS notification with the title, the description and a per-notification tag', async () => {
    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(1)
    const [shown] = FakeNotification.shown
    expect(shown.title).toBe('Invoice overdue')
    expect(shown.options.body).toBe('Invoice 42 is 30 days late')
    expect(shown.options.tag).toBe('system:notification:11')
  })

  it('shows nothing while notifications are muted', async () => {
    const { receive, store } = setup()
    store.toggleMuted()

    await receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(0)
    store.toggleMuted()
  })

  it.each(['default', 'denied'] as const)('shows nothing while permission is %s', async perm => {
    FakeNotification.permission = perm

    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(0)
  })

  it('shows nothing while this tab has focus', async () => {
    hasFocus.mockReturnValue(true)

    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(0)
  })

  it('shows nothing while another Human tab has focus', async () => {
    claim('another-tab', Date.now())

    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(0)
  })

  it('ignores the claim of a tab that stopped renewing it', async () => {
    claim('another-tab', Date.now() - FOCUS_CLAIM_TTL_MS - 1)

    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(1)
  })

  it('ignores a claim it cannot read', async () => {
    localStorage.setItem(FOCUSED_TAB_KEY, 'another-tab')

    await setup().receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(1)
  })

  it('renews its own claim for as long as it has focus', () => {
    vi.useFakeTimers()
    setup()
    window.dispatchEvent(new Event('focus'))
    const first = JSON.parse(localStorage.getItem(FOCUSED_TAB_KEY)!).at

    vi.advanceTimersByTime(FOCUS_CLAIM_TTL_MS * 2)

    expect(JSON.parse(localStorage.getItem(FOCUSED_TAB_KEY)!).at).toBeGreaterThan(
      first + FOCUS_CLAIM_TTL_MS,
    )
    window.dispatchEvent(new Event('blur'))
    vi.useRealTimers()
  })

  it('is not held back by a focus claim this tab left behind', async () => {
    const { receive } = setup()
    window.dispatchEvent(new Event('focus'))
    expect(localStorage.getItem(FOCUSED_TAB_KEY)).toBeTruthy()

    await receive(RECORD)

    expect(FakeNotification.shown).toHaveLength(1)
    window.dispatchEvent(new Event('blur'))
  })

  it('stops holding other tabs back once this tab loses focus', async () => {
    setup()
    window.dispatchEvent(new Event('focus'))
    window.dispatchEvent(new Event('blur'))

    expect(localStorage.getItem(FOCUSED_TAB_KEY)).toBeNull()
  })

  it('titles an untitled notification with its description, then a generic title', async () => {
    const { receive } = setup()

    await receive({ ...SIMPLE, config: { simple: { title: '', description: 'Only a body' } } })
    await receive({
      ...SIMPLE,
      notificationID: '13',
      config: { simple: { title: '', description: '' } },
    })

    expect(FakeNotification.shown.map(n => [n.title, n.options.body])).toEqual([
      ['Only a body', ''],
      ['notifications.newNotification', ''],
    ])
  })

  it('clicking a record notification focuses the tab, marks it read and opens the record', async () => {
    const focus = vi.spyOn(window, 'focus').mockImplementation(() => {})
    const { receive, api } = setup()
    await receive(RECORD)

    FakeNotification.shown[0].onclick!()

    expect(focus).toHaveBeenCalled()
    expect(FakeNotification.shown[0].close).toHaveBeenCalled()
    expect(api.notificationMarkAsRead).toHaveBeenCalledWith({ notificationID: '11' })
    expect(openNotification).toHaveBeenCalledWith(
      expect.objectContaining({ notificationID: '11', kind: 'record' }),
    )
    focus.mockRestore()
  })

  it('shows the icon it is given, and none when the icon fails', async () => {
    await setup({ icon: async () => 'data:image/png;base64,ICON' }).receive(RECORD)
    await setup({
      icon: async () => {
        throw new Error('no icon')
      },
    }).receive({ ...RECORD, notificationID: '14' })

    expect(FakeNotification.shown.map(n => n.options.icon)).toEqual([
      'data:image/png;base64,ICON',
      undefined,
    ])
  })

  it('clicking a simple notification opens the notifications panel instead', async () => {
    const { receive, rightSidebar, api } = setup()
    await receive(SIMPLE)

    FakeNotification.shown[0].onclick!()

    expect(rightSidebar.isOpen('notifications')).toBe(true)
    expect(openNotification).not.toHaveBeenCalled()
    expect(api.notificationMarkAsRead).toHaveBeenCalledWith({ notificationID: '12' })
  })
})

describe('useSystemNotificationPermission', () => {
  it('asks the browser while it has not been asked, and records the answer', async () => {
    FakeNotification.permission = 'default'
    FakeNotification.requestPermission.mockImplementationOnce(async () => 'granted')
    const { permission, request } = useSystemNotificationPermission()

    await request()

    expect(FakeNotification.requestPermission).toHaveBeenCalledTimes(1)
    expect(permission.value).toBe('granted')
  })

  it.each(['granted', 'denied'] as const)('never asks again once the answer is %s', async perm => {
    FakeNotification.permission = perm
    const { permission, request } = useSystemNotificationPermission()

    await request()

    expect(FakeNotification.requestPermission).not.toHaveBeenCalled()
    expect(permission.value).toBe(perm)
  })
})
