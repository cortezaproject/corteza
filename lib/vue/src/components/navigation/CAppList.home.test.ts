import { describe, it, expect, vi, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import PrimeVue from 'primevue/config'
import CAppList from './CAppList.vue'
import { useRBACStore } from '../../composables/useRBAC'
import { useApplicationsStore } from '../../stores/useApplicationsStore'

const toastAdd = vi.fn()
vi.mock('primevue/usetoast', () => ({ useToast: () => ({ add: toastAdd }) }))

// Each tile's menu picks the user's own home application, and is offered only
// to a user holding `application.flag.self`.

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
  window.matchMedia = (q: string) =>
    ({
      matches: false,
      media: q,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
    }) as unknown as MediaQueryList
})

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/:p(.*)*', component: { template: '<div />' } }],
})

const APPS = [
  {
    applicationID: '1',
    name: 'Projects',
    enabled: true,
    unify: { listed: true, url: 'project/', home: true },
  },
  {
    applicationID: '2',
    name: 'TAQ',
    enabled: true,
    meta: { description: 'Workflows and triggers' },
    unify: { listed: true, url: 'taq/' },
  },
]

async function mountList({ canPick, ownHome = '' }: { canPick: boolean; ownHome?: string }) {
  toastAdd.mockClear()
  const pinia = createPinia()
  setActivePinia(pinia)

  const user = { userID: '9', meta: { homeApplicationID: ownHome } }
  const $SystemAPI = {
    baseURL: '',
    applicationList: vi.fn().mockResolvedValue({ set: APPS }),
    permissionsEffective: vi
      .fn()
      .mockResolvedValue(
        canPick
          ? [{ resource: 'corteza::system/', operation: 'application.flag.self', allow: true }]
          : [],
      ),
    userUpdate: vi.fn().mockResolvedValue({}),
  }

  const wrapper = mount(CAppList, {
    attachTo: document.body,
    global: {
      plugins: [pinia, router, PrimeVue],
      provide: { $SystemAPI, $Auth: { user } },
    },
  })

  await useRBACStore().load()
  await useApplicationsStore().fetchApplications()
  await flushPromises()

  return { wrapper, $SystemAPI, user }
}

async function pick(wrapper: Awaited<ReturnType<typeof mountList>>['wrapper'], tile: number) {
  await wrapper.findAll('[data-test-id="app-tile-menu"]')[tile].trigger('click')
  await flushPromises()
  const label = document.body.querySelector('[data-pc-section="itemlabel"]')?.textContent?.trim()
  document.body.querySelector<HTMLElement>('[data-pc-section="itemcontent"]')?.click()
  await flushPromises()
  return label
}

describe('CAppList home menu', () => {
  it('shows an application description under its name', async () => {
    const { wrapper } = await mountList({ canPick: false })
    const described = wrapper.findAll('[data-test-id="app-tile-description"]')
    expect(described.map(d => d.text())).toEqual(['Workflows and triggers'])
    wrapper.unmount()
  })

  it('is not offered without application.flag.self', async () => {
    const { wrapper } = await mountList({ canPick: false })
    expect(wrapper.findAll('a')).toHaveLength(2)
    expect(wrapper.findAll('[data-test-id="app-tile-menu"]')).toHaveLength(0)
    wrapper.unmount()
  })

  it('makes an application the user own home', async () => {
    const { wrapper, $SystemAPI, user } = await mountList({ canPick: true })
    expect(wrapper.findAll('[data-test-id="app-tile-menu"]')).toHaveLength(2)

    // The instance-wide home carries the badge until the user picks their own.
    const badged = () =>
      wrapper
        .findAll('[data-drag-item]')
        .map(t => t.find('[data-test-id="app-tile-home"]').exists())
    expect(badged()).toEqual([true, false])

    expect(await pick(wrapper, 1)).toBe('Set as Home')
    expect($SystemAPI.userUpdate).toHaveBeenCalledWith(user)
    expect(user.meta.homeApplicationID).toBe('2')
    expect(useApplicationsStore().ownHomeID).toBe('2')
    expect(badged()).toEqual([false, true])
    expect(toastAdd).toHaveBeenCalledWith(
      expect.objectContaining({
        severity: 'success',
        summary: 'Home page set',
        detail: 'TAQ now opens when you sign in or select Home.',
      }),
    )
    wrapper.unmount()
  })

  it('clears the pick on the application that holds it', async () => {
    const { wrapper, user } = await mountList({ canPick: true, ownHome: '2' })

    expect(await pick(wrapper, 1)).toBe('Remove as Home')
    expect(user.meta.homeApplicationID).toBe('0')
    expect(useApplicationsStore().ownHomeID).toBe('')
    expect(toastAdd).toHaveBeenCalledWith(
      expect.objectContaining({
        severity: 'success',
        summary: 'Home page removed',
        detail: 'TAQ no longer opens when you sign in or select Home.',
      }),
    )
    wrapper.unmount()
  })
})
