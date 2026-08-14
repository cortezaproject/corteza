import { describe, it, expect, vi, beforeEach } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'
import { useHistoryBack } from './useHistoryBack'

// Back has to land somewhere. The screen is the first history entry whenever it
// was deep-linked, bookmarked, opened in a fresh tab, or reached by a redirect
// that replaced the entry it came from — and there router.back() either does
// nothing or leaves the app.

const history = { state: {} as { back?: string | null } }
const router = {
  back: vi.fn(),
  push: vi.fn(),
  options: { history },
}

vi.mock('vue-router', () => ({
  useRouter: () => router,
}))

const FALLBACK = { name: 'admin.modules', params: { slug: 'ns' } }

function callGoBack() {
  let goBack: (to: unknown) => void
  mount(
    defineComponent({
      setup() {
        goBack = useHistoryBack()
        return () => null
      },
    }),
  )
  goBack!(FALLBACK)
}

beforeEach(() => {
  router.back.mockClear()
  router.push.mockClear()
})

describe('useHistoryBack', () => {
  it('returns to the previous screen when there is one', () => {
    history.state = { back: '/compose/namespace/ns/admin/modules' }

    callGoBack()

    expect(router.back).toHaveBeenCalled()
    expect(router.push).not.toHaveBeenCalled()
  })

  it('pushes the fallback when the screen is the first history entry', () => {
    history.state = { back: null }

    callGoBack()

    expect(router.push).toHaveBeenCalledWith(FALLBACK)
    expect(router.back).not.toHaveBeenCalled()
  })

  it('pushes the fallback when there is no history state at all', () => {
    history.state = {}

    callGoBack()

    expect(router.push).toHaveBeenCalledWith(FALLBACK)
    expect(router.back).not.toHaveBeenCalled()
  })
})
