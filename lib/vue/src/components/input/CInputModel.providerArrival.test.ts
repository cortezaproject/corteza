import { describe, it, expect, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CInputModel from './CInputModel.vue'

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
  llmProviderModelsCancellable: () => ({
    response: async () => ({ models: ['mistral-large-latest', 'mistral-small-latest'] }),
    cancel: () => {},
  }),
}

function mountModel(props: Record<string, unknown>) {
  return mount(CInputModel, {
    props,
    global: { provide: { $SystemAPI: api } },
  })
}

describe('CInputModel provider changes', () => {
  // An agent's provider is resolved after its model is already bound — from
  // the loaded resource, or stood in for when the instance has one provider.
  // Treating that arrival as a change cleared a model nobody touched, and the
  // next save persisted the blank.
  it('keeps the model when the provider arrives', async () => {
    const w = mountModel({ modelValue: 'mistral-large-latest', llmProviderID: '' })
    await w.setProps({ llmProviderID: '4965' })
    await flushPromises()

    expect(w.emitted('update:modelValue')).toBeUndefined()
  })

  it('keeps the model when the provider arrives over the unset sentinel', async () => {
    const w = mountModel({ modelValue: 'mistral-large-latest', llmProviderID: '0' })
    await w.setProps({ llmProviderID: '4965' })
    await flushPromises()

    expect(w.emitted('update:modelValue')).toBeUndefined()
  })

  // Switching between two real providers is a real change: the models on offer
  // differ, so the old choice cannot stand.
  it('clears the model when one provider replaces another', async () => {
    const w = mountModel({ modelValue: 'mistral-large-latest', llmProviderID: '4965' })
    await w.setProps({ llmProviderID: '8877' })
    await flushPromises()

    expect(w.emitted('update:modelValue')).toEqual([[null]])
  })
})
