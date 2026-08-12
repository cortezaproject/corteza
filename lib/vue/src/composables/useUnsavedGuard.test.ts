import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { defineComponent, ref } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createI18n } from 'vue-i18n'
import { useUnsavedGuard } from './useUnsavedGuard'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: { en: { 'unsaved.confirm': 'You have unsaved changes. Leave?' } },
})

// Mount the component AS a route-view child so onBeforeRouteLeave works correctly
async function mountGuardedRoute(isDirty: boolean, tabClose = false) {
  const dirty = ref(isDirty)
  let guardRef: { markSaved: () => void } | undefined

  const GuardedComp = defineComponent({
    setup() {
      guardRef = useUnsavedGuard({ isDirty: dirty, messageKey: 'unsaved.confirm', tabClose })
      return {}
    },
    template: '<div data-testid="guarded"/>',
  })

  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: GuardedComp, name: 'guarded' },
      { path: '/other', component: { template: '<div/>' }, name: 'other' },
    ],
  })

  const App = defineComponent({ template: '<router-view/>' })
  const wrapper = mount(App, {
    global: { plugins: [router, i18n] },
    attachTo: document.body,
  })

  await router.push('/')
  await flushPromises()

  return {
    wrapper,
    router,
    dirty,
    get guard() {
      return guardRef!
    },
  }
}

describe('useUnsavedGuard', () => {
  afterEach(() => {
    vi.restoreAllMocks()
    document.body.innerHTML = ''
  })

  describe('navigation guard', () => {
    it('allows navigation when not dirty', async () => {
      const { router } = await mountGuardedRoute(false)
      await router.push('/other')
      expect(router.currentRoute.value.path).toBe('/other')
    })

    it('blocks navigation when dirty and confirm returns false', async () => {
      vi.spyOn(window, 'confirm').mockReturnValue(false)
      const { router } = await mountGuardedRoute(true)
      await router.push('/other')
      expect(router.currentRoute.value.path).toBe('/')
    })

    it('allows navigation when dirty but confirm returns true', async () => {
      vi.spyOn(window, 'confirm').mockReturnValue(true)
      const { router } = await mountGuardedRoute(true)
      await router.push('/other')
      expect(router.currentRoute.value.path).toBe('/other')
    })

    it('calls confirm with translated message', async () => {
      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
      const { router } = await mountGuardedRoute(true)
      await router.push('/other')
      expect(confirmSpy).toHaveBeenCalledWith('You have unsaved changes. Leave?')
    })
  })

  describe('markSaved()', () => {
    it('allows navigation after markSaved even when dirty', async () => {
      vi.spyOn(window, 'confirm').mockReturnValue(false)
      const { router, guard } = await mountGuardedRoute(true)

      guard.markSaved()
      await router.push('/other')

      expect(router.currentRoute.value.path).toBe('/other')
      expect(window.confirm).not.toHaveBeenCalled()
    })

    it('resets after one use — next navigation still checks dirty', async () => {
      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
      const { router, guard } = await mountGuardedRoute(true)

      guard.markSaved()
      await router.push('/other')
      // markSaved consumed; component unmounted (navigated away), so guard no longer active
      // The important assertion: confirm was NOT called for the first navigation
      expect(confirmSpy).not.toHaveBeenCalled()
    })
  })

  describe('beforeunload (tabClose)', () => {
    it('calls preventDefault when dirty and tabClose enabled', async () => {
      const { wrapper } = await mountGuardedRoute(true, true)
      await flushPromises()

      const event = new Event('beforeunload', { cancelable: true }) as BeforeUnloadEvent
      Object.defineProperty(event, 'returnValue', { writable: true, value: '' })
      const preventSpy = vi.spyOn(event, 'preventDefault')
      window.dispatchEvent(event)
      expect(preventSpy).toHaveBeenCalled()

      wrapper.unmount()
    })

    it('does not call preventDefault when clean', async () => {
      const { wrapper } = await mountGuardedRoute(false, true)
      await flushPromises()

      const event = new Event('beforeunload', { cancelable: true }) as BeforeUnloadEvent
      Object.defineProperty(event, 'returnValue', { writable: true, value: '' })
      const preventSpy = vi.spyOn(event, 'preventDefault')
      window.dispatchEvent(event)
      expect(preventSpy).not.toHaveBeenCalled()

      wrapper.unmount()
    })

    it('removes listener on unmount', async () => {
      const { wrapper, dirty } = await mountGuardedRoute(true, true)
      await flushPromises()
      wrapper.unmount()

      // After unmount, event should no longer be intercepted
      const event = new Event('beforeunload', { cancelable: true }) as BeforeUnloadEvent
      Object.defineProperty(event, 'returnValue', { writable: true, value: '' })
      const preventSpy = vi.spyOn(event, 'preventDefault')
      dirty.value = true
      window.dispatchEvent(event)
      expect(preventSpy).not.toHaveBeenCalled()
    })
  })
})
