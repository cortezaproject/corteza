import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'

// Resource translations are an instance-level feature: the server collapses
// `resourceTranslations.languages` to a single entry when they are switched
// off, and every translator entry point in compose hangs off this one flag.
// A regression here puts a translate button in front of someone whose instance
// has nowhere to put the translation.

let canManage = true

vi.mock('vue-i18n', () => ({ useI18n: () => ({ locale: { value: 'en' } }) }))
vi.mock('@planetcrust/human-vue', () => ({
  useRBACStore: () => ({ can: () => canManage }),
}))

import { useResourceTranslations } from './useResourceTranslations'

function probe(languages: unknown, manage = true) {
  canManage = manage
  const Probe = defineComponent({
    setup() {
      return useResourceTranslations()
    },
    template: '<div />',
  })
  return mount(Probe, {
    global: {
      provide: { $Settings: { get: () => languages } },
    },
  }).vm as unknown as ReturnType<typeof useResourceTranslations>
}

describe('useResourceTranslations', () => {
  it('hides the translator when the instance has one language', () => {
    const vm = probe(['en'])
    expect(vm.resourceTranslationsEnabled).toBe(false)
    expect(vm.showTranslatorButton).toBe(false)
  })

  it('hides the translator when the setting is unset or empty', () => {
    expect(probe(undefined).showTranslatorButton).toBe(false)
    expect(probe([]).showTranslatorButton).toBe(false)
    expect(probe('en,fr').showTranslatorButton).toBe(false)
  })

  it('shows the translator with more than one language', () => {
    const vm = probe(['en', 'fr', 'sl'])
    expect(vm.resourceTranslationsEnabled).toBe(true)
    expect(vm.showTranslatorButton).toBe(true)
  })

  it('still hides it without the manage permission', () => {
    expect(probe(['en', 'fr'], false).showTranslatorButton).toBe(false)
  })
})
