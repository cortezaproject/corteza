import { describe, it, expect, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import CInputUser from './CInputUser.vue'

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

const user = (userID: string) => ({ userID, name: `User ${userID}` })

// The page the selector is currently showing. Deliberately never contains the
// already-picked user, which is the case pinning exists for.
const PAGE = [user('2'), user('3')]

let reads: string[] = []

const api = {
  userListCancellable: () => ({
    response: async () => ({ set: PAGE, filter: {} }),
    cancel: () => {},
  }),
  userRead: async ({ userID }: { userID: string }) => {
    reads.push(userID)
    return user(userID)
  },
}

type VM = { options: Array<{ userID: string; label: string }>; selectedUsers: string[] }

function mountInput(props: Record<string, unknown>) {
  return mount(CInputUser, {
    props,
    global: { plugins: [createPinia()], provide: { $SystemAPI: api } },
  })
}

describe('CInputUser multiple', () => {
  beforeAll(() => {
    reads = []
  })

  it('emits the picked users as an array of IDs', async () => {
    const w = mountInput({ multiple: true, modelValue: [] })
    await flushPromises()

    const multi = w.findComponent({ name: 'MultiSelect' })
    expect(multi.exists()).toBe(true)

    multi.vm.$emit('update:modelValue', ['2', '3'])
    await flushPromises()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual([['2', '3']])
  })

  it('keeps a picked user in the options when the page it came from is gone', async () => {
    const w = mountInput({ multiple: true, modelValue: ['1'] })
    await flushPromises()

    const vm = w.vm as unknown as VM
    expect(vm.selectedUsers).toEqual(['1'])

    // Opening the dropdown re-fetches, and the page it gets back does not
    // contain '1'. Each fetch replaces the options wholesale, so without
    // pinning the chip for '1' loses the name it renders from.
    w.findComponent({ name: 'MultiSelect' }).vm.$emit('show')
    await flushPromises()

    expect(vm.options.map(u => u.userID).sort()).toEqual(['1', '2', '3'])
    expect(vm.options.find(u => u.userID === '1')?.label).toBe('User 1')
  })

  it('accepts a scalar model value and treats it as one selection', async () => {
    const w = mountInput({ multiple: true, modelValue: '1' })
    await flushPromises()

    expect((w.vm as unknown as VM).selectedUsers).toEqual(['1'])
  })

  it('still emits a scalar in single-select mode', async () => {
    const w = mountInput({ modelValue: null })
    await flushPromises()

    const select = w.findComponent({ name: 'Select' })
    expect(select.exists()).toBe(true)

    select.vm.$emit('update:modelValue', '2')
    await flushPromises()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual(['2'])
  })
})
