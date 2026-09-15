import { describe, it, expect, beforeAll } from 'vitest'
import { h } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import CInputUser from './CInputUser.vue'

beforeAll(() => {
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

const api = {
  userListCancellable: () => ({
    response: async () => ({ set: [{ userID: '2', name: 'User 2' }], filter: {} }),
    cancel: () => {},
  }),
  userRead: async ({ userID }: { userID: string }) => ({ userID, name: `User ${userID}` }),
}

function mountInput(multiple: boolean, slots = {}) {
  return mount(CInputUser, {
    props: { multiple },
    slots,
    global: { plugins: [createPinia()], provide: { $SystemAPI: api } },
  })
}

describe('CInputUser option slot', () => {
  for (const [mode, name] of [
    [false, 'Select'],
    [true, 'MultiSelect'],
  ] as const) {
    it(`hands the option slot to ${name}`, async () => {
      const w = mountInput(mode, {
        option: ({ option }: { option: { name: string } }) => h('i', option.name),
      })
      await flushPromises()

      const slot = w.findComponent({ name }).vm.$slots.option
      expect(slot).toBeTypeOf('function')
      expect(slot!({ option: { name: 'User 2' } })[0].children).toEqual([
        expect.objectContaining({ children: 'User 2' }),
      ])
    })

    it(`leaves ${name} its own option rendering without one`, async () => {
      const w = mountInput(mode)
      await flushPromises()

      expect(w.findComponent({ name }).vm.$slots.option).toBeUndefined()
    })
  }
})
