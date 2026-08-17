import { describe, it, expect, vi, afterEach } from 'vitest'
import { defineComponent, ref, type Ref } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createI18n } from 'vue-i18n'
import { useDraftGuard } from './useDraftGuard'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: { en: { 'unsaved.confirm': 'Leave without saving?' } },
})

type Draft = { name: string; nested?: { n: number } } | null

// Mounted as a route-view child, the only place onBeforeRouteLeave takes effect.
async function mountGuarded(initial: Draft = { name: 'a', nested: { n: 1 } }) {
  const draft = ref<Draft>(initial)
  const busy = ref(false)
  const extra = ref<string[]>([])
  let guard: ReturnType<typeof useDraftGuard<Draft>> | undefined

  const Guarded = defineComponent({
    setup() {
      guard = useDraftGuard<Draft>({
        draft,
        busy,
        extra: () => [...extra.value],
        messageKey: 'unsaved.confirm',
        tabClose: false,
      })
      return {}
    },
    template: '<div/>',
  })

  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: Guarded, name: 'guarded' },
      { path: '/other', component: { template: '<div/>' }, name: 'other' },
    ],
  })

  const wrapper = mount(defineComponent({ template: '<router-view/>' }), {
    global: { plugins: [router, i18n] },
    attachTo: document.body,
  })

  await router.push('/')
  await flushPromises()

  return {
    router,
    draft: draft as Ref<Draft>,
    busy,
    extra,
    wrapper,
    get guard() {
      return guard!
    },
  }
}

describe('useDraftGuard', () => {
  afterEach(() => {
    vi.restoreAllMocks()
    document.body.innerHTML = ''
  })

  it('stays quiet before a baseline is captured', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, draft } = await mountGuarded()

    // Data still arriving: the draft changes, but nothing has been shown as
    // saved yet, so there is nothing to lose.
    draft.value = { name: 'loaded from the server' }
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/other')
  })

  it('stays quiet when the draft is untouched since capture', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard } = await mountGuarded()

    guard.capture()
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/other')
  })

  it('warns once the draft differs from the baseline', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.name = 'edited'
    await router.push('/other')

    expect(window.confirm).toHaveBeenCalledWith('Leave without saving?')
    expect(router.currentRoute.value.path).toBe('/')
  })

  it('compares deeply, not by reference', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.nested!.n = 2
    await router.push('/other')

    expect(router.currentRoute.value.path).toBe('/')
  })

  it('holds a snapshot, so mutating the draft cannot rewrite the baseline', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.nested!.n = 99
    await router.push('/other')

    // A baseline holding a live reference would compare equal to itself here.
    expect(confirmSpy).toHaveBeenCalled()
  })

  it('warns when state held outside the draft changes', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, extra } = await mountGuarded()

    guard.capture()
    extra.value = ['member-1']
    await router.push('/other')

    expect(router.currentRoute.value.path).toBe('/')
  })

  it('stays quiet while busy, however dirty the draft is', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft, busy } = await mountGuarded()

    guard.capture()
    draft.value!.name = 'edited'
    busy.value = true
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/other')
  })

  it('is clean again after re-capturing a saved draft', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.name = 'edited'
    guard.capture()
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/other')
  })

  it('goes quiet again after reset', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.name = 'edited'
    guard.reset()
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
  })

  it('lets markSaved through one navigation without asking', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value!.name = 'edited'
    guard.markSaved()
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/other')
  })

  it('treats a null draft as nothing to lose', async () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)
    const { router, guard, draft } = await mountGuarded()

    guard.capture()
    draft.value = null
    await router.push('/other')

    expect(confirmSpy).not.toHaveBeenCalled()
  })
})
