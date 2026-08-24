import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (k: string) => k }) }))

import CFieldUserEditor from './CFieldUserEditor.vue'

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

const api = {
  userListCancellable: () => ({
    response: async () => ({ set: [{ userID: '1', name: 'One' }], filter: {} }),
    cancel: () => {},
  }),
  userRead: async ({ userID }: { userID: string }) => ({ userID, name: `User ${userID}` }),
}

const field = (over: Record<string, unknown> = {}) => ({
  name: 'u',
  kind: 'User',
  isMulti: false,
  options: {},
  ...over,
})

function mountEditor(props: Record<string, unknown>, auth: unknown = null) {
  return mount(CFieldUserEditor, {
    props,
    global: { plugins: [createPinia()], provide: { $SystemAPI: api, $Auth: auth } },
  })
}

describe('CFieldUserEditor', () => {
  it('renders the combined picker only for a multi field asking for it', async () => {
    const w = mountEditor({
      field: field({ isMulti: true, options: { selectType: 'multiple' } }),
      modelValue: [],
    })
    await flushPromises()

    expect(w.findComponent({ name: 'MultiSelect' }).exists()).toBe(true)
    expect(w.findComponent({ name: 'Select' }).exists()).toBe(false)
  })

  it('edits one slot at a time otherwise', async () => {
    const w = mountEditor({ field: field({ isMulti: true, options: {} }), modelValue: '' })
    await flushPromises()

    expect(w.findComponent({ name: 'MultiSelect' }).exists()).toBe(false)
    expect(w.findComponent({ name: 'Select' }).exists()).toBe(true)
  })

  it('passes the field roles through as the picker scope', async () => {
    const w = mountEditor({
      field: field({ options: { roles: ['7'] } }),
      modelValue: '',
    })
    await flushPromises()

    expect(w.findComponent({ name: 'CInputUser' }).props('roleID')).toEqual(['7'])
  })

  it('drops repeats when the field wants unique values', async () => {
    const w = mountEditor({
      field: field({
        isMulti: true,
        options: { selectType: 'multiple', isUniqueMultiValue: true },
      }),
      modelValue: [],
    })
    await flushPromises()

    w.findComponent({ name: 'CInputUser' }).vm.$emit('update:modelValue', ['1', '2', '1'])
    await flushPromises()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual([['1', '2']])
  })

  it('keeps repeats when the field allows them', async () => {
    const w = mountEditor({
      field: field({ isMulti: true, options: { selectType: 'multiple' } }),
      modelValue: [],
    })
    await flushPromises()

    w.findComponent({ name: 'CInputUser' }).vm.$emit('update:modelValue', ['1', '2', '1'])
    await flushPromises()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual([['1', '2', '1']])
  })

  it('presets the authenticated user on an empty field that asks for it', async () => {
    const w = mountEditor(
      { field: field({ options: { presetWithAuthenticated: true } }), modelValue: '' },
      { user: { userID: '9' } },
    )
    await flushPromises()

    expect(w.emitted('update:modelValue')?.[0]).toEqual(['9'])
  })

  it('leaves a field that already has a value alone', async () => {
    const w = mountEditor(
      { field: field({ options: { presetWithAuthenticated: true } }), modelValue: '4' },
      { user: { userID: '9' } },
    )
    await flushPromises()

    expect(w.emitted('update:modelValue')).toBeUndefined()
  })
})
